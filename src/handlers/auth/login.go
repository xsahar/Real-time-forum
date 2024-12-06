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

func cleanExpiredSessions() error {
    _, err := database.DB.Exec("DELETE FROM sessions WHERE expires_at < datetime('now')")
    return err
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if err := cleanExpiredSessions(); err != nil {
        log.Printf("Failed to clean expired sessions: %v", err)
    }

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

    log.Printf("Login attempt for identifier: %s", creds.Identifier)

    if creds.Identifier == "" || creds.Password == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Identifier and password are required",
        })
        return
    }

    var user models.User
    err := database.DB.QueryRow(`
        SELECT id, email, nickname, password, first_name, last_name 
        FROM users 
        WHERE email = ? OR nickname = ?`,
        creds.Identifier, creds.Identifier,
    ).Scan(&user.ID, &user.Email, &user.Nickname, &user.Password, &user.FirstName, &user.LastName)

    if err != nil {
        log.Printf("User lookup failed: %v", err)
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Invalid credentials",
        })
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
        log.Printf("Password verification failed for user: %s", user.Email)
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Invalid credentials",
        })
        return
    }

    tx, err := database.DB.Begin()
    if err != nil {
        log.Printf("Transaction start failed: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Internal server error",
        })
        return
    }
    defer tx.Rollback()

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

    sessionToken := uuid.NewString()
    expiresAt := time.Now().Add(24 * time.Hour)
    
    _, err = tx.Exec(
        "INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
        user.ID, sessionToken, expiresAt,
    )
    if err != nil {
        log.Printf("Session creation failed: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Session creation failed",
        })
        return
    }

    if err = tx.Commit(); err != nil {
        log.Printf("Transaction commit failed: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Session commitment failed",
        })
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    sessionToken,
        Expires:  expiresAt,
        HttpOnly: true,
        Path:     "/",
        SameSite: http.SameSiteStrictMode,
    })

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(AuthResponse{
        Status: "success",
        Data: map[string]interface{}{
            "token": sessionToken,
            "user": map[string]interface{}{
                "id":        user.ID,
                "email":     user.Email,
                "nickname":  user.Nickname,
                "firstName": user.FirstName,
                "lastName":  user.LastName,
            },
        },
        Message: "Login successful",
    })
    
    log.Printf("User logged in successfully: %s", user.Email)
}