package updates

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/common"
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

func newTestUpdater(t *testing.T, server *httptest.Server, currentVersion string, installation Installation) *Updater {
	t.Helper()

	version, err := types.ParseVersionNumber(currentVersion)
	assert.NoError(t, err)

	return &Updater{
		Common:         common.NewDummyCommon(),
		installation:   installation,
		currentVersion: version,
		httpClient:     server.Client(),
		repoURL:        server.URL,
	}
}
