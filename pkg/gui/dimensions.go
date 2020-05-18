package gui

import "math"

type dimensions struct {
	x0 int
	x1 int
	y0 int
	y1 int
}

const (
	ROW = iota
	COLUMN
)

// to give a high-level explanation of what's going on here. We layout our views by arranging a bunch of boxes in the window.
// If a box has children, it needs to specify how it wants to arrange those children: ROW or COLUMN.
// If a box represents a view, you can put the view name in the viewName field.
// When determining how to divvy-up the available height (for row children) or width (for column children), we first
// give the boxes with a static `size` the space that they want. Then we apportion
// the remaining space based on the weights of the dynamic boxes (you can't define
// both size and weight at the same time: you gotta pick one). If there are two
// boxes, one with weight 1 and the other with weight 2, the first one gets 33%
// of the available space and the second one gets the remaining 66%

type box struct {
	// direction decides how the children boxes are laid out. ROW means the children will each form a row i.e. that they will be stacked on top of eachother.
	direction int // ROW or COLUMN

	// function which takes the width and height assigned to the box and decides which orientation it will have
	conditionalDirection func(width int, height int) int

	children []*box

	// function which takes the width and height assigned to the box and decides the layout of the children.
	conditionalChildren func(width int, height int) []*box

	// viewName refers to the name of the view this box represents, if there is one
	viewName string

	// static size. If parent box's direction is ROW this refers to height, otherwise width
	size int

	// dynamic size. Once all statically sized children have been considered, weight decides how much of the remaining space will be taken up by the box
	// TODO: consider making there be one int and a type enum so we can't have size and weight simultaneously defined
	weight int
}

func (b *box) isStatic() bool {
	return b.size > 0
}

func (b *box) getDirection(width int, height int) int {
	if b.conditionalDirection != nil {
		return b.conditionalDirection(width, height)
	}
	return b.direction
}

func (b *box) getChildren(width int, height int) []*box {
	if b.conditionalChildren != nil {
		return b.conditionalChildren(width, height)
	}
	return b.children
}

func (gui *Gui) arrangeViews(root *box, x0, y0, width, height int) map[string]dimensions {
	children := root.getChildren(width, height)
	if len(children) == 0 {
		// leaf node
		if root.viewName != "" {
			dimensionsForView := dimensions{x0: x0, y0: y0, x1: x0 + width - 1, y1: y0 + height - 1}
			return map[string]dimensions{root.viewName: dimensionsForView}
		}
		return map[string]dimensions{}
	}

	direction := root.getDirection(width, height)

	var availableSize int
	if direction == COLUMN {
		availableSize = width
	} else {
		availableSize = height
	}

	// work out size taken up by children
	reservedSize := 0
	totalWeight := 0
	for _, child := range children {
		// assuming either size or weight are non-zero
		reservedSize += child.size
		totalWeight += child.weight
	}

	remainingSize := availableSize - reservedSize
	if remainingSize < 0 {
		remainingSize = 0
	}

	unitSize := 0
	extraSize := 0
	if totalWeight > 0 {
		unitSize = remainingSize / totalWeight
		extraSize = remainingSize % totalWeight
	}

	result := map[string]dimensions{}
	offset := 0
	for _, child := range children {
		var boxSize int
		if child.isStatic() {
			boxSize = child.size
		} else {
			// TODO: consider more evenly distributing the remainder
			boxSize = unitSize * child.weight
			boxExtraSize := int(math.Min(float64(extraSize), float64(child.weight)))
			boxSize += boxExtraSize
			extraSize -= boxExtraSize
		}

		var resultForChild map[string]dimensions
		if direction == COLUMN {
			resultForChild = gui.arrangeViews(child, x0+offset, y0, boxSize, height)
		} else {
			resultForChild = gui.arrangeViews(child, x0, y0+offset, width, boxSize)
		}

		result = gui.mergeDimensionMaps(result, resultForChild)
		offset += boxSize
	}

	return result
}

func (gui *Gui) mergeDimensionMaps(a map[string]dimensions, b map[string]dimensions) map[string]dimensions {
	result := map[string]dimensions{}
	for _, dimensionMap := range []map[string]dimensions{a, b} {
		for k, v := range dimensionMap {
			result[k] = v
		}
	}
	return result
}
