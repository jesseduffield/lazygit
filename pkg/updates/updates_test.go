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
		})
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
		httpClient:        server.Client(),
		repoURL:           server.URL,
		goos:              "linux",
		goarch:            "amd64",
	}
}

type fakeAppStateStore struct {
	state *config.AppState
}

func (self *fakeAppStateStore) UpdateAppState(update func(state *config.AppState)) error {
	update(self.state)
	return nil
}
