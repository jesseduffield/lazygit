package helpers

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	appTypes "github.com/jesseduffield/lazygit/pkg/app/types"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/direnv"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type onNewRepoFn func(startArgs appTypes.StartArgs, contextKey types.ContextKey) error

// helps switch back and forth between repos
type ReposHelper struct {
	c                     *HelperCommon
	recordDirectoryHelper *RecordDirectoryHelper
	onNewRepo             onNewRepoFn
}

func NewRecentReposHelper(
	c *HelperCommon,
	recordDirectoryHelper *RecordDirectoryHelper,
	onNewRepo onNewRepoFn,
) *ReposHelper {
	return &ReposHelper{
		c:                     c,
		recordDirectoryHelper: recordDirectoryHelper,
		onNewRepo:             onNewRepo,
	}
}

func (self *ReposHelper) EnterSubmodule(submodule *models.SubmoduleConfig) error {
	// Check before pushing onto the repo-path stack, so a refused switch
	// doesn't leave a stale entry there (which escape would later switch back
	// to, needlessly reloading the current repo).
	if self.switchRefusedBecauseBusy() {
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	self.c.State().GetRepoPathStack().Push(types.RepoLocation{
		Path:               wd,
		GitLocationEnvVars: self.c.Git().RepoPaths.GitLocationEnvVars(),
	})

	return self.switchTo(submodule.FullPath(), self.c.Tr.ErrRepositoryMovedOrDeleted, context.NO_CONTEXT)
}

// What a repo has checked out. Exactly one of the two fields is set.
type headInfo struct {
	// The name of the checked-out branch.
	branch string
	// The commit HEAD is detached at.
	hash string
}

// gitDirOfRepo returns the directory that holds the git data of the repo at
// repoPath, and whether it could be found. An ordinary repo keeps that data in
// a .git directory; a worktree and a submodule have a .git file naming the
// directory instead.
func gitDirOfRepo(repoPath string) (string, bool) {
	gitDirPath := filepath.Join(repoPath, ".git")

	stat, err := os.Stat(gitDirPath)
	if err != nil {
		return "", false
	}
	if stat.IsDir() {
		return gitDirPath, true
	}

	content, err := os.ReadFile(gitDirPath)
	if err != nil {
		return "", false
	}
	gitDir, ok := strings.CutPrefix(strings.TrimSpace(string(content)), "gitdir: ")
	if !ok {
		return "", false
	}
	// A relative name is relative to the repo. Git writes one for a submodule,
	// and for a worktree created with --relative-paths.
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(repoPath, gitDir)
	}
	return gitDir, true
}

// The branch git names in the HEAD file of a repo that keeps its refs in a
// reftable. The refs live in a binary table there, and the name in HEAD
// resolves nowhere, so that a reader of the file gets an error instead of a
// stale answer.
const reftablePlaceholderBranch = ".invalid"

// readHeadInfo reads the HEAD file of the repo at repoPath to find out what it
// has checked out, and reports whether that worked.
func readHeadInfo(repoPath string) (headInfo, bool) {
	gitDir, ok := gitDirOfRepo(repoPath)
	if !ok {
		return headInfo{}, false
	}

	content, err := os.ReadFile(filepath.Join(gitDir, "HEAD"))
	if err != nil {
		return headInfo{}, false
	}

	head := strings.TrimSpace(string(content))
	if branch, ok := strings.CutPrefix(head, "ref: refs/heads/"); ok {
		if branch == reftablePlaceholderBranch {
			return headInfo{}, false
		}
		return headInfo{branch: branch}, true
	}
	return headInfo{hash: head}, true
}

// askGitForHeadInfo asks git what the repo at repoPath has checked out. This
// costs a process, so only the repos that readHeadInfo can't answer for go
// through here.
func (self *ReposHelper) askGitForHeadInfo(repoPath string) (headInfo, bool) {
	// symbolic-ref names the branch even when it has no commit yet, and
	// rev-parse resolves HEAD when it is detached. Neither can do the other's
	// job, so ask for the branch first and only then for the commit.
	if branch, ok := self.askGit(repoPath, "symbolic-ref", "--short", "--quiet", "HEAD"); ok {
		return headInfo{branch: branch}, true
	}
	if hash, ok := self.askGit(repoPath, "rev-parse", "HEAD"); ok {
		return headInfo{hash: hash}, true
	}
	return headInfo{}, false
}

