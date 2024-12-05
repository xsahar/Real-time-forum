package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Real-time-forum/src/database"
	"Real-time-forum/src/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

      if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        log.Printf("Invalid method %s", r.Method)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Method not allowed",
        })
        return
    }

    // Read and log request body
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        log.Printf("Failed to decode request body: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Invalid request body",
        })
        return
    }

		log.Printf("Received user data: %+v", user)

    // Validate all required fields
    if user.Email == "" || user.Password == "" || user.Nickname == "" || 
       user.FirstName == "" || user.LastName == "" || user.Gender == "" || 
       user.Age < 13 {
        log.Printf("Missing or invalid fields: %+v", user)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Missing required fields",
        })
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Internal server error",
        })
        return
    }

    // Begin transaction
    tx, err := database.DB.Begin()
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Database error",
        })
        return
    }
    defer tx.Rollback()

    // Insert user
    result, err := tx.Exec(`
        INSERT INTO users (email, nickname, password, age, gender, first_name, last_name, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
        user.Email, user.Nickname, hashedPassword, user.Age, user.Gender, user.FirstName, user.LastName)
    
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Failed to create user",
        })
        return
    }

    userID, _ := result.LastInsertId()
    sessionToken := uuid.NewString()

    // Create session
    _, err = tx.Exec(
        "INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
        userID, sessionToken, time.Now().Add(24*time.Hour))
    
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Session creation failed",
        })
        return
    }

    if err = tx.Commit(); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Transaction failed",
        })
        return
    }

    // Set session cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    sessionToken,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
    })

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(AuthResponse{
        Status:  "success",
        Message: "Registration successful",
    })
}