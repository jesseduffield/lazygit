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

func TestBrokenAppStateFile(t *testing.T) {
	configDir := useTempConfigDir(t)
	err := os.WriteFile(filepath.Join(configDir, stateFileName), make([]byte, 4096), 0o644)
	assert.NoError(t, err)

	appConfig, err := newAppConfig(t)
	assert.NoError(t, err)

	appState, err := appConfig.ReadAppState()
	assert.ErrorContains(t, err, "control characters are not allowed")
	assert.Equal(t, getDefaultAppState(), appState)

	assert.NoError(t, appConfig.UpdateAppState(addToHistory("ls")))
	appState, err = appConfig.ReadAppState()
	assert.NoError(t, err)
	assert.Equal(t, []string{"ls"}, appState.ShellCommandsHistory)
}

func addToHistory(command string) func(state *AppState) {
	return func(state *AppState) {
		state.ShellCommandsHistory = append([]string{command}, state.ShellCommandsHistory...)
	}
}

func TestAppStateReadSeesUpdatesOfOtherInstances(t *testing.T) {
	useTempConfigDir(t)
	instanceA, err := newAppConfig(t)
	assert.NoError(t, err)
	instanceB, err := newAppConfig(t)
	assert.NoError(t, err)

	assert.NoError(t, instanceB.UpdateAppState(addToHistory("from B")))

	appState, err := instanceA.ReadAppState()
	assert.NoError(t, err)
	assert.Equal(t, []string{"from B"}, appState.ShellCommandsHistory)
}

func TestAppStateUpdateKeepsUpdatesOfOtherInstances(t *testing.T) {
	useTempConfigDir(t)
	instanceA, err := newAppConfig(t)
	assert.NoError(t, err)
	instanceB, err := newAppConfig(t)
	assert.NoError(t, err)

	assert.NoError(t, instanceA.UpdateAppState(addToHistory("from A")))
	assert.NoError(t, instanceB.UpdateAppState(addToHistory("from B")))

	instanceC, err := newAppConfig(t)
	assert.NoError(t, err)
	appState, err := instanceC.ReadAppState()
	assert.NoError(t, err)
	assert.Equal(t, []string{"from B", "from A"}, appState.ShellCommandsHistory)
}
