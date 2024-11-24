package posts

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Real-time-forum/src/database"
	"Real-time-forum/src/handlers/posts"
	"Real-time-forum/src/models"
)

func TestMain(m *testing.M) {
    // Setup
    if err := database.InitDB(); err != nil {
        log.Fatal("Failed to init DB:", err)
    }

    // Create tables
    _, err := database.DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT UNIQUE,
            nickname TEXT UNIQUE,
            password TEXT,
            age INTEGER,
            gender TEXT,
            first_name TEXT,
            last_name TEXT
        )`)
    if err != nil {
        log.Fatal("Failed to create users table:", err)
    }

    _, err = database.DB.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER,
            token TEXT UNIQUE,
            expires_at TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`)
    if err != nil {
        log.Fatal("Failed to create sessions table:", err)
    }

    _, err = database.DB.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`)
    if err != nil {
        log.Fatal("Failed to create posts table:", err)
    }

    // Run tests
    code := m.Run()

    // Cleanup
    
    os.Exit(code)
}

func TestCreatePost(t *testing.T) {
	setupTestUser()
    userID, sessionToken := setupTestUser()

    testCases := []struct {
        name          string
        input         models.Post
        expectedCode  int
        expectedResp  struct {
            Status  string `json:"status"`
            Message string `json:"message,omitempty"`
        }
    }{
        {
            name: "Valid post",
            input: models.Post{
							UserID: 	 userID,
                Title:      "Test Post",
                Content:    "Test Content",
                Categories: []string{"test"},
            },
            expectedCode: http.StatusCreated,
            expectedResp: struct {
                Status  string `json:"status"`
                Message string `json:"message,omitempty"`
            }{
                Status:  "success",
                Message: "Post created successfully",
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            body, _ := json.Marshal(tc.input)
            req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            addAuthCookie(req, sessionToken)
            w := httptest.NewRecorder()

            posts.CreatePostHandler(w, req)

            if w.Code != tc.expectedCode {
                t.Errorf("Expected status %d, got %d", tc.expectedCode, w.Code)
            }

            var resp struct {
                Status  string `json:"status"`
                Message string `json:"message,omitempty"`
            }
            json.NewDecoder(w.Body).Decode(&resp)

            if resp != tc.expectedResp {
                t.Errorf("Expected response %+v, got %+v", tc.expectedResp, resp)
            }
        })
    }
}