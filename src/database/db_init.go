package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)


func InitializeDB() error {
    var err error
    DB, err = sql.Open("sqlite3", "./forum.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }

    if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }

    if err = CreateTables(); err != nil {
        return fmt.Errorf("failed to create tables: %v", err)
    }

    return nil
}


// CreateTables initializes database schema
func CreateTables() error {

    if DB == nil {
        return fmt.Errorf("database connection is nil")
    }

    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT NOT NULL UNIQUE,
            nickname TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            age INTEGER NOT NULL,
            gender TEXT NOT NULL,
            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
               `CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE CHECK(length(name) > 0 AND length(name) <= 255),
            description TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        `CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        )`,
        `CREATE TABLE IF NOT EXISTS post_categories (
            post_id INTEGER,
            category_id INTEGER,
            PRIMARY KEY (post_id, category_id),
            FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
            FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
        )`,
        `CREATE TABLE IF NOT EXISTS comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER,
            user_id INTEGER,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        )`,
        `CREATE TABLE IF NOT EXISTS likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER,
            user_id INTEGER,
            value INTEGER CHECK(value IN (-1,1)),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(post_id, user_id),
            FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        )`,
        `CREATE TABLE IF NOT EXISTS messages(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        from_id INTEGER NOT NULL,
        to_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (from_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (to_id) REFERENCES users(id) ON DELETE CASCADE
        )`,
        `CREATE INDEX IF NOT EXISTS idx_messages_from_to ON messages(from_id, to_id)`,

        `CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            token TEXT UNIQUE NOT NULL,
            expires_at DATETIME NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id) 
            )`,
    }

    for _, query := range queries {
        if _, err := DB.Exec(query); err != nil {
            return err
        }
    }

    // Initialize categories after tables are created
    if err := InitCategories(); err != nil {
        return err
    }

    return nil
}

// InitCategories ensures all default categories exist
func InitCategories() error {

    if DB == nil {
        return fmt.Errorf("database connection is nil")
    }
    txt, err := DB.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer txt.Rollback()

    categories := []string{
        "Science", "Technology", "Sports", "Business", "Blog", "Friendship", 
        "Cooking", "Hobby", "Event", "Health", "Travel", "Education", 
        "Entertainment", "Art", "Music", "Literature", "Fitness", "Gaming",
        "Nature", "History", "Politics", "Finance", "Lifestyle", "Fashion",
        "Photography", "Parenting", "DIY", "Crafts", "Pets", "Food",
    }

    tx, err := DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(`INSERT OR IGNORE INTO categories (name) VALUES (?)`)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, category := range categories {
        if _, err := stmt.Exec(category); err != nil {
            log.Printf("Error inserting category %s: %v", category, err)
            continue
        }
    }

    return tx.Commit()
}

