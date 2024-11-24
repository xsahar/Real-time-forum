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
    "Real-time-forum/src/middleware"

    _ "github.com/mattn/go-sqlite3"
)

func setupRoutes() {
    // API routes with middleware configuration
    apiRoutes := []struct {
        path         string
        handler      http.HandlerFunc
        useRateLimit bool
    }{
        {"/api/posts", posts.ListPostsHandler, false},      // List all posts
        {"/api/post/{id}", posts.GetPostHandler, false},    // Get single post
        {"/api/login", auth.LoginHandler, false},
        {"/api/register", auth.RegisterHandler, false},
        {"/api/logout", auth.LogoutHandler, false},
        {"/api/create-post", posts.CreatePostHandler, false},
        {"/api/post/{id}/comment", posts.PostAndCommentHandler, false},
        {"/api/like/{id}", posts.LikeHandler, true},        
        {"/api/dislike/{id}", posts.DislikeHandler, true},  
        {"/api/categories", categories.CategoriesHandler, false},
        {"/api/add-category", categories.AddCategoryHandler, false},
    }

    // Register API routes with middleware
    for _, route := range apiRoutes {
        handler := core.RecoveryHandler(route.handler)
        if route.useRateLimit {
            handler = middleware.RateLimit(handler)
        }
        http.HandleFunc(route.path, handler)
    }

    // SPA route - serve index.html for non-API routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if !strings.HasPrefix(r.URL.Path, "/api/") {
            http.ServeFile(w, r, "templates/index.html")
        }
    })

    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
    // Initialize database
    if err := database.InitDB(); err != nil {
        log.Fatal("Database initialization failed:", err)
    }
    defer database.DB.Close()

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