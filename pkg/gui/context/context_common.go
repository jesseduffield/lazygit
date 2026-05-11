package context

import (
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ContextCommon struct {
	*common.Common
	types.IGuiCommon
}
