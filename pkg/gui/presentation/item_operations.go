package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

func ItemOperationToString(itemOperation types.ItemOperation, tr *i18n.TranslationSet) string {
	switch itemOperation {
	case types.ItemOperationNone:
		return ""
	case types.ItemOperationPushing:
		return tr.PushingStatus
	case types.ItemOperationPulling:
		return tr.PullingStatus
	case types.ItemOperationFastForwarding:
		return tr.FastForwarding
	case types.ItemOperationDeleting:
		return tr.DeletingStatus
	case types.ItemOperationFetching:
		return tr.FetchingStatus
	case types.ItemOperationCheckingOut:
		return tr.CheckingOutStatus
	}

	return ""
}
