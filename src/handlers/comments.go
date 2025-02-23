package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/src/models"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetCommentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract postID from URL path
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		postID := pathParts[3]

		// Validate postID
		if postID == "" {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Query comments with author info
		rows, err := db.Query(`
            SELECT 
                c.id, c.post_id, c.user_id, c.content, c.created_at,
                u.nickname, u.first_name, u.last_name, u.email
            FROM comments c
            JOIN users u ON c.user_id = u.id
            WHERE c.post_id = ?
            ORDER BY c.created_at DESC
        `, postID)

		if err != nil {
			log.Printf("Error querying comments: %v", err)
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		comments := []models.Comment{}
		for rows.Next() {
			var c models.Comment
			var author models.User

			if err := rows.Scan(
				&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt,
				&author.Nickname, &author.FirstName, &author.LastName, &author.Email,
			); err != nil {
				log.Printf("Error scanning comment row: %v", err)
				continue
			}

			c.Author = &author
			comments = append(comments, c)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error iterating comments: %v", err)
			http.Error(w, "Error processing comments", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

func CreateCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract postID from URL path
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		postID := pathParts[3]

		// Validate postID
		if postID == "" {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Parse request body
		var comment models.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate comment content
		if comment.Content == "" {
			http.Error(w, "Comment content is required", http.StatusBadRequest)
			return
		}

		// Add user ID validation
		if comment.UserID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		// Add post existence check
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
		if err != nil {
			log.Printf("Error checking post existence: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Set comment metadata
		comment.ID = uuid.New().String()
		comment.PostID = postID
		comment.CreatedAt = time.Now()

		// Insert comment
		_, err = db.Exec(`
			INSERT INTO comments (id, post_id, user_id, content, created_at)
			VALUES (?, ?, ?, ?, ?)
		`, comment.ID, comment.PostID, comment.UserID, comment.Content, comment.CreatedAt)

		if err != nil {
			log.Printf("Error creating comment: %v", err)
			http.Error(w, "Error creating comment", http.StatusInternalServerError)
			return
		}

		// Fetch the created comment with author info
		var author models.User
		err = db.QueryRow(`
            SELECT u.nickname, u.first_name, u.last_name, u.email
            FROM users u
            WHERE u.id = ?
        `, comment.UserID).Scan(&author.Nickname, &author.FirstName, &author.LastName, &author.Email)

		if err != nil {
			log.Printf("Error fetching comment author: %v", err)
			http.Error(w, "Comment created but failed to fetch author details", http.StatusInternalServerError)
			return
		}

		comment.Author = &author

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comment)
	}
}
