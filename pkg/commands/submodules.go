package commands

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

func (c *GitCommand) GetSubmoduleConfigs() ([]*models.SubmoduleConfig, error) {
	file, err := os.Open(".gitmodules")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

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

func (c *GitCommand) SubmoduleStash(submodule *models.SubmoduleConfig) error {
	// if the path does not exist then it hasn't yet been initialized so we'll swallow the error
	// because the intention here is to have no dirty worktree state
	if _, err := os.Stat(submodule.Path); os.IsNotExist(err) {
		c.Log.Infof("submodule path %s does not exist, returning", submodule.Path)
		return nil
	}

	return c.Run(c.NewCmdObj("git -C " + c.OSCommand.Quote(submodule.Path) + " stash --include-untracked"))
}

func (c *GitCommand) SubmoduleReset(submodule *models.SubmoduleConfig) error {
	return c.Run(c.NewCmdObj("git submodule update --init --force -- " + c.OSCommand.Quote(submodule.Path)))
}

func (c *GitCommand) SubmoduleUpdateAll() error {
	// not doing an --init here because the user probably doesn't want that
	return c.Run(c.NewCmdObj("git submodule update --force"))
}

func (c *GitCommand) SubmoduleDelete(submodule *models.SubmoduleConfig) error {
	// based on https://gist.github.com/myusuf3/7f645819ded92bda6677

	if err := c.Run(c.NewCmdObj("git submodule deinit --force -- " + c.OSCommand.Quote(submodule.Path))); err != nil {
		if strings.Contains(err.Error(), "did not match any file(s) known to git") {
			if err := c.Run(c.NewCmdObj("git config --file .gitmodules --remove-section submodule." + c.OSCommand.Quote(submodule.Name))); err != nil {
				return err
			}

			if err := c.Run(c.NewCmdObj("git config --remove-section submodule." + c.OSCommand.Quote(submodule.Name))); err != nil {
				return err
			}

			// if there's an error here about it not existing then we'll just continue to do `git rm`
		} else {
			return err
		}
	}

	if err := c.Run(c.NewCmdObj("git rm --force -r " + submodule.Path)); err != nil {
		// if the directory isn't there then that's fine
		c.Log.Error(err)
	}

	return os.RemoveAll(filepath.Join(c.DotGitDir, "modules", submodule.Path))
}

func (c *GitCommand) SubmoduleAdd(name string, path string, url string) error {
	return c.OSCommand.Run(
		c.OSCommand.NewCmdObj(
			fmt.Sprintf(
				"git submodule add --force --name %s -- %s %s ",
				c.OSCommand.Quote(name),
				c.OSCommand.Quote(url),
				c.OSCommand.Quote(path),
			)))
}

func (c *GitCommand) SubmoduleUpdateUrl(name string, path string, newUrl string) error {
	// the set-url command is only for later git versions so we're doing it manually here
	if err := c.Run(c.NewCmdObj("git config --file .gitmodules submodule." + c.OSCommand.Quote(name) + ".url " + c.OSCommand.Quote(newUrl))); err != nil {
		return err
	}

	if err := c.Run(c.NewCmdObj("git submodule sync -- " + c.OSCommand.Quote(path))); err != nil {
		return err
	}

	return nil
}

func (c *GitCommand) SubmoduleInit(path string) error {
	return c.Run(c.NewCmdObj("git submodule init -- " + c.OSCommand.Quote(path)))
}

func (c *GitCommand) SubmoduleUpdate(path string) error {
	return c.Run(c.NewCmdObj("git submodule update --init -- " + c.OSCommand.Quote(path)))
}

func (c *GitCommand) SubmoduleBulkInitCmdObj() oscommands.ICmdObj {
	return c.NewCmdObj("git submodule init")
}

func (c *GitCommand) SubmoduleBulkUpdateCmdObj() oscommands.ICmdObj {
	return c.NewCmdObj("git submodule update")
}

func (c *GitCommand) SubmoduleForceBulkUpdateCmdObj() oscommands.ICmdObj {
	return c.NewCmdObj("git submodule update --force")
}

func (c *GitCommand) SubmoduleBulkDeinitCmdObj() oscommands.ICmdObj {
	return c.NewCmdObj("git submodule deinit --all --force")
}

func (c *GitCommand) ResetSubmodules(submodules []*models.SubmoduleConfig) error {
	for _, submodule := range submodules {
		if err := c.SubmoduleStash(submodule); err != nil {
			return err
		}
	}

	return c.SubmoduleUpdateAll()
}
