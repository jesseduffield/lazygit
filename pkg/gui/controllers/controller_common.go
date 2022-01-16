package controllers

import "github.com/jesseduffield/lazygit/pkg/common"

// if Go let me do private struct embedding of structs with public fields (which it should)
// I would just do that. But alas.
type ControllerCommon struct {
	*common.Common
	IGuiCommon
}
