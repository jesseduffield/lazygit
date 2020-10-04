package gui

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func recordingEvents() bool {
	return os.Getenv("RECORD_EVENTS") == "true"
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

	var leeway int64 = 1000

	for _, event := range events {
		for range ticker.C {
			now := gui.timeSinceStart() - leeway
			if gui.g != nil && now >= event.Timestamp {
				gui.g.ReplayedEvents <- *event.Event
				break
			}
		}
	}
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

	return ioutil.WriteFile("recorded_events.json", jsonEvents, 0600)
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
