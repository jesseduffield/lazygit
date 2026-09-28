package status

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

// displayedStatus is what the status line shows for the current statuses.
type displayedStatus struct {
	text  string
	color gocui.Attribute
}

func getDisplayedStatus(statusManager *StatusManager) displayedStatus {
	userConfig := config.GetDefaultConfig()
	// A single frame makes the spinner after a waiting status independent of
	// the time at which the test runs.
	userConfig.Gui.Spinner.Frames = []string{"*"}

	text, color := statusManager.GetStatusString(userConfig)
	return displayedStatus{text: text, color: color}
}

// expireToast does what the goroutine started by AddToastStatus does once the
// toast's time is up, so that the tests don't have to wait for seconds.
func expireToast(statusManager *StatusManager, toastId int) {
	statusManager.removeStatus(toastId)
}

func noRender() {}

func TestStatusLineIsEmptyWithoutStatuses(t *testing.T) {
	statusManager := NewStatusManager()

	assert.Equal(t, displayedStatus{text: "", color: gocui.ColorDefault}, getDisplayedStatus(statusManager))
}

func TestWaitingStatusIsShownWithSpinnerWhileItsOperationRuns(t *testing.T) {
	statusManager := NewStatusManager()

	assert.NoError(t, statusManager.WithWaitingStatus("Fetching", noRender, func(*WaitingStatusHandle) error {
		assert.Equal(t, displayedStatus{text: "Fetching *", color: gocui.ColorCyan}, getDisplayedStatus(statusManager))
		return nil
	}))

	assert.Equal(t, displayedStatus{text: "", color: gocui.ColorDefault}, getDisplayedStatus(statusManager))
}

func TestNewestWaitingStatusIsShown(t *testing.T) {
	statusManager := NewStatusManager()

	assert.NoError(t, statusManager.WithWaitingStatus("Fetching", noRender, func(*WaitingStatusHandle) error {
		assert.NoError(t, statusManager.WithWaitingStatus("Rebasing", noRender, func(*WaitingStatusHandle) error {
			assert.Equal(t, displayedStatus{text: "Rebasing *", color: gocui.ColorCyan}, getDisplayedStatus(statusManager))
			return nil
		}))

		assert.Equal(t, displayedStatus{text: "Fetching *", color: gocui.ColorCyan}, getDisplayedStatus(statusManager))
		return nil
	}))
}

func TestNewestToastIsShown(t *testing.T) {
	statusManager := NewStatusManager()

	statusManager.AddToastStatus("Something went wrong", types.ToastKindError)
	toastId := statusManager.AddToastStatus("Copied to clipboard", types.ToastKindStatus)
	assert.Equal(t, displayedStatus{text: "Copied to clipboard", color: gocui.ColorCyan}, getDisplayedStatus(statusManager))

	expireToast(statusManager, toastId)
	assert.Equal(t, displayedStatus{text: "Something went wrong", color: gocui.ColorRed}, getDisplayedStatus(statusManager))
}

func TestToastIsShownAheadOfOlderWaitingStatus(t *testing.T) {
	statusManager := NewStatusManager()

	assert.NoError(t, statusManager.WithWaitingStatus("Fetching", noRender, func(*WaitingStatusHandle) error {
		toastId := statusManager.AddToastStatus("Copied to clipboard", types.ToastKindStatus)
		assert.Equal(t, displayedStatus{text: "Copied to clipboard", color: gocui.ColorCyan}, getDisplayedStatus(statusManager))

		expireToast(statusManager, toastId)
		assert.Equal(t, displayedStatus{text: "Fetching *", color: gocui.ColorCyan}, getDisplayedStatus(statusManager))
		return nil
	}))
}
