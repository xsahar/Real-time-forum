package ws

import (
	// "database/sql"
	"database/sql"
	"time"
)


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