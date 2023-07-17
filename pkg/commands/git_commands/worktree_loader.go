package git_commands

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/samber/lo"
)

type WorktreeLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewWorktreeLoader(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *WorktreeLoader {
	return &WorktreeLoader{
		Common: common,
		cmd:    cmd,
	}
}

func (self *WorktreeLoader) GetWorktrees() ([]*models.Worktree, error) {
	cmdArgs := NewGitCmd("worktree").Arg("list", "--porcelain", "-z").ToArgv()
	worktreesOutput, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	splitLines := strings.Split(worktreesOutput, "\x00")

	var worktrees []*models.Worktree
	var current *models.Worktree
	for _, splitLine := range splitLines {
		if len(splitLine) == 0 && current != nil {
			worktrees = append(worktrees, current)
			current = nil
			continue
		}
		if strings.HasPrefix(splitLine, "worktree ") {
			path := strings.SplitN(splitLine, " ", 2)[1]
			current = &models.Worktree{
				IsMain: len(worktrees) == 0,
				Path:   path,
			}
		} else if strings.HasPrefix(splitLine, "branch ") {
			branch := strings.SplitN(splitLine, " ", 2)[1]
			current.Branch = strings.TrimPrefix(branch, "refs/heads/")
		}
	}

	names := getUniqueNamesFromPaths(lo.Map(worktrees, func(worktree *models.Worktree, _ int) string {
		return worktree.Path
	}))

	for index, worktree := range worktrees {
		worktree.NameField = names[index]
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// move current worktree to the top
	for i, worktree := range worktrees {
		if EqualPath(worktree.Path, pwd) {
			worktrees = append(worktrees[:i], worktrees[i+1:]...)
			worktrees = append([]*models.Worktree{worktree}, worktrees...)
			break
		}
	}

	// Some worktrees are on a branch but are mid-rebase, and in those cases,
	// `git worktree list` will not show the branch name. We can get the branch
	// name from the `rebase-merge/head-name` file (if it exists) in the folder
	// for the worktree in the parent repo's .git/worktrees folder.
	for _, worktree := range worktrees {
		// No point checking if we already have a branch name
		if worktree.Branch != "" {
			continue
		}

		rebaseBranch, ok := rebaseBranch(worktree.Path)
		if ok {
			worktree.Branch = rebaseBranch
		}
	}

	return worktrees, nil
}

func rebaseBranch(worktreePath string) (string, bool) {
	// need to find the actual path of the worktree in the .git dir
	gitPath, ok := WorktreeGitPath(worktreePath)
	if !ok {
		return "", false
	}

	// now we look inside that git path for a file `rebase-merge/head-name`
	// if it exists, we update the worktree to say that it has that for a head
	headNameContents, err := os.ReadFile(filepath.Join(gitPath, "rebase-merge", "head-name"))
	if err != nil {
		return "", false
	}

	headName := strings.TrimSpace(string(headNameContents))
	shortHeadName := strings.TrimPrefix(headName, "refs/heads/")

	return shortHeadName, true
}

func WorktreeGitPath(worktreePath string) (string, bool) {
	// first we get the path of the worktree, then we look at the contents of the `.git` file in that path
	// then we look for the line that says `gitdir: /path/to/.git/worktrees/<worktree-name>`
	// then we return that path
	gitFileContents, err := os.ReadFile(filepath.Join(worktreePath, ".git"))
	if err != nil {
		return "", false
	}

	gitDirLine := lo.Filter(strings.Split(string(gitFileContents), "\n"), func(line string, _ int) bool {
		return strings.HasPrefix(line, "gitdir: ")
	})

	if len(gitDirLine) == 0 {
		return "", false
	}

	gitDir := strings.TrimPrefix(gitDirLine[0], "gitdir: ")
	return gitDir, true
}

type pathWithIndexT struct {
	path  string
	index int
}

type nameWithIndexT struct {
	name  string
	index int
}

func getUniqueNamesFromPaths(paths []string) []string {
	pathsWithIndex := lo.Map(paths, func(path string, index int) pathWithIndexT {
		return pathWithIndexT{path, index}
	})

	namesWithIndex := getUniqueNamesFromPathsAux(pathsWithIndex, 0)

	// now sort based on index
	result := make([]string, len(namesWithIndex))
	for _, nameWithIndex := range namesWithIndex {
		result[nameWithIndex.index] = nameWithIndex.name
	}

	return result
}

func getUniqueNamesFromPathsAux(paths []pathWithIndexT, depth int) []nameWithIndexT {
	// If we have no paths, return an empty array
	if len(paths) == 0 {
		return []nameWithIndexT{}
	}

	// If we have only one path, return the last segment of the path
	if len(paths) == 1 {
		path := paths[0]
		return []nameWithIndexT{{index: path.index, name: sliceAtDepth(path.path, depth)}}
	}

	// group the paths by their value at the specified depth
	groups := make(map[string][]pathWithIndexT)
	for _, path := range paths {
		value := valueAtDepth(path.path, depth)
		groups[value] = append(groups[value], path)
	}

	result := []nameWithIndexT{}
	for _, group := range groups {
		if len(group) == 1 {
			path := group[0]
			result = append(result, nameWithIndexT{index: path.index, name: sliceAtDepth(path.path, depth)})
		} else {
			result = append(result, getUniqueNamesFromPathsAux(group, depth+1)...)
		}
	}

	return result
}

// if the path is /a/b/c/d, and the depth is 0, the value is 'd'. If the depth is 1, the value is 'c', etc
func valueAtDepth(path string, depth int) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	// Split the path into segments
	segments := strings.Split(path, "/")

	// Get the length of segments
	length := len(segments)

	// If the depth is greater than the length of segments, return an empty string
	if depth >= length {
		return ""
	}

	// Return the segment at the specified depth from the end of the path
	return segments[length-1-depth]
}

// if the path is /a/b/c/d, and the depth is 0, the value is 'd'. If the depth is 1, the value is 'b/c', etc
func sliceAtDepth(path string, depth int) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	// Split the path into segments
	segments := strings.Split(path, "/")

	// Get the length of segments
	length := len(segments)

	// If the depth is greater than or equal to the length of segments, return an empty string
	if depth >= length {
		return ""
	}

	// Join the segments from the specified depth till end of the path
	return strings.Join(segments[length-1-depth:], "/")
}
