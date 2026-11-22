package updates

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestCheckForUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tag_name":"v0.67.0"}`))
	}))
	defer server.Close()

	scenarios := []struct {
		currentVersion string
		expected       *types.VersionNumber
	}{
		{currentVersion: "0.66.2", expected: &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}},
		{currentVersion: "0.67.0", expected: nil},
		{currentVersion: "0.67.1", expected: nil},
	}

	for _, s := range scenarios {
		t.Run(s.currentVersion, func(t *testing.T) {
			updater := newTestUpdater(t, server, s.currentVersion, Installation{Kind: InstallKindRelease})
			newVersion, err := updater.CheckForUpdate()
			assert.NoError(t, err)
			assert.Equal(t, s.expected, newVersion)

			state := updater.readAppState()
			assert.Equal(t, "v0.67.0", state.Updates.LatestVersion)
			assert.WithinDuration(t, time.Now(), time.Unix(state.Updates.LastCheck, 0), time.Minute)
		})
	}
}

func TestCheckForUpdateIfDue(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte(`{"tag_name":"v0.67.0"}`))
	}))
	defer server.Close()

	now := time.Date(2026, 11, 3, 9, 0, 0, 0, time.Local)
	scenarios := []struct {
		name          string
		version       string
		method        string
		days          int64
		lastCheck     time.Time
		latestVersion string
		expectCheck   bool
	}{
		{
			name:        "never checked",
			version:     "0.66.2",
			method:      "prompt",
			days:        1,
			expectCheck: true,
		},
		{
			name:          "checked earlier today",
			version:       "0.66.2",
			method:        "prompt",
			days:          1,
			lastCheck:     time.Date(2026, 11, 3, 0, 1, 0, 0, time.Local),
			latestVersion: "v0.66.3",
			expectCheck:   false,
		},
		{
			name:        "checked yesterday, less than 24 hours ago",
			version:     "0.66.2",
			method:      "prompt",
			days:        1,
			lastCheck:   time.Date(2026, 11, 2, 10, 0, 0, 0, time.Local),
			expectCheck: true,
		},
		{
			name:          "checked six days ago, checking once a week",
			version:       "0.66.2",
			method:        "prompt",
			days:          7,
			lastCheck:     time.Date(2026, 10, 28, 9, 0, 0, 0, time.Local),
			latestVersion: "v0.66.2",
			expectCheck:   false,
		},
		{
			name:        "checked seven days ago, checking once a week",
			version:     "0.66.2",
			method:      "prompt",
			days:        7,
			lastCheck:   time.Date(2026, 10, 27, 23, 0, 0, 0, time.Local),
			expectCheck: true,
		},
		{
			name:        "checked earlier today, checking at every start",
			version:     "0.66.2",
			method:      "prompt",
			days:        0,
			lastCheck:   time.Date(2026, 11, 3, 8, 59, 0, 0, time.Local),
			expectCheck: true,
		},
		{
			name:        "checked in the future, when the clock was wrong",
			version:     "0.66.2",
			method:      "prompt",
			days:        1,
			lastCheck:   time.Date(2027, 5, 1, 9, 0, 0, 0, time.Local),
			expectCheck: true,
		},
		{
			name:        "update.method never",
			version:     "0.66.2",
			method:      "never",
			days:        1,
			expectCheck: true,
		},
		{
			name:        "dev build",
			version:     "0123abcd",
			method:      "prompt",
			days:        1,
			expectCheck: false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			requests = 0
			updater := newTestUpdater(t, server, "0.66.2", Installation{Kind: InstallKindRelease})
			updater.installation = DetectInstallation(BuildSourceRelease, s.version, "/home/user/.local/bin/lazygit")
			updater.UserConfig().Update.Method = s.method
			updater.UserConfig().Update.Days = s.days
			updater.now = func() time.Time { return now }
			if !s.lastCheck.IsZero() {
				updater.appState.(*fakeAppStateStore).state.Updates.LastCheck = s.lastCheck.Unix()
			}
			updater.appState.(*fakeAppStateStore).state.Updates.LatestVersion = s.latestVersion

			assert.NoError(t, updater.CheckForUpdateIfDue())
			assert.Equal(t, s.expectCheck, requests == 1)
		})
	}
}

func TestStartOfDayAfter(t *testing.T) {
	scenarios := []struct {
		time     time.Time
		days     int
		expected time.Time
	}{
		{time.Date(2026, 11, 3, 9, 5, 0, 0, time.UTC), 0, time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 11, 3, 9, 5, 0, 0, time.UTC), 1, time.Date(2026, 11, 4, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC), 1, time.Date(2026, 11, 4, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 11, 3, 23, 59, 0, 0, time.UTC), 7, time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 11, 30, 23, 59, 0, 0, time.UTC), 1, time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)},
		{time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC), 1, time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, s := range scenarios {
		assert.Equal(t, s.expected, startOfDayAfter(s.time, s.days))
	}
}

