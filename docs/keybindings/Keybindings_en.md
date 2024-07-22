_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Global keybindings

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Switch to a recent repo |  |
| `` <pgup> (fn+up/shift+k) `` | Scroll up main window |  |
| `` <pgdown> (fn+down/shift+j) `` | Scroll down main window |  |
| `` @ `` | View command log options | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | Push | Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` p `` | Pull | Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` } `` | Increase diff context size | Increase the amount of the context shown around changes in the diff view. |
| `` { `` | Decrease diff context size | Decrease the amount of the context shown around changes in the diff view. |
| `` : `` | Execute custom command | Bring up a prompt where you can enter a shell command to execute. Not to be confused with pre-configured custom commands. |
| `` <c-p> `` | View custom patch options |  |
| `` m `` | View merge/rebase options | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | Refresh | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | Next screen mode (normal/half/fullscreen) |  |
| `` _ `` | Prev screen mode |  |
| `` ? `` | Open keybindings menu |  |
| `` <c-s> `` | View filter options | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | View diffing options | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | View diffing options | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | Quit |  |
| `` <esc> `` | Cancel |  |
| `` <c-w> `` | Toggle whitespace | Toggle whether or not whitespace changes are shown in the diff view. |
| `` z `` | Undo | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | Redo | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |

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
| `` <c-o> `` | Copy path to clipboard |  |
| `` c `` | Checkout | Checkout file. This replaces the file in your working tree with the version from the selected commit. |
| `` d `` | Remove | Discard this commit's changes to this file. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes this file. |
| `` o `` | Open file | Open file in default application. |
| `` e `` | Edit | Open file in external editor. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle file included in patch | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Toggle all files | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Enter file / Toggle directory collapsed | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | Toggle file tree view | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` / `` | Search the current view by text |  |

## Commit summary

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Confirm |  |
| `` <esc> `` | Close |  |

## Commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit hash to clipboard |  |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` b `` | View bisect options |  |
| `` s `` | Squash | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | Fixup | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` r `` | Reword | Reword the selected commit's message. |
| `` R `` | Reword with editor |  |
| `` d `` | Drop | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | Edit (start interactive rebase) | Edit the selected commit. Use this to start an interactive rebase from the selected commit. When already mid-rebase, this will mark the selected commit for editing, which means that upon continuing the rebase, the rebase will pause at the selected commit to allow you to make changes. |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick | Mark the selected commit to be picked (when mid-rebase). This means that the commit will be retained upon continuing the rebase. |
| `` F `` | Create fixup commit | Create 'fixup!' commit for the selected commit. Later on, you can press `S` on this same commit to apply all above fixup commits. |
| `` S `` | Apply fixup commits | Squash all 'fixup!' commits, either above the selected commit, or all in current branch (autosquash). |
| `` <c-j> `` | Move commit down one |  |
| `` <c-k> `` | Move commit up one |  |
| `` V `` | Paste (cherry-pick) |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Amend | Amend commit with staged changes. If the selected commit is the HEAD commit, this will perform `git commit --amend`. Otherwise the commit will be amended via a rebase. |
| `` a `` | Amend commit attribute | Set/Reset commit author or set co-author. |
| `` t `` | Revert | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | Tag commit | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <c-l> `` | View log options | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` <space> `` | Checkout | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Copy (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View files |  |
| `` w `` | View worktree options |  |
| `` / `` | Search the current view by text |  |

## Confirmation panel

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Confirm |  |
| `` <esc> `` | Close/Cancel |  |

## Files

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy path to clipboard |  |
| `` <space> `` | Stage | Toggle staged for selected file. |
| `` <c-b> `` | Filter files by status |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Commit | Commit staged changes. |
| `` w `` | Commit changes without pre-commit hook |  |
| `` A `` | Amend last commit |  |
| `` C `` | Commit changes using git editor |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit | Open file in external editor. |
| `` o `` | Open file | Open file in default application. |
| `` i `` | Ignore or exclude file |  |
| `` r `` | Refresh files |  |
| `` s `` | Stash | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | View stash options | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | Stage all | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | Stage lines / Collapse directory | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | Discard | View options for discarding changes to the selected file. |
| `` g `` | View upstream reset options |  |
| `` D `` | Reset | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | Toggle file tree view | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Open external merge tool | Run `git mergetool`. |
| `` f `` | Fetch | Fetch changes from remote. |
| `` / `` | Search the current view by text |  |

## Local branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` i `` | Show git-flow options |  |
| `` <space> `` | Checkout | Checkout selected item. |
| `` n `` | New branch |  |
| `` o `` | Create pull request |  |
| `` O `` | View create pull request options |  |
| `` <c-y> `` | Copy pull request URL to clipboard |  |
| `` c `` | Checkout by name | Checkout by name. In the input box you can enter '-' to switch to the last branch. |
| `` F `` | Force checkout | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | Delete | View delete options for local/remote branch. |
| `` r `` | Rebase | Rebase the checked-out branch onto the selected branch. |
| `` M `` | Merge | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` f `` | Fast-forward | Fast-forward selected branch from its upstream. |
| `` T `` | New tag |  |
| `` s `` | Sort order |  |
| `` g `` | Reset |  |
| `` R `` | Rename branch |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Main panel (merging)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Pick hunk |  |
| `` b `` | Pick all hunks |  |
| `` <up> `` | Previous hunk |  |
| `` <down> `` | Next hunk |  |
| `` <left> `` | Previous conflict |  |
| `` <right> `` | Next conflict |  |
| `` z `` | Undo | Undo last merge conflict resolution. |
| `` e `` | Edit file | Open file in external editor. |
| `` o `` | Open file | Open file in default application. |
| `` M `` | Open external merge tool | Run `git mergetool`. |
| `` <esc> `` | Return to files panel |  |

## Main panel (normal)

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | Scroll down |  |
| `` mouse wheel up (fn+down) `` | Scroll up |  |

## Main panel (patch building)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Go to previous hunk |  |
| `` <right> `` | Go to next hunk |  |
| `` v `` | Toggle range select |  |
| `` a `` | Select hunk | Toggle hunk selection mode. |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` o `` | Open file | Open file in default application. |
| `` e `` | Edit file | Open file in external editor. |
| `` <space> `` | Toggle lines in patch |  |
| `` <esc> `` | Exit custom patch builder |  |
| `` / `` | Search the current view by text |  |

## Main panel (staging)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Go to previous hunk |  |
| `` <right> `` | Go to next hunk |  |
| `` v `` | Toggle range select |  |
| `` a `` | Select hunk | Toggle hunk selection mode. |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` <space> `` | Stage | Toggle selection staged / unstaged. |
| `` d `` | Discard | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | Open file | Open file in default application. |
| `` e `` | Edit file | Open file in external editor. |
| `` <esc> `` | Return to files panel |  |
| `` <tab> `` | Switch view | Switch to other view (staged/unstaged changes). |
| `` E `` | Edit hunk | Edit selected hunk in external editor. |
| `` c `` | Commit | Commit staged changes. |
| `` w `` | Commit changes without pre-commit hook |  |
| `` C `` | Commit changes using git editor |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
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
| `` <c-o> `` | Copy commit hash to clipboard |  |
| `` <space> `` | Checkout | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Copy (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Remote branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` <space> `` | Checkout | Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head. |
| `` n `` | New branch |  |
| `` M `` | Merge | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` r `` | Rebase | Rebase the checked-out branch onto the selected branch. |
| `` d `` | Delete | Delete the remote branch from the remote. |
| `` u `` | Set as upstream | Set the selected remote branch as the upstream of the checked-out branch. |
| `` s `` | Sort order |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Remotes

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | New remote |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Edit | Edit the selected remote's name or URL. |
| `` f `` | Fetch | Fetch updates from the remote repository. This retrieves new commits and branches without merging them into your local branches. |
| `` / `` | Filter the current view by text |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Apply | Apply the stash entry to your working directory. |
| `` g `` | Pop | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | Drop | Remove the stash entry from the stash list. |
| `` n `` | New branch | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | Rename stash |  |
| `` <enter> `` | View files |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Open config file | Open file in default application. |
| `` e `` | Edit config file | Open file in external editor. |
| `` u `` | Check for update |  |
| `` <enter> `` | Switch to a recent repo |  |
| `` a `` | Show/cycle all branch logs |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit hash to clipboard |  |
| `` <space> `` | Checkout | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Copy (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View files |  |
| `` w `` | View worktree options |  |
| `` / `` | Search the current view by text |  |

## Submodules

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy submodule name to clipboard |  |
| `` <enter> `` | Enter | Enter submodule. After entering the submodule, you can press `<esc>` to escape back to the parent repo. |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | Update selected submodule. |
| `` n `` | New submodule |  |
| `` e `` | Update submodule URL |  |
| `` i `` | Initialize | Initialize the selected submodule to prepare for fetching. You probably want to follow this up by invoking the 'update' action to fetch the submodule. |
| `` b `` | View bulk submodule options |  |
| `` / `` | Filter the current view by text |  |

## Tags

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Checkout | Checkout the selected tag tag as a detached HEAD. |
| `` n `` | New tag | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Delete | View delete options for local/remote tag. |
| `` P `` | Push tag | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | Open in editor |  |
| `` d `` | Remove | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | Filter the current view by text |  |
