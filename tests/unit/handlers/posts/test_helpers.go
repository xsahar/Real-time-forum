// tests/unit/handlers/posts/test_helpers.go
package posts

import (
	"Real-time-forum/src/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func setupTestUser() (int, string) {

	database.DB.Exec("DELETE FROM likes")
	database.DB.Exec("DELETE FROM comments")
	database.DB.Exec("DELETE FROM posts")
	database.DB.Exec("DELETE FROM sessions")
	database.DB.Exec("DELETE FROM users")

      // Create test user with unique values
			result, err := database.DB.Exec(`
			INSERT INTO users (email, nickname, password)
			VALUES (?, ?, ?)`,
			fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),  // Unique email
			fmt.Sprintf("testuser%d", time.Now().UnixNano()),          // Unique nickname
			"hashedpassword")
	if err != nil {
			log.Fatalf("Failed to create test user: %v", err)
	}
	userID, _ := result.LastInsertId()

	// Create session
	sessionToken := uuid.NewString()
	_, err = database.DB.Exec(`
			INSERT INTO sessions (user_id, token, expires_at)
			VALUES (?, ?, ?)`,
			userID, sessionToken, time.Now().Add(24*time.Hour))
	if err != nil {
			log.Fatalf("Failed to create session: %v", err)
	}

    return int(userID), sessionToken
}

func addAuthCookie(r *http.Request, sessionToken string) {
    cookie := &http.Cookie{
        Name:     "session_token",
        Value:    sessionToken,
        Expires:  time.Now().Add(24 * time.Hour),
        Path:     "/",
        HttpOnly: true,
    }
    r.AddCookie(cookie)
}