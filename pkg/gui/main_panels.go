package gui

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
)

func (gui *Gui) runTaskForView(view *gocui.View, task types.UpdateTask) error {
	switch v := task.(type) {
	case *types.RenderStringTask:
		return gui.newStringTask(view, v.Str)

	case *types.RenderStringWithoutScrollTask:
		return gui.newStringTaskWithoutScroll(view, v.Str)

	case *types.RenderStringWithScrollTask:
		return gui.newStringTaskWithScroll(view, v.Str, v.OriginX, v.OriginY)

	case *types.RunCommandTask:
		return gui.newCmdTask(view, v.Cmd, v.Prefix)

	case *types.RunPtyTask:
		return gui.newPtyTask(view, v.Cmd, v.Prefix)
	}

	return nil
}

func (gui *Gui) moveMainContextPairToTop(pair types.MainContextPair) {
	gui.moveMainContextToTop(pair.Main)
	if pair.Secondary != nil {
		gui.moveMainContextToTop(pair.Secondary)
	}
}

func (gui *Gui) moveMainContextToTop(context types.Context) {
	gui.helpers.Window.SetWindowContext(context)

	view := context.GetView()

	topView := gui.helpers.Window.TopViewInWindow(context.GetWindowName(), true)

	if topView != nil && topView != view {
		// We need to copy the content to avoid a flicker effect: If we're flicking
		// through files in the files panel, we use a different view to render the
		// files vs the directories, and if you select dir A, then file B, then dir
		// C, you'll briefly see dir A's contents again before the view is updated.
		// So here we're copying the content from the top window to avoid that
		// flicker effect.
		gui.g.CopyContent(topView, view)

		if err := gui.g.SetViewOnTopOf(view.Name(), topView.Name()); err != nil {
			gui.Log.Error(err)
		}
	}
}

func (gui *Gui) RefreshMainView(opts *types.ViewUpdateOpts, context types.Context) {
	view := context.GetView()

	if opts.Title != "" {
		view.Title = opts.Title
	}

	view.Subtitle = opts.SubTitle

	if err := gui.runTaskForView(view, opts.Task); err != nil {
		gui.c.Log.Error(err)
	}
}

func (gui *Gui) normalMainContextPair() types.MainContextPair {
	return types.NewMainContextPair(
		gui.State.Contexts.Normal,
		gui.State.Contexts.NormalSecondary,
	)
}

func (gui *Gui) mergingMainContextPair() types.MainContextPair {
	return types.NewMainContextPair(
		gui.State.Contexts.MergeConflicts,
		nil,
	)
}

func (gui *Gui) allMainContextPairs() []types.MainContextPair {
	return []types.MainContextPair{
		gui.normalMainContextPair(),
		gui.mergingMainContextPair(),
	}
}

func (gui *Gui) refreshMainViews(opts types.RefreshMainOpts) {
	panes := mainPanesFor(opts)

	// Before the render is triggered, so that the pane the focus moves into can be
	// told where to put its selection as it renders.
	gui.followFocusIntoWorkablePane(opts)
	gui.keepDiffSelectionAcrossACommitRewrite(opts)

	gui.moveMainContextPairToTop(opts.Pair)

	gui.handOverMainSection(opts.Pair, panes)

	if opts.Main != nil {
		gui.RefreshMainView(opts.Main, opts.Pair.Main)
	} else {
		gui.clearMainView(opts.Pair.Main)
	}

	if opts.Secondary != nil {
		gui.RefreshMainView(opts.Secondary, opts.Pair.Secondary)
	} else if opts.Pair.Secondary != nil {
		gui.clearMainView(opts.Pair.Secondary)
	}

	// Reset the scroll positions of all the other main views. We do this after
	// moving this pair to the top (which copies the previously-shown view's
	// content into the now-visible one to avoid a blank frame): resetting first
	// would zero that source view's scroll before it gets copied, forcing the
	// placeholder to the top instead of leaving it where the screen already was.
	for _, pair := range gui.allMainContextPairs() {
		if pair.Main != opts.Pair.Main {
			pair.Main.GetView().SetOrigin(0, 0)
		}
		if pair.Secondary != nil && pair.Secondary != opts.Pair.Secondary {
			pair.Secondary.GetView().SetOrigin(0, 0)
		}
	}

	gui.setMainPanes(panes)
}

