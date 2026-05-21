package git_commands

// "*|feat/detect-purge|origin/feat/detect-purge|[ahead 1]"
import (
	"strconv"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestObtainBranch(t *testing.T) {
	type scenario struct {
		testName                 string
		input                    []string
		storeCommitDateAsRecency bool
		expectedBranch           *models.Branch
	}

	// Use a time stamp of 2 1/2 hours ago, resulting in a recency string of "2h"
	now := time.Now().Unix()
	timeStamp := strconv.Itoa(int(now - 2.5*60*60))

	scenarios := []scenario{
		{
			testName:                 "TrimHeads",
			input:                    []string{"", "heads/a_branch", "", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:          "a_branch",
				AheadForPull:  "?",
				BehindForPull: "?",
				AheadForPush:  "?",
				BehindForPush: "?",
				Head:          false,
				Subject:       "subject",
				CommitHash:    "123",
			},
		},
		{
			testName:                 "NoUpstream",
			input:                    []string{"", "a_branch", "", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:          "a_branch",
				AheadForPull:  "?",
				BehindForPull: "?",
				AheadForPush:  "?",
				BehindForPush: "?",
				Head:          false,
				Subject:       "subject",
				CommitHash:    "123",
			},
		},
		{
			testName:                 "IsHead",
			input:                    []string{"*", "a_branch", "", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:          "a_branch",
				AheadForPull:  "?",
				BehindForPull: "?",
				AheadForPush:  "?",
				BehindForPush: "?",
				Head:          true,
				Subject:       "subject",
				CommitHash:    "123",
			},
		},
		{
			testName:                 "IsBehindAndAhead",
			input:                    []string{"", "a_branch", "a_remote/a_branch", "[behind 2, ahead 3]", "[behind 2, ahead 3]", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:          "a_branch",
				AheadForPull:  "3",
				BehindForPull: "2",
				AheadForPush:  "3",
				BehindForPush: "2",
				Head:          false,
				Subject:       "subject",
				CommitHash:    "123",
			},
		},
		{
			testName:                 "RemoteBranchIsGone",
			input:                    []string{"", "a_branch", "a_remote/a_branch", "[gone]", "[gone]", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
			expectedBranch: &models.Branch{
				Name:          "a_branch",
				UpstreamGone:  true,
				AheadForPull:  "?",
				BehindForPull: "?",
				AheadForPush:  "?",
				BehindForPush: "?",
				Head:          false,
				Subject:       "subject",
				CommitHash:    "123",
			},
		},
		{
			testName:                 "WithCommitDateAsRecency",
			input:                    []string{"", "a_branch", "", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: true,
			expectedBranch: &models.Branch{
				Name:          "a_branch",
				Recency:       "2h",
				AheadForPull:  "?",
				BehindForPull: "?",
				AheadForPush:  "?",
				BehindForPush: "?",
				Head:          false,
				Subject:       "subject",
				CommitHash:    "123",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			branch := obtainBranch(s.input, s.storeCommitDateAsRecency)
			assert.EqualValues(t, s.expectedBranch, branch)
		})
	}
}

func TestParseAheadBehindForEachRefOutput(t *testing.T) {
	type scenario struct {
		testName string
		input    string
		numBases int
		expected []branchAheadBehind
	}

	scenarios := []scenario{
		{
			testName: "single branch single base",
			input:    "refs/heads/feat\x002 5\n",
			numBases: 1,
			expected: []branchAheadBehind{
				{
					refName:      "refs/heads/feat",
					aheadBehinds: []aheadBehind{{ahead: 2, behind: 5, valid: true}},
				},
			},
		},
		{
			testName: "multiple branches multiple bases",
			input: "refs/heads/feat\x002 5\x0010 1\n" +
				"refs/heads/main\x000 0\x000 0\n",
			numBases: 2,
			expected: []branchAheadBehind{
				{
					refName: "refs/heads/feat",
					aheadBehinds: []aheadBehind{
						{ahead: 2, behind: 5, valid: true},
						{ahead: 10, behind: 1, valid: true},
					},
				},
				{
					refName: "refs/heads/main",
					aheadBehinds: []aheadBehind{
						{ahead: 0, behind: 0, valid: true},
						{ahead: 0, behind: 0, valid: true},
					},
				},
			},
		},
		{
			testName: "empty ahead-behind field for unreachable base",
			input:    "refs/heads/feat\x00\x002 5\n",
			numBases: 2,
			expected: []branchAheadBehind{
				{
					refName: "refs/heads/feat",
					aheadBehinds: []aheadBehind{
						{valid: false},
						{ahead: 2, behind: 5, valid: true},
					},
				},
			},
		},
		{
			testName: "ref name containing slashes and dashes",
			input:    "refs/heads/feat/foo-bar\x001 2\n",
			numBases: 1,
			expected: []branchAheadBehind{
				{
					refName:      "refs/heads/feat/foo-bar",
					aheadBehinds: []aheadBehind{{ahead: 1, behind: 2, valid: true}},
				},
			},
		},
		{
			testName: "trailing newline and blank lines are ignored",
			input:    "refs/heads/feat\x001 2\n\n",
			numBases: 1,
			expected: []branchAheadBehind{
				{
					refName:      "refs/heads/feat",
					aheadBehinds: []aheadBehind{{ahead: 1, behind: 2, valid: true}},
				},
			},
		},
		{
			testName: "line with wrong column count is skipped",
			input: "refs/heads/good\x001 2\n" +
				"refs/heads/bad\n" +
				"refs/heads/also_good\x003 4\n",
			numBases: 1,
			expected: []branchAheadBehind{
				{
					refName:      "refs/heads/good",
					aheadBehinds: []aheadBehind{{ahead: 1, behind: 2, valid: true}},
				},
				{
					refName:      "refs/heads/also_good",
					aheadBehinds: []aheadBehind{{ahead: 3, behind: 4, valid: true}},
				},
			},
		},
		{
			testName: "malformed ahead-behind field becomes invalid but line is kept",
			input:    "refs/heads/feat\x00not_a_number\n",
			numBases: 1,
			expected: []branchAheadBehind{
				{
					refName:      "refs/heads/feat",
					aheadBehinds: []aheadBehind{{valid: false}},
				},
			},
		},
		{
			testName: "empty input",
			input:    "",
			numBases: 1,
			expected: nil,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			result := parseAheadBehindForEachRefOutput(s.input, s.numBases)
			assert.Equal(t, s.expected, result)
		})
	}
}

func TestSelectBaseForBranch(t *testing.T) {
	type scenario struct {
		testName           string
		aheadBehinds       []aheadBehind
		mainRefs           []string
		expectedWinner     string
		expectedBehind     int
		expectedCandidates []string
	}

	scenarios := []scenario{
		{
			testName:           "single base, valid value",
			aheadBehinds:       []aheadBehind{{ahead: 3, behind: 7, valid: true}},
			mainRefs:           []string{"refs/heads/master"},
			expectedWinner:     "refs/heads/master",
			expectedBehind:     7,
			expectedCandidates: []string{"refs/heads/master"},
		},
		{
			testName: "multi-base, clear winner by ahead",
			aheadBehinds: []aheadBehind{
				{ahead: 50, behind: 10, valid: true}, // master
				{ahead: 5, behind: 2, valid: true},   // develop  ← smallest ahead
			},
			mainRefs:           []string{"refs/heads/master", "refs/heads/develop"},
			expectedWinner:     "refs/heads/develop",
			expectedBehind:     2,
			expectedCandidates: []string{"refs/heads/develop"},
		},
		{
			testName: "develop forked from master case (ancestor-of-each-other)",
			// feat-x has 5 commits since fork from develop.
			// develop is 50 commits ahead of master.
			// ahead vs master = 5 + 50 = 55; behind vs master = 0
			// ahead vs develop = 5;          behind vs develop = 5
			aheadBehinds: []aheadBehind{
				{ahead: 55, behind: 0, valid: true}, // master
				{ahead: 5, behind: 5, valid: true},  // develop  ← smallest ahead
			},
			mainRefs:           []string{"refs/heads/master", "refs/heads/develop"},
			expectedWinner:     "refs/heads/develop",
			expectedBehind:     5,
			expectedCandidates: []string{"refs/heads/develop"},
		},
		{
			testName: "tie on ahead - first base wins (config order)",
			aheadBehinds: []aheadBehind{
				{ahead: 5, behind: 10, valid: true}, // first
				{ahead: 5, behind: 99, valid: true}, // second, same ahead
			},
			mainRefs:       []string{"refs/heads/main", "refs/heads/develop"},
			expectedWinner: "refs/heads/main",
			expectedBehind: 10,
			expectedCandidates: []string{
				"refs/heads/main",
				"refs/heads/develop",
			},
		},
		{
			testName: "first base invalid, second valid",
			aheadBehinds: []aheadBehind{
				{valid: false},
				{ahead: 3, behind: 8, valid: true},
			},
			mainRefs:           []string{"refs/heads/master", "refs/heads/develop"},
			expectedWinner:     "refs/heads/develop",
			expectedBehind:     8,
			expectedCandidates: []string{"refs/heads/develop"},
		},
		{
			testName: "all invalid - returns empty",
			aheadBehinds: []aheadBehind{
				{valid: false},
				{valid: false},
			},
			mainRefs:           []string{"refs/heads/master", "refs/heads/develop"},
			expectedWinner:     "",
			expectedBehind:     0,
			expectedCandidates: nil,
		},
		{
			testName:           "empty - returns empty",
			aheadBehinds:       nil,
			mainRefs:           nil,
			expectedWinner:     "",
			expectedBehind:     0,
			expectedCandidates: nil,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			winner, behind, candidates := selectBaseForBranch(s.aheadBehinds, s.mainRefs)
			assert.Equal(t, s.expectedWinner, winner)
			assert.Equal(t, s.expectedBehind, behind)
			assert.Equal(t, s.expectedCandidates, candidates)
		})
	}
}

func TestBuildAheadBehindForEachRefArgs(t *testing.T) {
	type scenario struct {
		testName       string
		mainBranchRefs []string
		expected       []string
	}

	scenarios := []scenario{
		{
			testName:       "single base",
			mainBranchRefs: []string{"refs/heads/master"},
			expected: []string{
				"git",
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:refs/heads/master)",
				"refs/heads",
			},
		},
		{
			testName:       "two bases",
			mainBranchRefs: []string{"refs/heads/master", "refs/remotes/origin/develop"},
			expected: []string{
				"git",
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:refs/heads/master)%00%(ahead-behind:refs/remotes/origin/develop)",
				"refs/heads",
			},
		},
		{
			testName:       "four bases",
			mainBranchRefs: []string{"refs/heads/a", "refs/heads/b", "refs/heads/c", "refs/heads/d"},
			expected: []string{
				"git",
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:refs/heads/a)%00%(ahead-behind:refs/heads/b)%00%(ahead-behind:refs/heads/c)%00%(ahead-behind:refs/heads/d)",
				"refs/heads",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			result := buildAheadBehindForEachRefArgs(s.mainBranchRefs)
			assert.Equal(t, s.expected, result)
		})
	}
}

