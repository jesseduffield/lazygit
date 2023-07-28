package git_commands

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
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

func (self *SubmoduleCommands) GetConfigs() ([]*models.SubmoduleConfig, error) {
	file, err := os.Open(".gitmodules")
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
		} else {
			return "", false
		}
	}

	configs := []*models.SubmoduleConfig{}
	for scanner.Scan() {
		line := scanner.Text()

		if name, ok := firstMatch(line, `\[submodule "(.*)"\]`); ok {
			configs = append(configs, &models.SubmoduleConfig{Name: name})
			continue
		}

		if len(configs) > 0 {
			lastConfig := configs[len(configs)-1]

			if path, ok := firstMatch(line, `\s*path\s*=\s*(.*)\s*`); ok {
				lastConfig.Path = path
			} else if url, ok := firstMatch(line, `\s*url\s*=\s*(.*)\s*`); ok {
				lastConfig.Url = url
			}
		}
	}

	return configs, nil
}

func (self *SubmoduleCommands) Stash(submodule *models.SubmoduleConfig) error {
	// if the path does not exist then it hasn't yet been initialized so we'll swallow the error
	// because the intention here is to have no dirty worktree state
	if _, err := os.Stat(submodule.Path); os.IsNotExist(err) {
		self.Log.Infof("submodule path %s does not exist, returning", submodule.Path)
		return nil
	}

	cmdArgs := NewGitCmd("stash").
		Dir(submodule.Path).
		Arg("--include-untracked").
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *SubmoduleCommands) Reset(submodule *models.SubmoduleConfig) error {
	cmdArgs := NewGitCmd("submodule").
		Arg("update", "--init", "--force", "--", submodule.Path).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *SubmoduleCommands) UpdateAll() error {
	// not doing an --init here because the user probably doesn't want that
	cmdArgs := NewGitCmd("submodule").Arg("update", "--force").ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *SubmoduleCommands) Delete(submodule *models.SubmoduleConfig) error {
	// based on https://gist.github.com/myusuf3/7f645819ded92bda6677

	if err := self.cmd.New(
		NewGitCmd("submodule").
			Arg("deinit", "--force", "--", submodule.Path).ToArgv(),
	).Run(); err != nil {
		if !strings.Contains(err.Error(), "did not match any file(s) known to git") {
			return err
		}

		if err := self.cmd.New(
			NewGitCmd("config").
				Arg("--file", ".gitmodules", "--remove-section", "submodule."+submodule.Path).
				ToArgv(),
		).Run(); err != nil {
			return err
		}

		if err := self.cmd.New(
			NewGitCmd("config").
				Arg("--remove-section", "submodule."+submodule.Path).
				ToArgv(),
		).Run(); err != nil {
			return err
		}
	}

	if err := self.cmd.New(
		NewGitCmd("rm").Arg("--force", "-r", submodule.Path).ToArgv(),
	).Run(); err != nil {
		// if the directory isn't there then that's fine
		self.Log.Error(err)
	}

	// We may in fact want to use the repo's git dir path but git docs say not to
	// mix submodules and worktrees anyway.
	return os.RemoveAll(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "modules", submodule.Path))
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

func (self *SubmoduleCommands) UpdateUrl(name string, path string, newUrl string) error {
	setUrlCmdStr := NewGitCmd("config").
		Arg(
			"--file", ".gitmodules", "submodule."+name+".url", newUrl,
		).
		ToArgv()

	// the set-url command is only for later git versions so we're doing it manually here
	if err := self.cmd.New(setUrlCmdStr).Run(); err != nil {
		return err
	}

	syncCmdStr := NewGitCmd("submodule").Arg("sync", "--", path).
		ToArgv()

	if err := self.cmd.New(syncCmdStr).Run(); err != nil {
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

func (self *SubmoduleCommands) BulkInitCmdObj() oscommands.ICmdObj {
	cmdArgs := NewGitCmd("submodule").Arg("init").
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *SubmoduleCommands) BulkUpdateCmdObj() oscommands.ICmdObj {
	cmdArgs := NewGitCmd("submodule").Arg("update").
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *SubmoduleCommands) ForceBulkUpdateCmdObj() oscommands.ICmdObj {
	cmdArgs := NewGitCmd("submodule").Arg("update", "--force").
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *SubmoduleCommands) BulkDeinitCmdObj() oscommands.ICmdObj {
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
