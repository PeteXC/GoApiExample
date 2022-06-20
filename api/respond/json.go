package respond

import (
	"encoding/json"
	"net/http"
)

func WithJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		http.Error(w, "failed to encode", http.StatusInternalServerError)
	}
}

func WithJSONStatus(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		http.Error(w, "failed to encode", http.StatusInternalServerError)
	}
}
