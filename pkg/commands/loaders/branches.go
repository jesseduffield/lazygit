package loaders

import (
	"regexp"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/go-git/v5/config"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
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

type BranchLoaderConfigCommands interface {
	Branches() (map[string]*config.Branch, error)
}

// BranchLoader returns a list of Branch objects for the current repo
type BranchLoader struct {
	*common.Common
	getRawBranches       func() (string, error)
	getCurrentBranchName func() (string, string, error)
	config               BranchLoaderConfigCommands
}

func NewBranchLoader(
	cmn *common.Common,
	getRawBranches func() (string, error),
	getCurrentBranchName func() (string, string, error),
	config BranchLoaderConfigCommands,
) *BranchLoader {
	return &BranchLoader{
		Common:               cmn,
		getRawBranches:       getRawBranches,
		getCurrentBranchName: getCurrentBranchName,
		config:               config,
	}
}

// Load the list of branches for the current repo
func (self *BranchLoader) Load(reflogCommits []*models.Commit) ([]*models.Branch, error) {
	branches := self.obtainBranches()

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
				branches = slices.Remove(branches, j)
				continue outer
			}
		}
	}

	branches = slices.Prepend(branches, branchesWithRecency...)

	foundHead := false
	for i, branch := range branches {
		if branch.Head {
			foundHead = true
			branch.Recency = "  *"
			branches = slices.Move(branches, i, 0)
			break
		}
	}
	if !foundHead {
		currentBranchName, currentBranchDisplayName, err := self.getCurrentBranchName()
		if err != nil {
			return nil, err
		}
		branches = slices.Prepend(branches, &models.Branch{Name: currentBranchName, DisplayName: currentBranchDisplayName, Head: true, Recency: "  *"})
	}

	configBranches, err := self.config.Branches()
	if err != nil {
		return nil, err
	}

	for _, branch := range branches {
		match := configBranches[branch.Name]
		if match != nil {
			branch.UpstreamRemote = match.Remote
			branch.UpstreamBranch = match.Merge.Short()
		}
	}

	return branches, nil
}

func (self *BranchLoader) obtainBranches() []*models.Branch {
	output, err := self.getRawBranches()
	if err != nil {
		panic(err)
	}

	trimmedOutput := strings.TrimSpace(output)
	outputLines := strings.Split(trimmedOutput, "\n")

	return slices.FilterMap(outputLines, func(line string) (*models.Branch, bool) {
		if line == "" {
			return nil, false
		}

		split := strings.Split(line, "\x00")
		if len(split) != 4 {
			// Ignore line if it isn't separated into 4 parts
			// This is probably a warning message, for more info see:
			// https://github.com/jesseduffield/lazygit/issues/1385#issuecomment-885580439
			return nil, false
		}

		return obtainBranch(split), true
	})
}

// Obtain branch information from parsed line output of getRawBranches()
// split contains the '|' separated tokens in the line of output
func obtainBranch(split []string) *models.Branch {
	name := strings.TrimPrefix(split[1], "heads/")
	branch := &models.Branch{
		Name:      name,
		Pullables: "?",
		Pushables: "?",
		Head:      split[0] == "*",
	}

	upstreamName := split[2]
	if upstreamName == "" {
		// if we're here then it means we do not have a local version of the remote.
		// The branch might still be tracking a remote though, we just don't know
		// how many commits ahead/behind it is
		return branch
	}

	track := split[3]
	if track == "[gone]" {
		branch.UpstreamGone = true
	} else {
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
	}

	return branch
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
