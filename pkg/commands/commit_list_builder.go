package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

// context:
// here we get the commits from git log but format them to show whether they're
// unpushed/pushed/merged into the base branch or not, or if they're yet to
// be processed as part of a rebase (these won't appear in git log but we
// grab them from the rebase-related files in the .git directory to show them

// if we find out we need to use one of these functions in the git.go file, we
// can just pull them out of here and put them there and then call them from in here

const SEPARATION_CHAR = "|"

// CommitListBuilder returns a list of Branch objects for the current repo
type CommitListBuilder struct {
	Log                 *logrus.Entry
	GitCommand          *GitCommand
	OSCommand           *OSCommand
	Tr                  *i18n.Localizer
	CherryPickedCommits []*Commit
}

// NewCommitListBuilder builds a new commit list builder
func NewCommitListBuilder(log *logrus.Entry, gitCommand *GitCommand, osCommand *OSCommand, tr *i18n.Localizer, cherryPickedCommits []*Commit) (*CommitListBuilder, error) {
	return &CommitListBuilder{
		Log:                 log,
		GitCommand:          gitCommand,
		OSCommand:           osCommand,
		Tr:                  tr,
		CherryPickedCommits: cherryPickedCommits,
	}, nil
}

// extractCommitFromLine takes a line from a git log and extracts the sha, message, date, and tag if present
// then puts them into a commit object
// example input:
// 8ad01fe32fcc20f07bc6693f87aa4977c327f1e1|10 hours ago|Jesse Duffield| (HEAD -> master, tag: v0.15.2)|refresh commits when adding a tag
func (c *CommitListBuilder) extractCommitFromLine(line string) *Commit {
	split := strings.Split(line, SEPARATION_CHAR)

	sha := split[0]
	unixTimestamp := split[1]
	author := split[2]
	extraInfo := strings.TrimSpace(split[3])
	message := strings.Join(split[4:], SEPARATION_CHAR)
	tags := []string{}

	if extraInfo != "" {
		re := regexp.MustCompile(`tag: ([^,\)]+)`)
		tagMatch := re.FindStringSubmatch(extraInfo)
		if len(tagMatch) > 1 {
			tags = append(tags, tagMatch[1])
		}
	}

	unitTimestampInt, _ := strconv.Atoi(unixTimestamp)

	return &Commit{
		Sha:           sha,
		Name:          message,
		Tags:          tags,
		ExtraInfo:     extraInfo,
		UnixTimestamp: int64(unitTimestampInt),
		Author:        author,
	}
}

type GetCommitsOptions struct {
	Limit      bool
	FilterPath string
}

