// tests/unit/handlers/ws/websocket_test.go
package websocket_test

import (
    "testing"
    // "net/http"
    "net/http/httptest"
    "strings"
    "github.com/gorilla/websocket"
    
    "Real-time-forum/src/handlers/ws"
)

func TestWebSocketHandler(t *testing.T) {
    hub := ws.NewHub()
    server := httptest.NewServer(ws.WebSocketHandler(hub))
    defer server.Close()

    tests := []struct {
        name    string
        message string
        wantErr bool
    }{
        {
            name:    "Valid message",
            message: "Hello",
            wantErr: false,
        },
        {
            name:    "Empty message",
            message: "",
            wantErr: true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
            conn, _, err := websocket.DefaultDialer.Dial(url, nil)
            if err != nil {
                t.Fatalf("Failed to connect: %v", err)
            }
            defer conn.Close()

            err = conn.WriteMessage(websocket.TextMessage, []byte(tc.message))
            if (err != nil) != tc.wantErr {
                t.Errorf("WriteMessage() error = %v, wantErr %v", err, tc.wantErr)
            }
        })
    }
}