package auth

// import (
//     "Real-time-forum/src/database"
//     "Real-time-forum/src/models"
// )

type AuthResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}