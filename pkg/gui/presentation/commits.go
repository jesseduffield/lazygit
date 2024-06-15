package presentation

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/authors"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/graph"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/kyokomi/emoji/v2"
	"github.com/samber/lo"
	"github.com/sasha-s/go-deadlock"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type pipeSetCacheKey struct {
	commitHash  string
	commitCount int
	divergence  models.Divergence
}

var (
	pipeSetCache = make(map[pipeSetCacheKey][][]*graph.Pipe)
	mutex        deadlock.Mutex
)

type bisectBounds struct {
	newIndex int
	oldIndex int
}

func GetCommitListDisplayStrings(
	common *common.Common,
	commits []*models.Commit,
	branches []*models.Branch,
	currentBranchName string,
	hasRebaseUpdateRefsConfig bool,
	fullDescription bool,
	cherryPickedCommitHashSet *set.Set[string],
	diffName string,
	markedBaseCommit string,
	timeFormat string,
	shortTimeFormat string,
	now time.Time,
	parseEmoji bool,
	selectedCommitHash string,
	startIdx int,
	endIdx int,
	showGraph bool,
	bisectInfo *git_commands.BisectInfo,
	showYouAreHereLabel bool,
) [][]string {
	mutex.Lock()
	defer mutex.Unlock()

	if len(commits) == 0 {
		return nil
	}

	if startIdx > len(commits) {
		return nil
	}

	// this is where my non-TODO commits begin
	rebaseOffset := min(indexOfFirstNonTODOCommit(commits), endIdx)

	filteredCommits := commits[startIdx:endIdx]

	bisectBounds := getbisectBounds(commits, bisectInfo)

	// function expects to be passed the index of the commit in terms of the `commits` slice
	var getGraphLine func(int) string
	if showGraph {
		if len(commits) > 0 && commits[0].Divergence != models.DivergenceNone {
			// Showing a divergence log; we know we don't have any rebasing
			// commits in this case. But we need to render separate graphs for
			// the Local and Remote sections.
			allGraphLines := []string{}

			_, localSectionStart, found := lo.FindIndexOf(
				commits, func(c *models.Commit) bool { return c.Divergence == models.DivergenceLeft })
			if !found {
				localSectionStart = len(commits)
			}

			if localSectionStart > 0 {
				// we have some remote commits
				pipeSets := loadPipesets(commits[:localSectionStart])
				if startIdx < localSectionStart {
					// some of the remote commits are visible
					start := startIdx
					end := min(endIdx, localSectionStart)
					graphPipeSets := pipeSets[start:end]
					graphCommits := commits[start:end]
					graphLines := graph.RenderAux(
						graphPipeSets,
						graphCommits,
						selectedCommitHash,
					)
					allGraphLines = append(allGraphLines, graphLines...)
				}
			}
			if localSectionStart < len(commits) {
				// we have some local commits
				pipeSets := loadPipesets(commits[localSectionStart:])
				if localSectionStart < endIdx {
					// some of the local commits are visible
					graphOffset := max(startIdx, localSectionStart)
					pipeSetOffset := max(startIdx-localSectionStart, 0)
					graphPipeSets := pipeSets[pipeSetOffset : endIdx-localSectionStart]
					graphCommits := commits[graphOffset:endIdx]
					graphLines := graph.RenderAux(
						graphPipeSets,
						graphCommits,
						selectedCommitHash,
					)
					allGraphLines = append(allGraphLines, graphLines...)
				}
			}

			getGraphLine = func(idx int) string {
				return allGraphLines[idx-startIdx]
			}
		} else {
			// this is where the graph begins (may be beyond the TODO commits depending on startIdx,
			// but we'll never include TODO commits as part of the graph because it'll be messy)
			graphOffset := max(startIdx, rebaseOffset)

			pipeSets := loadPipesets(commits[rebaseOffset:])
			pipeSetOffset := max(startIdx-rebaseOffset, 0)
			graphPipeSets := pipeSets[pipeSetOffset:max(endIdx-rebaseOffset, 0)]
			graphCommits := commits[graphOffset:endIdx]
			graphLines := graph.RenderAux(
				graphPipeSets,
				graphCommits,
				selectedCommitHash,
			)
			getGraphLine = func(idx int) string {
				if idx >= graphOffset {
					return graphLines[idx-graphOffset]
				} else {
					return ""
				}
			}
		}
	} else {
		getGraphLine = func(int) string { return "" }
	}

	// Determine the hashes of the local branches for which we want to show a
	// branch marker in the commits list. We only want to do this for branches
	// that are not the current branch, and not any of the main branches. The
	// goal is to visualize stacks of local branches, so anything that doesn't
	// contribute to a branch stack shouldn't show a marker.
	//
	// If there are other branches pointing to the current head commit, we only
	// want to show the marker if the rebase.updateRefs config is on.
	branchHeadsToVisualize := set.NewFromSlice(lo.FilterMap(branches,
		func(b *models.Branch, index int) (string, bool) {
			return b.CommitHash,
				// Don't consider branches that don't have a commit hash. As far
				// as I can see, this happens for a detached head, so filter
				// these out
				b.CommitHash != "" &&
					// Don't show a marker for the current branch
					b.Name != currentBranchName &&
					// Don't show a marker for main branches
					!lo.Contains(common.UserConfig.Git.MainBranches, b.Name) &&
					// Don't show a marker for the head commit unless the
					// rebase.updateRefs config is on
					(hasRebaseUpdateRefsConfig || b.CommitHash != commits[0].Hash)
		}))

	lines := make([][]string, 0, len(filteredCommits))
	var bisectStatus BisectStatus
	willBeRebased := markedBaseCommit == ""
	for i, commit := range filteredCommits {
		unfilteredIdx := i + startIdx
		bisectStatus = getBisectStatus(unfilteredIdx, commit.Hash, bisectInfo, bisectBounds)
		isYouAreHereCommit := false
		if showYouAreHereLabel && (commit.Action == models.ActionConflict || unfilteredIdx == rebaseOffset) {
			isYouAreHereCommit = true
			showYouAreHereLabel = false
		}
		isMarkedBaseCommit := commit.Hash != "" && commit.Hash == markedBaseCommit
		if isMarkedBaseCommit {
			willBeRebased = true
		}
		lines = append(lines, displayCommit(
			common,
			commit,
			branchHeadsToVisualize,
			hasRebaseUpdateRefsConfig,
			cherryPickedCommitHashSet,
			isMarkedBaseCommit,
			willBeRebased,
			diffName,
			timeFormat,
			shortTimeFormat,
			now,
			parseEmoji,
			getGraphLine(unfilteredIdx),
			fullDescription,
			bisectStatus,
			bisectInfo,
			isYouAreHereCommit,
		))
	}
	return lines
}

