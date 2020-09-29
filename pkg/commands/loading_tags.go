package commands

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

const semverRegex = `v?((\d+\.?)+)([^\d]?.*)`

func convertToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func (c *GitCommand) GetTags() ([]*models.Tag, error) {
	// get remote branches
	remoteBranchesStr, err := c.OSCommand.RunCommandWithOutput(`git tag --list`)
	if err != nil {
		return nil, err
	}

	content := utils.TrimTrailingNewline(remoteBranchesStr)
	if content == "" {
		return nil, nil
	}

	split := strings.Split(content, "\n")

	// first step is to get our remotes from go-git
	tags := make([]*models.Tag, len(split))
	for i, tagName := range split {

		tags[i] = &models.Tag{
			Name: tagName,
		}
	}

	// now lets sort our tags by name numerically
	re := regexp.MustCompile(semverRegex)

	// the reason  this is complicated is because we're both sorting alphabetically
	// and when we're dealing with semver strings
	sort.Slice(tags, func(i, j int) bool {
		a := tags[i].Name
		b := tags[j].Name

		matchA := re.FindStringSubmatch(a)
		matchB := re.FindStringSubmatch(b)

		if len(matchA) > 0 && len(matchB) > 0 {
			numbersA := strings.Split(matchA[1], ".")
			numbersB := strings.Split(matchB[1], ".")
			k := 0
			for {
				if len(numbersA) == k && len(numbersB) == k {
					break
				}
				if len(numbersA) == k {
					return true
				}
				if len(numbersB) == k {
					return false
				}
				if convertToInt(numbersA[k]) < convertToInt(numbersB[k]) {
					return true
				}
				if convertToInt(numbersA[k]) > convertToInt(numbersB[k]) {
					return false
				}
				k++
			}

			return strings.ToLower(matchA[3]) < strings.ToLower(matchB[3])
		}

		return strings.ToLower(a) < strings.ToLower(b)
	})

	return tags, nil
}
