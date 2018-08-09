# Keybindings:

## Global:

      ← → ↑ ↓/h j k l:  navigate
      PgUp/PgDn:        scroll diff panel (use fn+up/down on osx)
      q:                quit
      p:                pull
      shift+P:          push

## Files Panel:

      space:    toggle staged
      c:        commit changes
      shift+S:  stash files
      t:        add patched (i.e. pick chunks of a file to add)
      o:        open
      e:        edit
      s:        open in sublime (requires 'subl' command)
      v:        open in vscode (requires 'code' command)
      i:        add to .gitignore
      d:        delete if untracked checkout if tracked (aka go away)
      shift+R:  refresh files

## Branches Panel:

      space:    checkout branch
      f:        force checkout branch
      m:        merge into currently checked out branch
      c:        checkout by name
      n:        new branch

## Commits Panel:

      s:       squash down (only available for topmost commit)
      r:       rename commit
      g:       reset to this commit

## Stash Panel:

      space:   apply
      g:       pop
      d:       drop

## Popup Panel:

      esc:     close/cancel
      enter:   confirm

## Resolving Merge Conflicts (Diff Panel):

      ← →/h l:    navigate conflicts
      ↑ ↓/ k j:   select hunk
      space:      pick hunk
      b:          pick both hunks
      z:          undo (only available while still inside diff panel)
