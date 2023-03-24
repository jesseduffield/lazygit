package gui

import (
	"log"
	"os"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type IntegrationTest interface {
	Run(guiAdapter *GuiDriver)
}

func (gui *Gui) handleTestMode(test integrationTypes.IntegrationTest) {
	if os.Getenv(components.SANDBOX_ENV_VAR) == "true" {
		return
	}

	if test != nil {
		go func() {
			time.Sleep(time.Millisecond * 100)

			test.Run(&GuiDriver{gui: gui})

			gui.g.Update(func(*gocui.Gui) error {
				return gocui.ErrQuit
			})

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
