package handlers

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/cody-smarty/calc-lib"
)

type Calculator interface {
	Calculate(a, b int) int
}

type CLIHandler struct {
	stdout     io.Writer
	calculator Calculator
}

func NewCLIHandler(stdout io.Writer, calculator *calc.Addition) *CLIHandler {
	return &CLIHandler{
		stdout:     stdout,
		calculator: calculator,
	}
}

func (hand *CLIHandler) Handle(args []string) error {
	if len(args) != 2 {
		return errWrongArgCount
	}
	a, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("%w: '%s'", errInvalidArgument, args[0])
	}
	b, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("%w: '%s'", errInvalidArgument, args[1])
	}
	result := hand.calculator.Calculate(a, b)
	_, err = fmt.Fprint(hand.stdout, result)
	if err != nil {
		return fmt.Errorf("%w: '%s'", errOutputFailure, args[0])
	}
	return err
}

var (
	errWrongArgCount   = errors.New("usage: calculator <a> <b>")
	errInvalidArgument = errors.New("invalid argument")
	errOutputFailure   = errors.New("output failure")
)
