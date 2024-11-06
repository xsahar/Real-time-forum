package main

import (
    "fmt"
    h "Real-time-forum/src/handlers"
    "log"
    "net/http"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Serve the index.html for the main route and other relevant routes
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" || r.URL.Path == "/login" || r.URL.Path == "/register" || r.URL.Path == "/new-post" {
            http.ServeFile(w, r, "templates/index.html")
            return
        }
        h.NotFoundHandler(w)
    })

    // API endpoints
    http.HandleFunc("/api/login", h.RecoveryHandler(h.LoginHandler))
    http.HandleFunc("/api/register", h.RecoveryHandler(h.RegisterHandler))
    http.HandleFunc("/api/logout", h.RecoveryHandler(h.LogoutHandler))
    http.HandleFunc("/api/create-post", h.RecoveryHandler(h.CreatePostHandler))
    http.HandleFunc("/api/post/", h.RecoveryHandler(h.PostAndCommentHandler))
    http.HandleFunc("/api/like/", h.RecoveryHandler(h.LikeHandler))
    http.HandleFunc("/api/dislike/", h.RecoveryHandler(h.DislikeHandler))
    http.HandleFunc("/api/new-post", h.RecoveryHandler(h.CreatePostHandler))
    http.HandleFunc("/api/categories", h.RecoveryHandler(h.CategoriesHandler))
    http.HandleFunc("/api/add-category", h.RecoveryHandler(h.AddCategoryHandler))

    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    message := fmt.Sprintf("Server started at http://localhost:%s/\n", port)
    log.Print(message)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}