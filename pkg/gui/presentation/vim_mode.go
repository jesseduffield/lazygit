package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

// VimModeSubTitle returns the padded mode indicator segment (" NORMAL ")
// for a view with vim-style editing attached, and "" for views without one.
// Like vim's single mode line, only the focused view shows an indicator;
// with the commit summary and description each holding their own editor,
// two simultaneous indicators would suggest a shared mode that isn't there.
func VimModeSubTitle(tr *i18n.TranslationSet, view *gocui.View, focused bool) string {
	vim := view.VimEditor()
	if vim == nil || !focused {
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
