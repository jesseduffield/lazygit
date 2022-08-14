package gocui

import "math"

// returns start and height of scrollbar
// `max` is the maximum possible value of `position`
func calcScrollbar(listSize int, pageSize int, position int, scrollAreaSize int) (int, int) {
	height := calcScrollbarHeight(listSize, pageSize, scrollAreaSize)
	// assume we can't scroll past the last item
	maxPosition := listSize - pageSize
	if maxPosition <= 0 {
		return 0, height
	}
	if position == maxPosition {
		return scrollAreaSize - height, height
	}
	// we only want to show the scrollbar at the top or bottom positions if we're at the end. Hence the .Ceil (for moving the scrollbar once we scroll down) and the -1 (for pretending there's a smaller range than we actually have, with the above condition ensuring we snap to the bottom once we're at the end of the list)
	start := int(math.Ceil(((float64(position) / float64(maxPosition)) * float64(scrollAreaSize-height-1))))
	return start, height
}

func calcScrollbarHeight(listSize int, pageSize int, scrollAreaSize int) int {
	if pageSize >= listSize {
		return scrollAreaSize
	}
	height := int((float64(pageSize) / float64(listSize)) * float64(scrollAreaSize))
	minHeight := 2
	if height < minHeight {
		return minHeight
	}

	return height
}
