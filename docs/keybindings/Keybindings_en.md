_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit Keybindings

## Global Keybindings

<pre>
  <kbd>ctrl+r</kbd>: switch to a recent repo
  <kbd>pgup</kbd>: scroll up main panel (fn+up/shift+k)
  <kbd>pgdown</kbd>: scroll down main panel (fn+down/shift+j)
  <kbd>m</kbd>: view merge/rebase options
  <kbd>ctrl+p</kbd>: view custom patch options
  <kbd>R</kbd>: refresh
  <kbd>x</kbd>: open menu
  <kbd>+</kbd>: next screen mode (normal/half/fullscreen)
  <kbd>_</kbd>: prev screen mode
  <kbd>ctrl+s</kbd>: view filter-by-path options
  <kbd>W</kbd>: open diff menu
  <kbd>ctrl+e</kbd>: open diff menu
  <kbd>@</kbd>: open command log menu
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: execute custom command
  <kbd>z</kbd>: undo (via reflog) (experimental)
  <kbd>ctrl+z</kbd>: redo (via reflog) (experimental)
  <kbd>P</kbd>: push
  <kbd>p</kbd>: pull
</pre>

## List Panel Navigation

<pre>
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
  <kbd>]</kbd>: next tab
  <kbd>[</kbd>: previous tab
</pre>

## Commit Files

<pre>
  <kbd>ctrl+o</kbd>: copy the committed file name to the clipboard
  <kbd>c</kbd>: checkout file
  <kbd>d</kbd>: discard this commit's changes to this file
  <kbd>o</kbd>: open file
  <kbd>e</kbd>: edit file
  <kbd>space</kbd>: toggle file included in patch
  <kbd>a</kbd>: toggle all files included in patch
  <kbd>enter</kbd>: enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: toggle file tree view
</pre>

## Commits

<pre>
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>b</kbd>: view bisect options
  <kbd>s</kbd>: squash down
  <kbd>f</kbd>: fixup commit
  <kbd>r</kbd>: reword commit
  <kbd>R</kbd>: reword commit with editor
  <kbd>d</kbd>: delete commit
  <kbd>e</kbd>: edit commit
  <kbd>p</kbd>: pick commit (when mid-rebase)
  <kbd>F</kbd>: create fixup commit for this commit
  <kbd>S</kbd>: squash all 'fixup!' commits above selected commit (autosquash)
  <kbd>ctrl+j</kbd>: move commit down one
  <kbd>ctrl+k</kbd>: move commit up one
  <kbd>v</kbd>: paste commits (cherry-pick)
  <kbd>A</kbd>: amend commit with staged changes
  <kbd>a</kbd>: reset commit author
  <kbd>t</kbd>: revert commit
  <kbd>T</kbd>: tag commit
  <kbd>ctrl+l</kbd>: open log menu
  <kbd>space</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: copy commit (cherry-pick)
  <kbd>C</kbd>: copy commit range (cherry-pick)
  <kbd>enter</kbd>: view selected item's files
</pre>

## Files

<pre>
  <kbd>ctrl+o</kbd>: copy the file name to the clipboard
  <kbd>ctrl+w</kbd>: Toggle whether or not whitespace changes are shown in the diff view
  <kbd>d</kbd>: view 'discard changes' options
  <kbd>space</kbd>: toggle staged
  <kbd>ctrl+b</kbd>: Filter files (staged/unstaged)
  <kbd>c</kbd>: commit changes
  <kbd>w</kbd>: commit changes without pre-commit hook
  <kbd>A</kbd>: amend last commit
  <kbd>C</kbd>: commit changes using git editor
  <kbd>e</kbd>: edit file
  <kbd>o</kbd>: open file
  <kbd>i</kbd>: ignore or exclude file
  <kbd>r</kbd>: refresh files
  <kbd>s</kbd>: stash all changes
  <kbd>S</kbd>: view stash options
  <kbd>a</kbd>: stage/unstage all
  <kbd>enter</kbd>: stage individual hunks/lines for file, or collapse/expand for directory
  <kbd>g</kbd>: view upstream reset options
  <kbd>D</kbd>: view reset options
  <kbd>`</kbd>: toggle file tree view
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>f</kbd>: fetch
</pre>

## Local Branches

<pre>
  <kbd>ctrl+o</kbd>: copy branch name to clipboard
  <kbd>i</kbd>: show git-flow options
  <kbd>space</kbd>: checkout
  <kbd>n</kbd>: new branch
  <kbd>o</kbd>: create pull request
  <kbd>O</kbd>: create pull request options
  <kbd>ctrl+y</kbd>: copy pull request URL to clipboard
  <kbd>c</kbd>: checkout by name
  <kbd>F</kbd>: force checkout
  <kbd>d</kbd>: delete branch
  <kbd>r</kbd>: rebase checked-out branch onto this branch
  <kbd>M</kbd>: merge into currently checked out branch
  <kbd>f</kbd>: fast-forward this branch from its upstream
  <kbd>g</kbd>: view reset options
  <kbd>R</kbd>: rename branch
  <kbd>u</kbd>: set/unset upstream
  <kbd>enter</kbd>: view commits
</pre>

## Main Panel (Merging)

<pre>
  <kbd>e</kbd>: edit file
  <kbd>o</kbd>: open file
  <kbd>◄</kbd>: select previous conflict
  <kbd>►</kbd>: select next conflict
  <kbd>▲</kbd>: select previous hunk
  <kbd>▼</kbd>: select next hunk
  <kbd>z</kbd>: undo
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>space</kbd>: pick hunk
  <kbd>b</kbd>: pick all hunks
  <kbd>esc</kbd>: return to files panel
</pre>

## Main Panel (Normal)

<pre>
  <kbd>mouse wheel ▼</kbd>: scroll down (fn+up)
  <kbd>mouse wheel ▲</kbd>: scroll up (fn+down)
</pre>

## Main Panel (Patch Building)

<pre>
  <kbd>◄</kbd>: select previous hunk
  <kbd>►</kbd>: select next hunk
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>ctrl+o</kbd>: copy the selected text to the clipboard
  <kbd>o</kbd>: open file
  <kbd>e</kbd>: edit file
  <kbd>space</kbd>: add/remove line(s) to patch
  <kbd>esc</kbd>: exit custom patch builder
</pre>

## Main Panel (Staging)

<pre>
  <kbd>◄</kbd>: select previous hunk
  <kbd>►</kbd>: select next hunk
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>ctrl+o</kbd>: copy the selected text to the clipboard
  <kbd>o</kbd>: open file
  <kbd>e</kbd>: edit file
  <kbd>esc</kbd>: return to files panel
  <kbd>tab</kbd>: switch to other panel (staged/unstaged changes)
  <kbd>space</kbd>: toggle line staged / unstaged
  <kbd>d</kbd>: delete change (git reset)
  <kbd>E</kbd>: edit hunk
  <kbd>c</kbd>: commit changes
  <kbd>w</kbd>: commit changes without pre-commit hook
  <kbd>C</kbd>: commit changes using git editor
</pre>

## Reflog

<pre>
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
  <kbd>space</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: copy commit (cherry-pick)
  <kbd>C</kbd>: copy commit range (cherry-pick)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>enter</kbd>: view commits
</pre>

## Remote Branches

<pre>
  <kbd>space</kbd>: checkout
  <kbd>n</kbd>: new branch
  <kbd>M</kbd>: merge into currently checked out branch
  <kbd>r</kbd>: rebase checked-out branch onto this branch
  <kbd>d</kbd>: delete branch
  <kbd>u</kbd>: set as upstream of checked-out branch
  <kbd>esc</kbd>: Return to remotes list
  <kbd>g</kbd>: view reset options
  <kbd>enter</kbd>: view commits
</pre>

## Remotes

<pre>
  <kbd>f</kbd>: fetch remote
  <kbd>n</kbd>: add new remote
  <kbd>d</kbd>: remove remote
  <kbd>e</kbd>: edit remote
</pre>

## Stash

<pre>
  <kbd>space</kbd>: apply
  <kbd>g</kbd>: pop
  <kbd>d</kbd>: drop
  <kbd>n</kbd>: new branch
  <kbd>r</kbd>: rename stash
  <kbd>enter</kbd>: view selected item's files
</pre>

## Status

<pre>
  <kbd>e</kbd>: edit config file
  <kbd>o</kbd>: open config file
  <kbd>u</kbd>: check for update
  <kbd>enter</kbd>: switch to a recent repo
  <kbd>a</kbd>: show all branch logs
</pre>

## Sub-commits

<pre>
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
  <kbd>space</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: copy commit (cherry-pick)
  <kbd>C</kbd>: copy commit range (cherry-pick)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>enter</kbd>: view selected item's files
</pre>

## Submodules

<pre>
  <kbd>ctrl+o</kbd>: copy submodule name to clipboard
  <kbd>enter</kbd>: enter submodule
  <kbd>d</kbd>: remove submodule
  <kbd>u</kbd>: update submodule
  <kbd>n</kbd>: add new submodule
  <kbd>e</kbd>: update submodule URL
  <kbd>i</kbd>: initialize submodule
  <kbd>b</kbd>: view bulk submodule options
</pre>

## Tags

<pre>
  <kbd>space</kbd>: checkout
  <kbd>d</kbd>: delete tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: create tag
  <kbd>g</kbd>: view reset options
  <kbd>enter</kbd>: view commits
</pre>
