package context

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

type SearchTrait struct {
	c *ContextCommon
	*SearchHistory

	searchString string
}

func NewSearchTrait(c *ContextCommon) *SearchTrait {
	return &SearchTrait{
		c:             c,
		SearchHistory: NewSearchHistory(),
	}
}

func (self *SearchTrait) GetSearchString() string {
	return self.searchString
}

func (self *SearchTrait) SetSearchString(searchString string) {
	self.searchString = searchString
}

func (self *SearchTrait) ClearSearchString() {
	self.SetSearchString("")
}

// used for type switch
func (self *SearchTrait) IsSearchableContext() {}

func (self *SearchTrait) onSelectItemWrapper(innerFunc func(int) error) func(int, int, int) error {
	return func(selectedLineIdx int, index int, total int) error {
		self.RenderSearchStatus(index, total)

		if total != 0 {
			if err := innerFunc(selectedLineIdx); err != nil {
				return err
			}
		}

		return nil
	}
}

func (self *SearchTrait) RenderSearchStatus(index int, total int) {
	keybindingConfig := self.c.UserConfig.Keybinding

	if total == 0 {
		self.c.SetViewContent(
			self.c.Views().Search,
			fmt.Sprintf(
				self.c.Tr.NoMatchesFor,
				self.searchString,
				theme.OptionsFgColor.Sprintf(self.c.Tr.ExitSearchMode, keybindings.Label(keybindingConfig.Universal.Return)),
			),
		)
	} else {
		self.c.SetViewContent(
			self.c.Views().Search,
			fmt.Sprintf(
				self.c.Tr.MatchesFor,
				self.searchString,
				index+1,
				total,
				theme.OptionsFgColor.Sprintf(
					self.c.Tr.SearchKeybindings,
					keybindings.Label(keybindingConfig.Universal.NextMatch),
					keybindings.Label(keybindingConfig.Universal.PrevMatch),
					keybindings.Label(keybindingConfig.Universal.Return),
				),
			),
		)
	}
}

func (self *SearchTrait) IsSearching() bool {
	return self.searchString != ""
}
