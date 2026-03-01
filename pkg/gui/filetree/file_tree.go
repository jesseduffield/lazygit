package filetree

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

type FileTreeDisplayFilter int

const (
	DisplayAll FileTreeDisplayFilter = iota
	DisplayStaged
	DisplayUnstaged
	DisplayTracked
	DisplayUntracked
	// this shows files with merge conflicts
	DisplayConflicted
)

type ITree[T any] interface {
	InTreeMode() bool
	ExpandToPath(path string)
	ToggleShowTree()
	GetIndexForPath(path string) (int, bool)
	Len() int
	GetItem(index int) types.HasUrn
	SetTree()
	IsCollapsed(path string) bool
	ToggleCollapsed(path string)
	CollapsedPaths() *CollapsedPaths
	CollapseAll()
	ExpandAll()
}

type IFileTree interface {
	ITree[models.File]

	FilterFiles(test func(*models.File) bool) []*models.File
	SetStatusFilter(filter FileTreeDisplayFilter)
	SetPathFilter(filter string, useFuzzySearch bool)
	ReApplyPathFilter(useFuzzySearch bool)
	EnsurePathFilterInitialized(defaultFilter string, useFuzzySearch bool)
	ForceShowUntracked() bool
	Get(index int) *FileNode
	GetFile(path string) *models.File
	GetAllItems() []*FileNode
	GetAllFiles() []*models.File
	GetStatusFilter() FileTreeDisplayFilter
	GetPathFilter() string
	GetRoot() *FileNode
}

type FileTree struct {
	getFiles              func() []*models.File
	tree                  *Node[models.File]
	showTree              bool
	common                *common.Common
	statusFilter          FileTreeDisplayFilter
	pathFilter            string
	useFuzzyPathFilter    bool
	pathFilterInitialized bool
	collapsedPaths        *CollapsedPaths
}

var _ IFileTree = &FileTree{}

func NewFileTree(getFiles func() []*models.File, common *common.Common, showTree bool) *FileTree {
	return &FileTree{
		getFiles:       getFiles,
		common:         common,
		showTree:       showTree,
		statusFilter:   DisplayAll,
		collapsedPaths: NewCollapsedPaths(),
	}
}

func (self *FileTree) InTreeMode() bool {
	return self.showTree
}

func (self *FileTree) ExpandToPath(path string) {
	self.collapsedPaths.ExpandToPath(path)
}

func (self *FileTree) getFilesForDisplay() []*models.File {
	var files []*models.File

	switch self.statusFilter {
	case DisplayAll:
		files = self.getFiles()
	case DisplayStaged:
		files = self.FilterFiles(func(file *models.File) bool { return file.HasStagedChanges })
	case DisplayUnstaged:
		files = self.FilterFiles(func(file *models.File) bool { return file.HasUnstagedChanges })
	case DisplayTracked:
		// untracked but staged files are technically not tracked by git
		// but including such files in the filtered mode helps see what files are getting committed
		files = self.FilterFiles(func(file *models.File) bool { return file.Tracked || file.HasStagedChanges })
	case DisplayUntracked:
		files = self.FilterFiles(func(file *models.File) bool { return !(file.Tracked || file.HasStagedChanges) })
	case DisplayConflicted:
		files = self.FilterFiles(func(file *models.File) bool { return file.HasMergeConflicts })
	default:
		panic(fmt.Sprintf("Unexpected files display filter: %d", self.statusFilter))
	}

	return self.pathFilteredFiles(files)
}

func (self *FileTree) ForceShowUntracked() bool {
	return self.statusFilter == DisplayUntracked
}

func (self *FileTree) FilterFiles(test func(*models.File) bool) []*models.File {
	return lo.Filter(self.getFiles(), func(file *models.File, _ int) bool { return test(file) })
}

func (self *FileTree) SetStatusFilter(filter FileTreeDisplayFilter) {
	self.statusFilter = filter
	self.SetTree()
}

func (self *FileTree) SetPathFilter(filter string, useFuzzySearch bool) {
	self.pathFilter = filter
	self.useFuzzyPathFilter = useFuzzySearch
	self.SetTree()
}

func (self *FileTree) ReApplyPathFilter(useFuzzySearch bool) {
	self.useFuzzyPathFilter = useFuzzySearch
	self.SetTree()
}

