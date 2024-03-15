package git_commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mgutz/str"
)

type BranchCommands struct {
	*GitCommon
}

func NewBranchCommands(gitCommon *GitCommon) *BranchCommands {
	return &BranchCommands{
		GitCommon: gitCommon,
	}
}

// New creates a new branch
func (self *BranchCommands) New(name string, base string) error {
	cmdArgs := NewGitCmd("checkout").
		Arg("-b", name, base).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// CreateWithUpstream creates a new branch with a given upstream, but without
// checking it out
func (self *BranchCommands) CreateWithUpstream(name string, upstream string) error {
	cmdArgs := NewGitCmd("branch").
		Arg("--track").
		Arg(name, upstream).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// CurrentBranchInfo get the current branch information.
func (self *BranchCommands) CurrentBranchInfo() (BranchInfo, error) {
	branchName, err := self.cmd.New(
		NewGitCmd("symbolic-ref").
			Arg("--short", "HEAD").
			ToArgv(),
	).DontLog().RunWithOutput()
	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return BranchInfo{
			RefName:      trimmedBranchName,
			DisplayName:  trimmedBranchName,
			DetachedHead: false,
		}, nil
	}
	output, err := self.cmd.New(
		NewGitCmd("branch").
			Arg("--points-at=HEAD", "--format=%(HEAD)%00%(objectname)%00%(refname)").
			ToArgv(),
	).DontLog().RunWithOutput()
	if err != nil {
		return BranchInfo{}, err
	}
	for _, line := range utils.SplitLines(output) {
		split := strings.Split(strings.TrimRight(line, "\r\n"), "\x00")
		if len(split) == 3 && split[0] == "*" {
			sha := split[1]
			displayName := split[2]
			return BranchInfo{
				RefName:      sha,
				DisplayName:  displayName,
				DetachedHead: true,
			}, nil
		}
	}
	return BranchInfo{
		RefName:      "HEAD",
		DisplayName:  "HEAD",
		DetachedHead: true,
	}, nil
}

// CurrentBranchName get name of current branch
func (self *BranchCommands) CurrentBranchName() (string, error) {
	cmdArgs := NewGitCmd("rev-parse").
		Arg("--abbrev-ref").
		Arg("--verify").
		Arg("HEAD").
		ToArgv()

	output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err == nil {
		return strings.TrimSpace(output), nil
	}
	return "", err
}

// LocalDelete delete branch locally
func (self *BranchCommands) LocalDelete(branch string, force bool) error {
	cmdArgs := NewGitCmd("branch").
		ArgIfElse(force, "-D", "-d").
		Arg(branch).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// Checkout checks out a branch (or commit), with --force if you set the force arg to true
type CheckoutOptions struct {
	Force   bool
	EnvVars []string
}

func (self *BranchCommands) Checkout(branch string, options CheckoutOptions) error {
	cmdArgs := NewGitCmd("checkout").
		ArgIf(options.Force, "--force").
		Arg(branch).
		ToArgv()

	return self.cmd.New(cmdArgs).
		// prevents git from prompting us for input which would freeze the program
		// TODO: see if this is actually needed here
		AddEnvVars("GIT_TERMINAL_PROMPT=0").
		AddEnvVars(options.EnvVars...).
		Run()
}

// GetGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (self *BranchCommands) GetGraph(branchName string) (string, error) {
	return self.GetGraphCmdObj(branchName).DontLog().RunWithOutput()
}

func (self *BranchCommands) GetGraphCmdObj(branchName string) oscommands.ICmdObj {
	branchLogCmdTemplate := self.UserConfig.Git.BranchLogCmd
	templateValues := map[string]string{
		"branchName": self.cmd.Quote(branchName),
	}

	resolvedTemplate := utils.ResolvePlaceholderString(branchLogCmdTemplate, templateValues)

	return self.cmd.New(str.ToArgv(resolvedTemplate)).DontLog()
}

func (self *BranchCommands) SetCurrentBranchUpstream(remoteName string, remoteBranchName string) error {
	cmdArgs := NewGitCmd("branch").
		Arg(fmt.Sprintf("--set-upstream-to=%s/%s", remoteName, remoteBranchName)).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *BranchCommands) SetUpstream(remoteName string, remoteBranchName string, branchName string) error {
	cmdArgs := NewGitCmd("branch").
		Arg(fmt.Sprintf("--set-upstream-to=%s/%s", remoteName, remoteBranchName)).
		Arg(branchName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *BranchCommands) UnsetUpstream(branchName string) error {
	cmdArgs := NewGitCmd("branch").Arg("--unset-upstream", branchName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *BranchCommands) GetCurrentBranchUpstreamDifferenceCount() (string, string) {
	return self.GetCommitDifferences("HEAD", "HEAD@{u}")
}

func (self *BranchCommands) GetUpstreamDifferenceCount(branchName string) (string, string) {
	return self.GetCommitDifferences(branchName, branchName+"@{u}")
}

// GetCommitDifferences checks how many pushables/pullables there are for the
// current branch
func (self *BranchCommands) GetCommitDifferences(from, to string) (string, string) {
	pushableCount, err := self.countDifferences(to, from)
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := self.countDifferences(from, to)
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

func (self *BranchCommands) countDifferences(from, to string) (string, error) {
	cmdArgs := NewGitCmd("rev-list").
		Arg(fmt.Sprintf("%s..%s", from, to)).
		Arg("--count").
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

func (self *BranchCommands) IsHeadDetached() bool {
	cmdArgs := NewGitCmd("symbolic-ref").Arg("-q", "HEAD").ToArgv()

	err := self.cmd.New(cmdArgs).DontLog().Run()
	return err != nil
}

func (self *BranchCommands) Rename(oldName string, newName string) error {
	cmdArgs := NewGitCmd("branch").
		Arg("--move", oldName, newName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

type MergeOpts struct {
	FastForwardOnly bool
}

func (self *BranchCommands) Merge(branchName string, opts MergeOpts) error {
	cmdArgs := NewGitCmd("merge").
		Arg("--no-edit").
		Arg(strings.Fields(self.UserConfig.Git.Merging.Args)...).
		ArgIf(opts.FastForwardOnly, "--ff-only").
		Arg(branchName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *BranchCommands) AllBranchesLogCmdObj() oscommands.ICmdObj {
	return self.cmd.New(str.ToArgv(self.UserConfig.Git.AllBranchesLogCmd)).DontLog()
}
