package git_commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type StashCommands struct {
	*GitCommon
	fileLoader  *FileLoader
	workingTree *WorkingTreeCommands
}

func NewStashCommands(
	gitCommon *GitCommon,
	fileLoader *FileLoader,
	workingTree *WorkingTreeCommands,
) *StashCommands {
	return &StashCommands{
		GitCommon:   gitCommon,
		fileLoader:  fileLoader,
		workingTree: workingTree,
	}
}

func (self *StashCommands) DropNewest() error {
	cmdArgs := NewGitCmd("stash").Arg("drop").ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *StashCommands) Drop(index int) error {
	cmdArgs := NewGitCmd("stash").Arg("drop", fmt.Sprintf("stash@{%d}", index)).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *StashCommands) Pop(index int) error {
	cmdArgs := NewGitCmd("stash").Arg("pop", fmt.Sprintf("stash@{%d}", index)).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *StashCommands) Apply(index int) error {
	cmdArgs := NewGitCmd("stash").Arg("apply", fmt.Sprintf("stash@{%d}", index)).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// Push push stash
func (self *StashCommands) Push(message string) error {
	cmdArgs := NewGitCmd("stash").Arg("push", "-m", message).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *StashCommands) Store(hash string, message string) error {
	trimmedMessage := strings.Trim(message, " \t")

	cmdArgs := NewGitCmd("stash").Arg("store").
		ArgIf(trimmedMessage != "", "-m", trimmedMessage).
		Arg(hash).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *StashCommands) Hash(index int) (string, error) {
	cmdArgs := NewGitCmd("rev-parse").
		Arg(fmt.Sprintf("refs/stash@{%d}", index)).
		ToArgv()

	hash, _, err := self.cmd.New(cmdArgs).DontLog().RunWithOutputs()
	return strings.Trim(hash, "\r\n"), err
}

func (self *StashCommands) ShowStashEntryCmdObj(index int) oscommands.ICmdObj {
	cmdArgs := NewGitCmd("stash").Arg("show").
		Arg("-p").
		Arg("--stat").
		Arg(fmt.Sprintf("--color=%s", self.UserConfig.Git.Paging.ColorArg)).
		Arg(fmt.Sprintf("--unified=%d", self.AppState.DiffContextSize)).
		ArgIf(self.AppState.IgnoreWhitespaceInDiffView, "--ignore-all-space").
		Arg(fmt.Sprintf("--find-renames=%d%%", self.AppState.RenameSimilarityThreshold)).
		Arg(fmt.Sprintf("stash@{%d}", index)).
		Dir(self.repoPaths.worktreePath).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog()
}

func (self *StashCommands) StashAndKeepIndex(message string) error {
	cmdArgs := NewGitCmd("stash").Arg("push", "--keep-index", "-m", message).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *StashCommands) StashUnstagedChanges(message string) error {
	if err := self.cmd.New(
		NewGitCmd("commit").
			Arg("--no-verify", "-m", "[lazygit] stashing unstaged changes").
			ToArgv(),
	).Run(); err != nil {
		return err
	}
	if err := self.Push(message); err != nil {
		return err
	}

	if err := self.cmd.New(
		NewGitCmd("reset").Arg("--soft", "HEAD^").ToArgv(),
	).Run(); err != nil {
		return err
	}
	return nil
}

// SaveStagedChanges stashes only the currently staged changes.
func (self *StashCommands) SaveStagedChanges(message string) error {
	if self.version.IsAtLeast(2, 35, 0) {
		return self.cmd.New(NewGitCmd("stash").Arg("push").Arg("--staged").Arg("-m", message).ToArgv()).Run()
	}

	// Git versions older than 2.35.0 don't support the --staged flag, so we
	// need to fall back to a more complex solution.
	// Shoutouts to Joe on https://stackoverflow.com/questions/14759748/stashing-only-staged-changes-in-git-is-it-possible
	//
	// Note that this method has a few bugs:
	// - it fails when there are *only* staged changes
	// - it fails when staged and unstaged changes within a single file are too close together
	// We don't bother fixing these, because users can simply update git when
	// they are affected by these issues.

	// wrap in 'writing', which uses a mutex
	if err := self.cmd.New(
		NewGitCmd("stash").Arg("--keep-index").ToArgv(),
	).Run(); err != nil {
		return err
	}

	if err := self.Push(message); err != nil {
		return err
	}

	if err := self.cmd.New(
		NewGitCmd("stash").Arg("apply", "stash@{1}").ToArgv(),
	).Run(); err != nil {
		return err
	}

	if err := self.os.PipeCommands(
		self.cmd.New(NewGitCmd("stash").Arg("show", "-p").ToArgv()),
		self.cmd.New(NewGitCmd("apply").Arg("-R").ToArgv()),
	); err != nil {
		return err
	}

	if err := self.cmd.New(
		NewGitCmd("stash").Arg("drop", "stash@{1}").ToArgv(),
	).Run(); err != nil {
		return err
	}

	// if you had staged an untracked file, that will now appear as 'AD' in git status
	// meaning it's deleted in your working tree but added in your index. Given that it's
	// now safely stashed, we need to remove it.
	files := self.fileLoader.
		GetStatusFiles(GetStatusFileOptions{})

	for _, file := range files {
		if file.ShortStatus == "AD" {
			if err := self.workingTree.UnStageFile(file.Names(), false); err != nil {
				return err
			}
		}
	}

	return nil
}

func (self *StashCommands) StashIncludeUntrackedChanges(message string) error {
	return self.cmd.New(
		NewGitCmd("stash").Arg("push", "--include-untracked", "-m", message).
			ToArgv(),
	).Run()
}

func (self *StashCommands) Rename(index int, message string) error {
	hash, err := self.Hash(index)
	if err != nil {
		return err
	}

	if err := self.Drop(index); err != nil {
		return err
	}

	err = self.Store(hash, message)
	if err != nil {
		return err
	}

	return nil
}
