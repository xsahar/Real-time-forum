package response

import (
    "encoding/json"
    "log"
    "net/http"
)

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SendError(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    if err := json.NewEncoder(w).Encode(APIResponse{
        Status: "error",
        Error:  message,
    }); err != nil {
        log.Printf("Error encoding response: %v", err)
    }
}

func SendResponse(w http.ResponseWriter, response APIResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding response: %v", err)
    }
}

func SendSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(APIResponse{
        Status:  "success",
        Message: message,
        Data:    data,
    }); err != nil {
        log.Printf("Error encoding response: %v", err)
    }
}

type CategoryResponse struct {

	Status  string      `json:"status"`

	Message string      `json:"message"`

	Data    interface{} `json:"data"`

}
