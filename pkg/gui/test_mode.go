package gui

import (
	"log"
	"os"
	"time"

	"github.com/jesseduffield/gocui"
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

			test.Run(&GuiDriver{gui: gui, isIdleChan: isIdleChan})

			gui.g.Update(func(*gocui.Gui) error {
				return gocui.ErrQuit
			})

			waitUntilIdle()

			time.Sleep(time.Second * 1)

			log.Fatal("gocui should have already exited")
		}()

		go utils.Safe(func() {
			time.Sleep(time.Second * 40)
			log.Fatal("40 seconds is up, lazygit recording took too long to complete")
		})
	}
}

func Headless() bool {
	return os.Getenv("HEADLESS") != ""
}
