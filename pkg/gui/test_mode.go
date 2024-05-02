package gui

import (
	"log"
	"os"
	"time"

	"github.com/jesseduffield/gocui"
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
		isIdleChan := make(chan struct{})

		gui.c.GocuiGui().AddIdleListener(isIdleChan)

		waitUntilIdle := func() {
			<-isIdleChan
		}

		go func() {
			waitUntilIdle()

			toastChan := make(chan string, 100)
			gui.PopupHandler.(*popup.PopupHandler).SetToastFunc(
				func(message string, kind types.ToastKind) { toastChan <- message })

			test.Run(&GuiDriver{gui: gui, isIdleChan: isIdleChan, toastChan: toastChan})

			gui.g.Update(func(*gocui.Gui) error {
				return gocui.ErrQuit
			})

			waitUntilIdle()

			time.Sleep(time.Second * 1)

			log.Fatal("gocui should have already exited")
		}()

		if os.Getenv(components.WAIT_FOR_DEBUGGER_ENV_VAR) == "" {
			go utils.Safe(func() {
				time.Sleep(time.Second * 40)
				log.Fatal("40 seconds is up, lazygit recording took too long to complete")
			})
		}
	}
}

func Headless() bool {
	return os.Getenv("HEADLESS") != ""
}
