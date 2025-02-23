package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type Message struct {
	ID         string `json:"id"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdAt"`
}

// GetMessagesHandler handles fetching messages between two users
func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context with nil check
		userId, ok := r.Context().Value("userId").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get other user's ID from URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		otherUserId := parts[3]

		// Query messages
		rows, err := db.Query(`
            SELECT id, sender_id, receiver_id, content, created_at 
            FROM messages 
            WHERE (sender_id = ? AND receiver_id = ?) 
               OR (sender_id = ? AND receiver_id = ?)
            ORDER BY created_at DESC LIMIT 50
        `, userId, otherUserId, otherUserId, userId)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var messages []Message
		for rows.Next() {
			var msg Message
			err := rows.Scan(&msg.ID, &msg.SenderId, &msg.ReceiverId, &msg.Content, &msg.CreatedAt)
			if err != nil {
				continue
			}
			messages = append(messages, msg)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

// SendMessageHandler handles sending a new message
func SendMessageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get sender ID from context with nil check
		senderId, ok := r.Context().Value("userId").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse request body
		var msg struct {
			ReceiverId string `json:"receiverId"`
			Content    string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Insert message
		_, err := db.Exec(`
            INSERT INTO messages (sender_id, receiver_id, content) 
            VALUES (?, ?, ?)
        `, senderId, msg.ReceiverId, msg.Content)

		if err != nil {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
