package ws

import (
    // "encoding/json"
    "log"
    "time"
    "github.com/gorilla/websocket"
)


func NewClient(userID int, hub *Hub, conn *websocket.Conn) *Client {
    return &Client{
        UserID: userID,
        Hub:    hub,
        Send:   make(chan Message),
        conn:   conn,
    }
}

func (c *Client) SendPrivateMessage(to int, content string) {
    message := Message{
        Type: MessageTypeChat,
        Data: ChatMessage{
            Content:   content,
            From:     c.Hub.OnlineUsers[c.UserID],
            FromID:   c.UserID,
            Timestamp: time.Now(),
        },
        From: c.UserID,
        To:   to,
    }
    c.Hub.Broadcast <- message
}

func (c *Client) ReadPump() {
    defer func() {
        c.Hub.Unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        var message Message
        err := c.conn.ReadJSON(&message)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        message.From = c.UserID
        c.Hub.Broadcast <- message
    }
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // Hub closed the channel
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            err := c.conn.WriteJSON(message)
            if err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}