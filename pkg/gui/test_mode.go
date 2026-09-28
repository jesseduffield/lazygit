package gui

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type IntegrationTest interface {
	Run(*GuiDriver)
}

// captureToastsForIntegrationTest sends the toasts of an integration test to a
// channel, from which the test checks them with ExpectToast. NewGui calls it
// before anything can show a toast, so that the test also gets the toasts
// shown while lazygit starts up.
func (gui *Gui) captureToastsForIntegrationTest(test integrationTypes.IntegrationTest) {
	if test == nil || os.Getenv(components.SANDBOX_ENV_VAR) == "true" {
		return
	}

	toastChan := make(chan string, 100)
	gui.PopupHandler.(*popup.PopupHandler).SetToastFunc(
		func(message string, kind types.ToastKind) { toastChan <- message })
	gui.testToastChan = toastChan
}

func (gui *Gui) handleTestMode() {
	test := gui.integrationTest
	if os.Getenv(components.SANDBOX_ENV_VAR) == "true" {
		return
	}

	if test != nil {
		waitUntilIdle := func() {
			gui.c.GocuiGui().WaitUntilIdle()
		}

		go func() {
			waitUntilIdle()

			test.Run(&GuiDriver{gui: gui, toastChan: gui.testToastChan, headless: Headless()})

			gui.g.Update(func(*gocui.Gui) error {
				return gocui.ErrQuit
			})

			// Wait for the event loop to actually exit.
			<-gui.g.LoopExited()
		}()

		if os.Getenv(components.WAIT_FOR_DEBUGGER_ENV_VAR) == "" {
			timeout := 40 * time.Second * testTimeoutMultiplier
			go utils.Safe(func() {
				time.Sleep(timeout)
				// Dump all goroutine stacks before dying, so a hung test shows
				// where it got stuck rather than just that it timed out. The
				// test harness surfaces this process's stderr on failure.
				_ = pprof.Lookup("goroutine").WriteTo(os.Stderr, 2)
				log.Fatalf("%v is up, lazygit integration test took too long to complete", timeout)
			})
		}
	}
}

func Headless() bool {
	return os.Getenv("LAZYGIT_HEADLESS") != ""
}
