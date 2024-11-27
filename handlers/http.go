package handlers

import (
	"log"
	"net/http"
	"strconv"
)

type HTTPHandler struct {
	*http.ServeMux
	logger      *log.Logger
	calculators map[string]Calculator
}

func NewHTTPHandler(logger *log.Logger, calculators map[string]Calculator) *HTTPHandler {
	handler := HTTPHandler{
		ServeMux:    http.NewServeMux(),
		logger:      logger,
		calculators: calculators,
	}
	return &handler
}

func (this *HTTPHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	op := req.URL.Path[1:]
	calc, ok := this.calculators[op]
	if !ok {
		http.Error(resp, "'"+op+"' operator not supported", http.StatusNotFound)
		return
	}
	query := req.URL.Query()
	a, err := strconv.Atoi(query.Get("a"))
	if err != nil {
		http.Error(resp, "'a' argument invalid: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	b, err := strconv.Atoi(query.Get("b"))
	if err != nil {
		http.Error(resp, "'b' argument invalid: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	c := calc.Calculate(a, b)
	if _, err = resp.Write([]byte(strconv.Itoa(c))); err != nil {
		this.logger.Printf("failure writing result to response: %v", err)
	}
}
