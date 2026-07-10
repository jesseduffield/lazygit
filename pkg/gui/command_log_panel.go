package gui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// our UI command log looks like this:
// Stage File:
// git add -- 'filename'
// Unstage File:
// git reset HEAD 'filename'
//
// The 'Stage File' and 'Unstage File' lines are actions i.e they group up a set
// of command logs (typically there's only one command under an action but there may be more).
// So we call logAction to log the 'Stage File' part and then we call logCommand to log the command itself.
// We pass logCommand to our OSCommand struct so that it can handle logging commands
// for us.
func (gui *Gui) LogAction(action string) {
	if gui.Views.Extras == nil {
		return
	}

	// LogAction and LogCommand are called both from the UI thread and from git
	// worker goroutines, so bounce the writes onto the UI thread: they touch the
	// view's autoscroll flag and the GuiLog slice, which the layout/draw code
	// reads. Ordering between successive log calls is preserved by the FIFO the
	// bounce enqueues onto. It's a background bounce because writing the command
	// log is incidental display work that must not count towards lazygit being
	// busy (otherwise it could block a repo switch).
	gui.onUIThreadBackground(func() error {
		gui.Views.Extras.Autoscroll = true

		gui.GuiLog = append(gui.GuiLog, action)
		fmt.Fprint(gui.Views.Extras, "\n"+style.FgYellow.Sprint(action))
		return nil
	})
}

func (gui *Gui) LogCommand(cmdStr string, commandLine bool) {
	if gui.Views.Extras == nil {
		return
	}

	textStyle := theme.DefaultTextColor
	if !commandLine {
		// if we're not dealing with a direct command that could be run on the command line,
		// we style it differently to communicate that
		textStyle = style.FgMagenta
	}
	indentedCmdStr := "  " + strings.ReplaceAll(cmdStr, "\n", "\n  ")

	// See the comment in LogAction: bounce onto the UI thread since we may be
	// called from a git worker, in the background so it can't block a repo switch.
	gui.onUIThreadBackground(func() error {
		gui.Views.Extras.Autoscroll = true

		gui.GuiLog = append(gui.GuiLog, cmdStr)
		fmt.Fprint(gui.Views.Extras, "\n"+textStyle.Sprint(indentedCmdStr))
		return nil
	})
}

func (gui *Gui) printCommandLogHeader() {
	introStr := fmt.Sprintf(
		gui.c.Tr.CommandLogHeader,
		gui.c.UserConfig().Keybinding.Universal.ExtrasMenu,
	)
	fmt.Fprintln(gui.Views.Extras, style.FgCyan.Sprint(introStr))

	if gui.c.UserConfig().Gui.ShowRandomTip {
		fmt.Fprintf(
			gui.Views.Extras,
			"%s: %s",
			style.FgYellow.Sprint(gui.c.Tr.RandomTip),
			style.FgGreen.Sprint(gui.getRandomTip()),
		)
	}
}

func (gui *Gui) getRandomTip() string {
	tips := randomTips(gui.c.Tr, gui.c.UserConfig().Keybinding)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return tips[rnd.Intn(len(tips))]
}

func randomTips(translationSet *i18n.TranslationSet, keybindings config.KeybindingConfig) []string {
	key := func(binding config.Keybinding) string {
		return binding.String()
	}

	values := map[string]string{
		"pushKey":                key(keybindings.Universal.Push),
		"filteringMenuKey":       key(keybindings.Universal.FilteringMenu),
		"editKey":                key(keybindings.Universal.Edit),
		"rebaseOptionsKey":       key(keybindings.Universal.CreateRebaseOptionsMenu),
		"toggleTreeViewKey":      key(keybindings.Files.ToggleTreeView),
		"undoKey":                key(keybindings.Universal.Undo),
		"redoKey":                key(keybindings.Universal.Redo),
		"undoingDocsLink":        constants.Links.Docs.Undoing,
		"resetOptionsKey":        key(keybindings.Commits.ViewResetOptions),
		"pushTagKey":             key(keybindings.Branches.PushTag),
		"goIntoKey":              key(keybindings.Universal.GoInto),
		"diffingMenuKey":         key(keybindings.Universal.DiffingMenu),
		"removeKey":              key(keybindings.Universal.Remove),
		"mergeOptionsKey":        key(keybindings.Files.OpenMergeOptions),
		"revertKey":              key(keybindings.Commits.RevertCommit),
		"returnKey":              key(keybindings.Universal.Return),
		"prevPageKey":            key(keybindings.Universal.PrevPage),
		"nextPageKey":            key(keybindings.Universal.NextPage),
		"gotoTopKey":             key(keybindings.Universal.GotoTop),
		"gotoBottomKey":          key(keybindings.Universal.GotoBottom),
		"amendToCommitKey":       key(keybindings.Commits.AmendToCommit),
		"amendLastCommitKey":     key(keybindings.Files.AmendLastCommit),
		"nextBlockAlt2Key":       key(keybindings.Universal.NextBlockAlt2),
		"prevBlockAlt2Key":       key(keybindings.Universal.PrevBlockAlt2),
		"customPagersDocsLink":   constants.Links.Docs.CustomPagers,
		"customCommandsDocsLink": constants.Links.Docs.CustomCommands,
		"issuesLink":             constants.Links.Issues,
	}

	templates := []string{
		translationSet.RandomTipForcePush,
		translationSet.RandomTipFilterCommitsByPath,
		translationSet.RandomTipStartInteractiveRebase,
		translationSet.RandomTipFlatFileView,
		translationSet.RandomTipJoinTeam,
		translationSet.RandomTipUndoRedo,
		translationSet.RandomTipHardReset,
		translationSet.RandomTipPushTag,
		translationSet.RandomTipViewStashFiles,
		translationSet.RandomTipDiffCommits,
		translationSet.RandomTipDropCommit,
		translationSet.RandomTipResolveMergeConflicts,
		translationSet.RandomTipRevertCommit,
		translationSet.RandomTipExitMode,
		translationSet.RandomTipPagePanel,
		translationSet.RandomTipJumpPanel,
		translationSet.RandomTipToggleDirectory,
		translationSet.RandomTipAmendToCommit,
		translationSet.RandomTipAmendLastCommit,
		translationSet.RandomTipNavigateSidePanels,
		translationSet.RandomTipBareRepo,
		translationSet.RandomTipCommitSaveGame,
		translationSet.RandomTipSeparateRefactors,
		translationSet.RandomTipExperimentBranch,
		translationSet.RandomTipReviewDiff,
		translationSet.RandomTipReflog,
		translationSet.RandomTipStashDebugSnippets,
		translationSet.RandomTipDelta,
		translationSet.RandomTipCustomCommands,
		translationSet.RandomTipReportBug,
	}

	tips := make([]string, len(templates))
	for index, tipTemplate := range templates {
		tips[index] = utils.ResolvePlaceholderString(tipTemplate, values)
	}

	return tips
}
