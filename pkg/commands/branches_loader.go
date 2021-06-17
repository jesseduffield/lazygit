package commands

import (
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

type getCurrentBranchNameFunc func() (string, string, error)

// BranchesLoader returns a list of Branch objects for the current repo
type BranchesLoader struct {
	*MgrCtx

	getCurrentBranchName getCurrentBranchNameFunc
}

func NewBranchesLoader(mgrCtx *MgrCtx, getCurrentBranchname getCurrentBranchNameFunc) *BranchesLoader {
	return &BranchesLoader{
		MgrCtx:               mgrCtx,
		getCurrentBranchName: getCurrentBranchname,
	}
}

// Build the list of branches for the current repo
func (b *BranchesLoader) Load(reflogCommits []*models.Commit) []*models.Branch {
	branches := b.obtainBranches()

	reflogBranches := b.obtainReflogBranches(reflogCommits)

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
				branches = append(branches[0:j], branches[j+1:]...)
				continue outer
			}
		}
	}

	branches = append(branchesWithRecency, branches...)

	foundHead := false
	for i, branch := range branches {
		if branch.Head {
			foundHead = true
			branch.Recency = "  *"
			branches = append(branches[0:i], branches[i+1:]...)
			branches = append([]*models.Branch{branch}, branches...)
			break
		}
	}
	if !foundHead {
		currentBranchName, currentBranchDisplayName, err := b.getCurrentBranchName()
		if err != nil {
			panic(err)
		}
		branches = append([]*models.Branch{{Name: currentBranchName, DisplayName: currentBranchDisplayName, Head: true, Recency: "  *"}}, branches...)
	}
	return branches
}

func (b *BranchesLoader) obtainBranches() []*models.Branch {
	output, err := b.RunWithOutput(
		b.BuildGitCmdObjFromStr(
			`for-each-ref --sort=-committerdate --format="%(HEAD)|%(refname:short)|%(upstream:short)|%(upstream:track)" refs/heads`,
		),
	)
	if err != nil {
		panic(err)
	}

	trimmedOutput := strings.TrimSpace(output)
	outputLines := strings.Split(trimmedOutput, "\n")
	branches := make([]*models.Branch, 0, len(outputLines))
	for _, line := range outputLines {
		if line == "" {
			continue
		}

		split := strings.Split(line, SEPARATION_CHAR)

		name := strings.TrimPrefix(split[1], "heads/")
		branch := &models.Branch{
			Name:      name,
			Pullables: "?",
			Pushables: "?",
			Head:      split[0] == "*",
		}

		upstreamName := split[2]
		if upstreamName == "" {
			branches = append(branches, branch)
			continue
		}

		branch.UpstreamName = upstreamName

		track := split[3]
		re := regexp.MustCompile(`ahead (\d+)`)
		match := re.FindStringSubmatch(track)
		if len(match) > 1 {
			branch.Pushables = match[1]
		} else {
			branch.Pushables = "0"
		}

		re = regexp.MustCompile(`behind (\d+)`)
		match = re.FindStringSubmatch(track)
		if len(match) > 1 {
			branch.Pullables = match[1]
		} else {
			branch.Pullables = "0"
		}

		branches = append(branches, branch)
	}

	return branches
}

// TODO: only look at the new reflog commits, and otherwise store the recencies in
// int form against the branch to recalculate the time ago
func (b *BranchesLoader) obtainReflogBranches(reflogCommits []*models.Commit) []*models.Branch {
	foundBranchesMap := map[string]bool{}
	re := regexp.MustCompile(`checkout: moving from ([\S]+) to ([\S]+)`)
	reflogBranches := make([]*models.Branch, 0, len(reflogCommits))
	for _, commit := range reflogCommits {
		if match := re.FindStringSubmatch(commit.Name); len(match) == 3 {
			recency := utils.UnixToTimeAgo(commit.UnixTimestamp)
			for _, branchName := range match[1:] {
				if !foundBranchesMap[branchName] {
					foundBranchesMap[branchName] = true
					reflogBranches = append(reflogBranches, &models.Branch{
						Recency: recency,
						Name:    branchName,
					})
				}
			}
		}
	}
	return reflogBranches
}
