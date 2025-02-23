package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
)

// GetOnlineUsersHandler returns a list of all users and their online status
func GetOnlineUsersHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // For now, return all users as offline
        rows, err := db.Query("SELECT id, nickname FROM users")
        if err != nil {
            http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var users []map[string]interface{}
        for rows.Next() {
            var id, nickname string
            if err := rows.Scan(&id, &nickname); err != nil {
                continue
            }
            users = append(users, map[string]interface{}{
                "id": id,
                "nickname": nickname,
                "online": false,
            })
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}