package helpers

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type DiffHelper struct {
	c              *HelperCommon
	diffLineHelper *DiffLineHelper
	filesHelper    *FilesHelper

	// Diffs of the messages of "amend!" commits, keyed by everything that
	// shapes them: the two commits, the diff renderer, and the values its
	// command was resolved with. An empty diff means that the commit doesn't
	// change the message. Only accessed on the UI thread, so a diff produced
	// on a render's goroutine is stored by way of OnUIThread.
	commitMessageDiffs map[string]string
}

func NewDiffHelper(c *HelperCommon, diffLineHelper *DiffLineHelper, filesHelper *FilesHelper) *DiffHelper {
	return &DiffHelper{
		c:                  c,
		diffLineHelper:     diffLineHelper,
		filesHelper:        filesHelper,
		commitMessageDiffs: make(map[string]string),
	}
}

func (self *DiffHelper) DiffArgs() []string {
	output := []string{"--stat", "-p", self.c.Modes().Diffing.Ref}

	right := self.currentDiffTerminal()
	if right != "" {
		output = append(output, right)
	}

	if self.c.Modes().Diffing.Reverse {
		output = append(output, "-R")
	}

	output = append(output, "--")

	file := self.currentlySelectedFilename()
	if file != "" {
		output = append(output, file)
	} else if self.c.Modes().Filtering.Active() {
		output = append(output, self.c.Modes().Filtering.GetPath())
	}

	return output
}

// Returns an update task that can be passed to RenderToMainViews to render a
// diff for the selected commit(s). We need to pass both the selected commit
// and the refRange for a range selection. If the refRange is nil (meaning that
// either there's no range, or it can't be diffed for some reason), then we want
// to fall back to rendering the diff for the single commit.
// In addition, we need to pass the list of all commits; this is needed for
// showing the commit message diff for "amend!" commits.
func (self *DiffHelper) GetUpdateTaskForRenderingCommitsDiff(
	commits []*models.Commit,
	commit *models.Commit,
	refRange *types.RefRange,
) types.UpdateTask {
	mode := self.diffLineHelper.MainViewDiffMode()

	if refRange != nil {
		from, to := refRange.From, refRange.To
		args := []string{from.ParentRefName(), to.RefName(), "--stat", "-p"}
		args = append(args, "--")
		if filterPath := self.c.Modes().Filtering.GetPath(); filterPath != "" {
			// If both refs are commits, filter by the union of their paths. This is useful for
			// example when diffing a range of commits in filter-by-path mode across a rename.
			fromCommit, ok1 := from.(*models.Commit)
			toCommit, ok2 := to.(*models.Commit)
			if ok1 && ok2 {
				paths := append(self.FilterPathsForCommit(fromCommit), self.FilterPathsForCommit(toCommit)...)
				args = append(args, lo.Uniq(paths)...)
			} else {
				// If either ref is not a commit (which is possible in sticky diff mode, when
				// diffing against a branch or tag), we just filter by the filter path; that's the
				// best we can do in this case.
				args = append(args, filterPath)
			}
		}
		cmdObj := self.c.Git().Diff.DiffCmdObj(args, mode)
		prefix := style.FgYellow.Sprintf("%s %s-%s\n\n", self.c.Tr.ShowingDiffForRange, from.ShortRefName(), to.ShortRefName())
		return types.NewMainViewDiffTaskWithPrefix(cmdObj.GetCmd(), types.StaticPrefix(prefix), mode)
	}

	cmdObj := self.c.Git().Commit.ShowCmdObj(commit.Hash(), self.FilterPathsForCommit(commit), mode)
	return types.NewMainViewDiffTaskWithPrefix(cmdObj.GetCmd(), self.commitMessageDiffPrefix(commits, commit), mode)
}

