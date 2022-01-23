package gui

import (
	"testing"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestCanDeactivatePopupContextsWithoutViews(t *testing.T) {
	contexts := []func(gui *Gui) types.Context{
		func(gui *Gui) types.Context { return gui.State.Contexts.Credentials },
		func(gui *Gui) types.Context { return gui.State.Contexts.Confirmation },
		func(gui *Gui) types.Context { return gui.State.Contexts.CommitMessage },
		func(gui *Gui) types.Context { return gui.State.Contexts.Search },
	}

	for _, c := range contexts {
		gui := NewDummyGui()
		context := c(gui)
		gui.g = &gocui.Gui{}

		_ = gui.deactivateContext(context)

		// This really only checks a prerequisit, not the effect of deactivateContext
		view, _ := gui.g.View(context.GetViewName())
		assert.Nil(t, view, string(context.GetKey()))
	}
}

func TestCanDeactivateCommitFilesContextsWithoutViews(t *testing.T) {
	gui := NewDummyGui()
	gui.g = &gocui.Gui{}

	_ = gui.deactivateContext(gui.State.Contexts.CommitFiles)

	// This really only checks a prerequisite, not the effect of deactivateContext
	view, _ := gui.g.View(gui.State.Contexts.CommitFiles.GetViewName())
	assert.Nil(t, view)
}
