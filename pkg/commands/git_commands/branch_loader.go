package git_commands

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

// context:
// we want to only show 'safe' branches (ones that haven't e.g. been deleted)
// which `git branch -a` gives us, but we also want the recency data that
// git reflog gives us.
// So we get the HEAD, then append get the reflog branches that intersect with
// our safe branches, then add the remaining safe branches, ensuring uniqueness
// along the way

// if we find out we need to use one of these functions in the git.go file, we
// can just pull them out of here and put them there and then call them from in here

type BranchLoaderConfigCommands interface {
	Branches(cmd oscommands.ICmdObjBuilder) map[string]*BranchConfig
}

type BranchInfo struct {
	RefName      string
	DisplayName  string // e.g. '(HEAD detached at 123asdf)'
	DetachedHead bool
}

// BranchLoader returns a list of Branch objects for the current repo
type BranchLoader struct {
	*common.Common
	*GitCommon
	cmd                  oscommands.ICmdObjBuilder
	getCurrentBranchInfo func() (BranchInfo, error)
	config               BranchLoaderConfigCommands
}

func NewBranchLoader(
	cmn *common.Common,
	gitCommon *GitCommon,
	cmd oscommands.ICmdObjBuilder,
	getCurrentBranchInfo func() (BranchInfo, error),
	config BranchLoaderConfigCommands,
) *BranchLoader {
	return &BranchLoader{
		Common:               cmn,
		GitCommon:            gitCommon,
		cmd:                  cmd,
		getCurrentBranchInfo: getCurrentBranchInfo,
		config:               config,
	}
}

// Load the list of branches for the current repo
func (self *BranchLoader) Load(reflogCommits []*models.Commit,
	mainBranches *MainBranches,
	oldBranches []*models.Branch,
	loadBehindCounts bool,
	onWorker func(func() error),
	renderFunc func(),
) ([]*models.Branch, error) {
	branches := self.obtainBranches()

	if self.UserConfig().Git.LocalBranchSortOrder == "recency" {
		reflogBranches := self.obtainReflogBranches(reflogCommits)
		// loop through reflog branches. If there is a match, merge them, then remove it from the branches and keep it in the reflog branches
		branchesWithRecency := make([]*models.Branch, 0)
	outer:
		for _, reflogBranch := range reflogBranches {
			for j, branch := range branches {
				if branch.Head {
					continue
				}
				if strings.EqualFold(reflogBranch.Name, branch.Name) {
					branch.Recency = reflogBranch.Recency
					branchesWithRecency = append(branchesWithRecency, branch)
					branches = utils.Remove(branches, j)
					continue outer
				}
			}
		}

		// Sort branches that don't have a recency value alphabetically
		// (we're really doing this for the sake of deterministic behaviour across git versions)
		slices.SortFunc(branches, func(a *models.Branch, b *models.Branch) int {
			return strings.Compare(a.Name, b.Name)
		})

		branches = utils.Prepend(branches, branchesWithRecency...)
	}

	foundHead := false
	for i, branch := range branches {
		if branch.Head {
			foundHead = true
			branch.Recency = "  *"
			branches = utils.Move(branches, i, 0)
			break
		}
	}
	if !foundHead {
		info, err := self.getCurrentBranchInfo()
		if err != nil {
			return nil, err
		}
		branches = utils.Prepend(branches, &models.Branch{Name: info.RefName, DisplayName: info.DisplayName, Head: true, DetachedHead: info.DetachedHead, Recency: "  *"})
	}

	configBranches := self.config.Branches(self.cmd)

	for _, branch := range branches {
		match := configBranches[branch.Name]
		if match != nil {
			branch.UpstreamRemote = match.Remote
			branch.UpstreamBranch = match.Merge
		}

		// If the branch already existed, take over its BehindBaseBranch value
		// to reduce flicker
		if oldBranch, found := lo.Find(oldBranches, func(b *models.Branch) bool {
			return b.Name == branch.Name
		}); found {
			branch.BehindBaseBranch.Store(oldBranch.BehindBaseBranch.Load())
		}
	}

	if loadBehindCounts && self.UserConfig().Gui.ShowDivergenceFromBaseBranch != "none" {
		onWorker(func() error {
			return self.GetBehindBaseBranchValuesForAllBranches(branches, mainBranches, renderFunc)
		})
	}

	return branches, nil
}

