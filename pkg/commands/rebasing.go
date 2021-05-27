package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/mgutz/str"
)

func (c *GitCommand) RewordCommit(commits []*models.Commit, index int) (*exec.Cmd, error) {
	todo, sha, err := c.GenerateGenericRebaseTodo(commits, index, "reword")
	if err != nil {
		return nil, err
	}

	return c.PrepareInteractiveRebaseCommand(sha, todo, false)
}

func (c *GitCommand) MoveCommitDown(commits []*models.Commit, index int) error {
	// we must ensure that we have at least two commits after the selected one
	if len(commits) <= index+2 {
		// assuming they aren't picking the bottom commit
		return errors.New(c.Tr.NoRoom)
	}

	todo := ""
	orderedCommits := append(commits[0:index], commits[index+1], commits[index])
	for _, commit := range orderedCommits {
		todo = "pick " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	cmd, err := c.PrepareInteractiveRebaseCommand(commits[index+2].Sha, todo, true)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

func (c *GitCommand) InteractiveRebase(commits []*models.Commit, index int, action string) error {
	todo, sha, err := c.GenerateGenericRebaseTodo(commits, index, action)
	if err != nil {
		return err
	}

	cmd, err := c.PrepareInteractiveRebaseCommand(sha, todo, true)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

// PrepareInteractiveRebaseCommand returns the cmd for an interactive rebase
// we tell git to run lazygit to edit the todo list, and we pass the client
// lazygit a todo string to write to the todo file
func (c *GitCommand) PrepareInteractiveRebaseCommand(baseSha string, todo string, overrideEditor bool) (*exec.Cmd, error) {
	ex := c.OSCommand.GetLazygitPath()

	debug := "FALSE"
	if c.OSCommand.Config.GetDebug() {
		debug = "TRUE"
	}

	cmdStr := fmt.Sprintf("git rebase --interactive --autostash --keep-empty %s", baseSha)
	c.Log.WithField("command", cmdStr).Info("RunCommand")
	splitCmd := str.ToArgv(cmdStr)

	cmd := c.OSCommand.Command(splitCmd[0], splitCmd[1:]...)

	gitSequenceEditor := ex
	if todo == "" {
		gitSequenceEditor = "true"
	} else {
		c.OSCommand.LogCommand(fmt.Sprintf("Creating TODO file for interactive rebase: \n\n%s", todo), false)
	}

	cmd.Env = os.Environ()
	cmd.Env = append(
		cmd.Env,
		"LAZYGIT_CLIENT_COMMAND=INTERACTIVE_REBASE",
		"LAZYGIT_REBASE_TODO="+todo,
		"DEBUG="+debug,
		"LANG=en_US.UTF-8",   // Force using EN as language
		"LC_ALL=en_US.UTF-8", // Force using EN as language
		"GIT_SEQUENCE_EDITOR="+gitSequenceEditor,
	)

	if overrideEditor {
		cmd.Env = append(cmd.Env, "GIT_EDITOR="+ex)
	}

	return cmd, nil
}

func (c *GitCommand) GenerateGenericRebaseTodo(commits []*models.Commit, actionIndex int, action string) (string, string, error) {
	baseIndex := actionIndex + 1

	if len(commits) <= baseIndex {
		return "", "", errors.New(c.Tr.CannotRebaseOntoFirstCommit)
	}

	if action == "squash" || action == "fixup" {
		baseIndex++

		if len(commits) <= baseIndex {
			return "", "", errors.New(c.Tr.CannotSquashOntoSecondCommit)
		}
	}

	todo := ""
	for i, commit := range commits[0:baseIndex] {
		var commitAction string
		if i == actionIndex {
			commitAction = action
		} else if commit.IsMerge {
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
func (c *GitCommand) AmendTo(sha string) error {
	if err := c.CreateFixupCommit(sha); err != nil {
		return err
	}

	return c.SquashAllAboveFixupCommits(sha)
}

// EditRebaseTodo sets the action at a given index in the git-rebase-todo file
func (c *GitCommand) EditRebaseTodo(index int, action string) error {
	fileName := filepath.Join(c.DotGitDir, "rebase-merge/git-rebase-todo")
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

func (c *GitCommand) getTodoCommitCount(content []string) int {
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
func (c *GitCommand) MoveTodoDown(index int) error {
	fileName := filepath.Join(c.DotGitDir, "rebase-merge/git-rebase-todo")
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
func (c *GitCommand) SquashAllAboveFixupCommits(sha string) error {
	return c.runSkipEditorCommand(
		fmt.Sprintf(
			"git rebase --interactive --autostash --autosquash %s^",
			sha,
		),
	)
}

// BeginInteractiveRebaseForCommit starts an interactive rebase to edit the current
// commit and pick all others. After this you'll want to call `c.GenericMergeOrRebaseAction("rebase", "continue")`
func (c *GitCommand) BeginInteractiveRebaseForCommit(commits []*models.Commit, commitIndex int) error {
	if len(commits)-1 < commitIndex {
		return errors.New("index outside of range of commits")
	}

	// we can make this GPG thing possible it just means we need to do this in two parts:
	// one where we handle the possibility of a credential request, and the other
	// where we continue the rebase
	if c.UsingGpg() {
		return errors.New(c.Tr.DisabledForGPG)
	}

	todo, sha, err := c.GenerateGenericRebaseTodo(commits, commitIndex, "edit")
	if err != nil {
		return err
	}

	cmd, err := c.PrepareInteractiveRebaseCommand(sha, todo, true)
	if err != nil {
		return err
	}

	if err := c.OSCommand.RunPreparedCommand(cmd); err != nil {
		return err
	}

	return nil
}

// RebaseBranch interactive rebases onto a branch
func (c *GitCommand) RebaseBranch(branchName string) error {
	cmd, err := c.PrepareInteractiveRebaseCommand(branchName, "", false)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

// GenericMerge takes a commandType of "merge" or "rebase" and a command of "abort", "skip" or "continue"
// By default we skip the editor in the case where a commit will be made
func (c *GitCommand) GenericMergeOrRebaseAction(commandType string, command string) error {
	err := c.runSkipEditorCommand(
		fmt.Sprintf(
			"git %s --%s",
			commandType,
			command,
		),
	)
	if err != nil {
		if !strings.Contains(err.Error(), "no rebase in progress") {
			return err
		}
		c.Log.Warn(err)
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

func (c *GitCommand) runSkipEditorCommand(command string) error {
	cmd := c.OSCommand.ExecutableFromString(command)
	lazyGitPath := c.OSCommand.GetLazygitPath()
	cmd.Env = append(
		cmd.Env,
		"LAZYGIT_CLIENT_COMMAND=EXIT_IMMEDIATELY",
		"GIT_EDITOR="+lazyGitPath,
		"EDITOR="+lazyGitPath,
		"VISUAL="+lazyGitPath,
	)
	return c.OSCommand.RunExecutable(cmd)
}
