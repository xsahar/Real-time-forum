package posts

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "Real-time-forum/src/database"
    "Real-time-forum/src/models"
    "Real-time-forum/src/handlers/auth"
)

type CommentResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := auth.ValidateSession(w, r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    if r.Method == http.MethodPost {
        // Add comment
        var comment models.Comment
        if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Validate content
        if comment.Content == "" {
            http.Error(w, "Comment content required", http.StatusBadRequest)
            return
        }

        // Insert comment
        result, err := database.DB.Exec(`
            INSERT INTO comments (post_id, user_id, content, created_at)
            VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
            comment.PostID, userID, comment.Content)
        if err != nil {
            http.Error(w, "Error creating comment", http.StatusInternalServerError)
            return
        }

        commentID, _ := result.LastInsertId()

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(CommentResponse{
            Status:  "success",
            Message: "Comment added successfully",
            Data:    commentID,
        })
        return
    }

    // GET method - fetch comments for a post
    postID := r.URL.Query().Get("post_id")
    if postID == "" {
        http.Error(w, "Post ID required", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(postID)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    rows, err := database.DB.Query(`
        SELECT c.id, c.content, u.nickname as author, c.created_at
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at DESC`, id)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var comments []models.Comment
    for rows.Next() {
        var comment models.Comment
        if err := rows.Scan(&comment.ID, &comment.Content, &comment.Author, &comment.CreatedAt); err != nil {
            continue
        }
        comments = append(comments, comment)
    }

    json.NewEncoder(w).Encode(CommentResponse{
        Status: "success",
        Data:   comments,
    })
}