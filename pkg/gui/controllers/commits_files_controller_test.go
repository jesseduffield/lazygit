package controllers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func Test_normalisedSelectedCommitFileNodes(t *testing.T) {
	scenarios := []struct {
		name          string
		files         []string
		selectedPaths []string
		expectedPaths []string
	}{
		{
			name:          "sibling directories whose names share a prefix are both kept",
			files:         []string{"foo/file", "foobar/file"},
			selectedPaths: []string{"foo", "foobar"},
			expectedPaths: []string{"foo", "foobar"},
		},
		{
			name:          "sibling files whose names share a prefix are both kept",
			files:         []string{"a.txt", "a.txt.orig"},
			selectedPaths: []string{"a.txt", "a.txt.orig"},
			expectedPaths: []string{"a.txt", "a.txt.orig"},
		},
		{
			name:          "a file inside a selected directory is dropped",
			files:         []string{"foo/file", "foobar/file"},
			selectedPaths: []string{"foo", "foo/file"},
			expectedPaths: []string{"foo"},
		},
		{
			name:          "a directory inside a selected directory is dropped",
			files:         []string{"foo/bar/file", "foo/baz/file", "foobar/file"},
			selectedPaths: []string{"foo", "foo/bar", "foobar"},
			expectedPaths: []string{"foo", "foobar"},
		},
		{
			name:          "the root item drops everything below it",
			files:         []string{"foo/file", "foobar/file"},
			selectedPaths: []string{".", "foo", "foobar/file"},
			expectedPaths: []string{"."},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			nodesByPath := commitFileNodesByPath(s.files)

			selectedNodes := lo.Map(s.selectedPaths, func(path string, _ int) *filetree.CommitFileNode {
				node, ok := nodesByPath[path]
				assert.True(t, ok, "no node for path %q", path)
				return node
			})

			actualPaths := lo.Map(normalisedSelectedCommitFileNodes(selectedNodes),
				func(node *filetree.CommitFileNode, _ int) string { return node.GetPath() })

			assert.Equal(t, s.expectedPaths, actualPaths)
		})
	}
}

// builds a commit file tree from the given paths and indexes every node by its
// logical path, so that scenarios can name the nodes they select
func commitFileNodesByPath(paths []string) map[string]*filetree.CommitFileNode {
	files := lo.Map(paths, func(path string, _ int) *models.CommitFile {
		return &models.CommitFile{Path: path}
	})

	root := filetree.BuildTreeFromCommitFiles(files, true,
		filetree.NodeSortComparator[models.CommitFile]("mixed", false))

	nodesByPath := map[string]*filetree.CommitFileNode{}
	var collect func(node *filetree.Node[models.CommitFile])
	collect = func(node *filetree.Node[models.CommitFile]) {
		if node.GetInternalPath() != "" {
			nodesByPath[node.GetPath()] = filetree.NewCommitFileNode(node)
		}
		for _, child := range node.Children {
			collect(child)
		}
	}
	collect(root)

	return nodesByPath
}
