package git_commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
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
	getRebaseMode func() (enums.RebaseMode, error),
	gitCommon *GitCommon,
) *CommitLoader {
	return &CommitLoader{
		Common:        cmn,
		cmd:           cmd,
		getRebaseMode: getRebaseMode,
		readFile:      os.ReadFile,
		walkFiles:     filepath.Walk,
		mainBranches:  nil,
		GitCommon:     gitCommon,
	}
}

type GetCommitsOptions struct {
	Limit                bool
	FilterPath           string
	FilterAuthor         string
	IncludeRebaseCommits bool
	RefName              string // e.g. "HEAD" or "my_branch"
	RefForPushedStatus   string // the ref to use for determining pushed/unpushed status
	// determines if we show the whole git graph i.e. pass the '--all' flag
	All bool
	// If non-empty, show divergence from this ref (left-right log)
	RefToShowDivergenceFrom string
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

	wg := sync.WaitGroup{}

	wg.Add(2)

	var logErr error
	go utils.Safe(func() {
		defer wg.Done()

		logErr = self.getLogCmd(opts).RunAndProcessLines(func(line string) (bool, error) {
			commit := self.extractCommitFromLine(line, opts.RefToShowDivergenceFrom != "")
			commits = append(commits, commit)
			return false, nil
		})
	})

	var ancestor string
	var remoteAncestor string
	go utils.Safe(func() {
		defer wg.Done()

		ancestor = self.getMergeBase(opts.RefName)
		if opts.RefToShowDivergenceFrom != "" {
			remoteAncestor = self.getMergeBase(opts.RefToShowDivergenceFrom)
		}
	})

	passedFirstPushedCommit := false
	// I can get this before
	firstPushedCommit, err := self.getFirstPushedCommit(opts.RefForPushedStatus)
	if err != nil {
		// must have no upstream branch so we'll consider everything as pushed
		passedFirstPushedCommit = true
	}

	wg.Wait()

	if logErr != nil {
		return nil, logErr
	}

	for _, commit := range commits {
		if commit.Hash == firstPushedCommit {
			passedFirstPushedCommit = true
		}
		if commit.Status != models.StatusRebasing {
			if passedFirstPushedCommit {
				commit.Status = models.StatusPushed
			} else {
				commit.Status = models.StatusUnpushed
			}
		}
	}

	if len(commits) == 0 {
		return commits, nil
	}

	if opts.RefToShowDivergenceFrom != "" {
		sort.SliceStable(commits, func(i, j int) bool {
			// In the divergence view we want incoming commits to come first
			return commits[i].Divergence > commits[j].Divergence
		})

		_, localSectionStart, found := lo.FindIndexOf(commits, func(commit *models.Commit) bool {
			return commit.Divergence == models.DivergenceLeft
		})
		if !found {
			localSectionStart = len(commits)
		}

		setCommitMergedStatuses(remoteAncestor, commits[:localSectionStart])
		setCommitMergedStatuses(ancestor, commits[localSectionStart:])
	} else {
		setCommitMergedStatuses(ancestor, commits)
	}

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

// extractCommitFromLine takes a line from a git log and extracts the hash, message, date, and tag if present
// then puts them into a commit object
// example input:
// 8ad01fe32fcc20f07bc6693f87aa4977c327f1e1|10 hours ago|Jesse Duffield| (HEAD -> master, tag: v0.15.2)|refresh commits when adding a tag
func (self *CommitLoader) extractCommitFromLine(line string, showDivergence bool) *models.Commit {
	split := strings.SplitN(line, "\x00", 8)

	hash := split[0]
	unixTimestamp := split[1]
	authorName := split[2]
	authorEmail := split[3]
	extraInfo := strings.TrimSpace(split[4])
	parentHashes := split[5]
	divergence := models.DivergenceNone
	if showDivergence {
		divergence = lo.Ternary(split[6] == "<", models.DivergenceLeft, models.DivergenceRight)
	}
	message := split[7]

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
		Hash:          hash,
		Name:          message,
		Tags:          tags,
		ExtraInfo:     extraInfo,
		UnixTimestamp: int64(unitTimestampInt),
		AuthorName:    authorName,
		AuthorEmail:   authorEmail,
		Parents:       parents,
		Divergence:    divergence,
	}
}

