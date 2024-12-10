package crud

import (
	"Real-time-forum/src/database"
	"log"
	"net/http"
    "Real-time-forum/src/utils/response"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
    // Add CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Content-Type", "application/json")

    // Handle preflight
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    if r.Method != http.MethodGet {
        response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    log.Printf("Fetching categories...")
    categories, err := database.GetCategories("")
    if err != nil {
        log.Printf("Error fetching categories: %v", err)
        response.SendError(w, "Failed to fetch categories", http.StatusInternalServerError)
        return
    }

    // Debug response data
    log.Printf("Found %d categories", len(categories))
    log.Printf("Response data: %+v", categories)

    response.SendResponse(w, response.APIResponse{
        Status: "success",
        Data:   categories,
    })
}