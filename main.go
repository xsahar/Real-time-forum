package main

import (
	"forum/src/database"
	"forum/src/handlers"
	"log"
	"net/http"
	"strings"
	"time"
)

// middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// middleware for CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize database
	db, err := database.InitDB("forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ensure online status column exists
	if err := database.EnsureOnlineColumn(db); err != nil {
		log.Printf("Warning: Failed to ensure online column exists: %v", err)
	} else {
		log.Printf("Database schema verified successfully")
	}

	// Create router
	mux := http.NewServeMux()

	// API routes
	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/")
		switch {
		// Auth routes
		case path == "auth/register":
			handlers.RegisterHandler(db).ServeHTTP(w, r)
		case path == "auth/login":
			handlers.LoginHandler(db).ServeHTTP(w, r)

		// Posts routes
		case strings.HasPrefix(path, "posts"):
			parts := strings.Split(path, "/")
			switch {
			case len(parts) == 1:
				switch r.Method {
				case http.MethodGet:
					handlers.GetPostsHandler(db).ServeHTTP(w, r)
				case http.MethodPost:
					handlers.CreatePostHandler(db).ServeHTTP(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			case len(parts) == 2:
				if r.Method != http.MethodGet {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					return
				}
				handlers.GetPostByIDHandler(db).ServeHTTP(w, r)
			case len(parts) == 3 && parts[2] == "comments":
				switch r.Method {
				case http.MethodGet:
					handlers.GetCommentsHandler(db).ServeHTTP(w, r)
				case http.MethodPost:
					handlers.CreateCommentHandler(db).ServeHTTP(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			default:
				http.NotFound(w, r)
			}
		// Categories route
		case path == "categories":
			handlers.GetCategoriesHandler(w, r)

		// Online users route
	case path == "users/online":
		handlers.AuthMiddleware(db, handlers.GetOnlineUsersHandler(db)).ServeHTTP(w, r)

		// Messages routes
		case strings.HasPrefix(path, "messages"):
			// Ensure authentication middleware is applied
			authHandler := handlers.AuthMiddleware(db, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					handlers.GetMessagesHandler(db).ServeHTTP(w, r)
				case http.MethodPost:
					handlers.SendMessageHandler(db).ServeHTTP(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			}))
			authHandler.ServeHTTP(w, r)

		// WebSocket route
		case path == "ws":
			// Apply auth middleware to WebSocket handler
			handlers.AuthMiddleware(db, handlers.WebSocketHandler(db)).ServeHTTP(w, r)

		default:
			http.NotFound(w, r)
		}
	})

	// Apply middleware to API routes
	mux.Handle("/api/", loggingMiddleware(corsMiddleware(apiHandler)))

	// Serve static files
	fs := http.FileServer(http.Dir("src/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html for all other routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/posts") ||
			strings.HasPrefix(r.URL.Path, "/login") || strings.HasPrefix(r.URL.Path, "/register") {
			http.ServeFile(w, r, "src/templates/index.html")
			return
		}
		http.NotFound(w, r)
	})

	// Start server
	log.Printf("Server starting on http://localhost:8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}
