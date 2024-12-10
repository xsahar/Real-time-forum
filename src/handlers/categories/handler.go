package categories

import (
    "net/http"
    "Real-time-forum/src/handlers/categories/crud"
    "Real-time-forum/src/utils/response"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        crud.GetCategoriesHandler(w, r)
    case http.MethodPost:
        crud.CreateCategoryHandler(w, r)
    case http.MethodPut:
        crud.UpdateCategoryHandler(w, r)
    case http.MethodDelete:
        crud.DeleteCategoryHandler(w, r)
    default:
        response.SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}