func TestInstall(t *testing.T) {
	server := newReleaseServer(t)

	dir := t.TempDir()
	executablePath := filepath.Join(dir, "lazygit")
	writeOldExecutable(t, executablePath)

	updater := newTestUpdater(t, server, "0.66.2", Installation{Kind: InstallKindRelease, ExecutablePath: executablePath})
	assert.NoError(t, updater.CheckCanInstall())
	assert.False(t, updater.IsUpdateInstalled())

	assert.NoError(t, updater.Install(&types.VersionNumber{Major: 0, Minor: 67, Patch: 0}))
	assertFileContent(t, executablePath, "new")
	assertDirContains(t, dir, "lazygit")
	assert.True(t, updater.IsUpdateInstalled())

	// After a restart, the installed version is the running one
	restartedUpdater := newTestUpdater(t, server, "0.67.0", Installation{Kind: InstallKindRelease, ExecutablePath: executablePath})
	assert.False(t, restartedUpdater.IsUpdateInstalled())
}

func TestInstallFails(t *testing.T) {
	server := newReleaseServer(t)

	dir := t.TempDir()
	executablePath := filepath.Join(dir, "lazygit")
	writeOldExecutable(t, executablePath)

	updater := newTestUpdater(t, server, "0.66.2", Installation{Kind: InstallKindRelease, ExecutablePath: executablePath})
	err := updater.Install(&types.VersionNumber{Major: 0, Minor: 68, Patch: 0})
	assert.ErrorContains(t, err, "lazygit_0.68.0_linux_x86_64.tar.gz returned 404 Not Found")
	assertFileContent(t, executablePath, "old")
	assertDirContains(t, dir, "lazygit")
	assert.False(t, updater.IsUpdateInstalled())
}

func TestInstallFromTwoInstances(t *testing.T) {
	server := newReleaseServer(t)

	dir := t.TempDir()
	executablePath := filepath.Join(dir, "lazygit")
	writeOldExecutable(t, executablePath)

	installation := Installation{Kind: InstallKindRelease, ExecutablePath: executablePath}
	updater := newTestUpdater(t, server, "0.66.2", installation)
	otherUpdater := newTestUpdater(t, server, "0.66.2", installation)

	// The other lazygit has downloaded the update and is about to replace the
	// executable when this one finishes installing it
	newPath, err := writeNewExecutable(executablePath, []byte("from the other lazygit"))
	assert.NoError(t, err)
	assert.NoError(t, updater.Install(&types.VersionNumber{Major: 0, Minor: 67, Patch: 0}))

	assert.True(t, otherUpdater.IsUpdateInstalled())
	assert.NoError(t, otherUpdater.installNewExecutable(newPath))
	assertFileContent(t, executablePath, "new")
	assertDirContains(t, dir, "lazygit")

	// Installing again doesn't download anything; the server doesn't have this
	// version
	assert.NoError(t, otherUpdater.Install(&types.VersionNumber{Major: 0, Minor: 68, Patch: 0}))
	assertFileContent(t, executablePath, "new")
}

// newReleaseServer returns a server that has the release archive of v0.67.0
// for linux on x86_64, with "new" as the content of the lazygit binary
func newReleaseServer(t *testing.T) *httptest.Server {
	t.Helper()

	archive := makeTarGz(t, map[string]string{"LICENSE": "license", "lazygit": "new"})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/releases/download/v0.67.0/lazygit_0.67.0_linux_x86_64.tar.gz" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write(archive)
	}))
	t.Cleanup(server.Close)
	return server
}

// writeOldExecutable writes an executable with "old" as its content. It
// backdates the file, because some file systems record the modification time
// with a granularity of several milliseconds, and installing an update right
// afterwards could otherwise give the new file the same time.
func writeOldExecutable(t *testing.T, path string) {
	t.Helper()

	assert.NoError(t, os.WriteFile(path, []byte("old"), 0o755))
	modTime := time.Now().Add(-time.Hour)
	assert.NoError(t, os.Chtimes(path, modTime, modTime))
}

func newTestUpdater(t *testing.T, server *httptest.Server, currentVersion string, installation Installation) *Updater {
	t.Helper()

	version, err := types.ParseVersionNumber(currentVersion)
	assert.NoError(t, err)

	return &Updater{
		Common:            common.NewDummyCommon(),
		appState:          &fakeAppStateStore{state: &config.AppState{}},
		installation:      installation,
		currentVersion:    version,
		executableModTime: fileModTime(installation.ExecutablePath),
		goos:              "linux",
		goarch:            "amd64",
		httpClient:        server.Client(),
		repoURL:           server.URL,
		now:               time.Now,
	}
}

type fakeAppStateStore struct {
	state *config.AppState
}

func (self *fakeAppStateStore) ReadAppState() (*config.AppState, error) {
	state := *self.state
	return &state, nil
}

func (self *fakeAppStateStore) UpdateAppState(update func(state *config.AppState)) error {
	update(self.state)
	return nil
}
