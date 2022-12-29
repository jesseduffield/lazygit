package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/snake"
)

func main() {
	game := snake.NewGame(10, 10, render)
	ctx := context.Background()
	game.Start(ctx)

	go func() {
		for {
			var input string
			fmt.Scanln(&input)

			switch input {
			case "w":
				game.SetDirection(snake.Up)
			case "s":
				game.SetDirection(snake.Down)
			case "a":
				game.SetDirection(snake.Left)
			case "d":
				game.SetDirection(snake.Right)
			}
		}
	}()

	time.Sleep(100 * time.Second)
}

func render(cells [][]snake.CellType, alive bool) {
	if !alive {
		log.Fatal("YOU DIED!")
	}

	writer := &strings.Builder{}

	width := len(cells[0])

	writer.WriteString(strings.Repeat("\n", 20))

	writer.WriteString(strings.Repeat("█", width+2) + "\n")

	for _, row := range cells {
		writer.WriteString("█")

		for _, cell := range row {
			switch cell {
			case snake.None:
				writer.WriteString(" ")
			case snake.Snake:
				writer.WriteString("X")
			case snake.Food:
				writer.WriteString("o")
			}
		}

		writer.WriteString("█")

		writer.WriteString("\n")
	}

	writer.WriteString(strings.Repeat("█", width+2))

	fmt.Println(writer.String())
}
