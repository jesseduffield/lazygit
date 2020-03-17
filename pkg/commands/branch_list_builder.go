package commands

import (
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"

	"github.com/sirupsen/logrus"
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

// BranchListBuilder returns a list of Branch objects for the current repo
type BranchListBuilder struct {
	Log        *logrus.Entry
	GitCommand *GitCommand
}

// NewBranchListBuilder builds a new branch list builder
func NewBranchListBuilder(log *logrus.Entry, gitCommand *GitCommand) (*BranchListBuilder, error) {
	return &BranchListBuilder{
		Log:        log,
		GitCommand: gitCommand,
	}, nil
}

func (b *BranchListBuilder) obtainCurrentBranchName() string {
	branchName, err := b.GitCommand.CurrentBranchName()
	if err != nil {
		panic(err.Error())
	}
	return branchName
}

func (b *BranchListBuilder) obtainBranches() []*Branch {
	cmdStr := `git branch --format="%(refname:short)|%(upstream:short)|%(upstream:track)"`
	output, err := b.GitCommand.OSCommand.RunCommandWithOutput(cmdStr)
	if err != nil {
		panic(err)
	}

	trimmedOutput := strings.TrimSpace(output)
	outputLines := strings.Split(trimmedOutput, "\n")
	branches := make([]*Branch, len(outputLines))
	for i, line := range outputLines {
		split := strings.Split(line, SEPARATION_CHAR)

		name := split[0]
		branches[i] = &Branch{
			Name:      name,
			Pullables: "?",
			Pushables: "?",
		}
		upstreamName := split[1]
		if upstreamName == "" {
			continue
		}

		branches[i].UpstreamName = upstreamName

		track := split[2]
		re := regexp.MustCompile(`ahead (\d+)`)
		match := re.FindStringSubmatch(track)
		if len(match) > 1 {
			branches[i].Pushables = match[1]
		} else {
			branches[i].Pushables = "0"
		}

		re = regexp.MustCompile(`behind (\d+)`)
		match = re.FindStringSubmatch(track)
		if len(match) > 1 {
			branches[i].Pullables = match[1]
		} else {
			branches[i].Pullables = "0"
		}
	}

	return branches
}

// Build the list of branches for the current repo
func (b *BranchListBuilder) Build() []*Branch {
	currentBranchName := b.obtainCurrentBranchName()
	branches := b.obtainBranches()
	reflogBranches := b.obtainReflogBranches()

	// loop through reflog branches. If there is a match, merge them, then remove it from the branches and keep it in the reflog branches
	branchesWithRecency := make([]*Branch, 0)
outer:
	for _, reflogBranch := range reflogBranches {
		for j, branch := range branches {
			if strings.EqualFold(reflogBranch.Name, branch.Name) {
				branch.Recency = reflogBranch.Recency
				branchesWithRecency = append(branchesWithRecency, branch)
				branches = append(branches[0:j], branches[j+1:]...)
				continue outer
			}
		}
	}

	branches = append(branchesWithRecency, branches...)

	// if it's the current branch we need to pull it up to the top
	for i, branch := range branches {
		if branch.Name == currentBranchName {
			branches = append(branches[0:i], branches[i+1:]...)
			branches = append([]*Branch{branch}, branches...)
			break
		}
	}

	if len(branches) == 0 || branches[0].Name != currentBranchName {
		branches = append([]*Branch{{Name: currentBranchName}}, branches...)
	}

	branches[0].Recency = "  *"

	return branches
}

// A line will have the form '10 days ago master' so we need to strip out the
// useful information from that into timeNumber, timeUnit, and branchName
func branchInfoFromLine(line string) (string, string) {
	// example line: HEAD@{12 minutes ago}|checkout: moving from pulling-from-forks to tim77-patch-1
	r := regexp.MustCompile(`HEAD\@\{([^\s]+) ([^\s]+) ago\}\|.*?([^\s]*)$`)
	matches := r.FindStringSubmatch(strings.TrimSpace(line))
	if len(matches) == 0 {
		return "", ""
	}
	since := matches[1]
	unit := matches[2]
	branchName := matches[3]
	return since + abbreviatedTimeUnit(unit), branchName
}

func abbreviatedTimeUnit(timeUnit string) string {
	r := regexp.MustCompile("s$")
	timeUnit = r.ReplaceAllString(timeUnit, "")
	timeUnitMap := map[string]string{
		"hour":   "h",
		"minute": "m",
		"second": "s",
		"week":   "w",
		"year":   "y",
		"day":    "d",
		"month":  "m",
	}
	return timeUnitMap[timeUnit]
}

func (b *BranchListBuilder) obtainReflogBranches() []*Branch {
	branches := make([]*Branch, 0)
	// if we directly put this string in RunCommandWithOutput the compiler complains because it thinks it's a format string
	unescaped := "git reflog --date=relative --pretty='%gd|%gs' --grep-reflog='checkout: moving' HEAD"
	rawString, err := b.GitCommand.OSCommand.RunCommandWithOutput(unescaped)
	if err != nil {
		return branches
	}

	branchNameMap := map[string]bool{}

	branchLines := utils.SplitLines(rawString)
	for _, line := range branchLines {
		recency, branchName := branchInfoFromLine(line)
		if branchName == "" {
			continue
		}
		if _, ok := branchNameMap[branchName]; ok {
			continue
		}
		branchNameMap[branchName] = true
		branch := &Branch{Name: branchName, Recency: recency}
		branches = append(branches, branch)
	}
	return branches
}
