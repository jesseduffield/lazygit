package helpers

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type WorktreeHelper struct {
	c                 *HelperCommon
	reposHelper       *ReposHelper
	refsHelper        *RefsHelper
	suggestionsHelper *SuggestionsHelper
}

func NewWorktreeHelper(c *HelperCommon, reposHelper *ReposHelper, refsHelper *RefsHelper, suggestionsHelper *SuggestionsHelper) *WorktreeHelper {
	return &WorktreeHelper{
		c:                 c,
		reposHelper:       reposHelper,
		refsHelper:        refsHelper,
		suggestionsHelper: suggestionsHelper,
	}
}

func (self *WorktreeHelper) GetMainWorktreeName() string {
	for _, worktree := range self.c.Model().Worktrees {
		if worktree.IsMain {
			return worktree.Name
		}
	}

	return ""
}

// If we're on the main worktree, we return an empty string
func (self *WorktreeHelper) GetLinkedWorktreeName() string {
	worktrees := self.c.Model().Worktrees
	if len(worktrees) == 0 {
		return ""
	}

	// worktrees always have the current worktree on top
	currentWorktree := worktrees[0]
	if currentWorktree.IsMain {
		return ""
	}

	return currentWorktree.Name
}

func (self *WorktreeHelper) NewWorktree() error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.NewWorktreeForBranchTitle,
		FindSuggestionsFunc: self.suggestionsHelper.GetWorktreeBranchNameSuggestionsFunc(),
		HandleConfirm: func(value string) error {
			return self.newWorktreeForPickerValue(value)
		},
	})

	return nil
}

// newWorktreeForPickerValue classifies the value the user picked or typed in the
// worktrees-panel picker and routes to the matching creation flow:
//   - an existing local branch -> a worktree that checks it out;
//   - a remote branch -> a new local tracking branch + worktree;
//   - anything else -> a new branch off the current ref + worktree.
//
// All three then feed the shared location menu. The picker filters out branches
// already checked out somewhere, but a verbatim type-in is still guarded here.
func (self *WorktreeHelper) newWorktreeForPickerValue(value string) error {
	if branch, ok := lo.Find(self.c.Model().Branches, func(branch *models.Branch) bool {
		return branch.Name == value
	}); ok {
		if worktree, ok := git_commands.WorktreeForBranch(branch, self.c.Model().Worktrees); ok {
			return errors.New(utils.ResolvePlaceholderString(self.c.Tr.BranchCheckedOutByWorktree,
				map[string]string{"branchName": branch.Name, "worktreeName": worktree.Name}))
		}

		prompt := utils.ResolvePlaceholderString(self.c.Tr.WorktreeLocationPromptCheckout,
			map[string]string{"branchName": branch.Name})
		return self.promptForWorktreeLocation(branch.Name, prompt, func(path string) error {
			return self.createWorktree(git_commands.NewWorktreeOpts{Path: path, Base: branch.RefName()}, context.WORKTREES_CONTEXT_KEY)
		})
	}

	if _, branchName, ok := self.refsHelper.ParseRemoteBranchName(value); ok {
		prompt := utils.ResolvePlaceholderString(self.c.Tr.WorktreeLocationPromptTrackingBranch,
			map[string]string{"name": branchName, "ref": value})
		return self.promptForWorktreeLocation(branchName, prompt, func(path string) error {
			return self.createWorktree(git_commands.NewWorktreeOpts{Path: path, Base: value, Branch: branchName}, context.WORKTREES_CONTEXT_KEY)
		})
	}

	name := SanitizedBranchName(value)
	base := self.refsHelper.GetCheckedOutRef().RefName()
	prompt := utils.ResolvePlaceholderString(self.c.Tr.WorktreeLocationPromptNewBranch,
		map[string]string{"name": name, "base": base})
	return self.promptForWorktreeLocation(name, prompt, func(path string) error {
		return self.createWorktree(git_commands.NewWorktreeOpts{Path: path, Base: base, Branch: name}, context.WORKTREES_CONTEXT_KEY)
	})
}

func (self *WorktreeHelper) Switch(worktree *models.Worktree, contextKey types.ContextKey) error {
	if worktree.IsCurrent {
		return errors.New(self.c.Tr.AlreadyInWorktree)
	}

	self.c.LogAction(self.c.Tr.SwitchToWorktree)

	return self.reposHelper.DispatchSwitchTo(worktree.Path, self.c.Tr.ErrWorktreeMovedOrRemoved, contextKey)
}

