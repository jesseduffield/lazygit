package git_commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
	cmdStr := NewGitCmd("checkout").
		Arg("-b", self.cmd.Quote(name), self.cmd.Quote(base)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

// CurrentBranchInfo get the current branch information.
func (self *BranchCommands) CurrentBranchInfo() (BranchInfo, error) {
	branchName, err := self.cmd.New(
		NewGitCmd("symbolic-ref").
			Arg("--short", "HEAD").
			ToString(),
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
			Arg("--points-at=HEAD", "--format=\"%(HEAD)%00%(objectname)%00%(refname)\"").
			ToString(),
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

// Delete delete branch
func (self *BranchCommands) Delete(branch string, force bool) error {
	cmdStr := NewGitCmd("branch").
		ArgIfElse(force, "-D", "-d").
		Arg(self.cmd.Quote(branch)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

// Checkout checks out a branch (or commit), with --force if you set the force arg to true
type CheckoutOptions struct {
	Force   bool
	EnvVars []string
}

func (self *BranchCommands) Checkout(branch string, options CheckoutOptions) error {
	cmdStr := NewGitCmd("checkout").
		ArgIf(options.Force, "--force").
		Arg(self.cmd.Quote(branch)).
		ToString()

	return self.cmd.New(cmdStr).
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
	return self.cmd.New(utils.ResolvePlaceholderString(branchLogCmdTemplate, templateValues)).DontLog()
}

func (self *BranchCommands) SetCurrentBranchUpstream(remoteName string, remoteBranchName string) error {
	cmdStr := NewGitCmd("branch").
		Arg(fmt.Sprintf("--set-upstream-to=%s/%s", self.cmd.Quote(remoteName), self.cmd.Quote(remoteBranchName))).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *BranchCommands) SetUpstream(remoteName string, remoteBranchName string, branchName string) error {
	cmdStr := NewGitCmd("branch").
		Arg(fmt.Sprintf("--set-upstream-to=%s/%s", self.cmd.Quote(remoteName), self.cmd.Quote(remoteBranchName))).
		Arg(self.cmd.Quote(branchName)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *BranchCommands) UnsetUpstream(branchName string) error {
	cmdStr := NewGitCmd("branch").Arg("--unset-upstream", self.cmd.Quote(branchName)).
		ToString()

	return self.cmd.New(cmdStr).Run()
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
	cmdStr := NewGitCmd("rev-list").
		Arg(fmt.Sprintf("%s..%s", from, to)).
		Arg("--count").
		ToString()

	return self.cmd.New(cmdStr).DontLog().RunWithOutput()
}

func (self *BranchCommands) IsHeadDetached() bool {
	cmdStr := NewGitCmd("symbolic-ref").Arg("-q", "HEAD").ToString()

	err := self.cmd.New(cmdStr).DontLog().Run()
	return err != nil
}

func (self *BranchCommands) Rename(oldName string, newName string) error {
	cmdStr := NewGitCmd("branch").
		Arg("--move", self.cmd.Quote(oldName), self.cmd.Quote(newName)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *BranchCommands) GetRawBranches() (string, error) {
	cmdStr := NewGitCmd("for-each-ref").
		Arg("--sort=-committerdate").
		Arg(`--format="%(HEAD)%00%(refname:short)%00%(upstream:short)%00%(upstream:track)"`).
		Arg("refs/heads").
		ToString()

	return self.cmd.New(cmdStr).DontLog().RunWithOutput()
}

type MergeOpts struct {
	FastForwardOnly bool
}

func (self *BranchCommands) Merge(branchName string, opts MergeOpts) error {
	command := NewGitCmd("merge").
		Arg("--no-edit").
		ArgIf(self.UserConfig.Git.Merging.Args != "", self.UserConfig.Git.Merging.Args).
		ArgIf(opts.FastForwardOnly, "--ff-only").
		Arg(self.cmd.Quote(branchName)).
		ToString()

	return self.cmd.New(command).Run()
}

func (self *BranchCommands) AllBranchesLogCmdObj() oscommands.ICmdObj {
	return self.cmd.New(self.UserConfig.Git.AllBranchesLogCmd).DontLog()
}
