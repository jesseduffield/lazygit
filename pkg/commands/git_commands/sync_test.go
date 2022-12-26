package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestSyncPush(t *testing.T) {
	type scenario struct {
		testName string
		version  *GitVersion
		opts     PushOpts
		test     func(oscommands.ICmdObj, error)
	}

	scenarios := []scenario{
		{
			testName: "Push with force disabled",
			version:  &GitVersion{2, 29, 3, ""},
			opts:     PushOpts{Force: false},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.ToString(), "git push")
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force enabled",
			version:  &GitVersion{2, 29, 3, ""},
			opts:     PushOpts{Force: true},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.ToString(), "git push --force-with-lease")
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force enabled (>= 2.30.0)",
			version:  &GitVersion{2, 30, 0, ""},
			opts:     PushOpts{Force: true},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.ToString(), "git push --force-with-lease --force-if-includes")
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force disabled, upstream supplied",
			version:  &GitVersion{2, 29, 3, ""},
			opts: PushOpts{
				Force:          false,
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
			},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.ToString(), `git push "origin" "master"`)
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force disabled, setting upstream",
			version:  &GitVersion{2, 29, 3, ""},
			opts: PushOpts{
				Force:          false,
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.ToString(), `git push --set-upstream "origin" "master"`)
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force enabled, setting upstream",
			version:  &GitVersion{2, 29, 3, ""},
			opts: PushOpts{
				Force:          true,
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.ToString(), `git push --force-with-lease --set-upstream "origin" "master"`)
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force enabled, setting upstream (>= 2.30.0)",
			version:  &GitVersion{2, 30, 0, ""},
			opts: PushOpts{
				Force:          true,
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.ToString(), `git push --force-with-lease --force-if-includes --set-upstream "origin" "master"`)
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with remote branch but no origin",
			version:  &GitVersion{2, 29, 3, ""},
			opts: PushOpts{
				Force:          true,
				UpstreamRemote: "",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Error(t, err)
				assert.EqualValues(t, "Must specify a remote if specifying a branch", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildSyncCommands(commonDeps{gitVersion: s.version})
			s.test(instance.PushCmdObj(s.opts))
		})
	}
}
