package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// GetOnlineUsersHandler returns a list of all users and their online status
func GetOnlineUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current user ID from context
		currentUserID, ok := r.Context().Value("userId").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Modified query to exclude current user and get last message timestamp
		rows, err := db.Query(`
            SELECT u.id, u.nickname, u.online, 
                (SELECT MAX(created_at) FROM messages 
                WHERE (sender_id = u.id AND receiver_id = ?) 
                OR (sender_id = ? AND receiver_id = u.id)) as last_message
            FROM users u
            WHERE u.id != ?
        `, currentUserID, currentUserID, currentUserID)
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []map[string]interface{}
		for rows.Next() {
			var id, nickname string
			var online bool
			var lastMessage sql.NullString

			if err := rows.Scan(&id, &nickname, &online, &lastMessage); err != nil {
				continue
			}

			users = append(users, map[string]interface{}{
				"id":              id,
				"nickname":        nickname,
				"online":          online,
				"hasMessages":     lastMessage.Valid,
				"lastMessageTime": lastMessage.String,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
