package git_commands

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type StashLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewStashLoader(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *StashLoader {
	return &StashLoader{
		Common: common,
		cmd:    cmd,
	}
}

func (self *StashLoader) GetStashEntries(filterPath string) []*models.StashEntry {
	if filterPath == "" {
		return self.getUnfilteredStashEntries()
	}

	cmdArgs := NewGitCmd("stash").Arg("list", "-z", "--name-only", "--pretty=%ct|%gs").ToArgv()
	rawString, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return self.getUnfilteredStashEntries()
	}
	stashEntries := []*models.StashEntry{}
	var currentStashEntry *models.StashEntry
	lines := utils.SplitNul(rawString)
	isAStash := func(line string) bool { return strings.HasPrefix(line, "stash@{") }
	re := regexp.MustCompile(`stash@\{(\d+)\}`)

outer:
	for i := 0; i < len(lines); i++ {
		if !isAStash(lines[i]) {
			continue
		}
		match := re.FindStringSubmatch(lines[i])
		idx, err := strconv.Atoi(match[1])
		if err != nil {
			return self.getUnfilteredStashEntries()
		}
		currentStashEntry = self.stashEntryFromLine(lines[i], idx)
		for i+1 < len(lines) && !isAStash(lines[i+1]) {
			i++
			if lines[i] == filterPath {
				stashEntries = append(stashEntries, currentStashEntry)
				continue outer
			}
		}
	}
	return stashEntries
}

func (self *StashLoader) getUnfilteredStashEntries() []*models.StashEntry {
	cmdArgs := NewGitCmd("stash").Arg("list", "-z", "--pretty=%ct|%gs").ToArgv()

	rawString, _ := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	return lo.Map(utils.SplitNul(rawString), func(line string, index int) *models.StashEntry {
		return self.stashEntryFromLine(line, index)
	})
}

func (c *StashLoader) stashEntryFromLine(line string, index int) *models.StashEntry {
	model := &models.StashEntry{
		Name:  line,
		Index: index,
	}

	tstr, msg, ok := strings.Cut(line, "|")
	if !ok {
		return model
	}

	t, err := strconv.ParseInt(tstr, 10, 64)
	if err != nil {
		return model
	}

	model.Name = msg
	model.Recency = utils.UnixToTimeAgo(t)

	return model
}
