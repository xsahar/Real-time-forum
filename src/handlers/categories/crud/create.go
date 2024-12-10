package crud

import (
	"Real-time-forum/src/database"
	"Real-time-forum/src/models"
	// "Real-time-forum/src/utils"
	"Real-time-forum/src/utils/response"
	"encoding/json"
	"net/http"
)

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    if r.Method != http.MethodPost {
        response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var category models.Category
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        response.SendError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    result, err := database.CreateCategory(category)
    if err != nil {
        response.SendError(w, "Failed to create category", http.StatusInternalServerError)
        return
    }

    response.SendResponse(w, response.APIResponse{
        Status:  "success",
        Message: "Category created successfully",
        Data:    result,
    })
}