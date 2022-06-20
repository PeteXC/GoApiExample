package respond

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message    string   `json:"message"`
	Issues     []string `json:"issues,omitempty"`
	StatusCode int64    `json:"statusCode"`
}

func WithError(w http.ResponseWriter, msg string, status int, issues ...string) {
	enc := json.NewEncoder(w)
	w.WriteHeader(status)
	err := enc.Encode(ErrorResponse{
		Message:    msg,
		Issues:     issues,
		StatusCode: int64(status),
	})
	if err != nil {
		http.Error(w, msg, status)
	}
}

func WithNotFound(w http.ResponseWriter, issues ...string) {
	WithError(w, "not found", http.StatusNotFound, issues...)
}

func WithInternalServerError(w http.ResponseWriter, issues ...string) {
	WithError(w, "something went wrong", http.StatusInternalServerError, issues...)
}

func WithBadRequest(w http.ResponseWriter, issues ...string) {
	WithError(w, "bad request", http.StatusBadRequest, issues...)
}

func WithForbidden(w http.ResponseWriter, issues ...string) {
	WithError(w, "forbiden", http.StatusForbidden, issues...)
}