func TestGetBehindBaseBranchValuesForAllBranches_FastPath(t *testing.T) {
	mainBranchRefs := []string{"refs/heads/master", "refs/remotes/origin/develop"}

	// Two branches: feat-x has clear divergence from develop; main matches master exactly.
	branches := []*models.Branch{
		{Name: "feat-x"},
		{Name: "main"},
	}

	expectedFormat := "%(refname)%00%(ahead-behind:refs/heads/master)%00%(ahead-behind:refs/remotes/origin/develop)"
	output := "refs/heads/feat-x\x0055 0\x005 5\n" + // picks develop (ahead=5 < 55), behind=5
		"refs/heads/main\x000 0\x000 0\n" // picks master (first, tie), behind=0

	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"for-each-ref", "--format=" + expectedFormat, "refs/heads"}, output, nil)

	gitCommon := buildGitCommon(commonDeps{
		runner:     runner,
		gitVersion: &GitVersion{2, 41, 0, ""},
	})

	loader := &BranchLoader{
		Common:    gitCommon.Common,
		GitCommon: gitCommon,
		cmd:       gitCommon.cmd,
	}

	mainBranches := &MainBranches{
		c:                    gitCommon.Common,
		cmd:                  gitCommon.cmd,
		existingMainBranches: mainBranchRefs,
		previousMainBranches: gitCommon.Common.UserConfig().Git.MainBranches,
	}

	rendered := false
	err := loader.GetBehindBaseBranchValuesForAllBranches(branches, mainBranches, func() { rendered = true })
	assert.NoError(t, err)
	assert.True(t, rendered, "renderFunc should have been called")

	assert.Equal(t, int32(5), branches[0].BehindBaseBranch.Load(), "feat-x should be behind develop by 5")
	assert.Equal(t, int32(0), branches[1].BehindBaseBranch.Load(), "main should be behind master by 0")

	runner.CheckForMissingCalls()
}

