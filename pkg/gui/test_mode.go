package gui

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/result"
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

			defer handlePanic()

			guiDriver := &GuiDriver{gui: gui}
			test.Run(guiDriver)

			// if we're here then the test must have passed: it panics upon failure
			if err := result.LogSuccess(); err != nil {
				panic(err)
			}

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

func handlePanic() {
	if r := recover(); r != nil {
		buf := make([]byte, 4096*4) // arbitrarily large buffer size
		stackSize := runtime.Stack(buf, false)
		stackTrace := string(buf[:stackSize])

		if err := result.LogFailure(fmt.Sprintf("%v\n%s", r, stackTrace)); err != nil {
			panic(err)
		}

		// Re-panic
		panic(r)
	}
}

func Headless() bool {
	return os.Getenv("HEADLESS") != ""
}
