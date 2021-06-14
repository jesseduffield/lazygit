package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
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
		c.log.Infof("submodule path %s does not exist, returning", submodule.Path)
		return nil
	}

	return c.RunGitCmdFromStr(fmt.Sprintf("-C %s stash --include-untracked", submodule.Path))
}

func (c *GitCommand) SubmoduleReset(submodule *models.SubmoduleConfig) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("submodule update --init --force %s", submodule.Path))
}

func (c *GitCommand) SubmoduleDelete(submodule *models.SubmoduleConfig) error {
	// based on https://gist.github.com/myusuf3/7f645819ded92bda6677

	if err := c.RunGitCmdFromStr(fmt.Sprintf("submodule deinit --force %s", submodule.Path)); err != nil {
		if strings.Contains(err.Error(), "did not match any file(s) known to git") {
			if err := c.RunGitCmdFromStr(fmt.Sprintf("config --file .gitmodules --remove-section submodule.%s", submodule.Name)); err != nil {
				return err
			}

			if err := c.RunGitCmdFromStr(fmt.Sprintf("config --remove-section submodule.%s", submodule.Name)); err != nil {
				return err
			}

			// if there's an error here about it not existing then we'll just continue to do `git rm`
		} else {
			return err
		}
	}

	if err := c.RunGitCmdFromStr(fmt.Sprintf("rm --force -r %s", submodule.Path)); err != nil {
		// if the directory isn't there then that's fine
		c.log.Error(err)
	}

	return os.RemoveAll(filepath.Join(c.dotGitDir, "modules", submodule.Path))
}

func (c *GitCommand) SubmoduleAdd(name string, path string, url string) error {
	return c.RunGitCmdFromStr(
		fmt.Sprintf(
			"submodule add --force --name %s -- %s %s ",
			c.GetOSCommand().Quote(name),
			c.GetOSCommand().Quote(url),
			c.GetOSCommand().Quote(path),
		),
	)
}

func (c *GitCommand) SubmoduleUpdateUrl(name string, path string, newUrl string) error {
	// the set-url command is only for later git versions so we're doing it manually here
	if err := c.RunGitCmdFromStr(fmt.Sprintf("config --file .gitmodules submodule.%s.url %s", name, newUrl)); err != nil {
		return err
	}

	if err := c.RunGitCmdFromStr(fmt.Sprintf("submodule sync %s", path)); err != nil {
		return err
	}

	return nil
}

func (c *GitCommand) SubmoduleInit(path string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("submodule init %s", path))
}

func (c *GitCommand) SubmoduleUpdate(path string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("submodule update --init %s", path))
}

func (c *GitCommand) SubmoduleBulkInitCmdObj() ICmdObj {
	return BuildGitCmdObjFromStr("submodule init")
}

func (c *GitCommand) SubmoduleBulkUpdateCmdObj() ICmdObj {
	return BuildGitCmdObjFromStr("submodule update")
}

func (c *GitCommand) SubmoduleForceBulkUpdateCmdObj() ICmdObj {
	// not doing an --init here because the user probably doesn't want that
	return BuildGitCmdObjFromStr("submodule update --force")
}

func (c *GitCommand) SubmoduleBulkDeinitCmdObj() ICmdObj {
	return BuildGitCmdObjFromStr("submodule deinit --all --force")
}

func (c *GitCommand) ResetSubmodules(submodules []*models.SubmoduleConfig) error {
	for _, submodule := range submodules {
		if err := c.SubmoduleStash(submodule); err != nil {
			return err
		}
	}

	return c.RunExecutable(c.SubmoduleForceBulkUpdateCmdObj())
}
