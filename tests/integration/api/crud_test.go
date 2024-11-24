// tests/integration/api/crud_test.go
package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"Real-time-forum/src/database"
	"Real-time-forum/src/handlers/auth"
	"Real-time-forum/src/handlers/posts"
	"Real-time-forum/src/models"
)

func setupTestDB(t *testing.T) {
	// Create tables
	queries := []string{
			`CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					email TEXT UNIQUE,
					nickname TEXT UNIQUE,
					password TEXT,
					age INTEGER,
					gender TEXT,
					first_name TEXT,
					last_name TEXT
			)`,
			`CREATE TABLE IF NOT EXISTS sessions (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER,
					token TEXT UNIQUE,
					expires_at TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id)
			)`,
			`CREATE TABLE IF NOT EXISTS posts (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER,
					      title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`,
    }

    for _, query := range queries {
        if _, err := database.DB.Exec(query); err != nil {
            t.Fatalf("Failed to create table: %v", err)
        }
    }
}

// Setup test user and session
func setupTestUserAndSession(t *testing.T) (int, string) {
    result, err := database.DB.Exec(`
        INSERT INTO users (email, nickname, password)
        VALUES (?, ?, ?)`,
        "test@example.com", "testuser", "hashedpassword")
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    userID, _ := result.LastInsertId()

    // Create session
    sessionToken := "test-session-token"
    _, err = database.DB.Exec(`
        INSERT INTO sessions (user_id, token, expires_at)
        VALUES (?, ?, ?)`,
        userID, sessionToken, time.Now().Add(24*time.Hour))
    if err != nil {
        t.Fatalf("Failed to create session: %v", err)
    }

    return int(userID), sessionToken
}

