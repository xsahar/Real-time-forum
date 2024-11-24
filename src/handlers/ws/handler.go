package ws

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "Real-time-forum/src/handlers/auth"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // In production, should check origin
        return true
    },
}

// WebSocketHandler manages WebSocket connections
func WebSocketHandler(hub *Hub) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Verify user is authenticated
        userID, ok := auth.ValidateSession(w, r)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Upgrade HTTP connection to WebSocket
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Printf("WebSocket upgrade failed: %v", err)
            return
        }

        // Create new client
        client := NewClient(userID, hub, conn)

        // Register client with hub
        hub.Register <- client

        // Start client message pumps
        go client.WritePump()
        go client.ReadPump()
    }
}

// GetClientByUserID finds a client by user ID
func (h *Hub) GetClientByUserID(userID int) *Client {
    h.Mu.RLock()
    defer h.Mu.RUnlock()
    
    for client := range h.Clients {
        if client.UserID == userID {
            return client
        }
    }
    return nil
}