func (self *CommitLoader) getHydratedRebasingCommits(rebaseMode enums.RebaseMode) ([]*models.Commit, error) {
	commits := self.getRebasingCommits(rebaseMode)

	if len(commits) == 0 {
		return nil, nil
	}

	commitHashes := lo.FilterMap(commits, func(commit *models.Commit, _ int) (string, bool) {
		return commit.Hash, commit.Hash != ""
	})

	// note that we're not filtering these as we do non-rebasing commits just because
	// I suspect that will cause some damage
	cmdObj := self.cmd.New(
		NewGitCmd("show").
			Config("log.showSignature=false").
			Arg("--no-patch", "--oneline", "--abbrev=20", prettyFormat).
			Arg(commitHashes...).
			ToArgv(),
	).DontLog()

	fullCommits := map[string]*models.Commit{}
	err := cmdObj.RunAndProcessLines(func(line string) (bool, error) {
		commit := self.extractCommitFromLine(line, false)
		fullCommits[commit.Hash] = commit
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	findFullCommit := lo.Ternary(self.version.IsOlderThan(2, 25, 2),
		func(hash string) *models.Commit {
			for s, c := range fullCommits {
				if strings.HasPrefix(s, hash) {
					return c
				}
			}
			return nil
		},
		func(hash string) *models.Commit {
			return fullCommits[hash]
		})

	hydratedCommits := make([]*models.Commit, 0, len(commits))
	for _, rebasingCommit := range commits {
		if rebasingCommit.Hash == "" {
			hydratedCommits = append(hydratedCommits, rebasingCommit)
		} else if commit := findFullCommit(rebasingCommit.Hash); commit != nil {
			commit.Action = rebasingCommit.Action
			commit.Status = rebasingCommit.Status
			hydratedCommits = append(hydratedCommits, commit)
		}
	}
	return hydratedCommits, nil
}

// getRebasingCommits obtains the commits that we're in the process of rebasing

// git-rebase-todo example:
// pick ac446ae94ee560bdb8d1d057278657b251aaef17 ac446ae
// pick afb893148791a2fbd8091aeb81deba4930c73031 afb8931
func (self *CommitLoader) getRebasingCommits(rebaseMode enums.RebaseMode) []*models.Commit {
	if rebaseMode != enums.REBASE_MODE_INTERACTIVE {
		return nil
	}

	bytesContent, err := self.readFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo"))
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred reading git-rebase-todo: %s", err.Error()))
		// we assume an error means the file doesn't exist so we just return
		return nil
	}

	commits := []*models.Commit{}

	todos, err := todo.Parse(bytes.NewBuffer(bytesContent), self.config.GetCoreCommentChar())
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred while parsing git-rebase-todo file: %s", err.Error()))
		return nil
	}

	// See if the current commit couldn't be applied because it conflicted; if
	// so, add a fake entry for it
	if conflictedCommitHash := self.getConflictedCommit(todos); conflictedCommitHash != "" {
		commits = append(commits, &models.Commit{
			Hash:   conflictedCommitHash,
			Name:   "",
			Status: models.StatusRebasing,
			Action: models.ActionConflict,
		})
	}

	for _, t := range todos {
		if t.Command == todo.UpdateRef {
			t.Msg = t.Ref
		} else if t.Commit == "" {
			// Command does not have a commit associated, skip
			continue
		}
		commits = utils.Prepend(commits, &models.Commit{
			Hash:   t.Commit,
			Name:   t.Msg,
			Status: models.StatusRebasing,
			Action: t.Command,
		})
	}

	return commits
}

