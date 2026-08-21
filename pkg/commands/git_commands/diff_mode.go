package git_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
)

// DiffMode says what a diff command's output is for, which is what decides whether the
// configured diff renderer produces it and whether it is coloured.
type DiffMode int

const (
	// DiffModeRendered is the diff as the user has arranged for it to look: through the
	// diff renderer, with the renderer's own arguments and its preference about colour.
	DiffModeRendered DiffMode = iota
	// DiffModeRaw is git's own coloured diff, for showing a diff whose rendered form
	// couldn't be acted on.
	DiffModeRaw
	// DiffModePlain is git's own uncoloured diff, which is what patches are built from and
	// text is copied out of, rather than anything to look at.
	DiffModePlain
)

// colorArg is what to pass to git's --color for this mode. Rendered output is coloured
// however the renderer wants its input; a raw diff is git's own colour, that being the
// point of it; a plain one is for reading as text, not for looking at.
func (self DiffMode) colorArg(diffRendererConfigManager *config.DiffRendererConfigManager) string {
	switch self {
	case DiffModeRendered:
		return diffRendererConfigManager.GetColorArg()
	case DiffModeRaw:
		return "always"
	default:
		return "never"
	}
}