func getbisectBounds(commits []*models.Commit, bisectInfo *git_commands.BisectInfo) *bisectBounds {
	if !bisectInfo.Bisecting() {
		return nil
	}

	bisectBounds := &bisectBounds{}

	for i, commit := range commits {
		if commit.Hash == bisectInfo.GetNewHash() {
			bisectBounds.newIndex = i
		}

		status, ok := bisectInfo.Status(commit.Hash)
		if ok && status == git_commands.BisectStatusOld {
			bisectBounds.oldIndex = i
			return bisectBounds
		}
	}

	// shouldn't land here
	return nil
}

// precondition: slice is not empty
func indexOfFirstNonTODOCommit(commits []*models.Commit) int {
	for i, commit := range commits {
		if !commit.IsTODO() {
			return i
		}
	}

	// shouldn't land here
	return 0
}

func loadPipesets(commits []*models.Commit) [][]*graph.Pipe {
	// given that our cache key is a commit hash and a commit count, it's very important that we don't actually try to render pipes
	// when dealing with things like filtered commits.
	cacheKey := pipeSetCacheKey{
		commitHash:  commits[0].Hash,
		commitCount: len(commits),
		divergence:  commits[0].Divergence,
	}

	pipeSets, ok := pipeSetCache[cacheKey]
	if !ok {
		// pipe sets are unique to a commit head. and a commit count. Sometimes we haven't loaded everything for that.
		// so let's just cache it based on that.
		getStyle := func(commit *models.Commit) style.TextStyle {
			return authors.AuthorStyle(commit.AuthorName)
		}
		pipeSets = graph.GetPipeSets(commits, getStyle)
		pipeSetCache[cacheKey] = pipeSets
	}

	return pipeSets
}

