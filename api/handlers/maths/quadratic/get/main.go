package main

import (
	"fmt"
	"os"

	"github.com/PeteXC/GoApiExample/api/handlers/maths/quadratic"
	"github.com/akrylysov/algnhsa"
	"github.com/joerdav/zapray"
	"github.com/pkg/errors"
)

func main() {
	log, err := zapray.NewProduction()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	offset, err := getMandatoryEnvironmentVariable("OFFSET")
	if err != nil {
		panic(errors.Wrap(err, "cannot retrieve offset"))
	}

	handler, err := quadratic.NewHandler(log, offset)
	if err != nil {
		panic(err)
	}
	xhandler := zapray.NewMiddleware("quadraticMathApi", handler)
	algnhsa.ListenAndServe(xhandler, nil)
}

func getMandatoryEnvironmentVariable(name string) (value string, err error) {
	value = os.Getenv(name)
	if value == "" {
		err = errors.New(fmt.Sprintf("Mandatory environment variable %q not set", name))
		return
	}
	return
}
