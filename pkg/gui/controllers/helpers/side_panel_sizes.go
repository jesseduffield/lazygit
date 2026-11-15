package helpers

import (
	"math"
	"slices"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/samber/lo"
)

// minSideSectionSize is the smallest size of the side section, and of the main
// section, when the user has set the size of the side panels. It leaves room
// for one line or column inside the frame.
const minSideSectionSize = 3

// sidePanelSizeForScreen returns the size of the side panels for a screen of the
// given size, based on the sizes that the user set at other screen sizes (sorted
// by ScreenSize). Between two of those screen sizes, it interpolates. Beyond the
// smallest or the largest one, it keeps the ratio of that one. It returns false
// if the user hasn't set any.
func sidePanelSizeForScreen(sizes []config.SidePanelSize, screenSize int) (int, bool) {
	if len(sizes) == 0 {
		return 0, false
	}

	sizeFromRatio := func(ratio float64) int {
		return int(math.Round(ratio * float64(screenSize)))
	}
	first, last := sizes[0], sizes[len(sizes)-1]
	if screenSize <= first.ScreenSize {
		return sizeFromRatio(first.Ratio), true
	}
	if screenSize >= last.ScreenSize {
		return sizeFromRatio(last.Ratio), true
	}

	// We interpolate the size rather than the ratio. Interpolating the ratio
	// can result in a size that is outside the range of the two sizes, e.g.
	// wider than both of them.
	i := slices.IndexFunc(sizes, func(size config.SidePanelSize) bool {
		return size.ScreenSize >= screenSize
	})
	lower, upper := sizes[i-1], sizes[i]
	lowerSize := lower.Ratio * float64(lower.ScreenSize)
	upperSize := upper.Ratio * float64(upper.ScreenSize)
	t := float64(screenSize-lower.ScreenSize) / float64(upper.ScreenSize-lower.ScreenSize)
	return int(math.Round(lowerSize + t*(upperSize-lowerSize))), true
}

// recordSidePanelSize returns sizes plus the size that the user set at the given
// screen size. It replaces the sizes for similar screen sizes, so that e.g.
// zooming the terminal's font in or out by a step doesn't keep adding sizes.
func recordSidePanelSize(sizes []config.SidePanelSize, screenSize int, size int) []config.SidePanelSize {
	result := lo.Reject(sizes, func(existing config.SidePanelSize, _ int) bool {
		return isSimilarScreenSize(existing.ScreenSize, screenSize)
	})
	// Four decimals are enough to get the same size back for any realistic
	// screen size, and they keep state.yml readable.
	ratio := math.Round(float64(size)/float64(screenSize)*10000) / 10000
	result = append(result, config.SidePanelSize{ScreenSize: screenSize, Ratio: ratio})
	slices.SortFunc(result, func(a, b config.SidePanelSize) int {
		return a.ScreenSize - b.ScreenSize
	})
	return result
}

// isSimilarScreenSize returns true if the two screen sizes differ by at most 15%.
func isSimilarScreenSize(a int, b int) bool {
	return float64(max(a, b)) <= float64(min(a, b))*1.15
}

func clampSideSectionSize(size int, screenSize int) int {
	return max(minSideSectionSize, min(size, screenSize-minSideSectionSize))
}

// ValidSidePanelSizes returns the side panel sizes from state.yml without the
// ones that can't be used (e.g. because the file was edited by hand), sorted by
// screen size as sidePanelSizeForScreen expects them.
func ValidSidePanelSizes(
	sizesByLayout map[config.SidePanelLayout][]config.SidePanelSize,
) map[config.SidePanelLayout][]config.SidePanelSize {
	result := map[config.SidePanelLayout][]config.SidePanelSize{}
	for layout, sizes := range sizesByLayout {
		valid := lo.Filter(sizes, func(size config.SidePanelSize, _ int) bool {
			return size.ScreenSize > 0 && size.Ratio > 0 && size.Ratio < 1
		})
		slices.SortStableFunc(valid, func(a, b config.SidePanelSize) int {
			return a.ScreenSize - b.ScreenSize
		})
		// Interpolating between two sizes for the same screen size would divide
		// by zero, so keep only the last one.
		valid = lo.Filter(valid, func(size config.SidePanelSize, i int) bool {
			return i == len(valid)-1 || valid[i+1].ScreenSize != size.ScreenSize
		})
		if len(valid) > 0 {
			result[layout] = valid
		}
	}
	return result
}