// similar to the git_commands.BisectStatus but more gui-focused
type BisectStatus int

const (
	BisectStatusNone BisectStatus = iota
	BisectStatusOld
	BisectStatusNew
	BisectStatusSkipped
	// adding candidate here which isn't present in the commands package because
	// we need to actually go through the commits to get this info
	BisectStatusCandidate
	// also adding this
	BisectStatusCurrent
)

func getBisectStatus(index int, commitHash string, bisectInfo *git_commands.BisectInfo, bisectBounds *bisectBounds) BisectStatus {
	if !bisectInfo.Started() {
		return BisectStatusNone
	}

	if bisectInfo.GetCurrentHash() == commitHash {
		return BisectStatusCurrent
	}

	status, ok := bisectInfo.Status(commitHash)
	if ok {
		switch status {
		case git_commands.BisectStatusNew:
			return BisectStatusNew
		case git_commands.BisectStatusOld:
			return BisectStatusOld
		case git_commands.BisectStatusSkipped:
			return BisectStatusSkipped
		}
	} else {
		if bisectBounds != nil && index >= bisectBounds.newIndex && index <= bisectBounds.oldIndex {
			return BisectStatusCandidate
		} else {
			return BisectStatusNone
		}
	}

	// should never land here
	return BisectStatusNone
}

func getBisectStatusText(bisectStatus BisectStatus, bisectInfo *git_commands.BisectInfo) string {
	if bisectStatus == BisectStatusNone {
		return ""
	}

	style := getBisectStatusColor(bisectStatus)

	switch bisectStatus {
	case BisectStatusNew:
		return style.Sprintf("<-- " + bisectInfo.NewTerm())
	case BisectStatusOld:
		return style.Sprintf("<-- " + bisectInfo.OldTerm())
	case BisectStatusCurrent:
		// TODO: i18n
		return style.Sprintf("<-- current")
	case BisectStatusSkipped:
		return style.Sprintf("<-- skipped")
	case BisectStatusCandidate:
		return style.Sprintf("?")
	case BisectStatusNone:
		return ""
	}

	return ""
}

