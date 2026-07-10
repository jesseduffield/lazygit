package gui

import (
	"log"
	"os"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type IntegrationTest interface {
	Run(*GuiDriver)
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

			toastChan := make(chan string, 100)
			gui.PopupHandler.(*popup.PopupHandler).SetToastFunc(
				func(message string, kind types.ToastKind) { toastChan <- message })

			test.Run(&GuiDriver{gui: gui, toastChan: toastChan, headless: Headless()})

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
				log.Fatalf("%v is up, lazygit integration test took too long to complete", timeout)
			})
		}
	}
}

func Headless() bool {
	return os.Getenv("LAZYGIT_HEADLESS") != ""
}
