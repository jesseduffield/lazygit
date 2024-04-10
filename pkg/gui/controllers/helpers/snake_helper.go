package helpers

import (
	"fmt"
	"strings"

	"github.com/lobes/lazytask/pkg/gui/style"
	"github.com/lobes/lazytask/pkg/snake"
)

type SnakeHelper struct {
	c    *HelperCommon
	game *snake.Game
}

func NewSnakeHelper(c *HelperCommon) *SnakeHelper {
	return &SnakeHelper{
		c: c,
	}
}

func (self *SnakeHelper) StartGame() {
	view := self.c.Views().Snake

	game := snake.NewGame(view.Width(), view.Height(), self.renderSnakeGame, self.c.LogAction)
	self.game = game
	game.Start()
}

func (self *SnakeHelper) ExitGame() {
	self.game.Exit()
}

func (self *SnakeHelper) SetDirection(direction snake.Direction) {
	self.game.SetDirection(direction)
}

func (self *SnakeHelper) renderSnakeGame(cells [][]snake.CellType, alive bool) {
	view := self.c.Views().Snake

	if !alive {
		_ = self.c.ErrorMsg(self.c.Tr.YouDied)
		return
	}

	output := self.drawSnakeGame(cells)

	view.Clear()
	fmt.Fprint(view, output)
	self.c.Render()
}

func (self *SnakeHelper) drawSnakeGame(cells [][]snake.CellType) string {
	writer := &strings.Builder{}

	for i, row := range cells {
		for _, cell := range row {
			switch cell {
			case snake.None:
				writer.WriteString(" ")
			case snake.Snake:
				writer.WriteString("█")
			case snake.Food:
				writer.WriteString(style.FgMagenta.Sprint("█"))
			}
		}

		if i < len(cells) {
			writer.WriteString("\n")
		}
	}

	output := writer.String()
	return output
}
