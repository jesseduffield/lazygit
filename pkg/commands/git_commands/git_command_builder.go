package git_commands

import "fmt"

// convenience struct for building git commands. Especially useful when
// including conditional args
type GitCommandBuilder struct {
	command string
}

func NewGitCmd(command string) *GitCommandBuilder {
	return &GitCommandBuilder{command: command}
}

func (self *GitCommandBuilder) Arg(flag string) *GitCommandBuilder {
	if flag == "" {
		return self
	}

	self.command += " " + flag

	return self
}

func (self *GitCommandBuilder) ArgIf(include bool, flag string) *GitCommandBuilder {
	if include {
		return self.Arg(flag)
	}

	return self
}

func (self *GitCommandBuilder) ArgIfElse(isTrue bool, onTrue string, onFalse string) *GitCommandBuilder {
	if isTrue {
		return self.Arg(onTrue)
	} else {
		return self.Arg(onFalse)
	}
}

func (self *GitCommandBuilder) Args(args []string) *GitCommandBuilder {
	for _, arg := range args {
		self.Arg(arg)
	}

	return self
}

func (self *GitCommandBuilder) Config(value string) *GitCommandBuilder {
	// config settings come before the command
	self.command = fmt.Sprintf("-c %s %s", value, self.command)

	return self
}

func (self *GitCommandBuilder) RepoPath(value string) *GitCommandBuilder {
	// repo path comes before the command
	self.command = fmt.Sprintf("-C %s %s", value, self.command)

	return self
}

func (self *GitCommandBuilder) ToString() string {
	return "git " + self.command
}
