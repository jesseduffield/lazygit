package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// all fields mandatory (except `CanRebase` because it's boolean)
type SwitchToCommitFilesContextOpts struct {
	RefName    string
	CanRebase  bool
	Context    types.Context
	WindowName string
}