// askGit runs a git command against the repo at repoPath and returns the one
// line it writes, or false if it fails or writes nothing.
func (self *ReposHelper) askGit(repoPath string, subcommand string, args ...string) (string, bool) {
	cmdObj := self.c.OS().Cmd.New(git_commands.NewGitCmd(subcommand).
		Dir(repoPath).
		Arg(args...).
		ToArgv()).DontLog()
	stdout, _, err := git_commands.ForOtherRepo(cmdObj).RunWithOutputs()
	if err != nil {
		return "", false
	}
	output := strings.TrimSpace(stdout)
	return output, output != ""
}

func (self *ReposHelper) getCurrentBranch(path string) string {
	head, ok := readHeadInfo(path)
	if !ok {
		head, ok = self.askGitForHeadInfo(path)
	}
	if !ok {
		return self.c.Tr.BranchUnknown
	}
	if head.branch != "" {
		return head.branch
	}
	return utils.ResolvePlaceholderString(self.c.Tr.HeadDetachedAt,
		map[string]string{"hash": utils.ShortHash(head.hash)})
}

// The most that the name and the branch column of the recent repos menu are
// allowed to take up. Each column is padded to the width of its widest entry,
// so without a limit one long entry pushes the columns after it off the right
// edge of the menu for every entry.
const (
	recentReposNameMaxWidth   = 30
	recentReposBranchMaxWidth = 30
)

// One entry of the recent repos menu.
type recentRepoEntry struct {
	// What the entry stands for
	path       string
	branchName string

	// The text of its three columns
	nameColumn   string
	branchColumn string
	dirColumn    string
}

func newRecentRepoEntry(path string, branchName string) recentRepoEntry {
	// The icon is part of the column, so it counts towards the column's width.
	branchColumn := branchName
	if icons.IsIconEnabled() {
		branchColumn = icons.BRANCH_ICON + " " + branchName
	}

	return recentRepoEntry{
		path:         path,
		branchName:   branchName,
		nameColumn:   filepath.Base(path),
		branchColumn: branchColumn,
		// The last segment of the path is already in the first column, so the
		// directory that contains the repo is enough to tell repos with the
		// same name apart.
		dirColumn: utils.ContractTilde(filepath.Dir(path)),
	}
}

// How wide the columns of the recent repos menu are allowed to get.
type recentRepoColumnWidths struct {
	name   int
	branch int
	dir    int
}

// Gives the name and the branch column as much as their entries need, up to
// their respective maximum, and the rest of the row to the directory column.
// Measuring the entries first matters because most users don't have names that
// long; truncating the directories as if they did would cut them short for no
// reason.
func (self *ReposHelper) fitRecentRepoColumns(entries []recentRepoEntry) recentRepoColumnWidths {
	// The menu appends a Cancel entry, whose label sits in the first column.
	nameWidth := utils.StringWidth(self.c.Tr.Cancel)
	branchWidth := 0
	for _, entry := range entries {
		nameWidth = max(nameWidth, utils.StringWidth(entry.nameColumn))
		branchWidth = max(branchWidth, utils.StringWidth(entry.branchColumn))
	}

	nameWidth = min(nameWidth, recentReposNameMaxWidth)
	branchWidth = min(branchWidth, recentReposBranchMaxWidth)

	return recentRepoColumnWidths{
		name:   nameWidth,
		branch: branchWidth,
		// The menu's frame takes up two columns, and two more separate the
		// three columns from each other. What's left is the room a whole row
		// has in a menu that is as wide as it gets.
		dir: menuMaxWidth - 2 - 2 - nameWidth - branchWidth,
	}
}

