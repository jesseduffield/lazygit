package types

import (
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/filtering"
)

type Modes struct {
	Filtering     filtering.Filtering
	CherryPicking *cherrypicking.CherryPicking
	Diffing       diffing.Diffing
}
