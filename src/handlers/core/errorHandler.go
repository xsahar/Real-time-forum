package core

import (
    "encoding/json"
    "log"
    "net/http"
)

type ErrorResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

func ServerError(w http.ResponseWriter, err error) {
    log.Printf("Server error: %v", err)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusInternalServerError)
    
    json.NewEncoder(w).Encode(ErrorResponse{
        Status:  "error",
        Message: "Internal server error",
    })
}

func NotFoundError(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)
    
    json.NewEncoder(w).Encode(ErrorResponse{
        Status:  "error",
        Message: "Not found",
    })
}

func RecoveryHandler(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic: %v", err)
                
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusInternalServerError)
                
                json.NewEncoder(w).Encode(ErrorResponse{
                    Status:  "error",
                    Message: "Internal server error",
                })
            }
        }()
        next.ServeHTTP(w, r)
    }
}