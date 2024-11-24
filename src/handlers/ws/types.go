// types.go
package ws

import (
	"database/sql"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Message Types
const (
    MessageTypeChat   = "chat"
    MessageTypeStatus = "status"
    MessageTypeTyping = "typing"
)

// WebSocket timing constants
const (
    writeWait  = 10 * time.Second
    pongWait   = 60 * time.Second
    pingPeriod = (pongWait * 9) / 10
)

// ChatMessage represents a single chat message
type ChatMessage struct {
    Content   string    `json:"content"`
    From     string    `json:"from"`      // Nickname
    FromID   int       `json:"from_id"`   // User ID
    Timestamp time.Time `json:"timestamp"`
}

// Message represents a WebSocket message
type Message struct {
    Type string      `json:"type"`    // Type of message (chat/status/typing)
    Data interface{} `json:"data"`    // Actual message content
    From int         `json:"from"`    // Who sent the message
    To   int        `json:"to"`       // Who should receive it
}

// Client represents a single WebSocket client
type Client struct {
    UserID int              
    Hub    *Hub             
    Send   chan Message     
    conn   *websocket.Conn  
}

// Hub maintains the set of active clients
type Hub struct {
    Clients     map[*Client]bool  
    Broadcast   chan Message      
    Register    chan *Client      
    Unregister  chan *Client     
    OnlineUsers map[int]string    
    Mu          sync.RWMutex      
}

// MessageStore handles message persistence
type MessageStore struct {
    DB *sql.DB
}