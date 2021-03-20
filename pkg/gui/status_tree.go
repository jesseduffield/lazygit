package gui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/sirupsen/logrus"
)

func GetTreeFromStatusFiles(files []*models.File, log *logrus.Entry) *models.StatusLineNode {
	root := &models.StatusLineNode{}

	var curr *models.StatusLineNode
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

			newChild := &models.StatusLineNode{
				Name: path, // TODO: Remove concept of name
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

func GetFlatTreeFromStatusFiles(files []*models.File) *models.StatusLineNode {
	root := &models.StatusLineNode{}
	for _, file := range files {
		root.Children = append(root.Children, &models.StatusLineNode{
			Name: file.GetPath(),
			Path: file.GetPath(),
			File: file,
		})
	}

	root.Sort()

	return root
}
