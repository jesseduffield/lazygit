package graph

import (
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

const mergeSymbol = '⏣'
const commitSymbol = '⎔'

type cellType int

const (
	CONNECTION cellType = iota
	COMMIT
	MERGE
)

type Cell struct {
	up, down, left, right bool
	cellType              cellType
	rightStyle            *style.TextStyle
	style                 style.TextStyle
}

func (cell *Cell) render() string {
	up, down, left, right := cell.up, cell.down, cell.left, cell.right

	first, second := getBoxDrawingChars(up, down, left, right)
	var adjustedFirst rune
	switch cell.cellType {
	case CONNECTION:
		adjustedFirst = first
	case COMMIT:
		adjustedFirst = commitSymbol
	case MERGE:
		adjustedFirst = mergeSymbol
	}

	var rightStyle *style.TextStyle
	if cell.rightStyle == nil {
		rightStyle = &cell.style
	} else {
		rightStyle = cell.rightStyle
	}

	// just doing this for the sake of easy testing, so that we don't need to
	// assert on the style of a space given a space has no styling (assuming we
	// stick to only using foreground styles)
	var styledSecondChar string
	if second == ' ' {
		styledSecondChar = " "
	} else {
		styledSecondChar = rightStyle.Sprint(string(second))
	}

	return cell.style.Sprint(string(adjustedFirst)) + styledSecondChar
}

func (cell *Cell) reset() {
	cell.up = false
	cell.down = false
	cell.left = false
	cell.right = false
}

func (cell *Cell) setUp(style style.TextStyle) *Cell {
	cell.up = true
	cell.style = style
	return cell
}

func (cell *Cell) setDown(style style.TextStyle) *Cell {
	cell.down = true
	cell.style = style
	return cell
}

func (cell *Cell) setLeft(style style.TextStyle) *Cell {
	cell.left = true
	if !cell.up && !cell.down {
		// vertical trumps left
		cell.style = style
	}
	return cell
}

func (cell *Cell) setRight(style style.TextStyle, override bool) *Cell {
	cell.right = true
	if cell.rightStyle == nil || override {
		cell.rightStyle = &style
	}
	return cell
}

func (cell *Cell) setStyle(style style.TextStyle) *Cell {
	cell.style = style
	return cell
}

func (cell *Cell) setType(cellType cellType) *Cell {
	cell.cellType = cellType
	return cell
}

func getBoxDrawingChars(up, down, left, right bool) (rune, rune) {
	if up && down && left && right {
		return '│', '─'
	} else if up && down && left && !right {
		return '│', ' '
	} else if up && down && !left && right {
		return '│', '─'
	} else if up && down && !left && !right {
		return '│', ' '
	} else if up && !down && left && right {
		return '┴', '─'
	} else if up && !down && left && !right {
		return '╯', ' '
	} else if up && !down && !left && right {
		return '╰', '─'
	} else if up && !down && !left && !right {
		return '╵', ' '
	} else if !up && down && left && right {
		return '┬', '─'
	} else if !up && down && left && !right {
		return '╮', ' '
	} else if !up && down && !left && right {
		return '╭', '─'
	} else if !up && down && !left && !right {
		return '╷', ' '
	} else if !up && !down && left && right {
		return '─', '─'
	} else if !up && !down && left && !right {
		return '─', ' '
	} else if !up && !down && !left && right {
		return '╶', '─'
	} else if !up && !down && !left && !right {
		return ' ', ' '
	} else {
		panic("should not be possible")
	}
}
