package posts

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "fmt"

    "Real-time-forum/src/database"
    "Real-time-forum/src/handlers/posts"
    "Real-time-forum/src/models"
)

func TestAddComment(t *testing.T) {
	
    userID, sessionToken := setupTestUser()
    
    result, err := database.DB.Exec(`
        INSERT INTO posts (user_id, title, content)
        VALUES (?, ?, ?)`,
        userID, "Test Post", "Test Content")
    if err != nil {
        t.Fatalf("Failed to create test post: %v", err)
    }
    postID, _ := result.LastInsertId()

    testCases := []struct {
        name         string
        comment      models.Comment
        expectedCode int
        expectedResp struct {
            Status  string `json:"status"`
            Message string `json:"message,omitempty"`
        }
    }{
        {
            name: "Valid comment",
            comment: models.Comment{
                PostID:  int(postID),
                Content: "Test comment",
            },
            expectedCode: http.StatusCreated,
            expectedResp: struct {
                Status  string `json:"status"`
                Message string `json:"message,omitempty"`
            }{
                Status:  "success",
                Message: "Comment added successfully",
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            body, _ := json.Marshal(tc.comment)
            req := httptest.NewRequest(http.MethodPost, 
                fmt.Sprintf("/api/post/%d/comment", tc.comment.PostID),
                bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            addAuthCookie(req, sessionToken)
            w := httptest.NewRecorder()

            posts.CommentHandler(w, req)

            if w.Code != tc.expectedCode {
                t.Errorf("Expected status code %d, got %d", tc.expectedCode, w.Code)
            }

            var resp struct {
                Status  string `json:"status"`
                Message string `json:"message,omitempty"`
            }
            if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
                t.Fatalf("Failed to decode response: %v", err)
            }

            if resp != tc.expectedResp {
                t.Errorf("Expected response %+v, got %+v", tc.expectedResp, resp)
            }
        })
    }
}