// For an "amend!" commit, returns a prefix with a diff of the commit message it
// sets against the message it replaces, to be shown above the commit's own diff.
// Returns nil for any other commit. The prefix is empty for an "amend!" commit
// that only changes the contents of the commit it applies to.
func (self *DiffHelper) commitMessageDiffPrefix(commits []*models.Commit, commit *models.Commit) types.Prefix {
	previousCommit, ok := findCommitWithPreviousMessage(commits, commit)
	if !ok {
		return nil
	}

	header := style.FgYellow.Sprintf("%s\n", utils.ResolvePlaceholderString(
		self.c.Tr.CommitMessageChanges,
		map[string]string{"hash": previousCommit.ShortHash()},
	))

	return func(width int) func() string {
		produceDiff := self.commitMessageDiff(previousCommit.Hash(), commit.Hash(), width)
		return func() string {
			diff := produceDiff()
			if diff == "" {
				return ""
			}

			return header + diff + strings.Repeat("─", width) + "\n"
		}
	}
}

// The names the two messages are diffed under. A diff renderer shows them as the
// names of the files being diffed, so they are what tells the reader which side
// is which. They are not translated because they end up as file names, and git
// mangles paths outside of ASCII when it states them in a diff.
const (
	oldMessageName = "old message"
	newMessageName = "new message"
)

// commitMessageDiff returns the function that produces the diff of the messages
// of the two commits, laid out to the given width. It is called on the UI
// thread, and the function it returns on the render's own goroutine (see
// types.Prefix).
func (self *DiffHelper) commitMessageDiff(previousHash string, hash string, width int) func() string {
	values := config.DiffRendererValues{
		Width:           width,
		DiffContext:     self.c.UserConfig().Git.DiffContextSize,
		LightBackground: self.c.TerminalHasLightBackground(),
	}

	// The diff renderer lays the diff out, and lays it out according to the
	// values its command is resolved with, so both belong in the key along with
	// the two messages.
	key := fmt.Sprintf("%s\x00%s\x00%+v\x00%s",
		hash,
		previousHash,
		values,
		self.c.State().GetDiffRendererConfigManager().Signature())
	if diff, ok := self.commitMessageDiffs[key]; ok {
		return func() string { return diff }
	}

	differ, err := self.c.Git().Diff.NewTextDiffer(values)
	if err != nil {
		self.c.Log.Error(err)
		return func() string { return "" }
	}
	commitCommands := self.c.Git().Commit

	return func() string {
		diff, err := renderCommitMessageDiff(commitCommands, differ, previousHash, hash)
		if err != nil {
			self.c.Log.Error(err)
			return ""
		}

		self.c.OnUIThread(func() error {
			self.commitMessageDiffs[key] = diff
			return nil
		})
		return diff
	}
}

// renderCommitMessageDiff returns the diff of the messages of the two commits, or
// an empty string if they are the same.
func renderCommitMessageDiff(
	commitCommands *git_commands.CommitCommands, differ *git_commands.TextDiffer, previousHash string, hash string,
) (string, error) {
	messages, err := commitCommands.GetCommitMessages([]string{previousHash, hash})
	if err != nil {
		return "", err
	}

	before := messageAfterAmending(messages[0])
	after := messageAfterAmending(messages[1])
	if before == after {
		return "", nil
	}

	return differ.RenderedDiff(
		git_commands.NamedText{Name: oldMessageName, Content: before},
		git_commands.NamedText{Name: newMessageName, Content: after})
}

// PlainDiffBetweenRefs returns the diff of the given files between two refs as git
// writes it, without colour or a diff renderer's involvement — what a panel showing
// a commit's diff hands out as the diff behind its rendering (see
// types.FocusedMainViewDiffSource). It honours diffing mode, so that the diff is of
// the same two ends the main view is showing.
func (self *DiffHelper) PlainDiffBetweenRefs(from string, to string, paths []string) string {
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)
	// An error means there is no diff to be had, which for our purposes is the same
	// as an empty one.
	diff, _ := self.c.Git().WorkingTree.ShowFileDiffCmdObj(from, to, reverse, paths, git_commands.DiffModePlain).RunWithOutput()
	return diff
}

