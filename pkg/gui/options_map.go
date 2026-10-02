package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type OptionsMapMgr struct {
	c                     *helpers.HelperCommon
	customCommandBindings []*types.Binding
}

func (gui *Gui) renderContextOptionsMap() {
	// In demos, we render our own content to this view
	if gui.integrationTest != nil && gui.integrationTest.IsDemo() {
		return
	}

	mgr := OptionsMapMgr{c: gui.c, customCommandBindings: gui.customCommandBindings}
	mgr.renderContextOptionsMap()
}

// Render the options available for the current context at the bottom of the screen
// STYLE GUIDE: we use the default options fg color for most keybindings. We can
// only use a different color if we're in a specific mode where the user is likely
// to want to press that key. For example, when in cherry-picking mode, we
// want to prominently show the keybinding for pasting commits.
func (self *OptionsMapMgr) renderContextOptionsMap() {
	currentContext := self.c.Context().Current()

	currentContextBindings := currentContext.GetKeybindings(self.c.KeybindingsOpts())
	globalBindings := self.c.Contexts().Global.GetKeybindings(self.c.KeybindingsOpts())
	customCommandBindings := self.customCommandBindings
	if currentContext.GetKey() == context.SEARCH_CONTEXT_KEY {
		customCommandBindings = nil
	}
	customGlobalBindings := []*types.Binding{}

	if len(customCommandBindings) > 0 {
		currentViewName := currentContext.GetViewName()
		customContextBindings := lo.Filter(customCommandBindings, func(b *types.Binding, _ int) bool {
			return b.ViewName == currentViewName
		})
		currentContextBindings = append(customContextBindings, currentContextBindings...)
	}

	currentContextKeys := set.NewFromSlice(
		lo.FlatMap(currentContextBindings, func(binding *types.Binding, _ int) []gocui.Key {
			return binding.Keys
		}))
	if len(customCommandBindings) > 0 {
		customGlobalBindings = lo.Filter(customCommandBindings, func(binding *types.Binding, _ int) bool {
			return binding.ViewName == "" && len(binding.Keys) > 0 && !currentContextKeys.Includes(binding.Keys[0])
		})
	}

	customGlobalKeys := set.NewFromSlice(
		lo.FlatMap(customGlobalBindings, func(binding *types.Binding, _ int) []gocui.Key {
			return binding.Keys
		}))

	allBindings := append(customGlobalBindings, currentContextBindings...)
	allBindings = append(allBindings, lo.Filter(globalBindings, func(binding *types.Binding, _ int) bool {
		return len(binding.Keys) > 0 && !currentContextKeys.Includes(binding.Keys[0]) && !customGlobalKeys.Includes(binding.Keys[0])
	})...)

	bindingsToDisplay := lo.Filter(allBindings, func(binding *types.Binding, _ int) bool {
		return len(binding.Keys) > 0 && binding.DisplayOnScreen && !binding.IsDisabled()
	})

	optionsMap := lo.Map(bindingsToDisplay, func(binding *types.Binding, _ int) bindingInfo {
		displayStyle := theme.OptionsFgColor
		if binding.DisplayStyle != nil {
			displayStyle = *binding.DisplayStyle
		}

		return bindingInfo{
			key:         config.LabelForKey(binding.Keys[0]),
			description: binding.GetShortDescription(),
			style:       displayStyle,
		}
	})

	// Mode-specific local keybindings
	if currentContext.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
		if self.c.Modes().CherryPicking.Active() {
			optionsMap = utils.Prepend(optionsMap, bindingInfo{
				key:         self.c.KeybindingsOpts().Config.Commits.PasteCommits.String(),
				description: self.c.Tr.PasteCommits,
				style:       style.FgCyan,
			})
		}

		if self.c.Model().BisectInfo.Started() {
			optionsMap = utils.Prepend(optionsMap, bindingInfo{
				key:         self.c.KeybindingsOpts().Config.Commits.ViewBisectOptions.String(),
				description: self.c.Tr.ViewBisectOptions,
				style:       style.FgGreen,
			})
		}
	}

	// Mode-specific global keybindings
	if state := self.c.Model().WorkingTreeStateAtLastCommitRefresh; state.Any() {
		optionsMap = utils.Prepend(optionsMap, bindingInfo{
			key:         self.c.KeybindingsOpts().Config.Universal.CreateRebaseOptionsMenu.String(),
			description: state.OptionsMapTitle(self.c.Tr),
			style:       style.FgYellow,
		})
	}

	if self.c.Git().Patch.PatchBuilder.Active() {
		optionsMap = utils.Prepend(optionsMap, bindingInfo{
			key:         self.c.KeybindingsOpts().Config.Universal.CreatePatchOptionsMenu.String(),
			description: self.c.Tr.ViewPatchOptions,
			style:       style.FgYellow,
		})
	}

	self.renderOptions(self.formatBindingInfos(optionsMap))
}

func (self *OptionsMapMgr) formatBindingInfos(bindingInfos []bindingInfo) string {
	width := self.c.Views().Options.InnerWidth() - 2 // -2 for some padding
	var builder strings.Builder
	ellipsis := "…"
	separator := " | "

	length := 0

	for i, info := range bindingInfos {
		plainText := fmt.Sprintf("%s: %s", info.description, info.key)

		// Check if adding the next formatted string exceeds the available width
		textLen := utils.StringWidth(plainText)
		if i > 0 && length+len(separator)+textLen > width {
			builder.WriteString(theme.OptionsFgColor.Sprint(separator + ellipsis))
			break
		}

		formatted := info.style.Sprintf(plainText)

		if i > 0 {
			builder.WriteString(theme.OptionsFgColor.Sprint(separator))
			length += len(separator)
		}
		builder.WriteString(formatted)
		length += textLen
	}

	return builder.String()
}

func (self *OptionsMapMgr) renderOptions(options string) {
	self.c.SetViewContent(self.c.Views().Options, options)
}

type bindingInfo struct {
	key         string
	description string
	style       style.TextStyle
}
