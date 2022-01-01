package commands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestGitCommandGetCommitDifferences(t *testing.T) {
	type scenario struct {
		testName          string
		runner            *oscommands.FakeCmdObjRunner
		expectedPushables string
		expectedPullables string
	}

	scenarios := []scenario{
		{
			"Can't retrieve pushable count",
			oscommands.NewFakeRunner(t).
				Expect("git rev-list @{u}..HEAD --count", "", errors.New("error")),
			"?", "?",
		},
		{
			"Can't retrieve pullable count",
			oscommands.NewFakeRunner(t).
				Expect("git rev-list @{u}..HEAD --count", "1\n", nil).
				Expect("git rev-list HEAD..@{u} --count", "", errors.New("error")),
			"?", "?",
		},
		{
			"Retrieve pullable and pushable count",
			oscommands.NewFakeRunner(t).
				Expect("git rev-list @{u}..HEAD --count", "1\n", nil).
				Expect("git rev-list HEAD..@{u} --count", "2\n", nil),
			"1", "2",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			pushables, pullables := gitCmd.Branch.GetCommitDifferences("HEAD", "@{u}")
			assert.EqualValues(t, s.expectedPushables, pushables)
			assert.EqualValues(t, s.expectedPullables, pullables)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandNewBranch(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		Expect(`git checkout -b "test" "master"`, "", nil)
	gitCmd := NewDummyGitCommandWithRunner(runner)

	assert.NoError(t, gitCmd.Branch.New("test", "master"))
	runner.CheckForMissingCalls()
}

func TestGitCommandDeleteBranch(t *testing.T) {
	type scenario struct {
		testName string
		force    bool
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			"Delete a branch",
			false,
			oscommands.NewFakeRunner(t).Expect(`git branch -d "test"`, "", nil),
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Force delete a branch",
			true,
			oscommands.NewFakeRunner(t).Expect(`git branch -D "test"`, "", nil),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)

			s.test(gitCmd.Branch.Delete("test", s.force))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandMerge(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		Expect(`git merge --no-edit "test"`, "", nil)
	gitCmd := NewDummyGitCommandWithRunner(runner)

	assert.NoError(t, gitCmd.Branch.Merge("test", MergeOpts{}))
	runner.CheckForMissingCalls()
}

func TestGitCommandCheckout(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
		force    bool
	}

	scenarios := []scenario{
		{
			"Checkout",
			oscommands.NewFakeRunner(t).Expect(`git checkout "test"`, "", nil),
			func(err error) {
				assert.NoError(t, err)
			},
			false,
		},
		{
			"Checkout forced",
			oscommands.NewFakeRunner(t).Expect(`git checkout --force "test"`, "", nil),
			func(err error) {
				assert.NoError(t, err)
			},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.Branch.Checkout("test", CheckoutOptions{Force: s.force}))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestGitCommandGetBranchGraph(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).ExpectGitArgs([]string{
		"log", "--graph", "--color=always", "--abbrev-commit", "--decorate", "--date=relative", "--pretty=medium", "test", "--",
	}, "", nil)
	gitCmd := NewDummyGitCommandWithRunner(runner)
	_, err := gitCmd.Branch.GetGraph("test")
	assert.NoError(t, err)
}

func TestGitCommandGetAllBranchGraph(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).ExpectGitArgs([]string{
		"log", "--graph", "--all", "--color=always", "--abbrev-commit", "--decorate", "--date=relative", "--pretty=medium",
	}, "", nil)
	gitCmd := NewDummyGitCommandWithRunner(runner)
	cmdStr := gitCmd.UserConfig.Git.AllBranchesLogCmd
	_, err := gitCmd.Cmd.New(cmdStr).RunWithOutput()
	assert.NoError(t, err)
}

func TestGitCommandCurrentBranchName(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(string, string, error)
	}

	scenarios := []scenario{
		{
			"says we are on the master branch if we are",
			oscommands.NewFakeRunner(t).Expect(`git symbolic-ref --short HEAD`, "master", nil),
			func(name string, displayname string, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "master", name)
				assert.EqualValues(t, "master", displayname)
			},
		},
		{
			"falls back to git `git branch --contains` if symbolic-ref fails",
			oscommands.NewFakeRunner(t).
				Expect(`git symbolic-ref --short HEAD`, "", errors.New("error")).
				Expect(`git branch --contains`, "* master", nil),
			func(name string, displayname string, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "master", name)
				assert.EqualValues(t, "master", displayname)
			},
		},
		{
			"handles a detached head",
			oscommands.NewFakeRunner(t).
				Expect(`git symbolic-ref --short HEAD`, "", errors.New("error")).
				Expect(`git branch --contains`, "* (HEAD detached at 123abcd)", nil),
			func(name string, displayname string, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "123abcd", name)
				assert.EqualValues(t, "(HEAD detached at 123abcd)", displayname)
			},
		},
		{
			"bubbles up error if there is one",
			oscommands.NewFakeRunner(t).
				Expect(`git symbolic-ref --short HEAD`, "", errors.New("error")).
				Expect(`git branch --contains`, "", errors.New("error")),
			func(name string, displayname string, err error) {
				assert.Error(t, err)
				assert.EqualValues(t, "", name)
				assert.EqualValues(t, "", displayname)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommandWithRunner(s.runner)
			s.test(gitCmd.Branch.CurrentBranchName())
			s.runner.CheckForMissingCalls()
		})
	}
}
