package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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

	case *types.RunDiffRendererTask:
		return gui.newRenderTask(view, v.Cmd, v.Prefix)
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

func (gui *Gui) stagingMainContextPair() types.MainContextPair {
	return types.NewMainContextPair(
		gui.State.Contexts.Staging,
		gui.State.Contexts.StagingSecondary,
	)
}

func (gui *Gui) patchBuildingMainContextPair() types.MainContextPair {
	return types.NewMainContextPair(
		gui.State.Contexts.CustomPatchBuilder,
		gui.State.Contexts.CustomPatchBuilderSecondary,
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
		gui.stagingMainContextPair(),
		gui.patchBuildingMainContextPair(),
		gui.mergingMainContextPair(),
	}
}

func (gui *Gui) refreshMainViews(opts types.RefreshMainOpts) {
	gui.moveMainContextPairToTop(opts.Pair)

	if opts.Main != nil {
		gui.RefreshMainView(opts.Main, opts.Pair.Main)
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

	gui.splitMainPanel(opts.Secondary != nil)
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
func (gui *Gui) clearMainView(mainContext types.Context) {
	mainContext.GetView().Clear()
	mainContext.SetHasSelectableContent(false)
	gui.State.ContextMgr.UpdateSelectionHighlights()
}

// updateDiffSelectionVisibility works out whether a main pane holds anything for a
// selection to sit on, from what it is now showing: only beneath a panel whose main
// view is a diff, and only while that diff holds something to select — never over a
// message like "No changed files", and never over a diff with nothing in it, such as a
// binary file's or an empty commit's. Whether the selection is then drawn, and drawn as
// the active one, follows from the context stack.
//
// It is asked wherever the pane's content changes: as a string is rendered, at the
// paint that reveals a command's output, with every further batch of that output, and
// once it has been read to the end. contentIsComplete tells those apart, since a render
// still being read can leave the question open (see diffPaneHasSomethingToSelect). The
// pane never answers from the render before it, and a render that leaves the question
// open is read on until it doesn't, so the answer is always about what is there.
func (gui *Gui) updateDiffSelectionVisibility(view *gocui.View, contentIsComplete bool) {
	mainContext := gui.mainContextForView(view)
	if mainContext == nil {
		return
	}

	gui.dropAnAnswerAboutAnotherRender(mainContext, view)

	if hasSomethingToSelect, known := gui.diffPaneHasSomethingToSelect(
		mainContext, view, contentIsComplete,
	); known {
		mainContext.SetHasSelectableContent(hasSomethingToSelect)
		gui.State.ContextMgr.UpdateSelectionHighlights()
	} else {
		gui.readOnUntilTheDiffPaneCanTell(view)
	}
}

// dropAnAnswerAboutAnotherRender takes away what the pane worked out about the content
// of an earlier render, so that this one starts from no answer rather than inheriting
// one. An answer about other content says nothing about this content: carried over, it
// shows a selection over a diff that may have nothing to select, or hides one over a
// diff that has.
//
// A re-render of the same content keeps its answer, and with it the selection drawn
// over it, since that answer is still about what the pane is showing.
func (gui *Gui) dropAnAnswerAboutAnotherRender(mainContext *context.MainContext, view *gocui.View) {
	manager := gui.getViewBufferManagerForView(view)
	if manager == nil {
		return
	}

	if key := manager.GetTaskKey(); key != mainContext.SelectableContentRenderKey() {
		mainContext.SetSelectableContentRenderKey(key)
		mainContext.SetHasSelectableContent(false)
		gui.State.ContextMgr.UpdateSelectionHighlights()
	}
}

// readOnUntilTheDiffPaneCanTell keeps a render going past the lines that were asked of
// it, while the pane still can't say whether there is anything in it to select.
//
// A render is asked for as many lines as the scrollbar needs (see
// linesToReadFromCmdTask), and what a commit's diff opens with can run past that: the
// diffstat of a commit touching thousands of files, or a commit message thousands of
// lines long. Without this the pane would be left with no answer until the user
// scrolled far enough to ask for the rest themselves, which is no way to find out
// whether a diff can be acted on. Another render's worth is asked for each time, so the
// reading stops soon after the first change line, and only runs to the end of a diff
// that has none.
func (gui *Gui) readOnUntilTheDiffPaneCanTell(view *gocui.View) {
	manager := gui.getViewBufferManagerForView(view)
	if manager == nil {
		return
	}

	step := gui.linesToReadFromCmdTask(view).Total
	if step < 0 {
		// A view that is being searched is already being read to the end.
		return
	}
	manager.ReadLinesAndWait(view.LinesHeight() + step)
}

// diffPaneHasSomethingToSelect answers whether the given main pane holds anything for a
// selection to sit on, from what it is showing so far. known is false while a render
// still being read leaves the question open.
func (gui *Gui) diffPaneHasSomethingToSelect(
	mainContext *context.MainContext, view *gocui.View, contentIsComplete bool,
) (bool, bool) {
	if _, showsDiff := gui.State.ContextMgr.CurrentSide().(types.DiffMainViewContext); !showsDiff {
		// Under a panel that shows no diff there is nothing to select whatever the pane
		// ends up holding, so this needs no content to answer. Answering it now matters,
		// because a render may never reach an end. The rest of a long commit log is read
		// only as far as the user scrolls, and until then the pane would go on showing
		// the selection it was left with under the panel before.
		return false, true
	}

	if !contentIsComplete && mainContext.HasSelectableContent() {
		// This render has already found something to select, and its content only grows
		// from here, so there is nothing to ask again — nor to read the diff for. An
		// answer the render before it gave has been dropped by now (see
		// dropAnAnswerAboutAnotherRender), so this really is about the content in hand.
		return true, true
	}

	hasChangeLines := gui.helpers.DiffLine.ViewHasChangeLines(view)

	// One change line among those read settles it. Finding none in a render that is
	// still going may only mean the changes are in the part still to come. A commit's
	// diff opens with a diffstat, and for a commit touching hundreds of files that runs
	// well past the screenful the first paint reveals, so that answer waits.
	return hasChangeLines, contentIsComplete || hasChangeLines
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

func (gui *Gui) splitMainPanel(splitMainPanel bool) {
	gui.State.SplitMainPanel = splitMainPanel
}

// reApplySearch runs a search the view holds again over the content a render has just
// finished putting there, so that the matches highlighted and the "x of y" status
// describe what the view shows now rather than what it showed when the search was
// typed. Call it once the content is final.
func (gui *Gui) reApplySearch(view *gocui.View) {
	// While the prompt is open, the search view holds what the user is typing, and the
	// status would be written over it.
	if gui.State.ContextMgr.Current().GetKey() == context.SEARCH_CONTEXT_KEY {
		return
	}

	view.RefreshSearch()
}