func (self *BranchLoader) GetBehindBaseBranchValuesForAllBranches(
	branches []*models.Branch,
	mainBranches *MainBranches,
	renderFunc func(),
) error {
	mainBranchRefs := mainBranches.Get()
	if len(mainBranchRefs) == 0 {
		return nil
	}

	if self.version.IsAtLeast(2, 41, 0) {
		return self.getBehindBaseBranchValuesFast(branches, mainBranches, renderFunc)
	}
	return self.getBehindBaseBranchValuesLegacy(branches, mainBranches, renderFunc)
}

func (self *BranchLoader) getBehindBaseBranchValuesLegacy(
	branches []*models.Branch,
	mainBranches *MainBranches,
	renderFunc func(),
) error {
	t := time.Now()
	errg := errgroup.Group{}

	for _, branch := range branches {
		errg.Go(func() error {
			_, behinds, err := self.baseBranchCandidatesAndBehinds(branch, mainBranches)
			if err != nil {
				return err
			}
			branch.BehindBaseBranch.Store(classifyBehind(behinds))
			return nil
		})
	}

	err := errg.Wait()
	self.Log.Debugf("time to get behind base branch values for all branches (legacy): %s", time.Since(t))
	renderFunc()
	return err
}

// Holds parsed values from a single %(ahead-behind:<base>) field. `valid`
// is false when the field failed to parse (e.g. the base was unreachable
// from this ref); the entry is preserved so that the slice stays index-
// aligned with the configured main branches.
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
// Individual malformed ahead-behind fields produce {valid: false} entries.
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

// selectBaseForBranch picks the closest base(s) for a branch given
// (ahead, behind) measurements against each configured main branch.
// "Closest" = smallest ahead value (fewest branch commits not in the
// base). Ties are broken by the order of mainRefs (i.e. config order).
//
// aheadBehinds must be index-aligned with mainRefs; invalid entries are
// skipped. Returns parallel slices: the refs tied at the minimum ahead
// (in config order) and their behind values. The caller picks
// candidates[0] for a single answer, or detects ambiguity via
// `len(candidates) > 1`.
func selectBaseForBranch(
	aheadBehinds []aheadBehind, mainRefs []string,
) (candidates []string, behinds []int) {
	bestAhead := -1
	for i, ab := range aheadBehinds {
		if !ab.valid {
			continue
		}
		switch {
		case bestAhead < 0 || ab.ahead < bestAhead:
			bestAhead = ab.ahead
			candidates = []string{mainRefs[i]}
			behinds = []int{ab.behind}
		case ab.ahead == bestAhead:
			candidates = append(candidates, mainRefs[i])
			behinds = append(behinds, ab.behind)
		}
	}
	return candidates, behinds
}

// classifyBehind condenses per-candidate behind values into the single
// number stored on Branch.BehindBaseBranch for column display. When the
// candidates all agree (possibly on 0), return that value; otherwise
// return one of the BehindBaseAmbiguous* sentinels so the renderer can
// show "?" or "↓?" instead of a misleadingly precise count.
func classifyBehind(behinds []int) int32 {
	if len(behinds) == 0 {
		return 0
	}
	first := behinds[0]
	allEqual := true
	anyZero := first == 0
	for _, b := range behinds[1:] {
		if b != first {
			allEqual = false
		}
		if b == 0 {
			anyZero = true
		}
	}
	if allEqual {
		return int32(first)
	}
	if anyZero {
		return models.BehindBaseAmbiguousMaybeUpToDate
	}
	return models.BehindBaseAmbiguousDefinitelyBehind
}

// The output format is:
//
//	<refname>\x00<ahead> <behind>\x00<ahead> <behind>...\n
//
// with one ahead-behind field per base, in the same order as mainBranchRefs.
//
// Requires git >= 2.41 (when %(ahead-behind:...) was added).
func buildAheadBehindForEachRefArgs(mainBranchRefs []string) []string {
	formatParts := make([]string, 0, 1+len(mainBranchRefs))
	formatParts = append(formatParts, "%(refname)")
	for _, ref := range mainBranchRefs {
		formatParts = append(formatParts, "%(ahead-behind:"+ref+")")
	}
	format := strings.Join(formatParts, "%00")

	return NewGitCmd("for-each-ref").
		Arg("--format=" + format).
		Arg("refs/heads").
		ToArgv()
}

