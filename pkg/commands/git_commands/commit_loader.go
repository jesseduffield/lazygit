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

	getWorkingTreeState func() models.WorkingTreeState
	readFile            func(filename string) ([]byte, error)
	walkFiles           func(root string, fn filepath.WalkFunc) error
	dotGitDir           string
	*GitCommon
}

// making our dependencies explicit for the sake of easier testing
func NewCommitLoader(
	cmn *common.Common,
	cmd oscommands.ICmdObjBuilder,
	getWorkingTreeState func() models.WorkingTreeState,
	gitCommon *GitCommon,
) *CommitLoader {
	return &CommitLoader{
		Common:              cmn,
		cmd:                 cmd,
		getWorkingTreeState: getWorkingTreeState,
		readFile:            os.ReadFile,
		walkFiles:           filepath.Walk,
		GitCommon:           gitCommon,
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
	MainBranches            *MainBranches
	HashPool                *utils.StringPool
}

// GetCommits obtains the commits of the current branch
func (self *CommitLoader) GetCommits(opts GetCommitsOptions) ([]*models.Commit, error) {
	commits := []*models.Commit{}

	if opts.IncludeRebaseCommits && opts.FilterPath == "" {
		var err error
		commits, err = self.MergeRebasingCommits(opts.HashPool, commits)
		if err != nil {
			return nil, err
		}
	}

	wg := sync.WaitGroup{}

	wg.Add(2)

	var logErr error
	go utils.Safe(func() {
		defer wg.Done()

		logErr = self.getLogCmd(opts).RunAndProcessLines(func(line string) (bool, error) {
			commit := self.extractCommitFromLine(opts.HashPool, line, opts.RefToShowDivergenceFrom != "")
			commits = append(commits, commit)
			return false, nil
		})
	})

	var ancestor string
	var remoteAncestor string
	go utils.Safe(func() {
		defer wg.Done()

		ancestor = opts.MainBranches.GetMergeBase(opts.RefName)
		if opts.RefToShowDivergenceFrom != "" {
			remoteAncestor = opts.MainBranches.GetMergeBase(opts.RefToShowDivergenceFrom)
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
		if commit.Hash() == firstPushedCommit {
			passedFirstPushedCommit = true
		}
		if !commit.IsTODO() {
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

func (self *CommitLoader) MergeRebasingCommits(hashPool *utils.StringPool, commits []*models.Commit) ([]*models.Commit, error) {
	// chances are we have as many commits as last time so we'll set the capacity to be the old length
	result := make([]*models.Commit, 0, len(commits))
	for i, commit := range commits {
		if !commit.IsTODO() { // removing the existing rebase commits so we can add the refreshed ones
			result = append(result, commits[i:]...)
			break
		}
	}

	workingTreeState := self.getWorkingTreeState()
	addConflictedRebasingCommit := true
	if workingTreeState.CherryPicking || workingTreeState.Reverting {
		sequencerCommits, err := self.getHydratedSequencerCommits(hashPool, workingTreeState)
		if err != nil {
			return nil, err
		}
		result = append(sequencerCommits, result...)
		addConflictedRebasingCommit = false
	}

	if workingTreeState.Rebasing {
		rebasingCommits, err := self.getHydratedRebasingCommits(hashPool, addConflictedRebasingCommit)
		if err != nil {
			return nil, err
		}
		if len(rebasingCommits) > 0 {
			result = append(rebasingCommits, result...)
		}
	}
	return result, nil
}

// extractCommitFromLine takes a line from a git log and extracts the hash, message, date, and tag if present
// then puts them into a commit object
// example input:
// 8ad01fe32fcc20f07bc6693f87aa4977c327f1e1|10 hours ago|Jesse Duffield| (HEAD -> master, tag: v0.15.2)|refresh commits when adding a tag
func (self *CommitLoader) extractCommitFromLine(hashPool *utils.StringPool, line string, showDivergence bool) *models.Commit {
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

	return models.NewCommit(hashPool, models.NewCommitOpts{
		Hash:          hash,
		Name:          message,
		Tags:          tags,
		ExtraInfo:     extraInfo,
		UnixTimestamp: int64(unitTimestampInt),
		AuthorName:    authorName,
		AuthorEmail:   authorEmail,
		Parents:       parents,
		Divergence:    divergence,
	})
}

func (self *CommitLoader) getHydratedRebasingCommits(hashPool *utils.StringPool, addConflictingCommit bool) ([]*models.Commit, error) {
	todoFileHasShortHashes := self.version.IsOlderThan(2, 25, 2)
	return self.getHydratedTodoCommits(hashPool, self.getRebasingCommits(hashPool, addConflictingCommit), todoFileHasShortHashes)
}

func (self *CommitLoader) getHydratedSequencerCommits(hashPool *utils.StringPool, workingTreeState models.WorkingTreeState) ([]*models.Commit, error) {
	commits := self.getSequencerCommits(hashPool)
	if len(commits) > 0 {
		// If we have any commits in .git/sequencer/todo, then the last one of
		// those is the conflicting one.
		commits[len(commits)-1].Status = models.StatusConflicted
	} else {
		// For single-commit cherry-picks and reverts, git apparently doesn't
		// use the sequencer; in that case, CHERRY_PICK_HEAD or REVERT_HEAD is
		// our conflicting commit, so synthesize it here.
		conflicedCommit := self.getConflictedSequencerCommit(hashPool, workingTreeState)
		if conflicedCommit != nil {
			commits = append(commits, conflicedCommit)
		}
	}

	return self.getHydratedTodoCommits(hashPool, commits, true)
}

func (self *CommitLoader) getHydratedTodoCommits(hashPool *utils.StringPool, todoCommits []*models.Commit, todoFileHasShortHashes bool) ([]*models.Commit, error) {
	if len(todoCommits) == 0 {
		return nil, nil
	}

	commitHashes := lo.FilterMap(todoCommits, func(commit *models.Commit, _ int) (string, bool) {
		return commit.Hash(), commit.Hash() != ""
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
		commit := self.extractCommitFromLine(hashPool, line, false)
		fullCommits[commit.Hash()] = commit
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	findFullCommit := lo.Ternary(todoFileHasShortHashes,
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

	hydratedCommits := make([]*models.Commit, 0, len(todoCommits))
	for _, rebasingCommit := range todoCommits {
		if rebasingCommit.Hash() == "" {
			hydratedCommits = append(hydratedCommits, rebasingCommit)
		} else if commit := findFullCommit(rebasingCommit.Hash()); commit != nil {
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
func (self *CommitLoader) getRebasingCommits(hashPool *utils.StringPool, addConflictingCommit bool) []*models.Commit {
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
	if addConflictingCommit {
		if conflictedCommit := self.getConflictedCommit(hashPool, todos); conflictedCommit != nil {
			commits = append(commits, conflictedCommit)
		}
	}

	for _, t := range todos {
		if t.Command == todo.UpdateRef {
			t.Msg = t.Ref
		} else if t.Command == todo.Exec {
			t.Msg = t.ExecCommand
		} else if t.Commit == "" {
			// Command does not have a commit associated, skip
			continue
		}
		commits = utils.Prepend(commits, models.NewCommit(hashPool, models.NewCommitOpts{
			Hash:   t.Commit,
			Name:   t.Msg,
			Status: models.StatusRebasing,
			Action: t.Command,
		}))
	}

	return commits
}

func (self *CommitLoader) getConflictedCommit(hashPool *utils.StringPool, todos []todo.Todo) *models.Commit {
	bytesContent, err := self.readFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/done"))
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred reading rebase-merge/done: %s", err.Error()))
		return nil
	}

	doneTodos, err := todo.Parse(bytes.NewBuffer(bytesContent), self.config.GetCoreCommentChar())
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred while parsing rebase-merge/done file: %s", err.Error()))
		return nil
	}

	amendFileExists, _ := self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/amend"))
	messageFileExists, _ := self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/message"))

	return self.getConflictedCommitImpl(hashPool, todos, doneTodos, amendFileExists, messageFileExists)
}

func (self *CommitLoader) getConflictedCommitImpl(hashPool *utils.StringPool, todos []todo.Todo, doneTodos []todo.Todo, amendFileExists bool, messageFileExists bool) *models.Commit {
	// Should never be possible, but just to be safe:
	if len(doneTodos) == 0 {
		self.Log.Error("no done entries in rebase-merge/done file")
		return nil
	}
	lastTodo := doneTodos[len(doneTodos)-1]
	if lastTodo.Command == todo.Break || lastTodo.Command == todo.Exec || lastTodo.Command == todo.Reword {
		return nil
	}

	// In certain cases, git reschedules commands that failed. One example is if
	// a patch would overwrite an untracked file (another one is an "exec" that
	// failed, but we don't care about that here because we dealt with exec
	// already above). To detect this, compare the last command of the "done"
	// file against the first command of "git-rebase-todo"; if they are the
	// same, the command was rescheduled.
	if len(doneTodos) > 0 && len(todos) > 0 && doneTodos[len(doneTodos)-1] == todos[0] {
		// Command was rescheduled, no need to display it
		return nil
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
		return nil
	}

	if lastTodo.Command == todo.Edit {
		if amendFileExists {
			// Special case for "edit": if the "amend" file exists, the "edit"
			// command was successful, otherwise it wasn't
			return nil
		}

		if !messageFileExists {
			// As an additional check, see if the "message" file exists; if it
			// doesn't, it must be because a multi-commit cherry-pick or revert
			// was performed in the meantime, which deleted both the amend file
			// and the message file.
			return nil
		}
	}

	// I don't think this is ever possible, but again, just to be safe:
	if lastTodo.Commit == "" {
		self.Log.Error("last command in rebase-merge/done file doesn't have a commit")
		return nil
	}

	// Any other todo that has a commit associated with it must have failed with
	// a conflict, otherwise we wouldn't have stopped the rebase:
	return models.NewCommit(hashPool, models.NewCommitOpts{
		Hash:   lastTodo.Commit,
		Action: lastTodo.Command,
		Status: models.StatusConflicted,
	})
}

func (self *CommitLoader) getSequencerCommits(hashPool *utils.StringPool) []*models.Commit {
	bytesContent, err := self.readFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "sequencer/todo"))
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred reading sequencer/todo: %s", err.Error()))
		// we assume an error means the file doesn't exist so we just return
		return nil
	}

	commits := []*models.Commit{}

	todos, err := todo.Parse(bytes.NewBuffer(bytesContent), self.config.GetCoreCommentChar())
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred while parsing sequencer/todo file: %s", err.Error()))
		return nil
	}

	for _, t := range todos {
		if t.Commit == "" {
			// Command does not have a commit associated, skip
			continue
		}
		commits = utils.Prepend(commits, models.NewCommit(hashPool, models.NewCommitOpts{
			Hash:   t.Commit,
			Name:   t.Msg,
			Status: models.StatusCherryPickingOrReverting,
			Action: t.Command,
		}))
	}

	return commits
}

func (self *CommitLoader) getConflictedSequencerCommit(hashPool *utils.StringPool, workingTreeState models.WorkingTreeState) *models.Commit {
	var shaFile string
	var action todo.TodoCommand
	if workingTreeState.CherryPicking {
		shaFile = "CHERRY_PICK_HEAD"
		action = todo.Pick
	} else if workingTreeState.Reverting {
		shaFile = "REVERT_HEAD"
		action = todo.Revert
	} else {
		return nil
	}
	bytesContent, err := self.readFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), shaFile))
	if err != nil {
		self.Log.Error(fmt.Sprintf("error occurred reading %s: %s", shaFile, err.Error()))
		// we assume an error means the file doesn't exist so we just return
		return nil
	}
	lines := strings.Split(string(bytesContent), "\n")
	if len(lines) == 0 {
		return nil
	}
	return models.NewCommit(hashPool, models.NewCommitOpts{
		Hash:   lines[0],
		Status: models.StatusConflicted,
		Action: action,
	})
}

func setCommitMergedStatuses(ancestor string, commits []*models.Commit) {
	if ancestor == "" {
		return
	}

	passedAncestor := false
	for i, commit := range commits {
		// some commits aren't really commits and don't have hashes, such as the update-ref todo
		if commit.Hash() != "" && strings.HasPrefix(ancestor, commit.Hash()) {
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

const prettyFormat = `--pretty=format:%H%x00%at%x00%aN%x00%ae%x00%D%x00%P%x00%m%x00%s`
