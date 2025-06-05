package git_commands

import (
	"strings"
)

// convenience struct for building git commands. Especially useful when
// including conditional args
type GitCommandBuilder struct {
	// command string
	args []string
}

func NewGitCmd(command string) *GitCommandBuilder {
	return &GitCommandBuilder{args: []string{command}}
}

func (self *GitCommandBuilder) Arg(args ...string) *GitCommandBuilder {
	self.args = append(self.args, args...)

	return self
}

func (self *GitCommandBuilder) ArgIf(condition bool, ifTrue ...string) *GitCommandBuilder {
	if condition {
		self.Arg(ifTrue...)
	}

	return self
}

func (self *GitCommandBuilder) ArgIfElse(condition bool, ifTrue string, ifFalse string) *GitCommandBuilder {
	if condition {
		return self.Arg(ifTrue)
	} else {
		return self.Arg(ifFalse)
	}
}

func (self *GitCommandBuilder) Config(value string) *GitCommandBuilder {
	// config settings come before the command
	self.args = append([]string{"-c", value}, self.args...)

	return self
}

func (self *GitCommandBuilder) ConfigIf(condition bool, ifTrue string) *GitCommandBuilder {
	if condition {
		self.Config(ifTrue)
	}

	return self
}

// the -C arg will make git do a `cd` to the directory before doing anything else
func (self *GitCommandBuilder) Dir(path string) *GitCommandBuilder {
	// repo path comes before the command
	self.args = append([]string{"-C", path}, self.args...)

	return self
}

func (self *GitCommandBuilder) DirIf(condition bool, path string) *GitCommandBuilder {
	if condition {
		return self.Dir(path)
	}

	return self
}

// Note, you may prefer to use the Dir method instead of this one
func (self *GitCommandBuilder) Worktree(path string) *GitCommandBuilder {
	// worktree arg comes before the command
	self.args = append([]string{"--work-tree", path}, self.args...)

	return self
}

func (self *GitCommandBuilder) WorktreePathIf(condition bool, path string) *GitCommandBuilder {
	if condition {
		return self.Worktree(path)
	}

	return self
}

// Note, you may prefer to use the Dir method instead of this one
func (self *GitCommandBuilder) GitDir(path string) *GitCommandBuilder {
	// git dir arg comes before the command
	self.args = append([]string{"--git-dir", path}, self.args...)

	return self
}

func (self *GitCommandBuilder) GitDirIf(condition bool, path string) *GitCommandBuilder {
	if condition {
		return self.GitDir(path)
	}

	return self
}

func (self *GitCommandBuilder) ToArgv() []string {
	return append([]string{"git"}, self.args...)
}

func (self *GitCommandBuilder) ToString() string {
	return strings.Join(self.ToArgv(), " ")
}
