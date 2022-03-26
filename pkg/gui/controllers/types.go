package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// all fields mandatory (except `CanRebase` because it's boolean)
type SwitchToCommitFilesContextOpts struct {
	// this is something like a commit or branch
	Ref types.Ref

	// from the local commits view we're allowed to do rebase stuff with any patch
	// we generate from the diff files context, but we don't have that same ability
	// with say the sub commits context or the reflog context.
	CanRebase bool

	Context types.Context
}
