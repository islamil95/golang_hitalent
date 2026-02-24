package handler

import (
	"encoding/json"
	"net/http"
)

// JSON отправляет JSON-ответ.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// ErrorBody — стандартный формат тела ошибки.
type ErrorBody struct {
	Error string `json:"error"`
}

// Err пишет ответ с ошибкой в формате JSON.
func Err(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorBody{Error: message})
}
