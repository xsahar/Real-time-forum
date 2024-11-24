package ws

// import (
// 	"sync"
// 	"time"
// )


// NewHub creates a new hub instance
func NewHub() *Hub {
    return &Hub{
        Clients:    make(map[*Client]bool),     // Empty map for new clients
        Broadcast:  make(chan Message),         // Channel for messages
        Register:   make(chan *Client),         // Channel for new connections
        Unregister: make(chan *Client),         // Channel for disconnections
    }
}

// Run starts the hub and handles all client operations in an infinite loop
func (h *Hub) Run() {
    for {
        // Select statement to handle multiple channels concurrently
        select {
        // Handle new client connections
        case client := <-h.Register:
            h.Mu.Lock()                         // Lock to prevent concurrent map access
            h.Clients[client] = true            // Add new client to clients map
            h.Mu.Unlock()                       // Release the lock
            h.broadcastStatus(client.UserID, "online")  // Notify all clients about new user

        // Handle client disconnections
        case client := <-h.Unregister:
            h.Mu.Lock()                         // Lock for thread-safe map access
            if _, ok := h.Clients[client]; ok { // Check if client exists
                delete(h.Clients, client)        // Remove client from map
                close(client.Send)              // Close client's message channel
                h.broadcastStatus(client.UserID, "offline") // Notify others about disconnection
            }
            h.Mu.Unlock()                       // Release the lock

        // Handle message broadcasting
        case message := <-h.Broadcast:
            h.Mu.RLock()                       // Read lock for concurrent read access
            // Iterate through all connected clients
            for client := range h.Clients {
                // Send to specific client (message.To != 0) or broadcast to all (message.To == 0)
                if message.To == 0 || message.To == client.UserID {
                    select {
                    case client.Send <- message:  // Try to send message to client
                    default:
                        // If client's channel is full/blocked, disconnect client
                        close(client.Send)
                        delete(h.Clients, client)
                    }
                }
            }
            h.Mu.RUnlock()                      // Release the read lock
        }
    }
}

// broadcastStatus creates and sends a status update message to all clients
func (h *Hub) broadcastStatus(userID int, status string) {
    // Create a new status message
    h.Broadcast <- Message{
        Type: "status",                         // Message type for status updates
        Data: map[string]string{"status": status}, // Status data (online/offline)
        From: userID,                           // User who changed status
        To:   0,                                // 0 means broadcast to all clients
    }
}

func (h *Hub) GetOnlineUsers()map[string]interface{}{
    h.Mu.RLock()
    defer h.Mu.RUnlock()
    
users := make(map[string]interface{})
for id, nickname := range h.OnlineUsers {
    users[nickname] = map[string]interface{}{
        "id": id,
        "nickname": nickname,
        "status": "online",
    }
}
    return users
}