func TestCRUDOperations(t *testing.T) {
    // Setup
    if err := database.InitDB(); err != nil {
        t.Fatalf("Failed to init DB: %v", err)
    }
    defer database.DB.Close()

    // Setup test database
    setupTestDB(t)

    // Clear test data
    database.DB.Exec("DELETE FROM posts")
    database.DB.Exec("DELETE FROM users")
    database.DB.Exec("DELETE FROM sessions")

    // Setup test user and session
    userID, sessionToken := setupTestUserAndSession(t)

		var postID int
    t.Run("Create post", func(t *testing.T) {
        body, _ := json.Marshal(models.Post{
            UserID:  userID,
            Title:   "Test Post",
            Content: "Test Content",
        })

        req := httptest.NewRequest("POST", "/api/create-post", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.AddCookie(&http.Cookie{
            Name:  "session_token",
            Value: sessionToken,
        })
        
        w := httptest.NewRecorder()
        posts.CreatePostHandler(w, req)

        var resp struct {
            Status  string `json:"status"`
            Message string `json:"message"`
            Data    int    `json:"data"`
        }
				json.NewDecoder(w.Body).Decode(&resp)
        postID = resp.Data
    })

    // Test cases
    tests := []struct {
        name         string
        method       string
        path         string
        body        interface{}
        wantStatus   int
        wantResponse map[string]interface{}
    }{
        {
            name:    "Create post",
            method:  "POST",
            path:    "/api/create-post",
            body:    models.Post{
                UserID:  userID,
                Title:   "Test Post",
                Content: "Test Content",
            },
            wantStatus: http.StatusCreated,
            wantResponse: map[string]interface{}{
                "status":  "success",
                "message": "Post created successfully",
            },
        },
       {
            name:    "Get post",
            method:  "GET",
            path:    fmt.Sprintf("/api/post/%d", postID),
            wantStatus: http.StatusOK,
            wantResponse: map[string]interface{}{
                "status": "success",
                "data": map[string]interface{}{
                    "title":   "Test Post",
                    "content": "Test Content",
                },
            },
        },
				{
					name:    "Update post",
					method:  "PUT",
					path:    fmt.Sprintf("/api/post/%d", postID),
					body:    models.Post{
							Title:   "Updated Post",
							Content: "Updated Content",
					},
					wantStatus: http.StatusOK,
					wantResponse: map[string]interface{}{
							"status":  "success",
							"message": "Post updated successfully",
					},
			},
			{
				name:    "Get updated post",
				method:  "GET",
				path:    fmt.Sprintf("/api/post/%d", postID),
				wantStatus: http.StatusOK,
				wantResponse: map[string]interface{}{
						"status": "success",
						"data": map[string]interface{}{
								"title":   "Updated Post",
								"content": "Updated Content",
						},
				},
		},
		{
			name:    "Delete post",
			method:  "DELETE",
			path:    fmt.Sprintf("/api/post/%d", postID),
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
					"status":  "success",
					"message": "Post deleted successfully",
			},
	},
	{
			name:    "Get deleted post",
			method:  "GET",
			path:    fmt.Sprintf("/api/post/%d", postID),
			wantStatus: http.StatusNotFound,
			wantResponse: map[string]interface{}{
					"status": "error",
					"message": "Post not found",
			},
	},
}
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // Prepare request
            var body []byte
            if tc.body != nil {
                var err error
                body, err = json.Marshal(tc.body)
                if err != nil {
                    t.Fatalf("Failed to marshal request body: %v", err)
                }
            }

            // Create request
            req := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            req.AddCookie(&http.Cookie{
                Name:  "session_token",
                Value: sessionToken,
            })
            
            w := httptest.NewRecorder()
            
            // Route request to appropriate handler
            switch tc.method {
            case "POST":
                posts.CreatePostHandler(w, req)
            case "GET":
                posts.PostAndCommentHandler(w, req)
            case "PUT":
                UpdatePostHandler(w, req)
            case "DELETE":
                DeletePostHandler(w, req)
            }

            // Debug output
            t.Logf("Request: %s %s", tc.method, tc.path)
            t.Logf("Response Status: %d", w.Code)
            t.Logf("Response Body: %s", w.Body.String())

            // Check status code
            if w.Code != tc.wantStatus {
                t.Errorf("got status %d, want %d", w.Code, tc.wantStatus)
            }

            // Verify response body
            if tc.wantResponse != nil {
                var got map[string]interface{}
                if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
                    t.Fatalf("Failed to decode response body: %v", err)
                }

                // Compare status
                if got["status"] != tc.wantResponse["status"] {
                    t.Errorf("got status %v, want %v", got["status"], tc.wantResponse["status"])
                }

                // Compare message if present
                if tc.wantResponse["message"] != nil {
                    if got["message"] != tc.wantResponse["message"] {
                        t.Errorf("got message %v, want %v", got["message"], tc.wantResponse["message"])
                    }
                }

                // Compare data if present
                if tc.wantResponse["data"] != nil {
                    wantData := tc.wantResponse["data"].(map[string]interface{})
                    if gotData, ok := got["data"].(map[string]interface{}); ok {
                        if gotData["title"] != wantData["title"] {
                            t.Errorf("got title %v, want %v", gotData["title"], wantData["title"])
                        }
                        if gotData["content"] != wantData["content"] {
                            t.Errorf("got content %v, want %v", gotData["content"], wantData["content"])
                        }
                    }
                }
            }
        })
    }
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Verify auth
    userID, ok := auth.ValidateSession(w, r)
    if !ok {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Unauthorized",
        })
        return
    }

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

    // Verify post exists and user owns it
    var exists bool
    err = database.DB.QueryRow(
        "SELECT EXISTS(SELECT 1 FROM posts WHERE id = ? AND user_id = ?)",
        postID, userID,
    ).Scan(&exists)

    if err != nil || !exists {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Post not found or unauthorized",
        })
        return
    }

    // Begin transaction
    tx, err := database.DB.Begin()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Database error",
        })
        return
    }
    defer tx.Rollback()

    // Delete associated data first
    _, err = tx.Exec("DELETE FROM likes WHERE post_id = ?", postID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Failed to delete likes",
        })
        return
    }

    _, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", postID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Failed to delete comments",
        })
        return
    }

    // Delete post
    _, err = tx.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, userID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Failed to delete post",
        })
        return
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Failed to commit changes",
        })
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "success",
        "message": "Post deleted successfully",
    })
}
func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Verify auth
    userID, ok := auth.ValidateSession(w, r)
    if !ok {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Unauthorized",
        })
        return
    }

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

    // Parse request body
    var post models.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Invalid request body",
        })
        return
    }

    // Verify post exists and user owns it
    var exists bool
    err = database.DB.QueryRow(
        "SELECT EXISTS(SELECT 1 FROM posts WHERE id = ? AND user_id = ?)",
        postID, userID,
    ).Scan(&exists)

    if err != nil || !exists {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Post not found or unauthorized",
        })
        return
    }

    // Update post
    _, err = database.DB.Exec(
        "UPDATE posts SET title = ?, content = ? WHERE id = ? AND user_id = ?",
        post.Title, post.Content, postID, userID,
    )
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "error",
            "error":  "Failed to update post",
        })
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "success",
        "message": "Post updated successfully",
    })
}