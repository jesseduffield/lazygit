package presentation

import (
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/authors"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/graph"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/kyokomi/emoji/v2"
)

type pipeSetCacheKey struct {
	commitSha   string
	commitCount int
}

var pipeSetCache = make(map[pipeSetCacheKey][][]*graph.Pipe)
var mutex sync.Mutex

type BisectProgress int

const (
	BeforeNewCommit BisectProgress = iota
	InbetweenCommits
	AfterOldCommit
)

func GetCommitListDisplayStrings(
	commits []*models.Commit,
	fullDescription bool,
	cherryPickedCommitShaMap map[string]bool,
	diffName string,
	parseEmoji bool,
	selectedCommitSha string,
	startIdx int,
	length int,
	showGraph bool,
	bisectInfo *git_commands.BisectInfo,
) [][]string {
	mutex.Lock()
	defer mutex.Unlock()

	if len(commits) == 0 {
		return nil
	}

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
			return authors.AuthorStyle(commit.Author)
		}
		pipeSets = graph.GetPipeSets(commits, getStyle)
		pipeSetCache[cacheKey] = pipeSets
	}

	if startIdx > len(commits) {
		return nil
	}
	end := startIdx + length
	if end > len(commits)-1 {
		end = len(commits) - 1
	}

	filteredCommits := commits[startIdx : end+1]

	var getGraphLine func(int) string
	if showGraph {
		filteredPipeSets := pipeSets[startIdx : end+1]
		graphLines := graph.RenderAux(filteredPipeSets, filteredCommits, selectedCommitSha)
		getGraphLine = func(idx int) string { return graphLines[idx] }
	} else {
		getGraphLine = func(idx int) string { return "" }
	}

	lines := make([][]string, 0, len(filteredCommits))
	bisectProgress := BeforeNewCommit
	var bisectStatus BisectStatus
	for i, commit := range filteredCommits {
		bisectStatus, bisectProgress = getBisectStatus(commit.Sha, bisectInfo, bisectProgress)
		lines = append(lines, displayCommit(
			commit,
			cherryPickedCommitShaMap,
			diffName,
			parseEmoji,
			getGraphLine(i),
			fullDescription,
			bisectStatus,
			bisectInfo,
		))
	}
	return lines
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

func getBisectStatus(commitSha string, bisectInfo *git_commands.BisectInfo, bisectProgress BisectProgress) (BisectStatus, BisectProgress) {
	if !bisectInfo.Started() {
		return BisectStatusNone, bisectProgress
	}

	if bisectInfo.GetCurrentSha() == commitSha {
		return BisectStatusCurrent, bisectProgress
	}

	status, ok := bisectInfo.Status(commitSha)
	if ok {
		switch status {
		case git_commands.BisectStatusNew:
			return BisectStatusNew, InbetweenCommits
		case git_commands.BisectStatusOld:
			return BisectStatusOld, AfterOldCommit
		case git_commands.BisectStatusSkipped:
			return BisectStatusSkipped, bisectProgress
		}
	} else {
		if bisectProgress == InbetweenCommits {
			return BisectStatusCandidate, bisectProgress
		} else {
			return BisectStatusNone, bisectProgress
		}
	}

	// should never land here
	return BisectStatusNone, bisectProgress
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
	}

	return ""
}

func displayCommit(
	commit *models.Commit,
	cherryPickedCommitShaMap map[string]bool,
	diffName string,
	parseEmoji bool,
	graphLine string,
	fullDescription bool,
	bisectStatus BisectStatus,
	bisectInfo *git_commands.BisectInfo,
) []string {
	shaColor := getShaColor(commit, diffName, cherryPickedCommitShaMap, bisectStatus, bisectInfo)
	bisectString := getBisectStatusText(bisectStatus, bisectInfo)

	actionString := ""
	if commit.Action != "" {
		actionString = actionColorMap(commit.Action).Sprint(commit.Action) + " "
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
	}

	name := commit.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	authorFunc := authors.ShortAuthor
	if fullDescription {
		authorFunc = authors.LongAuthor
	}

	cols := make([]string, 0, 5)
	cols = append(cols, shaColor.Sprint(commit.ShortSha()))
	cols = append(cols, bisectString)
	if fullDescription {
		cols = append(cols, style.FgBlue.Sprint(utils.UnixToDate(commit.UnixTimestamp)))
	}
	cols = append(
		cols,
		actionString,
		authorFunc(commit.Author),
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
	cherryPickedCommitShaMap map[string]bool,
	bisectStatus BisectStatus,
	bisectInfo *git_commands.BisectInfo,
) style.TextStyle {
	if bisectInfo.Started() {
		return getBisectStatusColor(bisectStatus)
	}

	diffed := commit.Sha == diffName
	shaColor := theme.DefaultTextColor
	switch commit.Status {
	case "unpushed":
		shaColor = style.FgRed
	case "pushed":
		shaColor = style.FgYellow
	case "merged":
		shaColor = style.FgGreen
	case "rebasing":
		shaColor = style.FgBlue
	case "reflog":
		shaColor = style.FgBlue
	}

	if diffed {
		shaColor = theme.DiffTerminalColor
	} else if cherryPickedCommitShaMap[commit.Sha] {
		shaColor = theme.CherryPickedCommitTextStyle
	}

	return shaColor
}

func actionColorMap(str string) style.TextStyle {
	switch str {
	case "pick":
		return style.FgCyan
	case "drop":
		return style.FgRed
	case "edit":
		return style.FgGreen
	case "fixup":
		return style.FgMagenta
	default:
		return style.FgYellow
	}
}
