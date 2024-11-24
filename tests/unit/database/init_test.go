package database_test

import (
    "testing"
    "os"
    "Real-time-forum/src/database"
)

func TestDatabaseInit(t *testing.T) {
	//1. Setup test database
    testDB := "./test_forum.db"
    os.Setenv("DB_PATH", testDB)
    defer os.Remove(testDB) // Cleanup after test

    // Test initialization
    err := database.InitDB()
    if err != nil {
        t.Fatalf("Failed to initialize database: %v", err)
    }
    defer database.CloseDB()

  // 2. Verify table creation
  tables := []string{
    "users", 
    "categories", 
    "posts", 
    "comments", 
    "likes", 
    "post_categories",
}
    for _, table := range tables {
        var name string
        err := database.DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
        if err != nil {
            t.Errorf("Table %s not created", table)
        }
    }

   // 3. Check categories initialization
    var count int
    err = database.DB.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
    if err != nil {
        t.Fatal("Failed to count categories")
    }
    if count == 0 {
        t.Error("Categories not initialized")
    }
}