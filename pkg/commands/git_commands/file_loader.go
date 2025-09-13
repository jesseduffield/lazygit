package git_commands

import (
	"fmt"
	"path/filepath"
	"strconv"
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
	// If true, we'll show untracked files even if the user has set the config to hide them.
	// This is useful for users with bare repos for dotfiles who default to hiding untracked files,
	// but want to occasionally see them to `git add` a new file.
	ForceShowUntracked bool
}

func (self *FileLoader) GetStatusFiles(opts GetStatusFileOptions) []*models.File {
	// check if config wants us ignoring untracked files
	untrackedFilesSetting := self.config.GetShowUntrackedFiles()

	if opts.ForceShowUntracked || untrackedFilesSetting == "" {
		untrackedFilesSetting = "all"
	}
	untrackedFilesArg := fmt.Sprintf("--untracked-files=%s", untrackedFilesSetting)

	statuses, err := self.gitStatus(GitStatusOptions{NoRenames: opts.NoRenames, UntrackedFilesArg: untrackedFilesArg})
	if err != nil {
		self.Log.Error(err)
	}
	files := []*models.File{}

	fileDiffs := map[string]FileDiff{}
	if self.GitCommon.Common.UserConfig().Gui.ShowNumstatInFilesView {
		fileDiffs, err = self.getFileDiffs()
		if err != nil {
			self.Log.Error(err)
		}
	}

	for _, status := range statuses {
		if strings.HasPrefix(status.StatusString, "warning") {
			self.Log.Warningf("warning when calling git status: %s", status.StatusString)
			continue
		}

		file := &models.File{
			Path:          status.Path,
			PreviousPath:  status.PreviousPath,
			DisplayString: status.StatusString,
		}

		if diff, ok := fileDiffs[status.Path]; ok {
			file.LinesAdded = diff.LinesAdded
			file.LinesDeleted = diff.LinesDeleted
		}

		models.SetStatusFields(file, status.Change)
		files = append(files, file)
	}

	// Go through the files to see if any of these files are actually worktrees
	// so that we can render them correctly
	worktreePaths := linkedWortkreePaths(self.Fs, self.repoPaths.RepoGitDirPath())
	for _, file := range files {
		for _, worktreePath := range worktreePaths {
			absFilePath, err := filepath.Abs(file.Path)
			if err != nil {
				self.Log.Error(err)
				continue
			}
			if absFilePath == worktreePath {
				file.IsWorktree = true
				// `git status` renders this worktree as a folder with a trailing slash but we'll represent it as a singular worktree
				// If we include the slash, it will be rendered as a folder with a null file inside.
				file.Path = strings.TrimSuffix(file.Path, "/")
				break
			}
		}
	}

	return files
}

type FileDiff struct {
	LinesAdded   int
	LinesDeleted int
}

func (self *FileLoader) getFileDiffs() (map[string]FileDiff, error) {
	diffs, err := self.gitDiffNumStat()
	if err != nil {
		return nil, err
	}

	splitLines := strings.Split(diffs, "\x00")

	fileDiffs := map[string]FileDiff{}
	for _, line := range splitLines {
		splitLine := strings.Split(line, "\t")
		if len(splitLine) != 3 {
			continue
		}

		linesAdded, err := strconv.Atoi(splitLine[0])
		if err != nil {
			continue
		}
		linesDeleted, err := strconv.Atoi(splitLine[1])
		if err != nil {
			continue
		}

		fileName := splitLine[2]
		fileDiffs[fileName] = FileDiff{
			LinesAdded:   linesAdded,
			LinesDeleted: linesDeleted,
		}
	}

	return fileDiffs, nil
}

// GitStatus returns the file status of the repo
type GitStatusOptions struct {
	NoRenames         bool
	UntrackedFilesArg string
}

type FileStatus struct {
	StatusString string
	Change       string // ??, MM, AM, ...
	Path         string
	PreviousPath string
}

func (self *FileLoader) gitDiffNumStat() (string, error) {
	return self.cmd.New(
		NewGitCmd("diff").
			Arg("--numstat").
			Arg("-z").
			Arg("HEAD").
			ToArgv(),
	).DontLog().RunWithOutput()
}

func (self *FileLoader) gitStatus(opts GitStatusOptions) ([]FileStatus, error) {
	cmdArgs := NewGitCmd("status").
		Arg(opts.UntrackedFilesArg).
		Arg("--porcelain").
		Arg("-z").
		ArgIfElse(
			opts.NoRenames,
			"--no-renames",
			fmt.Sprintf("--find-renames=%d%%", self.UserConfig().Git.RenameSimilarityThreshold),
		).
		ToArgv()

	statusLines, _, err := self.cmd.New(cmdArgs).DontLog().RunWithOutputs()
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
			Path:         original[3:],
			PreviousPath: "",
		}

		if strings.HasPrefix(status.Change, "R") || strings.HasPrefix(status.Change, "C") {
			// if a line starts with 'R' (rename) or 'C' (copy) then the next line is the original file.
			status.PreviousPath = splitLines[i+1]
			status.StatusString = fmt.Sprintf("%s %s -> %s", status.Change, status.PreviousPath, status.Path)
			i++
		}

		response = append(response, status)
	}

	return response, nil
}
