package categories_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "Real-time-forum/src/database"
    "Real-time-forum/src/handlers/categories"
    "Real-time-forum/src/models"
)

func TestMain(m *testing.M) {
    // Setup test database
    os.Setenv("DB_PATH", ":memory:")
    if err := database.InitDB(); err != nil {
        panic(err)
    }
    
    // Create tables
    if err := database.CreateTables(); err != nil {
        panic(err)
    }

    // Run tests
    code := m.Run()

    // Cleanup
    database.DB.Close()
    os.Exit(code)
}

func TestGetCategories(t *testing.T) {

	   // Clear existing categories first
		 database.DB.Exec("DELETE FROM categories")
    
		 // Insert test data
		 _, err := database.DB.Exec(
				 "INSERT INTO categories (name, description) VALUES (?, ?)",
				 "Test Category", "Test Description",
		 )
		 if err != nil {
				 t.Fatalf("Failed to insert test data: %v", err)
		 }

    // Create request
    req := httptest.NewRequest("GET", "/api/categories", nil)
    rr := httptest.NewRecorder()

    // Call handler
    categories.GetCategoriesHandler(rr, req)

    // Check status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Check response structure
    var response categories.CategoryResponse
    err = json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
        t.Fatalf("Could not decode response: %v", err)
    }

    if response.Status != "success" {
        t.Errorf("expected success status, got %v", response.Status)
    }

    // Verify categories exist in response
    cats, ok := response.Data.([]interface{})
    if !ok {
        t.Fatal("could not parse categories data")
    }

    if len(cats) == 0 {
        t.Error("expected at least one category")
    }
}

func TestAddCategory(t *testing.T) {
    category := models.Category{
        Name:        "New Test Category",
        Description: "New Test Description",
    }
    
    body, _ := json.Marshal(category)
    
    req := httptest.NewRequest("POST", "/api/categories", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    categories.AddCategoryHandler(rr, req)

    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)
    }

    var response categories.CategoryResponse
    err := json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
        t.Fatalf("Could not decode response: %v", err)
    }

    if response.Status != "success" {
        t.Errorf("expected success status, got %v", response.Status)
    }
}