package forum

import (
	t "Real-time-forum/src/database"
	// s "Real-time-forum/src/structs"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


func ValidateSession(w http.ResponseWriter, r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 0, false
	}

	sessionToken := cookie.Value
	var userID int
	err = t.Db.QueryRow("SELECT user_id FROM sessions WHERE token = ? AND expires_at > ?", sessionToken, time.Now()).Scan(&userID)
	if err != nil {
		// If session is invalid, clear the cookie and redirect to login
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now().Add(-1 * time.Hour),
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 0, false
	}

	return userID, true
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
			nickname := r.FormValue("nickname")
			email := r.FormValue("email")
			password := r.FormValue("password")
			var identifier string

			if strings.Contains(email, "@") {
					identifier = email
			} else {
					identifier = nickname
			}

			var dbPassword string
			var userID int

			err := t.Db.QueryRow("SELECT id, password FROM user WHERE uemail = ? OR nickname = ?", identifier).Scan(&userID, &dbPassword)
			if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "Invalid username or password"})
					return
			}

			err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
			if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "Invalid username or password"})
					return
			}

			sessionToken := uuid.New().String()
			expiresAt := time.Now().Add(24 * time.Hour)

			// Delete any existing sessions for this user
			_, err = t.Db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
			if err != nil {
					log.Println("Login error: Failed to delete previous sessions:", err)
					http.Error(w, "Error managing session", http.StatusInternalServerError)
					return
			}

			// Create a new session
			_, err = t.Db.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)", userID, sessionToken, expiresAt)
			if err != nil {
					log.Println("Login error: Session creation failed:", err)
					http.Error(w, "Error creating session", http.StatusInternalServerError)
					return
			}

			http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   sessionToken,
					Expires: expiresAt,
			})

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"redirect": "/"})
			return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
			email := r.FormValue("email")
			username := r.FormValue("username")
			password := r.FormValue("password")
			nickname := r.FormValue("nickname")
			ageStr := r.FormValue("age")
			gender := r.FormValue("gender")
			firstname := r.FormValue("first_name")
			lastname := r.FormValue("last_name")

			var errors []string

			// remove leading and trailing spaces to check if string is empty
			username = strings.TrimSpace(username)
			password = strings.TrimSpace(password)
			nickname = strings.TrimSpace(nickname)
			email = strings.TrimSpace(email)
			gender = strings.TrimSpace(gender)
			firstname = strings.TrimSpace(firstname)
			lastname = strings.TrimSpace(lastname)

			// Convert age from string to integer
			age, err := strconv.Atoi(ageStr)
			if err != nil {
					errors = append(errors, "Age must be a valid number")
			}

			// Validation checks
			if email == "" || username == "" || password == "" || nickname == "" || ageStr == "" || gender == "" || firstname == "" || lastname == "" {
					errors = append(errors, "All fields are required")
			}
			if len(email) < 15 {
					errors = append(errors, "Email must be at least 15 characters")
			}
			if len(username) < 5 {
					errors = append(errors, "Username must be at least 5 characters")
			}
			if len(password) < 8 {
					errors = append(errors, "Password must be at least 8 characters")
			}
			if gender != "male" && gender != "female" {
					errors = append(errors, "Choose one")
			}

			if len(errors) > 0 {
					data := struct {
							Errors   []string
							Email    string
							Username string
							Nickname string
							Age      string
							Gender   string
							FirstName string
							LastName  string
					}{
							Errors:   errors,
							Email:    email,
							Username: username,
							Nickname: nickname,
							Age:      ageStr,
							Gender:   gender,
							FirstName: firstname,
							LastName:  lastname,
					}
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(data)
					return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
			}

			// Insert into the database
			result, err := t.Db.Exec(`INSERT INTO user (email, username, password, nickname, age, gender, first_name, last_name ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, email, username, hashedPassword, nickname, age, gender, firstname, lastname)
			// Error Checking
			if err != nil {
					if err.Error() == "UNIQUE constraint failed: user.email" {
							errors = append(errors, "Email already exists")
					} else if err.Error() == "UNIQUE constraint failed: user.username" {
							errors = append(errors, "Username already exists")
					} else if err.Error() == "UNIQUE constraint failed: user.nickname" {
							errors = append(errors, "Nickname already exists")
					} else {
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							log.Println("Error inserting user:", err)
							return
					}
					data := struct {
							Errors   []string
							Email    string
							Username string
							Nickname string
							Age      string
							Gender   string
							FirstName string
							LastName  string
					}{
							Errors:   errors,
							Email:    email,
							Username: username,
							Nickname: nickname,
							Age:      ageStr,
							Gender:   gender,
							FirstName: firstname,
							LastName:  lastname,
					}
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(data)
					return
			}

			userID, err := result.LastInsertId()
			if err != nil {
					http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
					return
			}

			sessionToken := uuid.New().String()
			_, err = t.Db.Exec("INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)", userID, sessionToken, time.Now().Add(24*time.Hour))
			if err != nil {
					http.Error(w, "Error creating session", http.StatusInternalServerError)
					return
			}

			http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Value:   sessionToken,
					Expires: time.Now().Add(24 * time.Hour),
			})

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"redirect": "/"})
			return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sessionToken := cookie.Value
	_, err = t.Db.Exec("DELETE FROM sessions WHERE token = ?", sessionToken)
	if err != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}