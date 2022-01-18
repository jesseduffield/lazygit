package git_commands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestBranchGetCommitDifferences(t *testing.T) {
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
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildBranchCommands(commonDeps{runner: s.runner})
			pushables, pullables := instance.GetCommitDifferences("HEAD", "@{u}")
			assert.EqualValues(t, s.expectedPushables, pushables)
			assert.EqualValues(t, s.expectedPullables, pullables)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestBranchNewBranch(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		Expect(`git checkout -b "test" "master"`, "", nil)
	instance := buildBranchCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.New("test", "master"))
	runner.CheckForMissingCalls()
}

func TestBranchDeleteBranch(t *testing.T) {
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
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildBranchCommands(commonDeps{runner: s.runner})

			s.test(instance.Delete("test", s.force))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestBranchMerge(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		Expect(`git merge --no-edit "test"`, "", nil)
	instance := buildBranchCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.Merge("test", MergeOpts{}))
	runner.CheckForMissingCalls()
}

func TestBranchCheckout(t *testing.T) {
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
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildBranchCommands(commonDeps{runner: s.runner})
			s.test(instance.Checkout("test", CheckoutOptions{Force: s.force}))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestBranchGetBranchGraph(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).ExpectGitArgs([]string{
		"log", "--graph", "--color=always", "--abbrev-commit", "--decorate", "--date=relative", "--pretty=medium", "test", "--",
	}, "", nil)
	instance := buildBranchCommands(commonDeps{runner: runner})
	_, err := instance.GetGraph("test")
	assert.NoError(t, err)
}

func TestBranchGetAllBranchGraph(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).ExpectGitArgs([]string{
		"log", "--graph", "--all", "--color=always", "--abbrev-commit", "--decorate", "--date=relative", "--pretty=medium",
	}, "", nil)
	instance := buildBranchCommands(commonDeps{runner: runner})
	err := instance.AllBranchesLogCmdObj().Run()
	assert.NoError(t, err)
}

func TestBranchCurrentBranchName(t *testing.T) {
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
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildBranchCommands(commonDeps{runner: s.runner})
			s.test(instance.CurrentBranchName())
			s.runner.CheckForMissingCalls()
		})
	}
}
