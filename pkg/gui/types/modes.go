package types

import (
	"github.com/lobes/lazytask/pkg/gui/modes/cherrypicking"
	"github.com/lobes/lazytask/pkg/gui/modes/diffing"
	"github.com/lobes/lazytask/pkg/gui/modes/filtering"
	"github.com/lobes/lazytask/pkg/gui/modes/marked_base_commit"
)

type Modes struct {
	Filtering        filtering.Filtering
	CherryPicking    *cherrypicking.CherryPicking
	Diffing          diffing.Diffing
	MarkedBaseCommit marked_base_commit.MarkedBaseCommit
}
