_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Sneltoetsen

## Globale sneltoetsen

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+r> `` | Wissel naar een recente repo |  |
| `` <pgup>, K, <ctrl+u> (fn+up/shift+k) `` | Scroll naar beneden vanaf hoofdpaneel |  |
| `` <pgdown>, J, <ctrl+d> (fn+down/shift+j) `` | Scroll naar beneden vanaf hoofdpaneel |  |
| `` @ `` | View command log options | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | Push | Push de huidige branch naar de bijbehorende upstream-branch. Als er geen upstream is geconfigureerd wordt er gevraagd om een upstream-branch te configureren. |
| `` p `` | Pull | Pull wijzigingen van de remote voor de huidige branch. Als er geen upstream is geconfigureerd wordt er gevraagd om een upstream-branch te configureren. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` } `` | Increase diff context size | Increase the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` { `` | Decrease diff context size | Decrease the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` : `` | Execute shell command | Bring up a prompt where you can enter a shell command to execute. |
| `` <ctrl+p> `` | Bekijk aangepaste patch opties |  |
| `` m `` | Bekijk merge/rebase opties | Toon abort/continue/skip opties voor huidige merge/rebase. |
| `` R `` | Verversen | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | Volgende scherm modus (normaal/half/groot) |  |
| `` _ `` | Vorige scherm modus |  |
| `` \| `` | Cycle pagers | Choose the next pager in the list of configured pagers. |
| `` \ `` | Cycle pagers (reverse) | Choose the previous pager in the list of configured pagers. |
| `` <esc> `` | Annuleren |  |
| `` ? `` | Open menu |  |
| `` <ctrl+s> `` | Bekijk scoping opties | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W, <ctrl+e> `` | Open diff menu | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q, <ctrl+c> `` | Afsluiten |  |
| `` <ctrl+z> `` | Pauzeer de applicatie |  |
| `` <ctrl+w> `` | Toggle whitespace | Toggle whether or not whitespace changes are shown in the diff view.<br><br>The default can be changed in the config file with the key 'git.ignoreWhitespaceInDiffView'. |
| `` <alt+shift+c> `` | Verander config bestand | Open bestand in externe editor. |
| `` z `` | Ongedaan maken (via reflog) (experimenteel) | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` Z `` | Redo (via reflog) (experimenteel) | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |

## Lijstpaneel navigatie

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Vorige pagina |  |
| `` . `` | Volgende pagina |  |
| `` <, <home> `` | Scroll naar boven |  |
| `` >, <end> `` | Scroll naar beneden |  |
| `` v `` | Toggle drag selecteer |  |
| `` <shift+down> `` | Range select down |  |
| `` <shift+up> `` | Range select up |  |
| `` / `` | Start met zoeken |  |
| `` H `` | Scroll naar links |  |
| `` L `` | Scroll naar rechts |  |
| `` ] `` | Volgende tabblad |  |
| `` [ `` | Vorige tabblad |  |

## Bestanden

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopieer de bestandsnaam naar het klembord |  |
| `` <space> `` | Toggle staged | Toggle staged for selected file. |
| `` <ctrl+b> `` | Filter bestanden op status |  |
| `` y `` | Kopieer naar klembord |  |
| `` c `` | Commit veranderingen | Commit gestagede wijzigingen. |
| `` w `` | Commit veranderingen zonder pre-commit hook |  |
| `` A `` | Wijzig laatste commit |  |
| `` C `` | Commit veranderingen met de git editor |  |
| `` <ctrl+f> `` | Find base commit for fixup | Vind de commit waar je huidige wijzigingen bovenop zijn gebouwd met als doel die commit te amenden/fixen. Hierdoor hoef je dit niet met de hand te doen. Zie: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit | Open bestand in externe editor. |
| `` o `` | Open bestand | Open bestand in standaardapplicatie. |
| `` i `` | Ignore or exclude file |  |
| `` r `` | Refresh bestanden |  |
| `` s `` | Stash | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | Bekijk stash opties | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | Toggle staged alle | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | Stage individuele hunks/lijnen | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | Bekijk 'veranderingen ongedaan maken' opties | View options for discarding changes to the selected file. |
| `` g `` | Bekijk upstream reset opties |  |
| `` D `` | Resetten | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | Toggle bestandsboom weergave | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` <ctrl+t> `` | Open externe diff applicatie (git difftool) |  |
| `` M `` | Bekijk merge conflict opties | Bekijk opties voor het oplossen van mergeconflicten. |
| `` f `` | Fetch | Fetch changes from remote. |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Vouw alle bestanden uit | Vouw alle mappen in de bestandsstructuur uit |
| `` 0 `` | Focus main view |  |
| `` / `` | Filter the current view by text |  |

## Bevestigingspaneel

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Bevestig |  |
| `` <esc> `` | Sluiten |  |
| `` <ctrl+o> `` | Kopieer naar klembord |  |

## Branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopieer branch name naar klembord |  |
| `` i `` | Laat git-flow opties zien |  |
| `` <space> `` | Uitchecken | Geselecteerd item uitchecken. |
| `` n `` | Nieuwe branch |  |
| `` N `` | Verplaats commits naar nieuwe branch | Maak een nieuwe branch en verplaats niet-gepushte commits van de huidige branch hier naar toe. Gebruik dit in het geval dat je deze commits eigenlijk op een nieuwe branch had willen maken.<br><br>Let op dat de selectie genegeerd wordt. De nieuwe branch komt ofwel bovenop de main branch, of bovenop de huidige branch (je kan kiezen). |
| `` w `` | New worktree |  |
| `` o `` | Maak een pull-request |  |
| `` O `` | Bekijk opties voor pull-aanvraag |  |
| `` G `` | Open pull request in browser |  |
| `` <ctrl+y> `` | Kopieer de URL van het pull-verzoek naar het klembord |  |
| `` c `` | Uitchecken bij naam | Checkout by name. In the input box you can enter '-' to switch to the previous branch. |
| `` - `` | Vorige branch uitchecken |  |
| `` F `` | Forceer checkout | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | Verwijderen | View delete options for local/remote branch. |
| `` r `` | Rebase branch | Rebase de uitgecheckte branch bovenop de geselecteerde branch. |
| `` M `` | Merge in met huidige checked out branch | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` f `` | Fast-forward deze branch vanaf zijn upstream | Fast-forward selected branch from its upstream. |
| `` T `` | Creëer tag |  |
| `` s `` | Sort order |  |
| `` g `` | Bekijk reset opties |  |
| `` R `` | Hernoem branch |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <ctrl+t> `` | Open externe diff applicatie (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Bekijk commits |  |
| `` / `` | Filter the current view by text |  |

## Commit bericht

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Bevestig |  |
| `` <esc> `` | Sluiten |  |

## Commit bestanden

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopieer de bestandsnaam naar het klembord |  |
| `` y `` | Kopieer naar klembord |  |
| `` c `` | Uitchecken | Bestand uitchecken |
| `` d `` | Bekijk 'veranderingen ongedaan maken' opties | Uitsluit deze commit zijn veranderingen aan dit bestand |
| `` o `` | Open bestand | Open bestand in standaardapplicatie. |
| `` e `` | Edit | Open bestand in externe editor. |
| `` <ctrl+t> `` | Open externe diff applicatie (git difftool) |  |
| `` <space> `` | Toggle bestand inbegrepen in patch | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Toggle all files | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Enter bestand om geselecteerde regels toe te voegen aan de patch | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | Toggle bestandsboom weergave | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Vouw alle bestanden uit | Vouw alle mappen in de bestandsstructuur uit |
| `` 0 `` | Focus main view |  |
| `` / `` | Filter the current view by text |  |

## Commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <ctrl+r> `` | Reset cherry-picked (gekopieerde) commits selectie |  |
| `` b `` | View bisect options |  |
| `` s `` | Squash | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | Fixup | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` c `` | Set fixup message | Set the message option for the fixup commit. The -C option means to use this commit's message instead of the target commit's message. |
| `` r `` | Hernoem commit | Reword the selected commit's message. |
| `` R `` | Hernoem commit met editor |  |
| `` d `` | Verwijder commit | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | Edit (start interactive rebase) | Wijzig commit |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.<br>If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick | Kies commit (wanneer midden in rebase) |
| `` F `` | Creëer fixup commit | Creëer fixup commit |
| `` S `` | Apply fixup commits | Squash bovenstaande commits |
| `` <ctrl+j>, <alt+down> `` | Verplaats commit 1 naar beneden |  |
| `` <ctrl+k>, <alt+up> `` | Verplaats commit 1 naar boven |  |
| `` V `` | Plak commits (cherry-pick) |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Amend | Wijzig commit met staged veranderingen |
| `` a `` | Amend commit attribute | Set/Reset commit author or set co-author. |
| `` t `` | Revert | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | Tag commit | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <ctrl+l> `` | View log options | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` G `` | Open pull request in browser |  |
| `` <space> `` | Uitchecken | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Creëer nieuwe branch van commit |  |
| `` N `` | Verplaats commits naar nieuwe branch | Maak een nieuwe branch en verplaats niet-gepushte commits van de huidige branch hier naar toe. Gebruik dit in het geval dat je deze commits eigenlijk op een nieuwe branch had willen maken.<br><br>Let op dat de selectie genegeerd wordt. De nieuwe branch komt ofwel bovenop de main branch, of bovenop de huidige branch (je kan kiezen). |
| `` w `` | New worktree |  |
| `` g `` | Bekijk reset opties | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Kopieer commit (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <ctrl+t> `` | Open externe diff applicatie (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Bekijk gecommite bestanden |  |
| `` / `` | Start met zoeken |  |

## Input prompt

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Bevestig |  |
| `` <esc> `` | Sluiten |  |

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
| `` b `` | Pick both hunks |  |
| `` <up>, k `` | Selecteer bovenste hunk |  |
| `` <down>, j `` | Selecteer onderste hunk |  |
| `` <left>, h `` | Selecteer voorgaand conflict |  |
| `` <right>, l `` | Selecteer volgende conflict |  |
| `` z `` | Ongedaan maken | Undo last merge conflict resolution. |
| `` e `` | Verander bestand | Open bestand in externe editor. |
| `` o `` | Open bestand | Open bestand in standaardapplicatie. |
| `` M `` | Bekijk merge conflict opties | Bekijk opties voor het oplossen van mergeconflicten. |
| `` <esc> `` | Ga terug naar het bestanden paneel |  |

## Normaal

| Key | Action | Info |
|-----|--------|-------------|
| `` <mouse wheel down> (fn+up) `` | Scroll omlaag |  |
| `` <mouse wheel up> (fn+down) `` | Scroll omhoog |  |
| `` <tab> `` | Ga naar een ander paneel | Switch to other view (staged/unstaged changes). |
| `` <esc> `` | Exit back to side panel |  |
| `` <space> `` | Show/hide selection |  |
| `` / `` | Start met zoeken |  |

## Patch bouwen

| Key | Action | Info |
|-----|--------|-------------|
| `` <left>, h `` | Selecteer de vorige hunk |  |
| `` <right>, l `` | Selecteer de volgende hunk |  |
| `` v `` | Toggle drag selecteer |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <ctrl+o> `` | Copy selected text to clipboard |  |
| `` o `` | Open bestand | Open bestand in standaardapplicatie. |
| `` e `` | Verander bestand | Open bestand in externe editor. |
| `` <space> `` | Voeg toe/verwijder lijn(en) in patch |  |
| `` d `` | Remove lines from commit | Remove the selected lines from this commit. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes these lines. |
| `` <esc> `` | Sluit lijn-bij-lijn modus |  |
| `` / `` | Start met zoeken |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <space> `` | Uitchecken | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Creëer nieuwe branch van commit |  |
| `` N `` | Verplaats commits naar nieuwe branch | Maak een nieuwe branch en verplaats niet-gepushte commits van de huidige branch hier naar toe. Gebruik dit in het geval dat je deze commits eigenlijk op een nieuwe branch had willen maken.<br><br>Let op dat de selectie genegeerd wordt. De nieuwe branch komt ofwel bovenop de main branch, of bovenop de huidige branch (je kan kiezen). |
| `` w `` | New worktree |  |
| `` g `` | Bekijk reset opties | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Kopieer commit (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <ctrl+r> `` | Reset cherry-picked (gekopieerde) commits selectie |  |
| `` <ctrl+t> `` | Open externe diff applicatie (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Bekijk commits |  |
| `` / `` | Filter the current view by text |  |

## Remote branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopieer branch name naar klembord |  |
| `` <space> `` | Uitchecken | Geselecteerde remote branch uitchecken als nieuwe locale branch of als detached head. |
| `` n `` | Nieuwe branch |  |
| `` w `` | New worktree |  |
| `` M `` | Merge in met huidige checked out branch | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` r `` | Rebase branch | Rebase de uitgecheckte branch bovenop de geselecteerde branch. |
| `` d `` | Verwijderen | Delete the remote branch from the remote. |
| `` u `` | Instellen als upstream | Stel in als upstream van uitgecheckte branch |
| `` s `` | Sort order |  |
| `` g `` | Bekijk reset opties | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <ctrl+t> `` | Open externe diff applicatie (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Bekijk commits |  |
| `` / `` | Filter the current view by text |  |

## Remotes

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Bekijk branches |  |
| `` n `` | Voeg een nieuwe remote toe |  |
| `` d `` | Verwijderen | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Edit | Wijzig remote |
| `` f `` | Fetch | Fetch remote |
| `` F `` | Add fork remote | Quickly add a fork remote by replacing the owner in the origin URL and optionally check out a branch from new remote. |
| `` / `` | Filter the current view by text |  |

## Secondary

| Key | Action | Info |
|-----|--------|-------------|
| `` <tab> `` | Ga naar een ander paneel | Switch to other view (staged/unstaged changes). |
| `` <esc> `` | Exit back to side panel |  |
| `` <space> `` | Show/hide selection |  |
| `` / `` | Start met zoeken |  |

## Staging

| Key | Action | Info |
|-----|--------|-------------|
| `` <left>, h `` | Selecteer de vorige hunk |  |
| `` <right>, l `` | Selecteer de volgende hunk |  |
| `` v `` | Toggle drag selecteer |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <ctrl+o> `` | Copy selected text to clipboard |  |
| `` <space> `` | Toggle staged | Toggle lijnen staged / unstaged |
| `` d `` | Verwijdert change (git reset) | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | Open bestand | Open bestand in standaardapplicatie. |
| `` e `` | Verander bestand | Open bestand in externe editor. |
| `` <esc> `` | Ga terug naar het bestanden paneel |  |
| `` <tab> `` | Ga naar een ander paneel | Switch to other view (staged/unstaged changes). |
| `` E `` | Edit hunk | Edit selected hunk in external editor. |
| `` c `` | Commit veranderingen | Commit gestagede wijzigingen. |
| `` w `` | Commit veranderingen zonder pre-commit hook |  |
| `` C `` | Commit veranderingen met de git editor |  |
| `` <ctrl+f> `` | Find base commit for fixup | Vind de commit waar je huidige wijzigingen bovenop zijn gebouwd met als doel die commit te amenden/fixen. Hierdoor hoef je dit niet met de hand te doen. Zie: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | Start met zoeken |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Toepassen | Apply the stash entry to your working directory. |
| `` g `` | Pop | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | Laten vallen | Remove the stash entry from the stash list. |
| `` n `` | Nieuwe branch | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` w `` | New worktree |  |
| `` r `` | Hernoem stash |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Bekijk gecommite bestanden |  |
| `` / `` | Filter the current view by text |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` e `` | Verander config bestand | Open bestand in externe editor. |
| `` u `` | Check voor updates |  |
| `` <enter> `` | Wissel naar een recente repo |  |
| `` a `` | Show/cycle all branch logs |  |
| `` A `` | Show/cycle all branch logs (reverse) |  |
| `` 0 `` | Focus main view |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy abbreviated commit hash to clipboard |  |
| `` <space> `` | Uitchecken | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Creëer nieuwe branch van commit |  |
| `` N `` | Verplaats commits naar nieuwe branch | Maak een nieuwe branch en verplaats niet-gepushte commits van de huidige branch hier naar toe. Gebruik dit in het geval dat je deze commits eigenlijk op een nieuwe branch had willen maken.<br><br>Let op dat de selectie genegeerd wordt. De nieuwe branch komt ofwel bovenop de main branch, of bovenop de huidige branch (je kan kiezen). |
| `` w `` | New worktree |  |
| `` g `` | Bekijk reset opties | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Kopieer commit (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <ctrl+r> `` | Reset cherry-picked (gekopieerde) commits selectie |  |
| `` <ctrl+t> `` | Open externe diff applicatie (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Bekijk gecommite bestanden |  |
| `` / `` | Start met zoeken |  |

## Submodules

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Kopieer submodule naam naar klembord |  |
| `` <enter> `` | Enter | Enter submodule |
| `` d `` | Verwijderen | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | Update selected submodule. |
| `` n `` | Voeg nieuwe submodule toe |  |
| `` e `` | Update submodule URL |  |
| `` i `` | Initialize | Initialiseer submodule |
| `` b `` | Bekijk bulk submodule opties |  |
| `` / `` | Filter the current view by text |  |

## Tags

| Key | Action | Info |
|-----|--------|-------------|
| `` <ctrl+o> `` | Copy tag to clipboard |  |
| `` <space> `` | Uitchecken | Geselecteerde tag uitchecken als detached HEAD. |
| `` n `` | Creëer tag | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` w `` | New worktree |  |
| `` d `` | Verwijderen | View delete options for local/remote tag. |
| `` P `` | Push tag | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Resetten | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <ctrl+t> `` | Open externe diff applicatie (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Bekijk commits |  |
| `` / `` | Filter the current view by text |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | Openen in editor |  |
| `` d `` | Verwijderen | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | Filter the current view by text |  |
