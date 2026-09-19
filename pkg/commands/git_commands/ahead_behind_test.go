package git_commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
						{},
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
					aheadBehinds: []aheadBehind{{}},
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

func TestSelectBehindForBranch(t *testing.T) {
	type scenario struct {
		testName     string
		aheadBehinds []aheadBehind
		expected     int
	}

	scenarios := []scenario{
		{
			testName:     "single base, valid value",
			aheadBehinds: []aheadBehind{{ahead: 3, behind: 7, valid: true}},
			expected:     7,
		},
		{
			testName: "multi-base, clear winner by ahead",
			aheadBehinds: []aheadBehind{
				{ahead: 50, behind: 10, valid: true}, // master
				{ahead: 5, behind: 2, valid: true},   // develop  ← smallest ahead
			},
			expected: 2,
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
			expected: 5,
		},
		{
			testName: "tie on ahead - first base wins (config order)",
			aheadBehinds: []aheadBehind{
				{ahead: 5, behind: 10, valid: true}, // first
				{ahead: 5, behind: 99, valid: true}, // second, same ahead
			},
			expected: 10,
		},
		{
			testName: "first base invalid, second valid",
			aheadBehinds: []aheadBehind{
				{},
				{ahead: 3, behind: 8, valid: true},
			},
			expected: 8,
		},
		{
			testName:     "all invalid - returns 0",
			aheadBehinds: []aheadBehind{{}, {}},
			expected:     0,
		},
		{
			testName:     "empty - returns 0",
			aheadBehinds: nil,
			expected:     0,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			result := selectBehindForBranch(s.aheadBehinds)
			assert.Equal(t, s.expected, result)
		})
	}
}

func TestBuildAheadBehindForEachRefArgs(t *testing.T) {
	type scenario struct {
		testName    string
		bases       []string
		refPatterns []string
		expected    []string
	}

	scenarios := []scenario{
		{
			testName:    "single base",
			bases:       []string{"refs/heads/master"},
			refPatterns: []string{"refs/heads"},
			expected: []string{
				"git",
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:refs/heads/master)",
				"refs/heads",
			},
		},
		{
			testName:    "two bases",
			bases:       []string{"refs/heads/master", "refs/remotes/origin/develop"},
			refPatterns: []string{"refs/heads"},
			expected: []string{
				"git",
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:refs/heads/master)%00%(ahead-behind:refs/remotes/origin/develop)",
				"refs/heads",
			},
		},
		{
			testName:    "four bases",
			bases:       []string{"refs/heads/a", "refs/heads/b", "refs/heads/c", "refs/heads/d"},
			refPatterns: []string{"refs/heads"},
			expected: []string{
				"git",
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:refs/heads/a)%00%(ahead-behind:refs/heads/b)%00%(ahead-behind:refs/heads/c)%00%(ahead-behind:refs/heads/d)",
				"refs/heads",
			},
		},
		{
			testName:    "commit hashes as bases, individual refs as patterns",
			bases:       []string{"1234567", "89abcde"},
			refPatterns: []string{"refs/heads/a", "refs/remotes/origin/b"},
			expected: []string{
				"git",
				"for-each-ref",
				"--format=%(refname)%00%(ahead-behind:1234567)%00%(ahead-behind:89abcde)",
				"refs/heads/a",
				"refs/remotes/origin/b",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			result := buildAheadBehindForEachRefArgs(s.bases, s.refPatterns)
			assert.Equal(t, s.expected, result)
		})
	}
}
