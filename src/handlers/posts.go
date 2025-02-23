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

// CreatePostHandler handles new post creation
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var input struct {
            Title    string `json:"title"`
            Content  string `json:"content"`
            Category string `json:"category"`
            UserID   string `json:"userId"`
        }

        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Create new post with proper field order
				post := models.Post{
					ID:        uuid.New().String(),
					UserID:    input.UserID,
					Title:     input.Title,
					Content:   input.Content,
					Category:  input.Category,
					CreatedAt: time.Now(),
			}

			// Validate input
if post.Title == "" || post.Content == "" || post.Category == "" || post.UserID == "" {
	http.Error(w, "Missing required fields", http.StatusBadRequest)
	return
}


        // Insert post
				_, err := db.Exec(`
				INSERT INTO posts (id, user_id, title, content, category, created_at)
				VALUES (?, ?, ?, ?, ?, ?)
		`, post.ID, post.UserID, post.Title, post.Content, post.Category, post.CreatedAt)

        if err != nil {
            log.Printf("Error creating post: %v", err)
            http.Error(w, "Failed to create post", http.StatusInternalServerError)
            return
        }

        // Fetch complete post with author
        var author models.User
        err = db.QueryRow(`
            SELECT p.*, u.nickname, u.first_name, u.last_name 
            FROM posts p
            JOIN users u ON p.user_id = u.id
            WHERE p.id = ?
        `, post.ID).Scan(
            &post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.CreatedAt,
            &author.Nickname, &author.FirstName, &author.LastName,
        )

        if err != nil {
            log.Printf("Error fetching created post: %v", err)
            http.Error(w, "Post created but failed to fetch details", http.StatusInternalServerError)
            return
        }

        post.Author = &author
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(post)
    }
}

// GetPostsHandler retrieves all posts with optional category filter
func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")

		var rows *sql.Rows
		var err error

		if category != "" {
			rows, err = db.Query(`
                SELECT p.*, u.nickname, u.email 
                FROM posts p
                JOIN users u ON p.user_id = u.id
                WHERE p.category = ?
                ORDER BY p.created_at DESC
            `, category)
		} else {
			rows, err = db.Query(`
                SELECT p.*, u.nickname, u.email 
                FROM posts p
                JOIN users u ON p.user_id = u.id
                ORDER BY p.created_at DESC
            `)
		}

		if err != nil {
			log.Printf("Error fetching posts: %v", err)
			http.Error(w, "Error fetching posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var post models.Post
			var author models.User
			err := rows.Scan(
				&post.ID, &post.UserID, &post.Title, &post.Content,
				&post.Category, &post.CreatedAt, &author.Nickname, &author.Email,
			)
			if err != nil {
				log.Printf("Error scanning post: %v", err)
				continue
			}
			post.Author = &author
			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func GetPostByIDHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        parts := strings.Split(r.URL.Path, "/")
        if len(parts) != 4 { // /api/posts/{id}
            http.Error(w, "Invalid path", http.StatusBadRequest)
            return
        }
        postID := parts[3]

        var post models.Post
        var author models.User
				err := db.QueryRow(`
				SELECT 
						p.id, p.user_id, p.title, p.content, p.category, p.created_at,
						u.nickname, u.first_name, u.last_name
				FROM posts p
				JOIN users u ON p.user_id = u.id
				WHERE p.id = ?
		`, postID).Scan(
				&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.CreatedAt,
				&author.Nickname, &author.FirstName, &author.LastName,
		)

        if err == sql.ErrNoRows {
            http.Error(w, "Post not found", http.StatusNotFound)
            return
        }
        if err != nil {
            log.Printf("Error fetching post: %v", err)
            http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
            return
        }

        post.Author = &author
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(post)
    }
}