// handOverMainSection carries the content of the main section from the pane that has
// been showing it on its own to the pane about to, when a render moves the section's
// content from one to the other — a file's changes going from unstaged to staged, say.
//
// The section is one region of the screen to the user, so a change of which pane holds
// it has to look like that region re-rendering rather than blanking and filling in
// again: the incoming pane shows what the outgoing one was showing, where it was
// showing it, until its own render has read enough to be swapped in. It renders from
// the top when it does, the content it took over not being its own (see
// clearMainView).
func (gui *Gui) handOverMainSection(pair types.MainContextPair, panes types.MainPanes) {
	// The lower pane is always the same view, being the only one a render can leave
	// holding the section on its own; the upper one is whichever view of the main
	// window this render is for, which moveMainContextPairToTop has just given a copy
	// of what that window was showing.
	upper, lower := pair.Main.GetView(), gui.Views.Secondary

	var from, to *gocui.View
	switch {
	case gui.State.MainPanes == types.MainPaneOnly && panes == types.SecondaryPaneOnly:
		from, to = upper, lower
	case gui.State.MainPanes == types.SecondaryPaneOnly && panes == types.MainPaneOnly:
		from, to = lower, upper
	default:
		return
	}

	gui.g.CopyContent(from, to)
}

// mainPanesFor says which panes the given render occupies: the one it has content for,
// or both when it has content for both.
func mainPanesFor(opts types.RefreshMainOpts) types.MainPanes {
	switch {
	case opts.Secondary == nil:
		return types.MainPaneOnly
	case opts.Main == nil:
		return types.SecondaryPaneOnly
	default:
		return types.BothMainPanes
	}
}

// followFocusIntoWorkablePane moves the focus out of a main pane that the render about
// to happen leaves nothing to work on, and into the one it does.
//
// Each side of a file's diff has a pane of its own, and a pane holds something only
// while its side of the file does. So anything that empties the side the focus is on
// leaves that pane with nothing: staging the last unstaged change, committing what was
// staged, or either of those happening outside lazygit and arriving with a refresh.
// Usually the pane goes away with its content; configured to always split the diff it
// stays, empty. Either way the focus has nothing left to act on where it is.
//
// The pane moved into gets its selection once the render has finished and there is
// something to put one on, and shows none until then, so that the selection it was
// left with the last time it was used doesn't appear for a frame. A pane that has
// already been told where to put its selection — by the action that caused all this —
// keeps what it was told.
func (gui *Gui) followFocusIntoWorkablePane(opts types.RefreshMainOpts) {
	// The focused main view's two panes only: the staging and patch-building views
	// arrange theirs for themselves, and the merge-conflicts view has just the one.
	if opts.Pair.Main.GetKey() != context.NORMAL_MAIN_CONTEXT_KEY {
		return
	}

	current := gui.State.ContextMgr.CurrentStatic().GetKey()
	if current != opts.Pair.Main.GetKey() && current != opts.Pair.Secondary.GetKey() {
		return
	}
	pane := onlyWorkablePane(opts)
	if pane == nil || pane.GetKey() == current {
		return
	}

	target := gui.mainContextForView(pane.GetView())
	target.SetHasSelectableContent(false)
	gui.State.ContextMgr.UpdateSelectionHighlights()
	if manager := gui.getManager(target.GetView()); !manager.HasRestoreForNextTask() {
		manager.SetRestoreForNextTask(&tasks.RenderRestore{
			// The whole render is read before it is shown: where the selection goes
			// is decided from what is there, and a change line further down would
			// otherwise be missed.
			FirstPaintReady: func() bool { return false },
			Apply: func(swapIn func()) {
				swapIn()
				gui.helpers.DiffLine.EstablishSelection(target, -1)
			},
		})
	}
	gui.State.ContextMgr.Push(target, types.OnFocusOpts{})
}

