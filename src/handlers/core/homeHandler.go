package core

import (
    "encoding/json"
    "net/http"
    
    "Real-time-forum/src/database"
    "Real-time-forum/src/models"
)

type HomeResponse struct {
    Status  string        `json:"status"`
    Data    []models.Post `json:"data"`
    Error   string        `json:"error,omitempty"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get posts with author info
    rows, err := database.DB.Query(`
        SELECT p.id, p.title, p.content, u.nickname as author, p.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC`)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt)
        if err != nil {
            continue
        }
        // Get likes/dislikes count
        post.Likes, post.Dislikes = getPostLikesDislikes(post.ID)
        posts = append(posts, post)
    }

    // Return JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(HomeResponse{
        Status: "success",
        Data:   posts,
    })
}

func getPostLikesDislikes(postID int) (likes int, dislikes int) {
    err := database.DB.QueryRow(`
        SELECT 
            SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END) as likes,
            SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END) as dislikes
        FROM likes 
        WHERE post_id = ?`, postID).Scan(&likes, &dislikes)
    if err != nil {
        return 0, 0
    }
    return
}