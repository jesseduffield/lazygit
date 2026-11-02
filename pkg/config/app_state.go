package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ReadAppState returns the app state. If it can't be read, it returns the
// defaults together with the error. Don't change the returned struct; use
// UpdateAppState for that.
func (c *AppConfig) ReadAppState() (*AppState, error) {
	return c.appState, nil
}

// UpdateAppState calls update to change the app state, and saves it.
func (c *AppConfig) UpdateAppState(update func(state *AppState)) error {
	update(c.appState)
	return c.saveAppState()
}

// saveAppState marshalls the AppState struct and writes it to the disk
func (c *AppConfig) saveAppState() error {
	marshalledAppState, err := yaml.Marshal(c.appState)
	if err != nil {
		return err
	}

	filepath, err := stateFilePath(stateFileName)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, marshalledAppState, 0o644)
	if err != nil && os.IsPermission(err) {
		// apparently when people have read-only permissions they prefer us to fail silently
		return nil
	}

	return err
}

var stateFileName = "state.yml"

// loadAppState loads recorded AppState from file
func loadAppState() (*AppState, error) {
	appState := getDefaultAppState()

	filepath, err := stateFilePath(stateFileName)
	if err != nil {
		if os.IsPermission(err) {
			// apparently when people have read-only permissions they prefer us to fail silently
			return appState, nil
		}
		return nil, err
	}

	appStateBytes, err := os.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if len(appStateBytes) == 0 {
		return appState, nil
	}

	err = yaml.Unmarshal(appStateBytes, appState)
	if err != nil {
		// Nothing in state.yml is important enough to keep lazygit from
		// starting. So if it can't be parsed (users have seen it filled with
		// zero bytes after a crash), start with the defaults. The next save
		// replaces the file.
		return getDefaultAppState(), nil
	}

	return appState, nil
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
