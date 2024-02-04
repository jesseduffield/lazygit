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
		} else {
			return "", false
		}
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

	return configs, nil
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

func (self *SubmoduleCommands) Delete(submodule *models.SubmoduleConfig) error {
	// based on https://gist.github.com/myusuf3/7f645819ded92bda6677

	if submodule.ParentModule != nil {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		err = os.Chdir(submodule.ParentModule.FullPath())
		if err != nil {
			return err
		}

		defer func() { _ = os.Chdir(wd) }()
	}

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
	if submodule.ParentModule != nil {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		err = os.Chdir(submodule.ParentModule.FullPath())
		if err != nil {
			return err
		}

		defer func() { _ = os.Chdir(wd) }()
	}

	setUrlCmdStr := NewGitCmd("config").
		Arg(
			"--file", ".gitmodules", "submodule."+submodule.Name+".url", newUrl,
		).
		ToArgv()

	// the set-url command is only for later git versions so we're doing it manually here
	if err := self.cmd.New(setUrlCmdStr).Run(); err != nil {
		return err
	}

	syncCmdStr := NewGitCmd("submodule").Arg("sync", "--", submodule.Path).
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
