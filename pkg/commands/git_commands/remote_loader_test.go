package git_commands

import (
	"errors"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestGetRemotesFromConfig(t *testing.T) {
	configArgs := []string{"config", "--local", "--get-regexp", `^remote\.[^.]+\.(url|pushurl)$`}

	scenarios := []struct {
		testName        string
		runner          *oscommands.FakeCmdObjRunner
		expectedRemotes []*models.Remote
	}{
		{
			testName: "no remotes configured",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(configArgs, "", errors.New("exit status 1")),
			expectedRemotes: nil,
		},
		{
			testName: "single remote with one url",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(configArgs,
					"remote.origin.url https://github.com/foo/bar.git\n",
					nil),
			expectedRemotes: []*models.Remote{
				{Name: "origin", Urls: []string{"https://github.com/foo/bar.git"}},
			},
		},
		{
			testName: "mirror remote with multiple urls",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(configArgs,
					"remote.origin.url https://github.com/foo/bar.git\n"+
						"remote.origin.url git@github.com:foo/bar.git\n",
					nil),
			expectedRemotes: []*models.Remote{
				{Name: "origin", Urls: []string{
					"https://github.com/foo/bar.git",
					"git@github.com:foo/bar.git",
				}},
			},
		},
		{
			testName: "remote with both url and pushurl",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(configArgs,
					"remote.origin.url https://github.com/foo/bar.git\n"+
						"remote.origin.pushurl git@github.com:foo/bar.git\n",
					nil),
			expectedRemotes: []*models.Remote{
				{
					Name:     "origin",
					Urls:     []string{"https://github.com/foo/bar.git"},
					PushUrls: []string{"git@github.com:foo/bar.git"},
				},
			},
		},
		{
			testName: "multiple remotes",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(configArgs,
					"remote.origin.url https://github.com/foo/bar.git\n"+
						"remote.upstream.url https://github.com/baz/bar.git\n"+
						"remote.upstream.pushurl git@github.com:baz/bar.git\n",
					nil),
			expectedRemotes: []*models.Remote{
				{Name: "origin", Urls: []string{"https://github.com/foo/bar.git"}},
				{
					Name:     "upstream",
					Urls:     []string{"https://github.com/baz/bar.git"},
					PushUrls: []string{"git@github.com:baz/bar.git"},
				},
			},
		},
		{
			testName: "remote name containing dots is preserved",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(configArgs,
					"remote.my.fork.url https://github.com/foo/bar.git\n",
					nil),
			expectedRemotes: []*models.Remote{
				{Name: "my.fork", Urls: []string{"https://github.com/foo/bar.git"}},
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.testName, func(t *testing.T) {
			loader := &RemoteLoader{
				Common: common.NewDummyCommon(),
				cmd:    oscommands.NewDummyCmdObjBuilder(scenario.runner),
			}

			// map iteration order is non-deterministic, so compare unordered
			assert.ElementsMatch(t, scenario.expectedRemotes, loader.getRemotesFromConfig())

			scenario.runner.CheckForMissingCalls()
		})
	}
}
