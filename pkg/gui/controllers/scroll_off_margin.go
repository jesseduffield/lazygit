package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// To be called after pressing up-arrow; checks whether the cursor entered the
// top scroll-off margin, and so the view needs to be scrolled up one line
func checkScrollUp(view types.IViewTrait, userConfig *config.UserConfig, lineIdxBefore int, lineIdxAfter int) {
	if userConfig.Gui.ScrollOffBehavior != "jump" {
		viewPortStart, viewPortHeight := view.ViewPortYBounds()

		linesToScroll := calculateLinesToScrollUp(
			viewPortStart, viewPortHeight, userConfig.Gui.ScrollOffMargin, lineIdxBefore, lineIdxAfter)
		if linesToScroll != 0 {
			view.ScrollUp(linesToScroll)
		}
	}
}

// To be called after pressing down-arrow; checks whether the cursor entered the
// bottom scroll-off margin, and so the view needs to be scrolled down one line
func checkScrollDown(view types.IViewTrait, userConfig *config.UserConfig, lineIdxBefore int, lineIdxAfter int) {
	if userConfig.Gui.ScrollOffBehavior != "jump" {
		viewPortStart, viewPortHeight := view.ViewPortYBounds()

		linesToScroll := calculateLinesToScrollDown(
			viewPortStart, viewPortHeight, userConfig.Gui.ScrollOffMargin, lineIdxBefore, lineIdxAfter)
		if linesToScroll != 0 {
			view.ScrollDown(linesToScroll)
		}
	}
}

func calculateLinesToScrollUp(viewPortStart int, viewPortHeight int, scrollOffMargin int, lineIdxBefore int, lineIdxAfter int) int {
	// Cap the margin to half the view height. This allows setting the config to
	// a very large value to keep the cursor always in the middle of the screen.
	// Use +.5 so that if the height is even, the top margin is one line higher
	// than the bottom margin.
	scrollOffMargin = utils.Min(scrollOffMargin, int((float64(viewPortHeight)+.5)/2))

	// Scroll only if the "before" position was visible (this could be false if
	// the scroll wheel was used to scroll the selected line out of view) ...
	if lineIdxBefore >= viewPortStart && lineIdxBefore < viewPortStart+viewPortHeight {
		marginEnd := viewPortStart + scrollOffMargin
		// ... and the "after" position is within the top margin (or before it)
		if lineIdxAfter < marginEnd {
			return marginEnd - lineIdxAfter
		}
	}

	return 0
}

func calculateLinesToScrollDown(viewPortStart int, viewPortHeight int, scrollOffMargin int, lineIdxBefore int, lineIdxAfter int) int {
	// Cap the margin to half the view height. This allows setting the config to
	// a very large value to keep the cursor always in the middle of the screen.
	// Use -.5 so that if the height is even, the bottom margin is one line lower
	// than the top margin.
	scrollOffMargin = utils.Min(scrollOffMargin, int((float64(viewPortHeight)-.5)/2))

	// Scroll only if the "before" position was visible (this could be false if
	// the scroll wheel was used to scroll the selected line out of view) ...
	if lineIdxBefore >= viewPortStart && lineIdxBefore < viewPortStart+viewPortHeight {
		marginStart := viewPortStart + viewPortHeight - scrollOffMargin - 1
		// ... and the "after" position is within the bottom margin (or after it)
		if lineIdxAfter > marginStart {
			return lineIdxAfter - marginStart
		}
	}

	return 0
}
