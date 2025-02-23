package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"forum/src/models"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("Error decoding registration request: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request format"})
			return
		}

		// Validate password
		if len(user.Password) < 6 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Password must be at least 6 characters"})
			return
		}

		// Check if email or nickname already exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR nickname = ?)",
			user.Email, user.Nickname).Scan(&exists)

		if err != nil {
			log.Printf("Error checking existing user: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Error during registration"})
			return
		}

		if exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Email or nickname already exists"})
			return
		}

		// Hash password
		// log.Printf("Registration: Password length before hashing: %d", len(user.Password))
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Error processing password"})
			return
		}
		// log.Printf("Registration: Hashed password length: %d", len(hashedPassword))

		// Generate UUID for user
		user.ID = uuid.New().String()

		// Insert user into database
		_, err = db.Exec(`
            INSERT INTO users (id, nickname, age, gender, first_name, last_name, email, password)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        `, user.ID, user.Nickname, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, hashedPassword)

		if err != nil {
			log.Printf("Error creating user: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Error creating user"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var login LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			log.Printf("Login error: Failed to decode request: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request format"})
			return
		}

		log.Printf("Login attempt for identifier: %s", login.Identifier)

		var user models.User
		var storedPassword string // Change to string to match database storage

		// Query user data
		err := db.QueryRow(`
            SELECT id, nickname, email, first_name, last_name, password 
            FROM users 
            WHERE LOWER(email) = LOWER(?) OR LOWER(nickname) = LOWER(?)
        `, login.Identifier, login.Identifier).Scan(
			&user.ID,
			&user.Nickname,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&storedPassword,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("Login error: User not found: %s", login.Identifier)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"message": "Invalid credentials"})
				return
			}
			log.Printf("Login error: Database error: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Server error"})
			return
		}

		// Debug logging
		log.Printf("Found user: %s", user.Nickname)
		log.Printf("Stored password length: %d", len(storedPassword))
		log.Printf("Input password length: %d", len(login.Password))

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(login.Password))
		if err != nil {
			log.Printf("Login error: Password comparison failed: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid credentials"})
			return
		}

		// Generate new session token
		token := uuid.New().String()
		log.Printf("Login successful for user: %s", user.Nickname)

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthResponse{
			Token: token,
			User:  user,
		})
	}
}

func AuthMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow WebSocket upgrade requests without token in header
		if r.URL.Path == "/api/ws" {
			token := r.URL.Query().Get("token")
			if token == "" {
				http.Error(w, "No token provided", http.StatusUnauthorized)
				return
			}

			var userId string
			err := db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", token).Scan(&userId)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Regular API auth handling
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		var userId string
		err := db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", tokenString).Scan(&userId)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
