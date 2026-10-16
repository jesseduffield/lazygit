package git_commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// DiffInputs is what the diffs of some files of the working tree are made from, as far
// as it can be told without reading the files: each file's entries in the index, the
// commit HEAD is at for the staged side of the diffs, and the type, permissions, size
// and modification time of each file in the working tree for the unstaged side.
//
// Taken again later and compared, it tells whether the files' diffs may have changed in
// between. Some changes alter the inputs without changing a diff: HEAD moving to a
// commit with the same version of a file, or a file written with the content it has.
// And a rewrite of a file that keeps its size and modification time goes unnoticed.
type DiffInputs struct {
	// The commit HEAD is at, if taken for the staged side; "" in a repo without commits.
	head string

	// Each file's entries in the index, as `git ls-files -s` prints them, by repo-relative
	// path; "" for a file that has none.
	indexEntries map[string]string

	// What a stat of each file in the working tree says, by repo-relative path, if taken
	// for the unstaged side; "" for a file that doesn't exist.
	files map[string]string
}

// GetDiffInputs takes the DiffInputs of the files at the given repo-relative paths, for
// the staged side of their diffs, the unstaged side, or both.
func (self *WorkingTreeCommands) GetDiffInputs(paths []string, stagedSide bool, unstagedSide bool) (*DiffInputs, error) {
	inputs := &DiffInputs{indexEntries: map[string]string{}}

	for _, path := range paths {
		inputs.indexEntries[path] = ""
	}
	for _, batch := range pathsPerCommandLine(paths) {
		output, err := self.cmd.New(
			NewGitCmd("ls-files").Arg("-s", "-z", "--").Arg(batch...).Dir(self.repoPaths.worktreePath).ToArgv(),
		).DontLog().RunWithOutput()
		if err != nil {
			return nil, err
		}
		// An entry is "<mode> <blob> <stage>\t<path>"; a file with a merge conflict has
		// one for each stage.
		for entry := range strings.SplitSeq(output, "\x00") {
			if _, path, ok := strings.Cut(entry, "\t"); ok {
				inputs.indexEntries[path] += entry + "\n"
			}
		}
	}

	if stagedSide {
		// In a repo without commits there is no HEAD to verify, and git prints nothing.
		// That is a state like any other, so the error is no reason to fail.
		head, _ := self.cmd.New(NewGitCmd("rev-parse").Arg("--verify", "--quiet", "HEAD").ToArgv()).
			DontLog().RunWithOutput()
		inputs.head = strings.TrimSpace(head)
	}

	if unstagedSide {
		inputs.files = map[string]string{}
		for _, path := range paths {
			inputs.files[path] = self.statWorkingTreeFile(filepath.Join(self.repoPaths.worktreePath, path))
		}
	}

	return inputs, nil
}

// statWorkingTreeFile describes a file of the working tree by what a stat of it says,
// without following a symlink: git diffs a symlink as the path it points to. It is ""
// for a file that doesn't exist.
func (self *WorkingTreeCommands) statWorkingTreeFile(path string) string {
	var info os.FileInfo
	var err error
	if lstater, ok := self.Fs.(afero.Lstater); ok {
		info, _, err = lstater.LstatIfPossible(path)
	} else {
		info, err = self.Fs.Stat(path)
	}
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%v %d %d", info.Mode(), info.Size(), info.ModTime().UnixNano())
}

// SameFor reports whether the inputs of the diffs of the files at the given
// repo-relative paths, on the given side, are the same in both.
func (self *DiffInputs) SameFor(other *DiffInputs, paths []string, stagedSide bool) bool {
	for _, path := range paths {
		if self.indexEntries[path] != other.indexEntries[path] {
			return false
		}
		if !stagedSide && self.files[path] != other.files[path] {
			return false
		}
	}
	return !stagedSide || self.head == other.head
}

// Paths returns the repo-relative paths of the files whose inputs were taken.
func (self *DiffInputs) Paths() []string {
	paths := make([]string, 0, len(self.indexEntries))
	for path := range self.indexEntries {
		paths = append(paths, path)
	}
	return paths
}
