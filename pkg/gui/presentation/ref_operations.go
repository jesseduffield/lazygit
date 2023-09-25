package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

func refOperationToString(refOperation types.RefOperation, tr *i18n.TranslationSet) string {
	switch refOperation {
	case types.RefOperationNone:
		return ""
	case types.RefOperationPushing:
		return tr.PushingStatus
	case types.RefOperationPulling:
		return tr.PullingStatus
	case types.RefOperationFastForwarding:
		return tr.FastForwardingOperation
	case types.RefOperationDeleting:
		return tr.DeletingStatus
	}

	return ""
}
