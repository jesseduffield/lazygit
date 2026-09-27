package git_commands

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// A branch as the tests below describe it: a name, the hash of its tip, and the
// committer date of its tip
type testRef struct {
	name string
	hash string
	date string
}

func buildTestRefs(refs []testRef) ([]*models.Branch, map[string]refTip) {
	branches := lo.Map(refs, func(ref testRef, _ int) *models.Branch {
		return &models.Branch{Name: ref.name, CommitHash: ref.hash}
	})
	tips := make(map[string]refTip, len(refs))
	for _, ref := range refs {
		tips["refs/heads/"+ref.name] = refTip{hash: ref.hash, committerDate: ref.date}
	}
	return branches, tips
}

func branchNames(branches []*models.Branch) []string {
	return lo.Map(branches, func(branch *models.Branch, _ int) string {
		return branch.Name
	})
}

func TestGroupRefsWithEqualDates(t *testing.T) {
	scenarios := []struct {
		testName string
		refs     []testRef
		expected [][]string
	}{
		{
			testName: "no refs",
			refs:     []testRef{},
			expected: [][]string{},
		},
		{
			testName: "all dates distinct",
			refs: []testRef{
				{name: "a", hash: "1", date: "300"},
				{name: "b", hash: "2", date: "200"},
				{name: "c", hash: "3", date: "100"},
			},
			expected: [][]string{},
		},
		{
			testName: "all dates equal",
			refs: []testRef{
				{name: "a", hash: "1", date: "100"},
				{name: "b", hash: "2", date: "100"},
				{name: "c", hash: "3", date: "100"},
			},
			expected: [][]string{{"a", "b", "c"}},
		},
		{
			testName: "two groups, with a lone ref between and after them",
			refs: []testRef{
				{name: "a", hash: "1", date: "300"},
				{name: "b", hash: "2", date: "300"},
				{name: "c", hash: "3", date: "200"},
				{name: "d", hash: "4", date: "100"},
				{name: "e", hash: "5", date: "100"},
				{name: "f", hash: "6", date: "50"},
			},
			expected: [][]string{{"a", "b"}, {"d", "e"}},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			branches, tips := buildTestRefs(s.refs)
			groups := groupRefsWithEqualDates(branches, (*models.Branch).FullRefName, tips)
			assert.Equal(t, s.expected, lo.Map(groups,
				func(group []*models.Branch, _ int) []string {
					return branchNames(group)
				}))
		})
	}
}

