package git_commands

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
)

// .gitmodules looks like this:
// [submodule "mysubmodule"]
//   path = blah/mysubmodule
//   url = git@github.com:subbo.git

type SubmoduleCommands struct {
	*GitCommon
}

func NewSubmoduleCommands(gitCommon *GitCommon) *SubmoduleCommands {
	return &SubmoduleCommands{
		GitCommon: gitCommon,
	}
}

func (self *SubmoduleCommands) GetConfigs(parentModule *models.SubmoduleConfig) ([]*models.SubmoduleConfig, error) {
	gitModulesPath := ".gitmodules"
	if parentModule != nil {
		gitModulesPath = filepath.Join(parentModule.FullPath(), gitModulesPath)
	}
	file, err := os.Open(gitModulesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	firstMatch := func(str string, regex string) (string, bool) {
		re := regexp.MustCompile(regex)
		matches := re.FindStringSubmatch(str)

		if len(matches) > 0 {
			return matches[1], true
		}
		return "", false
	}

	configs := []*models.SubmoduleConfig{}
	lastConfigIdx := -1
	for scanner.Scan() {
		line := scanner.Text()

		if name, ok := firstMatch(line, `\[submodule "(.*)"\]`); ok {
			configs = append(configs, &models.SubmoduleConfig{
				Name: name, ParentModule: parentModule,
			})
			lastConfigIdx = len(configs) - 1
			continue
		}

		if lastConfigIdx != -1 {
			if path, ok := firstMatch(line, `\s*path\s*=\s*(.*)\s*`); ok {
				configs[lastConfigIdx].Path = path
				nestedConfigs, err := self.GetConfigs(configs[lastConfigIdx])
				if err == nil {
					configs = append(configs, nestedConfigs...)
				}
			} else if url, ok := firstMatch(line, `\s*url\s*=\s*(.*)\s*`); ok {
				configs[lastConfigIdx].Url = url
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

// AnyHaveStageableChanges reports whether any of the given submodule paths has
// a checked-out commit that differs from the one recorded in the
// superproject's index, i.e. a change that `git add <path>` would actually
// stage. A submodule that only has dirty or untracked content (with no new
// commit) can't be staged from the superproject, so it won't be reported here.
func (self *SubmoduleCommands) AnyHaveStageableChanges(paths []string) (bool, error) {
	if len(paths) == 0 {
		return false, nil
	}

	cmdArgs := NewGitCmd("submodule").Arg("status", "--").Arg(paths...).ToArgv()
	output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return false, err
	}

	// Each line looks like "<prefix><sha> <path> (<describe>)". A '+' prefix
	// means the checked-out commit differs from the index, i.e. there's a
	// commit change to stage.
	return lo.SomeBy(strings.Split(output, "\n"), func(line string) bool {
		return strings.HasPrefix(line, "+")
	}), nil
}

// GetConflictCommits returns the three gitlink commits of a conflicted submodule
// from the index: the merge base, our (current) commit, and their (incoming)
// commit. Any of them can be empty if that stage is absent (e.g. a submodule
// that was added on only one side). The path is relative to the repo root.
func (self *SubmoduleCommands) GetConflictCommits(path string) (base string, ours string, theirs string, err error) {
	cmdArgs := NewGitCmd("ls-files").Arg("-u", "-z", "--", path).ToArgv()
	output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return "", "", "", err
	}

	// Each NUL-terminated entry looks like "<mode> <sha> <stage>\t<path>".
	for _, entry := range strings.Split(output, "\x00") {
		// fields are split on the tab and the spaces, so the leading three are
		// always mode, sha, stage regardless of what the path contains.
		fields := strings.Fields(entry)
		if len(fields) < 3 {
			continue
		}
		switch fields[2] {
		case "1":
			base = fields[1]
		case "2":
			ours = fields[1]
		case "3":
			theirs = fields[1]
		}
	}

	return base, ours, theirs, nil
}

// GetCommitSummary returns "<short-sha> <subject>" for a commit inside the
// submodule at the given path, for display in the conflict menu.
func (self *SubmoduleCommands) GetCommitSummary(path string, sha string) (string, error) {
	cmdArgs := NewGitCmd("log").
		Dir(path).
		Arg("--format=%h %s", "--max-count=1", sha).
		Config("log.showsignature=false").
		ToArgv()

	summary, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	return strings.TrimSpace(summary), err
}

// CheckoutConflictCommit resolves a submodule conflict by checking the submodule
// out at the given commit. `git checkout --ours/--theirs` is a no-op on
// gitlinks, so we check out the chosen commit in the submodule itself; the
// caller then stages the submodule to record the resolution.
func (self *SubmoduleCommands) CheckoutConflictCommit(path string, sha string) error {
	cmdArgs := NewGitCmd("checkout").Dir(path).Arg(sha).ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// ConflictSideLog returns a oneline log, run inside the submodule, of the commits
// that `side` has but `otherSide` does not (i.e. `otherSide..side`) — the commits
// unique to one side of a commit conflict, relative to their common ancestor. It
// is empty if `side` is an ancestor of `otherSide` (e.g. that side was rewound).
func (self *SubmoduleCommands) ConflictSideLog(path string, side string, otherSide string) (string, error) {
	cmdArgs := NewGitCmd("log").Dir(path).
		Arg("--oneline", "--color=always", otherSide+".."+side).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

func (self *SubmoduleCommands) Stash(submodule *models.SubmoduleConfig) error {
	// if the path does not exist then it hasn't yet been initialized so we'll swallow the error
	// because the intention here is to have no dirty worktree state
	if _, err := os.Stat(submodule.Path); os.IsNotExist(err) {
		self.Log.Infof("submodule path %s does not exist, returning", submodule.FullPath())
		return nil
	}

	cmdArgs := NewGitCmd("stash").
		Dir(submodule.FullPath()).
		Arg("--include-untracked").
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *SubmoduleCommands) Reset(submodule *models.SubmoduleConfig) error {
	parentDir := ""
	if submodule.ParentModule != nil {
		parentDir = submodule.ParentModule.FullPath()
	}
	cmdArgs := NewGitCmd("submodule").
		Arg("update", "--init", "--force", "--", submodule.Path).
		DirIf(parentDir != "", parentDir).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *SubmoduleCommands) UpdateAll() error {
	// not doing an --init here because the user probably doesn't want that
	cmdArgs := NewGitCmd("submodule").Arg("update", "--force").ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// runInParentModule runs the given command in the submodule's parent module's
// directory when the submodule is nested: its path arguments (and the
// .gitmodules file the config commands touch) are relative to the parent
// module. The directory is set on the command itself rather than by
// temporarily chdir-ing the process there, which would leak the parent
// module's directory into whatever other commands run concurrently (e.g. a
// background refresh's).
func (self *SubmoduleCommands) runInParentModule(submodule *models.SubmoduleConfig, cmdObj *oscommands.CmdObj) error {
	if submodule.ParentModule != nil {
		cmdObj.SetWd(submodule.ParentModule.FullPath())
	}
	return cmdObj.Run()
}

func (self *SubmoduleCommands) Delete(submodule *models.SubmoduleConfig) error {
	// based on https://gist.github.com/myusuf3/7f645819ded92bda6677

	if err := self.runInParentModule(submodule, self.cmd.New(
		NewGitCmd("submodule").
			Arg("deinit", "--force", "--", submodule.Path).ToArgv(),
	)); err != nil {
		if !strings.Contains(err.Error(), "did not match any file(s) known to git") {
			return err
		}

		if err := self.runInParentModule(submodule, self.cmd.New(
			NewGitCmd("config").
				Arg("--file", ".gitmodules", "--remove-section", "submodule."+submodule.Path).
				ToArgv(),
		)); err != nil {
			return err
		}

		if err := self.runInParentModule(submodule, self.cmd.New(
			NewGitCmd("config").
				Arg("--remove-section", "submodule."+submodule.Path).
				ToArgv(),
		)); err != nil {
			return err
		}
	}

	if err := self.runInParentModule(submodule, self.cmd.New(
		NewGitCmd("rm").Arg("--force", "-r", submodule.Path).ToArgv(),
	)); err != nil {
		// if the directory isn't there then that's fine
		self.Log.Error(err)
	}

	// We may in fact want to use the repo's git dir path but git docs say not to
	// mix submodules and worktrees anyway.
	return os.RemoveAll(submodule.GitDirPath(self.repoPaths.repoGitDirPath))
}

func (self *SubmoduleCommands) Add(name string, path string, url string) error {
	cmdArgs := NewGitCmd("submodule").
		Arg("add").
		Arg("--force").
		Arg("--name").
		Arg(name).
		Arg("--").
		Arg(url).
		Arg(path).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *SubmoduleCommands) UpdateUrl(submodule *models.SubmoduleConfig, newUrl string) error {
	setUrlCmdStr := NewGitCmd("config").
		Arg(
			"--file", ".gitmodules", "submodule."+submodule.Name+".url", newUrl,
		).
		ToArgv()

	// the set-url command is only for later git versions so we're doing it manually here
	if err := self.runInParentModule(submodule, self.cmd.New(setUrlCmdStr)); err != nil {
		return err
	}

	syncCmdStr := NewGitCmd("submodule").Arg("sync", "--", submodule.Path).
		ToArgv()

	if err := self.runInParentModule(submodule, self.cmd.New(syncCmdStr)); err != nil {
		return err
	}

	return nil
}

func (self *SubmoduleCommands) Init(path string) error {
	cmdArgs := NewGitCmd("submodule").Arg("init", "--", path).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *SubmoduleCommands) Update(path string) error {
	cmdArgs := NewGitCmd("submodule").Arg("update", "--init", "--", path).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *SubmoduleCommands) BulkInitCmdObj() *oscommands.CmdObj {
	cmdArgs := NewGitCmd("submodule").Arg("init").
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *SubmoduleCommands) BulkUpdateCmdObj() *oscommands.CmdObj {
	cmdArgs := NewGitCmd("submodule").Arg("update").
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *SubmoduleCommands) ForceBulkUpdateCmdObj() *oscommands.CmdObj {
	cmdArgs := NewGitCmd("submodule").Arg("update", "--force").
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *SubmoduleCommands) BulkUpdateRecursivelyCmdObj() *oscommands.CmdObj {
	cmdArgs := NewGitCmd("submodule").Arg("update", "--init", "--recursive").
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *SubmoduleCommands) BulkDeinitCmdObj() *oscommands.CmdObj {
	cmdArgs := NewGitCmd("submodule").Arg("deinit", "--all", "--force").
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *SubmoduleCommands) ResetSubmodules(submodules []*models.SubmoduleConfig) error {
	for _, submodule := range submodules {
		if err := self.Stash(submodule); err != nil {
			return err
		}
	}

	return self.UpdateAll()
}
