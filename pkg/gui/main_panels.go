package gui

import (
	"github.com/jesseduffield/gocui"
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

	case *types.RunPtyTask:
		return gui.newPtyTask(view, v.Cmd, v.Prefix)
	}

	return nil
}

func (gui *Gui) moveMainContextPairToTop(pair types.MainContextPair) {
	gui.setWindowContext(pair.Main)
	gui.moveToTopOfWindow(pair.Main)
	if pair.Secondary != nil {
		gui.setWindowContext(pair.Secondary)
		gui.moveToTopOfWindow(pair.Secondary)
	}
}

func (gui *Gui) RefreshMainView(opts *types.ViewUpdateOpts, context types.Context) error {
	view := context.GetView()

	if opts.Title != "" {
		view.Title = opts.Title
	}

	if err := gui.runTaskForView(view, opts.Task); err != nil {
		gui.c.Log.Error(err)
		return nil
	}

	return nil
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

func (gui *Gui) refreshMainViews(opts types.RefreshMainOpts) error {
	// need to reset scroll positions of all other main views
	for _, pair := range gui.allMainContextPairs() {
		if pair.Main != opts.Pair.Main {
			_ = pair.Main.GetView().SetOrigin(0, 0)
		}
		if pair.Secondary != nil && pair.Secondary != opts.Pair.Secondary {
			_ = pair.Secondary.GetView().SetOrigin(0, 0)
		}
	}

	if opts.Main != nil {
		if err := gui.RefreshMainView(opts.Main, opts.Pair.Main); err != nil {
			return err
		}
	}

	if opts.Secondary != nil {
		if err := gui.RefreshMainView(opts.Secondary, opts.Pair.Secondary); err != nil {
			return err
		}
	} else if opts.Pair.Secondary != nil {
		opts.Pair.Secondary.GetView().Clear()
	}

	gui.moveMainContextPairToTop(opts.Pair)

	gui.splitMainPanel(opts.Secondary != nil)

	return nil
}

func (gui *Gui) splitMainPanel(splitMainPanel bool) {
	gui.State.SplitMainPanel = splitMainPanel
}

func (gui *Gui) isMainPanelSplit() bool {
	return gui.State.SplitMainPanel
}
