package main

import (
	"net/http"

	"github.com/PeteXC/GoApiExample/api/handlers/maths/add"
	"github.com/PeteXC/GoApiExample/api/handlers/maths/quadratic"
	"github.com/gorilla/mux"
	"github.com/joerdav/zapray"
)

func main() {
	log, err := zapray.NewProduction()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	r := mux.NewRouter()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	mathsAdditionHandler, _ := add.NewHandler(log)
	r.Handle("/maths/add", mathsAdditionHandler)

	mathsQuadraticHandler, _ := quadratic.NewHandler(log, "yes")
	r.Handle("/maths/quadratic", mathsQuadraticHandler)

	log.Info("Starting local api")
	err = http.ListenAndServe("127.0.0.1:8080", r)
	if err != nil {
		panic(err)
	}
}
