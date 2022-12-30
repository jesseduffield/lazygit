package snake

import (
	"math/rand"
	"time"

	"github.com/samber/lo"
)

type Game struct {
	// width/height of the board
	width  int
	height int

	// function for rendering the game. If alive is false, the cells are expected
	// to be ignored.
	render func(cells [][]CellType, alive bool)

	// closed when the game is exited
	exit chan (struct{})

	// channel for specifying the direction the player wants the snake to go in
	setNewDir chan (Direction)

	// allows logging for debugging
	logger func(string)

	// putting this on the struct for deterministic testing
	randIntFn func(int) int
}

type State struct {
	// first element is the head, final element is the tail
	snakePositions []Position

	foodPosition Position

	// direction of the snake
	direction Direction
	// direction as of the end of the last tick. We hold onto this so that
	// the snake can't do a 180 turn inbetween ticks
	lastTickDirection Direction
}

type Position struct {
	x int
	y int
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type CellType int

const (
	None CellType = iota
	Snake
	Food
)

func NewGame(width, height int, render func(cells [][]CellType, alive bool), logger func(string)) *Game {
	return &Game{
		width:     width,
		height:    height,
		render:    render,
		randIntFn: rand.Intn,
		exit:      make(chan struct{}),
		logger:    logger,
		setNewDir: make(chan Direction),
	}
}

func (self *Game) Start() {
	go self.gameLoop()
}

func (self *Game) Exit() {
	close(self.exit)
}

func (self *Game) SetDirection(direction Direction) {
	self.setNewDir <- direction
}

func (self *Game) gameLoop() {
	state := self.initializeState()
	var alive bool

	self.render(self.getCells(state), true)

	ticker := time.NewTicker(time.Duration(75) * time.Millisecond)

	for {
		select {
		case <-self.exit:
			return
		case dir := <-self.setNewDir:
			state.direction = self.newDirection(state, dir)
		case <-ticker.C:
			state, alive = self.tick(state)
			self.render(self.getCells(state), alive)
			if !alive {
				return
			}
		}
	}
}

func (self *Game) initializeState() State {
	centerOfScreen := Position{self.width / 2, self.height / 2}
	snakePositions := []Position{centerOfScreen}

	state := State{
		snakePositions: snakePositions,
		direction:      Right,
		foodPosition:   self.newFoodPos(snakePositions),
	}

	return state
}

func (self *Game) newFoodPos(snakePositions []Position) Position {
	// arbitrarily setting a limit of attempts to place food
	attemptLimit := 1000

	for i := 0; i < attemptLimit; i++ {
		newFoodPos := Position{self.randIntFn(self.width), self.randIntFn(self.height)}

		if !lo.Contains(snakePositions, newFoodPos) {
			return newFoodPos
		}
	}

	panic("SORRY, BUT I WAS TOO LAZY TO MAKE THE SNAKE GAME SMART ENOUGH TO PUT THE FOOD SOMEWHERE SENSIBLE NO MATTER WHAT, AND I ALSO WAS TOO LAZY TO ADD A WIN CONDITION")
}

// returns whether the snake is alive
func (self *Game) tick(currentState State) (State, bool) {
	nextState := currentState // copy by value
	newHeadPos := nextState.snakePositions[0]

	nextState.lastTickDirection = nextState.direction

	switch nextState.direction {
	case Up:
		newHeadPos.y--
	case Down:
		newHeadPos.y++
	case Left:
		newHeadPos.x--
	case Right:
		newHeadPos.x++
	}

	outOfBounds := newHeadPos.x < 0 || newHeadPos.x >= self.width || newHeadPos.y < 0 || newHeadPos.y >= self.height
	eatingOwnTail := lo.Contains(nextState.snakePositions, newHeadPos)

	if outOfBounds || eatingOwnTail {
		return State{}, false
	}

	nextState.snakePositions = append([]Position{newHeadPos}, nextState.snakePositions...)

	if newHeadPos == nextState.foodPosition {
		nextState.foodPosition = self.newFoodPos(nextState.snakePositions)
	} else {
		nextState.snakePositions = nextState.snakePositions[:len(nextState.snakePositions)-1]
	}

	return nextState, true
}

func (self *Game) getCells(state State) [][]CellType {
	cells := make([][]CellType, self.height)

	setCell := func(pos Position, value CellType) {
		cells[pos.y][pos.x] = value
	}

	for i := 0; i < self.height; i++ {
		cells[i] = make([]CellType, self.width)
	}

	for _, pos := range state.snakePositions {
		setCell(pos, Snake)
	}

	setCell(state.foodPosition, Food)

	return cells
}

func (self *Game) newDirection(state State, direction Direction) Direction {
	// don't allow the snake to turn 180 degrees
	if (state.lastTickDirection == Up && direction == Down) ||
		(state.lastTickDirection == Down && direction == Up) ||
		(state.lastTickDirection == Left && direction == Right) ||
		(state.lastTickDirection == Right && direction == Left) {
		return state.direction
	}

	return direction
}
