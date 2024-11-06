package forum

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// Initialize the database connection and create tables
func InitDB() {
    var err error
    Db, err = sql.Open("sqlite3", "./forum.db")
    if err != nil {
        log.Fatal(err)
    }
    createTables()
    insertCategories()
}

func createTables() {
    schema := `
    CREATE TABLE IF NOT EXISTS user (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL, 
        Nickname TEXT UNIQUE NOT NULL,
        age INTEGER NOT NULL,
        gender TEXT NOT NULL, 
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL 
    );

    CREATE TABLE IF NOT EXISTS post (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user (id)
    );

    CREATE TABLE IF NOT EXISTS category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL,
        user_id INTEGER,
        FOREIGN KEY (user_id) REFERENCES user (id)
    );

    CREATE TABLE IF NOT EXISTS post_category (
        post_id INTEGER NOT NULL,
        category_id INTEGER NOT NULL,
        PRIMARY KEY (post_id, category_id),
        FOREIGN KEY (post_id) REFERENCES post (id),
        FOREIGN KEY (category_id) REFERENCES category (id)
    );

    CREATE TABLE IF NOT EXISTS comments (
        comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        FOREIGN KEY (post_id) REFERENCES post (id),
        FOREIGN KEY (user_id) REFERENCES user (id)
    );

    CREATE TABLE IF NOT EXISTS post_likes_dislikes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        type TEXT CHECK(type IN ('like', 'dislike')) NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user (id),
        FOREIGN KEY (post_id) REFERENCES post (id)
    );

    CREATE TABLE IF NOT EXISTS comment_likes_dislikes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        comment_id INTEGER NOT NULL,
        type TEXT CHECK(type IN ('like', 'dislike')) NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user (id),
        FOREIGN KEY (comment_id) REFERENCES comments (comment_id)
    );

    CREATE TABLE IF NOT EXISTS sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        token TEXT NOT NULL UNIQUE,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        expires_at DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user(id)
    );
    `

    _, err := Db.Exec(schema)
    if err != nil {
        log.Fatal("Error creating tables: ", err)
    }
}

func insertCategories() {
    categories := []string{
        "Science", "Technology", "Sports", "Business", "Blog", "Friendship", 
        "Cooking", "Hobby", "Event", "Health", "Travel", "Education", 
        "Entertainment", "Art", "Music", "Literature", "Fitness", "Gaming",
        "Nature", "History", "Politics", "Finance", "Lifestyle", "Fashion",
        "Photography", "Parenting", "DIY", "Crafts", "Pets", "Food",
    }

    for _, category := range categories {
        var exists bool
        err := Db.QueryRow("SELECT EXISTS(SELECT 1 FROM category WHERE name=?)", category).Scan(&exists)
        if err != nil {
            log.Printf("Error checking category %s : %v", category, err)
            continue
        }

        if !exists {
            _, err := Db.Exec("INSERT INTO category (name) VALUES (?)", category)
            if err != nil {
                log.Printf("Error inserting category %s : %v", category, err)
            }
        }
    }
}