package filetree

import (
	"sort"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func BuildTreeFromFiles(files []*models.File) *FileNode {
	root := &FileNode{}

	var curr *FileNode
	for _, file := range files {
		splitPath := split(file.Name)
		curr = root
	outer:
		for i := range splitPath {
			var setFile *models.File
			isFile := i == len(splitPath)-1
			if isFile {
				setFile = file
			}

			path := join(splitPath[:i+1])
			for _, existingChild := range curr.Children {
				if existingChild.Path == path {
					curr = existingChild
					continue outer
				}
			}

			newChild := &FileNode{
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

func BuildFlatTreeFromCommitFiles(files []*models.CommitFile) *CommitFileNode {
	rootAux := BuildTreeFromCommitFiles(files)
	sortedFiles := rootAux.GetLeaves()

	return &CommitFileNode{Children: sortedFiles}
}

func BuildTreeFromCommitFiles(files []*models.CommitFile) *CommitFileNode {
	root := &CommitFileNode{}

	var curr *CommitFileNode
	for _, file := range files {
		splitPath := split(file.Name)
		curr = root
	outer:
		for i := range splitPath {
			var setFile *models.CommitFile
			isFile := i == len(splitPath)-1
			if isFile {
				setFile = file
			}

			path := join(splitPath[:i+1])

			for _, existingChild := range curr.Children {
				if existingChild.Path == path {
					curr = existingChild
					continue outer
				}
			}

			newChild := &CommitFileNode{
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

func BuildFlatTreeFromFiles(files []*models.File) *FileNode {
	rootAux := BuildTreeFromFiles(files)
	sortedFiles := rootAux.GetLeaves()

	// Move merge conflicts to top. This is the one way in which sorting
	// differs between flat mode and tree mode
	sort.SliceStable(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].File != nil && sortedFiles[i].File.HasMergeConflicts && !(sortedFiles[j].File != nil && sortedFiles[j].File.HasMergeConflicts)
	})

	return &FileNode{Children: sortedFiles}
}

func split(str string) []string {
	return strings.Split(str, "/")
}

func join(strs []string) string {
	return strings.Join(strs, "/")
}