func (self *ReposHelper) recentRepoMenuItem(entry recentRepoEntry, widths recentRepoColumnWidths) *types.MenuItem {
	displayedName := utils.TruncateWithEllipsis(entry.nameColumn, widths.name)
	displayedBranch := utils.TruncateWithEllipsis(entry.branchColumn, widths.branch)
	// The beginning and the end of a directory are both worth seeing, so it
	// loses its middle rather than its end when it doesn't fit.
	displayedDir := utils.TruncateWithEllipsisInMiddle(entry.dirColumn, widths.dir)

	// Spell out whatever the columns show in truncated form
	type tooltipField struct {
		label string
		value string
	}
	fields := []tooltipField{}
	addTooltipField := func(label string, value string) {
		fields = append(fields, tooltipField{label: label, value: value})
	}

	if displayedName != entry.nameColumn {
		addTooltipField(self.c.Tr.RecentReposRepoLabel, entry.nameColumn)
	}
	if displayedBranch != entry.branchColumn {
		addTooltipField(self.c.Tr.RecentReposBranchLabel, entry.branchName)
	}
	if displayedDir != entry.dirColumn {
		addTooltipField(self.c.Tr.RecentReposPathLabel, entry.dirColumn)
	}

	// Line the values up behind the widest of the labels that are there
	labelWidth := utils.MaxFn(fields, func(field tooltipField) int {
		return utils.StringWidth(field.label)
	})
	tooltipLines := lo.Map(fields, func(field tooltipField, _ int) string {
		return utils.WithPadding(field.label, labelWidth, utils.AlignLeft) + " " + field.value
	})

	return &types.MenuItem{
		LabelColumns: []string{
			displayedName,
			style.FgCyan.Sprint(displayedBranch),
			style.FgMagenta.Sprint(displayedDir),
		},
		// Filtering matches the full text, including the parts that the columns
		// above truncate or leave out.
		FilterColumns: []string{entry.nameColumn, entry.branchName, entry.path},
		Tooltip:       strings.Join(tooltipLines, "\n"),
		OnPress: func() error {
			// Check before clearing the stack, so a refused switch doesn't
			// forget the submodule breadcrumb (which would leave escape
			// unable to return to the parent repo).
			if self.switchRefusedBecauseBusy() {
				return nil
			}
			// if we were in a submodule, we want to forget about that stack of repos
			// so that hitting escape in the new repo does nothing
			self.c.State().GetRepoPathStack().Clear()
			return self.switchTo(entry.path, self.c.Tr.ErrRepositoryMovedOrDeleted, context.NO_CONTEXT)
		},
	}
}

func (self *ReposHelper) CreateRecentReposMenu() error {
	// we'll show an empty panel if there are no recent repos
	recentRepoPaths := []string{}
	if len(self.c.GetAppState().RecentRepos) > 0 {
		// we skip the first one because we're currently in it
		recentRepoPaths = self.c.GetAppState().RecentRepos[1:]
	}

	currentBranches := make([]string, len(recentRepoPaths))

	wg := sync.WaitGroup{}
	wg.Add(len(recentRepoPaths))

	for i, path := range recentRepoPaths {
		go func() {
			defer wg.Done()
			currentBranches[i] = self.getCurrentBranch(path)
		}()
	}

	wg.Wait()

	entries := lo.Map(recentRepoPaths, func(path string, i int) recentRepoEntry {
		return newRecentRepoEntry(path, currentBranches[i])
	})

	columnWidths := self.fitRecentRepoColumns(entries)
	menuItems := lo.Map(entries, func(entry recentRepoEntry, _ int) *types.MenuItem {
		return self.recentRepoMenuItem(entry, columnWidths)
	})

	return self.c.Menu(types.CreateMenuOptions{
		Title:           self.c.Tr.RecentRepos,
		Items:           menuItems,
		FilterAsYouType: true,
	})
}

// SwitchToParentRepo switches back to the repo the current submodule was
// entered from (the top of the repo-path stack). Like the other callers that do
// work before switching, it checks for an in-flight operation *before* popping
// the stack, so a refused switch leaves the stack intact — otherwise the entry
// would be consumed and escape would no longer return to the parent once the
// operation finished. The caller must only call this when the stack is
// non-empty.
func (self *ReposHelper) SwitchToParentRepo() error {
	if self.switchRefusedBecauseBusy() {
		return nil
	}
	return self.switchToLocation(self.c.State().GetRepoPathStack().Pop(), self.c.Tr.ErrRepositoryMovedOrDeleted, context.NO_CONTEXT)
}

func (self *ReposHelper) DispatchSwitchTo(path string, errMsg string, contextKey types.ContextKey) error {
	if self.switchRefusedBecauseBusy() {
		return nil
	}
	return self.switchTo(path, errMsg, contextKey)
}

// switchRefusedBecauseBusy reports (and shows a toast) whether a repo switch
// must be refused because a foreground git operation is in flight. Switching
// reassigns gui.git and the process cwd, so switching mid-operation would run
// the operation's remaining git commands against the wrong repo. Callers that
// do work before the switch (creating a worktree, recording the repo-path
// stack) check this up front, so they don't do that work only to have the
// switch refused; the switch itself (switchTo) is then unguarded.
func (self *ReposHelper) switchRefusedBecauseBusy() bool {
	if self.c.GocuiGui().Busy() {
		self.c.ErrorToast(self.c.Tr.CantSwitchWhileOperationInProgress)
		return true
	}
	return false
}

