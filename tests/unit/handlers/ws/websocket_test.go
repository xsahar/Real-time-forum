// websocket_test.go
package ws_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"Real-time-forum/src/handlers/ws"

	"github.com/gorilla/websocket"
)

func TestNewHub(t *testing.T) {
    hub := ws.NewHub()
    if hub.Clients == nil {
        t.Error("Expected clients map to be initialized")
    }
    if hub.Broadcast == nil {
        t.Error("Expected broadcast channel to be initialized")
    }
}
func TestClientLifecycle(t *testing.T) {
    hub := ws.NewHub()
    go hub.Run()

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        wsHandler := ws.WebSocketHandler(hub)
        wsHandler.ServeHTTP(w, r)
    }))
    defer server.Close()

    t.Setenv("TEST_MODE", "true")

    url := "ws" + strings.TrimPrefix(server.URL, "http")
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    time.Sleep(100 * time.Millisecond)

    clientCount := len(hub.Clients)
    if clientCount != 1 {
        t.Errorf("Expected 1 client, got %d", clientCount)
    }
}

func TestMessageBroadcast(t *testing.T) {
    hub := ws.NewHub()
    go hub.Run()

    msg := ws.Message{
        Type: ws.MessageTypeChat,
        Data: ws.ChatMessage{
            Content: "Test message",
            FromID:  1,
        },
        From: 1,
        To:   0,
    }

    // Create test client
    client := ws.NewClient(1, hub, nil)
    hub.Register <- client

    // Send message
    hub.Broadcast <- msg

    // Verify message received
    select {
    case received := <-client.Send:
        if received.Type != ws.MessageTypeChat {
            t.Errorf("Expected message type %s, got %s", ws.MessageTypeChat, received.Type)
        }
    case <-time.After(time.Second):
        t.Error("Timeout waiting for message")
    }
}

// websocket_test.go
func TestMessageStore(t *testing.T) {
    // Create in-memory SQLite database
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }
    defer db.Close()

    // Create messages table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            from_id INTEGER,
            to_id INTEGER,
            content TEXT,
            created_at TIMESTAMP,
            FOREIGN KEY (from_id) REFERENCES users(id),
            FOREIGN KEY (to_id) REFERENCES users(id)
        )
    `)
    if err != nil {
        t.Fatalf("Failed to create messages table: %v", err)
    }

    // Initialize message store
    store := ws.NewMessageStore(db)

    // Test message
    msg := ws.ChatMessage{
        Content:   "Test message",
        From:     "user1",
        FromID:   1,
        Timestamp: time.Now(),
    }

    // Test save
    err = store.SaveMessage(1, 2, msg.Content)
    if err != nil {
        t.Errorf("Failed to save message: %v", err)
    }

    // Test retrieval
    messages, err := store.GetMessages(1, 2, 10, 0)
    if err != nil {
        t.Errorf("Failed to get messages: %v", err)
    }

    if len(messages) != 1 {
        t.Errorf("Expected 1 message, got %d", len(messages))
    }
}

func TestClientSendPrivateMessage(t *testing.T) {
    hub := ws.NewHub()
    go hub.Run()

    client1 := ws.NewClient(1, hub, nil)
    client2 := ws.NewClient(2, hub, nil)

    hub. Register <- client1
    hub. Register <- client2

    client1.SendPrivateMessage(2, "Test private message")

    select {
    case msg := <-client2.Send:
        if msg.To != 2 {
            t.Errorf("Expected message to user 2, got %d", msg.To)
        }
    case <-time.After(time.Second):
        t.Error("Timeout waiting for private message")
    }
}