package config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

var stateFileName = "state.yml"

// appStateFilePath returns the path of state.yml. It returns "" if we lack the
// permissions to create its folder; we don't persist the app state then.
func appStateFilePath() (string, error) {
	path, err := stateFilePath(stateFileName)
	if err != nil {
		if os.IsPermission(err) {
			// apparently when people have read-only permissions they prefer us to fail silently
			return "", nil
		}
		return "", err
	}
	return path, nil
}

// ReadAppState reads state.yml and returns its contents. If it can't be read,
// it returns the defaults together with the error. Changing the returned struct
// has no effect; use UpdateAppState for that.
func (c *AppConfig) ReadAppState() (*AppState, error) {
	content, err := c.readAppStateFile()
	if err != nil {
		return getDefaultAppState(), err
	}
	return parseAppState(content)
}

// UpdateAppState reads state.yml, calls update to change what it found, and
// writes the result back. Other lazygit instances can change the file at any
// time, so update must apply its change to the state that it gets, e.g. add an
// entry to a list. It must not set a field to a value that the caller read
// earlier, because that would undo changes that other instances made in the
// meantime.
func (c *AppConfig) UpdateAppState(update func(state *AppState)) error {
	content, err := c.readAppStateFile()
	if err != nil {
		return err
	}
	// If the file can't be parsed, we start from the defaults and replace it.
	appState, _ := parseAppState(content)
	update(appState)

	newContent, err := yaml.Marshal(appState)
	if err != nil {
		return err
	}
	if bytes.Equal(newContent, content) {
		return nil
	}
	return c.writeAppStateFile(newContent)
}

func (c *AppConfig) readAppStateFile() ([]byte, error) {
	if c.appStatePath == "" {
		return nil, nil
	}

	content, err := os.ReadFile(c.appStatePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return content, nil
}

// parseAppState returns the defaults together with the error if content can't
// be parsed. Nothing in state.yml is important enough to keep lazygit from
// working, and users have seen the file filled with zero bytes after a crash.
func parseAppState(content []byte) (*AppState, error) {
	appState := getDefaultAppState()
	if err := yaml.Unmarshal(content, appState); err != nil {
		return getDefaultAppState(), err
	}
	return appState, nil
}

func (c *AppConfig) writeAppStateFile(content []byte) error {
	if c.appStatePath == "" {
		return nil
	}

	err := os.WriteFile(c.appStatePath, content, 0o644)
	if err != nil && os.IsPermission(err) {
		// apparently when people have read-only permissions they prefer us to fail silently
		return nil
	}

	return err
}

// AppState stores data between runs of the app like when the last update check
// was performed and which other repos have been checked out
type AppState struct {
	LastUpdateCheck     int64
	RecentRepos         []string
	StartupPopupVersion int
	LastVersion         string // this is the last version the user was using, for the purpose of showing release notes

	// these are for shell commands typed in directly, not for custom commands in the lazygit config.
	// For backwards compatibility we keep the old name in yaml files.
	ShellCommandsHistory []string `yaml:"customcommandshistory"`

	HideCommandLog bool
}

func getDefaultAppState() *AppState {
	return &AppState{}
}
