package commands

import (
	"bufio"
	"os"
	"regexp"
)

// .gitmodules looks like this:
// [submodule "mysubmodule"]
//   path = blah/mysubmodule
//   url = git@github.com:subbo.git

func (c *GitCommand) GetSubmoduleConfigs() ([]*SubmoduleConfig, error) {
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

	configs := []*SubmoduleConfig{}
	for scanner.Scan() {
		line := scanner.Text()

		if name, ok := firstMatch(line, `\[submodule "(.*)"\]`); ok {
			configs = append(configs, &SubmoduleConfig{Name: name})
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

func (c *GitCommand) SubmoduleStash(config *SubmoduleConfig) error {
	return c.OSCommand.RunCommand("git -C %s stash --include-untracked", config.Path)
}

func (c *GitCommand) SubmoduleReset(config *SubmoduleConfig) error {
	return c.OSCommand.RunCommand("git submodule update --force %s", config.Name)
}