func (self *DiffHelper) FilterPathsForCommit(commit *models.Commit) []string {
	filterPath := self.c.Modes().Filtering.GetPath()
	if filterPath != "" {
		if len(commit.FilterPaths) > 0 {
			return commit.FilterPaths
		}
		return []string{filterPath}
	}
	return nil
}

func (self *DiffHelper) ExitDiffMode() error {
	self.c.Modes().Diffing = diffing.New()
	self.c.Refresh(types.RefreshOptions{})
	return nil
}

// RenderToMainAgain renders the current side panel into the main view again, if
// that is what the main view shows. This is for when something that the
// rendering depends on has changed, such as the diff renderer.
func (self *DiffHelper) RenderToMainAgain() {
	currentSide := self.c.Context().CurrentSide()
	currentKey := self.c.Context().Current().GetKey()
	if currentSide.GetKey() == currentKey ||
		currentKey == context.NORMAL_MAIN_CONTEXT_KEY ||
		currentKey == context.NORMAL_SECONDARY_CONTEXT_KEY {
		// Whatever changed can make the diff come out differently, such as a new
		// renderer laying it out its own way, so the line you were looking at could
		// end up anywhere in the view; keep it in front of you.
		self.diffLineHelper.PreserveDiffPositionOnRerender(self.c.Contexts().Normal.GetView())
		self.diffLineHelper.PreserveDiffPositionOnRerender(self.c.Contexts().NormalSecondary.GetView())
		currentSide.HandleRenderToMain()
	}
}

func (self *DiffHelper) RenderDiff() {
	args := self.DiffArgs()
	cmdObj := self.c.Git().Diff.DiffCmdObj(args, git_commands.DiffModeRendered)
	prefix := style.FgMagenta.Sprintf(
		"%s %s\n\n",
		self.c.Tr.ShowingGitDiff,
		"git diff "+strings.Join(args, " "),
	)
	task := types.NewMainViewDiffTaskWithPrefix(cmdObj.GetCmd(), types.StaticPrefix(prefix), git_commands.DiffModeRendered)

	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title:    "Diff",
			SubTitle: self.IgnoringWhitespaceSubTitle(),
			Task:     task,
		},
	})
}

// CurrentDiffTerminals returns the current diff terminals of the currently selected item.
// in the case of a branch it returns both the branch and it's upstream name,
// which becomes an option when you bring up the diff menu, but when you're just
// flicking through branches it will be using the local branch name.
func (self *DiffHelper) CurrentDiffTerminals() []string {
	c := self.c.Context().CurrentSide()

	if c.GetKey() == "" {
		return nil
	}

	switch v := c.(type) {
	case types.DiffableContext:
		return v.GetDiffTerminals()
	}

	return nil
}

