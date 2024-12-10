package database

import (
    "database/sql"
    "fmt"
    "log"
    // "time"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
    log.Printf("Initializing database...")
    var err error
    DB, err = sql.Open("sqlite3", "./forum.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }

    if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }

    if _, err = DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
        return fmt.Errorf("failed to enable foreign keys: %v", err)
    }

    if err = CreateTables(); err != nil {
        return fmt.Errorf("failed to create tables: %v", err)
    }

    if err = validateTables(); err != nil {
        return fmt.Errorf("table validation failed: %v", err)
    }

    log.Printf("Database initialized successfully")
    return nil
}

func CreateTables() error {
    if DB == nil {
        return fmt.Errorf("database connection is nil")
    }

    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT NOT NULL UNIQUE,
            nickname TEXT NOT NULL UNIQUE CHECK(length(nickname) >= 3),
            password TEXT NOT NULL CHECK(length(password) >= 6),
            age INTEGER NOT NULL CHECK(age >= 13),
            gender TEXT NOT NULL CHECK(gender IN ('male', 'female')),
            first_name TEXT NOT NULL CHECK(length(first_name) > 0),
            last_name TEXT NOT NULL CHECK(length(last_name) > 0),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE CHECK(length(name) > 0 AND length(name) <= 255),
            description TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,

        `CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            title TEXT NOT NULL CHECK(length(title) > 0),
            content TEXT NOT NULL CHECK(length(content) > 0),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
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
            post_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            content TEXT NOT NULL CHECK(length(content) > 0),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        )`,

        `CREATE TABLE IF NOT EXISTS likes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            post_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            value INTEGER CHECK(value IN (-1,1)),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(post_id, user_id),
            FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        )`,

        `CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            from_id INTEGER NOT NULL,
            to_id INTEGER NOT NULL,
            content TEXT NOT NULL CHECK(length(content) > 0),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (from_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (to_id) REFERENCES users(id) ON DELETE CASCADE
        )`,

        `CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            token TEXT UNIQUE NOT NULL,
            expires_at DATETIME NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        )`,

        // Indexes
        `CREATE INDEX IF NOT EXISTS idx_messages_from_to ON messages(from_id, to_id)`,
        `CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id)`,
        `CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id)`,
        `CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id)`,
        `CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id)`,
        `CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id)`,
        `CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)`,
        `CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token)`,
    }

    for _, query := range queries {
        if _, err := DB.Exec(query); err != nil {
            return fmt.Errorf("failed to execute query: %v\nQuery: %s", err, query)
        }
    }

    if err := InitCategories(); err != nil {
        return fmt.Errorf("failed to initialize categories: %v", err)
    }

    return nil
}

func validateTables() error {
    tables := []string{
        "users", "sessions", "categories", "posts", 
        "post_categories", "comments", "likes", "messages",
    }
    
    for _, table := range tables {
        var count int
        err := DB.QueryRow(
            "SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", 
            table,
        ).Scan(&count)
        
        if err != nil {
            return fmt.Errorf("failed to check table %s: %v", table, err)
        }
        if count == 0 {
            return fmt.Errorf("table %s does not exist", table)
        }

        if table == "categories" {
            var categoryCount int
            err := DB.QueryRow("SELECT COUNT(*) FROM categories").Scan(&categoryCount)
            if err != nil {
                log.Printf("Failed to count categories: %v", err)
            } else {
                log.Printf("Found %d categories in the database", categoryCount)
            }
        }

        log.Printf("Table %s exists", table)
    }
    return nil
}

func InitCategories() error {
    if DB == nil {
        return fmt.Errorf("database connection is nil")
    }

    log.Printf("Initializing categories...")

    var count int
    err := DB.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
    if err != nil {
        return fmt.Errorf("failed to check existing categories: %v", err)
    }
    
    if count > 0 {
        log.Printf("Categories already initialized (%d categories exist)", count)
        return nil
    }

    tx, err := DB.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }
    defer tx.Rollback()

    categories := []string{
        "Science", "Technology", "Sports", "Business", "Blog", "Friendship",
        "Cooking", "Hobby", "Event", "Health", "Travel", "Education",
        "Entertainment", "Art", "Music", "Literature", "Fitness", "Gaming",
        "Nature", "History", "Politics", "Finance", "Lifestyle", "Fashion",
        "Photography", "Parenting", "DIY", "Crafts", "Pets", "Food",
    }

    stmt, err := tx.Prepare(`INSERT OR IGNORE INTO categories (name) VALUES (?)`)
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %v", err)
    }
    defer stmt.Close()

    for _, category := range categories {
        if _, err := stmt.Exec(category); err != nil {
            return fmt.Errorf("failed to insert category %s: %v", category, err)
        }
        log.Printf("Category %s inserted", category)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit categories: %v", err)
    }

    log.Printf("Successfully initialized %d categories", len(categories))
    return nil
}

func CleanupExpiredSessions() error {
    result, err := DB.Exec("DELETE FROM sessions WHERE expires_at < datetime('now')")
    if err != nil {
        return fmt.Errorf("failed to cleanup sessions: %v", err)
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows: %v", err)
    }
    
    if rows > 0 {
        log.Printf("Cleaned up %d expired sessions", rows)
    }
    return nil
}

func CloseDB() error {
    if DB != nil {
        if err := DB.Close(); err != nil {
            return fmt.Errorf("error closing database: %v", err)
        }
        log.Printf("Database connection closed")
    }
    return nil
}