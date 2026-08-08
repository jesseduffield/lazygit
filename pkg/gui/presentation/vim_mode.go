package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

// VimModeSubTitle returns the padded mode indicator segment (" NORMAL ")
// for a view with vim-style editing attached, and "" for views without one.
func VimModeSubTitle(tr *i18n.TranslationSet, view *gocui.View) string {
	vim := view.VimEditor()
	if vim == nil {
		return ""
	}
	var label string
	switch vim.Mode() {
	case gocui.VimModeNormal:
		label = tr.VimModeNormalSubTitle
	case gocui.VimModeVisual:
		label = tr.VimModeVisualSubTitle
	default:
		label = tr.VimModeInsertSubTitle
	}
	return " " + label + " "
}
