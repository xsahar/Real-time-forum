package categories

import (
    "encoding/json"
    "net/http"
    
    "Real-time-forum/src/database"
    "Real-time-forum/src/models"
)

// Response struct for category operations
type CategoryResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// GetCategoriesHandler returns all categories
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    rows, err := database.DB.Query("SELECT id, name, description FROM categories")
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var cat models.Category
        if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
            continue
        }
        categories = append(categories, cat)
    }

    json.NewEncoder(w).Encode(CategoryResponse{
        Status: "success",
        Data:   categories,
    })
}

// AddCategoryHandler creates a new category
func AddCategoryHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var category models.Category
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if category.Name == "" {
        http.Error(w, "Category name is required", http.StatusBadRequest)
        return
    }

    result, err := database.DB.Exec(
        "INSERT INTO categories (name, description) VALUES (?, ?)",
        category.Name, category.Description,
    )
    if err != nil {
        http.Error(w, "Failed to create category", http.StatusInternalServerError)
        return
    }

    id, _ := result.LastInsertId()
    category.ID = int(id)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(CategoryResponse{
        Status:  "success",
        Message: "Category created successfully",
        Data:    category,
    })
}

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        GetCategoriesHandler(w, r)
    } else if r.Method == http.MethodPost {
        AddCategoryHandler(w, r)
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}