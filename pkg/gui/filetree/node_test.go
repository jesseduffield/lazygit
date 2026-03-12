package filetree

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestGetVisualDepthAtIndex(t *testing.T) {
	scenarios := []struct {
		name           string
		files          []*models.File
		showRootItem   bool
		collapsedPaths []string
		expectedDepths []int // one per visible node, skipping root
	}{
		{
			name: "flat files with root item",
			files: []*models.File{
				{Path: "a"},
				{Path: "b"},
			},
			showRootItem: true,
			// Displayed as:
			//   index 0: ▼ /         (depth 0, the "." root dir)
			//   index 1:   a          (depth 1)
			//   index 2:   b          (depth 1)
			expectedDepths: []int{0, 1, 1},
		},
		{
			name: "flat files without root item",
			files: []*models.File{
				{Path: "a"},
				{Path: "b"},
			},
			showRootItem: false,
			// Displayed as:
			//   index 0: a            (depth 0)
			//   index 1: b            (depth 0)
			expectedDepths: []int{0, 0},
		},
		{
			name: "nested directories with root item",
			files: []*models.File{
				{Path: "dir/a"},
				{Path: "dir/b"},
				{Path: "c"},
			},
			showRootItem: true,
			// Displayed as:
			//   index 0: ▼ /         (depth 0)
			//   index 1:   ▼ dir     (depth 1)
			//   index 2:     a        (depth 2)
			//   index 3:     b        (depth 2)
			//   index 4:   c          (depth 1)
			expectedDepths: []int{0, 1, 2, 2, 1},
		},
		{
			name: "compressed paths with root item",
			files: []*models.File{
				{Path: "dir1/dir3/a"},
				{Path: "dir2/dir4/b"},
			},
			showRootItem: true,
			// Tree compresses dir1/dir3 and dir2/dir4 into single nodes.
			// Displayed as:
			//   index 0: ▼ /            (depth 0)
			//   index 1:   ▼ dir1/dir3  (depth 1, compressed)
			//   index 2:     a           (depth 2)
			//   index 3:   ▼ dir2/dir4  (depth 1, compressed)
			//   index 4:     b           (depth 2)
			expectedDepths: []int{0, 1, 2, 1, 2},
		},
		{
			name: "compressed paths without root item",
			files: []*models.File{
				{Path: "dir1/dir3/a"},
				{Path: "dir2/dir4/b"},
			},
			showRootItem: false,
			// Displayed as:
			//   index 0: ▼ dir1/dir3  (depth 0, compressed)
			//   index 1:   a           (depth 1)
			//   index 2: ▼ dir2/dir4  (depth 0, compressed)
			//   index 3:   b           (depth 1)
			expectedDepths: []int{0, 1, 0, 1},
		},
		{
			name: "collapsed directory hides children",
			files: []*models.File{
				{Path: "dir/a"},
				{Path: "dir/b"},
				{Path: "c"},
			},
			showRootItem:   true,
			collapsedPaths: []string{"./dir"},
			// Displayed as:
			//   index 0: ▼ /         (depth 0)
			//   index 1:   ▶ dir     (depth 1, collapsed)
			//   index 2:   c          (depth 1)
			expectedDepths: []int{0, 1, 1},
		},
		{
			name: "out of range returns -1",
			files: []*models.File{
				{Path: "a"},
			},
			showRootItem:   false,
			expectedDepths: []int{0},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			tree := BuildTreeFromFiles(s.files, s.showRootItem)
			collapsedPaths := NewCollapsedPaths()
			for _, p := range s.collapsedPaths {
				collapsedPaths.Collapse(p)
			}

			for i, expectedDepth := range s.expectedDepths {
				// +1 to skip the invisible root node, matching what FileTree.GetVisualDepth does
				actualDepth := tree.GetVisualDepthAtIndex(i+1, collapsedPaths)
				assert.Equal(t, expectedDepth, actualDepth,
					"index %d: expected depth %d, got %d", i, expectedDepth, actualDepth)
			}

			// Verify out-of-range returns -1
			outOfRange := tree.GetVisualDepthAtIndex(len(s.expectedDepths)+1, collapsedPaths)
			assert.Equal(t, -1, outOfRange, "out of range index should return -1")
		})
	}
}
