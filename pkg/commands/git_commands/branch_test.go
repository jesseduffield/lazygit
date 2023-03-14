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
		Expect(`git checkout -b "test" "refs/heads/master"`, "", nil)
	instance := buildBranchCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.New("test", "refs/heads/master"))
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

func TestBranchCurrentBranchInfo(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(BranchInfo, error)
	}

	scenarios := []scenario{
		{
			"says we are on the master branch if we are",
			oscommands.NewFakeRunner(t).Expect(`git symbolic-ref --short HEAD`, "master", nil),
			func(info BranchInfo, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "master", info.RefName)
				assert.EqualValues(t, "master", info.DisplayName)
				assert.False(t, info.DetachedHead)
			},
		},
		{
			"falls back to git `git branch --points-at=HEAD` if symbolic-ref fails",
			oscommands.NewFakeRunner(t).
				Expect(`git symbolic-ref --short HEAD`, "", errors.New("error")).
				Expect(`git branch --points-at=HEAD --format="%(HEAD)%00%(objectname)%00%(refname)"`, "*\x006f71c57a8d4bd6c11399c3f55f42c815527a73a4\x00(HEAD detached at 6f71c57a)\n", nil),
			func(info BranchInfo, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "6f71c57a8d4bd6c11399c3f55f42c815527a73a4", info.RefName)
				assert.EqualValues(t, "(HEAD detached at 6f71c57a)", info.DisplayName)
				assert.True(t, info.DetachedHead)
			},
		},
		{
			"handles a detached head (LANG=zh_CN.UTF-8)",
			oscommands.NewFakeRunner(t).
				Expect(`git symbolic-ref --short HEAD`, "", errors.New("error")).
				Expect(
					`git branch --points-at=HEAD --format="%(HEAD)%00%(objectname)%00%(refname)"`,
					"*\x00679b0456f3db7c505b398def84e7d023e5b55a8d\x00（头指针在 679b0456 分离）\n"+
						" \x00679b0456f3db7c505b398def84e7d023e5b55a8d\x00refs/heads/master\n",
					nil),
			func(info BranchInfo, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "679b0456f3db7c505b398def84e7d023e5b55a8d", info.RefName)
				assert.EqualValues(t, "（头指针在 679b0456 分离）", info.DisplayName)
				assert.True(t, info.DetachedHead)
			},
		},
		{
			"bubbles up error if there is one",
			oscommands.NewFakeRunner(t).
				Expect(`git symbolic-ref --short HEAD`, "", errors.New("error")).
				Expect(`git branch --points-at=HEAD --format="%(HEAD)%00%(objectname)%00%(refname)"`, "", errors.New("error")),
			func(info BranchInfo, err error) {
				assert.Error(t, err)
				assert.EqualValues(t, "", info.RefName)
				assert.EqualValues(t, "", info.DisplayName)
				assert.False(t, info.DetachedHead)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildBranchCommands(commonDeps{runner: s.runner})
			s.test(instance.CurrentBranchInfo())
			s.runner.CheckForMissingCalls()
		})
	}
}
