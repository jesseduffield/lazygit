package controllers

import (
	"fmt"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// togglePatchFromFocusedMainView toggles the selected diff line(s) — a single line, a
// range, or a hunk — into or out of the custom patch being built for (from, to,
// reverse). It is the patch-building counterpart of the files panel's staging handler,
// shared by the panels that build a patch from a focused main view: the commit files
// panel (the per-file diff) and the commits / sub-commits / stash panels (the
// whole-commit diff). The selection is resolved to change-line identities (shared with
// staging) and mapped to the patch builder's per-file line indices.
//
// Building a patch needs the patch builder started for this target; if a patch for a
// different target is active, we confirm before discarding it (as entering the patch
// builder does). The diff itself is unchanged by a toggle, so refresh — supplied by the
// caller, since what needs re-rendering differs per panel — re-renders the same diff
// command (keeping its scroll) and the secondary patch view, and the inclusion gutter
// rides that re-render. That re-render (and its layout re-wrap when the secondary view
// first appears) is in view-line space, which moves the selection, so we re-establish it
// afterwards by its width-independent change-line ordinal.
func togglePatchFromFocusedMainView(
	c *ControllerCommon,
	mainViewName string,
	firstLineIdx int,
	lastLineIdx int,
	from string,
	to string,
	reverse bool,
	canRebase bool,
	refresh func(),
) error {
	infos := c.Helpers().Staging.ChangeLinesInViewRange(mainViewName, firstLineIdx, lastLineIdx)
	if len(infos) == 0 {
		return nil
	}

	patchBuilder := c.Git().Patch.PatchBuilder
	mustDiscardPatch := patchBuilder.Active() && patchBuilder.NewPatchRequired(from, to, reverse)
	return c.ConfirmIf(mustDiscardPatch, types.ConfirmOpts{
		Title:  c.Tr.DiscardPatch,
		Prompt: c.Tr.DiscardPatchConfirm,
		HandleConfirm: func() error {
			if mustDiscardPatch {
				patchBuilder.Reset()
			}
			if !patchBuilder.Active() {
				patchBuilder.Start(from, to, reverse, canRebase)
			}

			if err := togglePatchLines(c, infos); err != nil {
				return err
			}

			refresh()
			// A toggle doesn't change the diff, so source and target are the same pane;
			// the re-render (and the layout re-wrap when the secondary view first appears)
			// still moves the selection in view-line space, so re-establish it.
			revealSelectionAfterPrimaryAction(c, mainViewName, mainViewName, firstLineIdx)
			return nil
		},
	})
}

// togglePatchLines toggles the change lines identified by infos into or out of the
// custom patch. The direction is decided once, from the first selected change line (in
// view order) — already in the patch means remove the whole selection, otherwise add it
// — and applied to every file the selection spans, mirroring how the patch explorer
// toggles a range.
func togglePatchLines(c *ControllerCommon, infos []types.DiffLineInfo) error {
	patchBuilder := c.Git().Patch.PatchBuilder

	// Group the selected change lines by file (in view order), then resolve each file's
	// identities to patch-line indices in one pass.
	var order []string
	identitiesByFile := map[string][]patch.LineIdentity{}
	for _, info := range infos {
		filename := patchFilename(c, info.Path)
		if filename == "" {
			continue
		}
		if _, seen := identitiesByFile[filename]; !seen {
			order = append(order, filename)
		}
		lineNumber, isDeletion := info.PatchSelectLine()
		identitiesByFile[filename] = append(identitiesByFile[filename],
			patch.LineIdentity{LineNumber: lineNumber, IsDeletion: isDeletion})
	}
	if len(order) == 0 {
		return nil
	}

	indicesByFile := map[string][]int{}
	for filename, identities := range identitiesByFile {
		indices, err := patchBuilder.PatchLineIndicesForLines(filename, identities)
		if err != nil {
			return err
		}
		indicesByFile[filename] = indices
	}

	// Decide the direction from the first selected change line.
	firstFile := order[0]
	included, err := patchBuilder.GetFileIncLineIndices(firstFile, "")
	if err != nil {
		return err
	}
	remove := len(indicesByFile[firstFile]) > 0 && lo.Contains(included, indicesByFile[firstFile][0])
	toggle := patchBuilder.AddFileLineRange
	if remove {
		toggle = patchBuilder.RemoveFileLineRange
	}

	for _, filename := range order {
		indices := indicesByFile[filename]
		if len(indices) == 0 {
			continue
		}
		if err := toggle(filename, "", indices); err != nil {
			return err
		}
	}
	return nil
}

// discardFromCommitDisabledReason reports why discarding lines from the commit shown in
// the focused main view is unavailable, or nil when it's available. It mirrors the patch
// builder's own guard: the diff must be a local commit's (canRebase), no rebase may be in
// progress, and the diff context size must be non-zero (a patch can't be built from a
// zero-context diff). Stash and sub-commits of another branch are never rebaseable, so
// they always get the local-commits reason.
func discardFromCommitDisabledReason(c *ControllerCommon, canRebase bool) *types.DisabledReason {
	if !canRebase {
		return &types.DisabledReason{Text: c.Tr.CanOnlyDiscardFromLocalCommits}
	}
	if c.Git().Status.WorkingTreeState().Any() {
		return &types.DisabledReason{Text: c.Tr.CantPatchWhileRebasingError}
	}
	if c.UserConfig().Git.DiffContextSize == 0 {
		return &types.DisabledReason{Text: fmt.Sprintf(c.Tr.Actions.NotEnoughContextToRemoveLines,
			c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView)}
	}
	return nil
}

// discardSelectionFromCommit discards the selected diff line(s) — a single line, a
// range, or a hunk — from the commit whose diff the focused main view shows, by building
// a one-off patch from the selection and removing it from the commit via a rebase. It is
// the patch-building counterpart of the working-tree discard the files panel does, and
// mirrors the patch builder's own "discard lines from commit": any in-progress custom
// patch is reset first (the confirm prompt warns when one exists), then a fresh patch
// holding exactly the selection is built for (from, to, reverse) and removed from the
// commit. The caller has already established (via canRebase) that the commit is on a
// local branch; the target commit is the one identified by `to`.
func discardSelectionFromCommit(
	c *ControllerCommon,
	mainViewName string,
	firstLineIdx int,
	lastLineIdx int,
	from string,
	to string,
	reverse bool,
	canRebase bool,
) error {
	infos := c.Helpers().Staging.ChangeLinesInViewRange(mainViewName, firstLineIdx, lastLineIdx)
	if len(infos) == 0 {
		return nil
	}

	commitIndex := -1
	for i, commit := range c.Model().Commits {
		if commit.Hash() == to {
			commitIndex = i
			break
		}
	}
	if commitIndex == -1 {
		return nil
	}

	patchBuilder := c.Git().Patch.PatchBuilder
	prompt := lo.Ternary(patchBuilder.Active(),
		c.Tr.DiscardLinesFromCommitPromptWithReset,
		c.Tr.DiscardLinesFromCommitPrompt)

	c.Confirm(types.ConfirmOpts{
		Title:  c.Tr.DiscardLinesFromCommitTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			// Build a fresh patch holding exactly the selection: reset any active patch,
			// then add the selected lines (togglePatchLines adds them, the patch being empty).
			if patchBuilder.Active() {
				patchBuilder.Reset()
			}
			patchBuilder.Start(from, to, reverse, canRebase)
			if err := togglePatchLines(c, infos); err != nil {
				return err
			}
			if patchBuilder.IsEmpty() {
				return nil
			}

			return c.WithWaitingStatusBlockingInput(c.Tr.RebasingStatus, func(gocui.Task) error {
				c.LogAction(c.Tr.Actions.RemovePatchFromCommit)
				err := c.Git().Patch.DeletePatchesFromCommit(c.Model().Commits, commitIndex)
				// The rebase rewrites the commit, so the focused main view's selection is
				// re-established on the next surviving change as the diff re-renders (see
				// preserveFocusedMainViewSelectionAcrossContentChange).
				return c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
			})
		},
	})
	return nil
}

// patchFilename maps a diff line's absolute path to the key the patch builder stores
// the file under — its repo-relative, slash-separated path.
func patchFilename(c *ControllerCommon, absPath string) string {
	relPath, err := filepath.Rel(c.Git().RepoPaths.WorktreePath(), absPath)
	if err != nil {
		return ""
	}
	return filepath.ToSlash(relPath)
}
