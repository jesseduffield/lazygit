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
			assert.Equal(t, s.expected, updater.Status().AvailableVersion)

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

	now := time.Now()
	scenarios := []struct {
		name          string
		version       string
		method        string
		days          int64
		lastCheck     time.Time
		latestVersion string
		expectCheck   bool
		expectStatus  *types.VersionNumber
	}{
		{
			name:         "never checked",
			version:      "0.66.2",
			method:       "prompt",
			days:         1,
			expectCheck:  true,
			expectStatus: &types.VersionNumber{Major: 0, Minor: 67, Patch: 0},
		},
		{
			name:          "checked recently",
			version:       "0.66.2",
			method:        "prompt",
			days:          1,
			lastCheck:     now.Add(-23 * time.Hour),
			latestVersion: "v0.66.3",
			expectCheck:   false,
			expectStatus:  &types.VersionNumber{Major: 0, Minor: 66, Patch: 3},
		},
		{
			name:         "checked a day ago",
			version:      "0.66.2",
			method:       "prompt",
			days:         1,
			lastCheck:    now.Add(-25 * time.Hour),
			expectCheck:  true,
			expectStatus: &types.VersionNumber{Major: 0, Minor: 67, Patch: 0},
		},
		{
			name:          "checked two days ago, checking once a week",
			version:       "0.66.2",
			method:        "prompt",
			days:          7,
			lastCheck:     now.Add(-48 * time.Hour),
			latestVersion: "v0.66.2",
			expectCheck:   false,
			expectStatus:  nil,
		},
		{
			name:         "update.method never",
			version:      "0.66.2",
			method:       "never",
			days:         1,
			expectCheck:  true,
			expectStatus: &types.VersionNumber{Major: 0, Minor: 67, Patch: 0},
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
			if !s.lastCheck.IsZero() {
				updater.appState.(*fakeAppStateStore).state.Updates.LastCheck = s.lastCheck.Unix()
			}
			updater.appState.(*fakeAppStateStore).state.Updates.LatestVersion = s.latestVersion

			assert.NoError(t, updater.CheckForUpdateIfDue())
			assert.Equal(t, s.expectCheck, requests == 1)
			assert.Equal(t, s.expectStatus, updater.Status().AvailableVersion)
		})
	}
}

func TestRefreshStatus(t *testing.T) {
	updater := newTestUpdater(t, nil, "0.66.2", Installation{Kind: InstallKindRelease, ExecutablePath: "/bin/lazygit"})
	state := updater.appState.(*fakeAppStateStore).state

	assert.Equal(t, Status{}, updater.RefreshStatus())

	state.Updates.LatestVersion = "v0.66.2"
	assert.Equal(t, Status{}, updater.RefreshStatus())

	state.Updates.LatestVersion = "v0.67.0"
	assert.Equal(t, Status{AvailableVersion: &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}}, updater.RefreshStatus())

	state.Updates.InstalledVersion = "v0.67.0"
	state.Updates.InstalledPath = "/other/lazygit"
	assert.Equal(t, Status{AvailableVersion: &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}}, updater.RefreshStatus())

	state.Updates.InstalledPath = "/bin/lazygit"
	assert.Equal(t, Status{
		AvailableVersion: &types.VersionNumber{Major: 0, Minor: 67, Patch: 0},
		InstalledVersion: &types.VersionNumber{Major: 0, Minor: 67, Patch: 0},
	}, updater.RefreshStatus())
}

func TestInstall(t *testing.T) {
	archive := makeTarGz(t, map[string]string{"LICENSE": "license", "lazygit": "new"})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/releases/download/v0.67.0/lazygit_0.67.0_linux_x86_64.tar.gz" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write(archive)
	}))
	defer server.Close()

	dir := t.TempDir()
	executablePath := filepath.Join(dir, "lazygit")
	assert.NoError(t, os.WriteFile(executablePath, []byte("old"), 0o755))

	updater := newTestUpdater(t, server, "0.66.2", Installation{Kind: InstallKindRelease, ExecutablePath: executablePath})
	assert.NoError(t, updater.CheckCanInstall())
	assert.Nil(t, updater.RefreshStatus().InstalledVersion)

	assert.NoError(t, updater.Install(&types.VersionNumber{Major: 0, Minor: 67, Patch: 0}))
	assertFileContent(t, executablePath, "new")
	assertDirContains(t, dir, "lazygit")
	assert.Equal(t, &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}, updater.Status().InstalledVersion)

	err := updater.Install(&types.VersionNumber{Major: 0, Minor: 68, Patch: 0})
	assert.ErrorContains(t, err, "lazygit_0.68.0_linux_x86_64.tar.gz returned 404 Not Found")
	assertFileContent(t, executablePath, "new")
	assert.Equal(t, &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}, updater.RefreshStatus().InstalledVersion)

	// Another lazygit, running from a different executable, doesn't see the
	// installed version
	otherUpdater := newTestUpdater(t, server, "0.66.2", Installation{Kind: InstallKindRelease, ExecutablePath: "/other/lazygit"})
	otherUpdater.appState = updater.appState
	assert.Nil(t, otherUpdater.RefreshStatus().InstalledVersion)

	// After a restart, the installed version is the running one
	restartedUpdater := newTestUpdater(t, server, "0.67.0", Installation{Kind: InstallKindRelease, ExecutablePath: executablePath})
	restartedUpdater.appState = updater.appState
	assert.Nil(t, restartedUpdater.RefreshStatus().InstalledVersion)
}

func newTestUpdater(t *testing.T, server *httptest.Server, currentVersion string, installation Installation) *Updater {
	t.Helper()

	version, err := types.ParseVersionNumber(currentVersion)
	assert.NoError(t, err)

	updater := &Updater{
		Common:         common.NewDummyCommon(),
		appState:       &fakeAppStateStore{state: &config.AppState{}},
		installation:   installation,
		currentVersion: version,
		goos:           "linux",
		goarch:         "amd64",
	}
	if server != nil {
		updater.httpClient = server.Client()
		updater.repoURL = server.URL
	}
	return updater
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
