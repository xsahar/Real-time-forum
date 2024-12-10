package crud

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"

    "Real-time-forum/src/database"
    "Real-time-forum/src/handlers/auth"
)

type LikeResponse struct {
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
    Data    struct {
        Likes    int `json:"likes"`
        Dislikes int `json:"dislikes"`
    } `json:"data,omitempty"`
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Parse URL path first for logging
    parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    if len(parts) != 3 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "message": "Invalid URL format",
        })
        return
    }

    log.Printf("LikeHandler: Processing request for post %s", parts[2])
    
    userID, ok := auth.ValidateSession(w, r)
    if !ok {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "message": "Unauthorized",
        })
        return
    }

    // Get post ID
    postID, err := strconv.Atoi(parts[2])
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "message": "Invalid post ID",
        })
        return
    }

    // Determine action type
    value := 1
    action := "Like"
    if parts[1] == "dislike" {
        value = -1
        action = "Dislike"
    }

    // Start transaction
    tx, err := database.DB.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "message": "Database error",
        })
        return
    }
    defer tx.Rollback()

    // Check if post exists
    var exists bool
    err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
    if err != nil {
        log.Printf("Error checking post existence: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "message": "Database error",
        })
        return
    }
    if !exists {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "message": "Post not found",
        })
        return
    }

    // Update likes within transaction
    _, err = tx.Exec(`
        INSERT INTO likes (post_id, user_id, value)
        VALUES (?, ?, ?)
        ON CONFLICT(post_id, user_id) DO UPDATE SET value = ?`,
        postID, userID, value, value)
    
    if err != nil {
        log.Printf("Error updating like: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "message": "Failed to update like",
        })
        return
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "message": "Database error",
        })
        return
    }

    // Get updated counts
    likes, dislikes := getPostLikesDislikes(postID)

    // Success response
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "success",
        "message": fmt.Sprintf("%s added successfully", action),
        "data": map[string]int{
            "likes":    likes,
            "dislikes": dislikes,
        },
    })
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
    // Extract post ID from original path
    parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    if len(parts) != 3 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "message": "Invalid URL format",
        })
        return
    }

    // Construct new path with dislike action
    postID := parts[2]
    r.URL.Path = fmt.Sprintf("/api/dislike/%s", postID)
    
    // Call LikeHandler
    LikeHandler(w, r)
}

func getPostLikesDislikes(postID int) (likes int, dislikes int) {
    err := database.DB.QueryRow(`
        SELECT 
            COUNT(CASE WHEN value = 1 THEN 1 END) as likes,
            COUNT(CASE WHEN value = -1 THEN 1 END) as dislikes
        FROM likes 
        WHERE post_id = ?`, postID).Scan(&likes, &dislikes)
    if err != nil {
        log.Printf("Error getting like counts: %v", err)
        return 0, 0
    }
    return
}