package crud

import (
    "encoding/json"
    "log"
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
    // Set headers
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Handle preflight
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    // Validate authentication
    userID, ok := auth.ValidateSession(w, r)
    if !ok {
        log.Printf("CommentHandler: Unauthorized access attempt")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    switch r.Method {
    case http.MethodPost:
        handleAddComment(w, r, userID)
    case http.MethodGet:
        handleGetComments(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleAddComment(w http.ResponseWriter, r *http.Request, userID int) {
    log.Printf("HandleAddComment: Processing new comment")
    
    var comment models.Comment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        log.Printf("HandleAddComment: Invalid request body: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if comment.Content == "" {
        log.Printf("HandleAddComment: Empty comment content")
        http.Error(w, "Comment content required", http.StatusBadRequest)
        return
    }

    commentID, err := database.AddComment(comment.PostID, userID, comment.Content)
    if err != nil {
        log.Printf("HandleAddComment: Database error: %v", err)
        http.Error(w, "Error creating comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(CommentResponse{
        Status:  "success",
        Message: "Comment added successfully",
        Data:    commentID,
    })
}

func handleGetComments(w http.ResponseWriter, r *http.Request) {
    log.Printf("HandleGetComments: Fetching comments")
    
    postID := r.URL.Query().Get("post_id")
    if postID == "" {
        log.Printf("HandleGetComments: Missing post_id parameter")
        http.Error(w, "Post ID required", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(postID)
    if err != nil {
        log.Printf("HandleGetComments: Invalid post_id: %v", err)
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    comments, err := database.GetComments(id)
    if err != nil {
        log.Printf("HandleGetComments: Database error: %v", err)
        http.Error(w, "Error fetching comments", http.StatusInternalServerError)
        return
    }

    log.Printf("HandleGetComments: Found %d comments for post %d", len(comments), id)
    json.NewEncoder(w).Encode(CommentResponse{
        Status: "success",
        Data:   comments,
    })
}