func TestSortRefsWithEqualDatesByAncestry(t *testing.T) {
	// A stack of three branches, all committed within the same second, as git
	// returns them: sorted by name. bottom is the base of middle, which is the
	// base of top.
	stack := []testRef{
		{name: "middle", hash: "m", date: "100"},
		{name: "top", hash: "t", date: "100"},
		{name: "bottom", hash: "b", date: "100"},
	}

	// One %(ahead-behind:<hash>) field per tip, in the order the tips appear in
	// the branch list
	stackOutput := "refs/heads/middle\x000 0\x000 1\x001 0\n" +
		"refs/heads/top\x001 0\x000 0\x002 0\n" +
		"refs/heads/bottom\x000 1\x000 2\x000 0\n"

	stackArgs := []string{
		"for-each-ref",
		"--format=%(refname)%00%(ahead-behind:m)%00%(ahead-behind:t)%00%(ahead-behind:b)",
		"refs/heads/middle", "refs/heads/top", "refs/heads/bottom",
	}

	scenarios := []struct {
		testName      string
		refs          []testRef
		gitVersion    *GitVersion
		expectedArgs  []string
		output        string
		outputErr     error
		expectedOrder []string
		expectedErr   string
	}{
		{
			testName:      "a stack of branches is sorted from the top down",
			refs:          stack,
			gitVersion:    &GitVersion{2, 41, 0, ""},
			expectedArgs:  stackArgs,
			output:        stackOutput,
			expectedOrder: []string{"top", "middle", "bottom"},
		},
		{
			testName:      "git too old for %(ahead-behind:...), so nothing to do",
			refs:          stack,
			gitVersion:    &GitVersion{2, 40, 0, ""},
			expectedOrder: []string{"middle", "top", "bottom"},
		},
		{
			testName: "no two refs share a date, so nothing to do",
			refs: []testRef{
				{name: "middle", hash: "m", date: "300"},
				{name: "top", hash: "t", date: "200"},
				{name: "bottom", hash: "b", date: "100"},
			},
			gitVersion:    &GitVersion{2, 41, 0, ""},
			expectedOrder: []string{"middle", "top", "bottom"},
		},
		{
			testName:      "the command fails",
			refs:          stack,
			gitVersion:    &GitVersion{2, 41, 0, ""},
			expectedArgs:  stackArgs,
			outputErr:     errors.New("fatal: failed to find 'm'"),
			expectedOrder: []string{"middle", "top", "bottom"},
			expectedErr:   "fatal: failed to find 'm'",
		},
		{
			testName: "a branch unrelated to the stack is placed by name",
			refs: []testRef{
				{name: "middle", hash: "m", date: "100"},
				{name: "other", hash: "o", date: "100"},
				{name: "top", hash: "t", date: "100"},
				{name: "bottom", hash: "b", date: "100"},
			},
			gitVersion: &GitVersion{2, 41, 0, ""},
			expectedArgs: []string{
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:m)%00%(ahead-behind:o)%00%(ahead-behind:t)%00%(ahead-behind:b)",
				"refs/heads/middle", "refs/heads/other", "refs/heads/top", "refs/heads/bottom",
			},
			output: "refs/heads/middle\x000 0\x002 1\x000 1\x001 0\n" +
				"refs/heads/other\x001 2\x000 0\x001 3\x001 1\n" +
				"refs/heads/top\x001 0\x003 1\x000 0\x002 0\n" +
				"refs/heads/bottom\x000 1\x001 1\x000 2\x000 0\n",
			expectedOrder: []string{"other", "top", "middle", "bottom"},
		},
		{
			testName: "a stack stays together when a branch sorts into the middle of it",
			refs: []testRef{
				{name: "add-tests", hash: "t", date: "100"},
				{name: "cleanup", hash: "c", date: "100"},
				{name: "fix-parser", hash: "p", date: "100"},
			},
			gitVersion: &GitVersion{2, 41, 0, ""},
			expectedArgs: []string{
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:t)%00%(ahead-behind:c)%00%(ahead-behind:p)",
				"refs/heads/add-tests", "refs/heads/cleanup", "refs/heads/fix-parser",
			},
			// add-tests is based on fix-parser, and cleanup is on a line of its
			// own, but its name sorts between the two
			output: "refs/heads/add-tests\x000 0\x002 1\x001 0\n" +
				"refs/heads/cleanup\x001 2\x000 0\x001 1\n" +
				"refs/heads/fix-parser\x000 1\x001 1\x000 0\n",
			expectedOrder: []string{"add-tests", "fix-parser", "cleanup"},
		},
		{
			testName: "each group is sorted on its own",
			refs: []testRef{
				{name: "middle", hash: "m", date: "200"},
				{name: "top", hash: "t", date: "200"},
				{name: "lone", hash: "l", date: "150"},
				{name: "base", hash: "a", date: "100"},
				{name: "derived", hash: "d", date: "100"},
			},
			gitVersion: &GitVersion{2, 41, 0, ""},
			expectedArgs: []string{
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:m)%00%(ahead-behind:t)%00%(ahead-behind:a)%00%(ahead-behind:d)",
				"refs/heads/middle", "refs/heads/top", "refs/heads/base", "refs/heads/derived",
			},
			output: "refs/heads/middle\x000 0\x000 1\x002 0\x001 0\n" +
				"refs/heads/top\x001 0\x000 0\x003 0\x002 0\n" +
				"refs/heads/base\x000 2\x000 3\x000 0\x000 1\n" +
				"refs/heads/derived\x000 1\x000 2\x001 0\x000 0\n",
			expectedOrder: []string{"top", "middle", "lone", "derived", "base"},
		},
		{
			testName: "a group whose branches are all on one commit is left alone",
			refs: []testRef{
				{name: "a", hash: "x", date: "100"},
				{name: "b", hash: "x", date: "100"},
			},
			gitVersion:    &GitVersion{2, 41, 0, ""},
			expectedOrder: []string{"a", "b"},
		},
		{
			testName: "branches on the same commit are asked about once",
			refs: []testRef{
				{name: "mirror-1", hash: "s", date: "100"},
				{name: "mirror-2", hash: "s", date: "100"},
				{name: "stack-bottom", hash: "u", date: "100"},
				{name: "stack-top", hash: "t", date: "100"},
			},
			gitVersion: &GitVersion{2, 41, 0, ""},
			expectedArgs: []string{
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:s)%00%(ahead-behind:u)%00%(ahead-behind:t)",
				"refs/heads/mirror-1", "refs/heads/stack-bottom", "refs/heads/stack-top",
			},
			output: "refs/heads/mirror-1\x000 0\x001 1\x001 2\n" +
				"refs/heads/stack-bottom\x001 1\x000 0\x000 1\n" +
				"refs/heads/stack-top\x002 1\x001 0\x000 0\n",
			expectedOrder: []string{"mirror-1", "mirror-2", "stack-top", "stack-bottom"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t)
			if s.expectedArgs != nil {
				runner.ExpectGitArgs(s.expectedArgs, s.output, s.outputErr)
			}
			gitCommon := buildGitCommon(commonDeps{runner: runner, gitVersion: s.gitVersion})

			branches, tips := buildTestRefs(s.refs)
			err := sortRefsWithEqualDatesByAncestry(gitCommon.cmd, gitCommon.version,
				branches, (*models.Branch).FullRefName, tips)

			if s.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, s.expectedErr)
			}
			assert.Equal(t, s.expectedOrder, branchNames(branches))
			runner.CheckForMissingCalls()
		})
	}
}

