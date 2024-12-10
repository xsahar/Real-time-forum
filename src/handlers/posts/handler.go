package posts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	// "time"

	"Real-time-forum/src/database"
	"Real-time-forum/src/handlers/auth"
	"Real-time-forum/src/models"
    "Real-time-forum/src/handlers/posts/crud"
    
)

type PostResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}



func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    log.Printf("ListPostsHandler: Received %s request", r.Method)

    if r.Method != http.MethodGet {
        json.NewEncoder(w).Encode(PostResponse{
            Status: "error",
            Error:  "Method not allowed",
        })
        return
    }

    // Use GetAllPosts from database package
    posts, err := database.GetAllPosts()
    if err != nil {
        log.Printf("ListPostsHandler: Database error: %v", err)
        json.NewEncoder(w).Encode(PostResponse{
            Status: "error",
            Error:  fmt.Sprintf("Failed to fetch posts: %v", err),
        })
        return
    }

    response := PostResponse{
        Status: "success",
        Data:   posts,
    }
    
    if len(posts) == 0 {
        log.Printf("ListPostsHandler: No posts found")
        response.Message = "No posts available"
    }

    json.NewEncoder(w).Encode(response)
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
    post, err := database.GetPostWithDetails(id)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(PostResponse{
        Status: "success",
        Data:   post,
    })
}

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

    // Handle methods
    switch r.Method {
    case http.MethodGet:
        GetPostHandler(w, r)
    case http.MethodPost:
        crud.CommentHandler(w, r)  // Use the imported crud package
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
