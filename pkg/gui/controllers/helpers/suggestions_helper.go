package helpers

import (
	"fmt"
	"os"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/jesseduffield/minimal/gitignore"
	"github.com/samber/lo"
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

type ISuggestionsHelper interface {
	GetRemoteSuggestionsFunc() func(string) []*types.Suggestion
	GetBranchNameSuggestionsFunc() func(string) []*types.Suggestion
	GetFilePathSuggestionsFunc() func(string) []*types.Suggestion
	GetRemoteBranchesSuggestionsFunc(separator string) func(string) []*types.Suggestion
	GetRefsSuggestionsFunc() func(string) []*types.Suggestion
}

type SuggestionsHelper struct {
	c *types.HelperCommon

	model                *types.Model
	refreshSuggestionsFn func()
}

var _ ISuggestionsHelper = &SuggestionsHelper{}

func NewSuggestionsHelper(
	c *types.HelperCommon,
	model *types.Model,
	refreshSuggestionsFn func(),
) *SuggestionsHelper {
	return &SuggestionsHelper{
		c:                    c,
		model:                model,
		refreshSuggestionsFn: refreshSuggestionsFn,
	}
}

func (self *SuggestionsHelper) getRemoteNames() []string {
	return slices.Map(self.model.Remotes, func(remote *models.Remote) string {
		return remote.Name
	})
}

func matchesToSuggestions(matches []string) []*types.Suggestion {
	return slices.Map(matches, func(match string) *types.Suggestion {
		return &types.Suggestion{
			Value: match,
			Label: match,
		}
	})
}

func (self *SuggestionsHelper) GetRemoteSuggestionsFunc() func(string) []*types.Suggestion {
	remoteNames := self.getRemoteNames()

	return FuzzySearchFunc(remoteNames)
}

func (self *SuggestionsHelper) getBranchNames() []string {
	return slices.Map(self.model.Branches, func(branch *models.Branch) string {
		return branch.Name
	})
}

func (self *SuggestionsHelper) GetBranchNameSuggestionsFunc() func(string) []*types.Suggestion {
	branchNames := self.getBranchNames()

	return func(input string) []*types.Suggestion {
		var matchingBranchNames []string
		if input == "" {
			matchingBranchNames = branchNames
		} else {
			matchingBranchNames = utils.FuzzySearch(input, branchNames)
		}

		return slices.Map(matchingBranchNames, func(branchName string) *types.Suggestion {
			return &types.Suggestion{
				Value: branchName,
				Label: presentation.GetBranchTextStyle(branchName).Sprint(branchName),
			}
		})
	}
}

// here we asynchronously fetch the latest set of paths in the repo and store in
// self.model.FilesTrie. On the main thread we'll be doing a fuzzy search via
// self.model.FilesTrie. So if we've looked for a file previously, we'll start with
// the old trie and eventually it'll be swapped out for the new one.
// Notably, unlike other suggestion functions we're not showing all the options
// if nothing has been typed because there'll be too much to display efficiently
func (self *SuggestionsHelper) GetFilePathSuggestionsFunc() func(string) []*types.Suggestion {
	_ = self.c.WithWaitingStatus(self.c.Tr.LcLoadingFileSuggestions, func() error {
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
		self.model.FilesTrie = trie

		self.refreshSuggestionsFn()

		return err
	})

	return func(input string) []*types.Suggestion {
		matchingNames := []string{}
		_ = self.model.FilesTrie.VisitFuzzy(patricia.Prefix(input), true, func(prefix patricia.Prefix, item patricia.Item, skipped int) error {
			matchingNames = append(matchingNames, item.(string))
			return nil
		})

		// doing another fuzzy search for good measure
		matchingNames = utils.FuzzySearch(input, matchingNames)

		return matchesToSuggestions(matchingNames)
	}
}

func (self *SuggestionsHelper) getRemoteBranchNames(separator string) []string {
	return slices.FlatMap(self.model.Remotes, func(remote *models.Remote) []string {
		return slices.Map(remote.Branches, func(branch *models.RemoteBranch) string {
			return fmt.Sprintf("%s%s%s", remote.Name, separator, branch.Name)
		})
	})
}

func (self *SuggestionsHelper) GetRemoteBranchesSuggestionsFunc(separator string) func(string) []*types.Suggestion {
	return FuzzySearchFunc(self.getRemoteBranchNames(separator))
}

func (self *SuggestionsHelper) getTagNames() []string {
	return slices.Map(self.model.Tags, func(tag *models.Tag) string {
		return tag.Name
	})
}

func (self *SuggestionsHelper) GetRefsSuggestionsFunc() func(string) []*types.Suggestion {
	remoteBranchNames := self.getRemoteBranchNames("/")
	localBranchNames := self.getBranchNames()
	tagNames := self.getTagNames()
	additionalRefNames := []string{"HEAD", "FETCH_HEAD", "MERGE_HEAD", "ORIG_HEAD"}

	refNames := append(append(append(remoteBranchNames, localBranchNames...), tagNames...), additionalRefNames...)

	return FuzzySearchFunc(refNames)
}

func (self *SuggestionsHelper) GetAuthorsSuggestionsFunc() func(string) []*types.Suggestion {
	authors := lo.Uniq(slices.Map(self.model.Commits, func(commit *models.Commit) string {
		return fmt.Sprintf("%s <%s>", commit.AuthorName, commit.AuthorEmail)
	}))

	return FuzzySearchFunc(authors)
}

func FuzzySearchFunc(options []string) func(string) []*types.Suggestion {
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
