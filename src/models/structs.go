package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    // UserID    int       `json:"user_id"` 
    Nickname  string    `json:"nickname"`
    Password  string    `json:"-"`
    Age       int       `json:"age"`
    Gender    string    `json:"gender"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    CreatedAt time.Time `json:"created_at"`
}

type Post struct {
    ID         int       `json:"id"`
    UserID     int       `json:"user_id"`  // Add this field
    Title      string    `json:"title"`
    Content    string    `json:"content"`
    Author     string    `json:"author"`
    Categories []string  `json:"categories"`
    CreatedAt  time.Time `json:"created_at"`
    Likes      int       `json:"likes"`
    Dislikes   int       `json:"dislikes"`
}

type Comment struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    AuthorID  int       `json:"author_id"`
    Author    string    `json:"author"`
    Content   string    `json:"content"`
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