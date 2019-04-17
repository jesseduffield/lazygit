package cleanhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPrintablePathCheckHandler(t *testing.T) {
	getTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	})

	cases := map[string]struct {
		path       string
		expectCode int
		input      *HandlerInput
	}{
		"valid nil input": {
			path:       "/valid",
			expectCode: http.StatusOK,
			input:      nil,
		},

		"valid empty error status": {
			path:       "/valid",
			expectCode: http.StatusOK,
			input:      &HandlerInput{},
		},

		"invalid newline": {
			path:       "/invalid\n",
			expectCode: http.StatusBadRequest,
		},

		"invalid carriage return": {
			path:       "/invalid\r",
			expectCode: http.StatusBadRequest,
		},

		"invalid null": {
			path:       "/invalid\x00",
			expectCode: http.StatusBadRequest,
		},

		"invalid alternate status": {
			path:       "/invalid\n",
			expectCode: http.StatusInternalServerError,
			input: &HandlerInput{
				ErrStatus: http.StatusInternalServerError,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Create test HTTP server
			ts := httptest.NewServer(PrintablePathCheckHandler(getTestHandler, tc.input))
			defer ts.Close()

			res, err := http.Get(ts.URL + tc.path)
			if err != nil {
				t.Fatal(err)
			}

			if tc.expectCode != res.StatusCode {
				t.Fatalf("expected %d, got :%d", tc.expectCode, res.StatusCode)
			}
		})
	}
}
