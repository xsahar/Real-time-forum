package auth

import (
    // "bytes"
    "encoding/json"
    "fmt"
    // "io"
    "log"
    "net/http"
    "strings"
    "time"

    "Real-time-forum/src/database"
    "Real-time-forum/src/models"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

// Helper function for safe logging
func sanitizeUserDataForLog(user models.User) models.User {
    sanitized := user
    sanitized.Password = "[REDACTED]"
    return sanitized
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Method not allowed",
        })
        return
    }

    // Decode request body
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        log.Printf("Failed to decode request body: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Invalid JSON format",
        })
        return
    }

    // Single log with sanitized data
    log.Printf("Processing registration request: %+v", sanitizeUserDataForLog(user))

    // Validate required fields
    if err := validateUserData(&user); err != nil {
        log.Printf("Validation failed: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  err.Error(),
        })
        return
    }

    // Check if user exists
    exists, err := checkUserExists(user.Email, user.Nickname)
    if err != nil {
        log.Printf("Database error checking user existence: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Internal server error",
        })
        return
    }
    if exists {
        w.WriteHeader(http.StatusConflict)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "User with this email or nickname already exists",
        })
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Password hashing failed: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Failed to process registration",
        })
        return
    }

    // Begin transaction for user creation
    tx, err := database.DB.Begin()
    if err != nil {
        log.Printf("Failed to begin transaction: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Database error",
        })
        return
    }
    defer tx.Rollback()

    // Insert user (let SQLite handle timestamp)
    result, err := tx.Exec(`
        INSERT INTO users (email, nickname, password, age, gender, first_name, last_name)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
        user.Email, user.Nickname, hashedPassword, user.Age, user.Gender, user.FirstName, user.LastName,
    )

    if err != nil {
        log.Printf("Failed to insert user: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Failed to create user",
        })
        return
    }

    userID, err := result.LastInsertId()
    if err != nil {
        log.Printf("Failed to get user ID: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(AuthResponse{
            Status: "error",
            Error:  "Failed to create user",
        })
        return
    }

    // Create session
    sessionToken := uuid.NewString()
    expiresAt := time.Now().Add(24 * time.Hour)

    _, err = tx.Exec(
        "INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
        userID, sessionToken, expiresAt)

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
            Error:  "Failed to complete registration",
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

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(AuthResponse{
        Status: "success",
        Data: map[string]interface{}{
            "token": sessionToken,
        },
        Message: "Registration successful",
    })
}

func validateUserData(user *models.User) error {
    var missingFields []string

    if user.Email == "" || !strings.Contains(user.Email, "@") {
        missingFields = append(missingFields, "valid email")
    }
    if len(user.Password) < 6 {
        missingFields = append(missingFields, "password (min 6 characters)")
    }
    if len(user.Nickname) < 3 {
        missingFields = append(missingFields, "nickname (min 3 characters)")
    }
    if user.FirstName == "" {
        missingFields = append(missingFields, "firstName")
    }
    if user.LastName == "" {
        missingFields = append(missingFields, "lastName")
    }
    if user.Gender == "" {
        missingFields = append(missingFields, "gender")
    }
    if user.Age < 13 {
        missingFields = append(missingFields, "age (must be 13 or older)")
    }

    if len(missingFields) > 0 {
        return fmt.Errorf("missing or invalid fields: %s", strings.Join(missingFields, ", "))
    }
    return nil
}

func checkUserExists(email, nickname string) (bool, error) {
    var exists bool
    err := database.DB.QueryRow(
        "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR nickname = ?)",
        email, nickname,
    ).Scan(&exists)
    return exists, err
}