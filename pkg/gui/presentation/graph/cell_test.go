package graph

import (
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/stretchr/testify/assert"
)

func TestResetRGBCacheDropsEntriesOfOldStyles(t *testing.T) {
	oldStyle := style.New().SetFg(style.NewRGBColor(color.RGB(0, 255, 0)))
	newStyle := style.New().SetFg(style.NewRGBColor(color.RGB(255, 0, 0)))

	cachedSprint(oldStyle, "│")
	ResetRGBCache()
	cachedSprint(newStyle, "│")

	// Only what was rendered after the reset may be left, whatever was cached
	// before it, including entries that other tests of this package added.
	expected := map[rgbCacheKey]string{
		{newStyle.Style.(*color.RGBStyle), "│"}: newStyle.Sprint("│"),
	}
	assert.Equal(t, expected, rgbCache)
}
