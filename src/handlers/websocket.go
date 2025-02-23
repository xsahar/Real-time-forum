package handlers

// import (
// 	"database/sql"
// 	// "encoding/json"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// var (
// 	upgrader = websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin:     func(r *http.Request) bool { return true },
// 	}

// 	// Store active connections
// 	clients    = make(map[string]*websocket.Conn)
// 	clientsMux sync.RWMutex
// )

// type WSMessage struct {
// 	Type    string      `json:"type"`
// 	Payload interface{} `json:"payload"`
// }

// func WebSocketHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userId, ok := r.Context().Value("userId").(string)
// 		if !ok {
// 			log.Printf("WebSocket: no user ID in context")
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// Upgrade HTTP connection to WebSocket
// 		conn, err := upgrader.Upgrade(w, r, nil)
// 		if err != nil {
// 			log.Printf("WebSocket upgrade error: %v", err)
// 			return
// 		}
// 		defer conn.Close()

// 		// Store connection
// 		clientsMux.Lock()
// 		clients[userId] = conn
// 		clientsMux.Unlock()

// 		// Update user's online status
// 		db.Exec("UPDATE users SET online = true WHERE id = ?", userId)
// 		broadcastUserStatus(userId, true)

// 		// Clean up on disconnect
// 		defer func() {
// 			clientsMux.Lock()
// 			delete(clients, userId)
// 			clientsMux.Unlock()
// 			db.Exec("UPDATE users SET online = false WHERE id = ?", userId)
// 			broadcastUserStatus(userId, false)
// 		}()

// 		// Message handling loop
// 		for {
// 			var msg WSMessage
// 			err := conn.ReadJSON(&msg)
// 			if err != nil {
// 				break
// 			}

// 			handleWSMessage(db, userId, msg)
// 		}
// 	}
// }

// func handleWSMessage(db *sql.DB, senderId string, msg WSMessage) {
// 	switch msg.Type {
// 	case "message":
// 		if payload, ok := msg.Payload.(map[string]interface{}); ok {
// 			receiverId := payload["receiverId"].(string)
// 			content := payload["content"].(string)

// 			// Store message in database
// 			_, err := db.Exec(`
//                 INSERT INTO messages (sender_id, receiver_id, content)
//                 VALUES (?, ?, ?)
//             `, senderId, receiverId, content)

// 			if err != nil {
// 				log.Printf("Error storing message: %v", err)
// 				return
// 			}

// 			// Forward message to receiver if online
// 			clientsMux.RLock()
// 			if conn, ok := clients[receiverId]; ok {
// 				conn.WriteJSON(WSMessage{
// 					Type: "message",
// 					Payload: map[string]interface{}{
// 						"senderId": senderId,
// 						"content":  content,
// 						"time":     time.Now(),
// 					},
// 				})
// 			}
// 			clientsMux.RUnlock()
// 		}
// 	}
// }

// func broadcastUserStatus(userId string, online bool) {
// 	status := WSMessage{
// 		Type: "status",
// 		Payload: map[string]interface{}{
// 			"userId": userId,
// 			"online": online,
// 		},
// 	}

// 	clientsMux.RLock()
// 	for _, conn := range clients {
// 		conn.WriteJSON(status)
// 	}
// 	clientsMux.RUnlock()
// }
