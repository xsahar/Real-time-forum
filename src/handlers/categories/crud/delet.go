package crud

import (
	"Real-time-forum/src/database"
	// "Real-time-forum/src/utils"
	"Real-time-forum/src/utils/response"
	"net/http"
	"strconv"
	"strings"
)

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    if r.Method != http.MethodDelete {
        response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    parts := strings.Split(r.URL.Path, "/")
    id, err := strconv.Atoi(parts[len(parts)-1])
    if err != nil {
        response.SendError(w, "Invalid category ID", http.StatusBadRequest)
        return
    }

    if err := database.DeleteCategory(id); err != nil {
        response.SendError(w, "Failed to delete category", http.StatusInternalServerError)
        return
    }

    response.SendResponse(w, response.APIResponse{
        Status: "success",
        Message: "Category deleted successfully",
    })
}