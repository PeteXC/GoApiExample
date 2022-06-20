package main

import (
	"github.com/PeteXC/GoApiExample/api/handlers/maths/add"
	"github.com/akrylysov/algnhsa"
	"github.com/joerdav/zapray"
)

func main() {
	log, err := zapray.NewProduction()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	handler, err := add.NewHandler(log)
	if err != nil {
		panic(err)
	}
	xhandler := zapray.NewMiddleware("addMathApi", handler)
	algnhsa.ListenAndServe(xhandler, nil)
}
