_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Global Keybindings

<pre>
  <kbd>&lt;c-r&gt;</kbd>: switch to a recent repo
  <kbd>&lt;pgup&gt;</kbd>: scroll up main panel (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: scroll down main panel (fn+down/shift+j)
  <kbd>@</kbd>: open command log menu
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: execute custom command
  <kbd>&lt;c-p&gt;</kbd>: view custom patch options
  <kbd>m</kbd>: view merge/rebase options
  <kbd>R</kbd>: refresh
  <kbd>+</kbd>: next screen mode (normal/half/fullscreen)
  <kbd>_</kbd>: prev screen mode
  <kbd>?</kbd>: open menu
  <kbd>&lt;c-s&gt;</kbd>: view filter-by-path options
  <kbd>W</kbd>: open diff menu
  <kbd>&lt;c-e&gt;</kbd>: open diff menu
  <kbd>&lt;c-w&gt;</kbd>: Toggle whether or not whitespace changes are shown in the diff view
  <kbd>z</kbd>: undo (via reflog) (experimental)
  <kbd>&lt;c-z&gt;</kbd>: redo (via reflog) (experimental)
  <kbd>P</kbd>: push
  <kbd>p</kbd>: pull
</pre>

## List Panel Navigation

<pre>
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd>&lt;</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>&gt;</kbd>: scroll to bottom
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
  <kbd>]</kbd>: next tab
  <kbd>[</kbd>: previous tab
</pre>

## Commit Files

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy the committed file name to the clipboard
  <kbd>c</kbd>: checkout file
  <kbd>d</kbd>: discard this commit's changes to this file
  <kbd>o</kbd>: open file
  <kbd>e</kbd>: edit file
  <kbd>&lt;space&gt;</kbd>: toggle file included in patch
  <kbd>a</kbd>: toggle all files included in patch
  <kbd>&lt;enter&gt;</kbd>: enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: toggle file tree view
</pre>

## Commit Summary

<pre>
  <kbd>&lt;enter&gt;</kbd>: confirm
  <kbd>&lt;esc&gt;</kbd>: close
</pre>

## Commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy commit SHA to clipboard
  <kbd>&lt;c-r&gt;</kbd>: reset cherry-picked (copied) commits selection
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
  <kbd>&lt;c-j&gt;</kbd>: move commit down one
  <kbd>&lt;c-k&gt;</kbd>: move commit up one
  <kbd>v</kbd>: paste commits (cherry-pick)
  <kbd>A</kbd>: amend commit with staged changes
  <kbd>a</kbd>: reset commit author
  <kbd>t</kbd>: revert commit
  <kbd>T</kbd>: tag commit
  <kbd>&lt;c-l&gt;</kbd>: open log menu
  <kbd>&lt;space&gt;</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: copy commit (cherry-pick)
  <kbd>C</kbd>: copy commit range (cherry-pick)
  <kbd>&lt;enter&gt;</kbd>: view selected item's files
</pre>

## Confirmation Panel

<pre>
  <kbd>&lt;enter&gt;</kbd>: confirm
  <kbd>&lt;esc&gt;</kbd>: close/cancel
</pre>

## Files

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy the file name to the clipboard
  <kbd>d</kbd>: view 'discard changes' options
  <kbd>&lt;space&gt;</kbd>: toggle staged
  <kbd>&lt;c-b&gt;</kbd>: Filter files (staged/unstaged)
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
  <kbd>&lt;enter&gt;</kbd>: stage individual hunks/lines for file, or collapse/expand for directory
  <kbd>g</kbd>: view upstream reset options
  <kbd>D</kbd>: view reset options
  <kbd>`</kbd>: toggle file tree view
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>f</kbd>: fetch
</pre>

## Local Branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy branch name to clipboard
  <kbd>i</kbd>: show git-flow options
  <kbd>&lt;space&gt;</kbd>: checkout
  <kbd>n</kbd>: new branch
  <kbd>o</kbd>: create pull request
  <kbd>O</kbd>: create pull request options
  <kbd>&lt;c-y&gt;</kbd>: copy pull request URL to clipboard
  <kbd>c</kbd>: checkout by name
  <kbd>F</kbd>: force checkout
  <kbd>d</kbd>: delete branch
  <kbd>r</kbd>: rebase checked-out branch onto this branch
  <kbd>M</kbd>: merge into currently checked out branch
  <kbd>f</kbd>: fast-forward this branch from its upstream
  <kbd>T</kbd>: create tag
  <kbd>g</kbd>: view reset options
  <kbd>R</kbd>: rename branch
  <kbd>u</kbd>: set/unset upstream
  <kbd>&lt;enter&gt;</kbd>: view commits
</pre>

## Main Panel (Merging)

<pre>
  <kbd>e</kbd>: edit file
  <kbd>o</kbd>: open file
  <kbd>&lt;left&gt;</kbd>: select previous conflict
  <kbd>&lt;right&gt;</kbd>: select next conflict
  <kbd>&lt;up&gt;</kbd>: select previous hunk
  <kbd>&lt;down&gt;</kbd>: select next hunk
  <kbd>z</kbd>: undo
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>&lt;space&gt;</kbd>: pick hunk
  <kbd>b</kbd>: pick all hunks
  <kbd>&lt;esc&gt;</kbd>: return to files panel
</pre>

## Main Panel (Normal)

<pre>
  <kbd>mouse wheel down</kbd>: scroll down (fn+up)
  <kbd>mouse wheel up</kbd>: scroll up (fn+down)
</pre>

## Main Panel (Patch Building)

<pre>
  <kbd>&lt;left&gt;</kbd>: select previous hunk
  <kbd>&lt;right&gt;</kbd>: select next hunk
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: copy the selected text to the clipboard
  <kbd>o</kbd>: open file
  <kbd>e</kbd>: edit file
  <kbd>&lt;space&gt;</kbd>: add/remove line(s) to patch
  <kbd>&lt;esc&gt;</kbd>: exit custom patch builder
</pre>

## Main Panel (Staging)

<pre>
  <kbd>&lt;left&gt;</kbd>: select previous hunk
  <kbd>&lt;right&gt;</kbd>: select next hunk
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: copy the selected text to the clipboard
  <kbd>o</kbd>: open file
  <kbd>e</kbd>: edit file
  <kbd>&lt;esc&gt;</kbd>: return to files panel
  <kbd>&lt;tab&gt;</kbd>: switch to other panel (staged/unstaged changes)
  <kbd>&lt;space&gt;</kbd>: toggle line staged / unstaged
  <kbd>d</kbd>: delete change (git reset)
  <kbd>E</kbd>: edit hunk
  <kbd>c</kbd>: commit changes
  <kbd>w</kbd>: commit changes without pre-commit hook
  <kbd>C</kbd>: commit changes using git editor
</pre>

## Menu

<pre>
  <kbd>&lt;enter&gt;</kbd>: execute
  <kbd>&lt;esc&gt;</kbd>: close
</pre>

## Reflog

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy commit SHA to clipboard
  <kbd>&lt;space&gt;</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: copy commit (cherry-pick)
  <kbd>C</kbd>: copy commit range (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: view commits
</pre>

## Remote Branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy branch name to clipboard
  <kbd>&lt;space&gt;</kbd>: checkout
  <kbd>n</kbd>: new branch
  <kbd>M</kbd>: merge into currently checked out branch
  <kbd>r</kbd>: rebase checked-out branch onto this branch
  <kbd>d</kbd>: delete branch
  <kbd>u</kbd>: set as upstream of checked-out branch
  <kbd>&lt;esc&gt;</kbd>: Return to remotes list
  <kbd>g</kbd>: view reset options
  <kbd>&lt;enter&gt;</kbd>: view commits
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
  <kbd>&lt;space&gt;</kbd>: apply
  <kbd>g</kbd>: pop
  <kbd>d</kbd>: drop
  <kbd>n</kbd>: new branch
  <kbd>r</kbd>: rename stash
  <kbd>&lt;enter&gt;</kbd>: view selected item's files
</pre>

## Status

<pre>
  <kbd>o</kbd>: open config file
  <kbd>e</kbd>: edit config file
  <kbd>u</kbd>: check for update
  <kbd>&lt;enter&gt;</kbd>: switch to a recent repo
  <kbd>a</kbd>: show all branch logs
</pre>

## Sub-commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy commit SHA to clipboard
  <kbd>&lt;space&gt;</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: copy commit (cherry-pick)
  <kbd>C</kbd>: copy commit range (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: view selected item's files
</pre>

## Submodules

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy submodule name to clipboard
  <kbd>&lt;enter&gt;</kbd>: enter submodule
  <kbd>d</kbd>: remove submodule
  <kbd>u</kbd>: update submodule
  <kbd>n</kbd>: add new submodule
  <kbd>e</kbd>: update submodule URL
  <kbd>i</kbd>: initialize submodule
  <kbd>b</kbd>: view bulk submodule options
</pre>

## Tags

<pre>
  <kbd>&lt;space&gt;</kbd>: checkout
  <kbd>d</kbd>: delete tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: create tag
  <kbd>g</kbd>: view reset options
  <kbd>&lt;enter&gt;</kbd>: view commits
</pre>
