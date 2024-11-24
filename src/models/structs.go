package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    Nickname  string    `json:"nickname"`
    Password  string    `json:"password"` 
    Age       int       `json:"age"`
    Gender    string    `json:"gender"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    CreatedAt time.Time `json:"created_at"`
}

type Post struct {
    ID         int       `json:"id"`
    UserID     int       `json:"user_id"`
    Title      string    `json:"title"`
    Content    string    `json:"content"`
    Author     string    `json:"author"`
    Likes      int       `json:"likes"`
    Dislikes   int       `json:"dislikes"`
    Categories []string  `json:"categories"`
    CreatedAt  time.Time `json:"created_at"`
}

type Comment struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    Author    string    `json:"author"`      // Add author nickname
    CreatedAt time.Time `json:"created_at"`
}

type Category struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}

type Like struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    Value     int       `json:"value"` // 1 for like, -1 for dislike
    CreatedAt time.Time `json:"created_at"`
}