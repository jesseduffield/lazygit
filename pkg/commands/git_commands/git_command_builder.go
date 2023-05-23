package git_commands

import "strings"

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

func (self *GitCommandBuilder) RepoPath(value string) *GitCommandBuilder {
	// repo path comes before the command
	self.args = append([]string{"-C", value}, self.args...)

	return self
}

func (self *GitCommandBuilder) ToArgv() []string {
	return append([]string{"git"}, self.args...)
}

func (self *GitCommandBuilder) ToString() string {
	return strings.Join(self.ToArgv(), " ")
}
