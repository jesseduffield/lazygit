package snake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnake(t *testing.T) {
	scenarios := []struct {
		state         State
		expectedState State
		expectedAlive bool
	}{
		{
			state: State{
				snakePositions:    []Position{{x: 5, y: 5}},
				direction:         Right,
				lastTickDirection: Right,
				foodPosition:      Position{x: 9, y: 9},
			},
			expectedState: State{
				snakePositions:    []Position{{x: 6, y: 5}},
				direction:         Right,
				lastTickDirection: Right,
				foodPosition:      Position{x: 9, y: 9},
			},
			expectedAlive: true,
		},
		{
			state: State{
				snakePositions:    []Position{{x: 5, y: 5}, {x: 4, y: 5}, {x: 4, y: 4}, {x: 5, y: 4}},
				direction:         Up,
				lastTickDirection: Up,
				foodPosition:      Position{x: 9, y: 9},
			},
			expectedState: State{},
			expectedAlive: false,
		},
		{
			state: State{
				snakePositions:    []Position{{x: 5, y: 5}},
				direction:         Right,
				lastTickDirection: Right,
				foodPosition:      Position{x: 6, y: 5},
			},
			expectedState: State{
				snakePositions:    []Position{{x: 6, y: 5}, {x: 5, y: 5}},
				direction:         Right,
				lastTickDirection: Right,
				foodPosition:      Position{x: 8, y: 8},
			},
			expectedAlive: true,
		},
	}

	for _, scenario := range scenarios {
		game := NewGame(10, 10, nil, func(string) {})
		game.randIntFn = func(int) int { return 8 }
		state, alive := game.tick(scenario.state)
		assert.Equal(t, scenario.expectedAlive, alive)
		assert.EqualValues(t, scenario.expectedState, state)
	}
}
