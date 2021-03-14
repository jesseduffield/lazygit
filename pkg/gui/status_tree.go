package gui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func GetTreeFromStatusFiles(files []*models.File) *models.StatusLineNode {
	root := &models.StatusLineNode{}

	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	var curr *models.StatusLineNode
	for _, file := range files {
		split := strings.Split(file.Name, string(os.PathSeparator))
		curr = root
	outer:
		for i, dir := range split {
			var setFile *models.File
			if i == len(split)-1 {
				setFile = file
			}
			for _, existingChild := range curr.Children {
				if existingChild.Name == dir {
					curr = existingChild
					continue outer
				}
			}
			newChild := &models.StatusLineNode{
				Name: dir,
				Path: filepath.Join(split[:i+1]...),
				File: setFile,
			}
			curr.Children = append(curr.Children, newChild)

			curr = newChild
		}
	}

	root.SortTree()

	return root
}
