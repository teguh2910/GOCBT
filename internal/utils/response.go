package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// WriteJSONResponse writes a JSON response
func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteSuccessResponse writes a successful JSON response
func WriteSuccessResponse(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	WriteJSONResponse(w, response, http.StatusOK)
}

// WriteErrorResponse writes an error JSON response
func WriteErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	WriteJSONResponse(w, response, statusCode)
}

// WriteCreatedResponse writes a created response
func WriteCreatedResponse(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	WriteJSONResponse(w, response, http.StatusCreated)
}

// WriteNoContentResponse writes a no content response
func WriteNoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
