package presentation

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
)

func GetSpiceStackDisplayStrings(
	items []*models.SpiceStackItem,
	getItemOperation func(item types.HasUrn) types.ItemOperation,
	diffName string,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
) [][]string {
	if len(items) == 0 {
		return [][]string{}
	}

	continuing := make(map[int]bool) // tracks which depth levels have more siblings

	return lo.Map(items, func(item *models.SpiceStackItem, idx int) []string {
		// Check if this is a commit item
		if item.IsCommit {
			prefix := buildCommitPrefix(item, idx, items, continuing)
			commitText := style.FgCyan.Sprint(item.CommitSha) + " " +
				style.FgDefault.Sprint(item.CommitSubject)
			return []string{prefix + commitText, "", ""}
		}

		// Regular branch item
		prefix := buildTreePrefix(item, idx, items, continuing)
		name := formatBranchName(item, diffName)
		status := formatStatus(item)
		pr := formatPR(item)

		return []string{prefix + name, status, pr}
	})
}

func buildTreePrefix(item *models.SpiceStackItem, idx int, items []*models.SpiceStackItem, continuing map[int]bool) string {
	if item.Depth == 0 {
		return "" // trunk has no prefix
	}

	var parts []string

	// Build vertical pipes for ancestor levels (git-spice fliptree approach)
	// Draw a pipe at depth d if the ancestor at that depth has siblingIndex > 0
	// In our reversed list (leaves first), ancestors come AFTER the current item
	for d := 1; d < item.Depth; d++ {
		ancestor := findAncestorAtDepthAfter(idx, d, items)
		if ancestor != nil && ancestor.SiblingIndex > 0 {
			parts = append(parts, "│ ")
		} else {
			parts = append(parts, "  ")
		}
	}

	// Joint for this node (lighter weight box-drawing characters)
	var joint string
	if item.SiblingIndex > 0 {
		// Not the first sibling, use middle branch connector
		joint = "├─"
	} else {
		// First sibling, use topmost branch connector
		joint = "┌─"
	}

	// Add horizontal-up if this node has children (items at depth+1 appear before this)
	hasChildren := hasItemsAtDepthBefore(idx, item.Depth+1, items)
	if hasChildren {
		joint += "┴"
	}

	// Add branch indicator: ● for current branch, ◯ for others
	var indicator string
	if item.Current {
		indicator = "●"
	} else {
		indicator = "◯"
	}

	// Always add one space after the indicator for consistent spacing before branch name
	parts = append(parts, joint+indicator+" ")
	return strings.Join(parts, "")
}

// findAncestorAtDepthAfter finds the ancestor at the given depth
// In our reversed list (leaves first), ancestors come AFTER the current index
func findAncestorAtDepthAfter(idx int, depth int, items []*models.SpiceStackItem) *models.SpiceStackItem {
	for i := idx + 1; i < len(items); i++ {
		if items[i].Depth == depth {
			return items[i]
		}
		if items[i].Depth < depth {
			// Reached shallower depth, stop searching
			break
		}
	}
	return nil
}

// hasItemsAtDepthBefore checks if there are items at the given depth before this index
func hasItemsAtDepthBefore(idx int, depth int, items []*models.SpiceStackItem) bool {
	for i := 0; i < idx; i++ {
		if items[i].Depth == depth {
			return true
		}
	}
	return false
}

func buildCommitPrefix(item *models.SpiceStackItem, idx int, items []*models.SpiceStackItem, continuing map[int]bool) string {
	var parts []string

	// Find the parent branch (most recent non-commit item before this commit)
	var parentBranch *models.SpiceStackItem
	for i := idx - 1; i >= 0; i-- {
		if !items[i].IsCommit {
			parentBranch = items[i]
			break
		}
	}

	if parentBranch == nil {
		// Shouldn't happen, but handle gracefully
		return ""
	}

	// For depths 1 to parentBranch.Depth-1: check ancestors like the branch would
	for d := 1; d < parentBranch.Depth; d++ {
		ancestor := findAncestorAtDepthAfter(idx, d, items)
		if ancestor != nil && ancestor.SiblingIndex > 0 {
			parts = append(parts, "│ ")
		} else {
			parts = append(parts, "  ")
		}
	}

	// At the parent branch's depth: always draw a pipe
	parts = append(parts, "│ ")

	// Add spacing to align with commit content (adjust based on typical joint width)
	parts = append(parts, "  ")

	return strings.Join(parts, "")
}

func formatBranchName(item *models.SpiceStackItem, diffName string) string {
	name := item.Name

	if item.Current {
		name = style.FgGreen.SetBold().Sprint(name)
	} else if item.Name == diffName {
		name = theme.DiffTerminalColor.Sprint(name)
	}

	return name
}

func formatStatus(item *models.SpiceStackItem) string {
	var parts []string

	if item.Current {
		parts = append(parts, style.FgGreen.Sprint("✓"))
	}
	if item.NeedsRestack {
		parts = append(parts, style.FgYellow.Sprint("⟳ restack"))
	}
	if item.NeedsPush {
		parts = append(parts, style.FgCyan.Sprint("↑ push"))
	}
	if item.Behind > 0 {
		parts = append(parts, style.FgRed.Sprintf("↓%d", item.Behind))
	}

	return strings.Join(parts, " ")
}

func formatPR(item *models.SpiceStackItem) string {
	if item.PRNumber == "" {
		return ""
	}

	var statusStyle style.TextStyle
	switch item.PRStatus {
	case "open":
		statusStyle = style.FgGreen
	case "merged":
		statusStyle = style.FgMagenta
	case "closed":
		statusStyle = style.FgRed
	default:
		statusStyle = theme.DefaultTextColor
	}

	return statusStyle.Sprint(item.PRNumber)
}
