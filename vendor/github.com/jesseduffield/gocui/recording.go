package gocui

import (
	"log"
	"time"
)

func (g *Gui) replayRecording() {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	// The playback could be paused at any time because integration tests run concurrently.
	// Therefore we can't just check for a given event whether we've passed its timestamp,
	// or else we'll have an explosion of keypresses after the test is resumed.
	// We need to check if we've waited long enough since the last event was replayed.
	// Only handling key events for now.
	for i, event := range g.Recording.KeyEvents {
		var prevEventTimestamp int64 = 0
		if i > 0 {
			prevEventTimestamp = g.Recording.KeyEvents[i-1].Timestamp
		}
		timeToWait := (event.Timestamp - prevEventTimestamp) / int64(g.RecordingConfig.Speed)
		if i == 0 {
			timeToWait += int64(g.RecordingConfig.Leeway)
		}
		var timeWaited int64 = 0
	middle:
		for {
			select {
			case <-ticker.C:
				timeWaited += 1
				if g != nil && timeWaited >= timeToWait {
					g.ReplayedEvents.keys <- event
					break middle
				}
			case <-g.stop:
				return
			}
		}
	}

	// leaving some time for any handlers to execute before quitting
	time.Sleep(time.Second * 1)

	g.Update(func(*Gui) error {
		return ErrQuit
	})

	time.Sleep(time.Second * 1)

	log.Fatal("gocui should have already exited")
}
