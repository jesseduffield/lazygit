package commands

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
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

	return c.OSCommand.RunCommand("git -C %s stash --include-untracked", submodule.Path)
}

func (c *GitCommand) SubmoduleReset(submodule *models.SubmoduleConfig) error {
	return c.OSCommand.RunCommand("git submodule update --init --force %s", submodule.Name)
}

func (c *GitCommand) SubmoduleUpdateAll() error {
	// not doing an --init here because the user probably doesn't want that
	return c.OSCommand.RunCommand("git submodule update --force")
}

func (c *GitCommand) SubmoduleDelete(submodule *models.SubmoduleConfig) error {
	// based on https://gist.github.com/myusuf3/7f645819ded92bda6677

	if err := c.OSCommand.RunCommand("git submodule deinit --force %s", submodule.Path); err != nil {
		return err
	}

	if err := c.OSCommand.RunCommand("git rm --force %s", submodule.Path); err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join(c.DotGitDir, "modules", submodule.Path))
}

func (c *GitCommand) AddSubmodule(name string, path string, url string) error {
	return c.OSCommand.RunCommand(
		"git submodule add --force --name %s -- %s %s ",
		c.OSCommand.Quote(name),
		c.OSCommand.Quote(url),
		c.OSCommand.Quote(path),
	)
}
