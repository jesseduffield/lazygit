package main

import (
	"errors"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

var (
	// ErrNoCommits : When no commits are found for the branch
	ErrNoCommits = errors.New("No commits for this branch")
)

func refreshCommits(g *gocui.Gui) error {
	g.Update(func(*gocui.Gui) error {
		state.Commits = getCommits()
		v, err := g.View("commits")
		if err != nil {
			panic(err)
		}
		v.Clear()
		red := color.New(color.FgRed)
		yellow := color.New(color.FgYellow)
		white := color.New(color.FgWhite)
		shaColor := white
		for _, commit := range state.Commits {
			if commit.Pushed {
				shaColor = red
			} else {
				shaColor = yellow
			}
			shaColor.Fprint(v, commit.Sha+" ")
			white.Fprintln(v, commit.Name)
		}
		refreshStatus(g)
		return nil
	})
	return nil
}

func handleResetToCommit(g *gocui.Gui, commitView *gocui.View) error {
	return createConfirmationPanel(g, commitView, "Reset To Commit", "Are you sure you want to reset to this commit?", func(g *gocui.Gui, v *gocui.View) error {
		commit, err := getSelectedCommit(g)
		devLog(commit)
		if err != nil {
			panic(err)
		}
		if output, err := gitResetToCommit(commit.Sha); err != nil {
			return createErrorPanel(g, output)
		}
		if err := refreshCommits(g); err != nil {
			panic(err)
		}
		if err := refreshFiles(g); err != nil {
			panic(err)
		}
		resetOrigin(commitView)
		return handleCommitSelect(g, nil)
	}, nil)
}

func renderCommitsOptions(g *gocui.Gui) error {
	return renderOptionsMap(g, map[string]string{
		"s":       "squash down",
		"r":       "rename",
		"g":       "reset to this commit",
		"← → ↑ ↓": "navigate",
	})
}

func handleCommitSelect(g *gocui.Gui, v *gocui.View) error {
	if err := renderCommitsOptions(g); err != nil {
		return err
	}
	commit, err := getSelectedCommit(g)
	if err != nil {
		if err != ErrNoCommits {
			return err
		}
		return renderString(g, "main", "No commits for this branch")
	}
	commitText := gitShow(commit.Sha)
	return renderString(g, "main", commitText)
}

func handleCommitSquashDown(g *gocui.Gui, v *gocui.View) error {
	if getItemPosition(v) != 0 {
		return createErrorPanel(g, "Can only squash topmost commit")
	}
	if len(state.Commits) == 1 {
		return createErrorPanel(g, "You have no commits to squash with")
	}
	commit, err := getSelectedCommit(g)
	if err != nil {
		return err
	}
	if output, err := gitSquashPreviousTwoCommits(commit.Name); err != nil {
		return createErrorPanel(g, output)
	}
	if err := refreshCommits(g); err != nil {
		panic(err)
	}
	refreshStatus(g)
	return handleCommitSelect(g, v)
}

func handleRenameCommit(g *gocui.Gui, v *gocui.View) error {
	if getItemPosition(v) != 0 {
		return createErrorPanel(g, "Can only rename topmost commit")
	}
	createPromptPanel(g, v, "Rename Commit", func(g *gocui.Gui, v *gocui.View) error {
		if output, err := gitRenameCommit(v.Buffer()); err != nil {
			return createErrorPanel(g, output)
		}
		if err := refreshCommits(g); err != nil {
			panic(err)
		}
		return handleCommitSelect(g, v)
	})
	return nil
}

func getSelectedCommit(g *gocui.Gui) (Commit, error) {
	v, err := g.View("commits")
	if err != nil {
		panic(err)
	}
	if len(state.Commits) == 0 {
		return Commit{}, ErrNoCommits
	}
	lineNumber := getItemPosition(v)
	if lineNumber > len(state.Commits)-1 {
		colorLog(color.FgRed, "potential error in getSelected Commit (mismatched ui and state)", state.Commits, lineNumber)
		return state.Commits[len(state.Commits)-1], nil
	}
	return state.Commits[lineNumber], nil
}
