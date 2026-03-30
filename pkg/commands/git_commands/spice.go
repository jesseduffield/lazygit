package git_commands

import (
	"os/exec"
)

type SpiceCommands struct {
	*GitCommon
	initializedCache *bool // Cache for IsInitialized check
}

func NewSpiceCommands(gitCommon *GitCommon) *SpiceCommands {
	return &SpiceCommands{
		GitCommon:        gitCommon,
		initializedCache: nil,
	}
}

func (self *SpiceCommands) IsAvailable() bool {
	_, err := exec.LookPath("gs")
	return err == nil
}

// IsInitialized checks if repo has been initialized with git-spice
// Result is cached to avoid repeated command execution
func (self *SpiceCommands) IsInitialized() bool {
	if self.initializedCache != nil {
		return *self.initializedCache
	}

	if !self.IsAvailable() {
		result := false
		self.initializedCache = &result
		return false
	}
	// Try running a simple log command - will succeed if initialized, fail otherwise
	cmdArgs := []string{"gs", "log", "short"}
	result := self.cmd.New(cmdArgs).DontLog().Run() == nil
	self.initializedCache = &result
	return result
}

func (self *SpiceCommands) GetStackBranches(format string) (string, error) {
	if format != "short" && format != "long" {
		format = "short" // Fallback if invalid
	}
	cmdArgs := []string{"gs", "log", format, "--json", "-a", "-S"}
	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

func (self *SpiceCommands) Restack(branchName string) error {
	if branchName == "" {
		cmdArgs := []string{"gs", "stack", "restack"}
		return self.cmd.New(cmdArgs).Run()
	}
	cmdArgs := []string{"gs", "branch", "restack", "--branch", branchName}
	return self.cmd.New(cmdArgs).Run()
}

// SubmitOpts contains options for the Submit command
type SubmitOpts struct {
	NoPublish  bool
	UpdateOnly bool
}

func (self *SpiceCommands) Submit(branchName string, opts SubmitOpts) error {
	var cmdArgs []string
	if branchName == "" {
		cmdArgs = []string{"gs", "stack", "submit"}
	} else {
		cmdArgs = []string{"gs", "branch", "submit", "--branch", branchName}
	}
	if opts.NoPublish {
		cmdArgs = append(cmdArgs, "--no-publish")
	}
	if opts.UpdateOnly {
		cmdArgs = append(cmdArgs, "--update-only")
	}
	return self.cmd.New(cmdArgs).Run()
}

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

func (self *SpiceCommands) CreateBranch(branchName string, targetBranch string) error {
	cmdArgs := []string{"gs", "branch", "create"}
	if targetBranch != "" {
		cmdArgs = append(cmdArgs, "--target", targetBranch)
	}
	cmdArgs = append(cmdArgs, branchName)
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) CreateCommit() error {
	cmdArgs := []string{"gs", "commit", "create"}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) DeleteBranch(branchName string) error {
	cmdArgs := []string{"gs", "branch", "delete", branchName}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) MoveBranchUp(branchName string) error {
	cmdArgs := []string{"gs", "branch", "up", "--branch", branchName}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) MoveBranchDown(branchName string) error {
	cmdArgs := []string{"gs", "branch", "down", "--branch", branchName}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) Init(trunk string) error {
	cmdArgs := []string{"gs", "repo", "init", "--trunk", trunk}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) Sync(restack bool) error {
	cmdArgs := []string{"gs", "repo", "sync"}
	if restack {
		cmdArgs = append(cmdArgs, "--restack")
	}
	return self.cmd.New(cmdArgs).Run()
}

func (self *SpiceCommands) ClearInitializedCache() {
	self.initializedCache = nil
}
