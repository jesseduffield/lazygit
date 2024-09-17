package filetree

import (
	"sort"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func BuildTreeFromFiles(files []*models.File) *Node[models.File] {
	root := &Node[models.File]{}

	childrenMapsByNode := make(map[*Node[models.File]]map[string]*Node[models.File])

	var curr *Node[models.File]
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

			var currNodeChildrenMap map[string]*Node[models.File]
			var isCurrNodeMapped bool

			if currNodeChildrenMap, isCurrNodeMapped = childrenMapsByNode[curr]; !isCurrNodeMapped {
				currNodeChildrenMap = make(map[string]*Node[models.File])
				childrenMapsByNode[curr] = currNodeChildrenMap
			}

			child, doesCurrNodeHaveChildAlready := currNodeChildrenMap[path]
			if doesCurrNodeHaveChildAlready {
				curr = child
				continue outer
			}

			newChild := &Node[models.File]{
				Path: path,
				File: setFile,
			}
			curr.Children = append(curr.Children, newChild)

			currNodeChildrenMap[path] = newChild

			curr = newChild
		}
	}

	root.Sort()
	root.Compress()

	return root
}

func BuildFlatTreeFromCommitFiles(files []*models.CommitFile) *Node[models.CommitFile] {
	rootAux := BuildTreeFromCommitFiles(files)
	sortedFiles := rootAux.GetLeaves()

	return &Node[models.CommitFile]{Children: sortedFiles}
}

func BuildTreeFromCommitFiles(files []*models.CommitFile) *Node[models.CommitFile] {
	root := &Node[models.CommitFile]{}

	var curr *Node[models.CommitFile]
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

			newChild := &Node[models.CommitFile]{
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

func BuildFlatTreeFromFiles(files []*models.File) *Node[models.File] {
	rootAux := BuildTreeFromFiles(files)
	sortedFiles := rootAux.GetLeaves()

	// from top down we have merge conflict files, then tracked file, then untracked
	// files. This is the one way in which sorting differs between flat mode and
	// tree mode
	sort.SliceStable(sortedFiles, func(i, j int) bool {
		iFile := sortedFiles[i].File
		jFile := sortedFiles[j].File

		// never going to happen but just to be safe
		if iFile == nil || jFile == nil {
			return false
		}

		if iFile.HasMergeConflicts && !jFile.HasMergeConflicts {
			return true
		}

		if jFile.HasMergeConflicts && !iFile.HasMergeConflicts {
			return false
		}

		if iFile.Tracked && !jFile.Tracked {
			return true
		}

		if jFile.Tracked && !iFile.Tracked {
			return false
		}

		return false
	})

	return &Node[models.File]{Children: sortedFiles}
}

func split(str string) []string {
	return strings.Split(str, "/")
}

func join(strs []string) string {
	return strings.Join(strs, "/")
}
