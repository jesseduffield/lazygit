package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir failed: %v", err)
	}

	scenarios := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "empty stays empty",
			input:    "",
			expected: "",
		},
		{
			name:     "no tilde left untouched",
			input:    filepath.Join("path", "to", "thing"),
			expected: filepath.Join("path", "to", "thing"),
		},
		{
			name:     "absolute path left untouched",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "tilde alone expands to home",
			input:    "~",
			expected: home,
		},
		{
			name:     "tilde slash expands to home/rest",
			input:    "~/foo/bar",
			expected: filepath.Join(home, "foo", "bar"),
		},
		{
			name:     "tilde without separator (~user style) is rejected",
			input:    "~bob/foo",
			expected: "~bob/foo",
			wantErr:  true,
		},
		{
			name:     "tilde mid-path is not expanded",
			input:    filepath.Join("foo", "~", "bar"),
			expected: filepath.Join("foo", "~", "bar"),
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			got, err := ExpandHomeDir(s.input)
			if s.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, s.expected, got)
		})
	}
}
