package git_commands

import (
	"bufio"
	"fmt"
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

	return self.cmd.New("git -C " + self.cmd.Quote(submodule.Path) + " stash --include-untracked").Run()
}

func (self *SubmoduleCommands) Reset(submodule *models.SubmoduleConfig) error {
	return self.cmd.New("git submodule update --init --force -- " + self.cmd.Quote(submodule.Path)).Run()
}

func (self *SubmoduleCommands) UpdateAll() error {
	// not doing an --init here because the user probably doesn't want that
	return self.cmd.New("git submodule update --force").Run()
}

func (self *SubmoduleCommands) Delete(submodule *models.SubmoduleConfig) error {
	// based on https://gist.github.com/myusuf3/7f645819ded92bda6677

	if err := self.cmd.New("git submodule deinit --force -- " + self.cmd.Quote(submodule.Path)).Run(); err != nil {
		if strings.Contains(err.Error(), "did not match any file(s) known to git") {
			if err := self.cmd.New("git config --file .gitmodules --remove-section submodule." + self.cmd.Quote(submodule.Name)).Run(); err != nil {
				return err
			}

			if err := self.cmd.New("git config --remove-section submodule." + self.cmd.Quote(submodule.Name)).Run(); err != nil {
				return err
			}

			// if there's an error here about it not existing then we'll just continue to do `git rm`
		} else {
			return err
		}
	}

	if err := self.cmd.New("git rm --force -r " + submodule.Path).Run(); err != nil {
		// if the directory isn't there then that's fine
		self.Log.Error(err)
	}

	return os.RemoveAll(filepath.Join(self.dotGitDir, "modules", submodule.Path))
}

func (self *SubmoduleCommands) Add(name string, path string, url string) error {
	return self.cmd.
		New(
			fmt.Sprintf(
				"git submodule add --force --name %s -- %s %s ",
				self.cmd.Quote(name),
				self.cmd.Quote(url),
				self.cmd.Quote(path),
			)).
		Run()
}

func (self *SubmoduleCommands) UpdateUrl(name string, path string, newUrl string) error {
	// the set-url command is only for later git versions so we're doing it manually here
	if err := self.cmd.New("git config --file .gitmodules submodule." + self.cmd.Quote(name) + ".url " + self.cmd.Quote(newUrl)).Run(); err != nil {
		return err
	}

	if err := self.cmd.New("git submodule sync -- " + self.cmd.Quote(path)).Run(); err != nil {
		return err
	}

	return nil
}

func (self *SubmoduleCommands) Init(path string) error {
	return self.cmd.New("git submodule init -- " + self.cmd.Quote(path)).Run()
}

func (self *SubmoduleCommands) Update(path string) error {
	return self.cmd.New("git submodule update --init -- " + self.cmd.Quote(path)).Run()
}

func (self *SubmoduleCommands) BulkInitCmdObj() oscommands.ICmdObj {
	return self.cmd.New("git submodule init")
}

func (self *SubmoduleCommands) BulkUpdateCmdObj() oscommands.ICmdObj {
	return self.cmd.New("git submodule update")
}

func (self *SubmoduleCommands) ForceBulkUpdateCmdObj() oscommands.ICmdObj {
	return self.cmd.New("git submodule update --force")
}

func (self *SubmoduleCommands) BulkDeinitCmdObj() oscommands.ICmdObj {
	return self.cmd.New("git submodule deinit --all --force")
}

func (self *SubmoduleCommands) ResetSubmodules(submodules []*models.SubmoduleConfig) error {
	for _, submodule := range submodules {
		if err := self.Stash(submodule); err != nil {
			return err
		}
	}

	return self.UpdateAll()
}
