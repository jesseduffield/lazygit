package commands

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

func (c *Git) GetRewordCommitCmdObj(commits []*models.Commit, index int) (ICmdObj, error) {
	todo, sha, err := c.GenerateGenericRebaseTodo(commits, index, "reword")
	if err != nil {
		return nil, err
	}

	return c.InteractiveRebaseCmdObj(sha, todo, false), nil
}

func (c *Git) MoveCommitDown(commits []*models.Commit, index int) error {
	// we must ensure that we have at least two commits after the selected one
	if len(commits) <= index+2 {
		// assuming they aren't picking the bottom commit
		return errors.New(c.tr.NoRoom)
	}

	todo := ""
	orderedCommits := append(commits[0:index], commits[index+1], commits[index])
	for _, commit := range orderedCommits {
		todo = "pick " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	return c.Run(
		c.InteractiveRebaseCmdObj(commits[index+2].Sha, todo, true),
	)
}

func (c *Git) InteractiveRebase(commits []*models.Commit, index int, action string) error {
	todo, sha, err := c.GenerateGenericRebaseTodo(commits, index, action)
	if err != nil {
		return err
	}

	return c.Run(c.InteractiveRebaseCmdObj(sha, todo, true))
}

// InteractiveRebaseCmdObj returns the command object for an interactive rebase
// we tell git to run lazygit to edit the todo list, and we pass the client
// lazygit a todo string to write to the todo file
func (c *Git) InteractiveRebaseCmdObj(baseSha string, todo string, overrideEditor bool) ICmdObj {
	ex := c.GetOS().GetLazygitPath()

	debug := "FALSE"
	if c.config.GetDebug() {
		debug = "TRUE"
	}

	cmdObj := BuildGitCmdObjFromStr(fmt.Sprintf("rebase --interactive --autostash --keep-empty %s", baseSha))
	c.log.WithField("command", cmdObj.ToString()).Info("RunCommand")

	gitSequenceEditor := ex
	if todo == "" {
		gitSequenceEditor = "true"
	} else {
		c.GetOS().LogCommand(fmt.Sprintf("Creating TODO file for interactive rebase: \n\n%s", todo), false)
	}

	SetDefaultEnvVars(cmdObj)
	cmdObj.AddEnvVars(
		"LAZYGIT_CLIENT_COMMAND=INTERACTIVE_REBASE",
		"LAZYGIT_REBASE_TODO="+todo,
		"DEBUG="+debug,
		"LANG=en_US.UTF-8",   // Force using EN as language
		"LC_ALL=en_US.UTF-8", // Force using EN as language
		"GIT_SEQUENCE_EDITOR="+gitSequenceEditor,
	)

	if overrideEditor {
		cmdObj.AddEnvVars("GIT_EDITOR=" + ex)
	}

	return cmdObj
}

func (c *Git) GenerateGenericRebaseTodo(commits []*models.Commit, actionIndex int, action string) (string, string, error) {
	baseIndex := actionIndex + 1

	if len(commits) <= baseIndex {
		return "", "", errors.New(c.tr.CannotRebaseOntoFirstCommit)
	}

	if action == "squash" || action == "fixup" {
		baseIndex++

		if len(commits) <= baseIndex {
			return "", "", errors.New(c.tr.CannotSquashOntoSecondCommit)
		}
	}

	todo := ""
	for i, commit := range commits[0:baseIndex] {
		var commitAction string
		if i == actionIndex {
			commitAction = action
		} else if commit.IsMerge() {
			// your typical interactive rebase will actually drop merge commits by default. Damn git CLI, you scary!
			// doing this means we don't need to worry about rebasing over merges which always causes problems.
			// you typically shouldn't be doing rebases that pass over merge commits anyway.
			commitAction = "drop"
		} else {
			commitAction = "pick"
		}
		todo = commitAction + " " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	return todo, commits[baseIndex].Sha, nil
}

// AmendTo amends the given commit with whatever files are staged
func (c *Git) AmendTo(sha string) error {
	if err := c.Commits().CreateFixupCommit(sha); err != nil {
		return err
	}

	return c.SquashAllAboveFixupCommits(sha)
}

// EditRebaseTodo sets the action at a given index in the git-rebase-todo file
func (c *Git) EditRebaseTodo(index int, action string) error {
	fileName := filepath.Join(c.dotGitDir, "rebase-merge/git-rebase-todo")
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	content := strings.Split(string(bytes), "\n")
	commitCount := c.getTodoCommitCount(content)

	// we have the most recent commit at the bottom whereas the todo file has
	// it at the bottom, so we need to subtract our index from the commit count
	contentIndex := commitCount - 1 - index
	splitLine := strings.Split(content[contentIndex], " ")
	content[contentIndex] = action + " " + strings.Join(splitLine[1:], " ")
	result := strings.Join(content, "\n")

	return ioutil.WriteFile(fileName, []byte(result), 0644)
}

func (c *Git) getTodoCommitCount(content []string) int {
	// count lines that are not blank and are not comments
	commitCount := 0
	for _, line := range content {
		if line != "" && !strings.HasPrefix(line, "#") {
			commitCount++
		}
	}
	return commitCount
}

// MoveTodoDown moves a rebase todo item down by one position
func (c *Git) MoveTodoDown(index int) error {
	fileName := filepath.Join(c.dotGitDir, "rebase-merge/git-rebase-todo")
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	content := strings.Split(string(bytes), "\n")
	commitCount := c.getTodoCommitCount(content)
	contentIndex := commitCount - 1 - index

	rearrangedContent := append(content[0:contentIndex-1], content[contentIndex], content[contentIndex-1])
	rearrangedContent = append(rearrangedContent, content[contentIndex+1:]...)
	result := strings.Join(rearrangedContent, "\n")

	return ioutil.WriteFile(fileName, []byte(result), 0644)
}

// SquashAllAboveFixupCommits squashes all fixup! commits above the given one
func (c *Git) SquashAllAboveFixupCommits(sha string) error {
	cmdObj := oscommands.NewCmdObjFromStr(fmt.Sprintf("rebase --interactive --autostash --autosquash %s^", sha))
	c.SkipEditor(cmdObj)

	return c.Run(cmdObj)
}

// BeginInteractiveRebaseForCommit starts an interactive rebase to edit the current
// commit and pick all others. After this you'll want to call `c.ContinueRebase`
func (c *Git) BeginInteractiveRebaseForCommit(commits []*models.Commit, commitIndex int) error {
	if len(commits)-1 < commitIndex {
		return errors.New("index outside of range of commits")
	}

	// we can make this GPG thing possible it just means we need to do this in two parts:
	// one where we handle the possibility of a credential request, and the other
	// where we continue the rebase
	if c.UsingGpg() {
		return errors.New(c.tr.DisabledForGPG)
	}

	todo, sha, err := c.GenerateGenericRebaseTodo(commits, commitIndex, "edit")
	if err != nil {
		return err
	}

	if err := c.Run(c.InteractiveRebaseCmdObj(sha, todo, true)); err != nil {
		return err
	}

	return nil
}

// RebaseBranch interactive rebases onto a branch
func (c *Git) RebaseBranch(branchName string) error {
	return c.Run(c.InteractiveRebaseCmdObj(branchName, "", false))
}

func (c *Git) AbortRebase() error {
	return c.GenericMergeOrRebaseAction("rebase", "abort")
}

func (c *Git) ContinueRebase() error {
	return c.GenericMergeOrRebaseAction("rebase", "continue")
}

func (c *Git) MergeOrRebase() string {
	if c.WorkingTreeState() == "merging" {
		return "merge"
	}

	return "rebase"
}

// GenericMerge takes a commandType of "merge" or "rebase" and a command of "abort", "skip" or "continue"
// By default we skip the editor in the case where a commit will be made
func (c *Git) GenericMergeOrRebaseAction(commandType string, command string) error {
	cmdObj := oscommands.NewCmdObjFromStr(fmt.Sprintf(
		"%s --%s",
		commandType,
		command,
	))
	c.SkipEditor(cmdObj)

	err := c.Run(cmdObj)

	if err != nil {
		if !strings.Contains(err.Error(), "no rebase in progress") {
			return err
		}
		c.log.Warn(err)
	}

	// sometimes we need to do a sequence of things in a rebase but the user needs to
	// fix merge conflicts along the way. When this happens we queue up the next step
	// so that after the next successful rebase continue we can continue from where we left off
	if commandType == "rebase" && command == "continue" && c.onSuccessfulContinue != nil {
		f := c.onSuccessfulContinue
		c.onSuccessfulContinue = nil
		return f()
	}
	if command == "abort" {
		c.onSuccessfulContinue = nil
	}
	return nil
}

// CherryPickCommits begins an interactive rebase with the given shas being cherry picked onto HEAD
func (c *Git) CherryPickCommits(commits []*models.Commit) error {
	todo := ""
	for _, commit := range commits {
		todo = "pick " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	return c.Run(c.InteractiveRebaseCmdObj("HEAD", todo, false))
}
