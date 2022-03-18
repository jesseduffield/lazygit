package gui

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/jesseduffield/gocui"
)

func recordingEvents() bool {
	return recordEventsTo() != ""
}

func recordEventsTo() string {
	return os.Getenv("RECORD_EVENTS_TO")
}

func replaying() bool {
	return os.Getenv("REPLAY_EVENTS_FROM") != ""
}

func headless() bool {
	return os.Getenv("HEADLESS") != ""
}

func getRecordingSpeed() float64 {
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

func (gui *Gui) loadRecording() (*gocui.Recording, error) {
	path := os.Getenv("REPLAY_EVENTS_FROM")

	data, err := ioutil.ReadFile(path)
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

func (gui *Gui) saveRecording(recording *gocui.Recording) error {
	if !recordingEvents() {
		return nil
	}

	jsonEvents, err := json.Marshal(recording)
	if err != nil {
		return err
	}

	path := recordEventsTo()

	return ioutil.WriteFile(path, jsonEvents, 0o600)
}
