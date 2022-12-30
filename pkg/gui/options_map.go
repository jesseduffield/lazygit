package gui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jesseduffield/generics/maps"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type OptionsMapMgr struct {
	c *types.HelperCommon
}

func (gui *Gui) renderContextOptionsMap(c types.Context) {
	mgr := OptionsMapMgr{c: gui.c}
	mgr.renderContextOptionsMap(c)
}

// render the options available for the current context at the bottom of the screen
func (self *OptionsMapMgr) renderContextOptionsMap(c types.Context) {
	optionsMap := c.GetOptionsMap()
	if optionsMap == nil {
		optionsMap = self.globalOptionsMap()
	}

	self.renderOptions(self.optionsMapToString(optionsMap))
}

func (self *OptionsMapMgr) optionsMapToString(optionsMap map[string]string) string {
	options := maps.MapToSlice(optionsMap, func(key string, description string) string {
		return key + ": " + description
	})
	sort.Strings(options)
	return strings.Join(options, ", ")
}

func (self *OptionsMapMgr) renderOptions(options string) {
	self.c.SetViewContent(self.c.Views().Options, options)
}

func (self *OptionsMapMgr) globalOptionsMap() map[string]string {
	keybindingConfig := self.c.UserConfig.Keybinding

	return map[string]string{
		fmt.Sprintf("%s/%s", keybindings.Label(keybindingConfig.Universal.ScrollUpMain), keybindings.Label(keybindingConfig.Universal.ScrollDownMain)):                                                                                                               self.c.Tr.LcScroll,
		fmt.Sprintf("%s %s %s %s", keybindings.Label(keybindingConfig.Universal.PrevBlock), keybindings.Label(keybindingConfig.Universal.NextBlock), keybindings.Label(keybindingConfig.Universal.PrevItem), keybindings.Label(keybindingConfig.Universal.NextItem)): self.c.Tr.LcNavigate,
		keybindings.Label(keybindingConfig.Universal.Return):         self.c.Tr.LcCancel,
		keybindings.Label(keybindingConfig.Universal.Quit):           self.c.Tr.LcQuit,
		keybindings.Label(keybindingConfig.Universal.OptionMenuAlt1): self.c.Tr.LcMenu,
		fmt.Sprintf("%s-%s", keybindings.Label(keybindingConfig.Universal.JumpToBlock[0]), keybindings.Label(keybindingConfig.Universal.JumpToBlock[len(keybindingConfig.Universal.JumpToBlock)-1])): self.c.Tr.LcJump,
		fmt.Sprintf("%s/%s", keybindings.Label(keybindingConfig.Universal.ScrollLeft), keybindings.Label(keybindingConfig.Universal.ScrollRight)):                                                    self.c.Tr.LcScrollLeftRight,
	}
}