// edge case where a failure would leave artifacts from prior load
func TestGetBehindBaseBranchValuesForAllBranches_FastPath_ClearsStaleValueWhenBranchMissingFromOutput(t *testing.T) {
	mainBranchRefs := []string{"refs/heads/master"}

	feat := &models.Branch{Name: "feat-x"}
	feat.BehindBaseBranch.Store(99) // stale value from a prior load
	ghost := &models.Branch{Name: "ghost"}
	ghost.BehindBaseBranch.Store(42) // stale value from a prior load

	expectedFormat := "%(refname)%00%(ahead-behind:refs/heads/master)"
	output := "refs/heads/feat-x\x003 5\n" // ghost is intentionally absent

	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"for-each-ref", "--format=" + expectedFormat, "refs/heads"}, output, nil)

	gitCommon := buildGitCommon(commonDeps{
		runner:     runner,
		gitVersion: &GitVersion{2, 41, 0, ""},
	})

	loader := &BranchLoader{
		Common:    gitCommon.Common,
		GitCommon: gitCommon,
		cmd:       gitCommon.cmd,
	}

	mainBranches := &MainBranches{
		c:                    gitCommon.Common,
		cmd:                  gitCommon.cmd,
		existingMainBranches: mainBranchRefs,
		previousMainBranches: gitCommon.Common.UserConfig().Git.MainBranches,
	}

	err := loader.GetBehindBaseBranchValuesForAllBranches(
		[]*models.Branch{feat, ghost}, mainBranches, func() {})
	assert.NoError(t, err)

	assert.Equal(t, int32(5), feat.BehindBaseBranch.Load(), "feat-x should be updated to fresh value")
	assert.Equal(t, int32(0), ghost.BehindBaseBranch.Load(), "ghost should be reset to 0 since it has no fresh data")

	runner.CheckForMissingCalls()
}

