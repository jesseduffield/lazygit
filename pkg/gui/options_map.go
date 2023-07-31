package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type OptionsMapMgr struct {
	c *helpers.HelperCommon
}

func (gui *Gui) renderContextOptionsMap(c types.Context) {
	// In demos, we render our own content to this view
	if gui.integrationTest != nil && gui.integrationTest.IsDemo() {
		return
	}
	mgr := OptionsMapMgr{c: gui.c}
	mgr.renderContextOptionsMap(c)
}

// render the options available for the current context at the bottom of the screen
func (self *OptionsMapMgr) renderContextOptionsMap(c types.Context) {
	bindingsToDisplay := lo.Filter(c.GetKeybindings(self.c.KeybindingsOpts()), func(binding *types.Binding, _ int) bool {
		return binding.Display
	})

	var optionsMap []bindingInfo
	if len(bindingsToDisplay) == 0 {
		optionsMap = self.globalOptions()
	} else {
		optionsMap = lo.Map(bindingsToDisplay, func(binding *types.Binding, _ int) bindingInfo {
			return bindingInfo{
				key:         keybindings.LabelFromKey(binding.Key),
				description: binding.Description,
			}
		})
	}

	self.renderOptions(self.formatBindingInfos(optionsMap))
}

func (self *OptionsMapMgr) formatBindingInfos(bindingInfos []bindingInfo) string {
	return strings.Join(
		lo.Map(bindingInfos, func(bindingInfo bindingInfo, _ int) string {
			return fmt.Sprintf("%s: %s", bindingInfo.key, bindingInfo.description)
		}),
		", ")
}

func (self *OptionsMapMgr) renderOptions(options string) {
	self.c.SetViewContent(self.c.Views().Options, options)
}

func (self *OptionsMapMgr) globalOptions() []bindingInfo {
	keybindingConfig := self.c.UserConfig.Keybinding

	return []bindingInfo{
		{
			key:         fmt.Sprintf("%s/%s", keybindings.Label(keybindingConfig.Universal.ScrollUpMain), keybindings.Label(keybindingConfig.Universal.ScrollDownMain)),
			description: self.c.Tr.Scroll,
		},
		{
			key:         keybindings.Label(keybindingConfig.Universal.Return),
			description: self.c.Tr.Cancel,
		},
		{
			key:         keybindings.Label(keybindingConfig.Universal.Quit),
			description: self.c.Tr.Quit,
		},
		{
			key:         keybindings.Label(keybindingConfig.Universal.OptionMenuAlt1),
			description: self.c.Tr.Keybindings,
		},
		{
			key:         fmt.Sprintf("%s-%s", keybindings.Label(keybindingConfig.Universal.JumpToBlock[0]), keybindings.Label(keybindingConfig.Universal.JumpToBlock[len(keybindingConfig.Universal.JumpToBlock)-1])),
			description: self.c.Tr.Jump,
		},
		{
			key:         fmt.Sprintf("%s/%s", keybindings.Label(keybindingConfig.Universal.ScrollLeft), keybindings.Label(keybindingConfig.Universal.ScrollRight)),
			description: self.c.Tr.ScrollLeftRight,
		},
	}
}

type bindingInfo struct {
	key         string
	description string
}
