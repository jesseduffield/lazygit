package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

func refOperationToString(refOperation types.ItemOperation, tr *i18n.TranslationSet) string {
	switch refOperation {
	case types.ItemOperationNone:
		return ""
	case types.ItemOperationPushing:
		return tr.PushingStatus
	case types.ItemOperationPulling:
		return tr.PullingStatus
	case types.ItemOperationFastForwarding:
		return tr.FastForwardingOperation
	case types.ItemOperationDeleting:
		return tr.DeletingStatus
	}

	return ""
}