func (self *CommitLoader) getConflictedCommit(todos []todo.Todo) string {
	bytesContent, err := self.readFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/done"))
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
	if _, err := os.Stat(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/amend")); err == nil {
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

func setCommitMergedStatuses(ancestor string, commits []*models.Commit) {
	if ancestor == "" {
		return
	}

	passedAncestor := false
	for i, commit := range commits {
		// some commits aren't really commits and don't have hashes, such as the update-ref todo
		if commit.Hash != "" && strings.HasPrefix(ancestor, commit.Hash) {
			passedAncestor = true
		}
		if commit.Status != models.StatusPushed && commit.Status != models.StatusUnpushed {
			continue
		}
		if passedAncestor {
			commits[i].Status = models.StatusMerged
		}
	}
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
	var existingBranches []string
	var wg sync.WaitGroup

	mainBranches := self.UserConfig.Git.MainBranches
	existingBranches = make([]string, len(mainBranches))

	for i, branchName := range mainBranches {
		wg.Add(1)
		go utils.Safe(func() {
			defer wg.Done()

			// Try to determine upstream of local main branch
			if ref, err := self.cmd.New(
				NewGitCmd("rev-parse").Arg("--symbolic-full-name", branchName+"@{u}").ToArgv(),
			).DontLog().RunWithOutput(); err == nil {
				existingBranches[i] = strings.TrimSpace(ref)
				return
			}

			// If this failed, a local branch for this main branch doesn't exist or it
			// has no upstream configured. Try looking for one in the "origin" remote.
			ref := "refs/remotes/origin/" + branchName
			if err := self.cmd.New(
				NewGitCmd("rev-parse").Arg("--verify", "--quiet", ref).ToArgv(),
			).DontLog().Run(); err == nil {
				existingBranches[i] = ref
				return
			}

			// If this failed as well, try if we have the main branch as a local
			// branch. This covers the case where somebody is using git locally
			// for something, but never pushing anywhere.
			ref = "refs/heads/" + branchName
			if err := self.cmd.New(
				NewGitCmd("rev-parse").Arg("--verify", "--quiet", ref).ToArgv(),
			).DontLog().Run(); err == nil {
				existingBranches[i] = ref
			}
		})
	}

	wg.Wait()

	existingBranches = lo.Filter(existingBranches, func(branch string, _ int) bool {
		return branch != ""
	})

	return existingBranches
}

func ignoringWarnings(commandOutput string) string {
	trimmedOutput := strings.TrimSpace(commandOutput)
	split := strings.Split(trimmedOutput, "\n")
	// need to get last line in case the first line is a warning about how the error is ambiguous.
	// At some point we should find a way to make it unambiguous
	lastLine := split[len(split)-1]

	return lastLine
}

// getFirstPushedCommit returns the first commit hash which has been pushed to the ref's upstream.
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
	gitLogOrder := self.AppState.GitLogOrder

	refSpec := opts.RefName
	if opts.RefToShowDivergenceFrom != "" {
		refSpec += "..." + opts.RefToShowDivergenceFrom
	}

	cmdArgs := NewGitCmd("log").
		Arg(refSpec).
		ArgIf(gitLogOrder != "default", "--"+gitLogOrder).
		ArgIf(opts.All, "--all").
		Arg("--oneline").
		Arg(prettyFormat).
		Arg("--abbrev=40").
		ArgIf(opts.FilterAuthor != "", "--author="+opts.FilterAuthor).
		ArgIf(opts.Limit, "-300").
		ArgIf(opts.FilterPath != "", "--follow").
		Arg("--no-show-signature").
		ArgIf(opts.RefToShowDivergenceFrom != "", "--left-right").
		Arg("--").
		ArgIf(opts.FilterPath != "", opts.FilterPath).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog()
}

const prettyFormat = `--pretty=format:%H%x00%at%x00%aN%x00%ae%x00%D%x00%p%x00%m%x00%s`
