package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// useTempConfigDir points the config and state folder at a new temp dir and
// returns it.
func useTempConfigDir(t *testing.T) string {
	t.Helper()

	configDir := t.TempDir()
	t.Setenv("CONFIG_DIR", configDir)
	t.Setenv("LG_CONFIG_FILE", "")
	return configDir
}

func newAppConfig(t *testing.T) (*AppConfig, error) {
	t.Helper()

	return NewAppConfig("lazygit", "unversioned", "", "", "", false, t.TempDir())
}

func TestStartupWithBrokenAppStateFile(t *testing.T) {
	configDir := useTempConfigDir(t)
	err := os.WriteFile(filepath.Join(configDir, stateFileName), make([]byte, 4096), 0o644)
	assert.NoError(t, err)

	_, err = newAppConfig(t)
	assert.NoError(t, err)
}
