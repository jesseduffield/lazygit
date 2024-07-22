_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Sneltoetsen

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Globale sneltoetsen

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Wissel naar een recente repo |  |
| `` <pgup> (fn+up/shift+k) `` | Scroll naar beneden vanaf hoofdpaneel |  |
| `` <pgdown> (fn+down/shift+j) `` | Scroll naar beneden vanaf hoofdpaneel |  |
| `` @ `` | View command log options | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | Push | Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` p `` | Pull | Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` } `` | Increase diff context size | Increase the amount of the context shown around changes in the diff view. |
| `` { `` | Decrease diff context size | Decrease the amount of the context shown around changes in the diff view. |
| `` : `` | Voer aangepaste commando uit | Bring up a prompt where you can enter a shell command to execute. Not to be confused with pre-configured custom commands. |
| `` <c-p> `` | Bekijk aangepaste patch opties |  |
| `` m `` | Bekijk merge/rebase opties | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | Verversen | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | Volgende scherm modus (normaal/half/groot) |  |
| `` _ `` | Vorige scherm modus |  |
| `` ? `` | Open menu |  |
| `` <c-s> `` | Bekijk scoping opties | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | Open diff menu | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | Open diff menu | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | Quit |  |
| `` <esc> `` | Annuleren |  |
| `` <c-w> `` | Toggle whitespace | Toggle whether or not whitespace changes are shown in the diff view. |
| `` z `` | Ongedaan maken (via reflog) (experimenteel) | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | Redo (via reflog) (experimenteel) | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |

## Lijstpaneel navigatie

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Vorige pagina |  |
| `` . `` | Volgende pagina |  |
| `` < `` | Scroll naar boven |  |
| `` > `` | Scroll naar beneden |  |
| `` v `` | Toggle drag selecteer |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | Start met zoeken |  |
| `` H `` | Scroll left |  |
| `` L `` | Scroll right |  |
| `` ] `` | Volgende tabblad |  |
| `` [ `` | Vorige tabblad |  |

## Bestanden

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer de bestandsnaam naar het klembord |  |
| `` <space> `` | Toggle staged | Toggle staged for selected file. |
| `` <c-b> `` | Filter files by status |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Commit veranderingen | Commit staged changes. |
| `` w `` | Commit veranderingen zonder pre-commit hook |  |
| `` A `` | Wijzig laatste commit |  |
| `` C `` | Commit veranderingen met de git editor |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit | Open file in external editor. |
| `` o `` | Open bestand | Open file in default application. |
| `` i `` | Ignore or exclude file |  |
| `` r `` | Refresh bestanden |  |
| `` s `` | Stash | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | Bekijk stash opties | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | Toggle staged alle | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | Stage individuele hunks/lijnen | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | Bekijk 'veranderingen ongedaan maken' opties | View options for discarding changes to the selected file. |
| `` g `` | Bekijk upstream reset opties |  |
| `` D `` | Reset | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | Toggle bestandsboom weergave | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Open external merge tool | Run `git mergetool`. |
| `` f `` | Fetch | Fetch changes from remote. |
| `` / `` | Start met zoeken |  |

## Bevestigingspaneel

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Bevestig |  |
| `` <esc> `` | Sluiten |  |

## Branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer branch name naar klembord |  |
| `` i `` | Laat git-flow opties zien |  |
| `` <space> `` | Uitchecken | Checkout selected item. |
| `` n `` | Nieuwe branch |  |
| `` o `` | Maak een pull-request |  |
| `` O `` | Bekijk opties voor pull-aanvraag |  |
| `` <c-y> `` | Kopieer de URL van het pull-verzoek naar het klembord |  |
| `` c `` | Uitchecken bij naam | Checkout by name. In the input box you can enter '-' to switch to the last branch. |
| `` F `` | Forceer checkout | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | Delete | View delete options for local/remote branch. |
| `` r `` | Rebase branch | Rebase the checked-out branch onto the selected branch. |
| `` M `` | Merge in met huidige checked out branch | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` f `` | Fast-forward deze branch vanaf zijn upstream | Fast-forward selected branch from its upstream. |
| `` T `` | Creëer tag |  |
| `` s `` | Sort order |  |
| `` g `` | Bekijk reset opties |  |
| `` R `` | Hernoem branch |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Bekijk commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Commit bericht

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Bevestig |  |
| `` <esc> `` | Sluiten |  |

## Commit bestanden

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer de bestandsnaam naar het klembord |  |
| `` c `` | Uitchecken | Bestand uitchecken |
| `` d `` | Remove | Uitsluit deze commit zijn veranderingen aan dit bestand |
| `` o `` | Open bestand | Open file in default application. |
| `` e `` | Edit | Open file in external editor. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle bestand inbegrepen in patch | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Toggle all files | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Enter bestand om geselecteerde regels toe te voegen aan de patch | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | Toggle bestandsboom weergave | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` / `` | Start met zoeken |  |

## Commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer commit hash naar klembord |  |
| `` <c-r> `` | Reset cherry-picked (gekopieerde) commits selectie |  |
| `` b `` | View bisect options |  |
| `` s `` | Squash | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | Fixup | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` r `` | Hernoem commit | Reword the selected commit's message. |
| `` R `` | Hernoem commit met editor |  |
| `` d `` | Verwijder commit | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | Edit (start interactive rebase) | Wijzig commit |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick | Kies commit (wanneer midden in rebase) |
| `` F `` | Creëer fixup commit | Creëer fixup commit |
| `` S `` | Apply fixup commits | Squash bovenstaande commits |
| `` <c-j> `` | Verplaats commit 1 naar beneden |  |
| `` <c-k> `` | Verplaats commit 1 naar boven |  |
| `` V `` | Plak commits (cherry-pick) |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Amend | Wijzig commit met staged veranderingen |
| `` a `` | Amend commit attribute | Set/Reset commit author or set co-author. |
| `` t `` | Revert | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | Tag commit | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <c-l> `` | View log options | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` <space> `` | Uitchecken | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Creëer nieuwe branch van commit |  |
| `` g `` | Bekijk reset opties | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Kopieer commit (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Bekijk gecommite bestanden |  |
| `` w `` | View worktree options |  |
| `` / `` | Start met zoeken |  |

## Menu

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Uitvoeren |  |
| `` <esc> `` | Sluiten |  |
| `` / `` | Filter the current view by text |  |

## Mergen

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Kies stuk |  |
| `` b `` | Kies beide stukken |  |
| `` <up> `` | Selecteer bovenste hunk |  |
| `` <down> `` | Selecteer onderste hunk |  |
| `` <left> `` | Selecteer voorgaand conflict |  |
| `` <right> `` | Selecteer volgende conflict |  |
| `` z `` | Ongedaan maken | Undo last merge conflict resolution. |
| `` e `` | Verander bestand | Open file in external editor. |
| `` o `` | Open bestand | Open file in default application. |
| `` M `` | Open external merge tool | Run `git mergetool`. |
| `` <esc> `` | Ga terug naar het bestanden paneel |  |

## Normaal

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | Scroll omlaag |  |
| `` mouse wheel up (fn+down) `` | Scroll omhoog |  |

## Patch bouwen

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Selecteer de vorige hunk |  |
| `` <right> `` | Selecteer de volgende hunk |  |
| `` v `` | Toggle drag selecteer |  |
| `` a `` | Toggle selecteer hunk | Toggle hunk selection mode. |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` o `` | Open bestand | Open file in default application. |
| `` e `` | Verander bestand | Open file in external editor. |
| `` <space> `` | Voeg toe/verwijder lijn(en) in patch |  |
| `` <esc> `` | Sluit lijn-bij-lijn modus |  |
| `` / `` | Start met zoeken |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer commit hash naar klembord |  |
| `` <space> `` | Uitchecken | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Creëer nieuwe branch van commit |  |
| `` g `` | Bekijk reset opties | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Kopieer commit (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset cherry-picked (gekopieerde) commits selectie |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Bekijk commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Remote branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer branch name naar klembord |  |
| `` <space> `` | Uitchecken | Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head. |
| `` n `` | Nieuwe branch |  |
| `` M `` | Merge in met huidige checked out branch | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` r `` | Rebase branch | Rebase the checked-out branch onto the selected branch. |
| `` d `` | Delete | Delete the remote branch from the remote. |
| `` u `` | Set as upstream | Stel in als upstream van uitgecheckte branch |
| `` s `` | Sort order |  |
| `` g `` | Bekijk reset opties | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Bekijk commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Remotes

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | Voeg een nieuwe remote toe |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Edit | Wijzig remote |
| `` f `` | Fetch | Fetch remote |
| `` / `` | Filter the current view by text |  |

## Staging

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Selecteer de vorige hunk |  |
| `` <right> `` | Selecteer de volgende hunk |  |
| `` v `` | Toggle drag selecteer |  |
| `` a `` | Toggle selecteer hunk | Toggle hunk selection mode. |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` <space> `` | Toggle staged | Toggle lijnen staged / unstaged |
| `` d `` | Verwijdert change (git reset) | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | Open bestand | Open file in default application. |
| `` e `` | Verander bestand | Open file in external editor. |
| `` <esc> `` | Ga terug naar het bestanden paneel |  |
| `` <tab> `` | Ga naar een ander paneel | Switch to other view (staged/unstaged changes). |
| `` E `` | Edit hunk | Edit selected hunk in external editor. |
| `` c `` | Commit veranderingen | Commit staged changes. |
| `` w `` | Commit veranderingen zonder pre-commit hook |  |
| `` C `` | Commit veranderingen met de git editor |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | Start met zoeken |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Toepassen | Apply the stash entry to your working directory. |
| `` g `` | Pop | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | Laten vallen | Remove the stash entry from the stash list. |
| `` n `` | Nieuwe branch | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | Rename stash |  |
| `` <enter> `` | Bekijk gecommite bestanden |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Open config bestand | Open file in default application. |
| `` e `` | Verander config bestand | Open file in external editor. |
| `` u `` | Check voor updates |  |
| `` <enter> `` | Wissel naar een recente repo |  |
| `` a `` | Alle logs van de branch laten zien |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer commit hash naar klembord |  |
| `` <space> `` | Uitchecken | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Creëer nieuwe branch van commit |  |
| `` g `` | Bekijk reset opties | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Kopieer commit (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset cherry-picked (gekopieerde) commits selectie |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Bekijk gecommite bestanden |  |
| `` w `` | View worktree options |  |
| `` / `` | Start met zoeken |  |

## Submodules

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer submodule naam naar klembord |  |
| `` <enter> `` | Enter | Enter submodule |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | Update selected submodule. |
| `` n `` | Voeg nieuwe submodule toe |  |
| `` e `` | Update submodule URL |  |
| `` i `` | Initialize | Initialiseer submodule |
| `` b `` | Bekijk bulk submodule opties |  |
| `` / `` | Filter the current view by text |  |

## Tags

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Uitchecken | Checkout the selected tag tag as a detached HEAD. |
| `` n `` | Creëer tag | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Delete | View delete options for local/remote tag. |
| `` P `` | Push tag | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Bekijk commits |  |
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