func TestGetBehindBaseBranchValuesForAllBranches_LegacyPath(t *testing.T) {
	mainBranchRefs := []string{"refs/heads/master"}

	branches := []*models.Branch{
		{Name: "feat-x"},
	}

	// In legacy path: per-branch GetBaseBranchCandidates (merge-base +
	// for-each-ref --contains) then rev-list --left-right --count.
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"merge-base", "refs/heads/feat-x", "refs/heads/master"}, "abc123\n", nil).
		ExpectGitArgs([]string{"for-each-ref", "--contains", "abc123", "--format=%(refname)", "refs/heads/master"}, "refs/heads/master\n", nil).
		ExpectGitArgs([]string{"rev-list", "--left-right", "--count", "refs/heads/feat-x...refs/heads/master"}, "5\t7\n", nil)

	gitCommon := buildGitCommon(commonDeps{
		runner:     runner,
		gitVersion: &GitVersion{2, 34, 0, ""}, // pre-2.41, forces legacy
	})

	loader := &BranchLoader{
		Common:    gitCommon.Common,
		GitCommon: gitCommon,
		cmd:       gitCommon.cmd,
	}

	mainBranches := &MainBranches{
		c:                    gitCommon.Common,
		cmd:                  gitCommon.cmd,
		existingMainBranches: mainBranchRefs,
		previousMainBranches: gitCommon.Common.UserConfig().Git.MainBranches,
	}

	rendered := false
	err := loader.GetBehindBaseBranchValuesForAllBranches(branches, mainBranches, func() { rendered = true })
	assert.NoError(t, err)
	assert.True(t, rendered)
	assert.Equal(t, int32(7), branches[0].BehindBaseBranch.Load())

	runner.CheckForMissingCalls()
}

