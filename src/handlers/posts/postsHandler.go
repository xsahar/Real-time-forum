package posts

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    // "time"

    "Real-time-forum/src/database"
    "Real-time-forum/src/handlers/auth"
    "Real-time-forum/src/models"
)

type PostResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}



func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Only allow GET method
    if r.Method != http.MethodGet {
        json.NewEncoder(w).Encode(PostResponse{
            Status: "error",
            Error:  "Method not allowed",
        })
        return
    }

    // Query posts with details
    query := `
        SELECT 
            p.id, p.user_id, p.title, p.content, p.created_at,
            u.nickname as author,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = 1) as likes,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = -1) as dislikes
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC
    `

    rows, err := database.DB.Query(query)
    if err != nil {
        json.NewEncoder(w).Encode(PostResponse{
            Status: "error",
            Error:  "Failed to fetch posts",
        })
        return
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        err := rows.Scan(
            &post.ID, 
            &post.UserID, 
            &post.Title, 
            &post.Content, 
            &post.CreatedAt,
            &post.Author, 
            &post.Likes, 
            &post.Dislikes,
        )
        if err != nil {
            continue
        }

        // Get categories for each post
        categories, err := getPostCategories(post.ID)
        if err == nil {
            post.Categories = categories
        }

        posts = append(posts, post)
    }

    // Return success response
    json.NewEncoder(w).Encode(PostResponse{
        Status: "success",
        Data:   posts,
    })
}

func getPostCategories(postID int) ([]string, error) {
    rows, err := database.DB.Query(`
        SELECT c.name 
        FROM categories c
        JOIN post_categories pc ON c.id = pc.category_id
        WHERE pc.post_id = ?
    `, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []string
    for rows.Next() {
        var category string
        if err := rows.Scan(&category); err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}

// CreatePostHandler creates a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Verify user is authenticated
    userID, ok := auth.ValidateSession(w, r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var post models.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate post
    if post.Title == "" || post.Content == "" {
        http.Error(w, "Title and content are required", http.StatusBadRequest)
        return
    }

    // Begin transaction
    tx, err := database.DB.Begin()
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Insert post
    result, err := tx.Exec(`
        INSERT INTO posts (user_id, title, content, created_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
        userID, post.Title, post.Content)
    if err != nil {
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

    postID, _ := result.LastInsertId()

    // Link categories
    for _, catID := range post.Categories {
        _, err = tx.Exec(
            "INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)",
            postID, catID)
        if err != nil {
            continue // Skip failed category associations
        }
    }

    if err = tx.Commit(); err != nil {
        http.Error(w, "Error saving post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(PostResponse{
        Status:  "success",
        Message: "Post created successfully",
        Data:    postID,
    })
}

// GetPostHandler retrieves a specific post
func GetPostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Extract post ID from URL
    postID := r.URL.Path[len("/api/post/"):]
    id, err := strconv.Atoi(postID)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    // Get post with categories and comments
    post, err := getPostWithDetails(id)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(PostResponse{
        Status: "success",
        Data:   post,
    })
}

func getPostWithDetails(postID int) (*models.Post, error) {
    var post models.Post
    err := database.DB.QueryRow(`
        SELECT 
            p.id, p.user_id, p.title, p.content, p.created_at,
            u.nickname as author,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = 1) as likes,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = -1) as dislikes
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = ?`, postID).Scan(
        &post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt,
        &post.Author, &post.Likes, &post.Dislikes)
    if err != nil {
        return nil, err
    }

    post.Categories, _ = getPostCategories(post.ID)
    return &post, nil
}

// func getPostCategories(postID int) ([]string, error) {
//     rows, err := database.DB.Query(`
//         SELECT c.name
//         FROM categories c
//         JOIN post_categories pc ON c.id = pc.category_id
//         WHERE pc.post_id = ?`, postID)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()

//     var categories []string
//     for rows.Next() {
//         var category string
//         if err := rows.Scan(&category); err != nil {
//             continue
//         }
//         categories = append(categories, category)
//     }
//     return categories, nil
// }

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    userID, ok := auth.ValidateSession(w, r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    parts := strings.Split(r.URL.Path, "/")
    postID, err := strconv.Atoi(parts[len(parts)-1])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    // Check if user owns the post
    var postUserID int
    err = database.DB.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&postUserID)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    if postUserID != userID {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Begin transaction
    tx, err := database.DB.Begin()
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Delete post and related data
    _, err = tx.Exec("DELETE FROM post_categories WHERE post_id = ?", postID)
    if err != nil {
        http.Error(w, "Error deleting post", http.StatusInternalServerError)
        return
    }

    _, err = tx.Exec("DELETE FROM likes WHERE post_id = ?", postID)
    if err != nil {
        http.Error(w, "Error deleting post", http.StatusInternalServerError)
        return
    }

    _, err = tx.Exec("DELETE FROM posts WHERE id = ?", postID)
    if err != nil {
        http.Error(w, "Error deleting post", http.StatusInternalServerError)
        return
    }

    if err = tx.Commit(); err != nil {
        http.Error(w, "Error deleting post", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(PostResponse{
        Status:  "success",
        Message: "Post deleted successfully",
    })
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    userID, ok := auth.ValidateSession(w, r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var post models.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Verify post ownership
    var postUserID int
    err := database.DB.QueryRow("SELECT user_id FROM posts WHERE id = ?", post.ID).Scan(&postUserID)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    if postUserID != userID {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Update post
    _, err = database.DB.Exec(`
        UPDATE posts 
        SET title = ?, content = ?
        WHERE id = ?`,
        post.Title, post.Content, post.ID)
    if err != nil {
        http.Error(w, "Error updating post", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(PostResponse{
        Status:  "success",
        Message: "Post updated successfully",
    })
}

func PostAndCommentHandler(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")

    // Get post ID from URL
    parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    if len(parts) < 3 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Invalid URL format",
        })
        return
    }

    postID, err := strconv.Atoi(parts[2])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Invalid post ID",
        })
        return
    }

    // Check if post exists
    var exists bool
    err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
    if err != nil || !exists {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{
            "status":  "error",
            "message": "Post not found",
        })
        return
    }

    if r.Method == http.MethodGet {
        GetPostHandler(w, r)
    } else if r.Method == http.MethodPost {
        CommentHandler(w, r)
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
