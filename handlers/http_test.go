package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/cody-smarty/calc-lib"
)

func TestHTTPHandler(t *testing.T) {
	tests := map[string]struct {
		request    string
		respStatus int
		respBody   string
	}{
		"Valid addition": {
			request:    "/add?a=3&b=4",
			respStatus: http.StatusOK,
			respBody:   "7",
		},
		"Valid subtraction": {
			request:    "/sub?b=4&a=3",
			respStatus: http.StatusOK,
			respBody:   "-1",
		},
		"Invalid operator": {
			request:    "/baz?a=3&4",
			respStatus: http.StatusNotFound,
			respBody:   "'baz' operator not supported",
		},
		"Bare operator": {
			request:    "/add",
			respStatus: http.StatusUnprocessableEntity,
			respBody:   "argument invalid:",
		},
		"Missing first arg": {
			request:    "/add?a=3",
			respStatus: http.StatusUnprocessableEntity,
			respBody:   "argument invalid:",
		},
		"Missing second arg": {
			request:    "/add?a=3",
			respStatus: http.StatusUnprocessableEntity,
			respBody:   "argument invalid:",
		},
		"Duplicate arg": {
			request:    "/add?a=3&b=4&a=5",
			respStatus: http.StatusOK,
			respBody:   "7",
		},
		"Invalid first args": {
			request:    "/add?a=bad&b=4&",
			respStatus: http.StatusUnprocessableEntity,
			respBody:   "argument invalid:",
		},
		"Invalid second args": {
			request:    "/add?a=&b=bad&",
			respStatus: http.StatusUnprocessableEntity,
			respBody:   "argument invalid:",
		},
	}

	for name, args := range tests {
		t.Run(name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "[TEST] ", 0)
			calculators := map[string]Calculator{
				"add": &calc.Addition{},
				"sub": &calc.Subtraction{},
			}
			handler := NewHTTPHandler(logger, calculators)
			req := httptest.NewRequest(http.MethodGet, args.request, nil)
			resp := httptest.NewRecorder()

			// For debugging
			reqDump, err := httputil.DumpRequest(req, true)
			assertErr(t, err, nil)
			t.Logf("request dump:\n%s", string(reqDump))

			handler.ServeHTTP(resp, req)

			// For debugging
			respDump, err := httputil.DumpResponse(resp.Result(), true)
			assertErr(t, err, nil)
			t.Logf("response dump:\n%s", string(respDump))

			if resp.Code != args.respStatus {
				t.Errorf("want: '%v', got: '%v'", args.respStatus, resp.Code)
			}

			// TODO -- Better way to compare error messages than hard coded string comparison
			if !strings.Contains(resp.Body.String(), args.respBody) {
				t.Errorf("want: '%v', got: '%v'", args.respBody, resp.Body.String())
			}
		})
	}
}
