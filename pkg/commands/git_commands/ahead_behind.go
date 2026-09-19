package git_commands

import (
	"strconv"
	"strings"

	"github.com/samber/lo"
)

// Holds parsed values from a single %(ahead-behind:<base>) field.
type aheadBehind struct {
	ahead, behind int
	valid         bool
}

type branchAheadBehind struct {
	refName      string
	aheadBehinds []aheadBehind
}

// Parses output produced by:
//
//	git for-each-ref --format='%(refname)\x00%(ahead-behind:<base1>)\x00...' refs/heads
//
// Lines whose NUL-split column count doesn't match (1 + numBases) are dropped.
// Blank lines are ignored.
// Individual malformed ahead-behind fields produce {valid: false} entries, so
// that the entries of a line stay aligned with the bases.
func parseAheadBehindForEachRefOutput(
	output string,
	numBases int, // number of %(ahead-behind:...) tokens
) []branchAheadBehind {
	if output == "" {
		return nil
	}
	lines := strings.Split(output, "\n")
	result := make([]branchAheadBehind, 0, len(lines))
	for _, line := range lines {
		cols := strings.Split(line, "\x00")
		if len(cols) != numBases+1 {
			continue
		}
		refName := cols[0]
		aheadBehinds := lo.Map(cols[1:], func(col string, _ int) aheadBehind {
			return parseAheadBehindField(col)
		})
		entry := branchAheadBehind{
			refName:      refName,
			aheadBehinds: aheadBehinds,
		}
		result = append(result, entry)
	}
	return result
}

func parseAheadBehindField(s string) aheadBehind {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return aheadBehind{}
	}
	ahead, err1 := strconv.Atoi(parts[0])
	behind, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return aheadBehind{}
	}
	return aheadBehind{ahead: ahead, behind: behind, valid: true}
}

// Picks the "closest" base by smallest ahead value (commits the branch
// has that the base doesn't = roughly "since fork point") and returns
// its behind value.
// Ties are broken by index order
func selectBehindForBranch(aheadBehinds []aheadBehind) int {
	validOnes := lo.Filter(aheadBehinds, func(ab aheadBehind, _ int) bool {
		return ab.valid
	})
	return lo.MinBy(validOnes, func(a, b aheadBehind) bool {
		return a.ahead < b.ahead
	}).behind
}

// Builds a for-each-ref command that reports, for each ref matched by one of
// refPatterns, how far it is ahead and behind each of the bases. A base is a
// ref name or a commit hash. The output format is:
//
//	<refname>\x00<ahead> <behind>\x00<ahead> <behind>...\n
//
// with one ahead-behind field per base, in the same order as bases.
//
// Requires git >= 2.41 (when %(ahead-behind:...) was added).
func buildAheadBehindForEachRefArgs(bases []string, refPatterns []string) []string {
	formatParts := make([]string, 0, 1+len(bases))
	formatParts = append(formatParts, "%(refname)")
	for _, base := range bases {
		formatParts = append(formatParts, "%(ahead-behind:"+base+")")
	}
	format := strings.Join(formatParts, "%00")

	return NewGitCmd("for-each-ref").
		Arg("--format=" + format).
		Arg(refPatterns...).
		ToArgv()
}
