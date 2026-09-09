package response

import (
	"encoding/json"
	"net/http"
)

// Envelope is the standard JSON shape every API response follows, success
// or failure, so clients only ever need one parser.
type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON writes any payload with the given status code.
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Success writes a 2xx envelope wrapping data.
func Success(w http.ResponseWriter, status int, data interface{}) {
	JSON(w, status, Envelope{Success: true, Data: data})
}

// Error writes an error envelope with the given status and message.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Envelope{Success: false, Error: message})
}
