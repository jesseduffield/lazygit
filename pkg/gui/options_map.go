package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type OptionsMapMgr struct {
	c *helpers.HelperCommon
}

func (gui *Gui) renderContextOptionsMap() {
	// In demos, we render our own content to this view
	if gui.integrationTest != nil && gui.integrationTest.IsDemo() {
		return
	}
	mgr := OptionsMapMgr{c: gui.c}
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

	allBindings := append(currentContextBindings, globalBindings...)

	bindingsToDisplay := lo.Filter(allBindings, func(binding *types.Binding, _ int) bool {
		return binding.DisplayOnScreen && !binding.IsDisabled()
	})

	optionsMap := lo.Map(bindingsToDisplay, func(binding *types.Binding, _ int) bindingInfo {
		displayStyle := theme.OptionsFgColor
		if binding.DisplayStyle != nil {
			displayStyle = *binding.DisplayStyle
		}

		description := binding.Description
		if binding.ShortDescription != "" {
			description = binding.ShortDescription
		}

		return bindingInfo{
			key:         keybindings.LabelFromKey(binding.Key),
			description: description,
			style:       displayStyle,
		}
	})

	// Mode-specific local keybindings
	if currentContext.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
		if self.c.Modes().CherryPicking.Active() {
			optionsMap = utils.Prepend(optionsMap, bindingInfo{
				key:         keybindings.Label(self.c.KeybindingsOpts().Config.Commits.PasteCommits),
				description: self.c.Tr.PasteCommits,
				style:       style.FgCyan,
			})
		}

		if self.c.Model().BisectInfo.Started() {
			optionsMap = utils.Prepend(optionsMap, bindingInfo{
				key:         keybindings.Label(self.c.KeybindingsOpts().Config.Commits.ViewBisectOptions),
				description: self.c.Tr.ViewBisectOptions,
				style:       style.FgGreen,
			})
		}
	}

	// Mode-specific global keybindings
	if self.c.Model().WorkingTreeStateAtLastCommitRefresh.IsRebasing() {
		optionsMap = utils.Prepend(optionsMap, bindingInfo{
			key:         keybindings.Label(self.c.KeybindingsOpts().Config.Universal.CreateRebaseOptionsMenu),
			description: self.c.Tr.ViewRebaseOptions,
			style:       style.FgYellow,
		})
	} else if self.c.Model().WorkingTreeStateAtLastCommitRefresh.IsMerging() {
		optionsMap = utils.Prepend(optionsMap, bindingInfo{
			key:         keybindings.Label(self.c.KeybindingsOpts().Config.Universal.CreateRebaseOptionsMenu),
			description: self.c.Tr.ViewMergeOptions,
			style:       style.FgYellow,
		})
	}

	if self.c.Git().Patch.PatchBuilder.Active() {
		optionsMap = utils.Prepend(optionsMap, bindingInfo{
			key:         keybindings.Label(self.c.KeybindingsOpts().Config.Universal.CreatePatchOptionsMenu),
			description: self.c.Tr.ViewPatchOptions,
			style:       style.FgYellow,
		})
	}

	self.renderOptions(self.formatBindingInfos(optionsMap))
}

func (self *OptionsMapMgr) formatBindingInfos(bindingInfos []bindingInfo) string {
	width := self.c.Views().Options.Width() - 4 // -4 for the padding
	var builder strings.Builder
	ellipsis := "…"
	separator := " | "

	length := 0

	for i, info := range bindingInfos {
		plainText := fmt.Sprintf("%s: %s", info.description, info.key)

		// Check if adding the next formatted string exceeds the available width
		if i > 0 && length+len(separator)+len(plainText) > width {
			builder.WriteString(theme.OptionsFgColor.Sprint(separator + ellipsis))
			break
		}

		formatted := info.style.Sprintf(plainText)

		if i > 0 {
			builder.WriteString(theme.OptionsFgColor.Sprint(separator))
			length += len(separator)
		}
		builder.WriteString(formatted)
		length += len(plainText)
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
