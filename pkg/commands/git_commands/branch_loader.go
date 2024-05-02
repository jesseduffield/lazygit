package git_commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/go-git/v5/config"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
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

type BranchInfo struct {
	RefName      string
	DisplayName  string // e.g. '(HEAD detached at 123asdf)'
	DetachedHead bool
}

// BranchLoader returns a list of Branch objects for the current repo
type BranchLoader struct {
	*common.Common
	cmd                  oscommands.ICmdObjBuilder
	getCurrentBranchInfo func() (BranchInfo, error)
	config               BranchLoaderConfigCommands
}

func NewBranchLoader(
	cmn *common.Common,
	cmd oscommands.ICmdObjBuilder,
	getCurrentBranchInfo func() (BranchInfo, error),
	config BranchLoaderConfigCommands,
) *BranchLoader {
	return &BranchLoader{
		Common:               cmn,
		cmd:                  cmd,
		getCurrentBranchInfo: getCurrentBranchInfo,
		config:               config,
	}
}

// Load the list of branches for the current repo
func (self *BranchLoader) Load(reflogCommits []*models.Commit) ([]*models.Branch, error) {
	branches := self.obtainBranches()

	if self.AppState.LocalBranchSortOrder == "recency" {
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
		slices.SortFunc(branches, func(a *models.Branch, b *models.Branch) bool {
			return a.Name < b.Name
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

		storeCommitDateAsRecency := self.AppState.LocalBranchSortOrder != "recency"
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
	switch strings.ToLower(self.AppState.LocalBranchSortOrder) {
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
	subject := split[4]
	commitHash := split[5]
	commitDate := split[6]

	name := strings.TrimPrefix(fullName, "heads/")
	pushables, pullables, gone := parseUpstreamInfo(upstreamName, track)

	recency := ""
	if storeCommitDateAsRecency {
		if unixTimestamp, err := strconv.ParseInt(commitDate, 10, 64); err == nil {
			recency = utils.UnixToTimeAgo(unixTimestamp)
		}
	}

	return &models.Branch{
		Name:         name,
		Recency:      recency,
		Pushables:    pushables,
		Pullables:    pullables,
		UpstreamGone: gone,
		Head:         headMarker == "*",
		Subject:      subject,
		CommitHash:   commitHash,
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

	pushables := parseDifference(track, `ahead (\d+)`)
	pullables := parseDifference(track, `behind (\d+)`)

	return pushables, pullables, false
}

func parseDifference(track string, regexStr string) string {
	re := regexp.MustCompile(regexStr)
	match := re.FindStringSubmatch(track)
	if len(match) > 1 {
		return match[1]
	} else {
		return "0"
	}
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
