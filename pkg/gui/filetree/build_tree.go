package filetree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func BuildTreeFromFiles(files []*models.File) *FileChangeNode {
	root := &FileChangeNode{}

	var curr *FileChangeNode
	for _, file := range files {
		split := strings.Split(file.Name, string(os.PathSeparator))
		curr = root
	outer:
		for i := range split {
			var setFile *models.File
			isFile := i == len(split)-1
			if isFile {
				setFile = file
			}

			path := filepath.Join(split[:i+1]...)

			for _, existingChild := range curr.Children {
				if existingChild.Path == path {
					curr = existingChild
					continue outer
				}
			}

			newChild := &FileChangeNode{
				Path: path,
				File: setFile,
			}
			curr.Children = append(curr.Children, newChild)

			curr = newChild
		}
	}

	root.Sort()
	root.Compress()

	return root
}

func BuildFlatTreeFromFiles(files []*models.File) *FileChangeNode {
	rootAux := BuildTreeFromFiles(files)
	sortedFiles := rootAux.GetLeaves()

	// Move merge conflicts to top. This is the one way in which sorting
	// differs between flat mode and tree mode
	sort.SliceStable(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].File != nil && sortedFiles[i].File.HasMergeConflicts && !(sortedFiles[j].File != nil && sortedFiles[j].File.HasMergeConflicts)
	})

	return &FileChangeNode{Children: sortedFiles}
}
