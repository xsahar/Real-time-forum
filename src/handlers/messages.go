package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

		// Extract other user's ID correctly from URL path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		otherUserId := parts[3]

		// Add pagination support
		limit := 10
		offset := 0

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
				limit = val
			}
		}

		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
				offset = val
			}
		}

		log.Printf("Fetching messages between %s and %s with limit: %d, offset: %d",
			userId, otherUserId, limit, offset)

		// Initialize empty message array
		var messages []map[string]interface{} = []map[string]interface{}{}

		// Check if the messages table exists and has records
		var tableExists int
		err := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='messages'").Scan(&tableExists)
		if err != nil {
			log.Printf("Error checking if messages table exists: %v", err)
			// Return empty array
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(messages)
			return
		}

		if tableExists == 0 {
			log.Printf("Messages table does not exist")
			// Return empty array
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(messages)
			return
		}
// Query messages with pagination - ORDER BY created_at DESC to get newest first
rows, err := db.Query(`
    SELECT id, sender_id, receiver_id, content, created_at, is_read 
    FROM messages 
    WHERE (sender_id = ? AND receiver_id = ?) 
       OR (sender_id = ? AND receiver_id = ?)
    ORDER BY created_at DESC  
    LIMIT ? OFFSET ?
`, userId, otherUserId, otherUserId, userId, limit, offset)

		if err != nil {
			log.Printf("Database error fetching messages: %v", err)
			// Return empty array instead of error
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(messages)
			return
		}
		defer rows.Close()

		// Process each row carefully
		for rows.Next() {
			// Use nullable types to handle NULL values
			var id, senderId, receiverId, content, createdAt sql.NullString
			var isRead sql.NullBool

			err := rows.Scan(&id, &senderId, &receiverId, &content, &createdAt, &isRead)
			if err != nil {
				log.Printf("Error scanning message row: %v", err)
				continue
			}

			// Convert nullable types to regular types with default values
			idStr := ""
			if id.Valid {
				idStr = id.String
			}

			senderIdStr := ""
			if senderId.Valid {
				senderIdStr = senderId.String
			}

			receiverIdStr := ""
			if receiverId.Valid {
				receiverIdStr = receiverId.String
			}

			contentStr := ""
			if content.Valid {
				contentStr = content.String
			}

			createdAtStr := ""
			if createdAt.Valid {
				createdAtStr = createdAt.String
			}

			isReadBool := false
			if isRead.Valid {
				isReadBool = isRead.Bool
			}

			messages = append(messages, map[string]interface{}{
				"id":         idStr,
				"senderId":   senderIdStr,
				"receiverId": receiverIdStr,
				"content":    contentStr,
				"createdAt":  createdAtStr,
				"isRead":     isReadBool,
			})
		}

		// Check for any error during iteration
		if err = rows.Err(); err != nil {
			log.Printf("Error iterating rows: %v", err)
			// Return partial results if we have them, otherwise empty array
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(messages)
			return
		}

		// Mark messages as read when fetched (only if we have valid messages)
		if len(messages) > 0 {
			_, err = db.Exec(`
					UPDATE messages 
					SET is_read = true 
					WHERE sender_id = ? AND receiver_id = ? AND is_read = false
				`, otherUserId, userId)

			// Don't return error if the update fails - still return messages
			if err != nil {
				log.Printf("Error marking messages as read: %v", err)
			}
		}

		// Always return a JSON array, even if empty
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
