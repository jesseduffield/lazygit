package gui

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jesseduffield/gocui"
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

	events, err := gui.loadRecordedEvents()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	// might need to add leeway if this ends up flakey
	var leeway int64 = 0
	// humans are slow so this speeds things up.
	var speed int64 = 5

	for _, event := range events {
		for range ticker.C {
			now := gui.timeSinceStart()*speed - leeway
			if gui.g != nil && now >= event.Timestamp {
				gui.g.ReplayedEvents <- *event.Event
				break
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