func (self *WorktreeHelper) Remove(worktree *models.Worktree, force bool) error {
	title := self.c.Tr.RemoveWorktreeTitle
	var templateStr string
	if force {
		templateStr = self.c.Tr.ForceRemoveWorktreePrompt
	} else {
		templateStr = self.c.Tr.RemoveWorktreePrompt
	}
	message := utils.ResolvePlaceholderString(
		templateStr,
		map[string]string{
			"worktreeName": worktree.Name,
		},
	)

	self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.RemovingWorktree, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.RemoveWorktree)
				if err := self.c.Git().Worktree.Delete(worktree.Path, force); err != nil {
					errMessage := err.Error()
					if !strings.Contains(errMessage, "--force") &&
						!strings.Contains(errMessage, "fatal: working trees containing submodules cannot be moved or removed") {
						return err
					}

					if !force {
						return self.Remove(worktree, true)
					}
					return err
				}
				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.WORKTREES, types.BRANCHES, types.FILES}})
				return nil
			})
		},
	})

	return nil
}

func (self *WorktreeHelper) Detach(worktree *models.Worktree) error {
	return self.c.WithWaitingStatus(self.c.Tr.DetachingWorktree, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.RemovingWorktree)

		err := self.c.Git().Worktree.Detach(worktree.Path)
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.WORKTREES, types.BRANCHES, types.FILES}})
		return nil
	})
}

// worktreeParentDirCandidates returns the candidate parent directories in which
// to create a new worktree, in priority order and de-duplicated:
//
//  1. the parent directory of each existing linked worktree (in worktree order);
//  2. the configured default path (relative paths are resolved against repoPath);
//  3. the repo's parent directory, if nothing else is available.
//
// repoPath is RepoPaths.RepoPath(), which is stable regardless of which worktree
// we're currently standing in. All returned paths are absolute.
func worktreeParentDirCandidates(repoPath string, linkedWorktreePaths []string, defaultPath string) []string {
	candidates := lo.Map(linkedWorktreePaths, func(path string, _ int) string {
		return filepath.Dir(path)
	})

	if defaultPath != "" {
		if filepath.IsAbs(defaultPath) {
			defaultPath = filepath.Clean(defaultPath)
		} else {
			defaultPath = filepath.Join(repoPath, defaultPath)
		}
		candidates = append(candidates, defaultPath)
	}

	candidates = lo.Uniq(candidates)

	if len(candidates) == 0 {
		candidates = append(candidates, filepath.Dir(repoPath))
	}

	return candidates
}

func (self *WorktreeHelper) NewWorktreeMenuForBranch(branch *models.Branch) error {
	return self.worktreeMenu(
		self.newBranchAndWorktreeItem(branch.Name, branch.RefName()),
		self.worktreeForBranchItem(branch),
		self.detachedWorktreeItem(branch.Name, branch.RefName(), branch.Name),
	)
}

func (self *WorktreeHelper) NewWorktreeMenuForCommit(commit *models.Commit) error {
	return self.worktreeMenu(
		self.newBranchAndWorktreeItem(commit.ShortHash(), commit.RefName()),
		self.detachedWorktreeItem(commit.ShortHash(), commit.RefName(), ""),
	)
}

func (self *WorktreeHelper) NewWorktreeMenuForTag(tag *models.Tag) error {
	return self.worktreeMenu(
		self.newBranchAndWorktreeItem(tag.Name, tag.RefName()),
		self.detachedWorktreeItem(tag.Name, tag.RefName(), ""),
	)
}

func (self *WorktreeHelper) NewWorktreeMenuForStash(stash *models.StashEntry) error {
	return self.worktreeMenu(
		self.newBranchAndWorktreeItem(stash.RefName(), stash.FullRefName()),
		self.detachedWorktreeItem(stash.RefName(), stash.FullRefName(), ""),
	)
}

