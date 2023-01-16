package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/snake"
)

func (gui *Gui) startSnake() {
	view := gui.Views.Snake

	game := snake.NewGame(view.Width(), view.Height(), gui.renderSnakeGame, gui.c.LogAction)
	gui.snakeGame = game
	game.Start()
}

func (gui *Gui) renderSnakeGame(cells [][]snake.CellType, alive bool) {
	view := gui.Views.Snake

	if !alive {
		_ = gui.c.ErrorMsg(gui.Tr.YouDied)
		return
	}

	output := drawSnakeGame(cells)

	view.Clear()
	fmt.Fprint(view, output)
	gui.c.Render()
}

func drawSnakeGame(cells [][]snake.CellType) string {
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
