package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

var (
	// DefaultTextColor is the default text color
	DefaultTextColor = style.FgDefault

	// GocuiDefaultTextColor does the same as DefaultTextColor but this one only colors gocui default text colors
	GocuiDefaultTextColor = gocui.ColorDefault

	// ActiveBorderColor is the border color of the active frame
	ActiveBorderColor gocui.Attribute

	// InactiveBorderColor is the border color of the inactive active frames
	InactiveBorderColor gocui.Attribute

	// FilteredActiveBorderColor is the border color of the active frame, when it's being searched/filtered
	SearchingActiveBorderColor gocui.Attribute

	// GocuiSelectedLineBgColor is the background color for the selected line in gocui
	GocuiSelectedLineBgColor gocui.Attribute
	// GocuiInactiveViewSelectedLineBgColor is the background color for the selected line in gocui if the view doesn't have focus
	GocuiInactiveViewSelectedLineBgColor gocui.Attribute

	OptionsColor gocui.Attribute

	// SelectedLineBgColor is the background color for the selected line
	SelectedLineBgColor = style.New()
	// InactiveViewSelectedLineBgColor is the background color for the selected line if the view doesn't have the focus
	InactiveViewSelectedLineBgColor = style.New()

	// CherryPickedCommitColor is the text style when cherry picking a commit
	CherryPickedCommitTextStyle = style.New()

	// MarkedBaseCommitTextStyle is the text style of the marked rebase base commit
	MarkedBaseCommitTextStyle = style.New()

	OptionsFgColor = style.New()

	DiffTerminalColor = style.FgMagenta

	UnstagedChangesColor = style.New()

	ModifiedColor  = style.New()
	DeletedColor   = style.New()
	UntrackedColor = style.New()
	AddedColor     = style.New()
	RenamedColor   = style.New()
	CopiedColor    = style.New()
)

// UpdateTheme updates all theme variables
func UpdateTheme(themeConfig config.ThemeConfig) {
	ActiveBorderColor = GetGocuiStyle(themeConfig.ActiveBorderColor)
	InactiveBorderColor = GetGocuiStyle(themeConfig.InactiveBorderColor)
	SearchingActiveBorderColor = GetGocuiStyle(themeConfig.SearchingActiveBorderColor)
	SelectedLineBgColor = GetTextStyle(themeConfig.SelectedLineBgColor, true)
	InactiveViewSelectedLineBgColor = GetTextStyle(themeConfig.InactiveViewSelectedLineBgColor, true)

	cherryPickedCommitBgTextStyle := GetTextStyle(themeConfig.CherryPickedCommitBgColor, true)
	cherryPickedCommitFgTextStyle := GetTextStyle(themeConfig.CherryPickedCommitFgColor, false)
	CherryPickedCommitTextStyle = cherryPickedCommitBgTextStyle.MergeStyle(cherryPickedCommitFgTextStyle)

	markedBaseCommitBgTextStyle := GetTextStyle(themeConfig.MarkedBaseCommitBgColor, true)
	markedBaseCommitFgTextStyle := GetTextStyle(themeConfig.MarkedBaseCommitFgColor, false)
	MarkedBaseCommitTextStyle = markedBaseCommitBgTextStyle.MergeStyle(markedBaseCommitFgTextStyle)

	unstagedChangesTextStyle := GetTextStyle(themeConfig.UnstagedChangesColor, false)
	UnstagedChangesColor = unstagedChangesTextStyle

	// Per-status colors — fall back to UnstagedChangesColor or hardcoded defaults if not set
	if len(themeConfig.ModifiedColor) > 0 {
		ModifiedColor = GetTextStyle(themeConfig.ModifiedColor, false)
	} else {
		ModifiedColor = UnstagedChangesColor
	}
	if len(themeConfig.DeletedColor) > 0 {
		DeletedColor = GetTextStyle(themeConfig.DeletedColor, false)
	} else {
		DeletedColor = UnstagedChangesColor
	}
	if len(themeConfig.UntrackedColor) > 0 {
		UntrackedColor = GetTextStyle(themeConfig.UntrackedColor, false)
	} else {
		UntrackedColor = UnstagedChangesColor
	}
	if len(themeConfig.AddedColor) > 0 {
		AddedColor = GetTextStyle(themeConfig.AddedColor, false)
	} else {
		AddedColor = style.FgGreen
	}
	if len(themeConfig.RenamedColor) > 0 {
		RenamedColor = GetTextStyle(themeConfig.RenamedColor, false)
	} else {
		RenamedColor = UnstagedChangesColor
	}
	if len(themeConfig.CopiedColor) > 0 {
		CopiedColor = GetTextStyle(themeConfig.CopiedColor, false)
	} else {
		CopiedColor = style.FgCyan
	}

	GocuiSelectedLineBgColor = GetGocuiStyle(themeConfig.SelectedLineBgColor)
	GocuiInactiveViewSelectedLineBgColor = GetGocuiStyle(themeConfig.InactiveViewSelectedLineBgColor)
	OptionsColor = GetGocuiStyle(themeConfig.OptionsTextColor)
	OptionsFgColor = GetTextStyle(themeConfig.OptionsTextColor, false)

	DefaultTextColor = GetTextStyle(themeConfig.DefaultFgColor, false)
	GocuiDefaultTextColor = GetGocuiStyle(themeConfig.DefaultFgColor)
}