func (self *WorktreeHelper) NewWorktreeMenuForRemoteBranch(remoteBranch *models.RemoteBranch) error {
	// e.g. "origin/foo" -> "foo": the local branch's (and worktree's) default name
	strippedName := strings.SplitAfterN(remoteBranch.RefName(), "/", 2)[1]

	return self.worktreeMenu(
		self.newLocalBranchAndWorktreeItem(remoteBranch, strippedName),
		self.detachedWorktreeItem(remoteBranch.FullName(), remoteBranch.FullName(), strippedName),
	)
}

func (self *WorktreeHelper) worktreeMenu(items ...*types.MenuItem) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.NewWorktree,
		Items: items,
	})
}

// newBranchAndWorktreeItem is the "new branch + worktree" action for a ref that
// isn't a remote branch (a branch, commit, tag or stash). ref is what we show the
// user (branch name, short hash, stash@{n}, ...); base is what we hand to
// `git worktree add`.
func (self *WorktreeHelper) newBranchAndWorktreeItem(ref string, base string) *types.MenuItem {
	return &types.MenuItem{
		Label: utils.ResolvePlaceholderString(self.c.Tr.NewBranchAndWorktreeFromRef, map[string]string{"ref": ref}),
		Keys:  menuKey('b'),
		OnPress: func() error {
			return self.startNewBranchWorktree("", base, func(name string) string {
				return utils.ResolvePlaceholderString(self.c.Tr.WorktreeLocationPromptNewBranch,
					map[string]string{"name": name, "base": ref})
			})
		},
	}
}

// newLocalBranchAndWorktreeItem is the "new branch + worktree" action for a remote
// branch: the new local branch tracks the remote one, and its name defaults to the
// remote branch name with the remote stripped off.
func (self *WorktreeHelper) newLocalBranchAndWorktreeItem(remoteBranch *models.RemoteBranch, strippedName string) *types.MenuItem {
	return &types.MenuItem{
		Label: utils.ResolvePlaceholderString(self.c.Tr.NewLocalBranchAndWorktreeFromRef, map[string]string{"ref": remoteBranch.FullName()}),
		Keys:  menuKey('b'),
		OnPress: func() error {
			return self.startNewBranchWorktree(strippedName, remoteBranch.FullName(), func(name string) string {
				return utils.ResolvePlaceholderString(self.c.Tr.WorktreeLocationPromptTrackingBranch,
					map[string]string{"name": name, "ref": remoteBranch.FullName()})
			})
		},
	}
}

// startNewBranchWorktree runs the shared name -> location -> create pipeline for the
// new-branch actions: prompt for the branch (and worktree) name, ask for the
// location, then create a worktree on a freshly created branch of that name.
// locationPrompt builds the location-menu prompt once the name is known.
func (self *WorktreeHelper) startNewBranchWorktree(nameInitialContent string, base string, locationPrompt func(name string) string) error {
	return self.promptForName(self.c.Tr.NewBranchAndWorktreeName, nameInitialContent, func(name string) error {
		return self.promptForWorktreeLocation(name, locationPrompt(name), func(path string) error {
			return self.createWorktree(git_commands.NewWorktreeOpts{Path: path, Base: base, Branch: name}, context.LOCAL_BRANCHES_CONTEXT_KEY)
		})
	})
}

// worktreeForBranchItem is the "check out an existing branch in a new worktree"
// action. It's disabled when the branch is already checked out somewhere.
func (self *WorktreeHelper) worktreeForBranchItem(branch *models.Branch) *types.MenuItem {
	return &types.MenuItem{
		Label: utils.ResolvePlaceholderString(self.c.Tr.WorktreeForRef, map[string]string{"ref": branch.Name}),
		Keys:  menuKey('w'),
		OnPress: func() error {
			prompt := utils.ResolvePlaceholderString(self.c.Tr.WorktreeLocationPromptCheckout,
				map[string]string{"branchName": branch.Name})
			return self.promptForWorktreeLocation(branch.Name, prompt, func(path string) error {
				return self.createWorktree(git_commands.NewWorktreeOpts{Path: path, Base: branch.RefName()}, context.LOCAL_BRANCHES_CONTEXT_KEY)
			})
		},
		DisabledReason: self.branchCheckedOutDisabledReason(branch),
	}
}

