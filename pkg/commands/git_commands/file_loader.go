package git_commands

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
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
	// When true, this status is part of an unattended background refresh, so it
	// keeps the default suppression of optional locks (avoiding index.lock
	// contention with git commands the user runs in a terminal, at the cost of
	// not persisting git's refreshed stat-cache). A foreground status opts back
	// in; see gitStatus.
	Background bool
}

func (self *FileLoader) GetStatusFiles(opts GetStatusFileOptions) []*models.File {
	// Decide how to pass --untracked-files to git status.
	//
	// If the user has explicitly set status.showUntrackedFiles, always honor it.
	// Otherwise default to "all" so that individual files inside newly-created
	// untracked directories show up in the files panel. The exception is very
	// large repos: git can only use its untracked cache in "normal" mode, so in
	// "all" mode it does a full recursive readdir of the entire worktree on every
	// status. In a repo with hundreds of thousands of files that takes several
	// seconds and holds index.lock long enough that a concurrent git command
	// (e.g. the checkout that triggered this refresh) can fail with
	// "Unable to create '.../index.lock': File exists". In "normal" mode git uses
	// the untracked cache and stays fast, at the cost of showing a brand-new
	// untracked directory as a single entry rather than listing each file within
	// it. See https://github.com/jesseduffield/lazygit/issues/5906.
	untrackedFilesSetting := self.config.GetShowUntrackedFiles()
	if untrackedFilesSetting == "" {
		untrackedFilesSetting = "all"
		if self.repoTooLargeForUntrackedFilesAll() {
			untrackedFilesSetting = "normal"
		}
	}
	// The "show untracked files" filter is an explicit, transient user request to
	// see untracked files, so honor it even in a large repo.
	if opts.ForceShowUntracked {
		untrackedFilesSetting = "all"
	}
	untrackedFilesArg := fmt.Sprintf("--untracked-files=%s", untrackedFilesSetting)

	statuses, err := self.gitStatus(GitStatusOptions{NoRenames: opts.NoRenames, UntrackedFilesArg: untrackedFilesArg, Background: opts.Background})
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

	self.setConflictMarkerSizes(files)

	return files
}

// Looks up how long the conflict markers in the conflicted files are. We ask
// git for all of them at once, because spawning a process per file would be
// painfully slow when hundreds of files are conflicted (especially on Windows).
func (self *FileLoader) setConflictMarkerSizes(files []*models.File) {
	conflictedFiles := lo.Filter(files, func(file *models.File, _ int) bool {
		return file.HasInlineMergeConflicts
	})
	if len(conflictedFiles) == 0 {
		return
	}

	paths := lo.Map(conflictedFiles, func(file *models.File, _ int) string {
		return file.Path
	})

	markerSizes, err := self.getConflictMarkerSizes(paths)
	if err != nil {
		self.Log.Error(err)
		return
	}

	for _, file := range conflictedFiles {
		file.ConflictMarkerSize = markerSizes[file.Path]
	}
}

func (self *FileLoader) getConflictMarkerSizes(paths []string) (map[string]int, error) {
	cmdArgs := NewGitCmd("check-attr").
		Arg("-z").
		Arg("--stdin").
		Arg("conflict-marker-size").
		ToArgv()

	// -z makes git both read the paths and write its output NUL-separated, so
	// that paths containing newlines don't throw us off.
	output, _, err := self.cmd.New(cmdArgs).
		SetStdin(strings.Join(paths, "\x00")).
		DontLog().
		RunWithOutputs()
	if err != nil {
		return nil, err
	}

	markerSizes := map[string]int{}
	fields := strings.Split(output, "\x00")
	// Each path yields a path/attribute/value triple; the value is either a
	// number or something like "unspecified", in which case we leave the marker
	// size at 0 to say that git's default applies.
	for i := 0; i+2 < len(fields); i += 3 {
		if markerSize, err := strconv.Atoi(fields[i+2]); err == nil && markerSize > 0 {
			markerSizes[fields[i]] = markerSize
		}
	}

	return markerSizes, nil
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

// untrackedFilesAllMaxTrackedFiles is the tracked-file count above which we stop
// defaulting --untracked-files to "all" (see GetStatusFiles): git's untracked
// cache is only used in "normal" mode, so "all" becomes prohibitively slow in
// very large repos. We use the number of entries in the index as a cheap proxy
// for the size of the working tree (and thus the cost of an "all" scan).
const untrackedFilesAllMaxTrackedFiles = 100_000

// repoTooLargeForUntrackedFilesAll reports whether the repository is large enough
// that we should prefer "--untracked-files=normal" over "all" by default. It
// reads the tracked-entry count straight from the index header, which is
// essentially free. If the count can't be determined it returns false, so we
// keep the previous default of "all".
func (self *FileLoader) repoTooLargeForUntrackedFilesAll() bool {
	if self.repoPaths == nil {
		return false
	}
	count, ok := trackedFileCountFromIndex(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "index"))
	return ok && count > untrackedFilesAllMaxTrackedFiles
}

// trackedFileCountFromIndex returns the number of entries recorded in the git
// index at indexPath, read from its 12-byte header: the 4-byte signature "DIRC",
// a 4-byte version, and a 4-byte big-endian entry count.
//
// Limitation: with a split index (core.splitIndex) this reads only the main
// index file, which holds a small delta.
func trackedFileCountFromIndex(indexPath string) (int, bool) {
	f, err := os.Open(indexPath)
	if err != nil {
		return 0, false
	}
	defer f.Close()

	var header [12]byte
	if _, err := io.ReadFull(f, header[:]); err != nil {
		return 0, false
	}
	if string(header[:4]) != "DIRC" {
		return 0, false
	}
	return int(binary.BigEndian.Uint32(header[8:12])), true
}

// GitStatus returns the file status of the repo
type GitStatusOptions struct {
	NoRenames         bool
	UntrackedFilesArg string
	Background        bool
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

	cmdObj := self.cmd.New(cmdArgs).DontLog()
	if !opts.Background {
		// Every git command suppresses optional locks by default (see
		// OptionalLocksEnvVar). A foreground refresh is the one exception: we let
		// it take the lock so it persists git's refreshed stat-cache, which keeps
		// subsequent status calls fast. Background refreshes leave it suppressed so
		// they can't contend for index.lock.
		cmdObj.RemoveEnvVar(OptionalLocksEnvVar)
	}

	statusLines, _, err := cmdObj.RunWithOutputs()
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
