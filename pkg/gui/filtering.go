package gui

import (
	"os"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/jesseduffield/minimal/gitignore"
	"gopkg.in/ozeidan/fuzzy-patricia.v3/patricia"
)

func (gui *Gui) validateNotInFilterMode() (bool, error) {
	if gui.State.Modes.Filtering.Active() {
		err := gui.ask(askOpts{
			title:         gui.Tr.MustExitFilterModeTitle,
			prompt:        gui.Tr.MustExitFilterModePrompt,
			handleConfirm: gui.exitFilterMode,
		})

		return false, err
	}
	return true, nil
}

func (gui *Gui) exitFilterMode() error {
	return gui.clearFiltering()
}

func (gui *Gui) clearFiltering() error {
	gui.State.Modes.Filtering.Reset()
	if gui.State.ScreenMode == SCREEN_HALF {
		gui.State.ScreenMode = SCREEN_NORMAL
	}

	return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{COMMITS}})
}

func (gui *Gui) setFiltering(path string) error {
	gui.State.Modes.Filtering.SetPath(path)
	if gui.State.ScreenMode == SCREEN_NORMAL {
		gui.State.ScreenMode = SCREEN_HALF
	}

	if err := gui.pushContext(gui.State.Contexts.BranchCommits); err != nil {
		return err
	}

	return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{COMMITS}, then: func() {
		gui.State.Contexts.BranchCommits.GetPanelState().SetSelectedLineIdx(0)
	}})
}

// here we asynchronously fetch the latest set of paths in the repo and store in
// gui.State.FilesTrie. On the main thread we'll be doing a fuzzy search via
// gui.State.FilesTrie. So if we've looked for a file previously, we'll start with
// the old trie and eventually it'll be swapped out for the new one.
func (gui *Gui) getFindSuggestionsForFilterPath() func(string) []*types.Suggestion {
	_ = gui.WithWaitingStatus(gui.Tr.LcLoadingFileSuggestions, func() error {
		trie := patricia.NewTrie()
		// load every non-gitignored file in the repo
		ignore, err := gitignore.FromGit()
		if err != nil {
			return err
		}

		err = ignore.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				trie.Insert(patricia.Prefix(path), path)
				return nil
			})
		// cache the trie for future use
		gui.State.FilesTrie = trie

		// refresh the selections view
		gui.suggestionsAsyncHandler.Do(func() func() {
			// assuming here that the confirmation view is what we're typing into.
			// This assumption may prove false over time
			suggestions := gui.findSuggestions(gui.Views.Confirmation.TextArea.GetContent())
			return func() { gui.setSuggestions(suggestions) }
		})

		return err
	})

	return func(input string) []*types.Suggestion {
		matchingNames := []string{}
		_ = gui.State.FilesTrie.VisitFuzzy(patricia.Prefix(input), true, func(prefix patricia.Prefix, item patricia.Item, skipped int) error {
			matchingNames = append(matchingNames, item.(string))
			return nil
		})

		// doing another fuzzy search for good measure
		matchingNames = utils.FuzzySearch(input, matchingNames)

		suggestions := make([]*types.Suggestion, len(matchingNames))
		for i, name := range matchingNames {
			suggestions[i] = &types.Suggestion{
				Value: name,
				Label: name,
			}
		}

		return suggestions
	}
}