// GetCommits obtains the commits of the current branch
func (c *CommitListBuilder) GetCommits(options GetCommitsOptions) ([]*Commit, error) {
	commits := []*Commit{}
	var rebasingCommits []*Commit
	rebaseMode, err := c.GitCommand.RebaseMode()
	if err != nil {
		return nil, err
	}
	if rebaseMode != "" && options.FilterPath == "" {
		// here we want to also prepend the commits that we're in the process of rebasing
		rebasingCommits, err = c.getRebasingCommits(rebaseMode)
		if err != nil {
			return nil, err
		}
		if len(rebasingCommits) > 0 {
			commits = append(commits, rebasingCommits...)
		}
	}

	unpushedCommits := c.getUnpushedCommits()
	cmd := c.getLogCmd(options)

	err = RunLineOutputCmd(cmd, func(line string) (bool, error) {
		if strings.Split(line, " ")[0] != "gpg:" {
			commit := c.extractCommitFromLine(line)
			_, unpushed := unpushedCommits[commit.ShortSha()]
			commit.Status = map[bool]string{true: "unpushed", false: "pushed"}[unpushed]
			commits = append(commits, commit)
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	if rebaseMode != "" {
		currentCommit := commits[len(rebasingCommits)]
		blue := color.New(color.FgYellow)
		youAreHere := blue.Sprintf("<-- %s ---", c.Tr.SLocalize("YouAreHere"))
		currentCommit.Name = fmt.Sprintf("%s %s", youAreHere, currentCommit.Name)
	}

	commits, err = c.setCommitMergedStatuses(commits)
	if err != nil {
		return nil, err
	}

	return commits, nil
}

// getRebasingCommits obtains the commits that we're in the process of rebasing
func (c *CommitListBuilder) getRebasingCommits(rebaseMode string) ([]*Commit, error) {
	switch rebaseMode {
	case "normal":
		return c.getNormalRebasingCommits()
	case "interactive":
		return c.getInteractiveRebasingCommits()
	default:
		return nil, nil
	}
}

func (c *CommitListBuilder) getNormalRebasingCommits() ([]*Commit, error) {
	rewrittenCount := 0
	bytesContent, err := ioutil.ReadFile(fmt.Sprintf("%s/rebase-apply/rewritten", c.GitCommand.DotGitDir))
	if err == nil {
		content := string(bytesContent)
		rewrittenCount = len(strings.Split(content, "\n"))
	}

	// we know we're rebasing, so lets get all the files whose names have numbers
	commits := []*Commit{}
	err = filepath.Walk(fmt.Sprintf("%s/rebase-apply", c.GitCommand.DotGitDir), func(path string, f os.FileInfo, err error) error {
		if rewrittenCount > 0 {
			rewrittenCount--
			return nil
		}
		if err != nil {
			return err
		}
		re := regexp.MustCompile(`^\d+$`)
		if !re.MatchString(f.Name()) {
			return nil
		}
		bytesContent, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(bytesContent)
		commit, err := c.commitFromPatch(content)
		if err != nil {
			return err
		}
		commits = append([]*Commit{commit}, commits...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return commits, nil
}

// git-rebase-todo example:
// pick ac446ae94ee560bdb8d1d057278657b251aaef17 ac446ae
// pick afb893148791a2fbd8091aeb81deba4930c73031 afb8931

// git-rebase-todo.backup example:
// pick 49cbba374296938ea86bbd4bf4fee2f6ba5cccf6 third commit on master
// pick ac446ae94ee560bdb8d1d057278657b251aaef17 blah  commit on master
// pick afb893148791a2fbd8091aeb81deba4930c73031 fourth commit on master

// getInteractiveRebasingCommits takes our git-rebase-todo and our git-rebase-todo.backup files
// and extracts out the sha and names of commits that we still have to go
// in the rebase:
func (c *CommitListBuilder) getInteractiveRebasingCommits() ([]*Commit, error) {
	bytesContent, err := ioutil.ReadFile(fmt.Sprintf("%s/rebase-merge/git-rebase-todo", c.GitCommand.DotGitDir))
	if err != nil {
		c.Log.Info(fmt.Sprintf("error occurred reading git-rebase-todo: %s", err.Error()))
		// we assume an error means the file doesn't exist so we just return
		return nil, nil
	}

	commits := []*Commit{}
	lines := strings.Split(string(bytesContent), "\n")
	for _, line := range lines {
		if line == "" || line == "noop" {
			return commits, nil
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		splitLine := strings.Split(line, " ")
		commits = append([]*Commit{{
			Sha:    splitLine[1],
			Name:   strings.Join(splitLine[2:], " "),
			Status: "rebasing",
			Action: splitLine[0],
		}}, commits...)
	}

	return commits, nil
}

// assuming the file starts like this:
// From e93d4193e6dd45ca9cf3a5a273d7ba6cd8b8fb20 Mon Sep 17 00:00:00 2001
// From: Lazygit Tester <test@example.com>
// Date: Wed, 5 Dec 2018 21:03:23 +1100
// Subject: second commit on master
func (c *CommitListBuilder) commitFromPatch(content string) (*Commit, error) {
	lines := strings.Split(content, "\n")
	sha := strings.Split(lines[0], " ")[1]
	name := strings.TrimPrefix(lines[3], "Subject: ")
	return &Commit{
		Sha:    sha,
		Name:   name,
		Status: "rebasing",
	}, nil
}

func (c *CommitListBuilder) setCommitMergedStatuses(commits []*Commit) ([]*Commit, error) {
	ancestor, err := c.getMergeBase()
	if err != nil {
		return nil, err
	}
	if ancestor == "" {
		return commits, nil
	}
	passedAncestor := false
	for i, commit := range commits {
		if strings.HasPrefix(ancestor, commit.Sha) {
			passedAncestor = true
		}
		if commit.Status != "pushed" {
			continue
		}
		if passedAncestor {
			commits[i].Status = "merged"
		}
	}
	return commits, nil
}

func (c *CommitListBuilder) getMergeBase() (string, error) {
	currentBranch, _, err := c.GitCommand.CurrentBranchName()
	if err != nil {
		return "", err
	}

	baseBranch := "master"
	if strings.HasPrefix(currentBranch, "feature/") {
		baseBranch = "develop"
	}

	// swallowing error because it's not a big deal; probably because there are no commits yet
	output, _ := c.OSCommand.RunCommandWithOutput("git merge-base HEAD %s", baseBranch)
	return output, nil
}

// getUnpushedCommits Returns the sha's of the commits that have not yet been pushed
// to the remote branch of the current branch, a map is returned to ease look up
func (c *CommitListBuilder) getUnpushedCommits() map[string]bool {
	pushables := map[string]bool{}
	o, err := c.OSCommand.RunCommandWithOutput("git rev-list @{u}..HEAD --abbrev-commit --abbrev=8")
	if err != nil {
		return pushables
	}
	for _, p := range utils.SplitLines(o) {
		pushables[p] = true
	}

	return pushables
}

// getLog gets the git log.
func (c *CommitListBuilder) getLogCmd(options GetCommitsOptions) *exec.Cmd {
	limitFlag := ""
	if options.Limit {
		limitFlag = "-300"
	}

	filterFlag := ""
	if options.FilterPath != "" {
		filterFlag = fmt.Sprintf(" --follow -- %s", c.OSCommand.Quote(options.FilterPath))
	}

	return c.OSCommand.ExecutableFromString(fmt.Sprintf("git log --oneline --pretty=format:\"%%H%s%%at%s%%aN%s%%d%s%%s\" %s --abbrev=%d --date=unix %s", SEPARATION_CHAR, SEPARATION_CHAR, SEPARATION_CHAR, SEPARATION_CHAR, limitFlag, 20, filterFlag))
}
