package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

const RENAME_SEPARATOR = " -> "

// LoadStatusFiles git status files
type LoadStatusFilesOpts struct {
	NoRenames bool
}

type StatusFilesLoader struct {
	ICommander

	config IGitConfigMgr
	log    *logrus.Entry
	os     oscommands.IOS
}

func NewStatusFilesLoader(commander ICommander, config IGitConfigMgr, log *logrus.Entry, os oscommands.IOS) *StatusFilesLoader {
	return &StatusFilesLoader{
		ICommander: commander,
		config:     config,
		log:        log,
		os:         os,
	}
}

func (c *StatusFilesLoader) Load(opts LoadStatusFilesOpts) []*models.File {
	cmdObj := c.buildCmdObj(opts)

	status, err := c.RunWithOutput(cmdObj)
	if err != nil {
		c.log.Error(err)
		return []*models.File{}
	}

	statusStrings := c.cleanGitStatus(status)
	files := make([]*models.File, 0, len(statusStrings))

	for _, statusString := range statusStrings {
		file := c.fileFromStatusString(statusString)
		if file == nil {
			continue
		}

		files = append(files, file)
	}

	return files
}

func (c *StatusFilesLoader) fileFromStatusString(statusString string) *models.File {
	if strings.HasPrefix(statusString, "warning") {
		c.log.Warningf("warning when calling git status: %s", statusString)
		return nil
	}
	change := statusString[0:2]
	stagedChange := change[0:1]
	unstagedChange := statusString[1:2]
	name := statusString[3:]
	untracked := utils.IncludesString([]string{"??", "A ", "AM"}, change)
	hasNoStagedChanges := utils.IncludesString([]string{" ", "U", "?"}, stagedChange)
	hasMergeConflicts := utils.IncludesString([]string{"DD", "AA", "UU", "AU", "UA", "UD", "DU"}, change)
	hasInlineMergeConflicts := utils.IncludesString([]string{"UU", "AA"}, change)
	previousName := ""
	if strings.Contains(name, RENAME_SEPARATOR) {
		split := strings.Split(name, RENAME_SEPARATOR)
		name = split[1]
		previousName = split[0]
	}

	return &models.File{
		Name:                    name,
		PreviousName:            previousName,
		DisplayString:           statusString,
		HasStagedChanges:        !hasNoStagedChanges,
		HasUnstagedChanges:      unstagedChange != " ",
		Tracked:                 !untracked,
		Deleted:                 unstagedChange == "D" || stagedChange == "D",
		Added:                   unstagedChange == "A" || untracked,
		HasMergeConflicts:       hasMergeConflicts,
		HasInlineMergeConflicts: hasInlineMergeConflicts,
		Type:                    c.os.FileType(name),
		ShortStatus:             change,
	}
}

func (c *StatusFilesLoader) buildCmdObj(opts LoadStatusFilesOpts) ICmdObj {
	// check if config wants us ignoring untracked files
	untrackedFilesSetting := c.config.GetConfigValue("status.showUntrackedFiles")

	if untrackedFilesSetting == "" {
		untrackedFilesSetting = "all"
	}
	untrackedFilesArg := fmt.Sprintf("--untracked-files=%s", untrackedFilesSetting)

	noRenamesFlag := ""
	if opts.NoRenames {
		noRenamesFlag = "--no-renames"
	}

	return c.BuildGitCmdObjFromStr(fmt.Sprintf("status %s --porcelain -z %s", untrackedFilesArg, noRenamesFlag))
}

func (*StatusFilesLoader) cleanGitStatus(statusLines string) []string {
	splitLines := strings.Split(statusLines, "\x00")
	// if a line starts with 'R' then the next line is the original file.
	for i := 0; i < len(splitLines)-1; i++ {
		original := splitLines[i]
		if strings.HasPrefix(original, "R  ") {
			next := splitLines[i+1]
			updated := "R  " + next + RENAME_SEPARATOR + strings.TrimPrefix(original, "R  ")
			splitLines[i] = updated
			splitLines = append(splitLines[0:i+1], splitLines[i+2:]...)
		}
	}

	return utils.SplitLines(strings.Join(splitLines, "\n"))
}
