// src/middleware/ratelimit.go
package middleware

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.Mutex
    window   time.Duration
    limit    int
}

var defaultLimiter *RateLimiter

func InitRateLimiter(window time.Duration, limit int) {
    defaultLimiter = &RateLimiter{
        requests: make(map[string][]time.Time),
        window:   window,
        limit:    limit,
    }
}

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if defaultLimiter == nil {
            // Use default values if not initialized
            InitRateLimiter(time.Minute, 10)
        }

        userID, ok := r.Context().Value("user_id").(int)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        userKey := fmt.Sprintf("%d", userID)
        now := time.Now()

        defaultLimiter.mu.Lock()
        defer defaultLimiter.mu.Unlock()

        // Clean old requests
        if requests, exists := defaultLimiter.requests[userKey]; exists {
            var valid []time.Time
            for _, t := range requests {
                if now.Sub(t) <= defaultLimiter.window {
                    valid = append(valid, t)
                }
            }
            defaultLimiter.requests[userKey] = valid
        }

        // Check rate limit
        if len(defaultLimiter.requests[userKey]) >= defaultLimiter.limit {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusTooManyRequests)
            json.NewEncoder(w).Encode(map[string]string{
                "status":  "error",
                "message": "Rate limit exceeded. Please try again later.",
            })
            return
        }

        // Add new request
        defaultLimiter.requests[userKey] = append(defaultLimiter.requests[userKey], now)

        next.ServeHTTP(w, r)
    }
}