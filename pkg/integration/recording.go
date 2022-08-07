package integration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/jesseduffield/gocui"
)

// this all relates to the old way of doing integration tests where you record yourself
// and then replay the events.

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

func SaveRecording(recording *gocui.Recording) error {
	if !RecordingEvents() {
		return nil
	}

	jsonEvents, err := json.Marshal(recording)
	if err != nil {
		return err
	}

	path := recordEventsTo()

	return ioutil.WriteFile(path, jsonEvents, 0o600)
}
