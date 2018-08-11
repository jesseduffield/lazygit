package main

import "github.com/jesseduffield/gocui"

func handleCommitConfirm(g *gocui.Gui, v *gocui.View) error {
	message := trimmedContent(v)
	if message == "" {
		return createErrorPanel(g, "You cannot commit without a commit message")
	}
	if output, err := gitCommit(g, message); err != nil {
		if err == errNoUsername {
			return createErrorPanel(g, err.Error())
		}
		return createErrorPanel(g, output)
	}
	refreshFiles(g)
	g.SetViewOnBottom("commitMessage")
	switchFocus(g, v, getFilesView(g))
	return refreshCommits(g)
}

func handleCommitClose(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnBottom("commitMessage")
	return switchFocus(g, v, getFilesView(g))
}

func handleNewlineCommitMessage(g *gocui.Gui, v *gocui.View) error {
	// resising ahead of time so that the top line doesn't get hidden to make
	// room for the cursor on the second line
	x0, y0, x1, y1 := getConfirmationPanelDimensions(g, v.Buffer())
	if _, err := g.SetView("commitMessage", x0, y0, x1, y1+1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	v.EditNewLine()
	return nil
}

func handleCommitFocused(g *gocui.Gui, v *gocui.View) error {
	return renderString(g, "options", "esc: close, enter: confirm")
}
