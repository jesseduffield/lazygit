package presentation

import (
	"fmt"
	"strings"

	"github.com/fsmiamoto/git-todo-parser/todo"
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
	"github.com/sasha-s/go-deadlock"
)

type pipeSetCacheKey struct {
	commitSha   string
	commitCount int
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
	fullDescription bool,
	cherryPickedCommitShaSet *set.Set[string],
	diffName string,
	timeFormat string,
	parseEmoji bool,
	selectedCommitSha string,
	startIdx int,
	length int,
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

	end := utils.Min(startIdx+length, len(commits))
	// this is where my non-TODO commits begin
	rebaseOffset := utils.Min(indexOfFirstNonTODOCommit(commits), end)

	filteredCommits := commits[startIdx:end]

	bisectBounds := getbisectBounds(commits, bisectInfo)

	// function expects to be passed the index of the commit in terms of the `commits` slice
	var getGraphLine func(int) string
	if showGraph {
		// this is where the graph begins (may be beyond the TODO commits depending on startIdx,
		// but we'll never include TODO commits as part of the graph because it'll be messy)
		graphOffset := utils.Max(startIdx, rebaseOffset)

		pipeSets := loadPipesets(commits[rebaseOffset:])
		pipeSetOffset := utils.Max(startIdx-rebaseOffset, 0)
		graphPipeSets := pipeSets[pipeSetOffset:utils.Max(end-rebaseOffset, 0)]
		graphCommits := commits[graphOffset:end]
		graphLines := graph.RenderAux(
			graphPipeSets,
			graphCommits,
			selectedCommitSha,
		)
		getGraphLine = func(idx int) string {
			if idx >= graphOffset {
				return graphLines[idx-graphOffset]
			} else {
				return ""
			}
		}
	} else {
		getGraphLine = func(idx int) string { return "" }
	}

	lines := make([][]string, 0, len(filteredCommits))
	var bisectStatus BisectStatus
	for i, commit := range filteredCommits {
		unfilteredIdx := i + startIdx
		bisectStatus = getBisectStatus(unfilteredIdx, commit.Sha, bisectInfo, bisectBounds)
		isYouAreHereCommit := showYouAreHereLabel && unfilteredIdx == rebaseOffset
		lines = append(lines, displayCommit(
			common,
			commit,
			cherryPickedCommitShaSet,
			diffName,
			timeFormat,
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
		if commit.Sha == bisectInfo.GetNewSha() {
			bisectBounds.newIndex = i
		}

		status, ok := bisectInfo.Status(commit.Sha)
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
	// given that our cache key is a commit sha and a commit count, it's very important that we don't actually try to render pipes
	// when dealing with things like filtered commits.
	cacheKey := pipeSetCacheKey{
		commitSha:   commits[0].Sha,
		commitCount: len(commits),
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

func getBisectStatus(index int, commitSha string, bisectInfo *git_commands.BisectInfo, bisectBounds *bisectBounds) BisectStatus {
	if !bisectInfo.Started() {
		return BisectStatusNone
	}

	if bisectInfo.GetCurrentSha() == commitSha {
		return BisectStatusCurrent
	}

	status, ok := bisectInfo.Status(commitSha)
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
	cherryPickedCommitShaSet *set.Set[string],
	diffName string,
	timeFormat string,
	parseEmoji bool,
	graphLine string,
	fullDescription bool,
	bisectStatus BisectStatus,
	bisectInfo *git_commands.BisectInfo,
	isYouAreHereCommit bool,
) []string {
	shaColor := getShaColor(commit, diffName, cherryPickedCommitShaSet, bisectStatus, bisectInfo)
	bisectString := getBisectStatusText(bisectStatus, bisectInfo)

	actionString := ""
	if commit.Action != models.ActionNone {
		todoString := commit.Action.String()
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
		} else if common.UserConfig.Gui.ExperimentalShowBranchHeads && commit.ExtraInfo != "" {
			tagString = style.FgMagenta.SetBold().Sprint("(*)") + " "
		}
	}

	name := commit.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	if isYouAreHereCommit {
		youAreHere := style.FgYellow.Sprintf("<-- %s ---", common.Tr.YouAreHere)
		name = fmt.Sprintf("%s %s", youAreHere, name)
	}

	authorFunc := authors.ShortAuthor
	if fullDescription {
		authorFunc = authors.LongAuthor
	}

	cols := make([]string, 0, 7)
	if icons.IsIconEnabled() {
		cols = append(cols, shaColor.Sprint(icons.IconForCommit(commit)))
	}
	cols = append(cols, shaColor.Sprint(commit.ShortSha()))
	cols = append(cols, bisectString)
	if fullDescription {
		cols = append(cols, style.FgBlue.Sprint(utils.UnixToDate(commit.UnixTimestamp, timeFormat)))
	}
	cols = append(
		cols,
		actionString,
		authorFunc(commit.AuthorName),
		graphLine+tagString+theme.DefaultTextColor.Sprint(name),
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

func getShaColor(
	commit *models.Commit,
	diffName string,
	cherryPickedCommitShaSet *set.Set[string],
	bisectStatus BisectStatus,
	bisectInfo *git_commands.BisectInfo,
) style.TextStyle {
	if bisectInfo.Started() {
		return getBisectStatusColor(bisectStatus)
	}

	diffed := commit.Sha != "" && commit.Sha == diffName
	shaColor := theme.DefaultTextColor
	switch commit.Status {
	case models.StatusUnpushed:
		shaColor = style.FgRed
	case models.StatusPushed:
		shaColor = style.FgYellow
	case models.StatusMerged:
		shaColor = style.FgGreen
	case models.StatusRebasing:
		shaColor = style.FgBlue
	case models.StatusReflog:
		shaColor = style.FgBlue
	default:
	}

	if diffed {
		shaColor = theme.DiffTerminalColor
	} else if cherryPickedCommitShaSet.Includes(commit.Sha) {
		shaColor = theme.CherryPickedCommitTextStyle
	}

	return shaColor
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
	default:
		return style.FgYellow
	}
}
