package updates

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

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
		})
	}
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
	assert.Nil(t, updater.InstalledVersion())

	assert.NoError(t, updater.Install(&types.VersionNumber{Major: 0, Minor: 67, Patch: 0}))
	assertFileContent(t, executablePath, "new")
	assertDirContains(t, dir, "lazygit")
	assert.Equal(t, &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}, updater.InstalledVersion())

	err := updater.Install(&types.VersionNumber{Major: 0, Minor: 68, Patch: 0})
	assert.ErrorContains(t, err, "lazygit_0.68.0_linux_x86_64.tar.gz returned 404 Not Found")
	assertFileContent(t, executablePath, "new")
	assert.Equal(t, &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}, updater.InstalledVersion())

	// Another lazygit, running from a different executable, doesn't see the
	// installed version
	otherUpdater := newTestUpdater(t, server, "0.66.2", Installation{Kind: InstallKindRelease, ExecutablePath: "/other/lazygit"})
	otherUpdater.appState = updater.appState
	assert.Nil(t, otherUpdater.InstalledVersion())

	// After a restart, the installed version is the running one
	restartedUpdater := newTestUpdater(t, server, "0.67.0", Installation{Kind: InstallKindRelease, ExecutablePath: executablePath})
	restartedUpdater.appState = updater.appState
	assert.Nil(t, restartedUpdater.InstalledVersion())
}

func newTestUpdater(t *testing.T, server *httptest.Server, currentVersion string, installation Installation) *Updater {
	t.Helper()

	version, err := types.ParseVersionNumber(currentVersion)
	assert.NoError(t, err)

	return &Updater{
		Common:         common.NewDummyCommon(),
		appState:       &fakeAppStateStore{state: &config.AppState{}},
		installation:   installation,
		currentVersion: version,
		httpClient:     server.Client(),
		repoURL:        server.URL,
		goos:           "linux",
		goarch:         "amd64",
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
