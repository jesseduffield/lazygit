package gui

import "os/exec"

type viewUpdateOpts struct {
	title string
	task  func() error
}

type refreshMainOpts struct {
	main      *viewUpdateOpts
	secondary *viewUpdateOpts
}

// constants for updateTask's kind field
const (
	RENDER_STRING = iota
	RUN_FUNCTION
	RUN_COMMAND
)

type updateTask struct {
	kind int
	str  string
	f    func(chan struct{}) error
	cmd  *exec.Cmd
}

func (gui *Gui) createRenderStringTask(str string) {

}

func (gui *Gui) refreshMain(opts refreshMainOpts) error {
	mainView := gui.getMainView()
	secondaryView := gui.getSecondaryView()

	if opts.main != nil {
		mainView.Title = opts.main.title
		if err := opts.main.task(); err != nil {
			gui.Log.Error(err)
			return nil
		}
	}

	gui.splitMainPanel(opts.secondary != nil)

	if opts.secondary != nil {
		secondaryView.Title = opts.secondary.title
		if err := opts.secondary.task(); err != nil {
			gui.Log.Error(err)
			return nil
		}
	}

	return nil
}

func (gui *Gui) splitMainPanel(splitMainPanel bool) {
	gui.State.SplitMainPanel = splitMainPanel

	// no need to set view on bottom when splitMainPanel is false: it will have zero size anyway thanks to our view arrangement code.
	if splitMainPanel {
		_, _ = gui.g.SetViewOnTop("secondary")
	}
}

func (gui *Gui) isMainPanelSplit() bool {
	return gui.State.SplitMainPanel
}
