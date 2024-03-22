# Fixup Commits

## Background

There's this common scenario that you have a PR in review, the reviewer is
requesting some changes, and you make those changes and would normally simply
squash them into the original commit that they came from. If you do that,
however, there's no way for the reviewer to see what you changed. You could just
make a separate commit with those changes at the end of the branch, but this is
not ideal because it results in a git history that is not very clean.

To help with this, git has a concept of fixup commits: you do make a separate
commit, but the subject of this commit is the string "fixup! " followed by the
original commit subject. This both tells the reviewer what's going on (you are
making a change that you later will squash into the designated commit), and it
provides an easy way to actually perform this squash operation when you are
ready to do that (before merging).

## Creating fixup commits

You could of course create fixup commits manually by typing in the commit
message with the prefix yourself. But lazygit has an easier way to do that:
in the Commits view, select the commit that you want to create a fixup for, and
press shift-F (for "Create fixup commit for this commit"). This automatically
creates a commit with the appropriate subject line.

Don't confuse this with the lowercase "f" command ("Fixup commit"); that one
squashes the selected commit into its parent, this is not what we want here.

## Creating amend commits

There's a special type of fixup commit that uses "amend!" instead of "fixup!" in
the commit message subject; in addition to fixing up the original commit with
changes it allows you to also (or only) change the commit message of the
original commit. The menu that appears when pressing shift-F has options for
both of these; they bring up a commit message panel similar to when you reword a
commit, but then create the "amend!" commit containing the new message. Note
that in that panel you only type the new message as you want it to be
eventually; lazygit then takes care of formatting the "amend!" commit
appropriately for you (with the subject of your new message moving into the body
of the "amend!" commit).

## Squashing fixup commits

When you're ready to merge the branch and want to squash all these fixup commits
that you created, that's very easy to do: select the first commit of your branch
and hit shift-S (for "Squash all 'fixup!' commits above selected commit
(autosquash)"). Boom, done.

## Finding the commit to create a fixup for

When you are making changes to code that you changed earlier in a long branch,
it can be tedious to find the commit to squash it into. Lazygit has a command to
help you with this, too: in the Files view, press ctrl-f to select the right
base commit in the Commits view automatically. From there, you can either press
shift-F to create a fixup commit for it, or shift-A to amend your changes into
the commit if you haven't published your branch yet.

This command works in many cases, and when it does it almost feels like magic,
but it's important to understand its limitations because it doesn't always work.
The way it works is that it looks at the deleted lines of your current
modifications, blames them to find out which commit those lines come from, and
if they all come from the same commit, it selects it. So here are cases where it
doesn't work:

- Your current diff has only added lines, but no deleted lines. In this case
  there's no way for lazygit to know which commit you want to add them to.
- The deleted lines belong to multiple different commits. In this case you can
  help lazygit by staging a set of files or hunks that all belong to the same
  commit; if some changes are staged, the ctrl-f command works only on those.
- The found commit is already on master; in this case, lazygit refuses to select
  it, because it doesn't make sense to create fixups for it, let alone amend to
  it.

To sum it up: the command works great if you are changing code again that you
changed or added earlier in the same branch. This is a common enough case to
make the command useful.
