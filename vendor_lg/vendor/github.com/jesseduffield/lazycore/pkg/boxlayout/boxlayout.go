package boxlayout

import (
	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/samber/lo"
)

type Dimensions struct {
	X0 int
	X1 int
	Y0 int
	Y1 int
}

type Direction int

const (
	ROW Direction = iota
	COLUMN
)

// to give a high-level explanation of what's going on here. We layout our windows by arranging a bunch of boxes in the available space.
// If a box has children, it needs to specify how it wants to arrange those children: ROW or COLUMN.
// If a box represents a window, you can put the window name in the Window field.
// When determining how to divvy-up the available height (for row children) or width (for column children), we first
// give the boxes with a static `size` the space that they want. Then we apportion
// the remaining space based on the weights of the dynamic boxes (you can't define
// both size and weight at the same time: you gotta pick one). If there are two
// boxes, one with weight 1 and the other with weight 2, the first one gets 33%
// of the available space and the second one gets the remaining 66%

type Box struct {
	// Direction decides how the children boxes are laid out. ROW means the children will each form a row i.e. that they will be stacked on top of eachother.
	Direction Direction

	// function which takes the width and height assigned to the box and decides which orientation it will have
	ConditionalDirection func(width int, height int) Direction

	Children []*Box

	// function which takes the width and height assigned to the box and decides the layout of the children.
	ConditionalChildren func(width int, height int) []*Box

	// Window refers to the name of the window this box represents, if there is one
	Window string

	// static Size. If parent box's direction is ROW this refers to height, otherwise width
	Size int

	// dynamic size. Once all statically sized children have been considered, Weight decides how much of the remaining space will be taken up by the box
	// TODO: consider making there be one int and a type enum so we can't have size and Weight simultaneously defined
	Weight int
}

func ArrangeWindows(root *Box, x0, y0, width, height int) map[string]Dimensions {
	children := root.getChildren(width, height)
	if len(children) == 0 {
		// leaf node
		if root.Window != "" {
			dimensionsForWindow := Dimensions{X0: x0, Y0: y0, X1: x0 + width - 1, Y1: y0 + height - 1}
			return map[string]Dimensions{root.Window: dimensionsForWindow}
		}
		return map[string]Dimensions{}
	}

	direction := root.getDirection(width, height)

	var availableSize int
	if direction == COLUMN {
		availableSize = width
	} else {
		availableSize = height
	}

	sizes := calcSizes(children, availableSize)

	result := map[string]Dimensions{}
	offset := 0
	for i, child := range children {
		boxSize := sizes[i]

		var resultForChild map[string]Dimensions
		if direction == COLUMN {
			resultForChild = ArrangeWindows(child, x0+offset, y0, boxSize, height)
		} else {
			resultForChild = ArrangeWindows(child, x0, y0+offset, width, boxSize)
		}

		result = mergeDimensionMaps(result, resultForChild)
		offset += boxSize
	}

	return result
}

func calcSizes(boxes []*Box, availableSpace int) []int {
	normalizedWeights := normalizeWeights(lo.Map(boxes, func(box *Box, _ int) int { return box.Weight }))

	totalWeight := 0
	reservedSpace := 0
	for i, box := range boxes {
		if box.isStatic() {
			reservedSpace += box.Size
		} else {
			totalWeight += normalizedWeights[i]
		}
	}

	dynamicSpace := utils.Max(0, availableSpace-reservedSpace)

	unitSize := 0
	extraSpace := 0
	if totalWeight > 0 {
		unitSize = dynamicSpace / totalWeight
		extraSpace = dynamicSpace % totalWeight
	}

	result := make([]int, len(boxes))
	for i, box := range boxes {
		if box.isStatic() {
			// assuming that only one static child can have a size greater than the
			// available space. In that case we just crop the size to what's available
			result[i] = utils.Min(availableSpace, box.Size)
		} else {
			result[i] = unitSize * normalizedWeights[i]
		}
	}

	// distribute the remainder across dynamic boxes.
	for extraSpace > 0 {
		for i, weight := range normalizedWeights {
			if weight > 0 {
				result[i]++
				extraSpace--
				normalizedWeights[i]--

				if extraSpace == 0 {
					break
				}
			}
		}
	}

	return result
}

// removes common multiple from weights e.g. if we get 2, 4, 4 we return 1, 2, 2.
func normalizeWeights(weights []int) []int {
	if len(weights) == 0 {
		return []int{}
	}

	// to spare us some computation we'll exit early if any of our weights is 1
	if lo.SomeBy(weights, func(weight int) bool { return weight == 1 }) {
		return weights
	}

	// map weights to factorSlices and find the lowest common factor
	positiveWeights := lo.Filter(weights, func(weight int, _ int) bool { return weight > 0 })
	factorSlices := lo.Map(positiveWeights, func(weight int, _ int) []int { return calcFactors(weight) })
	commonFactors := factorSlices[0]
	for _, factors := range factorSlices {
		commonFactors = lo.Intersect(commonFactors, factors)
	}

	if len(commonFactors) == 0 {
		return weights
	}

	newWeights := lo.Map(weights, func(weight int, _ int) int { return weight / commonFactors[0] })

	return normalizeWeights(newWeights)
}

func calcFactors(n int) []int {
	factors := []int{}
	for i := 2; i <= n; i++ {
		if n%i == 0 {
			factors = append(factors, i)
		}
	}
	return factors
}

func (b *Box) isStatic() bool {
	return b.Size > 0
}

func (b *Box) getDirection(width int, height int) Direction {
	if b.ConditionalDirection != nil {
		return b.ConditionalDirection(width, height)
	}
	return b.Direction
}

func (b *Box) getChildren(width int, height int) []*Box {
	if b.ConditionalChildren != nil {
		return b.ConditionalChildren(width, height)
	}
	return b.Children
}

func mergeDimensionMaps(a map[string]Dimensions, b map[string]Dimensions) map[string]Dimensions {
	result := map[string]Dimensions{}
	for _, dimensionMap := range []map[string]Dimensions{a, b} {
		for k, v := range dimensionMap {
			result[k] = v
		}
	}
	return result
}