func displayCommit(
	common *common.Common,
	commit *models.Commit,
	branchHeadsToVisualize *set.Set[string],
	hasRebaseUpdateRefsConfig bool,
	cherryPickedCommitHashSet *set.Set[string],
	isMarkedBaseCommit bool,
	willBeRebased bool,
	diffName string,
	timeFormat string,
	shortTimeFormat string,
	now time.Time,
	parseEmoji bool,
	graphLine string,
	fullDescription bool,
	bisectStatus BisectStatus,
	bisectInfo *git_commands.BisectInfo,
	isYouAreHereCommit bool,
) []string {
	bisectString := getBisectStatusText(bisectStatus, bisectInfo)

	hashString := ""
	hashColor := getHashColor(commit, diffName, cherryPickedCommitHashSet, bisectStatus, bisectInfo)
	hashLength := common.UserConfig.Gui.CommitHashLength
	if hashLength >= len(commit.Hash) {
		hashString = hashColor.Sprint(commit.Hash)
	} else if hashLength > 0 {
		hashString = hashColor.Sprint(commit.Hash[:hashLength])
	} else if !icons.IsIconEnabled() { // hashLength <= 0
		hashString = hashColor.Sprint("*")
	}

	divergenceString := ""
	if commit.Divergence != models.DivergenceNone {
		divergenceString = hashColor.Sprint(lo.Ternary(commit.Divergence == models.DivergenceLeft, "↑", "↓"))
	} else if icons.IsIconEnabled() {
		divergenceString = hashColor.Sprint(icons.IconForCommit(commit))
	}

	descriptionString := ""
	if fullDescription {
		descriptionString = style.FgBlue.Sprint(
			utils.UnixToDateSmart(now, commit.UnixTimestamp, timeFormat, shortTimeFormat),
		)
	}

	actionString := ""
	if commit.Action != models.ActionNone {
		todoString := lo.Ternary(commit.Action == models.ActionConflict, "conflict", commit.Action.String())
		actionString = actionColorMap(commit.Action).Sprint(todoString) + " "
	}

	tagString := ""
	if fullDescription {
		if commit.ExtraInfo != "" {
			tagString = style.FgMagenta.SetBold().Sprint(commit.ExtraInfo) + " "
		}
	} else {
		if len(commit.Tags) > 0 {
			tagString = theme.DiffTerminalColor.SetBold().Sprint(strings.Join(commit.Tags, " ")) + " "
		}

		if branchHeadsToVisualize.Includes(commit.Hash) &&
			// Don't show branch head on commits that are already merged to a main branch
			commit.Status != models.StatusMerged &&
			// Don't show branch head on a "pick" todo if the rebase.updateRefs config is on
			!(commit.IsTODO() && hasRebaseUpdateRefsConfig) {
			tagString = style.FgCyan.SetBold().Sprint(
				lo.Ternary(icons.IsIconEnabled(), icons.BRANCH_ICON, "*") + " " + tagString)
		}
	}

	name := commit.Name
	if commit.Action == todo.UpdateRef {
		name = strings.TrimPrefix(name, "refs/heads/")
	}
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	mark := ""
	if isYouAreHereCommit {
		color := lo.Ternary(commit.Action == models.ActionConflict, style.FgRed, style.FgYellow)
		youAreHere := color.Sprintf("<-- %s ---", common.Tr.YouAreHere)
		mark = fmt.Sprintf("%s ", youAreHere)
	} else if isMarkedBaseCommit {
		rebaseFromHere := style.FgYellow.Sprint(common.Tr.MarkedCommitMarker)
		mark = fmt.Sprintf("%s ", rebaseFromHere)
	} else if !willBeRebased {
		willBeRebased := style.FgYellow.Sprint("✓")
		mark = fmt.Sprintf("%s ", willBeRebased)
	}

	authorLength := common.UserConfig.Gui.CommitAuthorShortLength
	if fullDescription {
		authorLength = common.UserConfig.Gui.CommitAuthorLongLength
	}
	author := authors.AuthorWithLength(commit.AuthorName, authorLength)

	cols := make([]string, 0, 7)
	cols = append(
		cols,
		divergenceString,
		hashString,
		bisectString,
		descriptionString,
		actionString,
		author,
		graphLine+mark+tagString+theme.DefaultTextColor.Sprint(name),
	)

	return cols
}

func getBisectStatusColor(status BisectStatus) style.TextStyle {
	switch status {
	case BisectStatusNone:
		return style.FgBlack
	case BisectStatusNew:
		return style.FgRed
	case BisectStatusOld:
		return style.FgGreen
	case BisectStatusSkipped:
		return style.FgYellow
	case BisectStatusCurrent:
		return style.FgMagenta
	case BisectStatusCandidate:
		return style.FgBlue
	}

	// shouldn't land here
	return style.FgWhite
}

func getHashColor(
	commit *models.Commit,
	diffName string,
	cherryPickedCommitHashSet *set.Set[string],
	bisectStatus BisectStatus,
	bisectInfo *git_commands.BisectInfo,
) style.TextStyle {
	if bisectInfo.Started() {
		return getBisectStatusColor(bisectStatus)
	}

	diffed := commit.Hash != "" && commit.Hash == diffName
	hashColor := theme.DefaultTextColor
	switch commit.Status {
	case models.StatusUnpushed:
		hashColor = style.FgRed
	case models.StatusPushed:
		hashColor = style.FgYellow
	case models.StatusMerged:
		hashColor = style.FgGreen
	case models.StatusRebasing:
		hashColor = style.FgBlue
	case models.StatusReflog:
		hashColor = style.FgBlue
	default:
	}

	if diffed {
		hashColor = theme.DiffTerminalColor
	} else if cherryPickedCommitHashSet.Includes(commit.Hash) {
		hashColor = theme.CherryPickedCommitTextStyle
	} else if commit.Divergence == models.DivergenceRight && commit.Status != models.StatusMerged {
		hashColor = style.FgBlue
	}

	return hashColor
}

func actionColorMap(action todo.TodoCommand) style.TextStyle {
	switch action {
	case todo.Pick:
		return style.FgCyan
	case todo.Drop:
		return style.FgRed
	case todo.Edit:
		return style.FgGreen
	case todo.Fixup:
		return style.FgMagenta
	case models.ActionConflict:
		return style.FgRed
	default:
		return style.FgYellow
	}
}
