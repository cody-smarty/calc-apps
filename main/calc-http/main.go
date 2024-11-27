package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cody-smarty/calc-app/handlers"
	"github.com/cody-smarty/calc-lib"
)

const addr = "localhost:8080"

func main() {
	logger := log.New(os.Stderr, "http: ", 0)
	log.Printf("Server listening on %s...", addr)
	defer log.Println("Server stopped")
	router := handlers.NewHTTPHandler(logger, calculators)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

var calculators = map[string]handlers.Calculator{
	"add": &calc.Addition{},
	"sub": &calc.Subtraction{},
	"mul": &calc.Multiplication{},
	"div": &calc.Division{},
}