func (self *FileTree) EnsurePathFilterInitialized(defaultFilter string, useFuzzySearch bool) {
	if self.pathFilterInitialized {
		return
	}

	self.pathFilterInitialized = true
	if defaultFilter == "" {
		return
	}

	self.pathFilter = defaultFilter
	self.useFuzzyPathFilter = useFuzzySearch
	self.SetTree()
}

func (self *FileTree) ToggleShowTree() {
	self.showTree = !self.showTree
	self.SetTree()
}

func (self *FileTree) Get(index int) *FileNode {
	// need to traverse the tree depth first until we get to the index.
	return NewFileNode(self.tree.GetNodeAtIndex(index+1, self.collapsedPaths)) // ignoring root
}

func (self *FileTree) GetFile(path string) *models.File {
	for _, file := range self.getFiles() {
		if file.Path == path {
			return file
		}
	}

	return nil
}

func (self *FileTree) GetIndexForPath(path string) (int, bool) {
	index, found := self.tree.GetIndexForPath(path, self.collapsedPaths)
	return index - 1, found
}

// note: this gets all items when the filter is taken into consideration. There may
// be hidden files that aren't included here. Files off the screen however will
// be included
func (self *FileTree) GetAllItems() []*FileNode {
	if self.tree == nil {
		return nil
	}

	// ignoring root
	return lo.Map(self.tree.Flatten(self.collapsedPaths)[1:], func(node *Node[models.File], _ int) *FileNode {
		return NewFileNode(node)
	})
}

func (self *FileTree) Len() int {
	// -1 because we're ignoring the root
	return max(self.tree.Size(self.collapsedPaths)-1, 0)
}

func (self *FileTree) GetItem(index int) types.HasUrn {
	// Unimplemented because we don't yet need to show inlines statuses in commit file views
	return nil
}

func (self *FileTree) GetAllFiles() []*models.File {
	return self.getFiles()
}

func (self *FileTree) SetTree() {
	filesForDisplay := self.getFilesForDisplay()
	showRootItem := self.common.UserConfig().Gui.ShowRootItemInFileTree
	if self.showTree {
		self.tree = BuildTreeFromFiles(filesForDisplay, showRootItem)
	} else {
		self.tree = BuildFlatTreeFromFiles(filesForDisplay, showRootItem)
	}
}

func (self *FileTree) IsCollapsed(path string) bool {
	return self.collapsedPaths.IsCollapsed(path)
}

func (self *FileTree) ToggleCollapsed(path string) {
	self.collapsedPaths.ToggleCollapsed(path)
}

func (self *FileTree) CollapseAll() {
	dirPaths := lo.FilterMap(self.GetAllItems(), func(file *FileNode, index int) (string, bool) {
		return file.path, !file.IsFile()
	})

	for _, path := range dirPaths {
		self.collapsedPaths.Collapse(path)
	}
}

func (self *FileTree) ExpandAll() {
	self.collapsedPaths.ExpandAll()
}

func (self *FileTree) Tree() *FileNode {
	return NewFileNode(self.tree)
}

func (self *FileTree) GetRoot() *FileNode {
	return NewFileNode(self.tree)
}

func (self *FileTree) CollapsedPaths() *CollapsedPaths {
	return self.collapsedPaths
}

type fileFuzzySource struct {
	files []*models.File
}

func (self *fileFuzzySource) String(i int) string {
	fields := []string{self.files[i].Path}
	if self.files[i].PreviousPath != "" {
		fields = append(fields, self.files[i].PreviousPath)
	}
	return strings.Join(fields, " ")
}

func (self *fileFuzzySource) Len() int {
	return len(self.files)
}

func (self *FileTree) pathFilteredFiles(files []*models.File) []*models.File {
	if self.pathFilter == "" {
		return files
	}

	source := &fileFuzzySource{files: files}
	matches := utils.FindFrom(self.pathFilter, source, self.useFuzzyPathFilter)
	indices := lo.Map(matches, func(match fuzzy.Match, _ int) int {
		return match.Index
	})
	return utils.ValuesAtIndices(files, indices)
}

func (self *FileTree) GetStatusFilter() FileTreeDisplayFilter {
	return self.statusFilter
}

func (self *FileTree) GetPathFilter() string {
	return self.pathFilter
}
