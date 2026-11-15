package updates

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestFetchLatestVersion(t *testing.T) {
	scenarios := []struct {
		name            string
		status          int
		body            string
		expectedVersion *types.VersionNumber
		expectedError   string
	}{
		{
			name:            "release",
			status:          http.StatusOK,
			body:            `{"id":1,"tag_name":"v0.67.0","update_url":"/jesseduffield/lazygit/releases/tag/v0.67.0"}`,
			expectedVersion: &types.VersionNumber{Major: 0, Minor: 67, Patch: 0},
		},
		{
			name:          "not found",
			status:        http.StatusNotFound,
			body:          `{"error":"Not Found"}`,
			expectedError: "/releases/latest returned 404 Not Found",
		},
		{
			name:          "unexpected tag",
			status:        http.StatusOK,
			body:          `{"tag_name":"nightly"}`,
			expectedError: "unexpected version format: nightly",
		},
		{
			name:          "invalid json",
			status:        http.StatusOK,
			body:          `<html>`,
			expectedError: "invalid character '<' looking for beginning of value",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/releases/latest", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Accept"))
				w.WriteHeader(s.status)
				_, _ = w.Write([]byte(s.body))
			}))
			defer server.Close()

			version, err := fetchLatestVersion(context.Background(), server.Client(), server.URL)
			if s.expectedError != "" {
				assert.ErrorContains(t, err, s.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, s.expectedVersion, version)
			}
		})
	}
}
