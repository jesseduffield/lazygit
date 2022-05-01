package presentation

import (
	"strings"

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
	COLLAPSED_ARROW = "►"
)

// keeping these here as individual constants in case later on people want the old tree shape
const (
	INNER_ITEM = "  "
	LAST_ITEM  = "  "
	NESTED     = "  "
	NOTHING    = "  "
)

func RenderFileTree(
	tree filetree.IFileTree,
	diffName string,
	submoduleConfigs []*models.SubmoduleConfig,
) []string {
	return renderAux(tree.Tree(), tree.CollapsedPaths(), "", -1, func(n filetree.INode, depth int) string {
		castN := n.(*filetree.FileNode)
		return getFileLine(castN.GetHasUnstagedChanges(), castN.GetHasStagedChanges(), castN.NameAtDepth(depth), diffName, submoduleConfigs, castN.File)
	})
}

func RenderCommitFileTree(
	tree *filetree.CommitFileTreeViewModel,
	diffName string,
	patchManager *patch.PatchManager,
) []string {
	return renderAux(tree.Tree(), tree.CollapsedPaths(), "", -1, func(n filetree.INode, depth int) string {
		castN := n.(*filetree.CommitFileNode)

		// This is a little convoluted because we're dealing with either a leaf or a non-leaf.
		// But this code actually applies to both. If it's a leaf, the status will just
		// be whatever status it is, but if it's a non-leaf it will determine its status
		// based on the leaves of that subtree
		var status patch.PatchStatus
		if castN.EveryFile(func(file *models.CommitFile) bool {
			return patchManager.GetFileStatus(file.Name, tree.GetRef().RefName()) == patch.WHOLE
		}) {
			status = patch.WHOLE
		} else if castN.EveryFile(func(file *models.CommitFile) bool {
			return patchManager.GetFileStatus(file.Name, tree.GetRef().RefName()) == patch.UNSELECTED
		}) {
			status = patch.UNSELECTED
		} else {
			status = patch.PART
		}

		return getCommitFileLine(castN.NameAtDepth(depth), diffName, castN.File, status)
	})
}

func renderAux(
	s filetree.INode,
	collapsedPaths *filetree.CollapsedPaths,
	prefix string,
	depth int,
	renderLine func(filetree.INode, int) string,
) []string {
	if s == nil || s.IsNil() {
		return []string{}
	}

	isRoot := depth == -1

	if s.IsLeaf() {
		if isRoot {
			return []string{}
		}
		return []string{prefix + renderLine(s, depth)}
	}

	if collapsedPaths.IsCollapsed(s.GetPath()) {
		return []string{prefix + COLLAPSED_ARROW + " " + renderLine(s, depth)}
	}

	arr := []string{}
	if !isRoot {
		arr = append(arr, prefix+EXPANDED_ARROW+" "+renderLine(s, depth))
	}

	newPrefix := prefix
	if strings.HasSuffix(prefix, LAST_ITEM) {
		newPrefix = strings.TrimSuffix(prefix, LAST_ITEM) + NOTHING
	} else if strings.HasSuffix(prefix, INNER_ITEM) {
		newPrefix = strings.TrimSuffix(prefix, INNER_ITEM) + NESTED
	}

	for i, child := range s.GetChildren() {
		isLast := i == len(s.GetChildren())-1

		var childPrefix string
		if isRoot {
			childPrefix = newPrefix
		} else if isLast {
			childPrefix = newPrefix + LAST_ITEM
		} else {
			childPrefix = newPrefix + INNER_ITEM
		}

		arr = append(arr, renderAux(child, collapsedPaths, childPrefix, depth+1+s.GetCompressionLevel(), renderLine)...)
	}

	return arr
}

func getFileLine(hasUnstagedChanges bool, hasStagedChanges bool, name string, diffName string, submoduleConfigs []*models.SubmoduleConfig, file *models.File) string {
	// potentially inefficient to be instantiating these color
	// objects with each render
	partiallyModifiedColor := style.FgYellow

	restColor := style.FgGreen
	if name == diffName {
		restColor = theme.DiffTerminalColor
	} else if file == nil && hasStagedChanges && hasUnstagedChanges {
		restColor = partiallyModifiedColor
	} else if hasUnstagedChanges {
		restColor = theme.UnstagedChangesColor
	}

	output := ""
	if file != nil {
		// this is just making things look nice when the background attribute is 'reverse'
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

		output = firstCharCl.Sprint(firstChar)
		output += secondCharCl.Sprint(secondChar)
		output += restColor.Sprint(" ")
	}

	isSubmodule := file != nil && file.IsSubmodule(submoduleConfigs)
	isDirectory := file == nil

	if icons.IsIconEnabled() {
		output += restColor.Sprintf("%s ", icons.IconForFile(name, isSubmodule, isDirectory))
	}

	output += restColor.Sprint(utils.EscapeSpecialChars(name))

	if isSubmodule {
		output += theme.DefaultTextColor.Sprint(" (submodule)")
	}

	return output
}

func getCommitFileLine(name string, diffName string, commitFile *models.CommitFile, status patch.PatchStatus) string {
	var colour style.TextStyle
	if diffName == name {
		colour = theme.DiffTerminalColor
	} else {
		switch status {
		case patch.WHOLE:
			colour = style.FgGreen
		case patch.PART:
			colour = style.FgYellow
		case patch.UNSELECTED:
			colour = theme.DefaultTextColor
		}
	}

	output := ""

	name = utils.EscapeSpecialChars(name)
	if commitFile != nil {
		output += getColorForChangeStatus(commitFile.ChangeStatus).Sprint(commitFile.ChangeStatus) + " "
	}

	isSubmodule := false
	isDirectory := commitFile == nil

	if icons.IsIconEnabled() {
		output += colour.Sprintf("%s ", icons.IconForFile(name, isSubmodule, isDirectory))
	}

	output += colour.Sprint(name)
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
