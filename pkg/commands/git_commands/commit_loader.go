package git_commands

import (
	"bytes"
	"fmt"
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
	"github.com/samber/lo"
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

	getRebaseMode func() (enums.RebaseMode, error)
	readFile      func(filename string) ([]byte, error)
	walkFiles     func(root string, fn filepath.WalkFunc) error
	dotGitDir     string
	// List of main branches that exist in the repo.
	// We use these to obtain the merge base of the branch.
	// When nil, we're yet to obtain the list of existing main branches.
	// When an empty slice, we've obtained the list and it's empty.
	mainBranches []string
	*GitCommon
}

// making our dependencies explicit for the sake of easier testing
func NewCommitLoader(
	cmn *common.Common,
	cmd oscommands.ICmdObjBuilder,
	dotGitDir string,
	getRebaseMode func() (enums.RebaseMode, error),
	gitCommon *GitCommon,
) *CommitLoader {
	return &CommitLoader{
		Common:        cmn,
		cmd:           cmd,
		getRebaseMode: getRebaseMode,
		readFile:      os.ReadFile,
		walkFiles:     filepath.Walk,
		dotGitDir:     dotGitDir,
		mainBranches:  nil,
		GitCommon:     gitCommon,
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
		commit.Status = map[bool]models.CommitStatus{true: models.StatusUnpushed, false: models.StatusPushed}[!passedFirstPushedCommit]
		commits = append(commits, commit)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return commits, nil
	}

	commits = self.setCommitMergedStatuses(opts.RefName, commits)

	return commits, nil
}

