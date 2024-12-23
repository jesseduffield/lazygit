package helpers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type DiffHelper struct {
	c *HelperCommon
}

func NewDiffHelper(c *HelperCommon) *DiffHelper {
	return &DiffHelper{
		c: c,
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
func (self *DiffHelper) GetUpdateTaskForRenderingCommitsDiff(commit *models.Commit, refRange *types.RefRange) types.UpdateTask {
	if refRange != nil {
		from, to := refRange.From, refRange.To
		args := []string{from.ParentRefName(), to.RefName(), "--stat", "-p"}
		args = append(args, "--")
		if path := self.c.Modes().Filtering.GetPath(); path != "" {
			args = append(args, path)
		}
		cmdObj := self.c.Git().Diff.DiffCmdObj(args)
		task := types.NewRunPtyTask(cmdObj.GetCmd())
		task.Prefix = style.FgYellow.Sprintf("%s %s-%s\n\n", self.c.Tr.ShowingDiffForRange, from.ShortRefName(), to.ShortRefName())
		return task
	}

	cmdObj := self.c.Git().Commit.ShowCmdObj(commit.Hash, self.c.Modes().Filtering.GetPath())
	return types.NewRunPtyTask(cmdObj.GetCmd())
}

func (self *DiffHelper) ExitDiffMode() error {
	self.c.Modes().Diffing = diffing.New()
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *DiffHelper) RenderDiff() {
	args := self.DiffArgs()
	cmdObj := self.c.Git().Diff.DiffCmdObj(args)
	task := types.NewRunPtyTask(cmdObj.GetCmd())
	task.Prefix = style.FgMagenta.Sprintf(
		"%s %s\n\n",
		self.c.Tr.ShowingGitDiff,
		"git diff "+strings.Join(args, " "),
	)

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
	if self.c.GetAppState().IgnoreWhitespaceInDiffView {
		return self.c.Tr.IgnoreWhitespaceDiffViewSubTitle
	}

	return ""
}

func (self *DiffHelper) OpenDiffToolForRef(selectedRef types.Ref) error {
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

// AdjustLineNumber is used to adjust a line number in the diff that's currently
// being viewed, so that it corresponds to the line number in the actual working
// copy state of the file. It is used when clicking on a delta hyperlink in a
// diff, or when pressing `e` in the staging or patch building panels. It works
// by getting a diff of what's being viewed in the main view against the working
// copy, and then using that diff to adjust the line number.
// path is the file path of the file being viewed
// linenumber is the line number to adjust (one-based)
// viewname is the name of the view that shows the diff. We need to pass it
// because the diff adjustment is slightly different depending on which view is
// showing the diff.
func (self *DiffHelper) AdjustLineNumber(path string, linenumber int, viewname string) int {
	switch viewname {

	case "main", "patchBuilding":
		if diffableContext, ok := self.c.Context().CurrentSide().(types.DiffableContext); ok {
			ref := diffableContext.RefForAdjustingLineNumberInDiff()
			if len(ref) != 0 {
				return self.adjustLineNumber(linenumber, ref, "--", path)
			}
		}
		// if the type cast to DiffableContext returns false, we are in the
		// unstaged changes view of the Files panel; no need to adjust line
		// numbers in this case

	case "secondary", "stagingSecondary":
		return self.adjustLineNumber(linenumber, "--", path)
	}

	return linenumber
}

func (self *DiffHelper) adjustLineNumber(linenumber int, diffArgs ...string) int {
	args := append([]string{"--unified=0"}, diffArgs...)
	diff, err := self.c.Git().Diff.GetDiff(false, args...)
	if err != nil {
		return linenumber
	}
	patch := patch.Parse(diff)
	return patch.AdjustLineNumber(linenumber)
}