// detachedWorktreeItem is the "detached worktree at a ref" action. ref is shown to
// the user; base is handed to `git worktree add`. defaultDirName is the worktree
// directory name to use; when it's empty we prompt for one instead (commits, tags
// and stashes have no good name to derive).
func (self *WorktreeHelper) detachedWorktreeItem(ref string, base string, defaultDirName string) *types.MenuItem {
	prompt := utils.ResolvePlaceholderString(self.c.Tr.WorktreeLocationPromptDetached, map[string]string{"ref": ref})
	create := func(path string) error {
		return self.createWorktree(git_commands.NewWorktreeOpts{Path: path, Base: base, Detach: true}, context.LOCAL_BRANCHES_CONTEXT_KEY)
	}

	return &types.MenuItem{
		Label: utils.ResolvePlaceholderString(self.c.Tr.DetachedWorktreeAtRef, map[string]string{"ref": ref}),
		Keys:  menuKey('d'),
		OnPress: func() error {
			if defaultDirName != "" {
				return self.promptForWorktreeLocation(defaultDirName, prompt, create)
			}
			return self.promptForName(self.c.Tr.NewWorktreeName, "", func(name string) error {
				return self.promptForWorktreeLocation(name, prompt, create)
			})
		},
	}
}

func (self *WorktreeHelper) branchCheckedOutDisabledReason(branch *models.Branch) *types.DisabledReason {
	if worktree, ok := git_commands.WorktreeForBranch(branch, self.c.Model().Worktrees); ok {
		return &types.DisabledReason{
			Text: utils.ResolvePlaceholderString(self.c.Tr.BranchCheckedOutByWorktree,
				map[string]string{"branchName": branch.Name, "worktreeName": worktree.Name}),
		}
	}

	return nil
}

// promptForName asks for a branch/worktree name and sanitizes the response (most
// notably turning spaces into dashes so it's a valid branch name) before
// continuing.
func (self *WorktreeHelper) promptForName(title string, initialContent string, onConfirm func(name string) error) error {
	self.c.Prompt(types.PromptOpts{
		Title:          title,
		InitialContent: initialContent,
		HandleConfirm: func(response string) error {
			return onConfirm(SanitizedBranchName(response))
		},
	})

	return nil
}

// promptForWorktreeLocation shows the location menu: one item per candidate parent
// directory (each labelled with the absolute path the worktree would end up at),
// plus an "Other…" item that opens a free-form path prompt. The chosen absolute
// path is passed to onConfirm.
func (self *WorktreeHelper) promptForWorktreeLocation(dirName string, prompt string, onConfirm func(path string) error) error {
	linkedWorktreePaths := []string{}
	for _, worktree := range self.c.Model().Worktrees {
		if !worktree.IsMain {
			linkedWorktreePaths = append(linkedWorktreePaths, worktree.Path)
		}
	}
	parentDirs := worktreeParentDirCandidates(
		self.c.Git().RepoPaths.RepoPath(),
		linkedWorktreePaths,
		self.c.UserConfig().Worktree.DefaultPath,
	)

	targets := lo.Map(parentDirs, func(parentDir string, _ int) string {
		return filepath.Join(parentDir, dirName)
	})

	menuItems := lo.Map(targets, func(target string, _ int) *types.MenuItem {
		return &types.MenuItem{
			Label:   target,
			OnPress: func() error { return onConfirm(target) },
		}
	})

	menuItems = append(menuItems, &types.MenuItem{
		Label: self.c.Tr.WorktreeLocationOther,
		OnPress: func() error {
			self.c.Prompt(types.PromptOpts{
				Title:          self.c.Tr.NewWorktreePath,
				InitialContent: targets[0],
				HandleConfirm:  onConfirm,
			})
			return nil
		},
	})

	return self.c.Menu(types.CreateMenuOptions{
		Title:  self.c.Tr.WorktreeLocationTitle,
		Prompt: prompt,
		Items:  menuItems,
	})
}

func (self *WorktreeHelper) createWorktree(opts git_commands.NewWorktreeOpts, contextKey types.ContextKey) error {
	return self.c.WithWaitingStatus(self.c.Tr.AddingWorktree, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.AddWorktree)
		if err := self.c.Git().Worktree.New(opts); err != nil {
			return err
		}

		return self.reposHelper.DispatchSwitchTo(opts.Path, self.c.Tr.ErrWorktreeMovedOrRemoved, contextKey)
	})
}
