package git_commands

import (
	"testing"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestSyncPush(t *testing.T) {
	type scenario struct {
		testName string
		opts     PushOpts
		test     func(*oscommands.CmdObj, error)
	}

	scenarios := []scenario{
		{
			testName: "Push with force disabled",
			opts:     PushOpts{ForceWithLease: false},
			test: func(cmdObj *oscommands.CmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force-with-lease enabled",
			opts:     PushOpts{ForceWithLease: true},
			test: func(cmdObj *oscommands.CmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push", "--force-with-lease"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force enabled",
			opts:     PushOpts{Force: true},
			test: func(cmdObj *oscommands.CmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push", "--force"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force disabled, upstream supplied",
			opts: PushOpts{
				ForceWithLease: false,
				CurrentBranch:  "master",
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
			},
			test: func(cmdObj *oscommands.CmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push", "origin", "refs/heads/master:master"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force disabled, setting upstream",
			opts: PushOpts{
				ForceWithLease: false,
				CurrentBranch:  "master-local",
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj *oscommands.CmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push", "--set-upstream", "origin", "refs/heads/master-local:master"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force-with-lease enabled, setting upstream",
			opts: PushOpts{
				ForceWithLease: true,
				CurrentBranch:  "master",
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj *oscommands.CmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push", "--force-with-lease", "--set-upstream", "origin", "refs/heads/master:master"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with remote branch but no origin",
			opts: PushOpts{
				ForceWithLease: true,
				UpstreamRemote: "",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj *oscommands.CmdObj, err error) {
				assert.Error(t, err)
				assert.EqualValues(t, "Must specify a remote if specifying a branch", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildSyncCommands(commonDeps{})
			task := gocui.NewFakeTask()
			cmdObj, err := instance.PushCmdObj(task, s.opts)
			if err == nil {
				assert.True(t, cmdObj.ShouldLog())
				assert.Equal(t, cmdObj.GetCredentialStrategy(), oscommands.PROMPT)
				assert.False(t, cmdObj.ShouldSuppressOutputUnlessError())
			}
			s.test(cmdObj, err)
		})
	}
}

func TestSyncFetch(t *testing.T) {
	type scenario struct {
		testName       string
		fetchAllConfig bool
		test           func(*oscommands.CmdObj)
	}

	scenarios := []scenario{
		{
			testName:       "Fetch in foreground (all=false)",
			fetchAllConfig: false,
			test: func(cmdObj *oscommands.CmdObj) {
				assert.True(t, cmdObj.ShouldLog())
				assert.Equal(t, cmdObj.GetCredentialStrategy(), oscommands.PROMPT)
				assert.False(t, cmdObj.ShouldSuppressOutputUnlessError())
				assert.Equal(t, cmdObj.Args(), []string{"git", "fetch", "--no-write-fetch-head"})
			},
		},
		{
			testName:       "Fetch in foreground (all=true)",
			fetchAllConfig: true,
			test: func(cmdObj *oscommands.CmdObj) {
				assert.True(t, cmdObj.ShouldLog())
				assert.Equal(t, cmdObj.GetCredentialStrategy(), oscommands.PROMPT)
				assert.False(t, cmdObj.ShouldSuppressOutputUnlessError())
				assert.Equal(t, cmdObj.Args(), []string{"git", "fetch", "--all", "--no-write-fetch-head"})
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildSyncCommands(commonDeps{})
			instance.UserConfig().Git.FetchAll = s.fetchAllConfig
			task := gocui.NewFakeTask()
			s.test(instance.FetchCmdObj(task))
		})
	}
}

func TestSyncFetchBackground(t *testing.T) {
	type scenario struct {
		testName       string
		fetchAllConfig bool
		test           func(*oscommands.CmdObj)
	}

	scenarios := []scenario{
		{
			testName:       "Fetch in background (all=false)",
			fetchAllConfig: false,
			test: func(cmdObj *oscommands.CmdObj) {
				assert.False(t, cmdObj.ShouldLog())
				assert.Equal(t, cmdObj.GetCredentialStrategy(), oscommands.FAIL)
				assert.True(t, cmdObj.ShouldSuppressOutputUnlessError())
				assert.Equal(t, cmdObj.Args(), []string{"git", "fetch", "--no-write-fetch-head"})
			},
		},
		{
			testName:       "Fetch in background (all=true)",
			fetchAllConfig: true,
			test: func(cmdObj *oscommands.CmdObj) {
				assert.False(t, cmdObj.ShouldLog())
				assert.Equal(t, cmdObj.GetCredentialStrategy(), oscommands.FAIL)
				assert.True(t, cmdObj.ShouldSuppressOutputUnlessError())
				assert.Equal(t, cmdObj.Args(), []string{"git", "fetch", "--all", "--no-write-fetch-head"})
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildSyncCommands(commonDeps{})
			instance.UserConfig().Git.FetchAll = s.fetchAllConfig
			s.test(instance.FetchBackgroundCmdObj())
		})
	}
}
