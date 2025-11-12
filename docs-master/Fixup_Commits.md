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

If you have many modifications in your working copy, it is a good idea to stage
related changes that are meant to go into the same fixup commit; if no changes
are staged, ctrl-f works on all unstaged modifications, and then it might show
an error if it finds multiple different base commits. If you are interested in
what the command does to do its magic, and how you can help it work better, you
may want to read the [design document](dev/Find_Base_Commit_For_Fixup_Design.md)
that describes this.
