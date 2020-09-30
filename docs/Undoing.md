# Undo/Redo in lazygit

![Gif](../../assets/undo2.gif)

## Keybindings:
'z' to undo, 'ctrl+z' to redo

## How it works

If you're as clumsy as me you'll probably have felt the pain of botching an interactive rebase or doing a hard reset onto the wrong commit. Luckily, the reflog allows you to trace your steps and make things right again, but I personally can't stand trying to make sense of the reflog.

Lazygit can read through your reflog for you and walk back action by action so that you don't even need to read the reflog. If lazygit finds a reflog entry where you checked out a branch, we'll checkout the original branch. If the entry is from a commit being applied, we'll go back to the commit before that. If we hit an interactive rebase, we'll go back to the commit you were on just before you started it.

## You can even undo things you did outside of lazygit!

Because lazygit just uses the reflog to keep track of things, it doesn't matter whether you're trying to undo something you did in lazygit or directly on the command line. You can open lazygit for the first time and start undoing thing in your repo! Likewise, lazygit marks its undos/redos in the reflog so if you quit the application and come back, lazygit still knows where you're up to.

## Limitations

There are limitations: firstly, lazygit can only undo things that are recorded in the reflog. That means changes to your working tree or stash aren't covered. Secondly, anything permanent you do like pushing to a remote can't be undone. Thirdly, actions like creating a branch won't be undone, because they're not stored in the reflog.

If you are mid-rebase, undo/redo is not supported, because the reflog doesn't contain enough information about what specific things have happened inside that rebase. If you want to undo out of a rebase, it's best to abort the rebase (the default keybinding for bringing up rebase options is 'm').

Undo/Redo is a new feature so if you find a bug let us know. The worst case scenario is that you'll just need to look at your reflog and manually put yourself back on track.
