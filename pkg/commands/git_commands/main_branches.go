package git_commands

import (
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/sasha-s/go-deadlock"
)

type MainBranches struct {
	// List of main branches configured by the user. Just the bare names.
	configuredMainBranches []string
	// Which of these actually exist in the repository. Full ref names, and it
	// could be either "refs/heads/..." or "refs/remotes/origin/..." depending
	// on which one exists for a given bare name.
	existingMainBranches []string

	cmd   oscommands.ICmdObjBuilder
	mutex *deadlock.Mutex
}

func NewMainBranches(
	configuredMainBranches []string,
	cmd oscommands.ICmdObjBuilder,
) *MainBranches {
	return &MainBranches{
		configuredMainBranches: configuredMainBranches,
		existingMainBranches:   nil,
		cmd:                    cmd,
		mutex:                  &deadlock.Mutex{},
	}
}

// Get the list of main branches that exist in the repository. This is a list of
// full ref names.
func (self *MainBranches) Get() []string {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.existingMainBranches == nil {
		self.existingMainBranches = self.determineMainBranches()
	}

	return self.existingMainBranches
}

func (self *MainBranches) determineMainBranches() []string {
	var existingBranches []string
	var wg sync.WaitGroup

	existingBranches = make([]string, len(self.configuredMainBranches))

	for i, branchName := range self.configuredMainBranches {
		wg.Add(1)
		go utils.Safe(func() {
			defer wg.Done()

			// Try to determine upstream of local main branch
			if ref, err := self.cmd.New(
				NewGitCmd("rev-parse").Arg("--symbolic-full-name", branchName+"@{u}").ToArgv(),
			).DontLog().RunWithOutput(); err == nil {
				existingBranches[i] = strings.TrimSpace(ref)
				return
			}

			// If this failed, a local branch for this main branch doesn't exist or it
			// has no upstream configured. Try looking for one in the "origin" remote.
			ref := "refs/remotes/origin/" + branchName
			if err := self.cmd.New(
				NewGitCmd("rev-parse").Arg("--verify", "--quiet", ref).ToArgv(),
			).DontLog().Run(); err == nil {
				existingBranches[i] = ref
				return
			}

			// If this failed as well, try if we have the main branch as a local
			// branch. This covers the case where somebody is using git locally
			// for something, but never pushing anywhere.
			ref = "refs/heads/" + branchName
			if err := self.cmd.New(
				NewGitCmd("rev-parse").Arg("--verify", "--quiet", ref).ToArgv(),
			).DontLog().Run(); err == nil {
				existingBranches[i] = ref
			}
		})
	}

	wg.Wait()

	existingBranches = lo.Filter(existingBranches, func(branch string, _ int) bool {
		return branch != ""
	})

	return existingBranches
}
