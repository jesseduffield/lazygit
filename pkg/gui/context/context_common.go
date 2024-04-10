package context

import (
	"github.com/lobes/lazytask/pkg/common"
	"github.com/lobes/lazytask/pkg/gui/types"
)

type ContextCommon struct {
	*common.Common
	types.IGuiCommon
}