// keepDiffSelectionAcrossACommitRewrite arranges for a selection in the focused main
// view to come back on the same change of the diff when the render about to happen is
// of a different diff from the one on screen. That is what a commit rewritten under the
// user looks like — moving a patch out of it, discarding lines from it, undoing either —
// and it leaves the selection a position in a rendering that no longer exists.
//
// What is remembered is which change of the diff the selection was on rather than which
// line of which file, a rewrite being precisely a change to those lines: the change that
// takes its place is where the work carries on.
//
// It is asked of every render and stands down unless there is a selection to keep,
// nothing more precise is already waiting to be put back — the position preserves and
// the post-action reveals know better where their selection belongs — and the diff
// really is another one. A plain refresh re-renders the same diff, where the selection,
// possibly a range the user is in the middle of making, is still exactly right.
func (gui *Gui) keepDiffSelectionAcrossACommitRewrite(opts types.RefreshMainOpts) {
	// The focused main view's two panes only: no other pair has a diff selection.
	if opts.Pair.Main.GetKey() != context.NORMAL_MAIN_CONTEXT_KEY {
		return
	}

	current := gui.State.ContextMgr.CurrentStatic().GetKey()
	for _, pane := range []struct {
		context types.Context
		update  *types.ViewUpdateOpts
	}{
		{opts.Pair.Main, opts.Main},
		{opts.Pair.Secondary, opts.Secondary},
	} {
		if pane.update == nil || pane.context.GetKey() != current {
			continue
		}
		mainContext := gui.mainContextForView(pane.context.GetView())
		if mainContext == nil || !mainContext.GetView().Highlight {
			continue
		}
		manager := gui.getViewBufferManagerForView(mainContext.GetView())
		if manager == nil || manager.HasRestoreForNextTask() {
			continue
		}
		key, ok := diffTaskCommandKey(pane.update.Task)
		if !ok || key == manager.GetTaskKey() {
			continue
		}

		first, _ := mainContext.GetView().SelectedLineRange()
		gui.helpers.DiffLine.RevealSelectionAfterAction(mainContext, mainContext, first, 0, nil)
	}
}

// diffTaskCommandKey returns the key the given render will be remembered under, which
// is what says whether it is a render of the same diff as the one on screen. ok is
// false for a render that is a message rather than a diff.
func diffTaskCommandKey(task types.UpdateTask) (string, bool) {
	switch task := task.(type) {
	case *types.RunCommandTask:
		return strings.Join(task.Cmd.Args, " "), true
	case *types.RunPtyTask:
		return strings.Join(task.Cmd.Args, " "), true
	}
	return "", false
}

// onlyWorkablePane returns the main pane a render leaves as the only one worth having
// the focus in, or nil when that is true of both of them or of neither. Being shown is
// not the same as being worth working in: a pane the layout keeps around for the sake
// of always splitting the diff shows an empty side of the file.
func onlyWorkablePane(opts types.RefreshMainOpts) types.Context {
	main := opts.Main != nil && !opts.Main.NothingToActOn
	secondary := opts.Secondary != nil && !opts.Secondary.NothingToActOn
	if main == secondary {
		return nil
	}
	if main {
		return opts.Pair.Main
	}
	return opts.Pair.Secondary
}

// clampDiffSelectionToContent brings the focused main view's selection back onto the
// content when the render that just finished left the diff with fewer lines than the
// selection was on — a diff renderer that renders the same diff more compactly, a
// smaller context size. That selection lives in the view rather than in a model, so
// nothing else re-derives it, and past the end of the content it isn't drawn at all,
// which reads as having no selection until an arrow key brings it back.
//
// Called at end of input, when the content is final: doing it while the render is
// still loading would drag the selection to a line that only looks like the last one.
// Only these two views need it; every other view's selection is derived from a model
// as it renders, and so is clamped along with it.
func (gui *Gui) clampDiffSelectionToContent(view *gocui.View) {
	if gui.mainContextForView(view) == nil {
		return
	}
	if !view.Highlight {
		return
	}

	if lastLine := view.ViewLinesHeight() - 1; view.SelectedLineIdx() > lastLine {
		view.FocusPoint(0, max(0, lastLine), true)
	}
}

