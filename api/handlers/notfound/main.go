package main

import (
	"net/http"

	"github.com/PeteXC/GoApiExample/api/respond"
	"github.com/akrylysov/algnhsa"
)

func main() {
	notFoundHandler := func(w http.ResponseWriter, r *http.Request) {
		respond.WithError(w, "not found", http.StatusNotFound)
	}
	algnhsa.ListenAndServe(http.HandlerFunc(notFoundHandler), nil)
}
