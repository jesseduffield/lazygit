package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/i18n"
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
	Log        *logrus.Entry
	GitCommand *GitCommand
	OSCommand  *oscommands.OSCommand
	Tr         *i18n.TranslationSet
}

// NewCommitListBuilder builds a new commit list builder
func NewCommitListBuilder(
	log *logrus.Entry,
	gitCommand *GitCommand,
	osCommand *oscommands.OSCommand,
	tr *i18n.TranslationSet,
) *CommitListBuilder {
	return &CommitListBuilder{
		Log:        log,
		GitCommand: gitCommand,
		OSCommand:  osCommand,
		Tr:         tr,
	}
}

// extractCommitFromLine takes a line from a git log and extracts the sha, message, date, and tag if present
// then puts them into a commit object
// example input:
// 8ad01fe32fcc20f07bc6693f87aa4977c327f1e1|10 hours ago|Jesse Duffield| (HEAD -> master, tag: v0.15.2)|refresh commits when adding a tag
func (c *CommitListBuilder) extractCommitFromLine(line string) *models.Commit {
	split := strings.Split(line, SEPARATION_CHAR)

	sha := split[0]
	unixTimestamp := split[1]
	author := split[2]
	extraInfo := strings.TrimSpace(split[3])
	parentHashes := split[4]

	message := strings.Join(split[5:], SEPARATION_CHAR)
	tags := []string{}

	if extraInfo != "" {
		re := regexp.MustCompile(`tag: ([^,\)]+)`)
		tagMatch := re.FindStringSubmatch(extraInfo)
		if len(tagMatch) > 1 {
			tags = append(tags, tagMatch[1])
		}
	}

	unitTimestampInt, _ := strconv.Atoi(unixTimestamp)

	return &models.Commit{
		Sha:           sha,
		Name:          message,
		Tags:          tags,
		ExtraInfo:     extraInfo,
		UnixTimestamp: int64(unitTimestampInt),
		Author:        author,
		Parents:       strings.Split(parentHashes, " "),
	}
}

type GetCommitsOptions struct {
	Limit                bool
	FilterPath           string
	IncludeRebaseCommits bool
	RefName              string // e.g. "HEAD" or "my_branch"
	// determines if we show the whole git graph i.e. pass the '--all' flag
	All bool
}

func (c *CommitListBuilder) MergeRebasingCommits(commits []*models.Commit) ([]*models.Commit, error) {
	// chances are we have as many commits as last time so we'll set the capacity to be the old length
	result := make([]*models.Commit, 0, len(commits))
	for i, commit := range commits {
		if commit.Status != "rebasing" { // removing the existing rebase commits so we can add the refreshed ones
			result = append(result, commits[i:]...)
			break
		}
	}

	rebaseMode, err := c.GitCommand.RebaseMode()
	if err != nil {
		return nil, err
	}

	if rebaseMode == "" {
		// not in rebase mode so return original commits
		return result, nil
	}

	rebasingCommits, err := c.getHydratedRebasingCommits(rebaseMode)
	if err != nil {
		return nil, err
	}
	if len(rebasingCommits) > 0 {
		result = append(rebasingCommits, result...)
	}

	return result, nil
}

