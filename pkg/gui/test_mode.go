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
	if os.Getenv(components.SANDBOX_ENV_VAR) == "true" || test == nil {
		return
	}

	go func() {
		defer gui.handlePanicInTest(test)

		time.Sleep(time.Millisecond * 100)

		guiDriver := &GuiDriver{gui: gui}
		test.Run(guiDriver)

		// if we're here then the test must have passed: it panics upon failure
		gui.Log.Warnf("test %s logging success", test.Name())
		if err := result.LogSuccess(); err != nil {
			gui.Log.Warnf("test %s failed to log success!", test.Name())
			panic(err)
		}

		gui.g.Update(func(*gocui.Gui) error {
			return gocui.ErrQuit
		})

		time.Sleep(time.Second * 1)

		log.Fatal("gocui should have already exited")
	}()

	go utils.Safe(func() {
		defer gui.handlePanicInTest(test)

		time.Sleep(time.Second * 40)
		panic("40 seconds is up, lazygit recording took too long to complete")
	})
}

func (gui *Gui) handlePanicInTest(test integrationTypes.IntegrationTest) {
	if test == nil {
		return
	}

	if r := recover(); r != nil {
		buf := make([]byte, 4096*4) // arbitrarily large buffer size
		stackSize := runtime.Stack(buf, false)
		stackTrace := string(buf[:stackSize])

		gui.Log.Warnf("test %s panicked!", test.Name())

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