// clearMainView empties a pane that is being given nothing to show, selection and all.
//
// An emptied pane is showing nothing, so it also goes back to the top and stops
// claiming the render it was showing: whatever it is given next is content the user
// hasn't seen there, and is shown from the top like any other.
//
// A position waiting to be put back goes too: this pane is getting no render for it
// to ride, and whoever is waiting for the view to be back where it belongs has to
// hear that it never will be.
func (gui *Gui) clearMainView(mainContext types.Context) {
	view := mainContext.GetView()
	view.Clear()
	view.SetOrigin(0, 0)
	mainContext.SetHasSelectableContent(false)
	gui.State.ContextMgr.UpdateSelectionHighlights()
	if manager := gui.getViewBufferManagerForView(view); manager != nil {
		manager.ForgetRenderedContent()
		manager.DropRestoreForNextTask()
	}
}

// updateDiffPaneDecorations re-derives what is drawn over a main pane's content, rather
// than being part of it: whether a selection is shown, and which lines are marked as
// being in the custom patch.
//
// A pane holds something for a selection to sit on only beneath a panel whose main
// view is a diff, and only while that diff holds something to select — never over a
// message like "No changed files", and never over a diff with nothing in it, such as a
// binary file's or an empty commit's. Whether the selection is then drawn, and drawn as
// the active one, follows from the context stack.
//
// It is asked once a pane's content is final, which is when the last of those questions
// can be answered: as a string is rendered, and at the end of a command's output. At the
// moment a render is triggered there is nothing to read yet — and the content still on
// screen until the swap is the one the selection belongs to, so leaving it alone until
// then is right.
func (gui *Gui) updateDiffPaneDecorations(view *gocui.View) {
	mainContext := gui.mainContextForView(view)
	if mainContext == nil {
		return
	}

	_, showsDiff := gui.State.ContextMgr.CurrentSide().(types.DiffMainViewContext)
	mainContext.SetHasSelectableContent(showsDiff && gui.helpers.DiffLine.ViewHasChangeLines(view))
	gui.State.ContextMgr.UpdateSelectionHighlights()

	// The marks are over the diff, which is the upper pane's; the lower one shows the
	// patch they are marks of.
	if view == gui.Views.Main {
		gui.helpers.DiffLine.RefreshInclusionGutter()
	}
}

// mainContextForView returns the context of the main pane the given view is, or nil for
// any other view.
func (gui *Gui) mainContextForView(view *gocui.View) *context.MainContext {
	switch view {
	case gui.Views.Main:
		return gui.State.Contexts.Normal
	case gui.Views.Secondary:
		return gui.State.Contexts.NormalSecondary
	}
	return nil
}

func (gui *Gui) setMainPanes(panes types.MainPanes) {
	gui.State.MainPanes = panes

	// The label for the key that focuses the main view belongs on the pane that key
	// focuses, which is the secondary one while it is the only one shown.
	if panes == types.SecondaryPaneOnly {
		gui.showFocusMainViewJumpLabelOn(gui.Views.Secondary)
	} else {
		gui.showFocusMainViewJumpLabelOn(gui.Views.Main)
	}
}

// showFocusMainViewJumpLabelOn puts the main view's jump label on the given pane and
// takes it off the other one, so that only the pane the key focuses wears it.
func (gui *Gui) showFocusMainViewJumpLabelOn(view *gocui.View) {
	gui.Views.Main.TitlePrefix = ""
	gui.Views.Secondary.TitlePrefix = ""
	view.TitlePrefix = gui.focusMainViewJumpLabel
}
