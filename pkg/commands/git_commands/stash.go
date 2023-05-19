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
	cmdStr := NewGitCmd("stash").Arg("drop").ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *StashCommands) Drop(index int) error {
	cmdStr := NewGitCmd("stash").Arg("drop", fmt.Sprintf("stash@{%d}", index)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *StashCommands) Pop(index int) error {
	cmdStr := NewGitCmd("stash").Arg("pop", fmt.Sprintf("stash@{%d}", index)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *StashCommands) Apply(index int) error {
	cmdStr := NewGitCmd("stash").Arg("apply", fmt.Sprintf("stash@{%d}", index)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

// Save save stash
func (self *StashCommands) Save(message string) error {
	cmdStr := NewGitCmd("stash").Arg("save", self.cmd.Quote(message)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *StashCommands) Store(sha string, message string) error {
	trimmedMessage := strings.Trim(message, " \t")

	cmdStr := NewGitCmd("stash").Arg("store", self.cmd.Quote(sha)).
		ArgIf(trimmedMessage != "", "-m", self.cmd.Quote(trimmedMessage)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *StashCommands) Sha(index int) (string, error) {
	cmdStr := NewGitCmd("rev-parse").
		Arg(fmt.Sprintf("refs/stash@{%d}", index)).
		ToString()

	sha, _, err := self.cmd.New(cmdStr).DontLog().RunWithOutputs()
	return strings.Trim(sha, "\r\n"), err
}

func (self *StashCommands) ShowStashEntryCmdObj(index int, ignoreWhitespace bool) oscommands.ICmdObj {
	cmdStr := NewGitCmd("stash").Arg("show").
		Arg("-p").
		Arg("--stat").
		Arg(fmt.Sprintf("--color=%s", self.UserConfig.Git.Paging.ColorArg)).
		Arg(fmt.Sprintf("--unified=%d", self.UserConfig.Git.DiffContextSize)).
		ArgIf(ignoreWhitespace, "--ignore-all-space").
		Arg(fmt.Sprintf("stash@{%d}", index)).
		ToString()

	return self.cmd.New(cmdStr).DontLog()
}

func (self *StashCommands) StashAndKeepIndex(message string) error {
	cmdStr := NewGitCmd("stash").Arg("save", self.cmd.Quote(message), "--keep-index").
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *StashCommands) StashUnstagedChanges(message string) error {
	if err := self.cmd.New(
		NewGitCmd("commit").
			Arg("--no-verify", "-m", self.cmd.Quote("[lazygit] stashing unstaged changes")).
			ToString(),
	).Run(); err != nil {
		return err
	}
	if err := self.Save(message); err != nil {
		return err
	}

	if err := self.cmd.New(
		NewGitCmd("reset").Arg("--soft", "HEAD^").ToString(),
	).Run(); err != nil {
		return err
	}
	return nil
}

// SaveStagedChanges stashes only the currently staged changes. This takes a few steps
// shoutouts to Joe on https://stackoverflow.com/questions/14759748/stashing-only-staged-changes-in-git-is-it-possible
func (self *StashCommands) SaveStagedChanges(message string) error {
	// wrap in 'writing', which uses a mutex
	if err := self.cmd.New(
		NewGitCmd("stash").Arg("--keep-index").ToString(),
	).Run(); err != nil {
		return err
	}

	if err := self.Save(message); err != nil {
		return err
	}

	if err := self.cmd.New(
		NewGitCmd("stash").Arg("apply", "stash@{1}").ToString(),
	).Run(); err != nil {
		return err
	}

	if err := self.os.PipeCommands(
		NewGitCmd("stash").Arg("show", "-p").ToString(),
		NewGitCmd("apply").Arg("-R").ToString(),
	); err != nil {
		return err
	}

	if err := self.cmd.New(
		NewGitCmd("stash").Arg("drop", "stash@{1}").ToString(),
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
		NewGitCmd("stash").Arg("save", self.cmd.Quote(message), "--include-untracked").
			ToString(),
	).Run()
}

func (self *StashCommands) Rename(index int, message string) error {
	sha, err := self.Sha(index)
	if err != nil {
		return err
	}

	if err := self.Drop(index); err != nil {
		return err
	}

	err = self.Store(sha, message)
	if err != nil {
		return err
	}

	return nil
}
