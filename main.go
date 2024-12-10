package main

import (
    "log"
    "net/http"
    "os"
    "strings"

    "Real-time-forum/src/database"
    "Real-time-forum/src/handlers/auth"
    "Real-time-forum/src/handlers/categories"
    "Real-time-forum/src/handlers/core"
    "Real-time-forum/src/handlers/posts"
     postsCrud "Real-time-forum/src/handlers/posts/crud"
    "Real-time-forum/src/handlers/ws"
    "Real-time-forum/src/middleware"

    _ "github.com/mattn/go-sqlite3"
)

func setupCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}

func setupRoutes() {
    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // API routes
    apiRoutes := []struct {
        path         string
        handler      http.HandlerFunc
        useRateLimit bool
    }{
        {"/api/posts", posts.ListPostsHandler, false},
        {"/api/post/{id}", posts.GetPostHandler, false},
        {"/api/login", auth.LoginHandler, false},
        {"/api/register", auth.RegisterHandler, false},
        {"/api/logout", auth.LogoutHandler, false},
        {"/api/create-post", posts.CreatePostHandler, false},
        {"/api/post/{id}/comments", posts.PostAndCommentHandler, false},
        {"/api/like/{id}",postsCrud.LikeHandler, true},
        {"/api/chat/messages", ws.MessageHandler, false},
        {"/api/chat/ws", ws.WebSocketHandler(&ws.Hub{}), false},
        {"/api/categories", categories.CategoriesHandler, false},
    }

    // Apply middleware to API routes
    for _, route := range apiRoutes {
        handler := setupCORS(core.RecoveryHandler(route.handler))
        if route.useRateLimit {
            handler = middleware.RateLimit(handler)
        }
        http.HandleFunc(route.path, handler)
    }

    // SPA fallback - serve index.html for non-API routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if !strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/static/") {
            http.ServeFile(w, r, "templates/index.html")
        }
    })
}

func main() {
    // Initialize database
    if err := database.InitDB(); err != nil {
        log.Fatal("Database initialization failed:", err)
    }
    defer database.DB.Close()

    // Initialize WebSocket hub
    hub := ws.NewHub()
    go hub.Run()

    // Setup routes
    setupRoutes()

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    log.Printf("Server starting on http://localhost:%s\n", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}