# Keybindings:

## Global:

<pre>
  <kbd>←</kbd><kbd>→</kbd><kbd>↑</kbd><kbd>↓</kbd>/<kbd>h</kbd><kbd>j</kbd><kbd>k</kbd><kbd>l</kbd>:               navigate
  <kbd>PgUp</kbd>/<kbd>PgDn</kbd> or <kbd>ctrl</kbd>+<kbd>u</kbd>/<kbd>ctrl</kbd>+<kbd>d</kbd>:   scroll diff panel
                                     (for <kbd>PgUp</kbd> and <kbd>PgDn</kbd>, use <kbd>fn</kbd>+<kbd>up</kbd>/<kbd>fn</kbd>+<kbd>down</kbd> on osx)
  <kbd>q</kbd>:                                quit
  <kbd>p</kbd>:                                pull
  <kbd>shift</kbd>+<kbd>P</kbd>:                         push
</pre>

## Status Panel:

<pre>
  <kbd>e</kbd>:        edit config file
  <kbd>o</kbd>:        open config file
</pre>

## Files Panel:

<pre>
  <kbd>space</kbd>:    toggle staged
  <kbd>a</kbd>:        stage/unstage all
  <kbd>c</kbd>:        commit changes
  <kbd>shift</kbd>+<kbd>C</kbd>: commit using git editor
  <kbd>shift</kbd>+<kbd>S</kbd>: stash files
  <kbd>t</kbd>:        add patched (i.e. pick chunks of a file to add)
  <kbd>o</kbd>:        open
  <kbd>e</kbd>:        edit
  <kbd>s</kbd>:        open in sublime (requires 'subl' command)
  <kbd>v</kbd>:        open in vscode (requires 'code' command)
  <kbd>i</kbd>:        add to .gitignore
  <kbd>d</kbd>:        delete if untracked checkout if tracked (aka go away)
  <kbd>shift</kbd>+<kbd>R</kbd>: refresh files
  <kbd>shift</kbd>+<kbd>A</kbd>: abort merge
</pre>

## Branches Panel:

<pre>
  <kbd>space</kbd>:   checkout branch
  <kbd>f</kbd>:       force checkout branch
  <kbd>m</kbd>:       merge into currently checked out branch
  <kbd>c</kbd>:       checkout by name
  <kbd>n</kbd>:       new branch
  <kbd>d</kbd>:       delete branch
  <kbd>D</kbd>:       force delete branch
</pre>

## Commits Panel:

<pre>
  <kbd>s</kbd>:       squash down (only available for topmost commit)
  <kbd>r</kbd>:       rename commit
  <kbd>shift</kbd>+<kbd>R</kbd>: rename commit using git editor
  <kbd>g</kbd>:       reset to this commit
</pre>

## Stash Panel:

<pre>
  <kbd>space</kbd>:   apply
  <kbd>g</kbd>:       pop
  <kbd>d</kbd>:       drop
</pre>

## Popup Panel:

<pre>
  <kbd>esc</kbd>:     close/cancel
  <kbd>enter</kbd>:   confirm
  <kbd>tab</kbd>:     enter newline (if editing)
</pre>

## Resolving Merge Conflicts (Diff Panel):

<pre>
  <kbd>←</kbd><kbd>→</kbd>/<kbd>h</kbd><kbd>l</kbd>: navigate conflicts
  <kbd>↑</kbd><kbd>↓</kbd>/<kbd>k</kbd><kbd>j</kbd>: select hunk
  <kbd>space</kbd>:      pick hunk
  <kbd>b</kbd>:         pick both hunks
  <kbd>z</kbd>:         undo (only available while still inside diff panel)
</pre>
