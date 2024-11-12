package handlers

import (
	"bytes"
	"errors"
	"testing"

	"github.com/cody-smarty/calc-lib"
)

func TestHandler_NilCalculator(t *testing.T) {
	handler := NewCLIHandler(nil, nil)
	err := handler.Handle([]string{})
	if err != nil {

	}
}

func TestHandler_ResultWrittenToOutput(t *testing.T) {
	stdout := &bytes.Buffer{}
	handler := NewCLIHandler(stdout, &calc.Addition{})

	err := handler.Handle([]string{"2", "3"})

	if err != nil {
		t.Errorf("want: %v, got: %v", nil, err)
	}
	if stdout.String() != "5" {
		t.Errorf("want: %v, got: %v", "5", stdout.String())
	}
}

func TestHandler_TwoArgsRequired(t *testing.T) {
	handler := NewCLIHandler(nil, &calc.Addition{})

	err := handler.Handle(nil)

	assertErr(t, err, errWrongArgCount)
}

func TestHandler_FirstArgInvalid(t *testing.T) {
	handler := NewCLIHandler(nil, &calc.Addition{})

	err := handler.Handle([]string{"invalid", "42"})

	assertErr(t, err, errInvalidArgument)
}

func TestHandler_SecondArgInvalid(t *testing.T) {
	handler := NewCLIHandler(nil, &calc.Addition{})

	err := handler.Handle([]string{"42", "invalid"})

	assertErr(t, err, errInvalidArgument)
}

func TestHandler_OutputWriterError(t *testing.T) {
	boink := errors.New("boink")
	stdout := ErringWriter{err: boink}
	handler := NewCLIHandler(stdout, &calc.Addition{})

	err := handler.Handle([]string{"2", "3"})

	assertErr(t, err, errOutputFailure)
}

func assertErr(t *testing.T, actual error, targets ...error) {
	for _, target := range targets {
		if !errors.Is(actual, target) {
			t.Helper()
			t.Errorf("want: %v, got: %v", target, actual)
		}
	}
}

type ErringWriter struct {
	err error
}

func (this ErringWriter) Write(p []byte) (n int, err error) {
	return 0, this.err
}
