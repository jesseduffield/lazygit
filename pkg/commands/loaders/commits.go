package loaders

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

// context:
// here we get the commits from git log but format them to show whether they're
// unpushed/pushed/merged into the base branch or not, or if they're yet to
// be processed as part of a rebase (these won't appear in git log but we
// grab them from the rebase-related files in the .git directory to show them

// CommitLoader returns a list of Commit objects for the current repo
type CommitLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder

	getCurrentBranchName func() (string, string, error)
	getRebaseMode        func() (enums.RebaseMode, error)
	readFile             func(filename string) ([]byte, error)
	walkFiles            func(root string, fn filepath.WalkFunc) error
	dotGitDir            string
}

// making our dependencies explicit for the sake of easier testing
func NewCommitLoader(
	cmn *common.Common,
	cmd oscommands.ICmdObjBuilder,
	dotGitDir string,
	getCurrentBranchName func() (string, string, error),
	getRebaseMode func() (enums.RebaseMode, error),
) *CommitLoader {
	return &CommitLoader{
		Common:               cmn,
		cmd:                  cmd,
		getCurrentBranchName: getCurrentBranchName,
		getRebaseMode:        getRebaseMode,
		readFile:             ioutil.ReadFile,
		walkFiles:            filepath.Walk,
		dotGitDir:            dotGitDir,
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

// GetCommits obtains the commits of the current branch
func (self *CommitLoader) GetCommits(opts GetCommitsOptions) ([]*models.Commit, error) {
	commits := []*models.Commit{}
	var rebasingCommits []*models.Commit
	rebaseMode, err := self.getRebaseMode()
	if err != nil {
		return nil, err
	}

	if opts.IncludeRebaseCommits && opts.FilterPath == "" {
		var err error
		rebasingCommits, err = self.MergeRebasingCommits(commits)
		if err != nil {
			return nil, err
		}
		commits = append(commits, rebasingCommits...)
	}

	passedFirstPushedCommit := false
	firstPushedCommit, err := self.getFirstPushedCommit(opts.RefName)
	if err != nil {
		// must have no upstream branch so we'll consider everything as pushed
		passedFirstPushedCommit = true
	}

	err = self.getLogCmd(opts).RunAndProcessLines(func(line string) (bool, error) {
		commit := self.extractCommitFromLine(line)
		if commit.Sha == firstPushedCommit {
			passedFirstPushedCommit = true
		}
		commit.Status = map[bool]string{true: "unpushed", false: "pushed"}[!passedFirstPushedCommit]
		commits = append(commits, commit)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return commits, nil
	}

	if rebaseMode != enums.REBASE_MODE_NONE {
		currentCommit := commits[len(rebasingCommits)]
		youAreHere := style.FgYellow.Sprintf("<-- %s ---", self.Tr.YouAreHere)
		currentCommit.Name = fmt.Sprintf("%s %s", youAreHere, currentCommit.Name)
	}

	commits, err = self.setCommitMergedStatuses(opts.RefName, commits)
	if err != nil {
		return nil, err
	}

	return commits, nil
}

func (self *CommitLoader) MergeRebasingCommits(commits []*models.Commit) ([]*models.Commit, error) {
	// chances are we have as many commits as last time so we'll set the capacity to be the old length
	result := make([]*models.Commit, 0, len(commits))
	for i, commit := range commits {
		if commit.Status != "rebasing" { // removing the existing rebase commits so we can add the refreshed ones
			result = append(result, commits[i:]...)
			break
		}
	}

	rebaseMode, err := self.getRebaseMode()
	if err != nil {
		return nil, err
	}

	if rebaseMode == enums.REBASE_MODE_NONE {
		// not in rebase mode so return original commits
		return result, nil
	}

	rebasingCommits, err := self.getHydratedRebasingCommits(rebaseMode)
	if err != nil {
		return nil, err
	}
	if len(rebasingCommits) > 0 {
		result = append(rebasingCommits, result...)
	}

	return result, nil
}

// extractCommitFromLine takes a line from a git log and extracts the sha, message, date, and tag if present
// then puts them into a commit object
// example input:
// 8ad01fe32fcc20f07bc6693f87aa4977c327f1e1|10 hours ago|Jesse Duffield| (HEAD -> master, tag: v0.15.2)|refresh commits when adding a tag
func (self *CommitLoader) extractCommitFromLine(line string) *models.Commit {
	split := strings.SplitN(line, "\x00", 7)

	sha := split[0]
	unixTimestamp := split[1]
	authorName := split[2]
	authorEmail := split[3]
	extraInfo := strings.TrimSpace(split[4])
	parentHashes := split[5]
	message := split[6]

	tags := []string{}

	if extraInfo != "" {
		re := regexp.MustCompile(`tag: ([^,\)]+)`)
		tagMatch := re.FindStringSubmatch(extraInfo)
		if len(tagMatch) > 1 {
			tags = append(tags, tagMatch[1])
		}
	}

	unitTimestampInt, _ := strconv.Atoi(unixTimestamp)

	parents := []string{}
	if len(parentHashes) > 0 {
		parents = strings.Split(parentHashes, " ")
	}

	return &models.Commit{
		Sha:           sha,
		Name:          message,
		Tags:          tags,
		ExtraInfo:     extraInfo,
		UnixTimestamp: int64(unitTimestampInt),
		AuthorName:    authorName,
		AuthorEmail:   authorEmail,
		Parents:       parents,
	}
}

func (self *CommitLoader) getHydratedRebasingCommits(rebaseMode enums.RebaseMode) ([]*models.Commit, error) {
	commits, err := self.getRebasingCommits(rebaseMode)
	if err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, nil
	}

	commitShas := slices.Map(commits, func(commit *models.Commit) string {
		return commit.Sha
	})

	// note that we're not filtering these as we do non-rebasing commits just because
	// I suspect that will cause some damage
	cmdObj := self.cmd.New(
		fmt.Sprintf(
			"git -c log.showSignature=false show %s --no-patch --oneline %s --abbrev=%d",
			strings.Join(commitShas, " "),
			prettyFormat,
			20,
		),
	).DontLog()

	hydratedCommits := make([]*models.Commit, 0, len(commits))
	i := 0
	err = cmdObj.RunAndProcessLines(func(line string) (bool, error) {
		commit := self.extractCommitFromLine(line)
		matchingCommit := commits[i]
		commit.Action = matchingCommit.Action
		commit.Status = matchingCommit.Status
		hydratedCommits = append(hydratedCommits, commit)
		i++
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return hydratedCommits, nil
}

// getRebasingCommits obtains the commits that we're in the process of rebasing
func (self *CommitLoader) getRebasingCommits(rebaseMode enums.RebaseMode) ([]*models.Commit, error) {
	switch rebaseMode {
	case enums.REBASE_MODE_MERGING:
		return self.getNormalRebasingCommits()
	case enums.REBASE_MODE_INTERACTIVE:
		return self.getInteractiveRebasingCommits()
	default:
		return nil, nil
	}
}

func (self *CommitLoader) getNormalRebasingCommits() ([]*models.Commit, error) {
	rewrittenCount := 0
	bytesContent, err := self.readFile(filepath.Join(self.dotGitDir, "rebase-apply/rewritten"))
	if err == nil {
		content := string(bytesContent)
		rewrittenCount = len(strings.Split(content, "\n"))
	}

	// we know we're rebasing, so lets get all the files whose names have numbers
	commits := []*models.Commit{}
	err = self.walkFiles(filepath.Join(self.dotGitDir, "rebase-apply"), func(path string, f os.FileInfo, err error) error {
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
		bytesContent, err := self.readFile(path)
		if err != nil {
			return err
		}
		content := string(bytesContent)
		commit := self.commitFromPatch(content)
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
func (self *CommitLoader) getInteractiveRebasingCommits() ([]*models.Commit, error) {
	bytesContent, err := self.readFile(filepath.Join(self.dotGitDir, "rebase-merge/git-rebase-todo"))
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred reading git-rebase-todo: %s", err.Error()))
		// we assume an error means the file doesn't exist so we just return
		return nil, nil
	}

	commits := []*models.Commit{}

	todos, err := todo.Parse(bytes.NewBuffer(bytesContent))
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred while parsing git-rebase-todo file: %s", err.Error()))
		return nil, nil
	}

	for _, t := range todos {
		if t.Commit == "" {
			// Command does not have a commit associated, skip
			continue
		}
		commits = slices.Prepend(commits, &models.Commit{
			Sha:    t.Commit,
			Name:   t.Msg,
			Status: "rebasing",
			Action: t.Command.String(),
		})
	}

	return commits, nil
}

// assuming the file starts like this:
// From e93d4193e6dd45ca9cf3a5a273d7ba6cd8b8fb20 Mon Sep 17 00:00:00 2001
// From: Lazygit Tester <test@example.com>
// Date: Wed, 5 Dec 2018 21:03:23 +1100
// Subject: second commit on master
func (self *CommitLoader) commitFromPatch(content string) *models.Commit {
	lines := strings.Split(content, "\n")
	sha := strings.Split(lines[0], " ")[1]
	name := strings.TrimPrefix(lines[3], "Subject: ")
	return &models.Commit{
		Sha:    sha,
		Name:   name,
		Status: "rebasing",
	}
}

func (self *CommitLoader) setCommitMergedStatuses(refName string, commits []*models.Commit) ([]*models.Commit, error) {
	ancestor, err := self.getMergeBase(refName)
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

func (self *CommitLoader) getMergeBase(refName string) (string, error) {
	currentBranch, _, err := self.getCurrentBranchName()
	if err != nil {
		return "", err
	}

	baseBranch := "master"
	if strings.HasPrefix(currentBranch, "feature/") {
		baseBranch = "develop"
	}

	// swallowing error because it's not a big deal; probably because there are no commits yet
	output, _ := self.cmd.New(fmt.Sprintf("git merge-base %s %s", self.cmd.Quote(refName), self.cmd.Quote(baseBranch))).DontLog().RunWithOutput()
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
func (self *CommitLoader) getFirstPushedCommit(refName string) (string, error) {
	output, err := self.cmd.
		New(
			fmt.Sprintf("git merge-base %s %s@{u}", self.cmd.Quote(refName), self.cmd.Quote(refName)),
		).
		DontLog().
		RunWithOutput()
	if err != nil {
		return "", err
	}

	return ignoringWarnings(output), nil
}

// getLog gets the git log.
func (self *CommitLoader) getLogCmd(opts GetCommitsOptions) oscommands.ICmdObj {
	limitFlag := ""
	if opts.Limit {
		limitFlag = " -300"
	}

	filterFlag := ""
	if opts.FilterPath != "" {
		filterFlag = fmt.Sprintf(" --follow -- %s", self.cmd.Quote(opts.FilterPath))
	}

	config := self.UserConfig.Git.Log

	orderFlag := "--" + config.Order
	allFlag := ""
	if opts.All {
		allFlag = " --all"
	}

	return self.cmd.New(
		fmt.Sprintf(
			"git -c log.showSignature=false log %s %s %s --oneline %s%s --abbrev=%d%s",
			self.cmd.Quote(opts.RefName),
			orderFlag,
			allFlag,
			prettyFormat,
			limitFlag,
			40,
			filterFlag,
		),
	).DontLog()
}

var prettyFormat = fmt.Sprintf(
	"--pretty=format:\"%%H%s%%at%s%%aN%s%%ae%s%%d%s%%p%s%%s\"",
	NULL_CODE,
	NULL_CODE,
	NULL_CODE,
	NULL_CODE,
	NULL_CODE,
	NULL_CODE,
)

const NULL_CODE = "%x00"
