package auth

import (
    "encoding/json"
    "net/http"
    "time"

		"Real-time-forum/src/models" 
		"Real-time-forum/src/database"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

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

    var creds struct {
        Identifier    string `json:"identifier"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Invalid request body",
        })
        return
    }

    // Find user by email or nickname
    var user models.User
    err := database.DB.QueryRow(
        "SELECT id, password FROM users WHERE email = ? OR nickname = ?",
        creds.Identifier, creds.Identifier,
    ).Scan(&user.ID, &user.Password)

    if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Invalid credentials",
        })
        return
    }

     // Generate token
     sessionToken := uuid.NewString()

     // Create session with expiry
     expiresAt := time.Now().Add(24 * time.Hour)
     _, err = database.DB.Exec(
         "INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
         user.ID, sessionToken, expiresAt,
     )

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Session creation failed",
        })
        return
    }

     // Return success with token
     w.WriteHeader(http.StatusOK)
     json.NewEncoder(w).Encode(AuthResponse{
         Status: "success",
         Message: "Login successful",
         Data: map[string]string{
             "token": sessionToken,
         },
     })

    // Set session cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    sessionToken,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
    })

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(AuthResponse{
        Status:  "success",
        Message: "Login successful",
    })
}