func (self *BranchLoader) getBehindBaseBranchValuesFast(
	branches []*models.Branch,
	mainBranches *MainBranches,
	renderFunc func(),
) error {
	t := time.Now()

	mainBranchRefs := mainBranches.Get()
	output, err := self.cmd.New(
		buildAheadBehindForEachRefArgs(mainBranchRefs),
	).DontLog().RunWithOutput()
	if err != nil {
		return err
	}

	parsed := parseAheadBehindForEachRefOutput(output, len(mainBranchRefs))
	branchByRef := lo.KeyBy(branches, (*models.Branch).FullRefName)

	for _, p := range parsed {
		if branch, ok := branchByRef[p.refName]; ok {
			_, behinds := selectBaseForBranch(p.aheadBehinds, mainBranchRefs)
			branch.BehindBaseBranch.Store(classifyBehind(behinds))
			delete(branchByRef, p.refName)
		}
	}

	// Branches not in parse are default to 0
	for _, branch := range branchByRef {
		branch.BehindBaseBranch.Store(0)
	}

	self.Log.Debugf("time to get behind base branch values for all branches (fast): %s", time.Since(t))
	renderFunc()
	return nil
}

// GetBaseBranchCandidates returns the configured main branches that are the
// closest base for the given branch — typically a single ref, but more
// when the closeness rule (smallest ahead value) leaves a tie. Candidates
// are returned in config order, so callers wanting one answer can use
// candidates[0] as the config-order tiebreak. An empty slice (with nil
// error) means no configured main branch contains the branch's merge-base.
func (self *BranchLoader) GetBaseBranchCandidates(branch *models.Branch, mainBranches *MainBranches) ([]string, error) {
	candidates, _, err := self.baseBranchCandidatesAndBehinds(branch, mainBranches)
	return candidates, err
}

// baseBranchCandidatesAndBehinds is the full computation behind
// GetBaseBranchCandidates: it also reports the behind count for each
// returned candidate, which the legacy behind-base loader needs in order
// to classify the column display when the candidates disagree. Slices
// are parallel and in config order.
func (self *BranchLoader) baseBranchCandidatesAndBehinds(branch *models.Branch, mainBranches *MainBranches) ([]string, []int, error) {
	mergeBase := mainBranches.GetMergeBase(branch.FullRefName())
	if mergeBase == "" {
		return nil, nil, nil
	}

	mainBranchRefs := mainBranches.Get()
	output, err := self.cmd.New(
		NewGitCmd("for-each-ref").
			Arg("--contains").
			Arg(mergeBase).
			Arg("--format=%(refname)").
			Arg(mainBranchRefs...).
			ToArgv(),
	).DontLog().RunWithOutput()
	if err != nil {
		return nil, nil, err
	}
	trimmedOutput := strings.TrimSpace(output)
	if trimmedOutput == "" {
		return nil, nil, nil
	}
	contained := strings.Split(trimmedOutput, "\n")

	// for-each-ref sorts its output alphabetically by refname regardless of
	// the order we passed the refs in. Restore the user's configured order so
	// it can serve as the natural tiebreaker.
	containing := lo.Filter(mainBranchRefs, func(ref string, _ int) bool {
		return lo.Contains(contained, ref)
	})
	if len(containing) == 0 {
		return nil, nil, nil
	}

	// Measure ahead/behind against each containing ref and hand off to
	// selectBaseForBranch — the same selector the fast path uses — so
	// both paths agree on the closeness rule and the config-order
	// tiebreak. We do this even when there's only one containing ref,
	// because the legacy column display still needs the behind value.
	aheadBehinds := make([]aheadBehind, len(containing))
	for i, ref := range containing {
		revListOutput, err := self.cmd.New(
			NewGitCmd("rev-list").
				Arg("--left-right").
				Arg("--count").
				Arg(fmt.Sprintf("%s...%s", branch.FullRefName(), ref)).
				ToArgv(),
		).DontLog().RunWithOutput()
		if err != nil {
			return nil, nil, err
		}
		aheadBehinds[i] = parseAheadBehindField(strings.TrimSpace(revListOutput))
	}

	candidates, behinds := selectBaseForBranch(aheadBehinds, containing)
	if len(candidates) == 0 {
		// Every rev-list output was malformed; fall back to config order
		// with no reliable behinds.
		return containing, nil, nil
	}
	return candidates, behinds, nil
}

