
package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes database connection
func InitDB() error {
    var err error
    DB, err = sql.Open("sqlite3", "./forum.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }

    // Test connection
    if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }

    // Initialize schema
    if err = CreateTables(); err != nil {
        return fmt.Errorf("failed to create tables: %v", err)
    }

    return nil
}

// CloseDB closes database connection
func CloseDB() error {
    if DB != nil {
        return DB.Close()
    }
    return nil
}