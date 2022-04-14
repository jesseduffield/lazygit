package loaders

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/samber/lo"
)

type FileLoaderConfig interface {
	GetShowUntrackedFiles() string
}

type FileLoader struct {
	*common.Common
	cmd         oscommands.ICmdObjBuilder
	config      FileLoaderConfig
	getFileType func(string) string
}

func NewFileLoader(cmn *common.Common, cmd oscommands.ICmdObjBuilder, config FileLoaderConfig) *FileLoader {
	return &FileLoader{
		Common:      cmn,
		cmd:         cmd,
		getFileType: oscommands.FileType,
		config:      config,
	}
}

type GetStatusFileOptions struct {
	NoRenames bool
}

func (self *FileLoader) GetStatusFiles(opts GetStatusFileOptions) []*models.File {
	// check if config wants us ignoring untracked files
	untrackedFilesSetting := self.config.GetShowUntrackedFiles()

	if untrackedFilesSetting == "" {
		untrackedFilesSetting = "all"
	}
	untrackedFilesArg := fmt.Sprintf("--untracked-files=%s", untrackedFilesSetting)

	statuses, err := self.GitStatus(GitStatusOptions{NoRenames: opts.NoRenames, UntrackedFilesArg: untrackedFilesArg})
	if err != nil {
		self.Log.Error(err)
	}
	files := []*models.File{}

	for _, status := range statuses {
		if strings.HasPrefix(status.StatusString, "warning") {
			self.Log.Warningf("warning when calling git status: %s", status.StatusString)
			continue
		}
		change := status.Change
		stagedChange := change[0:1]
		unstagedChange := change[1:2]
		untracked := lo.Contains([]string{"??", "A ", "AM"}, change)
		hasNoStagedChanges := lo.Contains([]string{" ", "U", "?"}, stagedChange)
		hasInlineMergeConflicts := lo.Contains([]string{"UU", "AA"}, change)
		hasMergeConflicts := hasInlineMergeConflicts || lo.Contains([]string{"DD", "AU", "UA", "UD", "DU"}, change)

		file := &models.File{
			Name:                    status.Name,
			PreviousName:            status.PreviousName,
			DisplayString:           status.StatusString,
			HasStagedChanges:        !hasNoStagedChanges,
			HasUnstagedChanges:      unstagedChange != " ",
			Tracked:                 !untracked,
			Deleted:                 unstagedChange == "D" || stagedChange == "D",
			Added:                   unstagedChange == "A" || untracked,
			HasMergeConflicts:       hasMergeConflicts,
			HasInlineMergeConflicts: hasInlineMergeConflicts,
			Type:                    self.getFileType(status.Name),
			ShortStatus:             change,
		}
		files = append(files, file)
	}

	return files
}

// GitStatus returns the file status of the repo
type GitStatusOptions struct {
	NoRenames         bool
	UntrackedFilesArg string
}

type FileStatus struct {
	StatusString string
	Change       string // ??, MM, AM, ...
	Name         string
	PreviousName string
}

func (c *FileLoader) GitStatus(opts GitStatusOptions) ([]FileStatus, error) {
	noRenamesFlag := ""
	if opts.NoRenames {
		noRenamesFlag = " --no-renames"
	}

	statusLines, err := c.cmd.New(fmt.Sprintf("git status %s --porcelain -z%s", opts.UntrackedFilesArg, noRenamesFlag)).DontLog().RunWithOutput()
	if err != nil {
		return []FileStatus{}, err
	}

	splitLines := strings.Split(statusLines, "\x00")
	response := []FileStatus{}

	for i := 0; i < len(splitLines); i++ {
		original := splitLines[i]

		if len(original) < 3 {
			continue
		}

		status := FileStatus{
			StatusString: original,
			Change:       original[:2],
			Name:         original[3:],
			PreviousName: "",
		}

		if strings.HasPrefix(status.Change, "R") {
			// if a line starts with 'R' then the next line is the original file.
			status.PreviousName = splitLines[i+1]
			status.StatusString = fmt.Sprintf("%s %s -> %s", status.Change, status.PreviousName, status.Name)
			i++
		}

		response = append(response, status)
	}

	return response, nil
}
