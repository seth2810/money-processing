package routes

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func replyWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func replyWithError(w http.ResponseWriter, statusCode int, err error) {
	replyWithJSON(w, statusCode, errorResponse{err.Error()})
}
