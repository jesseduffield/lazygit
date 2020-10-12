package gui

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func recordingEvents() bool {
	return recordEventsTo() != ""
}

func recordEventsTo() string {
	return os.Getenv("RECORD_EVENTS_TO")
}

func (gui *Gui) timeSinceStart() int64 {
	return time.Since(gui.StartTime).Milliseconds()
}

func (gui *Gui) replayRecordedEvents() {
	if os.Getenv("REPLAY_EVENTS_FROM") == "" {
		return
	}

	go utils.Safe(func() {
		time.Sleep(time.Second * 20)
		log.Fatal("20 seconds is up, lazygit recording took too long to complete")
	})

	events, err := gui.loadRecordedEvents()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	// might need to add leeway if this ends up flakey
	var leeway int64 = 0
	// humans are slow so this speeds things up.
	speed := 1
	envReplaySpeed := os.Getenv("REPLAY_SPEED")
	if envReplaySpeed != "" {
		var err error
		speed, err = strconv.Atoi(envReplaySpeed)
		if err != nil {
			log.Fatal(err)
		}
	}

	// The playback could be paused at any time because integration tests run concurrently.
	// Therefore we can't just check for a given event whether we've passed its timestamp,
	// or else we'll have an explosion of keypresses after the test is resumed.
	// We need to check if we've waited long enough since the last event was replayed.
	for i, event := range events {
		var prevEventTimestamp int64 = 0
		if i > 0 {
			prevEventTimestamp = events[i-1].Timestamp
		}
		timeToWait := (event.Timestamp - prevEventTimestamp) / int64(speed)
		if i == 0 {
			timeToWait += leeway
		}
		var timeWaited int64 = 0
	middle:
		for {
			select {
			case <-ticker.C:
				timeWaited += 1
				if gui.g != nil && timeWaited >= timeToWait {
					gui.g.ReplayedEvents <- *event.Event
					break middle
				}
			case <-gui.stopChan:
				return
			}
		}
	}

	time.Sleep(time.Second * 1)

	gui.g.Update(func(*gocui.Gui) error {
		return gocui.ErrQuit
	})

	time.Sleep(time.Second * 1)

	log.Fatal("lazygit should have already exited")
}

func (gui *Gui) loadRecordedEvents() ([]RecordedEvent, error) {
	path := os.Getenv("REPLAY_EVENTS_FROM")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	events := []RecordedEvent{}

	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (gui *Gui) saveRecordedEvents() error {
	if !recordingEvents() {
		return nil
	}

	jsonEvents, err := json.Marshal(gui.RecordedEvents)
	if err != nil {
		return err
	}

	path := recordEventsTo()

	return ioutil.WriteFile(path, jsonEvents, 0600)
}

func (gui *Gui) recordEvents() {
	for event := range gui.g.RecordedEvents {
		recordedEvent := RecordedEvent{
			Timestamp: gui.timeSinceStart(),
			Event:     event,
		}

		gui.RecordedEvents = append(gui.RecordedEvents, recordedEvent)
	}
}