// When the branch's merge-base is contained in more than one configured main
// branch and the ahead counts are equal, the candidate list must preserve
// the user's configured order rather than the alphabetical order of
// for-each-ref's output.
func TestGetBaseBranchCandidates_AmbiguousReturnsAllInConfigOrder(t *testing.T) {
	mainBranchRefs := []string{"refs/heads/main", "refs/heads/develop"}
	branch := &models.Branch{Name: "feat-x"}

	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			[]string{"merge-base", "refs/heads/feat-x", "refs/heads/main", "refs/heads/develop"},
			"abc123\n", nil).
		ExpectGitArgs(
			[]string{
				"for-each-ref", "--contains", "abc123", "--format=%(refname)",
				"refs/heads/main", "refs/heads/develop",
			},
			"refs/heads/develop\nrefs/heads/main\n", nil).
		ExpectGitArgs(
			[]string{"rev-list", "--left-right", "--count", "refs/heads/feat-x...refs/heads/main"},
			"5\t10\n", nil).
		ExpectGitArgs(
			[]string{"rev-list", "--left-right", "--count", "refs/heads/feat-x...refs/heads/develop"},
			"5\t8\n", nil)

	gitCommon := buildGitCommon(commonDeps{runner: runner})

	loader := &BranchLoader{
		Common:    gitCommon.Common,
		GitCommon: gitCommon,
		cmd:       gitCommon.cmd,
	}

	mainBranches := &MainBranches{
		c:                    gitCommon.Common,
		cmd:                  gitCommon.cmd,
		existingMainBranches: mainBranchRefs,
		previousMainBranches: gitCommon.Common.UserConfig().Git.MainBranches,
	}

	candidates, err := loader.GetBaseBranchCandidates(branch, mainBranches)
	assert.NoError(t, err)
	assert.Equal(t, []string{"refs/heads/main", "refs/heads/develop"}, candidates)

	runner.CheckForMissingCalls()
}

// When a configured main branch has a strictly smaller ahead count than any
// other (e.g. the branch was forked off `main` after main's last merge into
// `develop`, so `develop` doesn't yet contain the fork point's recent main
// history), that base wins outright regardless of config order, so only
// that one ref is returned.
func TestGetBaseBranchCandidates_UnambiguousReturnsSmallestAheadOnly(t *testing.T) {
	mainBranchRefs := []string{"refs/heads/develop", "refs/heads/main"}
	branch := &models.Branch{Name: "feat-x"}

	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(
			[]string{"merge-base", "refs/heads/feat-x", "refs/heads/develop", "refs/heads/main"},
			"abc123\n", nil).
		ExpectGitArgs(
			[]string{
				"for-each-ref", "--contains", "abc123", "--format=%(refname)",
				"refs/heads/develop", "refs/heads/main",
			},
			"refs/heads/develop\nrefs/heads/main\n", nil).
		ExpectGitArgs(
			[]string{"rev-list", "--left-right", "--count", "refs/heads/feat-x...refs/heads/develop"},
			"8\t3\n", nil).
		ExpectGitArgs(
			[]string{"rev-list", "--left-right", "--count", "refs/heads/feat-x...refs/heads/main"},
			"5\t10\n", nil)

	gitCommon := buildGitCommon(commonDeps{runner: runner})

	loader := &BranchLoader{
		Common:    gitCommon.Common,
		GitCommon: gitCommon,
		cmd:       gitCommon.cmd,
	}

	mainBranches := &MainBranches{
		c:                    gitCommon.Common,
		cmd:                  gitCommon.cmd,
		existingMainBranches: mainBranchRefs,
		previousMainBranches: gitCommon.Common.UserConfig().Git.MainBranches,
	}

	candidates, err := loader.GetBaseBranchCandidates(branch, mainBranches)
	assert.NoError(t, err)
	assert.Equal(t, []string{"refs/heads/main"}, candidates)

	runner.CheckForMissingCalls()
}
