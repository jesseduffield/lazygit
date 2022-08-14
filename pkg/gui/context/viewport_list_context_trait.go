package context

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// This embeds a list context trait and adds logic to re-render the viewport
// whenever a line is focused. We use this in the commits panel because different
// sections of the log graph need to be highlighted depending on the currently selected line

type ViewportListContextTrait struct {
	*ListContextTrait
}

func (self *ViewportListContextTrait) FocusLine() {
	self.ListContextTrait.FocusLine()

	startIdx, length := self.GetViewTrait().ViewPortYBounds()
	displayStrings := self.ListContextTrait.getDisplayStrings(startIdx, length)
	content := utils.RenderDisplayStrings(displayStrings)
	self.GetViewTrait().SetViewPortContent(content)
}
