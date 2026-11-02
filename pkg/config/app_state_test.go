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
	/* EXPECTED:
	assert.Equal(t, []string{"from B"}, appState.ShellCommandsHistory)
	ACTUAL: */
	assert.Empty(t, appState.ShellCommandsHistory)
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
	/* EXPECTED:
	assert.Equal(t, []string{"from B", "from A"}, appState.ShellCommandsHistory)
	ACTUAL: */
	assert.Equal(t, []string{"from B"}, appState.ShellCommandsHistory)
}
