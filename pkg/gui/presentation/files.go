package presentation

import (
	"strings"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

const (
	EXPANDED_ARROW  = "▼"
	COLLAPSED_ARROW = "▶"
)

func RenderFileTree(
	tree filetree.IFileTree,
	submoduleConfigs []*models.SubmoduleConfig,
	showFileIcons bool,
) []string {
	collapsedPaths := tree.CollapsedPaths()
	return renderAux(tree.GetRoot().Raw(), collapsedPaths, -1, -1, func(node *filetree.Node[models.File], treeDepth int, visualDepth int, isCollapsed bool) string {
		fileNode := filetree.NewFileNode(node)

		return getFileLine(isCollapsed, fileNode.GetHasUnstagedChanges(), fileNode.GetHasStagedChanges(), treeDepth, visualDepth, showFileIcons, submoduleConfigs, node)
	})
}

func RenderCommitFileTree(
	tree *filetree.CommitFileTreeViewModel,
	patchBuilder *patch.PatchBuilder,
	showFileIcons bool,
) []string {
	collapsedPaths := tree.CollapsedPaths()
	return renderAux(tree.GetRoot().Raw(), collapsedPaths, -1, -1, func(node *filetree.Node[models.CommitFile], treeDepth int, visualDepth int, isCollapsed bool) string {
		status := commitFilePatchStatus(node, tree, patchBuilder)

		return getCommitFileLine(isCollapsed, treeDepth, visualDepth, node, status, showFileIcons)
	})
}

// Returns the status of a commit file in terms of its inclusion in the custom patch
func commitFilePatchStatus(node *filetree.Node[models.CommitFile], tree *filetree.CommitFileTreeViewModel, patchBuilder *patch.PatchBuilder) patch.PatchStatus {
	// This is a little convoluted because we're dealing with either a leaf or a non-leaf.
	// But this code actually applies to both. If it's a leaf, the status will just
	// be whatever status it is, but if it's a non-leaf it will determine its status
	// based on the leaves of that subtree
	if node.EveryFile(func(file *models.CommitFile) bool {
		return patchBuilder.GetFileStatus(file.Name, tree.GetRef().RefName()) == patch.WHOLE
	}) {
		return patch.WHOLE
	} else if node.EveryFile(func(file *models.CommitFile) bool {
		return patchBuilder.GetFileStatus(file.Name, tree.GetRef().RefName()) == patch.UNSELECTED
	}) {
		return patch.UNSELECTED
	} else {
		return patch.PART
	}
}

func renderAux[T any](
	node *filetree.Node[T],
	collapsedPaths *filetree.CollapsedPaths,
	// treeDepth is the depth of the node in the actual file tree. This is different to
	// visualDepth because some directory nodes are compressed e.g. 'pkg/gui/blah' takes
	// up two tree depths, but one visual depth. We need to track these separately,
	// because indentation relies on visual depth, whereas file path truncation
	// relies on tree depth.
	treeDepth int,
	visualDepth int,
	renderLine func(*filetree.Node[T], int, int, bool) string,
) []string {
	if node == nil {
		return []string{}
	}

	isRoot := treeDepth == -1

	if node.IsFile() {
		if isRoot {
			return []string{}
		}
		return []string{renderLine(node, treeDepth, visualDepth, false)}
	}

	arr := []string{}
	if !isRoot {
		isCollapsed := collapsedPaths.IsCollapsed(node.GetPath())
		arr = append(arr, renderLine(node, treeDepth, visualDepth, isCollapsed))
	}

	if collapsedPaths.IsCollapsed(node.GetPath()) {
		return arr
	}

	for _, child := range node.Children {
		arr = append(arr, renderAux(child, collapsedPaths, treeDepth+1+node.CompressionLevel, visualDepth+1, renderLine)...)
	}

	return arr
}

func getFileLine(
	isCollapsed bool,
	hasUnstagedChanges bool,
	hasStagedChanges bool,
	treeDepth int,
	visualDepth int,
	showFileIcons bool,
	submoduleConfigs []*models.SubmoduleConfig,
	node *filetree.Node[models.File],
) string {
	name := fileNameAtDepth(node, treeDepth)
	output := ""

	var nameColor style.TextStyle

	file := node.File

	indentation := strings.Repeat("  ", visualDepth)

	if hasStagedChanges && !hasUnstagedChanges {
		nameColor = style.FgGreen
	} else if hasStagedChanges {
		nameColor = style.FgYellow
	} else {
		nameColor = theme.DefaultTextColor
	}

	if file == nil {
		output += indentation + ""
		arrow := EXPANDED_ARROW
		if isCollapsed {
			arrow = COLLAPSED_ARROW
		}

		arrowStyle := nameColor

		output += arrowStyle.Sprint(arrow) + " "
	} else {
		// Sprinting the space at the end in the specific style is for the sake of
		// when a reverse style is used in the theme, which looks ugly if you just
		// use the default style
		output += indentation + formatFileStatus(file, nameColor) + nameColor.Sprint(" ")
	}

	isSubmodule := file != nil && file.IsSubmodule(submoduleConfigs)
	isLinkedWorktree := file != nil && file.IsWorktree
	isDirectory := file == nil

	if showFileIcons {
		icon := icons.IconForFile(name, isSubmodule, isLinkedWorktree, isDirectory)
		paint := color.C256(icon.Color, false)
		output += paint.Sprint(icon.Icon) + nameColor.Sprint(" ")
	}

	output += nameColor.Sprint(utils.EscapeSpecialChars(name))

	if isSubmodule {
		output += theme.DefaultTextColor.Sprint(" (submodule)")
	}

	return output
}

func formatFileStatus(file *models.File, restColor style.TextStyle) string {
	firstChar := file.ShortStatus[0:1]
	firstCharCl := style.FgGreen
	if firstChar == "?" {
		firstCharCl = theme.UnstagedChangesColor
	} else if firstChar == " " {
		firstCharCl = restColor
	}

	secondChar := file.ShortStatus[1:2]
	secondCharCl := theme.UnstagedChangesColor
	if secondChar == " " {
		secondCharCl = restColor
	}

	return firstCharCl.Sprint(firstChar) + secondCharCl.Sprint(secondChar)
}

func getCommitFileLine(
	isCollapsed bool,
	treeDepth int,
	visualDepth int,
	node *filetree.Node[models.CommitFile],
	status patch.PatchStatus,
	showFileIcons bool,
) string {
	indentation := strings.Repeat("  ", visualDepth)
	name := commitFileNameAtDepth(node, treeDepth)
	commitFile := node.File
	output := indentation

	isDirectory := commitFile == nil

	nameColor := theme.DefaultTextColor

	switch status {
	case patch.WHOLE:
		nameColor = style.FgGreen
	case patch.PART:
		nameColor = style.FgYellow
	case patch.UNSELECTED:
		nameColor = theme.DefaultTextColor
	}

	if isDirectory {
		arrow := EXPANDED_ARROW
		if isCollapsed {
			arrow = COLLAPSED_ARROW
		}

		output += nameColor.Sprint(arrow) + " "
	} else {
		var symbol string
		symbolStyle := nameColor

		switch status {
		case patch.WHOLE:
			symbol = "●"
		case patch.PART:
			symbol = "◐"
		case patch.UNSELECTED:
			symbol = commitFile.ChangeStatus
			symbolStyle = getColorForChangeStatus(symbol)
		}

		output += symbolStyle.Sprint(symbol) + " "
	}

	name = utils.EscapeSpecialChars(name)
	isSubmodule := false
	isLinkedWorktree := false

	if showFileIcons {
		icon := icons.IconForFile(name, isSubmodule, isLinkedWorktree, isDirectory)
		paint := color.C256(icon.Color, false)
		output += paint.Sprint(icon.Icon) + " "
	}

	output += nameColor.Sprint(name)
	return output
}

func getColorForChangeStatus(changeStatus string) style.TextStyle {
	switch changeStatus {
	case "A":
		return style.FgGreen
	case "M", "R":
		return style.FgYellow
	case "D":
		return theme.UnstagedChangesColor
	case "C":
		return style.FgCyan
	case "T":
		return style.FgMagenta
	default:
		return theme.DefaultTextColor
	}
}

func fileNameAtDepth(node *filetree.Node[models.File], depth int) string {
	splitName := split(node.Path)
	name := join(splitName[depth:])

	if node.File != nil && node.File.IsRename() {
		splitPrevName := split(node.File.PreviousName)

		prevName := node.File.PreviousName
		// if the file has just been renamed inside the same directory, we can shave off
		// the prefix for the previous path too. Otherwise we'll keep it unchanged
		sameParentDir := len(splitName) == len(splitPrevName) && join(splitName[0:depth]) == join(splitPrevName[0:depth])
		if sameParentDir {
			prevName = join(splitPrevName[depth:])
		}

		return prevName + " → " + name
	}

	return name
}

func commitFileNameAtDepth(node *filetree.Node[models.CommitFile], depth int) string {
	splitName := split(node.Path)
	name := join(splitName[depth:])

	return name
}

func split(str string) []string {
	return strings.Split(str, "/")
}

func join(strs []string) string {
	return strings.Join(strs, "/")
}
