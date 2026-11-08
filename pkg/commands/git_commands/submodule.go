package git_commands

import (
	"bufio"
	"fmt"
	"io"
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
	// Resolve the path against the repo this commands object was created for
	// rather than the process working directory, so that a read from a
	// still-running refresh keeps addressing that repo after the user
	// switched to another one.
	dir := self.repoPaths.WorktreePath()
	if parentModule != nil {
		dir = filepath.Join(dir, parentModule.FullPath())
	}
	gitModulesPath := filepath.Join(dir, ".gitmodules")
	file, err := os.Open(gitModulesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	parsedConfigs, err := parseGitmodules(file, parentModule)
	if err != nil {
		return nil, err
	}

	configs := []*models.SubmoduleConfig{}
	for _, config := range parsedConfigs {
		if isSuspiciousSubmoduleName(config.Name) {
			self.Log.Warnf("Ignoring suspicious submodule name: %s", config.Name)
			continue
		}

		if isSuspiciousSubmodulePath(dir, config.Path) {
			self.Log.Warnf("Ignoring suspicious submodule path: %s", config.Path)
			continue
		}

		configs = append(configs, config)

		nestedConfigs, err := self.GetConfigs(config)
		if err == nil {
			configs = append(configs, nestedConfigs...)
		}
	}

	return configs, nil
}

func parseGitmodules(reader io.Reader, parentModule *models.SubmoduleConfig) ([]*models.SubmoduleConfig, error) {
	scanner := bufio.NewScanner(reader)
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
	for scanner.Scan() {
		line := scanner.Text()

		if name, ok := firstMatch(line, `^\s*\[submodule "(.*)"\]`); ok {
			configs = append(configs, &models.SubmoduleConfig{
				Name: name, ParentModule: parentModule,
			})
			continue
		}

		if len(configs) > 0 {
			config := configs[len(configs)-1]
			if path, ok := firstMatch(line, `^\s*path\s*=\s*(.*)\s*`); ok {
				config.Path = path
			} else if url, ok := firstMatch(line, `^\s*url\s*=\s*(.*)\s*`); ok {
				config.Url = url
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

// isSuspiciousSubmoduleName applies the same rule as git's
// check_submodule_name(); git ignores a submodule whose name fails it. The
// git dir of a submodule is .git/modules/<name>, so a ".." component in the
// name could put it outside the repository. Like git, we treat both slashes
// and backslashes as separators on every platform.
func isSuspiciousSubmoduleName(name string) bool {
	if name == "" {
		return true
	}

	components := strings.FieldsFunc(name, func(r rune) bool {
		return r == '/' || r == '\\'
	})
	return lo.Contains(components, "..")
}

// isSuspiciousSubmodulePath reports whether a submodule path from the
// .gitmodules file in parentDir might lead somewhere other than a directory
// below parentDir. We run commands in the submodule's directory and read the
// .gitmodules file in it, so the path must not lead anywhere else.
//
// Git never checks out a file through a symlink, so a path that goes through
// a symlink can't be a submodule in the first place. We reject such paths
// even if the symlink points to a directory below parentDir.
func isSuspiciousSubmodulePath(parentDir string, path string) bool {
	if !filepath.IsLocal(path) || filepath.Clean(path) == "." {
		return true
	}

	dir := parentDir
	for _, component := range strings.Split(filepath.Clean(path), string(filepath.Separator)) {
		dir = filepath.Join(dir, component)
		fileInfo, err := os.Lstat(dir)
		if os.IsNotExist(err) {
			// Nothing below this point exists (e.g. because the submodule
			// isn't checked out), so it can't lead anywhere else either
			return false
		}
		if err != nil || !fileInfo.IsDir() {
			return true
		}
	}

	return false
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

	summary, err := ForOtherRepo(self.cmd.New(cmdArgs)).DontLog().RunWithOutput()
	return strings.TrimSpace(summary), err
}

// CheckoutConflictCommit resolves a submodule conflict by checking the submodule
// out at the given commit. `git checkout --ours/--theirs` is a no-op on
// gitlinks, so we check out the chosen commit in the submodule itself; the
// caller then stages the submodule to record the resolution.
func (self *SubmoduleCommands) CheckoutConflictCommit(path string, sha string) error {
	cmdArgs := NewGitCmd("checkout").Dir(path).Arg(sha).ToArgv()
	return ForOtherRepo(self.cmd.New(cmdArgs)).Run()
}

// ConflictSideLog returns a oneline log, run inside the submodule, of the commits
// that `side` has but `otherSide` does not (i.e. `otherSide..side`) — the commits
// unique to one side of a commit conflict, relative to their common ancestor. It
// is empty if `side` is an ancestor of `otherSide` (e.g. that side was rewound).
func (self *SubmoduleCommands) ConflictSideLog(path string, side string, otherSide string) (string, error) {
	cmdArgs := NewGitCmd("log").Dir(path).
		Arg("--oneline", "--color=always", otherSide+".."+side).
		ToArgv()

	return ForOtherRepo(self.cmd.New(cmdArgs)).DontLog().RunWithOutput()
}

func (self *SubmoduleCommands) Stash(submodule *models.SubmoduleConfig) error {
	// if the path does not exist then it hasn't yet been initialized so we'll swallow the error
	// because the intention here is to have no dirty worktree state
	if _, err := os.Stat(filepath.Join(self.repoPaths.WorktreePath(), submodule.FullPath())); os.IsNotExist(err) {
		self.Log.Infof("submodule path %s does not exist, returning", submodule.FullPath())
		return nil
	}

	cmdArgs := NewGitCmd("stash").
		Dir(submodule.FullPath()).
		Arg("--include-untracked").
		ToArgv()

	return ForOtherRepo(self.cmd.New(cmdArgs)).Run()
}

func (self *SubmoduleCommands) Reset(submodule *models.SubmoduleConfig) error {
	cmdArgs := NewGitCmd("submodule").
		Arg("update", "--init", "--force", "--", submodule.Path).
		ToArgv()

	return self.runInParentModule(submodule, self.cmd.New(cmdArgs))
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
//
// That directory is relative, so it resolves against the process working
// directory rather than against the repo directory the command builder
// otherwise pins commands to. Only foreground commands the user issued end up
// here, and lazygit won't switch repos while one of those is in flight, so the
// two are the same directory; don't call this from background work, where they
// need not be.
func (self *SubmoduleCommands) runInParentModule(submodule *models.SubmoduleConfig, cmdObj *oscommands.CmdObj) error {
	if submodule.ParentModule != nil {
		ForOtherRepo(cmdObj.SetWd(submodule.ParentModule.FullPath()))
	}
	return cmdObj.Run()
}

func (self *SubmoduleCommands) Delete(submodule *models.SubmoduleConfig) error {
	// based on https://gist.github.com/myusuf3/7f645819ded92bda6677

	// Git keeps the git dirs of submodules separately for each worktree
	gitDirPath := submodule.GitDirPath(self.repoPaths.WorktreeGitDirPath())

	// The name of the submodule comes from .gitmodules, so it is under the
	// control of whoever committed that file. Make sure that we never delete
	// anything other than a directory below the "modules" directory of the
	// worktree's git dir.
	modulesDirPath := filepath.Join(self.repoPaths.WorktreeGitDirPath(), "modules")
	if relPath, err := filepath.Rel(modulesDirPath, gitDirPath); err != nil || relPath == "." || !filepath.IsLocal(relPath) {
		return fmt.Errorf("refusing to delete %s because it is not inside %s", gitDirPath, modulesDirPath)
	}

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

	return os.RemoveAll(gitDirPath)
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
