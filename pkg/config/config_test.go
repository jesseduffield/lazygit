package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandHomeDir(t *testing.T) {
	type scenario struct {
		configDir    string
		expectedPath string
	}

	home, _ := os.UserHomeDir()
	scenarios := []scenario{
		{
			"~/.config/lazygit",
			filepath.Join(home, ".config/lazygit"),
		},
	}

	for _, s := range scenarios {
		t.Run(s.configDir, func(t *testing.T) {
			result := expandHomeDir(s.configDir)
			assert.Equal(t, result, s.expectedPath)
		})
	}
}
