package gui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
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

	gui.Views.Extras.Autoscroll = true

	gui.GuiLog = append(gui.GuiLog, action)
	fmt.Fprint(gui.Views.Extras, "\n"+style.FgYellow.Sprint(action))
}

func (gui *Gui) LogCommand(cmdStr string, commandLine bool) {
	if gui.Views.Extras == nil {
		return
	}

	gui.Views.Extras.Autoscroll = true

	textStyle := theme.DefaultTextColor
	if !commandLine {
		// if we're not dealing with a direct command that could be run on the command line,
		// we style it differently to communicate that
		textStyle = style.FgMagenta
	}
	gui.GuiLog = append(gui.GuiLog, cmdStr)
	indentedCmdStr := "  " + strings.Replace(cmdStr, "\n", "\n  ", -1)
	fmt.Fprint(gui.Views.Extras, "\n"+textStyle.Sprint(indentedCmdStr))
}

func (gui *Gui) printCommandLogHeader() {
	introStr := fmt.Sprintf(
		gui.c.Tr.CommandLogHeader,
		keybindings.Label(gui.c.UserConfig.Keybinding.Universal.ExtrasMenu),
	)
	fmt.Fprintln(gui.Views.Extras, style.FgCyan.Sprint(introStr))

	if gui.c.UserConfig.Gui.ShowRandomTip {
		fmt.Fprintf(
			gui.Views.Extras,
			"%s: %s",
			style.FgYellow.Sprint(gui.c.Tr.RandomTip),
			style.FgGreen.Sprint(gui.getRandomTip()),
		)
	}
}

func (gui *Gui) getRandomTip() string {
	config := gui.c.UserConfig.Keybinding

	formattedKey := func(key string) string {
		return keybindings.Label(key)
	}

	tips := []string{
		// keybindings and lazygit-specific advice
		fmt.Sprintf(
			"To force push, press '%s' and then if the push is rejected you will be asked if you want to force push",
			formattedKey(config.Universal.Push),
		),
		fmt.Sprintf(
			"To filter commits by path, press '%s'",
			formattedKey(config.Universal.FilteringMenu),
		),
		fmt.Sprintf(
			"To start an interactive rebase, press '%s' on a commit. You can always abort the rebase by pressing '%s' and selecting 'abort'",
			formattedKey(config.Universal.Edit),
			formattedKey(config.Universal.CreateRebaseOptionsMenu),
		),
		fmt.Sprintf(
			"In flat file view, merge conflicts are sorted to the top. To switch to flat file view press '%s'",
			formattedKey(config.Files.ToggleTreeView),
		),
		"If you want to learn Go and can think of ways to improve lazygit, join the team! Click 'Ask Question' and express your interest",
		fmt.Sprintf(
			"If you press '%s'/'%s' you can undo/redo your changes. Be wary though, this only applies to branches/commits, so only do this if your worktree is clear.\nDocs: %s",
			formattedKey(config.Universal.Undo),
			formattedKey(config.Universal.Redo),
			constants.Links.Docs.Undoing,
		),
		fmt.Sprintf(
			"to hard reset onto your current upstream branch, press '%s' in the files panel",
			formattedKey(config.Commits.ViewResetOptions),
		),
		fmt.Sprintf(
			"To push a tag, navigate to the tag in the tags tab and press '%s'",
			formattedKey(config.Branches.PushTag),
		),
		fmt.Sprintf(
			"You can view the individual files of a stash entry by pressing '%s'",
			formattedKey(config.Universal.GoInto),
		),
		fmt.Sprintf(
			"You can diff two commits by pressing '%s' on one commit and then navigating to the other. You can then press '%s' to view the files of the diff",
			formattedKey(config.Universal.DiffingMenu),
			formattedKey(config.Universal.GoInto),
		),
		fmt.Sprintf(
			"press '%s' on a commit to drop it (delete it)",
			formattedKey(config.Universal.Remove),
		),
		fmt.Sprintf(
			"If you need to pull out the big guns to resolve merge conflicts, you can press '%s' in the files panel to open 'git mergetool'",
			formattedKey(config.Files.OpenMergeTool),
		),
		fmt.Sprintf(
			"To revert a commit, press '%s' on that commit",
			formattedKey(config.Commits.RevertCommit),
		),
		fmt.Sprintf(
			"To escape a mode, for example cherry-picking, patch-building, diffing, or filtering mode, you can just spam the '%s' button. Unless of course you have `quitOnTopLevelReturn` enabled in your config",
			formattedKey(config.Universal.Return),
		),
		fmt.Sprintf(
			"You can page through the items of a panel using '%s' and '%s'",
			formattedKey(config.Universal.PrevPage),
			formattedKey(config.Universal.NextPage),
		),
		fmt.Sprintf(
			"You can jump to the top/bottom of a panel using '%s' and '%s'",
			formattedKey(config.Universal.GotoTop),
			formattedKey(config.Universal.GotoBottom),
		),
		fmt.Sprintf(
			"To collapse/expand a directory, press '%s'",
			formattedKey(config.Universal.GoInto),
		),
		fmt.Sprintf(
			"You can append your staged changes to an older commit by pressing '%s' on that commit",
			formattedKey(config.Commits.AmendToCommit),
		),
		fmt.Sprintf(
			"You can amend the last commit with your new file changes by pressing '%s' in the files panel",
			formattedKey(config.Files.AmendLastCommit),
		),
		fmt.Sprintf(
			"You can now navigate the side panels with '%s' and '%s'",
			formattedKey(config.Universal.NextBlockAlt2),
			formattedKey(config.Universal.PrevBlockAlt2),
		),

		"You can use lazygit with a bare repo by passing the --git-dir and --work-tree arguments as you would for the git CLI",

		// general advice
		"`git commit` is really just the programmer equivalent of saving your game. Always do it before embarking on an ambitious change!",
		"Try to separate commits that refactor code from commits that add new functionality: if they're squashed into the one commit, it can be hard to spot what's new.",
		"If you ever want to experiment, it's easy to create a new branch off your current one and go nuts, then delete it afterwards",
		"Always read through the diff of your changes before assigning somebody to review your code. Better for you to catch any silly mistakes than your colleagues!",
		"If something goes wrong, you can always checkout a commit from your reflog to return to an earlier state",
		"The stash is a good place to save snippets of code that you always find yourself adding when debugging.",

		// links
		fmt.Sprintf(
			"If you want a git diff with syntax colouring, check out lazygit's integration with delta:\n%s",
			constants.Links.Docs.CustomPagers,
		),
		fmt.Sprintf(
			"You can build your own custom menus and commands to run from within lazygit. For examples see:\n%s",
			constants.Links.Docs.CustomCommands,
		),
		fmt.Sprintf(
			"If you ever find a bug, do not hesitate to raise an issue on the repo:\n%s",
			constants.Links.Issues,
		),
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rnd.Intn(len(tips))
	return tips[randomIndex]
}