func (self *DiffHelper) currentDiffTerminal() string {
	names := self.CurrentDiffTerminals()
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func (self *DiffHelper) currentlySelectedFilename() string {
	currentContext := self.c.Context().Current()

	switch currentContext := currentContext.(type) {
	case types.IListContext:
		if lo.Contains([]types.ContextKey{context.FILES_CONTEXT_KEY, context.COMMIT_FILES_CONTEXT_KEY}, currentContext.GetKey()) {
			return currentContext.GetSelectedItemId()
		}
	}

	return ""
}

func (self *DiffHelper) WithDiffModeCheck(f func()) {
	if self.c.Modes().Diffing.Active() {
		self.RenderDiff()
	} else {
		f()
	}
}

func (self *DiffHelper) IgnoringWhitespaceSubTitle() string {
	if self.c.UserConfig().Git.IgnoreWhitespaceInDiffView {
		return self.c.Tr.IgnoreWhitespaceDiffViewSubTitle
	}

	return ""
}

func (self *DiffHelper) OpenDiffToolForRef(selectedRef models.Ref) error {
	to := selectedRef.RefName()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff("")
	_, err := self.c.RunSubprocess(self.c.Git().Diff.OpenDiffToolCmdObj(
		git_commands.DiffToolCmdOptions{
			Filepath:    ".",
			FromCommit:  from,
			ToCommit:    to,
			Reverse:     reverse,
			IsDirectory: true,
			Staged:      false,
		}))
	return err
}

// EditFilesOfDiff opens files of the diff that the current side panel shows in the
// editor. They are opened from the worktree of the diff (see otherWorktreeOfDiff).
// Files that don't exist there are left out, and if none of them exist, it returns an
// error (see existingPathsInWorktreeOfDiff).
func (self *DiffHelper) EditFilesOfDiff(paths []string) error {
	existingPaths, err := self.existingPathsInWorktreeOfDiff(paths)
	if err != nil {
		return err
	}
	return self.filesHelper.EditFiles(existingPaths)
}

// EditFileOfDiffAtLine opens a file of the diff shown in the given view in the editor,
// at a line of that diff. The line is carried forward to where it is in the file now
// (see AdjustLineNumber). Like EditFilesOfDiff, it opens the file from the worktree of
// the diff, and returns an error if the file doesn't exist there.
func (self *DiffHelper) EditFileOfDiffAtLine(path string, lineNumber int, viewName string) error {
	path, err := self.existingPathInWorktreeOfDiff(path)
	if err != nil {
		return err
	}
	lineNumber = self.AdjustLineNumber(path, lineNumber, viewName)
	return self.filesHelper.EditFileAtLine(path, lineNumber)
}

// OpenFileOfDiff opens a file of the diff that the current side panel shows with the
// default application for it. Like EditFilesOfDiff, it opens the file from the
// worktree of the diff, and returns an error if the file doesn't exist there.
func (self *DiffHelper) OpenFileOfDiff(path string) error {
	path, err := self.existingPathInWorktreeOfDiff(path)
	if err != nil {
		return err
	}
	return self.filesHelper.OpenFile(path)
}

// existingPathsInWorktreeOfDiff returns the absolute paths of those of the given files
// of the diff that exist in the worktree of the diff. A file of a diff may be missing
// there, for example if a later commit deleted it, or if the worktree has a branch
// checked out that doesn't have it. An editor would create such a file as a new, empty
// one, so it is left out.
//
// If none of the files exist, it returns an error. If only some of them are missing, it
// leaves them out silently. The user would see such an error only after returning from
// the editor.
func (self *DiffHelper) existingPathsInWorktreeOfDiff(paths []string) ([]string, error) {
	absPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		absPath, err := self.AbsolutePathOfFileOfDiff(path)
		if err != nil {
			return nil, err
		}
		absPaths = append(absPaths, absPath)
	}

	existingPaths := lo.Filter(absPaths, func(path string, _ int) bool {
		exists, err := self.c.OS().FileExists(path)
		// If we can't tell, the editor reports the problem
		return exists || err != nil
	})
	if len(existingPaths) > 0 {
		return existingPaths, nil
	}

	if len(absPaths) == 1 {
		return nil, errors.New(utils.ResolvePlaceholderString(
			self.c.Tr.FileOfDiffDoesNotExist, map[string]string{"path": absPaths[0]}))
	}
	worktreePath, ok := self.otherWorktreeOfDiff()
	if !ok {
		worktreePath = self.c.Git().RepoPaths.WorktreePath()
	}
	return nil, errors.New(utils.ResolvePlaceholderString(
		self.c.Tr.NoFileOfDiffExists, map[string]string{"worktreePath": worktreePath}))
}

// existingPathInWorktreeOfDiff is existingPathsInWorktreeOfDiff for a single file.
func (self *DiffHelper) existingPathInWorktreeOfDiff(path string) (string, error) {
	existingPaths, err := self.existingPathsInWorktreeOfDiff([]string{path})
	if err != nil {
		return "", err
	}
	return existingPaths[0], nil
}

