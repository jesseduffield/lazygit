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

type Cell struct {
	// strokes of the first character, which carries the vertical line
	up, down, left, right bool
	// the second character, a horizontal line joining the next cell
	connector  bool
	cellType   cellType
	rightStyle *style.TextStyle
	style      *style.TextStyle
}

func (cell *Cell) render(writer io.StringWriter) {
	first := getBoxDrawingChar(cell.up, cell.down, cell.left, cell.right)
	second := " "
	if cell.connector {
		second = "─"
	}
	var adjustedFirst string
	switch cell.cellType {
	case CONNECTION:
		adjustedFirst = first
	case COMMIT:
		adjustedFirst = string(CommitSymbol)
	case MERGE:
		adjustedFirst = string(MergeSymbol)
	}

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

	_, _ = writer.WriteString(cachedSprint(*cell.style, adjustedFirst))
	_, _ = writer.WriteString(styledSecondChar)
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

// reset keeps the connector: in the last cell a highlighted line spans, it
// leads out of that line, and the line redraws it in every other cell.
func (cell *Cell) reset() {
	cell.up = false
	cell.down = false
	cell.left = false
	cell.right = false
}

func (cell *Cell) setUp(style *style.TextStyle) *Cell {
	cell.up = true
	cell.style = style
	return cell
}

func (cell *Cell) setDown(style *style.TextStyle) *Cell {
	cell.down = true
	cell.style = style
	return cell
}

func (cell *Cell) setLeft(style *style.TextStyle) *Cell {
	cell.left = true
	if !cell.up && !cell.down {
		// vertical trumps left
		cell.style = style
	}
	return cell
}

//nolint:unparam
func (cell *Cell) setRight(style *style.TextStyle, override bool) *Cell {
	cell.right = true
	cell.connector = true
	if cell.rightStyle == nil || override {
		cell.rightStyle = style
	}
	return cell
}

func (cell *Cell) setStyle(style *style.TextStyle) *Cell {
	cell.style = style
	return cell
}

func (cell *Cell) setType(cellType cellType) *Cell {
	cell.cellType = cellType
	return cell
}

func getBoxDrawingChar(up, down, left, right bool) string {
	if up && down && left && right {
		return "│"
	} else if up && down && left && !right {
		return "│"
	} else if up && down && !left && right {
		return "│"
	} else if up && down && !left && !right {
		return "│"
	} else if up && !down && left && right {
		return "┴"
	} else if up && !down && left && !right {
		return "╯"
	} else if up && !down && !left && right {
		return "╰"
	} else if up && !down && !left && !right {
		return "╵"
	} else if !up && down && left && right {
		return "┬"
	} else if !up && down && left && !right {
		return "╮"
	} else if !up && down && !left && right {
		return "╭"
	} else if !up && down && !left && !right {
		return "╷"
	} else if !up && !down && left && right {
		return "─"
	} else if !up && !down && left && !right {
		return "─"
	} else if !up && !down && !left && right {
		return "╶"
	} else if !up && !down && !left && !right {
		return " "
	}

	panic("should not be possible")
}
