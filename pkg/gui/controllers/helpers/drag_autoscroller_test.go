package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDragAutoscrollZone(t *testing.T) {
	testCases := []struct {
		name              string
		pointerViewportY  int
		expectedDirection int
		expectedInterval  time.Duration
	}{
		{name: "above view", pointerViewportY: -1, expectedDirection: -1, expectedInterval: dragAutoscrollVeryFastInterval},
		{name: "top outer row", pointerViewportY: 0, expectedDirection: -1, expectedInterval: dragAutoscrollFastInterval},
		{name: "top inner row", pointerViewportY: 1, expectedDirection: -1, expectedInterval: dragAutoscrollSlowInterval},
		{name: "middle", pointerViewportY: 5},
		{name: "bottom inner row", pointerViewportY: 8, expectedDirection: 1, expectedInterval: dragAutoscrollSlowInterval},
		{name: "bottom outer row", pointerViewportY: 9, expectedDirection: 1, expectedInterval: dragAutoscrollFastInterval},
		{name: "below view", pointerViewportY: 10, expectedDirection: 1, expectedInterval: dragAutoscrollVeryFastInterval},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			direction, interval := dragAutoscrollZone(10, testCase.pointerViewportY)

			assert.Equal(t, testCase.expectedDirection, direction)
			assert.Equal(t, testCase.expectedInterval, interval)
		})
	}
}

func TestDragAutoscrollerDoesNotRestartWhenMovingToOuterEdge(t *testing.T) {
	self := &DragAutoscroller{
		generation: 1,
		direction:  1,
		interval:   dragAutoscrollSlowInterval,
	}

	generation, schedule := self.updateState(1, dragAutoscrollFastInterval)

	assert.Equal(t, uint64(1), generation)
	assert.False(t, schedule)
	assert.Equal(t, dragAutoscrollFastInterval, self.interval)
}

func TestStaleDragAutoscrollTickDoesNotCancelCurrentGeneration(t *testing.T) {
	self := &DragAutoscroller{
		generation: 4,
		direction:  1,
		interval:   dragAutoscrollFastInterval,
	}

	self.tick(2)

	assert.Equal(t, uint64(4), self.generation)
	assert.Equal(t, 1, self.direction)
	assert.Equal(t, dragAutoscrollFastInterval, self.interval)
}
