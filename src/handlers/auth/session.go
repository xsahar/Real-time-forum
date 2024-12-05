package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"Real-time-forum/src/database"
)

func ValidateSession(w http.ResponseWriter, r *http.Request) (int, bool) {
    // Check Authorization header first
    authHeader := r.Header.Get("Authorization")
    var token string
    
    if authHeader != "" {
        // Extract token from Bearer format
        if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
            token = authHeader[7:]
        }
    } else {
        // Fallback to cookie
        cookie, err := r.Cookie("session_token")
        if err != nil {
            return 0, false
        }
        token = cookie.Value
    }

    if token == "" {
        return 0, false
    }

    var userID int
    err := database.DB.QueryRow(
        "SELECT user_id FROM sessions WHERE token = ? AND expires_at > ?",
        token, time.Now(),
    ).Scan(&userID)

    if err != nil {
        return 0, false
    }

    return userID, true
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Get session cookie
    cookie, err := r.Cookie("session_token")
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Not logged in",
        })
        return
    }

    // Verify session exists
    var userID int
    err = database.DB.QueryRow(
        "SELECT user_id FROM sessions WHERE token = ?",
        cookie.Value,
    ).Scan(&userID)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Not logged in",
        })
        return
    }

    // Delete session
    _, err = database.DB.Exec("DELETE FROM sessions WHERE token = ?", cookie.Value)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Failed to logout",
        })
        return
    }

    // Clear cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    "",
        Expires:  time.Now().Add(-1 * time.Hour),
        HttpOnly: true,
    })

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(AuthResponse{
        Status:  "success",
        Message: "Logged out successfully",
    })
}