// switchTo switches lazygit to the repository (or worktree) at the given path,
// which git is expected to find from that path alone. That's true of every repo
// we switch to without having been there before.
func (self *ReposHelper) switchTo(path string, errMsg string, contextKey types.ContextKey) error {
	return self.switchToLocation(types.RepoLocation{Path: path}, errMsg, contextKey)
}

// switchToLocation switches lazygit to the repository (or worktree) at the
// given location. It runs synchronously on the UI thread: the switch swaps
// gui.State (in resetState) and reassigns gui.git and the process cwd, all of
// which the UI thread also reads, so doing it here rather than on a worker
// avoids racing those reads. The heavy data loading is still dispatched
// asynchronously by the refresh that onNewRepo kicks off.
//
// Everything from here on has to find the repo the way git does, from the
// directory we're about to change to, so the location's environment goes into
// the process env before we do. Usually that just clears whatever the repo
// we're leaving needed, but going back to a repo whose git dir isn't in its
// work tree (a dotfile repo opened with --git-dir/--work-tree, say) is the
// reason we remember the environment at all: nothing in the path leads to its
// git dir. On failure we put back what the repo we're staying in needs.
func (self *ReposHelper) switchToLocation(location types.RepoLocation, errMsg string, contextKey types.ContextKey) error {
	originalPath, err := os.Getwd()
	if err != nil {
		return nil
	}
	originalGitLocationEnvVars := env.GetGitLocationEnvVars()

	env.SetGitLocationEnvVars(location.GitLocationEnvVars)

	msg := utils.ResolvePlaceholderString(self.c.Tr.ChangingDirectoryTo, map[string]string{"path": location.Path})
	self.c.LogCommand(msg, false)

	if err := os.Chdir(location.Path); err != nil {
		env.SetGitLocationEnvVars(originalGitLocationEnvVars)
		if os.IsNotExist(err) {
			return errors.New(errMsg)
		}
		return err
	}

	if err := commands.VerifyInGitRepo(self.c.OS()); err != nil {
		env.SetGitLocationEnvVars(originalGitLocationEnvVars)
		if err := os.Chdir(originalPath); err != nil {
			return err
		}

		return err
	}

	direnvResult := self.logDirenvResult(direnv.Load(self.c.OS().Cmd))

	if err := self.recordDirectoryHelper.RecordCurrentDirectory(); err != nil {
		self.c.Log.Errorf("error recording current directory: %v", err)
	}

	if err := self.onNewRepo(appTypes.StartArgs{}, contextKey); err != nil {
		return err
	}

	if direnvResult.Blocked {
		self.promptDirenvApproval(direnvResult.EnvrcPath)
		return nil
	}

	return direnvResult.Err
}

// logDirenvResult writes whatever direnv emitted to the command log and the
// debug log; both happen for every load attempt regardless of outcome.
func (self *ReposHelper) logDirenvResult(result direnv.LoadResult) direnv.LoadResult {
	if result.Message != "" {
		self.c.LogCommand(result.Message, false)
	}
	if result.Err != nil {
		self.c.Log.WithError(result.Err).Warn("direnv load failed")
	}
	return result
}

// promptDirenvApproval shows the user the contents of an unapproved .envrc
// and offers to run `direnv allow` for them. On confirm, we approve the
// file and re-run Load so the new env reaches subprocesses; on cancel we
// leave the env as-is (the previous repo's vars are already unloaded by
// the initial Load call, which is the correct state).
func (self *ReposHelper) promptDirenvApproval(envrcPath string) {
	content, err := os.ReadFile(envrcPath)
	if err != nil {
		self.c.Log.WithError(err).Warn("could not read .envrc for approval prompt")
		return
	}

	indented := "  " + strings.ReplaceAll(strings.TrimRight(string(content), "\n"), "\n", "\n  ")
	prompt := utils.ResolvePlaceholderString(self.c.Tr.DirenvApprovalPrompt, map[string]string{
		"confirmKey": self.c.UserConfig().Keybinding.Universal.Confirm.String(),
		"cancelKey":  self.c.UserConfig().Keybinding.Universal.Return.String(),
		"content":    indented,
	})

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DirenvApprovalTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			if err := direnv.Allow(self.c.OS().Cmd, envrcPath); err != nil {
				return err
			}
			return self.logDirenvResult(direnv.Load(self.c.OS().Cmd)).Err
		},
	})
}