func (self *CommitLoader) MergeRebasingCommits(commits []*models.Commit) ([]*models.Commit, error) {
	// chances are we have as many commits as last time so we'll set the capacity to be the old length
	result := make([]*models.Commit, 0, len(commits))
	for i, commit := range commits {
		if !commit.IsTODO() { // removing the existing rebase commits so we can add the refreshed ones
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
		extraInfoFields := strings.Split(extraInfo, ",")
		for _, extraInfoField := range extraInfoFields {
			extraInfoField = strings.TrimSpace(extraInfoField)
			re := regexp.MustCompile(`tag: (.+)`)
			tagMatch := re.FindStringSubmatch(extraInfoField)
			if len(tagMatch) > 1 {
				tags = append(tags, tagMatch[1])
			}
		}

		extraInfo = "(" + extraInfo + ")"
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

	commitShas := slices.FilterMap(commits, func(commit *models.Commit) (string, bool) {
		return commit.Sha, commit.Sha != ""
	})

	// note that we're not filtering these as we do non-rebasing commits just because
	// I suspect that will cause some damage
	cmdObj := self.cmd.New(
		NewGitCmd("show").
			Config("log.showSignature=false").
			Arg("--no-patch", "--oneline", "--abbrev=20", prettyFormat).
			Arg(commitShas...).
			ToArgv(),
	).DontLog()

	fullCommits := map[string]*models.Commit{}
	err = cmdObj.RunAndProcessLines(func(line string) (bool, error) {
		commit := self.extractCommitFromLine(line)
		fullCommits[commit.Sha] = commit
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	findFullCommit := lo.Ternary(self.version.IsOlderThan(2, 25, 2),
		func(sha string) *models.Commit {
			for s, c := range fullCommits {
				if strings.HasPrefix(s, sha) {
					return c
				}
			}
			return nil
		},
		func(sha string) *models.Commit {
			return fullCommits[sha]
		})

	hydratedCommits := make([]*models.Commit, 0, len(commits))
	for _, rebasingCommit := range commits {
		if rebasingCommit.Sha == "" {
			hydratedCommits = append(hydratedCommits, rebasingCommit)
		} else if commit := findFullCommit(rebasingCommit.Sha); commit != nil {
			commit.Action = rebasingCommit.Action
			commit.Status = rebasingCommit.Status
			hydratedCommits = append(hydratedCommits, commit)
		}
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

	todos, err := todo.Parse(bytes.NewBuffer(bytesContent), self.config.GetCoreCommentChar())
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred while parsing git-rebase-todo file: %s", err.Error()))
		return nil, nil
	}

	// See if the current commit couldn't be applied because it conflicted; if
	// so, add a fake entry for it
	if conflictedCommitSha := self.getConflictedCommit(todos); conflictedCommitSha != "" {
		commits = append(commits, &models.Commit{
			Sha:    conflictedCommitSha,
			Name:   "",
			Status: models.StatusRebasing,
			Action: models.ActionConflict,
		})
	}

	for _, t := range todos {
		if t.Command == todo.UpdateRef {
			t.Msg = strings.TrimPrefix(t.Ref, "refs/heads/")
		} else if t.Commit == "" {
			// Command does not have a commit associated, skip
			continue
		}
		commits = slices.Prepend(commits, &models.Commit{
			Sha:    t.Commit,
			Name:   t.Msg,
			Status: models.StatusRebasing,
			Action: t.Command,
		})
	}

	return commits, nil
}

func (self *CommitLoader) getConflictedCommit(todos []todo.Todo) string {
	bytesContent, err := self.readFile(filepath.Join(self.dotGitDir, "rebase-merge/done"))
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred reading rebase-merge/done: %s", err.Error()))
		return ""
	}

	doneTodos, err := todo.Parse(bytes.NewBuffer(bytesContent), self.config.GetCoreCommentChar())
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred while parsing rebase-merge/done file: %s", err.Error()))
		return ""
	}

	amendFileExists := false
	if _, err := os.Stat(filepath.Join(self.dotGitDir, "rebase-merge/amend")); err == nil {
		amendFileExists = true
	}

	return self.getConflictedCommitImpl(todos, doneTodos, amendFileExists)
}

func (self *CommitLoader) getConflictedCommitImpl(todos []todo.Todo, doneTodos []todo.Todo, amendFileExists bool) string {
	// Should never be possible, but just to be safe:
	if len(doneTodos) == 0 {
		self.Log.Error("no done entries in rebase-merge/done file")
		return ""
	}
	lastTodo := doneTodos[len(doneTodos)-1]
	if lastTodo.Command == todo.Break || lastTodo.Command == todo.Exec || lastTodo.Command == todo.Reword {
		return ""
	}

	// In certain cases, git reschedules commands that failed. One example is if
	// a patch would overwrite an untracked file (another one is an "exec" that
	// failed, but we don't care about that here because we dealt with exec
	// already above). To detect this, compare the last command of the "done"
	// file against the first command of "git-rebase-todo"; if they are the
	// same, the command was rescheduled.
	if len(doneTodos) > 0 && len(todos) > 0 && doneTodos[len(doneTodos)-1] == todos[0] {
		// Command was rescheduled, no need to display it
		return ""
	}

	// Older versions of git have a bug whereby, if a command is rescheduled,
	// the last successful command is appended to the "done" file again. To
	// detect this, we need to compare the second-to-last done entry against the
	// first todo entry, and also compare the last done entry against the
	// last-but-two done entry; this latter check is needed for the following
	// case:
	//   pick A
	//   exec make test
	//   pick B
	//   exec make test
	// If pick B fails with conflicts, then the "done" file contains
	//   pick A
	//   exec make test
	//   pick B
	// and git-rebase-todo contains
	//   exec make test
	// Without the last condition we would erroneously treat this as the exec
	// command being rescheduled, so we wouldn't display our fake entry for
	// "pick B".
	if len(doneTodos) >= 3 && len(todos) > 0 && doneTodos[len(doneTodos)-2] == todos[0] &&
		doneTodos[len(doneTodos)-1] == doneTodos[len(doneTodos)-3] {
		// Command was rescheduled, no need to display it
		return ""
	}

	if lastTodo.Command == todo.Edit {
		if amendFileExists {
			// Special case for "edit": if the "amend" file exists, the "edit"
			// command was successful, otherwise it wasn't
			return ""
		}
	}

	// I don't think this is ever possible, but again, just to be safe:
	if lastTodo.Commit == "" {
		self.Log.Error("last command in rebase-merge/done file doesn't have a commit")
		return ""
	}

	// Any other todo that has a commit associated with it must have failed with
	// a conflict, otherwise we wouldn't have stopped the rebase:
	return lastTodo.Commit
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
		Status: models.StatusRebasing,
	}
}

func (self *CommitLoader) setCommitMergedStatuses(refName string, commits []*models.Commit) []*models.Commit {
	ancestor := self.getMergeBase(refName)
	if ancestor == "" {
		return commits
	}
	passedAncestor := false
	for i, commit := range commits {
		if strings.HasPrefix(ancestor, commit.Sha) {
			passedAncestor = true
		}
		if commit.Status != models.StatusPushed && commit.Status != models.StatusUnpushed {
			continue
		}
		if passedAncestor {
			commits[i].Status = models.StatusMerged
		}
	}
	return commits
}

func (self *CommitLoader) getMergeBase(refName string) string {
	if self.mainBranches == nil {
		self.mainBranches = self.getExistingMainBranches()
	}

	if len(self.mainBranches) == 0 {
		return ""
	}

	// We pass all configured main branches to the merge-base call; git will
	// return the base commit for the closest one.

	output, err := self.cmd.New(
		NewGitCmd("merge-base").Arg(refName).Arg(self.mainBranches...).
			ToArgv(),
	).DontLog().RunWithOutput()
	if err != nil {
		// If there's an error, it must be because one of the main branches that
		// used to exist when we called getExistingMainBranches() was deleted
		// meanwhile. To fix this for next time, throw away our cache.
		self.mainBranches = nil
	}
	return ignoringWarnings(output)
}

func (self *CommitLoader) getExistingMainBranches() []string {
	return lo.FilterMap(self.UserConfig.Git.MainBranches,
		func(branchName string, _ int) (string, bool) {
			// Try to determine upstream of local main branch
			if ref, err := self.cmd.New(
				NewGitCmd("rev-parse").Arg("--symbolic-full-name", branchName+"@{u}").ToArgv(),
			).DontLog().RunWithOutput(); err == nil {
				return strings.TrimSpace(ref), true
			}

			// If this failed, a local branch for this main branch doesn't exist or it
			// has no upstream configured. Try looking for one in the "origin" remote.
			ref := "refs/remotes/origin/" + branchName
			if err := self.cmd.New(
				NewGitCmd("rev-parse").Arg("--verify", "--quiet", ref).ToArgv(),
			).DontLog().Run(); err == nil {
				return ref, true
			}

			// If this failed as well, try if we have the main branch as a local
			// branch. This covers the case where somebody is using git locally
			// for something, but never pushing anywhere.
			ref = "refs/heads/" + branchName
			if err := self.cmd.New(
				NewGitCmd("rev-parse").Arg("--verify", "--quiet", ref).ToArgv(),
			).DontLog().Run(); err == nil {
				return ref, true
			}

			return "", false
		})
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
	output, err := self.cmd.New(
		NewGitCmd("merge-base").
			Arg(refName).
			Arg(strings.TrimPrefix(refName, "refs/heads/") + "@{u}").
			ToArgv(),
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
	config := self.UserConfig.Git.Log

	cmdArgs := NewGitCmd("log").
		Arg(opts.RefName).
		ArgIf(config.Order != "default", "--"+config.Order).
		ArgIf(opts.All, "--all").
		Arg("--oneline").
		Arg(prettyFormat).
		Arg("--abbrev=40").
		ArgIf(opts.Limit, "-300").
		ArgIf(opts.FilterPath != "", "--follow").
		Arg("--no-show-signature").
		Arg("--").
		ArgIf(opts.FilterPath != "", opts.FilterPath).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog()
}

const prettyFormat = `--pretty=format:%H%x00%at%x00%aN%x00%ae%x00%D%x00%p%x00%s`