func (self *BranchLoader) obtainBranches() []*models.Branch {
	output, err := self.getRawBranches()
	if err != nil {
		panic(err)
	}

	trimmedOutput := strings.TrimSpace(output)
	outputLines := strings.Split(trimmedOutput, "\n")

	return lo.FilterMap(outputLines, func(line string, _ int) (*models.Branch, bool) {
		if line == "" {
			return nil, false
		}

		split := strings.Split(line, "\x00")
		if len(split) != len(branchFields) {
			// Ignore line if it isn't separated into the expected number of parts
			// This is probably a warning message, for more info see:
			// https://github.com/jesseduffield/lazygit/issues/1385#issuecomment-885580439
			return nil, false
		}

		storeCommitDateAsRecency := self.UserConfig().Git.LocalBranchSortOrder != "recency"
		return obtainBranch(split, storeCommitDateAsRecency), true
	})
}

func (self *BranchLoader) getRawBranches() (string, error) {
	format := strings.Join(
		lo.Map(branchFields, func(thing string, _ int) string {
			return "%(" + thing + ")"
		}),
		"%00",
	)

	var sortOrder string
	switch strings.ToLower(self.UserConfig().Git.LocalBranchSortOrder) {
	case "recency", "date":
		sortOrder = "-committerdate"
	case "alphabetical":
		sortOrder = "refname"
	default:
		sortOrder = "refname"
	}

	cmdArgs := NewGitCmd("for-each-ref").
		Arg(fmt.Sprintf("--sort=%s", sortOrder)).
		Arg(fmt.Sprintf("--format=%s", format)).
		Arg("refs/heads").
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

var branchFields = []string{
	"HEAD",
	"refname:short",
	"upstream:short",
	"upstream:track",
	"push:track",
	"subject",
	"objectname",
	"committerdate:unix",
}

// Obtain branch information from parsed line output of getRawBranches()
func obtainBranch(split []string, storeCommitDateAsRecency bool) *models.Branch {
	headMarker := split[0]
	fullName := split[1]
	upstreamName := split[2]
	track := split[3]
	pushTrack := split[4]
	subject := split[5]
	commitHash := split[6]
	commitDate := split[7]

	name := strings.TrimPrefix(fullName, "heads/")
	aheadForPull, behindForPull, gone := parseUpstreamInfo(upstreamName, track)
	aheadForPush, behindForPush, _ := parseUpstreamInfo(upstreamName, pushTrack)

	recency := ""
	if storeCommitDateAsRecency {
		if unixTimestamp, err := strconv.ParseInt(commitDate, 10, 64); err == nil {
			recency = utils.UnixToTimeAgo(unixTimestamp)
		}
	}

	return &models.Branch{
		Name:          name,
		Recency:       recency,
		AheadForPull:  aheadForPull,
		BehindForPull: behindForPull,
		AheadForPush:  aheadForPush,
		BehindForPush: behindForPush,
		UpstreamGone:  gone,
		Head:          headMarker == "*",
		Subject:       subject,
		CommitHash:    commitHash,
	}
}

func parseUpstreamInfo(upstreamName string, track string) (string, string, bool) {
	if upstreamName == "" {
		// if we're here then it means we do not have a local version of the remote.
		// The branch might still be tracking a remote though, we just don't know
		// how many commits ahead/behind it is
		return "?", "?", false
	}

	if track == "[gone]" {
		return "?", "?", true
	}

	ahead := parseDifference(track, `ahead (\d+)`)
	behind := parseDifference(track, `behind (\d+)`)

	return ahead, behind, false
}

func parseDifference(track string, regexStr string) string {
	re := regexp.MustCompile(regexStr)
	match := re.FindStringSubmatch(track)
	if len(match) > 1 {
		return match[1]
	}
	return "0"
}

// TODO: only look at the new reflog commits, and otherwise store the recencies in
// int form against the branch to recalculate the time ago
func (self *BranchLoader) obtainReflogBranches(reflogCommits []*models.Commit) []*models.Branch {
	foundBranches := set.New[string]()
	re := regexp.MustCompile(`checkout: moving from ([\S]+) to ([\S]+)`)
	reflogBranches := make([]*models.Branch, 0, len(reflogCommits))

	for _, commit := range reflogCommits {
		match := re.FindStringSubmatch(commit.Name)
		if len(match) != 3 {
			continue
		}

		recency := utils.UnixToTimeAgo(commit.UnixTimestamp)
		for _, branchName := range match[1:] {
			if !foundBranches.Includes(branchName) {
				foundBranches.Add(branchName)
				reflogBranches = append(reflogBranches, &models.Branch{
					Recency: recency,
					Name:    branchName,
				})
			}
		}
	}
	return reflogBranches
}
