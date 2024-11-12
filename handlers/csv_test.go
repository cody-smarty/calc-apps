package handlers

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/cody-smarty/calc-lib"
)

func TestCSVHandler(t *testing.T) {
	tests := map[string]struct {
		input  string
		output string
	}{
		"Single Line": {
			input:  "1,+,2",
			output: "1,+,2,3\n",
		},
		"Multi Line": {
			input:  "1,+,2\n4,+,5\n",
			output: "1,+,2,3\n4,+,5,9\n",
		},
		"Too few args": {
			input:  "1,+,2\n1\n4,+,5\n",
			output: "1,+,2,3\n4,+,5,9\n",
		},
		"Too many args": {
			input:  "1,+,2\n9,+,8,7\n4,+,5\n",
			output: "1,+,2,3\n4,+,5,9\n",
		},
		"Invalid first args": {
			input:  "1,+,2\nNaN,+,8\n4,+,5\n",
			output: "1,+,2,3\n4,+,5,9\n",
		},
		"Invalid second args": {
			input:  "1,+,2\n9,+,NaN\n4,+,5\n",
			output: "1,+,2,3\n4,+,5,9\n",
		},
		"Invalid operand": {
			input:  "1,+,2\n9,Noop,8\n4,+,5\n",
			output: "1,+,2,3\n4,+,5,9\n",
		},
	}

	for name, args := range tests {
		t.Run(name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "[TEST] ", 0)
			input := strings.NewReader(args.input)
			var output bytes.Buffer
			caculators := map[string]Calculator{"+": &calc.Addition{}} // Only addition handled
			handler := NewCSVHandler(logger, input, &output, caculators)

			err := handler.Handle()

			assertErr(t, err, nil) // Errors not yet handled
			if output.String() != args.output {
				t.Errorf("want: '%v', got: '%v'", args.output, output.String())
			}
			if t.Failed() { // For debugging purposes
				t.Log(logBuf.String()) // TODO -- check the returned logs?
			}
		})
	}
}
