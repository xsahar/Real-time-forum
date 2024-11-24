package core_test

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"

    "Real-time-forum/src/handlers/core"
)

func TestServerError(t *testing.T) {
    tests := []struct {
        name           string
        err           error
        expectedCode  int
        expectedError string
    }{
        {
            name:           "Basic server error",
            err:           errors.New("test error"),
            expectedCode:  http.StatusInternalServerError,
            expectedError: "Internal server error",
        },
        {
            name:           "Nil error",
            err:           nil,
            expectedCode:  http.StatusInternalServerError,
            expectedError: "Internal server error",
        },
        {
            name:           "Custom error message",
            err:           errors.New("custom error"),
            expectedCode:  http.StatusInternalServerError,
            expectedError: "Internal server error",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            core.ServerError(w, tt.err)

            // Check status code
            if w.Code != tt.expectedCode {
                t.Errorf("Expected status %d, got %d", tt.expectedCode, w.Code)
            }

            // Check Content-Type header
            if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
                t.Errorf("Expected Content-Type application/json, got %s", contentType)
            }

            var response core.ErrorResponse
            if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
                t.Fatalf("Failed to decode response: %v", err)
            }

            if response.Message != tt.expectedError {
                t.Errorf("Expected message %q, got %q", tt.expectedError, response.Message)
            }
        })
    }
}

func TestNotFoundError(t *testing.T) {
    w := httptest.NewRecorder()
    core.NotFoundError(w)

    // Check status code
    if w.Code != http.StatusNotFound {
        t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
    }

    // Check Content-Type header
    if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
        t.Errorf("Expected Content-Type application/json, got %s", contentType)
    }

    var response core.ErrorResponse
    if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    expectedMessage := "Not found"
    if response.Message != expectedMessage {
        t.Errorf("Expected message %q, got %q", expectedMessage, response.Message)
    }
}

func TestRecoveryHandler(t *testing.T) {
    tests := []struct {
        name         string
        handler     http.HandlerFunc
        expectedCode int 
        checkBody   bool
    }{
        {
            name: "Handles panic",
            handler: func(w http.ResponseWriter, r *http.Request) {
                panic("test panic")
            },
            expectedCode: http.StatusInternalServerError,
            checkBody:    true,
        },
        {
            name: "Normal execution",
            handler: func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
            },
            expectedCode: http.StatusOK,
            checkBody:    false,
        },
        {
            name: "Handles panic with nil",
            handler: func(w http.ResponseWriter, r *http.Request) {
                panic(nil)
            },
            expectedCode: http.StatusInternalServerError,
            checkBody:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            wrapped := core.RecoveryHandler(tt.handler)
            w := httptest.NewRecorder()
            r := httptest.NewRequest("GET", "/test", nil)
            
            wrapped.ServeHTTP(w, r)

            if w.Code != tt.expectedCode {
                t.Errorf("Expected status %d, got %d", tt.expectedCode, w.Code)
            }

            if tt.checkBody {
                if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
                    t.Errorf("Expected Content-Type application/json, got %s", contentType)
                }

                var response core.ErrorResponse
                if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
                    t.Fatalf("Failed to decode response: %v", err)
                }
            }
        })
    }
}