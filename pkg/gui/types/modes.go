package types

import (
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/filtering"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/marked_base_commit"
)

type Modes struct {
	Filtering        filtering.Filtering
	CherryPicking    *cherrypicking.CherryPicking // shared by all repo states; see Gui.cherryPicking
	Diffing          diffing.Diffing
	MarkedBaseCommit marked_base_commit.MarkedBaseCommit
}
