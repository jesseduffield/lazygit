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
		// Check if more roots exist after this one
		hasMoreRoots := false
		for i := idx + 1; i < len(items); i++ {
			if items[i].Depth == 0 {
				hasMoreRoots = true
				break
			}
		}
		continuing[0] = hasMoreRoots
		return ""
	}

	var parts []string

	// Build vertical lines showing connection to parent below
	for d := 1; d < item.Depth; d++ {
		if continuing[d] {
			parts = append(parts, "│  ")
		} else {
			parts = append(parts, "   ")
		}
	}

	// In a top-down stack view, items point down to their parent
	// All items use ┌─ to show connection to parent below
	parts = append(parts, "┌─ ")
	continuing[item.Depth] = true

	return strings.Join(parts, "")
}

func buildCommitPrefix(item *models.SpiceStackItem, idx int, items []*models.SpiceStackItem, continuing map[int]bool) string {
	var parts []string

	// Vertical lines for ancestor levels (branch level)
	for d := 1; d < item.Depth; d++ {
		if continuing[d] {
			parts = append(parts, "│  ")
		} else {
			parts = append(parts, "   ")
		}
	}

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
