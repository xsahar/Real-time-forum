package auth

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "Real-time-forum/src/models"
    "Real-time-forum/src/database"
    
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type LoginCredentials struct {
    Identifier string `json:"identifier"`
    Password   string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Method not allowed",
        })
        return
    }

    var creds LoginCredentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        log.Printf("Failed to decode login request: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Invalid request format",
        })
        return
    }

    // Validate credentials
    if creds.Identifier == "" || creds.Password == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Identifier and password are required",
        })
        return
    }

    // Find user by email or nickname
    var user models.User
    err := database.DB.QueryRow(
        "SELECT id, password FROM users WHERE email = ? OR nickname = ?",
        creds.Identifier, creds.Identifier,
    ).Scan(&user.ID, &user.Password)

    if err != nil {
        log.Printf("User not found: %v", err)
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Invalid credentials",
        })
        return
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
        log.Printf("Password verification failed: %v", err)
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Invalid credentials",
        })
        return
    }

    // Start transaction for session management
    tx, err := database.DB.Begin()
    if err != nil {
        log.Printf("Failed to begin transaction: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Internal server error",
        })
        return
    }
    defer tx.Rollback()

    // Clear existing sessions for user
    _, err = tx.Exec("DELETE FROM sessions WHERE user_id = ?", user.ID)
    if err != nil {
        log.Printf("Failed to clear existing sessions: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Session management failed",
        })
        return
    }

    // Create new session
    sessionToken := uuid.NewString()
    expiresAt := time.Now().Add(24 * time.Hour)
    
    _, err = tx.Exec(
        "INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
        user.ID, sessionToken, expiresAt,
    )
    if err != nil {
        log.Printf("Failed to create session: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Session creation failed",
        })
        return
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        log.Printf("Failed to commit transaction: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Session commitment failed",
        })
        return
    }

    // Set session cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    sessionToken,
        Expires:  expiresAt,
        HttpOnly: true,
        Path:     "/",
        SameSite: http.SameSiteStrictMode,
    })

    // Single response with token
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(AuthResponse{
        Status:  "success",
        Message: "Login successful",
        Data: map[string]string{
            "token": sessionToken,
        },
    })
}

// Removed duplicate LogoutHandler function