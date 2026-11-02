package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/samber/lo"
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

func TestAppStateUpdateSkipsWritingUnchangedState(t *testing.T) {
	configDir := useTempConfigDir(t)
	statePath := filepath.Join(configDir, stateFileName)
	appConfig, err := newAppConfig(t)
	assert.NoError(t, err)
	assert.NoError(t, appConfig.UpdateAppState(addToHistory("ls")))
	past := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.NoError(t, os.Chtimes(statePath, past, past))

	assert.NoError(t, appConfig.UpdateAppState(func(state *AppState) {
		state.ShellCommandsHistory = lo.Uniq(append([]string{"ls"}, state.ShellCommandsHistory...))
	}))

	info, err := os.Stat(statePath)
	assert.NoError(t, err)
	assert.True(t, info.ModTime().Equal(past))
}

func TestAppStateConcurrentUpdatesOfSeveralInstances(t *testing.T) {
	useTempConfigDir(t)
	const instanceCount = 8
	const updatesPerInstance = 10

	instances := make([]*AppConfig, instanceCount)
	for i := range instances {
		instance, err := newAppConfig(t)
		assert.NoError(t, err)
		instances[i] = instance
	}

	var waitGroup sync.WaitGroup
	for i, instance := range instances {
		waitGroup.Go(func() {
			for j := range updatesPerInstance {
				err := instance.UpdateAppState(addToHistory(fmt.Sprintf("command %d-%d", i, j)))
				assert.NoError(t, err)
				_, err = instance.ReadAppState()
				assert.NoError(t, err)
			}
		})
	}
	waitGroup.Wait()

	appState, err := instances[0].ReadAppState()
	assert.NoError(t, err)
	assert.Len(t, appState.ShellCommandsHistory, instanceCount*updatesPerInstance)
}

func TestAppStateInReadOnlyFolder(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs file permissions that apply to the current user")
	}

	configDir := useTempConfigDir(t)
	appConfig, err := newAppConfig(t)
	assert.NoError(t, err)
	statePath := filepath.Join(configDir, stateFileName)
	assert.NoError(t, os.WriteFile(statePath, []byte("customcommandshistory:\n    - ls\n"), 0o444))
	assert.NoError(t, os.Chmod(configDir, 0o555))
	t.Cleanup(func() { _ = os.Chmod(configDir, 0o755) })

	assert.NoError(t, appConfig.UpdateAppState(addToHistory("pwd")))

	appState, err := appConfig.ReadAppState()
	assert.NoError(t, err)
	assert.Equal(t, []string{"ls"}, appState.ShellCommandsHistory)
}

func TestAppStateUpdateReplacesTheTargetOfASymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks needs special privileges on Windows")
	}

	configDir := useTempConfigDir(t)
	statePath := filepath.Join(configDir, stateFileName)
	targetPath := filepath.Join(t.TempDir(), "synced-state.yml")
	assert.NoError(t, os.WriteFile(targetPath, nil, 0o644))
	assert.NoError(t, os.Symlink(targetPath, statePath))
	appConfig, err := newAppConfig(t)
	assert.NoError(t, err)

	assert.NoError(t, appConfig.UpdateAppState(addToHistory("ls")))

	info, err := os.Lstat(statePath)
	assert.NoError(t, err)
	assert.Equal(t, os.ModeSymlink, info.Mode().Type())
	content, err := os.ReadFile(targetPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "- ls\n")
}

func TestAppStateUpdateKeepsFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs file permissions that apply to the current user")
	}

	configDir := useTempConfigDir(t)
	statePath := filepath.Join(configDir, stateFileName)
	assert.NoError(t, os.WriteFile(statePath, nil, 0o600))
	appConfig, err := newAppConfig(t)
	assert.NoError(t, err)

	assert.NoError(t, appConfig.UpdateAppState(addToHistory("ls")))
	info, err := os.Stat(statePath)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	// A read-only file stays as it is, without an error
	assert.NoError(t, os.Chmod(statePath, 0o400))
	assert.NoError(t, appConfig.UpdateAppState(addToHistory("pwd")))
	appState, err := appConfig.ReadAppState()
	assert.NoError(t, err)
	assert.Equal(t, []string{"ls"}, appState.ShellCommandsHistory)
	info, err = os.Stat(statePath)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0o400), info.Mode().Perm())
}

func TestAppStateUpdateReportsFileHeldOpenByAnotherProcess(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("only Windows refuses to replace a file that another process has open")
	}

	configDir := useTempConfigDir(t)
	appConfig, err := newAppConfig(t)
	assert.NoError(t, err)
	assert.NoError(t, appConfig.UpdateAppState(addToHistory("ls")))
	file, err := os.Open(filepath.Join(configDir, stateFileName))
	assert.NoError(t, err)
	defer file.Close()

	assert.Error(t, appConfig.UpdateAppState(addToHistory("pwd")))
}
