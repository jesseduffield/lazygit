package status

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestStatusStringShowsNewestToast(t *testing.T) {
	statusManager := NewStatusManager()
	statusManager.AddToastStatus("first toast", types.ToastKindStatus)
	statusManager.AddToastStatus("second toast", types.ToastKindStatus)

	message, _ := statusManager.GetStatusString(config.GetDefaultConfig())
	assert.Equal(t, "second toast", message)
}

func TestStatusStringShowsToastAboveLaterWaitingStatus(t *testing.T) {
	statusManager := NewStatusManager()
	statusManager.AddToastStatus("toast", types.ToastKindError)

	_ = statusManager.WithWaitingStatus("Fetching", func() {}, func(*WaitingStatusHandle) error {
		message, _ := statusManager.GetStatusString(config.GetDefaultConfig())
		assert.Equal(t, "toast", message)
		return nil
	})
}

func TestStatusStringShowsToastAboveEarlierWaitingStatus(t *testing.T) {
	statusManager := NewStatusManager()

	_ = statusManager.WithWaitingStatus("Fetching", func() {}, func(*WaitingStatusHandle) error {
		statusManager.AddToastStatus("toast", types.ToastKindError)

		message, _ := statusManager.GetStatusString(config.GetDefaultConfig())
		assert.Equal(t, "toast", message)
		return nil
	})
}

func TestStatusStringShowsNewestWaitingStatusWithoutToasts(t *testing.T) {
	statusManager := NewStatusManager()

	_ = statusManager.WithWaitingStatus("Fetching", func() {}, func(*WaitingStatusHandle) error {
		return statusManager.WithWaitingStatus("Rebasing", func() {}, func(*WaitingStatusHandle) error {
			message, _ := statusManager.GetStatusString(config.GetDefaultConfig())
			assert.Regexp(t, "^Rebasing ", message)
			return nil
		})
	})
}
