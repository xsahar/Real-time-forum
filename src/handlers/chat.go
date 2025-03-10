package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// GetOnlineUsersHandler returns a list of all users and their online status
func GetOnlineUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, nickname, online FROM users")
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []map[string]interface{}
		for rows.Next() {
			var id, nickname string
			var online bool
			if err := rows.Scan(&id, &nickname, &online); err != nil {
				continue
			}
			users = append(users, map[string]interface{}{
				"id":       id,
				"nickname": nickname,
				"online":   online,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
