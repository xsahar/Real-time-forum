package database

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAdvancedDatabaseOperations(t *testing.T) {
    // Initial setup
    setupTestDB(t)
    
    // Ensure database is created and tables exist
    if err := CreateTables(); err != nil {
        t.Fatalf("Failed to create tables: %v", err)
    }

    tests := []struct {
        name string
        fn   func(*testing.T)
    }{
        {"Database Initialization", func(t *testing.T) {
            // Save original DB
            origDB := DB
            defer func() { DB = origDB }()
            
            // Test nil database
            DB = nil
            if err := InitializeDB(); err == nil {
                t.Error("Expected error for nil database")
            }

            // Test invalid path
            DB, err := sql.Open("sqlite3", "/nonexistent/path/db.db")
            if err == nil {
                if err := DB.Ping(); err == nil {
                    t.Error("Expected error for invalid database path")
                }
            }
        }},

        {"Table Creation", func(t *testing.T) {
            origDB := DB
            defer func() { DB = origDB }()

            // Test nil database
            DB = nil
            if err := CreateTables(); err == nil {
                t.Error("Expected error for nil database")
            }

            // Restore and test invalid SQL
            DB = origDB
            _, err := DB.Exec("CREATE TABLE invalid)")
            if err == nil {
                t.Error("Expected error for invalid SQL")
            }
        }},

        {"Category Operations", func(t *testing.T) {
            // Clear existing data
            if _, err := DB.Exec("DELETE FROM categories"); err != nil {
                t.Fatalf("Failed to clear categories: %v", err)
            }

            // Insert test category
            _, err := DB.Exec(`
                INSERT INTO categories (name, description) 
                VALUES ('General', 'General Discussion')`)
            if err != nil {
                t.Fatalf("Failed to insert test category: %v", err)
            }

            testCases := []struct {
                name    string
                sql     string
                wantErr bool
            }{
                {
                    name: "Duplicate Name",
                    sql: `INSERT INTO categories (name, description) 
                          VALUES ('General', 'Another description')`,
                    wantErr: true,
                },
                {
                    name: "NULL Name",
                    sql: `INSERT INTO categories (description) 
                          VALUES ('Test description')`,
                    wantErr: true,
                },
                {
                    name: "Empty Name",
                    sql: `INSERT INTO categories (name, description)
                          VALUES ('', 'Empty name')`,
                    wantErr: true,
                },
                {
                    name: "Very Long Name",
                    sql: fmt.Sprintf(`
                        INSERT INTO categories (name, description)
                        VALUES ('%s', 'Too long')`, strings.Repeat("a", 256)),
                    wantErr: true,
                },
                {
                    name: "Special Characters",
                    sql: `INSERT INTO categories (name, description)
                          VALUES ('Test''s Category', 'Valid name')`,
                    wantErr: false,
                },
            }

            for _, tc := range testCases {
                t.Run(tc.name, func(t *testing.T) {
                    _, err := DB.Exec(tc.sql)
                    if tc.wantErr && err == nil {
                        t.Errorf("Expected error for %s", tc.name)
                    }
                    if !tc.wantErr && err != nil {
                        t.Errorf("Unexpected error for %s: %v", tc.name, err)
                    }
                })
            }
        }},

        {"Transaction Isolation", func(t *testing.T) {
            // Clear existing data
            if _, err := DB.Exec("DELETE FROM categories"); err != nil {
                t.Fatalf("Failed to clear categories: %v", err)
            }

            // First transaction
            tx1, err := DB.Begin()
            if err != nil {
                t.Fatalf("Failed to begin transaction 1: %v", err)
            }
            defer tx1.Rollback()

            _, err = tx1.Exec(`
                INSERT INTO categories (name, description) 
                VALUES ('Transaction1', 'Test Category')`)
            if err != nil {
                t.Fatalf("Failed insert in tx1: %v", err)
            }

            // Second transaction
            tx2, err := DB.Begin()
            if err != nil {
                t.Fatalf("Failed to begin transaction 2: %v", err)
            }
            defer tx2.Rollback()

            // Try duplicate insert
            _, err = tx2.Exec(`
                INSERT INTO categories (name, description) 
                VALUES ('Transaction1', 'Test Category')`)
            if err == nil {
                t.Error("Expected error for duplicate insert")
            }

            // Commit first transaction
            if err := tx1.Commit(); err != nil {
                t.Fatalf("Failed to commit tx1: %v", err)
            }
        }},
								// Add these test cases to TestAdvancedDatabaseOperations:
				
				{"Error Handling", func(t *testing.T) {
						// Test transaction rollback
						tx, err := DB.Begin()
						if err != nil {
								t.Fatalf("Failed to begin transaction: %v", err)
						}
				
						// Execute invalid SQL
						_, err = tx.Exec("INSERT INTO nonexistent VALUES (1)")
						if err == nil {
								t.Error("Expected error for invalid table")
						}
				
						// Verify rollback works
						if err := tx.Rollback(); err != nil {
								t.Errorf("Rollback failed: %v", err)
						}
				
						// Test double close
						if err := tx.Commit(); err == nil {
								t.Error("Expected error committing rolled back transaction")
						}
				}},
				
								{"Initialization Failures", func(t *testing.T) {
								setupTestDB(t) // Start with fresh DB
								
								origDB := DB
								defer func() { 
										DB = origDB
										setupTestDB(t) // Ensure DB is restored for next test
								}()
								
								if err := DB.Close(); err != nil {
										t.Fatalf("Failed to close database: %v", err)
								}
						
								// Try operations on closed database
								err := DB.Ping()
								if err == nil {
										t.Error("Expected error pinging closed database")
								}
						
								if err := InitializeDB(); err == nil {
										t.Error("Expected error initializing closed database")
								}
						
								if err := CreateTables(); err == nil {
										t.Error("Expected error creating tables on closed database")
								}
						
								if err := InitCategories(); err == nil {
										t.Error("Expected error initializing categories on closed database")
								}
						}},
						
						{"Constraint Testing", func(t *testing.T) {
								setupTestDB(t) // Start with fresh DB
								
								if err := CreateTables(); err != nil {
										t.Fatalf("Failed to create tables: %v", err)
								}
						
								// Insert test user
								_, err := DB.Exec(`
										INSERT INTO users (
												email, nickname, password, age, gender, 
												first_name, last_name
										) VALUES (
												'test@test.com', 'test', 'password123', 
												25, 'M', 'Test', 'User'
										)`)
								if err != nil {
										t.Fatalf("Failed to insert test user: %v", err)
								}
						
								// Test foreign key constraint
								_, err = DB.Exec(`
										INSERT INTO messages (from_id, to_id, content)
										VALUES (999, 999, 'test')`)
								if err == nil {
										t.Error("Expected foreign key constraint error")
								}
						
								// Test unique constraint
								_, err = DB.Exec(`
										INSERT INTO users (
												email, nickname, password, age, gender, 
												first_name, last_name
										) VALUES (
												'test@test.com', 'test2', 'password123', 
												30, 'F', 'Test2', 'User2'
										)`)
								if err == nil {
										t.Error("Expected unique constraint error")
								}
						
								// Test NOT NULL constraint
								_, err = DB.Exec(`
										INSERT INTO users (email, nickname) 
										VALUES ('new@test.com', 'newtest')`)
								if err == nil {
										t.Error("Expected NOT NULL constraint error")
								}
						}},
						
						{"Cleanup Verification", func(t *testing.T) {
								setupTestDB(t) // Start with fresh DB
								
								stats := DB.Stats()
								initialOpen := stats.OpenConnections
						
								for i := 0; i < 10; i++ {
										tx, err := DB.Begin()
										if err != nil {
												t.Fatalf("Failed to begin transaction: %v", err)
										}
										if err := tx.Rollback(); err != nil {
												t.Errorf("Failed to rollback transaction: %v", err)
										}
								}
						
								stats = DB.Stats()
								if stats.OpenConnections > initialOpen {
										t.Errorf("Connection leak detected: %d -> %d", 
												initialOpen, stats.OpenConnections)
								}
						}},

        				{"Concurrent Operations", func(t *testing.T) {
						// Start with fresh DB
						setupTestDB(t)
						
						// Create tables first
						if err := CreateTables(); err != nil {
								t.Fatalf("Failed to create tables: %v", err)
						}
						
						// Clear existing data
						if _, err := DB.Exec("DELETE FROM categories"); err != nil {
								t.Fatalf("Failed to clear categories: %v", err)
						}
				
						var wg sync.WaitGroup
						errs := make(chan error, 10)
						var mu sync.Mutex
						inserted := make(map[string]bool)
				
						// Start concurrent operations
						for i := 0; i < 10; i++ {
								wg.Add(1)
								go func(id int) {
										defer wg.Done()
										
										name := fmt.Sprintf("Category%d", id)
										_, err := DB.Exec(`
												INSERT INTO categories (name, description)
												VALUES (?, 'Concurrent test')`,
												name)
												
										if err != nil {
												select {
												case errs <- fmt.Errorf("insert %s failed: %v", name, err):
												default:
												}
												return
										}
										
										mu.Lock()
										inserted[name] = true
										mu.Unlock()
								}(i)
						}
				
						// Wait for all operations to complete
						done := make(chan struct{})
						go func() {
								wg.Wait()
								close(done)
								close(errs)
						}()
				
						// Check for timeout
						select {
						case <-done:
								// Success - continue
						case <-time.After(5 * time.Second):
								t.Fatal("Concurrent operations timed out")
						}
				
						// Check for errors
						for err := range errs {
								t.Error(err)
						}
				
						// Verify results
						var count int
						err := DB.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
						if err != nil {
								t.Fatalf("Failed to count categories: %v", err)
						}
						if count != 10 {
								t.Errorf("Expected 10 categories, got %d", count)
						}
				
						// Verify each category
						mu.Lock()
						defer mu.Unlock()
						for i := 0; i < 10; i++ {
								name := fmt.Sprintf("Category%d", i)
								if !inserted[name] {
										t.Errorf("Category %s was not inserted", name)
								}
						}
				}},

        {"Database Cleanup", func(t *testing.T) {
            // Test closing nil database
            origDB := DB
            DB = nil
            if err := CloseDB(); err != nil {
                t.Error("Expected nil error when closing nil database")
            }

            // Test closing actual database
            DB = origDB
            if err := CloseDB(); err != nil {
                t.Errorf("Failed to close database: %v", err)
            }

            // Verify database is closed
            _, err := DB.Exec("SELECT 1")
            if err == nil {
                t.Error("Expected error on closed database")
            }

            // Restore connection for next tests
            setupTestDB(t)
        }},
    }

    for _, tt := range tests {
        t.Run(tt.name, tt.fn)
    }
}