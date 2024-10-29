package main

import (
	"log"
	"os"

	"github.com/cody-smarty/calc-app/handlers"
	"github.com/cody-smarty/calc-lib"
)

type Calculator interface {
	Calculate(a, b int) int
}

func main() {
	handler := handlers.NewCLIHandler(os.Stdout, &calc.Addition{})
	err := handler.Handle(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
