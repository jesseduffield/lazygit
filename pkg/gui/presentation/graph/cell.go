package graph

import (
	"io"
	"sync"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

const (
	MergeSymbol  = '◎'
	CommitSymbol = '○'
)

type cellType int

const (
	CONNECTION cellType = iota
	COMMIT
	MERGE
)

// How a line that touches the top or bottom edge of a cell runs within it
type verticalLine uint8

const (
	noLine verticalLine = iota
	// On to the opposite edge, or into the commit symbol
	straightLine
	// Bends towards the left edge
	lineToLeft
	// Bends towards the right edge
	lineToRight
)

type Cell struct {
	// The lines that touch the top and bottom edges
	up, down verticalLine
	// Whether lines touch the left and right edges
	left, right bool
	// Whether a line passes through from the left edge to the right edge
	horizontal bool
	cellType   cellType
	rightStyle *style.TextStyle
	style      *style.TextStyle
}

func (cell *Cell) render(writer io.StringWriter) {
	first, second := cell.boxDrawingChars()

	var rightStyle *style.TextStyle
	if cell.rightStyle == nil {
		rightStyle = cell.style
	} else {
		rightStyle = cell.rightStyle
	}

	// just doing this for the sake of easy testing, so that we don't need to
	// assert on the style of a space given a space has no styling (assuming we
	// stick to only using foreground styles)
	var styledSecondChar string
	if second == " " {
		styledSecondChar = " "
	} else {
		styledSecondChar = cachedSprint(*rightStyle, second)
	}

	_, _ = writer.WriteString(cachedSprint(*cell.style, first))
	_, _ = writer.WriteString(styledSecondChar)
}

func (cell *Cell) boxDrawingChars() (string, string) {
	first, second := getBoxDrawingChars(cell.up != noLine, cell.down != noLine, cell.left, cell.right)
	switch cell.cellType {
	case COMMIT:
		return string(CommitSymbol), second
	case MERGE:
		return string(MergeSymbol), second
	default:
		return first, second
	}
}

type rgbCacheKey struct {
	*color.RGBStyle
	str string
}

var (
	rgbCache      = make(map[rgbCacheKey]string)
	rgbCacheMutex sync.RWMutex
)

func cachedSprint(style style.TextStyle, str string) string {
	switch v := style.Style.(type) {
	case *color.RGBStyle:
		rgbCacheMutex.RLock()
		key := rgbCacheKey{v, str}
		value, ok := rgbCache[key]
		rgbCacheMutex.RUnlock()
		if ok {
			return value
		}
		value = style.Sprint(str)
		rgbCacheMutex.Lock()
		rgbCache[key] = value
		rgbCacheMutex.Unlock()
		return value
	case color.Basic:
		return style.Sprint(str)
	case color.Style:
		value := style.Sprint(str)
		return value
	}
	return style.Sprint(str)
}

func (cell *Cell) reset() {
	cell.up = noLine
	cell.down = noLine
	cell.left = false
	cell.right = false
	cell.horizontal = false
}

func (cell *Cell) setUp(style *style.TextStyle, line verticalLine) *Cell {
	cell.up = line
	cell.style = style
	return cell
}

func (cell *Cell) setDown(style *style.TextStyle, line verticalLine) *Cell {
	cell.down = line
	cell.style = style
	return cell
}

func (cell *Cell) setLeft(style *style.TextStyle) *Cell {
	cell.left = true
	if cell.up == noLine && cell.down == noLine {
		// vertical trumps left
		cell.style = style
	}
	return cell
}

func (cell *Cell) setRight(style *style.TextStyle, override bool) *Cell {
	cell.right = true
	if cell.rightStyle == nil || override {
		cell.rightStyle = style
	}
	return cell
}

func (cell *Cell) setHorizontal(style *style.TextStyle, overrideRightStyle bool) *Cell {
	cell.horizontal = true
	return cell.setLeft(style).setRight(style, overrideRightStyle)
}

func (cell *Cell) setStyle(style *style.TextStyle) *Cell {
	cell.style = style
	return cell
}

func (cell *Cell) setType(cellType cellType) *Cell {
	cell.cellType = cellType
	return cell
}

func getBoxDrawingChars(up, down, left, right bool) (string, string) {
	if up && down && left && right {
		return "│", "─"
	} else if up && down && left && !right {
		return "│", " "
	} else if up && down && !left && right {
		return "│", "─"
	} else if up && down && !left && !right {
		return "│", " "
	} else if up && !down && left && right {
		return "┴", "─"
	} else if up && !down && left && !right {
		return "╯", " "
	} else if up && !down && !left && right {
		return "╰", "─"
	} else if up && !down && !left && !right {
		return "╵", " "
	} else if !up && down && left && right {
		return "┬", "─"
	} else if !up && down && left && !right {
		return "╮", " "
	} else if !up && down && !left && right {
		return "╭", "─"
	} else if !up && down && !left && !right {
		return "╷", " "
	} else if !up && !down && left && right {
		return "─", "─"
	} else if !up && !down && left && !right {
		return "─", " "
	} else if !up && !down && !left && right {
		return "╶", "─"
	} else if !up && !down && !left && !right {
		return " ", " "
	}

	panic("should not be possible")
}