func TestSortRefsWithEqualDatesByAncestry_TooManyTips(t *testing.T) {
	refs := lo.Map(lo.Range(maxTipsForAncestrySorting+1), func(i int, _ int) testRef {
		return testRef{
			name: fmt.Sprintf("branch-%03d", i),
			hash: fmt.Sprintf("hash-%03d", i),
			date: "100",
		}
	})

	// The runner fails the test if the command runs at all
	runner := oscommands.NewFakeRunner(t)
	gitCommon := buildGitCommon(commonDeps{runner: runner, gitVersion: &GitVersion{2, 41, 0, ""}})

	branches, tips := buildTestRefs(refs)
	err := sortRefsWithEqualDatesByAncestry(gitCommon.cmd, gitCommon.version,
		branches, (*models.Branch).FullRefName, tips)

	assert.NoError(t, err)
	assert.Equal(t, lo.Map(refs, func(ref testRef, _ int) string { return ref.name }),
		branchNames(branches))
	runner.CheckForMissingCalls()
}

// A repository with many remotes has the same branches under each of them, so
// a group can hold far more refs than the commits they point at
func TestSortRefsWithEqualDatesByAncestry_ManyRefsOnTwoTips(t *testing.T) {
	refs := lo.Map(lo.Range(4*maxTipsForAncestrySorting), func(i int, _ int) testRef {
		return testRef{
			name: fmt.Sprintf("branch-%03d", i),
			hash: lo.Ternary(i%2 == 0, "bottom", "top"),
			date: "100",
		}
	})

	runner := oscommands.NewFakeRunner(t).ExpectGitArgs([]string{
		"for-each-ref",
		"--format=%(refname)%00%(ahead-behind:bottom)%00%(ahead-behind:top)",
		"refs/heads/branch-000", "refs/heads/branch-001",
	},
		"refs/heads/branch-000\x000 0\x000 1\n"+
			"refs/heads/branch-001\x001 0\x000 0\n", nil)
	gitCommon := buildGitCommon(commonDeps{runner: runner, gitVersion: &GitVersion{2, 41, 0, ""}})

	branches, tips := buildTestRefs(refs)
	err := sortRefsWithEqualDatesByAncestry(gitCommon.cmd, gitCommon.version,
		branches, (*models.Branch).FullRefName, tips)

	assert.NoError(t, err)
	// the branches on the top commit first, each half still ordered by name
	expected := append(
		lo.FilterMap(refs, func(ref testRef, _ int) (string, bool) {
			return ref.name, ref.hash == "top"
		}),
		lo.FilterMap(refs, func(ref testRef, _ int) (string, bool) {
			return ref.name, ref.hash == "bottom"
		})...)
	assert.Equal(t, expected, branchNames(branches))
	runner.CheckForMissingCalls()
}
