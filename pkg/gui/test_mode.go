package gui

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jesseduffield/gocui"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type IntegrationTest interface {
	Run(guiAdapter *GuiDriver)
}

func (gui *Gui) handleTestMode(test integrationTypes.IntegrationTest) {
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

	if Replaying() {
		gui.g.RecordingConfig = gocui.RecordingConfig{
			Speed:  GetRecordingSpeed(),
			Leeway: 1000,
		}

		var err error
		gui.g.Recording, err = LoadRecording()
		if err != nil {
			panic(err)
		}

		go utils.Safe(func() {
			time.Sleep(time.Second * 40)
			log.Fatal("40 seconds is up, lazygit recording took too long to complete")
		})
	}
}

func Headless() bool {
	return os.Getenv("HEADLESS") != ""
}

// OLD integration test format stuff

func Replaying() bool {
	return os.Getenv("REPLAY_EVENTS_FROM") != ""
}

func RecordingEvents() bool {
	return recordEventsTo() != ""
}

func recordEventsTo() string {
	return os.Getenv("RECORD_EVENTS_TO")
}

func GetRecordingSpeed() float64 {
	// humans are slow so this speeds things up.
	speed := 1.0
	envReplaySpeed := os.Getenv("SPEED")
	if envReplaySpeed != "" {
		var err error
		speed, err = strconv.ParseFloat(envReplaySpeed, 64)
		if err != nil {
			log.Fatal(err)
		}
	}
	return speed
}

func LoadRecording() (*gocui.Recording, error) {
	path := os.Getenv("REPLAY_EVENTS_FROM")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	recording := &gocui.Recording{}

	err = json.Unmarshal(data, &recording)
	if err != nil {
		return nil, err
	}

	return recording, nil
}

func SaveRecording(recording *gocui.Recording) error {
	if !RecordingEvents() {
		return nil
	}

	jsonEvents, err := json.Marshal(recording)
	if err != nil {
		return err
	}

	path := recordEventsTo()

	return os.WriteFile(path, jsonEvents, 0o600)
}
