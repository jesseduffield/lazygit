package git_commands

import (
	"slices"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
)

// The tip commit of a ref, as far as sorting refs by ancestry needs it. The
// date is the raw %(committerdate:unix) field, and is only ever compared for
// equality.
type refTip struct {
	hash          string
	committerDate string
}

// Determining ancestry costs one %(ahead-behind:<tip>) field per tip, and git
// evaluates each of those for every tip we ask about, so the work grows with
// the square of the number of tips. Branches whose tips came out of one rebase
// are nowhere near this many, so leave the refs in the order git returned them
// once the tips add up to more than this.
const maxTipsForAncestrySorting = 100

// Sorts each group of refs that share a committer date so that a ref comes
// before the refs it is descended from. Refs that are not descended from one
// another keep the order they came in.
//
// git sorts refs with equal committer dates by name, and a stack of branches
// that gets rebased in one go ends up with the same committer date on all of
// its tips. Without this, such a stack appears in alphabetical order.
//
// refs must be sorted by committer date already, and tips must have an entry
// for each of them.
func sortRefsWithEqualDatesByAncestry[T any](
	cmd oscommands.ICmdObjBuilder,
	version *GitVersion,
	refs []T,
	fullRefName func(T) string,
	tips map[string]refTip,
) error {
	// %(ahead-behind:...) was added in git 2.41
	if !version.IsAtLeast(2, 41, 0) {
		return nil
	}

	groups := lo.Filter(groupRefsWithEqualDates(refs, fullRefName, tips),
		func(group []T, _ int) bool {
			return pointAtMoreThanOneCommit(group, fullRefName, tips)
		})
	if len(groups) == 0 {
		return nil
	}

	// Ancestry is a property of the tip commits, so ask git about each of them
	// once, however many refs point at it. A repository with several remotes
	// has the same branches under each of them, and they all share a date.
	tipHashes := []string{}
	refNames := []string{}
	seenTips := set.New[string]()
	for _, ref := range lo.Flatten(groups) {
		refName := fullRefName(ref)
		if hash := tips[refName].hash; !seenTips.Includes(hash) {
			seenTips.Add(hash)
			tipHashes = append(tipHashes, hash)
			refNames = append(refNames, refName)
		}
	}
	if len(tipHashes) > maxTipsForAncestrySorting {
		return nil
	}

	containedTips, err := loadContainedTips(cmd, tipHashes, refNames)
	if err != nil {
		return err
	}

	for _, group := range groups {
		sortGroupByAncestry(group, fullRefName, tips, containedTips)
	}

	return nil
}

// Refs that all point at the same commit have no order to be put in, so
// there is nothing to ask git about them
func pointAtMoreThanOneCommit[T any](
	refs []T,
	fullRefName func(T) string,
	tips map[string]refTip,
) bool {
	firstHash := tips[fullRefName(refs[0])].hash
	return lo.SomeBy(refs, func(ref T) bool {
		return tips[fullRefName(ref)].hash != firstHash
	})
}

// Returns the runs of consecutive refs that share a committer date, for runs of
// more than one ref. The returned slices share their backing array with refs,
// so sorting a run sorts that part of refs.
func groupRefsWithEqualDates[T any](
	refs []T,
	fullRefName func(T) string,
	tips map[string]refTip,
) [][]T {
	committerDate := func(ref T) string {
		return tips[fullRefName(ref)].committerDate
	}

	groups := [][]T{}
	start := 0
	for i := 1; i <= len(refs); i++ {
		if i < len(refs) && committerDate(refs[i]) == committerDate(refs[start]) {
			continue
		}
		if i-start > 1 {
			groups = append(groups, refs[start:i])
		}
		start = i
	}

	return groups
}

// For each of the given tips, which of the tips its history contains. Keyed by
// tip hash, with an entry for every tip passed in. refNames names, for each
// tip, a ref that points at it; git reports the values by ref name.
func loadContainedTips(
	cmd oscommands.ICmdObjBuilder,
	tipHashes []string,
	refNames []string,
) (map[string]*set.Set[string], error) {
	output, err := cmd.New(
		buildAheadBehindForEachRefArgs(tipHashes, refNames),
	).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	containedTips := make(map[string]*set.Set[string], len(tipHashes))
	for _, hash := range tipHashes {
		containedTips[hash] = set.New[string]()
	}
	tipByRefName := make(map[string]string, len(refNames))
	for i, refName := range refNames {
		tipByRefName[refName] = tipHashes[i]
	}

	for _, entry := range parseAheadBehindForEachRefOutput(output, len(tipHashes)) {
		hash, ok := tipByRefName[entry.refName]
		if !ok {
			continue
		}
		contained := containedTips[hash]
		for i, ab := range entry.aheadBehinds {
			// The tip's history contains the other tip if it has commits the
			// other one doesn't have, and the other one has none it doesn't
			// have.
			if ab.valid && ab.ahead > 0 && ab.behind == 0 {
				contained.Add(tipHashes[i])
			}
		}
	}

	return containedTips, nil
}

func sortGroupByAncestry[T any](
	group []T,
	fullRefName func(T) string,
	tips map[string]refTip,
	containedTips map[string]*set.Set[string],
) {
	hashes := lo.Map(group, func(ref T, _ int) string {
		return tips[fullRefName(ref)].hash
	})
	contained := lo.Map(hashes, func(hash string, _ int) *set.Set[string] {
		if contained, ok := containedTips[hash]; ok {
			return contained
		}
		return set.New[string]()
	})
	isDescendedFrom := func(descendant int, ancestor int) bool {
		return contained[descendant].Includes(hashes[ancestor])
	}

	// Indices into group, holding the refs we haven't placed yet
	remaining := lo.Range(len(group))
	// A ref can only go in once every ref that is descended from it is in
	canBePlaced := func(i int) bool {
		return !lo.SomeBy(remaining, func(j int) bool {
			return isDescendedFrom(j, i)
		})
	}

	sorted := make([]T, 0, len(group))
	lastPlaced := -1
	for len(remaining) > 0 {
		next := -1
		if lastPlaced != -1 {
			// Walk from the ref we placed last down to the refs it is based on,
			// so that the branches of a stack come out as one run even when an
			// unrelated branch sorts into the middle of them by name
			next = slices.IndexFunc(remaining, func(i int) bool {
				return isDescendedFrom(lastPlaced, i) && canBePlaced(i)
			})
		}
		if next == -1 {
			next = slices.IndexFunc(remaining, canBePlaced)
		}
		if next == -1 {
			// Commits can't descend from each other in a circle; only take the
			// first ref so that we don't spin if the values ever say otherwise.
			next = 0
		}
		lastPlaced = remaining[next]
		sorted = append(sorted, group[lastPlaced])
		remaining = slices.Delete(remaining, next, next+1)
	}

	copy(group, sorted)
}
