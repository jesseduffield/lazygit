package custom_commands

import (
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

type promptLoadingGui struct {
	types.IGuiCommon

	t *testing.T

	loadingActive   bool
	loadingMessages []string
	git             *commands.GitCommand
	os              *oscommands.OSCommand
	menuOpts        *types.CreateMenuOptions
}

func (self *promptLoadingGui) WithWaitingStatusSync(message string, f func() error) error {
	assert.False(self.t, self.loadingActive)

	self.loadingMessages = append(self.loadingMessages, message)
	self.loadingActive = true
	defer func() { self.loadingActive = false }()

	return f()
}

func (self *promptLoadingGui) Git() *commands.GitCommand {
	return self.git
}

func (self *promptLoadingGui) OS() *oscommands.OSCommand {
	return self.os
}

func (self *promptLoadingGui) Menu(opts types.CreateMenuOptions) error {
	assert.False(self.t, self.loadingActive)

	self.menuOpts = &opts
	return nil
}

func TestResolvePromptRunsTemplateResolutionWithPromptLoadingText(t *testing.T) {
	t.Parallel()

	gui := &promptLoadingGui{t: t}
	handler := newHandlerCreatorForPromptLoadingTest(gui)

	prompt := &config.CustomCommandPrompt{
		Type:         "input",
		Title:        "Pick value",
		LoadingText:  "Generating value",
		InitialValue: "{{ runCommand \"generate-value\" }}",
	}

	resolvedPrompt, err := handler.resolvePrompt(prompt, func(template string) (string, error) {
		assert.True(t, gui.loadingActive)
		return "resolved: " + template, nil
	})

	if !assert.NoError(t, err) || !assert.NotNil(t, resolvedPrompt) {
		return
	}
	assert.Equal(t, []string{"Generating value"}, gui.loadingMessages)
	assert.Equal(t, "Generating value", resolvedPrompt.LoadingText)
	assert.Equal(t, "resolved: {{ runCommand \"generate-value\" }}", resolvedPrompt.InitialValue)
}

func TestResolvePromptSkipsPromptLoadingWhenLoadingTextIsEmpty(t *testing.T) {
	t.Parallel()

	gui := &promptLoadingGui{t: t}
	handler := newHandlerCreatorForPromptLoadingTest(gui)

	prompt := &config.CustomCommandPrompt{
		Type:         "input",
		Title:        "Pick value",
		InitialValue: "static value",
	}

	resolvedPrompt, err := handler.resolvePrompt(prompt, func(template string) (string, error) {
		assert.False(t, gui.loadingActive)
		return template, nil
	})

	if !assert.NoError(t, err) || !assert.NotNil(t, resolvedPrompt) {
		return
	}
	assert.Empty(t, gui.loadingMessages)
	assert.Equal(t, "static value", resolvedPrompt.InitialValue)
}

func TestInputPromptRunsSuggestionsCommandWithPromptLoadingText(t *testing.T) {
	t.Parallel()

	gui := &promptLoadingGui{t: t}
	runner := oscommands.NewFakeRunner(t).ExpectFunc(
		"runs the suggestions command while the prompt loading status is active",
		func(cmdObj *oscommands.CmdObj) bool {
			return gui.loadingActive && strings.Contains(cmdObj.ToString(), "list-suggestions")
		},
		"one\ntwo",
		nil,
	)
	gui.os = oscommands.NewDummyOSCommandWithRunner(runner)

	handler := newHandlerCreatorForPromptLoadingTest(gui)
	findSuggestionsFn, err := handler.generateFindSuggestionsFunc(&config.CustomCommandPrompt{
		Type:        "input",
		Title:       "Pick value",
		LoadingText: "Loading suggestions",
		Suggestions: config.CustomCommandSuggestions{
			Command: "list-suggestions",
		},
	})

	if !assert.NoError(t, err) || !assert.NotNil(t, findSuggestionsFn) {
		return
	}
	runner.CheckForMissingCalls()
	assert.Equal(t, []string{"Loading suggestions"}, gui.loadingMessages)
	assert.Equal(t, []*types.Suggestion{
		{Value: "one", Label: "one"},
		{Value: "two", Label: "two"},
	}, findSuggestionsFn(""))
}

func TestMenuPromptFromCommandRunsCommandWithPromptLoadingText(t *testing.T) {
	t.Parallel()

	cmn := common.NewDummyCommon()
	gui := &promptLoadingGui{t: t}
	runner := oscommands.NewFakeRunner(t).ExpectFunc(
		"runs the menu command while the prompt loading status is active",
		func(cmdObj *oscommands.CmdObj) bool {
			return gui.loadingActive && strings.Join(cmdObj.GetCmd().Args, " ") == "list-options"
		},
		"first\nsecond",
		nil,
	)

	gui.git = &commands.GitCommand{
		Custom: git_commands.NewCustomCommands(git_commands.NewGitCommon(
			cmn,
			&git_commands.GitVersion{Major: 2},
			oscommands.NewDummyCmdObjBuilder(runner),
			oscommands.NewDummyOSCommandWithRunner(runner),
			git_commands.MockRepoPaths("."),
			git_commands.NewConfigCommands(cmn, git_config.NewFakeGitConfig(nil)),
			config.NewPagerConfig(cmn.UserConfig),
		)),
	}

	handler := newHandlerCreatorForPromptLoadingTest(gui)
	err := handler.menuPromptFromCommand(&config.CustomCommandPrompt{
		Type:        "menuFromCommand",
		Title:       "Choose an option",
		LoadingText: "Loading options",
		Command:     "list-options",
		Filter:      "(?P<option>.*)",
		ValueFormat: "{{ .option }}",
		LabelFormat: "{{ .option }}",
	}, func(string) error { return nil })

	assert.NoError(t, err)
	runner.CheckForMissingCalls()
	assert.Equal(t, []string{"Loading options"}, gui.loadingMessages)
	if !assert.NotNil(t, gui.menuOpts) {
		return
	}
	assert.Equal(t, "Choose an option", gui.menuOpts.Title)
	if !assert.Len(t, gui.menuOpts.Items, 2) {
		return
	}
	assert.Equal(t, []string{"first"}, gui.menuOpts.Items[0].LabelColumns)
	assert.Equal(t, []string{"second"}, gui.menuOpts.Items[1].LabelColumns)
}

func newHandlerCreatorForPromptLoadingTest(gui *promptLoadingGui) *HandlerCreator {
	cmn := common.NewDummyCommon()

	return &HandlerCreator{
		c: &helpers.HelperCommon{
			Common:     cmn,
			IGuiCommon: gui,
		},
		resolver:      NewResolver(cmn),
		menuGenerator: NewMenuGenerator(cmn),
	}
}
