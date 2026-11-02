package config

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"

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
	unlock, err := c.lockAppState(sharedLock)
	if err != nil {
		return getDefaultAppState(), err
	}
	defer unlock()

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
	unlock, err := c.lockAppState(exclusiveLock)
	if err != nil {
		return err
	}
	defer unlock()

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

type lockKind int

const (
	sharedLock lockKind = iota
	exclusiveLock
)

// lockAppState locks state.yml.lock, a file next to state.yml, and returns a
// function that releases the lock. Updates take an exclusive lock and reads a
// shared one, so a read never sees a half-written file, and two updates never
// start from the same contents.
//
// We don't lock state.yml itself, because writing it replaces the file, and
// because on Windows a lock keeps other processes from reading the locked part
// of a file. We never delete the lock file. Otherwise one process could hold a
// lock on the deleted file while another one locks a new file of the same name.
func (c *AppConfig) lockAppState(kind lockKind) (func(), error) {
	if c.appStatePath == "" {
		return func() {}, nil
	}

	path := c.appStatePath + ".lock"
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil && os.IsPermission(err) {
		// We can't create the lock file in a read-only folder. Lock it if it
		// exists, and do without it otherwise; apparently when people have
		// read-only permissions they prefer us to fail silently.
		file, err = os.Open(path)
		if err != nil {
			return func() {}, nil
		}
	}
	if err != nil {
		return nil, err
	}

	if err := lockFile(file, kind); err != nil {
		file.Close()
		if isLockingUnsupported(err) {
			// Some network file systems don't support locks. Doing without is
			// better than never reading or saving the state.
			return func() {}, nil
		}
		return nil, err
	}

	return func() {
		_ = unlockFile(file)
		file.Close()
	}, nil
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

	err := writeFileAtomically(c.appStatePath, content)
	if err != nil && os.IsPermission(err) {
		// apparently when people have read-only permissions they prefer us to fail silently
		return nil
	}

	return err
}

// writeFileAtomically replaces the file at path with a new one that has the
// given content. A reader that doesn't take the lock sees either the old file
// or the complete new one. The same holds after a crash, or if the disk is full.
// It returns a permission error only if we may not write to the file or create
// a file in its folder.
func writeFileAtomically(path string, content []byte) error {
	// If path is a symlink, replace the file that it points to
	if resolvedPath, err := filepath.EvalSymlinks(path); err == nil {
		path = resolvedPath
	}

	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		// Don't replace a file that we're not allowed to write to
		file, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		file.Close()
		mode = info.Mode().Perm()
	}

	tempFile, err := createTempFile(path, mode)
	if err != nil {
		return err
	}
	err = writeAndClose(tempFile, content)
	if err == nil {
		err = os.Rename(tempFile.Name(), path)
		if err != nil {
			// On Windows, the rename fails with ERROR_ACCESS_DENIED if another
			// process has the file open. Keep only the message, so that this
			// doesn't count as a permission error.
			err = errors.New(err.Error())
		}
	}
	if err != nil {
		_ = os.Remove(tempFile.Name())
	}
	return err
}

// createTempFile creates a new file next to path. Unlike os.CreateTemp, it
// uses the given mode, minus the umask.
func createTempFile(path string, mode os.FileMode) (*os.File, error) {
	for {
		tempPath := fmt.Sprintf("%s.%d.tmp", path, rand.Uint32())
		file, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
		if !os.IsExist(err) {
			return file, err
		}
	}
}

func writeAndClose(file *os.File, content []byte) error {
	_, err := file.Write(content)
	if err == nil {
		// Without this, the rename can reach the disk before the content, and
		// a crash leaves a file that is filled with zero bytes. Users on
		// Windows have seen this happen.
		err = file.Sync()
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
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
