package git_commands

import (
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Returns a set of branch names which are merged into a main branch.
// This
type MergedBranchLoader struct {
	c *GitCommon
}

func NewMergedBranchLoader(c *GitCommon) *MergedBranchLoader {
	return &MergedBranchLoader{c: c}
}

// TODO: check against upstreams, and share code with commit loader that determines main branches
func (self *MergedBranchLoader) Load() *set.Set[string] {
	set := set.New[string]()

	mainBranches := self.c.UserConfig.Git.MainBranches
	results := utils.ConcurrentMap(mainBranches, func(mainBranch string) []string {
		mergedBranches, err := self.GetMergedBranches(mainBranch)
		if err != nil {
			self.c.Log.Warnf("Failed to get merged branches for %s: %s", mainBranch, err)
			return nil
		}
		return mergedBranches
	})

	for i := 0; i < len(mainBranches); i++ {
		for _, mergedBranch := range results[i] {
			set.Add(mergedBranch)
		}
	}

	return set
}

func (self *MergedBranchLoader) GetMergedBranches(mainBranch string) ([]string, error) {
	// git for-each-ref --merged master --format '%(refname)' refs/heads/
	cmdArgs := NewGitCmd("for-each-ref").Arg(
		"--merged",
		mainBranch,
		"--format",
		"%(refname)",
		"refs/heads/",
	).ToArgv()

	output, err := self.c.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	branches := utils.SplitLines(output)

	return branches, nil
}
