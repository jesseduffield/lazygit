package utils

import "math"

// The configuration option for scrollheight uses partial pages when it is
// between (-1,1). This function calculates the scroll height based based of
// the window height.
func ScrollHeight(windowHeight int, scrollHeight float64) int {
	if scrollHeight == 0 || windowHeight <= 0 {
		// non-sensical value
		return 2
	} else if math.Abs(scrollHeight) <= 1 {
		// scroll partial pages
		var linesToScroll float64 = math.RoundToEven(math.Abs(scrollHeight) * float64(windowHeight))
		linesToScroll = math.Max(1, linesToScroll)
		if scrollHeight < 0 {
			return -1 * int(linesToScroll)
		} else {
			return int(linesToScroll)
		}
	} else {
		return int(math.Round(scrollHeight))
	}
}
