package crud

import (
	"Real-time-forum/src/database"
	"Real-time-forum/src/models"
	"Real-time-forum/src/utils/response"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    if r.Method != http.MethodPut {
        response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    parts := strings.Split(r.URL.Path, "/")
    id, err := strconv.Atoi(parts[len(parts)-1])
    if err != nil {
        response.SendError(w, "Invalid category ID", http.StatusBadRequest)
        return
    }

    var category models.Category
    if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
        response.SendError(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    category.ID = id

    result, err := database.UpdateCategory(category)
    if err != nil {
        response.SendError(w, "Failed to update category", http.StatusInternalServerError)
        return
    }

    response.SendResponse(w, response.APIResponse{
        Status: "success",
        Message: "Category updated successfully",
        Data: result,
    })
}