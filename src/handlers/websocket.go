// src/handlers/websocket.go
package handlers

import (
	"database/sql"
	// "encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	// Store active connections
	clients    = make(map[string]*websocket.Conn)
	clientsMux sync.RWMutex
)

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func WebSocketHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userId, ok := r.Context().Value("userId").(string)
		if !ok {
			log.Printf("WebSocket: no user ID in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("WebSocket: User %s connecting", userId)

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		// Store connection
		clientsMux.Lock()
		clients[userId] = conn
		clientsMux.Unlock()

		// Update user's online status
		_, err = db.Exec("UPDATE users SET online = true WHERE id = ?", userId)
		if err != nil {
			log.Printf("Error updating online status: %v", err)
		} else {
			log.Printf("User %s is now online", userId)
		}

		// Broadcast online status to all clients
		broadcastStatus(db, userId, true)

		// Cleanup on disconnect
		defer func() {
			clientsMux.Lock()
			delete(clients, userId)
			clientsMux.Unlock()

			// Update offline status
			_, err = db.Exec("UPDATE users SET online = false WHERE id = ?", userId)
			if err != nil {
				log.Printf("Error updating offline status: %v", err)
			} else {
				log.Printf("User %s is now offline", userId)
			}

			// Broadcast offline status
			broadcastStatus(db, userId, false)
		}()

		// Message handling loop
		for {
			var msg WSMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}

			log.Printf("Received message from %s: %v", userId, msg.Type)
			handleWSMessage(db, userId, msg)
		}
	}
}

func handleWSMessage(db *sql.DB, senderId string, msg WSMessage) {
	log.Printf("Handling WebSocket message type: %s", msg.Type)

	switch msg.Type {
	case "message":
		// Convert the payload from interface{} to map
		payloadMap, ok := msg.Payload.(map[string]interface{})
		if !ok {
			log.Printf("Invalid message payload format")
			return
		}

		// Extract receiverId and content
		receiverId, ok := payloadMap["receiverId"].(string)
		if !ok {
			log.Printf("Invalid or missing receiverId in message")
			return
		}

		content, ok := payloadMap["content"].(string)
		if !ok {
			log.Printf("Invalid or missing content in message")
			return
		}

		log.Printf("Message from %s to %s: %s", senderId, receiverId, content)

		// Generate a message ID and timestamp
		messageId := uuid.New().String()
		timestamp := time.Now().Format(time.RFC3339)

		// Store the message in the database
		_, err := db.Exec(
			"INSERT INTO messages (id, sender_id, receiver_id, content, created_at, is_read) VALUES (?, ?, ?, ?, ?, ?)",
			messageId, senderId, receiverId, content, timestamp, false,
		)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			return
		}

		// Create response message
		response := WSMessage{
			Type: "message",
			Payload: map[string]interface{}{
				"id":         messageId,
				"senderId":   senderId,
				"receiverId": receiverId,
				"content":    content,
				"createdAt":  timestamp,
				"isRead":     false,
			},
		}

		// Send to receiver if online
		clientsMux.RLock()
		if receiverConn, ok := clients[receiverId]; ok {
			if err := receiverConn.WriteJSON(response); err != nil {
				log.Printf("Error sending message to receiver: %v", err)
			} else {
				log.Printf("Message sent to receiver %s", receiverId)
			}
		} else {
			log.Printf("Receiver %s is not online, message saved for later", receiverId)
		}
		clientsMux.RUnlock()

		// Send confirmation back to sender
		clientsMux.RLock()
		if senderConn, ok := clients[senderId]; ok {
			if err := senderConn.WriteJSON(response); err != nil {
				log.Printf("Error sending confirmation to sender: %v", err)
			} else {
				log.Printf("Message confirmation sent to sender %s", senderId)
			}
		}
		clientsMux.RUnlock()
	}
}

// Broadcast status updates to all connected clients
func broadcastStatus(db *sql.DB, userId string, online bool) {
	// Get the user's nickname
	var nickname string
	err := db.QueryRow("SELECT nickname FROM users WHERE id = ?", userId).Scan(&nickname)
	if err != nil {
		log.Printf("Error getting user nickname: %v", err)
		return
	}

	statusMsg := WSMessage{
		Type: "status",
		Payload: map[string]interface{}{
			"userId":   userId,
			"nickname": nickname,
			"online":   online,
		},
	}

	// Send to all clients
	clientsMux.RLock()
	for _, conn := range clients {
		if err := conn.WriteJSON(statusMsg); err != nil {
			log.Printf("Error broadcasting status: %v", err)
		}
	}
	clientsMux.RUnlock()
}
