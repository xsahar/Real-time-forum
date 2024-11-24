package database

import (
    "database/sql"
    "os"
    "testing"
    _ "github.com/mattn/go-sqlite3"
)

var testDB *sql.DB

func setupTestDB(t *testing.T) {
    dbPath := "test_forum.db"
    var err error
    testDB, err = sql.Open("sqlite3", dbPath)
    if err != nil {
        t.Fatalf("Failed to open database: %v", err)
    }

    _, err = testDB.Exec("PRAGMA foreign_keys = ON")
    if err != nil {
        t.Fatalf("Failed to enable foreign keys: %v", err)
    }

    DB = testDB

    t.Cleanup(func() {
        testDB.Close()
        os.Remove(dbPath)
    })
}

func TestDatabaseInit(t *testing.T) {
    if err := InitDB(); err != nil {
        t.Fatalf("Failed to initialize database: %v", err)
    }
    
    if DB == nil {
        t.Error("Database not initialized")
    }
    
    if err := DB.Ping(); err != nil {
        t.Errorf("Database connection not working: %v", err)
    }
}

func TestBasicDatabaseOperations(t *testing.T) {
    setupTestDB(t)

    tests := []struct {
        name string
        fn   func(*testing.T)
    }{
        {"Create Tables", func(t *testing.T) {
            err := CreateTables()
            if err != nil {
                t.Fatalf("Failed to create tables: %v", err)
            }
        }},
        {"Verify Tables", func(t *testing.T) {
            tables := []string{
                "users", "categories", "posts", 
                "post_categories", "comments", 
                "likes", "messages",
            }

            for _, table := range tables {
                var name string
                err := testDB.QueryRow(`
                    SELECT name FROM sqlite_master 
                    WHERE type='table' AND name=?`, table).Scan(&name)
                if err != nil {
                    t.Errorf("Table %s not found: %v", table, err)
                }
            }
        }},
        {"Initialize Categories", func(t *testing.T) {
            err := InitCategories()
            if err != nil {
                t.Fatalf("Failed to initialize categories: %v", err)
            }
        }},
        {"Verify Categories", func(t *testing.T) {
            var count int
            err := testDB.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
            if err != nil {
                t.Fatalf("Failed to count categories: %v", err)
            }
            if count != 30 {
                t.Errorf("Expected 30 categories, got %d", count)
            }
        }},
    }

    for _, tt := range tests {
        t.Run(tt.name, tt.fn)
    }
}