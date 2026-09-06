package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Invoking lazygit as a daemon (e.g. as the editor of a rebase todo file) must
// not load or migrate the user config: writing a migrated config back to disk
// dirties the working tree in the middle of the git operation, aborting a
// rebase whose repo has the config file checked in (see #5998).
func TestNewAppConfigForDaemonDoesNotTouchUserConfig(t *testing.T) {
	stateDir := t.TempDir()
	t.Setenv("CONFIG_DIR", stateDir)

	configPath := filepath.Join(stateDir, ConfigFilename)
	preMigrationConfig := "git:\n  pagers:\n    - pager: less\n"
	if err := os.WriteFile(configPath, []byte(preMigrationConfig), 0o644); err != nil {
		t.Fatal(err)
	}

	appConfig, err := NewAppConfigForDaemon("lazygit", "test-version", "test-commit", "test-date", "test-build", false, t.TempDir())
	assert.NoError(t, err)
	if appConfig == nil {
		t.Fatal("expected non-nil app config")
	}

	// The config file must be untouched (not created, not migrated, not
	// rewritten).
	content, err := os.ReadFile(configPath)
	assert.NoError(t, err)
	assert.Equal(t, preMigrationConfig, string(content))

	// The daemon still gets a usable config object with defaults, and usable
	// app state.
	assert.NotNil(t, appConfig.GetUserConfig())
	assert.NotNil(t, appConfig.GetAppState())
}
