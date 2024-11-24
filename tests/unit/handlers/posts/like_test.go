package posts

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "Real-time-forum/src/database"
    "Real-time-forum/src/handlers/posts"
)

func TestLikeHandler(t *testing.T) {
    userID, sessionToken := setupTestUser()
    
    result, err := database.DB.Exec(`
        INSERT INTO posts (user_id, title, content)
        VALUES (?, ?, ?)`,
        userID, "Test Post", "Test Content")
    if err != nil {
        t.Fatalf("Failed to create test post: %v", err)
    }
    postID, _ := result.LastInsertId()

    // Verify post exists
    var exists bool
    err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
    if err != nil || !exists {
        t.Fatalf("Test post not created properly: %v", err)
    }
    t.Logf("Created test post with ID: %d", postID)

    type response struct {
        Status  string `json:"status"`
        Message string `json:"message,omitempty"`
    }

    testCases := []struct {
        name         string
        postID       int64
        action       string
        expectedCode int
        expectedResp response
    }{
        {
            name:    "Valid like",
            postID:  postID,
            action:  "like",
            expectedCode: http.StatusOK,
            expectedResp: response{
                Status:  "success",
                Message: "Like added successfully",
            },
        },
        {
            name:    "Valid dislike",
            postID:  postID,
            action:  "dislike",
            expectedCode: http.StatusOK,
            expectedResp: response{
                Status:  "success",
                Message: "Dislike added successfully",
            },
        },
        {
            name:    "Invalid post ID",
            postID:  9999,
            action:  "like",
            expectedCode: http.StatusNotFound,
            expectedResp: response{
                Status:  "error",
                Message: "Post not found",
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Verify post still exists
            if tc.postID == postID {
                var stillExists bool
                err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", tc.postID).Scan(&stillExists)
                if err != nil || !stillExists {
                    t.Fatalf("Test post %d no longer exists: %v", tc.postID, err)
                }
            }

            // Updated URL format to match handler expectation
            url := fmt.Sprintf("/api/%s/%d", tc.action, tc.postID)
            t.Logf("Request URL: %s", url)

            req := httptest.NewRequest(http.MethodPost, url, nil)
            req.Header.Set("Content-Type", "application/json")
            addAuthCookie(req, sessionToken)

            // Log URL components
            t.Logf("URL segments: %v", strings.Split(req.URL.Path, "/"))

            w := httptest.NewRecorder()

            if tc.action == "like" {
                posts.LikeHandler(w, req)
            } else {
                posts.DislikeHandler(w, req)
            }

            t.Logf("Response Status: %d", w.Code)
            t.Logf("Response Body: %s", w.Body.String())

            if w.Code != tc.expectedCode {
                t.Errorf("Expected status code %d, got %d", tc.expectedCode, w.Code)
                return
            }

            var resp response
            if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
                t.Fatalf("Failed to decode response: %v\nResponse body: %s", err, w.Body.String())
            }

            if resp != tc.expectedResp {
                t.Errorf("Expected response %+v, got %+v", tc.expectedResp, resp)
            }

            // Verify database for successful operations
            if w.Code == http.StatusOK {
                var count int
                query := `SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ? AND value = ?`
                value := 1
                if tc.action == "dislike" {
                    value = -1
                }
                err := database.DB.QueryRow(query, tc.postID, userID, value).Scan(&count)
                if err != nil {
                    t.Errorf("Failed to verify like in database: %v", err)
                }
                if count != 1 {
                    t.Errorf("Expected 1 like record, got %d", count)
                }
            }
        })
    }
}