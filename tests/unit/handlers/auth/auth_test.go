package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"Real-time-forum/src/database"
	"Real-time-forum/src/handlers/auth"
	"Real-time-forum/src/models"

	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
    // Setup test DB
    if err := database.InitDB(); err != nil {
        panic(err)
    }
    
    // Create tables including sessions
    _, err := database.DB.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER,
            token TEXT UNIQUE,
            expires_at TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )
    `)
    if err != nil {
        panic(err)
    }
    
    // Run tests
    code := m.Run()
    
    // Cleanup
    database.DB.Exec("DROP TABLE IF EXISTS sessions")
    database.DB.Exec("DROP TABLE IF EXISTS users")
    
    os.Exit(code)
}

type testResponse struct {
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
    Error   string `json:"error,omitempty"`
}

func TestRegistration(t *testing.T) {
    // Cleanup before each test
    database.DB.Exec("DELETE FROM users")
    database.DB.Exec("DELETE FROM sessions")

    testCases := []struct {
        name           string
        input         models.User
        expectedCode  int
        expectedResp  testResponse
    }{
        {
            name: "Valid registration",
            input: models.User{
                Email:     "test@example.com",
                Nickname:  "testuser",
                Password:  "password123",
                Age:      25,
                Gender:   "male",
                FirstName: "Test",
                LastName:  "User",
            },
            expectedCode: http.StatusCreated,
            expectedResp: testResponse{
                Status:  "success",
                Message: "Registration successful",
            },
        },
        {
            name: "Missing required fields",
            input: models.User{
                Email: "test@example.com",
            },
            expectedCode: http.StatusBadRequest,
            expectedResp: testResponse{
                Status: "error",
                Error:  "Missing required fields",
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create request with JSON body
            jsonData, err := json.Marshal(tc.input)
            if err != nil {
                t.Fatalf("Failed to marshal request body: %v", err)
            }

						t.Logf("Request body: %s", string(jsonData))

            // Create request with proper headers
            req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonData))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()

            // Execute handler
            auth.RegisterHandler(w, req)

            // Debug response
            t.Logf("Response Status: %d", w.Code)
            t.Logf("Response Body: %s", w.Body.String())

            // Verify status code
            if w.Code != tc.expectedCode {
                t.Errorf("Expected status %d, got %d", tc.expectedCode, w.Code)
            }

            // Verify response body
            var resp testResponse
            if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
                t.Fatalf("Failed to decode response: %v", err)
            }

            // Compare response
            if resp != tc.expectedResp {
                t.Errorf("Expected response %+v, got %+v", tc.expectedResp, resp)
            }

            // Verify session cookie for successful registration
            if tc.expectedCode == http.StatusCreated {
                cookies := w.Result().Cookies()
                if len(cookies) == 0 {
                    t.Error("No session cookie set")
                } else {
                    found := false
                    for _, cookie := range cookies {
                        if cookie.Name == "session_token" {
                            found = true
                            break
                        }
                    }
                    if !found {
                        t.Error("Session token cookie not found")
                    }
                }
            }
        })
    }
}

func TestLogin(t *testing.T) {
    // Clear and setup test data
    database.DB.Exec("DELETE FROM users")
    database.DB.Exec("DELETE FROM sessions")

    // Setup test user
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    database.DB.Exec(`
        INSERT INTO users (email, nickname, password)
        VALUES (?, ?, ?)`,
        "test@example.com", "testuser", hashedPassword)

    testCases := []struct {
        name           string
        credentials   struct {
            Login    string `json:"login"`
            Password string `json:"password"`
        }
        expectedCode  int
        expectedResp  testResponse
        checkCookie  bool
    }{
        {
            name: "Valid login with email",
            credentials: struct {
                Login    string `json:"login"`
                Password string `json:"password"`
            }{
                Login:    "test@example.com",
                Password: "password123",
            },
            expectedCode: http.StatusOK,
            expectedResp: testResponse{
                Status:  "success",
                Message: "Login successful",
            },
            checkCookie: true,
        },
        {
            name: "Valid login with nickname",
            credentials: struct {
                Login    string `json:"login"`
                Password string `json:"password"`
            }{
                Login:    "testuser",
                Password: "password123",
            },
            expectedCode: http.StatusOK,
            expectedResp: testResponse{
                Status:  "success",
                Message: "Login successful",
            },
            checkCookie: true,
        },
        {
            name: "Invalid password",
            credentials: struct {
                Login    string `json:"login"`
                Password string `json:"password"`
            }{
                Login:    "test@example.com",
                Password: "wrongpassword",
            },
            expectedCode: http.StatusUnauthorized,
            expectedResp: testResponse{
                Status: "error",
                Error:  "Invalid credentials",
            },
            checkCookie: false,
        },
        {
            name: "User not found",
            credentials: struct {
                Login    string `json:"login"`
                Password string `json:"password"`
            }{
                Login:    "nonexistent@example.com",
                Password: "password123",
            },
            expectedCode: http.StatusUnauthorized,
            expectedResp: testResponse{
                Status: "error",
                Error:  "Invalid credentials",
            },
            checkCookie: false,
        },
    }
// Add after test cases definition:
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create request with JSON body
            jsonData, err := json.Marshal(tc.credentials)
            if err != nil {
                t.Fatalf("Failed to marshal credentials: %v", err)
            }

            t.Logf("Request body: %s", string(jsonData))
            
            req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonData))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()

            // Execute handler
            auth.LoginHandler(w, req)

            // Debug response
            t.Logf("Response Status: %d", w.Code)
            t.Logf("Response Body: %s", w.Body.String())

            // Verify status code
            if w.Code != tc.expectedCode {
                t.Errorf("Expected status %d, got %d", tc.expectedCode, w.Code)
            }

            // Verify response body
            var resp testResponse
            if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
                t.Fatalf("Failed to decode response: %v", err)
            }

            if resp != tc.expectedResp {
                t.Errorf("Expected response %+v, got %+v", tc.expectedResp, resp)
            }

            // Check session cookie when expected
            if tc.checkCookie {
                cookies := w.Result().Cookies()
                var sessionCookie *http.Cookie
                for _, cookie := range cookies {
                    if cookie.Name == "session_token" {
                        sessionCookie = cookie
                        break
                    }
                }
                
                if sessionCookie == nil {
                    t.Error("No session cookie set")
                } else if !sessionCookie.HttpOnly {
                    t.Error("Session cookie should be HttpOnly")
                }
            }
        })
    }
}

func TestLogout(t *testing.T) {
    // Setup test data
    database.DB.Exec("DELETE FROM users")
    database.DB.Exec("DELETE FROM sessions")

    // Insert test user
    result, err := database.DB.Exec(`
        INSERT INTO users (email, nickname, password)
        VALUES (?, ?, ?)`,
        "test@example.com", "testuser", "hashedpassword")
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }

    userID, _ := result.LastInsertId()
    sessionToken := "test-session-token"

    // Create test session
    _, err = database.DB.Exec(`
        INSERT INTO sessions (user_id, token, expires_at)
        VALUES (?, ?, ?)`,
        userID, sessionToken, time.Now().Add(24*time.Hour))
    if err != nil {
        t.Fatalf("Failed to create test session: %v", err)
    }

    testCases := []struct {
        name          string
        sessionToken  string
        expectedCode int
        expectedResp testResponse
    }{
        {
            name:         "Successful logout",
            sessionToken: sessionToken,
            expectedCode: http.StatusOK,
            expectedResp: testResponse{
                Status:  "success",
                Message: "Logged out successfully",
            },
        },
        {
            name:         "No session token",
            sessionToken: "",
            expectedCode: http.StatusUnauthorized,
            expectedResp: testResponse{
                Status: "error",
                Error:  "Not logged in",
            },
        },
        {
            name:         "Invalid session token",
            sessionToken: "invalid-token",
            expectedCode: http.StatusUnauthorized,
            expectedResp: testResponse{
                Status: "error",
                Error:  "Not logged in",
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
            if tc.sessionToken != "" {
                req.AddCookie(&http.Cookie{
                    Name:  "session_token",
                    Value: tc.sessionToken,
                })
            }
            w := httptest.NewRecorder()

            // Execute handler
            auth.LogoutHandler(w, req)

            // Debug response
            t.Logf("Response Status: %d", w.Code)
            t.Logf("Response Body: %s", w.Body.String())

            // Verify status code
            if w.Code != tc.expectedCode {
                t.Errorf("Expected status %d, got %d", tc.expectedCode, w.Code)
            }

            // Verify response
            var resp testResponse
            if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
                t.Fatalf("Failed to decode response: %v", err)
            }

            if resp != tc.expectedResp {
                t.Errorf("Expected response %+v, got %+v", tc.expectedResp, resp)
            }

            // For successful logout, verify session was removed
            if tc.expectedCode == http.StatusOK {
                var count int
                err := database.DB.QueryRow("SELECT COUNT(*) FROM sessions WHERE token = ?", 
                    tc.sessionToken).Scan(&count)
                if err != nil {
                    t.Fatalf("Failed to check session: %v", err)
                }
                if count > 0 {
                    t.Error("Session was not removed")
                }
            }
        })
    }
}