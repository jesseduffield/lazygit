package gui

import (
	"fmt"
	"os"

	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/jesseduffield/minimal/gitignore"
	"gopkg.in/ozeidan/fuzzy-patricia.v3/patricia"
)

// Thinking out loud: I'm typically a staunch advocate of organising code by feature rather than type,
// because colocating code that relates to the same feature means far less effort
// to get all the context you need to work on any particular feature. But the one
// major benefit of grouping by type is that it makes it makes it less likely that
// somebody will re-implement the same logic twice, because they can quickly see
// if a certain method has been used for some use case, given that as a starting point
// they know about the type. In that vein, I'm including all our functions for
// finding suggestions in this file, so that it's easy to see if a function already
// exists for fetching a particular model.

func (gui *Gui) getRemoteNames() []string {
	result := make([]string, len(gui.State.Remotes))
	for i, remote := range gui.State.Remotes {
		result[i] = remote.Name
	}
	return result
}

func matchesToSuggestions(matches []string) []*types.Suggestion {
	suggestions := make([]*types.Suggestion, len(matches))
	for i, match := range matches {
		suggestions[i] = &types.Suggestion{
			Value: match,
			Label: match,
		}
	}
	return suggestions
}

func (gui *Gui) getRemoteSuggestionsFunc() func(string) []*types.Suggestion {
	remoteNames := gui.getRemoteNames()

	return fuzzySearchFunc(remoteNames)
}

func (gui *Gui) getBranchNames() []string {
	result := make([]string, len(gui.State.Branches))
	for i, branch := range gui.State.Branches {
		result[i] = branch.Name
	}
	return result
}

func (gui *Gui) getBranchNameSuggestionsFunc() func(string) []*types.Suggestion {
	branchNames := gui.getBranchNames()

	return func(input string) []*types.Suggestion {
		var matchingBranchNames []string
		if input == "" {
			matchingBranchNames = branchNames
		} else {
			matchingBranchNames = utils.FuzzySearch(input, branchNames)
		}

		suggestions := make([]*types.Suggestion, len(matchingBranchNames))
		for i, branchName := range matchingBranchNames {
			suggestions[i] = &types.Suggestion{
				Value: branchName,
				Label: presentation.GetBranchTextStyle(branchName).Sprint(branchName),
			}
		}

		return suggestions
	}
}

// here we asynchronously fetch the latest set of paths in the repo and store in
// gui.State.FilesTrie. On the main thread we'll be doing a fuzzy search via
// gui.State.FilesTrie. So if we've looked for a file previously, we'll start with
// the old trie and eventually it'll be swapped out for the new one.
// Notably, unlike other suggestion functions we're not showing all the options
// if nothing has been typed because there'll be too much to display efficiently
func (gui *Gui) getFilePathSuggestionsFunc() func(string) []*types.Suggestion {
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

func (gui *Gui) getRemoteBranchNames(separator string) []string {
	result := []string{}
	for _, remote := range gui.State.Remotes {
		for _, branch := range remote.Branches {
			result = append(result, fmt.Sprintf("%s%s%s", remote.Name, separator, branch.Name))
		}
	}
	return result
}

func (gui *Gui) getRemoteBranchesSuggestionsFunc(separator string) func(string) []*types.Suggestion {
	return fuzzySearchFunc(gui.getRemoteBranchNames(separator))
}

func (gui *Gui) getTagNames() []string {
	result := make([]string, len(gui.State.Tags))
	for i, tag := range gui.State.Tags {
		result[i] = tag.Name
	}
	return result
}

func (gui *Gui) getRefsSuggestionsFunc() func(string) []*types.Suggestion {
	remoteBranchNames := gui.getRemoteBranchNames("/")
	localBranchNames := gui.getBranchNames()
	tagNames := gui.getTagNames()
	additionalRefNames := []string{"HEAD", "FETCH_HEAD", "MERGE_HEAD", "ORIG_HEAD"}

	refNames := append(append(append(remoteBranchNames, localBranchNames...), tagNames...), additionalRefNames...)

	return fuzzySearchFunc(refNames)
}

func (gui *Gui) getCustomCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion {
	// reversing so that we display the latest command first
	history := utils.Reverse(gui.Config.GetAppState().CustomCommandsHistory)

	return fuzzySearchFunc(history)
}

func fuzzySearchFunc(options []string) func(string) []*types.Suggestion {
	return func(input string) []*types.Suggestion {
		var matches []string
		if input == "" {
			matches = options
		} else {
			matches = utils.FuzzySearch(input, options)
		}

		return matchesToSuggestions(matches)
	}
}