// AbsolutePathOfFileOfDiff returns the absolute path of a file of the diff that the
// current side panel shows, in the worktree that EditFilesOfDiff opens it from.
func (self *DiffHelper) AbsolutePathOfFileOfDiff(path string) (string, error) {
	return filepath.Abs(self.pathInWorktreeOfDiff(path))
}

// otherWorktreeOfDiff returns the path of the worktree to open files of the diff in
// the current side panel from, if that isn't the current worktree. If the diff is of a
// commit of a branch, and another worktree has that branch checked out, that is the
// other worktree. It has the files as the branch has them; the current worktree may
// have them in a different state, or not at all.
func (self *DiffHelper) otherWorktreeOfDiff() (string, bool) {
	branchDiffContext, ok := self.c.Context().CurrentSide().(types.BranchDiffContext)
	if !ok {
		return "", false
	}
	branch := branchDiffContext.BranchOfDiff()
	if branch == nil {
		return "", false
	}
	worktree, ok := git_commands.WorktreeForBranch(branch, self.c.Model().Worktrees)
	if !ok || worktree.IsCurrent || worktree.IsPathMissing {
		return "", false
	}
	return worktree.Path, true
}

// pathInWorktreeOfDiff returns the path to open a file of the diff from. The path
// passed in is either relative to the current worktree, or an absolute path in it, as
// a diff renderer may state paths.
func (self *DiffHelper) pathInWorktreeOfDiff(path string) string {
	worktreePath, ok := self.otherWorktreeOfDiff()
	if !ok {
		return path
	}
	if filepath.IsAbs(path) {
		relativePath := repoRelativePath(self.c.Git().RepoPaths.WorktreePath(), path)
		if relativePath == "" {
			return path
		}
		path = relativePath
	}
	return filepath.Join(worktreePath, path)
}

// AdjustLineNumber is used to adjust a line number in the diff that's currently
// being viewed, so that it corresponds to the line number in the actual working
// copy state of the file. It is used when clicking on a delta hyperlink in a
// diff, or when pressing `e` in a focused diff. It works
// by getting a diff of what's being viewed in the main view against the working
// copy, and then using that diff to adjust the line number. The working copy is that
// of the worktree the file is opened from (see otherWorktreeOfDiff).
// path is the file path of the file being viewed, in that worktree
// linenumber is the line number to adjust (one-based)
// viewname is the name of the view that shows the diff. We need to pass it
// because the diff adjustment is slightly different depending on which view is
// showing the diff.
func (self *DiffHelper) AdjustLineNumber(path string, linenumber int, viewname string) int {
	switch viewname {

	case "main":
		if diffableContext, ok := self.c.Context().CurrentSide().(types.DiffableContext); ok {
			ref := diffableContext.RefForAdjustingLineNumberInDiff()
			if len(ref) != 0 {
				return self.adjustLineNumber(linenumber, ref, "--", path)
			}
		}
		// if the type cast to DiffableContext returns false, we are in the
		// unstaged changes view of the Files panel; no need to adjust line
		// numbers in this case

	case "secondary":
		return self.adjustLineNumber(linenumber, "--", path)
	}

	return linenumber
}

func (self *DiffHelper) adjustLineNumber(linenumber int, diffArgs ...string) int {
	args := append([]string{"--unified=0"}, diffArgs...)
	var diff string
	var err error
	if worktreePath, ok := self.otherWorktreeOfDiff(); ok {
		diff, err = self.c.Git().Diff.GetDiffInOtherWorktree(worktreePath, args...)
	} else {
		diff, err = self.c.Git().Diff.GetDiff(false, args...)
	}
	if err != nil {
		return linenumber
	}
	patch := patch.Parse(diff)
	return patch.AdjustLineNumber(linenumber)
}
