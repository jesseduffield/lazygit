package git_commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type FileLoaderConfig interface {
	GetShowUntrackedFiles() string
}

type FileLoader struct {
	*GitCommon
	cmd         oscommands.ICmdObjBuilder
	config      FileLoaderConfig
	getFileType func(string) string
}

func NewFileLoader(gitCommon *GitCommon, cmd oscommands.ICmdObjBuilder, config FileLoaderConfig) *FileLoader {
	return &FileLoader{
		GitCommon:   gitCommon,
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

	statuses, err := self.gitStatus(GitStatusOptions{NoRenames: opts.NoRenames, UntrackedFilesArg: untrackedFilesArg})
	if err != nil {
		self.Log.Error(err)
	}
	files := []*models.File{}

	for _, status := range statuses {
		if strings.HasPrefix(status.StatusString, "warning") {
			self.Log.Warningf("warning when calling git status: %s", status.StatusString)
			continue
		}

		file := &models.File{
			Name:          status.Name,
			PreviousName:  status.PreviousName,
			DisplayString: status.StatusString,
		}

		models.SetStatusFields(file, status.Change)
		files = append(files, file)
	}

	// Go through the files to see if any of these files are actually worktrees
	// so that we can render them correctly
	worktreePaths := linkedWortkreePaths(self.Fs, self.repoPaths.RepoGitDirPath())
	for _, file := range files {
		for _, worktreePath := range worktreePaths {
			absFilePath, err := filepath.Abs(file.Name)
			if err != nil {
				self.Log.Error(err)
				continue
			}
			if absFilePath == worktreePath {
				file.IsWorktree = true
				// `git status` renders this worktree as a folder with a trailing slash but we'll represent it as a singular worktree
				// If we include the slash, it will be rendered as a folder with a null file inside.
				file.Name = strings.TrimSuffix(file.Name, "/")
				break
			}
		}
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

func (c *FileLoader) gitStatus(opts GitStatusOptions) ([]FileStatus, error) {
	cmdArgs := NewGitCmd("status").
		Arg(opts.UntrackedFilesArg).
		Arg("--porcelain").
		Arg("-z").
		ArgIf(opts.NoRenames, "--no-renames").
		ToArgv()

	statusLines, _, err := c.cmd.New(cmdArgs).DontLog().RunWithOutputs()
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
