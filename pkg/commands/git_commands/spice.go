package git_commands

import (
	"os/exec"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type SpiceCommands struct {
	*GitCommon
}

func NewSpiceCommands(gitCommon *GitCommon) *SpiceCommands {
	return &SpiceCommands{
		GitCommon: gitCommon,
	}
}

// IsAvailable checks if gs binary exists in PATH
func (self *SpiceCommands) IsAvailable() bool {
	_, err := exec.LookPath("gs")
	return err == nil
}

// IsInitialized checks if repo has been initialized with git-spice
func (self *SpiceCommands) IsInitialized() bool {
	if !self.IsAvailable() {
		return false
	}
	cmdArgs := []string{"gs", "repo", "status"}
	return self.cmd.New(cmdArgs).Run() == nil
}

// GetStackBranches runs gs log long --json -a and returns the raw output for parsing
func (self *SpiceCommands) GetStackBranches() (string, error) {
	cmdArgs := []string{"gs", "log", "long", "--json", "-a"}
	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

// Restack restacks a branch or entire stack
func (self *SpiceCommands) Restack(branchName string) error {
	if branchName == "" {
		cmdArgs := []string{"gs", "stack", "restack"}
		return self.cmd.New(cmdArgs).Run()
	}
	cmdArgs := []string{"gs", "branch", "restack", "--branch", branchName}
	return self.cmd.New(cmdArgs).Run()
}

// Submit submits PR for branch or stack
func (self *SpiceCommands) Submit(branchName string) error {
	if branchName == "" {
		cmdArgs := []string{"gs", "stack", "submit"}
		return self.cmd.New(cmdArgs).Run()
	}
	cmdArgs := []string{"gs", "branch", "submit", "--branch", branchName}
	return self.cmd.New(cmdArgs).Run()
}

// Navigation commands
func (self *SpiceCommands) NavigateUp() error {
	cmdArgs := []string{"gs", "up"}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) NavigateDown() error {
	cmdArgs := []string{"gs", "down"}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) NavigateTop() error {
	cmdArgs := []string{"gs", "top"}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) NavigateBottom() error {
	cmdArgs := []string{"gs", "bottom"}
	return self.cmd.New(cmdArgs).Run()
}

// Track/untrack branches
func (self *SpiceCommands) TrackBranch(branchName string) error {
	cmdArgs := []string{"gs", "branch", "track", branchName}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) UntrackBranch(branchName string) error {
	cmdArgs := []string{"gs", "branch", "untrack", branchName}
	return self.cmd.New(cmdArgs).Run()
}
