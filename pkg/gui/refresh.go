package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) Refresh(options types.RefreshOptions) error {
	return gui.helpers.Refresh.Refresh(options)
}
