package ws

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Handle GET request for message history
        w.Header().Set("Content-Type", "application/json")
        
        user1ID := r.URL.Query().Get("user1")
        user2ID := r.URL.Query().Get("user2")
        
        // Example usage of user1ID and user2ID
        if user1ID == "" || user2ID == "" {
            http.Error(w, "Missing user IDs", http.StatusBadRequest)
            return
        }
        
        // Implementation for getting messages
        
    } else if r.Method == "POST" {
        // Handle POST request for new message
        w.Header().Set("Content-Type", "application/json")
        
        var msg Message
        if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        // Implementation for saving message
    }

}
func (ms *MessageStore) SaveMessage(fromID, toID int, content string) error {
    _, err := ms.DB.Exec(`
        INSERT INTO messages (from_id, to_id, content, created_at)
        VALUES (?, ?, ?, ?)`,
        fromID, toID, content, time.Now())
    return err
}

func (ms *MessageStore) GetMessages(user1ID, user2ID int, limit, offset int) ([]ChatMessage, error) {
    rows, err := ms.DB.Query(`
        SELECT m.content, u.nickname, m.from_id, m.created_at
        FROM messages m
        JOIN users u ON m.from_id = u.id
        WHERE (m.from_id = ? AND m.to_id = ?) 
           OR (m.from_id = ? AND m.to_id = ?)
        ORDER BY m.created_at DESC
        LIMIT ? OFFSET ?`,
        user1ID, user2ID, user2ID, user1ID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []ChatMessage
    for rows.Next() {
        var msg ChatMessage
        err := rows.Scan(&msg.Content, &msg.From, &msg.FromID, &msg.Timestamp)
        if err != nil {
            continue
        }
        messages = append(messages, msg)
    }
    return messages, nil
}

// message_store.go
func NewMessageStore(db *sql.DB) *MessageStore {
    return &MessageStore{
        DB: db,  // Note: DB is now exported
    }
}