// GetCommits obtains the commits of the current branch
func (c *CommitListBuilder) GetCommits(opts GetCommitsOptions) ([]*models.Commit, error) {
	commits := []*models.Commit{}
	var rebasingCommits []*models.Commit
	rebaseMode, err := c.GitCommand.RebaseMode()
	if err != nil {
		return nil, err
	}

	if opts.IncludeRebaseCommits && opts.FilterPath == "" {
		var err error
		rebasingCommits, err = c.MergeRebasingCommits(commits)
		if err != nil {
			return nil, err
		}
		commits = append(commits, rebasingCommits...)
	}

	passedFirstPushedCommit := false
	firstPushedCommit, err := c.getFirstPushedCommit(opts.RefName)
	if err != nil {
		// must have no upstream branch so we'll consider everything as pushed
		passedFirstPushedCommit = true
	}

	cmdObj := c.getLogCmd(opts)

	err = c.OSCommand.RunLineOutputCmd(cmdObj, func(line string) (bool, error) {
		if canExtractCommit(line) {
			commit := c.extractCommitFromLine(line)
			if commit.Sha == firstPushedCommit {
				passedFirstPushedCommit = true
			}
			commit.Status = map[bool]string{true: "unpushed", false: "pushed"}[!passedFirstPushedCommit]
			commits = append(commits, commit)
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	if rebaseMode != "" {
		currentCommit := commits[len(rebasingCommits)]
		youAreHere := style.FgYellow.Sprintf("<-- %s ---", c.Tr.YouAreHere)
		currentCommit.Name = fmt.Sprintf("%s %s", youAreHere, currentCommit.Name)
	}

	commits, err = c.setCommitMergedStatuses(opts.RefName, commits)
	if err != nil {
		return nil, err
	}

	return commits, nil
}

func (c *CommitListBuilder) getHydratedRebasingCommits(rebaseMode string) ([]*models.Commit, error) {
	commits, err := c.getRebasingCommits(rebaseMode)
	if err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, nil
	}

	commitShas := make([]string, len(commits))
	for i, commit := range commits {
		commitShas[i] = commit.Sha
	}

	// note that we're not filtering these as we do non-rebasing commits just because
	// I suspect that will cause some damage
	cmdObj := c.OSCommand.NewCmdObj(
		fmt.Sprintf(
			"git show %s --no-patch --oneline %s --abbrev=%d",
			strings.Join(commitShas, " "),
			prettyFormat,
			20,
		),
	)

	hydratedCommits := make([]*models.Commit, 0, len(commits))
	i := 0
	err = c.OSCommand.RunLineOutputCmd(cmdObj, func(line string) (bool, error) {
		if canExtractCommit(line) {
			commit := c.extractCommitFromLine(line)
			matchingCommit := commits[i]
			commit.Action = matchingCommit.Action
			commit.Status = matchingCommit.Status
			hydratedCommits = append(hydratedCommits, commit)
			i++
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return hydratedCommits, nil
}

// getRebasingCommits obtains the commits that we're in the process of rebasing
func (c *CommitListBuilder) getRebasingCommits(rebaseMode string) ([]*models.Commit, error) {
	switch rebaseMode {
	case REBASE_MODE_MERGING:
		return c.getNormalRebasingCommits()
	case REBASE_MODE_INTERACTIVE:
		return c.getInteractiveRebasingCommits()
	default:
		return nil, nil
	}
}

func (c *CommitListBuilder) getNormalRebasingCommits() ([]*models.Commit, error) {
	rewrittenCount := 0
	bytesContent, err := ioutil.ReadFile(filepath.Join(c.GitCommand.DotGitDir, "rebase-apply/rewritten"))
	if err == nil {
		content := string(bytesContent)
		rewrittenCount = len(strings.Split(content, "\n"))
	}

	// we know we're rebasing, so lets get all the files whose names have numbers
	commits := []*models.Commit{}
	err = filepath.Walk(filepath.Join(c.GitCommand.DotGitDir, "rebase-apply"), func(path string, f os.FileInfo, err error) error {
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
		commits = append([]*models.Commit{commit}, commits...)
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
func (c *CommitListBuilder) getInteractiveRebasingCommits() ([]*models.Commit, error) {
	bytesContent, err := ioutil.ReadFile(filepath.Join(c.GitCommand.DotGitDir, "rebase-merge/git-rebase-todo"))
	if err != nil {
		c.Log.Error(fmt.Sprintf("error occurred reading git-rebase-todo: %s", err.Error()))
		// we assume an error means the file doesn't exist so we just return
		return nil, nil
	}

	commits := []*models.Commit{}
	lines := strings.Split(string(bytesContent), "\n")
	for _, line := range lines {
		if line == "" || line == "noop" {
			return commits, nil
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		splitLine := strings.Split(line, " ")
		commits = append([]*models.Commit{{
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
func (c *CommitListBuilder) commitFromPatch(content string) (*models.Commit, error) {
	lines := strings.Split(content, "\n")
	sha := strings.Split(lines[0], " ")[1]
	name := strings.TrimPrefix(lines[3], "Subject: ")
	return &models.Commit{
		Sha:    sha,
		Name:   name,
		Status: "rebasing",
	}, nil
}

func (c *CommitListBuilder) setCommitMergedStatuses(refName string, commits []*models.Commit) ([]*models.Commit, error) {
	ancestor, err := c.getMergeBase(refName)
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

func (c *CommitListBuilder) getMergeBase(refName string) (string, error) {
	currentBranch, _, err := c.GitCommand.CurrentBranchName()
	if err != nil {
		return "", err
	}

	baseBranch := "master"
	if strings.HasPrefix(currentBranch, "feature/") {
		baseBranch = "develop"
	}

	// swallowing error because it's not a big deal; probably because there are no commits yet
	output, _ := c.OSCommand.RunWithOutput(c.OSCommand.NewCmdObj(fmt.Sprintf("git merge-base %s %s", c.OSCommand.Quote(refName), c.OSCommand.Quote(baseBranch))))
	return ignoringWarnings(output), nil
}

func ignoringWarnings(commandOutput string) string {
	trimmedOutput := strings.TrimSpace(commandOutput)
	split := strings.Split(trimmedOutput, "\n")
	// need to get last line in case the first line is a warning about how the error is ambiguous.
	// At some point we should find a way to make it unambiguous
	lastLine := split[len(split)-1]

	return lastLine
}

// getFirstPushedCommit returns the first commit SHA which has been pushed to the ref's upstream.
// all commits above this are deemed unpushed and marked as such.
func (c *CommitListBuilder) getFirstPushedCommit(refName string) (string, error) {
	output, err := c.OSCommand.RunWithOutput(c.OSCommand.NewCmdObj(fmt.Sprintf("git merge-base %s %s@{u}", c.OSCommand.Quote(refName), c.OSCommand.Quote(refName))))
	if err != nil {
		return "", err
	}

	return ignoringWarnings(output), nil
}

// getLog gets the git log.
func (c *CommitListBuilder) getLogCmd(opts GetCommitsOptions) oscommands.ICmdObj {
	limitFlag := ""
	if opts.Limit {
		limitFlag = "-300"
	}

	filterFlag := ""
	if opts.FilterPath != "" {
		filterFlag = fmt.Sprintf(" --follow -- %s", c.OSCommand.Quote(opts.FilterPath))
	}

	config := c.GitCommand.UserConfig.Git.Log

	orderFlag := "--" + config.Order
	allFlag := ""
	if opts.All {
		allFlag = " --all"
	}

	return c.OSCommand.NewCmdObj(
		fmt.Sprintf(
			"git log %s %s %s --oneline %s %s --abbrev=%d %s",
			c.OSCommand.Quote(opts.RefName),
			orderFlag,
			allFlag,
			prettyFormat,
			limitFlag,
			20,
			filterFlag,
		),
	)
}

var prettyFormat = fmt.Sprintf(
	"--pretty=format:\"%%H%s%%at%s%%aN%s%%d%s%%p%s%%s\"",
	SEPARATION_CHAR,
	SEPARATION_CHAR,
	SEPARATION_CHAR,
	SEPARATION_CHAR,
	SEPARATION_CHAR,
)

func canExtractCommit(line string) bool {
	return strings.Split(line, " ")[0] != "gpg:"
}
