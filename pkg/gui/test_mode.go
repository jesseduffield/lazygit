package gui

import (
	"fmt"
	"log"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/integration"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type IntegrationTest interface {
	Run(guiAdapter *GuiAdapterImpl)
}

func (gui *Gui) handleTestMode() {
	if integration.PlayingIntegrationTest() {
		test, ok := integration.CurrentIntegrationTest()

		if !ok {
			panic(fmt.Sprintf("test %s not found", integration.IntegrationTestName()))
		}

		go func() {
			time.Sleep(time.Millisecond * 100)

			test.Run(&GuiAdapterImpl{gui: gui})

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

	if integration.Replaying() {
		gui.g.RecordingConfig = gocui.RecordingConfig{
			Speed:  integration.GetRecordingSpeed(),
			Leeway: 100,
		}

		var err error
		gui.g.Recording, err = integration.LoadRecording()
		if err != nil {
			panic(err)
		}

		go utils.Safe(func() {
			time.Sleep(time.Second * 40)
			log.Fatal("40 seconds is up, lazygit recording took too long to complete")
		})
	}
}
