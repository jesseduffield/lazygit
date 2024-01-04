_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Global keybindings

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Switch to a recent repo |  |
| `` <pgup> (fn+up/shift+k) `` | Scroll up main panel |  |
| `` <pgdown> (fn+down/shift+j) `` | Scroll down main panel |  |
| `` @ `` | Open command log menu |  |
| `` } `` | Increase the size of the context shown around changes in the diff view |  |
| `` { `` | Decrease the size of the context shown around changes in the diff view |  |
| `` : `` | Execute custom command |  |
| `` <c-p> `` | View custom patch options |  |
| `` m `` | View merge/rebase options |  |
| `` R `` | Refresh |  |
| `` + `` | Next screen mode (normal/half/fullscreen) |  |
| `` _ `` | Prev screen mode |  |
| `` ? `` | Open menu |  |
| `` <c-s> `` | View filter-by-path options |  |
| `` W `` | Open diff menu |  |
| `` <c-e> `` | Open diff menu |  |
| `` <c-w> `` | Toggle whether or not whitespace changes are shown in the diff view |  |
| `` z `` | Undo | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | Redo | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` P `` | Push |  |
| `` p `` | Pull |  |

## List panel navigation

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Previous page |  |
| `` . `` | Next page |  |
| `` < `` | Scroll to top |  |
| `` > `` | Scroll to bottom |  |
| `` v `` | Toggle range select |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | Search the current view by text |  |
| `` H `` | Scroll left |  |
| `` L `` | Scroll right |  |
| `` ] `` | Next tab |  |
| `` [ `` | Previous tab |  |

## Commit files

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy the committed file name to the clipboard |  |
| `` c `` | Checkout file |  |
| `` d `` | Discard this commit's changes to this file |  |
| `` o `` | Open file |  |
| `` e `` | Edit file |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle file included in patch |  |
| `` a `` | Toggle all files included in patch |  |
| `` <enter> `` | Enter file to add selected lines to the patch (or toggle directory collapsed) |  |
| `` ` `` | Toggle file tree view |  |
| `` / `` | Search the current view by text |  |

## Commit summary

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Confirm |  |
| `` <esc> `` | Close |  |

## Commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit SHA to clipboard |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` b `` | View bisect options |  |
| `` s `` | Squash down |  |
| `` f `` | Fixup commit |  |
| `` r `` | Reword commit |  |
| `` R `` | Reword commit with editor |  |
| `` d `` | Delete commit |  |
| `` e `` | Edit commit |  |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick commit (when mid-rebase) |  |
| `` F `` | Create fixup commit for this commit |  |
| `` S `` | Squash all 'fixup!' commits above selected commit (autosquash) |  |
| `` <c-j> `` | Move commit down one |  |
| `` <c-k> `` | Move commit up one |  |
| `` V `` | Paste commits (cherry-pick) |  |
| `` B `` | Mark commit as base commit for rebase | Select a base commit for the next rebase; this will effectively perform a 'git rebase --onto'. |
| `` A `` | Amend commit with staged changes |  |
| `` a `` | Set/Reset commit author |  |
| `` t `` | Revert commit |  |
| `` T `` | Tag commit |  |
| `` <c-l> `` | Open log menu |  |
| `` w `` | View worktree options |  |
| `` <space> `` | Checkout commit |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | View reset options |  |
| `` C `` | Copy commit (cherry-pick) |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View selected item's files |  |
| `` / `` | Search the current view by text |  |

## Confirmation panel

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Confirm |  |
| `` <esc> `` | Close/Cancel |  |

## Files

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy the file name to the clipboard |  |
| `` <space> `` | Toggle staged |  |
| `` <c-b> `` | Filter files by status |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Commit changes |  |
| `` w `` | Commit changes without pre-commit hook |  |
| `` A `` | Amend last commit |  |
| `` C `` | Commit changes using git editor |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit file |  |
| `` o `` | Open file |  |
| `` i `` | Ignore or exclude file |  |
| `` r `` | Refresh files |  |
| `` s `` | Stash all changes |  |
| `` S `` | View stash options |  |
| `` a `` | Stage/unstage all |  |
| `` <enter> `` | Stage individual hunks/lines for file, or collapse/expand for directory |  |
| `` d `` | View 'discard changes' options |  |
| `` g `` | View upstream reset options |  |
| `` D `` | View reset options |  |
| `` ` `` | Toggle file tree view |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Open external merge tool (git mergetool) |  |
| `` f `` | Fetch |  |
| `` / `` | Search the current view by text |  |

## Local branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` i `` | Show git-flow options |  |
| `` <space> `` | Checkout |  |
| `` n `` | New branch |  |
| `` o `` | Create pull request |  |
| `` O `` | Create pull request options |  |
| `` <c-y> `` | Copy pull request URL to clipboard |  |
| `` c `` | Checkout by name, enter '-' to switch to last |  |
| `` F `` | Force checkout |  |
| `` d `` | View delete options |  |
| `` r `` | Rebase checked-out branch onto this branch |  |
| `` M `` | Merge into currently checked out branch |  |
| `` f `` | Fast-forward this branch from its upstream |  |
| `` T `` | Create tag |  |
| `` s `` | Sort order |  |
| `` g `` | View reset options |  |
| `` R `` | Rename branch |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream |
| `` w `` | View worktree options |  |
| `` <enter> `` | View commits |  |
| `` / `` | Filter the current view by text |  |

## Main panel (merging)

| Key | Action | Info |
|-----|--------|-------------|
| `` e `` | Edit file |  |
| `` o `` | Open file |  |
| `` <left> `` | Select previous conflict |  |
| `` <right> `` | Select next conflict |  |
| `` <up> `` | Select previous hunk |  |
| `` <down> `` | Select next hunk |  |
| `` z `` | Undo |  |
| `` M `` | Open external merge tool (git mergetool) |  |
| `` <space> `` | Pick hunk |  |
| `` b `` | Pick all hunks |  |
| `` <esc> `` | Return to files panel |  |

## Main panel (normal)

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | Scroll down |  |
| `` mouse wheel up (fn+down) `` | Scroll up |  |

## Main panel (patch building)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Select previous hunk |  |
| `` <right> `` | Select next hunk |  |
| `` v `` | Toggle range select |  |
| `` a `` | Toggle select hunk |  |
| `` <c-o> `` | Copy the selected text to the clipboard |  |
| `` o `` | Open file |  |
| `` e `` | Edit file |  |
| `` <space> `` | Add/Remove line(s) to patch |  |
| `` <esc> `` | Exit custom patch builder |  |
| `` / `` | Search the current view by text |  |

## Main panel (staging)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Select previous hunk |  |
| `` <right> `` | Select next hunk |  |
| `` v `` | Toggle range select |  |
| `` a `` | Toggle select hunk |  |
| `` <c-o> `` | Copy the selected text to the clipboard |  |
| `` o `` | Open file |  |
| `` e `` | Edit file |  |
| `` <esc> `` | Return to files panel |  |
| `` <tab> `` | Switch to other panel (staged/unstaged changes) |  |
| `` <space> `` | Toggle line staged / unstaged |  |
| `` d `` | Discard change (git reset) |  |
| `` E `` | Edit hunk |  |
| `` c `` | Commit changes |  |
| `` w `` | Commit changes without pre-commit hook |  |
| `` C `` | Commit changes using git editor |  |
| `` / `` | Search the current view by text |  |

## Menu

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Execute |  |
| `` <esc> `` | Close |  |
| `` / `` | Filter the current view by text |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit SHA to clipboard |  |
| `` w `` | View worktree options |  |
| `` <space> `` | Checkout commit |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | View reset options |  |
| `` C `` | Copy commit (cherry-pick) |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` / `` | Filter the current view by text |  |

## Remote branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` <space> `` | Checkout |  |
| `` n `` | New branch |  |
| `` M `` | Merge into currently checked out branch |  |
| `` r `` | Rebase checked-out branch onto this branch |  |
| `` d `` | Delete remote tag |  |
| `` u `` | Set as upstream of checked-out branch |  |
| `` s `` | Sort order |  |
| `` g `` | View reset options |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | View commits |  |
| `` / `` | Filter the current view by text |  |

## Remotes

| Key | Action | Info |
|-----|--------|-------------|
| `` f `` | Fetch remote |  |
| `` n `` | Add new remote |  |
| `` d `` | Remove remote |  |
| `` e `` | Edit remote |  |
| `` / `` | Filter the current view by text |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Apply |  |
| `` g `` | Pop |  |
| `` d `` | Drop |  |
| `` n `` | New branch |  |
| `` r `` | Rename stash |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | View selected item's files |  |
| `` / `` | Filter the current view by text |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Open config file |  |
| `` e `` | Edit config file |  |
| `` u `` | Check for update |  |
| `` <enter> `` | Switch to a recent repo |  |
| `` a `` | Show all branch logs |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit SHA to clipboard |  |
| `` w `` | View worktree options |  |
| `` <space> `` | Checkout commit |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | View reset options |  |
| `` C `` | Copy commit (cherry-pick) |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View selected item's files |  |
| `` / `` | Search the current view by text |  |

## Submodules

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy submodule name to clipboard |  |
| `` <enter> `` | Enter submodule |  |
| `` <space> `` | Enter submodule |  |
| `` d `` | Remove submodule |  |
| `` u `` | Update submodule |  |
| `` n `` | Add new submodule |  |
| `` e `` | Update submodule URL |  |
| `` i `` | Initialize submodule |  |
| `` b `` | View bulk submodule options |  |
| `` / `` | Filter the current view by text |  |

## Tags

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Checkout |  |
| `` d `` | View delete options |  |
| `` P `` | Push tag |  |
| `` n `` | Create tag |  |
| `` g `` | View reset options |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | View commits |  |
| `` / `` | Filter the current view by text |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | Create worktree |  |
| `` <space> `` | Switch to worktree |  |
| `` <enter> `` | Switch to worktree |  |
| `` o `` | Open in editor |  |
| `` d `` | Remove worktree |  |
| `` / `` | Filter the current view by text |  |
