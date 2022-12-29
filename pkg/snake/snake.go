package snake

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/samber/lo"
)

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

type State struct {
	// first element is the head, final element is the tail
	snakePositions []Position
	direction      Direction
	foodPosition   Position
}

type Game struct {
	state State

	width  int
	height int
	render func(cells [][]CellType, alive bool)

	randIntFn func(int) int
}

func NewGame(width, height int, render func(cells [][]CellType, dead bool)) *Game {
	return &Game{
		width:     width,
		height:    height,
		render:    render,
		randIntFn: rand.Intn,
	}
}

func (self *Game) Start(ctx context.Context) {
	self.initializeState()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(500/self.getSpeed()) * time.Millisecond):
				fmt.Println("updating")

				alive := self.tick()
				self.render(self.getCells(), alive)
				if !alive {
					return
				}
			}
		}
	}()
}

func (self *Game) initializeState() {
	centerOfScreen := Position{self.width / 2, self.height / 2}

	self.state = State{
		snakePositions: []Position{centerOfScreen},
		direction:      Right,
	}

	self.state.foodPosition = self.setNewFoodPos()
}

// assume the player never actually wins, meaning we don't get stuck in a loop
func (self *Game) setNewFoodPos() Position {
	for i := 0; i < 1000; i++ {
		newFoodPos := Position{self.randIntFn(self.width), self.randIntFn(self.height)}

		if !lo.Contains(self.state.snakePositions, newFoodPos) {
			return newFoodPos
		}
	}

	panic("SORRY, BUT I WAS TOO LAZY TO MAKE THE SNAKE GAME SMART ENOUGH TO PUT THE FOOD SOMEWHERE SENSIBLE NO MATTER WHAT, AND I ALSO WAS TOO LAZY TO ADD A WIN CONDITION")
}

// returns whether the snake is alive
func (self *Game) tick() bool {
	newHeadPos := self.state.snakePositions[0]

	switch self.state.direction {
	case Up:
		newHeadPos.y--
	case Down:
		newHeadPos.y++
	case Left:
		newHeadPos.x--
	case Right:
		newHeadPos.x++
	}

	if newHeadPos.x < 0 || newHeadPos.x >= self.width || newHeadPos.y < 0 || newHeadPos.y >= self.height {
		return false
	}

	if lo.Contains(self.state.snakePositions, newHeadPos) {
		return false
	}

	self.state.snakePositions = append([]Position{newHeadPos}, self.state.snakePositions...)

	if newHeadPos == self.state.foodPosition {
		self.state.foodPosition = self.setNewFoodPos()
	} else {
		self.state.snakePositions = self.state.snakePositions[:len(self.state.snakePositions)-1]
	}

	return true
}

func (self *Game) getSpeed() int {
	return len(self.state.snakePositions)
}

func (self *Game) getCells() [][]CellType {
	cells := make([][]CellType, self.height)

	setCell := func(pos Position, value CellType) {
		cells[pos.y][pos.x] = value
	}

	for i := 0; i < self.height; i++ {
		cells[i] = make([]CellType, self.width)
	}

	for _, pos := range self.state.snakePositions {
		setCell(pos, Snake)
	}

	setCell(self.state.foodPosition, Food)

	return cells
}

func (self *Game) SetDirection(direction Direction) {
	self.state.direction = direction
}
