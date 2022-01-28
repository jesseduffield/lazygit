package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IGuiCommon interface {
	popup.IPopupHandler

	LogAction(string)
	Refresh(types.RefreshOptions) error
}
