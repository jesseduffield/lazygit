_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Global keybindings

<pre>
  <kbd>&lt;c-r&gt;</kbd>: Switch to a recent repo
  <kbd>&lt;pgup&gt;</kbd>: Scroll up main panel (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: Scroll down main panel (fn+down/shift+j)
  <kbd>@</kbd>: Open command log menu
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: Execute custom command
  <kbd>&lt;c-p&gt;</kbd>: View custom patch options
  <kbd>m</kbd>: View merge/rebase options
  <kbd>R</kbd>: Refresh
  <kbd>+</kbd>: Next screen mode (normal/half/fullscreen)
  <kbd>_</kbd>: Prev screen mode
  <kbd>?</kbd>: Open menu
  <kbd>&lt;c-s&gt;</kbd>: View filter-by-path options
  <kbd>W</kbd>: Open diff menu
  <kbd>&lt;c-e&gt;</kbd>: Open diff menu
  <kbd>&lt;c-w&gt;</kbd>: Toggle whether or not whitespace changes are shown in the diff view
  <kbd>z</kbd>: Undo
  <kbd>&lt;c-z&gt;</kbd>: Redo
  <kbd>P</kbd>: Push
  <kbd>p</kbd>: Pull
</pre>

## List panel navigation

<pre>
  <kbd>,</kbd>: Previous page
  <kbd>.</kbd>: Next page
  <kbd>&lt;</kbd>: Scroll to top
  <kbd>&gt;</kbd>: Scroll to bottom
  <kbd>/</kbd>: Search the current view by text
  <kbd>H</kbd>: Scroll left
  <kbd>L</kbd>: Scroll right
  <kbd>]</kbd>: Next tab
  <kbd>[</kbd>: Previous tab
</pre>

## Commit files

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy the committed file name to the clipboard
  <kbd>c</kbd>: Checkout file
  <kbd>d</kbd>: Discard this commit's changes to this file
  <kbd>o</kbd>: Open file
  <kbd>e</kbd>: Edit file
  <kbd>&lt;space&gt;</kbd>: Toggle file included in patch
  <kbd>a</kbd>: Toggle all files included in patch
  <kbd>&lt;enter&gt;</kbd>: Enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: Toggle file tree view
  <kbd>/</kbd>: Search the current view by text
</pre>

## Commit summary

<pre>
  <kbd>&lt;enter&gt;</kbd>: Confirm
  <kbd>&lt;esc&gt;</kbd>: Close
</pre>

## Commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy commit SHA to clipboard
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>b</kbd>: View bisect options
  <kbd>s</kbd>: Squash down
  <kbd>f</kbd>: Fixup commit
  <kbd>r</kbd>: Reword commit
  <kbd>R</kbd>: Reword commit with editor
  <kbd>d</kbd>: Delete commit
  <kbd>e</kbd>: Edit commit
  <kbd>p</kbd>: Pick commit (when mid-rebase)
  <kbd>F</kbd>: Create fixup commit for this commit
  <kbd>S</kbd>: Squash all 'fixup!' commits above selected commit (autosquash)
  <kbd>&lt;c-j&gt;</kbd>: Move commit down one
  <kbd>&lt;c-k&gt;</kbd>: Move commit up one
  <kbd>v</kbd>: Paste commits (cherry-pick)
  <kbd>B</kbd>: Mark commit as base commit for rebase
  <kbd>A</kbd>: Amend commit with staged changes
  <kbd>a</kbd>: Set/Reset commit author
  <kbd>t</kbd>: Revert commit
  <kbd>T</kbd>: Tag commit
  <kbd>&lt;c-l&gt;</kbd>: Open log menu
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Checkout commit
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: Open commit in browser
  <kbd>n</kbd>: Create new branch off of commit
  <kbd>g</kbd>: View reset options
  <kbd>c</kbd>: Copy commit (cherry-pick)
  <kbd>C</kbd>: Copy commit range (cherry-pick)
  <kbd>&lt;enter&gt;</kbd>: View selected item's files
  <kbd>/</kbd>: Search the current view by text
</pre>

## Confirmation panel

<pre>
  <kbd>&lt;enter&gt;</kbd>: Confirm
  <kbd>&lt;esc&gt;</kbd>: Close/Cancel
</pre>

## Files

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy the file name to the clipboard
  <kbd>d</kbd>: View 'discard changes' options
  <kbd>&lt;space&gt;</kbd>: Toggle staged
  <kbd>&lt;c-b&gt;</kbd>: Filter files by status
  <kbd>c</kbd>: Commit changes
  <kbd>w</kbd>: Commit changes without pre-commit hook
  <kbd>A</kbd>: Amend last commit
  <kbd>C</kbd>: Commit changes using git editor
  <kbd>e</kbd>: Edit file
  <kbd>o</kbd>: Open file
  <kbd>i</kbd>: Ignore or exclude file
  <kbd>r</kbd>: Refresh files
  <kbd>s</kbd>: Stash all changes
  <kbd>S</kbd>: View stash options
  <kbd>a</kbd>: Stage/unstage all
  <kbd>&lt;enter&gt;</kbd>: Stage individual hunks/lines for file, or collapse/expand for directory
  <kbd>g</kbd>: View upstream reset options
  <kbd>D</kbd>: View reset options
  <kbd>`</kbd>: Toggle file tree view
  <kbd>M</kbd>: Open external merge tool (git mergetool)
  <kbd>f</kbd>: Fetch
  <kbd>/</kbd>: Search the current view by text
</pre>

## Local branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy branch name to clipboard
  <kbd>i</kbd>: Show git-flow options
  <kbd>&lt;space&gt;</kbd>: Checkout
  <kbd>n</kbd>: New branch
  <kbd>o</kbd>: Create pull request
  <kbd>O</kbd>: Create pull request options
  <kbd>&lt;c-y&gt;</kbd>: Copy pull request URL to clipboard
  <kbd>c</kbd>: Checkout by name
  <kbd>F</kbd>: Force checkout
  <kbd>d</kbd>: Delete branch
  <kbd>r</kbd>: Rebase checked-out branch onto this branch
  <kbd>M</kbd>: Merge into currently checked out branch
  <kbd>f</kbd>: Fast-forward this branch from its upstream
  <kbd>T</kbd>: Create tag
  <kbd>g</kbd>: View reset options
  <kbd>R</kbd>: Rename branch
  <kbd>u</kbd>: Set/Unset upstream
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: View commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Main panel (merging)

<pre>
  <kbd>e</kbd>: Edit file
  <kbd>o</kbd>: Open file
  <kbd>&lt;left&gt;</kbd>: Select previous conflict
  <kbd>&lt;right&gt;</kbd>: Select next conflict
  <kbd>&lt;up&gt;</kbd>: Select previous hunk
  <kbd>&lt;down&gt;</kbd>: Select next hunk
  <kbd>z</kbd>: Undo
  <kbd>M</kbd>: Open external merge tool (git mergetool)
  <kbd>&lt;space&gt;</kbd>: Pick hunk
  <kbd>b</kbd>: Pick all hunks
  <kbd>&lt;esc&gt;</kbd>: Return to files panel
</pre>

## Main panel (normal)

<pre>
  <kbd>mouse wheel down</kbd>: Scroll down (fn+up)
  <kbd>mouse wheel up</kbd>: Scroll up (fn+down)
</pre>

## Main panel (patch building)

<pre>
  <kbd>&lt;left&gt;</kbd>: Select previous hunk
  <kbd>&lt;right&gt;</kbd>: Select next hunk
  <kbd>v</kbd>: Toggle drag select
  <kbd>V</kbd>: Toggle drag select
  <kbd>a</kbd>: Toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: Copy the selected text to the clipboard
  <kbd>o</kbd>: Open file
  <kbd>e</kbd>: Edit file
  <kbd>&lt;space&gt;</kbd>: Add/Remove line(s) to patch
  <kbd>&lt;esc&gt;</kbd>: Exit custom patch builder
  <kbd>/</kbd>: Search the current view by text
</pre>

## Main panel (staging)

<pre>
  <kbd>&lt;left&gt;</kbd>: Select previous hunk
  <kbd>&lt;right&gt;</kbd>: Select next hunk
  <kbd>v</kbd>: Toggle drag select
  <kbd>V</kbd>: Toggle drag select
  <kbd>a</kbd>: Toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: Copy the selected text to the clipboard
  <kbd>o</kbd>: Open file
  <kbd>e</kbd>: Edit file
  <kbd>&lt;esc&gt;</kbd>: Return to files panel
  <kbd>&lt;tab&gt;</kbd>: Switch to other panel (staged/unstaged changes)
  <kbd>&lt;space&gt;</kbd>: Toggle line staged / unstaged
  <kbd>d</kbd>: Discard change (git reset)
  <kbd>E</kbd>: Edit hunk
  <kbd>c</kbd>: Commit changes
  <kbd>w</kbd>: Commit changes without pre-commit hook
  <kbd>C</kbd>: Commit changes using git editor
  <kbd>/</kbd>: Search the current view by text
</pre>

## Menu

<pre>
  <kbd>&lt;enter&gt;</kbd>: Execute
  <kbd>&lt;esc&gt;</kbd>: Close
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Reflog

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy commit SHA to clipboard
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Checkout commit
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: Open commit in browser
  <kbd>n</kbd>: Create new branch off of commit
  <kbd>g</kbd>: View reset options
  <kbd>c</kbd>: Copy commit (cherry-pick)
  <kbd>C</kbd>: Copy commit range (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: View commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Remote branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy branch name to clipboard
  <kbd>&lt;space&gt;</kbd>: Checkout
  <kbd>n</kbd>: New branch
  <kbd>M</kbd>: Merge into currently checked out branch
  <kbd>r</kbd>: Rebase checked-out branch onto this branch
  <kbd>d</kbd>: Delete branch
  <kbd>u</kbd>: Set as upstream of checked-out branch
  <kbd>g</kbd>: View reset options
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: View commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Remotes

<pre>
  <kbd>f</kbd>: Fetch remote
  <kbd>n</kbd>: Add new remote
  <kbd>d</kbd>: Remove remote
  <kbd>e</kbd>: Edit remote
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Stash

<pre>
  <kbd>&lt;space&gt;</kbd>: Apply
  <kbd>g</kbd>: Pop
  <kbd>d</kbd>: Drop
  <kbd>n</kbd>: New branch
  <kbd>r</kbd>: Rename stash
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: View selected item's files
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Status

<pre>
  <kbd>o</kbd>: Open config file
  <kbd>e</kbd>: Edit config file
  <kbd>u</kbd>: Check for update
  <kbd>&lt;enter&gt;</kbd>: Switch to a recent repo
  <kbd>a</kbd>: Show all branch logs
</pre>

## Sub-commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy commit SHA to clipboard
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Checkout commit
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: Open commit in browser
  <kbd>n</kbd>: Create new branch off of commit
  <kbd>g</kbd>: View reset options
  <kbd>c</kbd>: Copy commit (cherry-pick)
  <kbd>C</kbd>: Copy commit range (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: View selected item's files
  <kbd>/</kbd>: Search the current view by text
</pre>

## Submodules

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy submodule name to clipboard
  <kbd>&lt;enter&gt;</kbd>: Enter submodule
  <kbd>&lt;space&gt;</kbd>: Enter submodule
  <kbd>d</kbd>: Remove submodule
  <kbd>u</kbd>: Update submodule
  <kbd>n</kbd>: Add new submodule
  <kbd>e</kbd>: Update submodule URL
  <kbd>i</kbd>: Initialize submodule
  <kbd>b</kbd>: View bulk submodule options
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Tags

<pre>
  <kbd>&lt;space&gt;</kbd>: Checkout
  <kbd>d</kbd>: Delete tag
  <kbd>P</kbd>: Push tag
  <kbd>n</kbd>: Create tag
  <kbd>g</kbd>: View reset options
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: View commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Worktrees

<pre>
  <kbd>n</kbd>: Create worktree
  <kbd>&lt;space&gt;</kbd>: Switch to worktree
  <kbd>&lt;enter&gt;</kbd>: Switch to worktree
  <kbd>o</kbd>: Open in editor
  <kbd>d</kbd>: Remove worktree
  <kbd>/</kbd>: Filter the current view by text
</pre>
