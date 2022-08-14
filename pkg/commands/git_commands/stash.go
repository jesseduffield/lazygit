package git_commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type StashCommands struct {
	*GitCommon
	fileLoader  *loaders.FileLoader
	workingTree *WorkingTreeCommands
}

func NewStashCommands(
	gitCommon *GitCommon,
	fileLoader *loaders.FileLoader,
	workingTree *WorkingTreeCommands,
) *StashCommands {
	return &StashCommands{
		GitCommon:   gitCommon,
		fileLoader:  fileLoader,
		workingTree: workingTree,
	}
}

func (self *StashCommands) DropNewest() error {
	return self.cmd.New("git stash drop").Run()
}

func (self *StashCommands) Drop(index int) error {
	return self.cmd.New(fmt.Sprintf("git stash drop stash@{%d}", index)).Run()
}

func (self *StashCommands) Pop(index int) error {
	return self.cmd.New(fmt.Sprintf("git stash pop stash@{%d}", index)).Run()
}

func (self *StashCommands) Apply(index int) error {
	return self.cmd.New(fmt.Sprintf("git stash apply stash@{%d}", index)).Run()
}

// Save save stash
func (self *StashCommands) Save(message string) error {
	return self.cmd.New("git stash save " + self.cmd.Quote(message)).Run()
}

func (self *StashCommands) ShowStashEntryCmdObj(index int) oscommands.ICmdObj {
	cmdStr := fmt.Sprintf("git stash show -p --stat --color=%s --unified=%d stash@{%d}", self.UserConfig.Git.Paging.ColorArg, self.UserConfig.Git.DiffContextSize, index)

	return self.cmd.New(cmdStr).DontLog()
}

func (self *StashCommands) StashAndKeepIndex(message string) error {
	return self.cmd.New(fmt.Sprintf("git stash save %s --keep-index", self.cmd.Quote(message))).Run()
}

func (self *StashCommands) StashUnstagedChanges(message string) error {
	if err := self.cmd.New("git commit --no-verify -m \"[lazygit] stashing unstaged changes\"").Run(); err != nil {
		return err
	}
	if err := self.Save(message); err != nil {
		return err
	}
	if err := self.cmd.New("git reset --soft HEAD^").Run(); err != nil {
		return err
	}
	return nil
}

// SaveStagedChanges stashes only the currently staged changes. This takes a few steps
// shoutouts to Joe on https://stackoverflow.com/questions/14759748/stashing-only-staged-changes-in-git-is-it-possible
func (self *StashCommands) SaveStagedChanges(message string) error {
	// wrap in 'writing', which uses a mutex
	if err := self.cmd.New("git stash --keep-index").Run(); err != nil {
		return err
	}

	if err := self.Save(message); err != nil {
		return err
	}

	if err := self.cmd.New("git stash apply stash@{1}").Run(); err != nil {
		return err
	}

	if err := self.os.PipeCommands("git stash show -p", "git apply -R"); err != nil {
		return err
	}

	if err := self.cmd.New("git stash drop stash@{1}").Run(); err != nil {
		return err
	}

	// if you had staged an untracked file, that will now appear as 'AD' in git status
	// meaning it's deleted in your working tree but added in your index. Given that it's
	// now safely stashed, we need to remove it.
	files := self.fileLoader.
		GetStatusFiles(loaders.GetStatusFileOptions{})

	for _, file := range files {
		if file.ShortStatus == "AD" {
			if err := self.workingTree.UnStageFile(file.Names(), false); err != nil {
				return err
			}
		}
	}

	return nil
}
