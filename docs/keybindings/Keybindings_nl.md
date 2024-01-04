_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Sneltoetsen

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Globale sneltoetsen

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Wissel naar een recente repo |  |
| `` <pgup> (fn+up/shift+k) `` | Scroll naar beneden vanaf hoofdpaneel |  |
| `` <pgdown> (fn+down/shift+j) `` | Scroll naar beneden vanaf hoofdpaneel |  |
| `` @ `` | Open command log menu |  |
| `` } `` | Increase the size of the context shown around changes in the diff view |  |
| `` { `` | Decrease the size of the context shown around changes in the diff view |  |
| `` : `` | Voer aangepaste commando uit |  |
| `` <c-p> `` | Bekijk aangepaste patch opties |  |
| `` m `` | Bekijk merge/rebase opties |  |
| `` R `` | Verversen |  |
| `` + `` | Volgende scherm modus (normaal/half/groot) |  |
| `` _ `` | Vorige scherm modus |  |
| `` ? `` | Open menu |  |
| `` <c-s> `` | Bekijk scoping opties |  |
| `` W `` | Open diff menu |  |
| `` <c-e> `` | Open diff menu |  |
| `` <c-w> `` | Toggle whether or not whitespace changes are shown in the diff view |  |
| `` z `` | Ongedaan maken (via reflog) (experimenteel) | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | Redo (via reflog) (experimenteel) | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` P `` | Push |  |
| `` p `` | Pull |  |

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
| `` <space> `` | Toggle staged |  |
| `` <c-b> `` | Filter files by status |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Commit veranderingen |  |
| `` w `` | Commit veranderingen zonder pre-commit hook |  |
| `` A `` | Wijzig laatste commit |  |
| `` C `` | Commit veranderingen met de git editor |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Verander bestand |  |
| `` o `` | Open bestand |  |
| `` i `` | Ignore or exclude file |  |
| `` r `` | Refresh bestanden |  |
| `` s `` | Stash-bestanden |  |
| `` S `` | Bekijk stash opties |  |
| `` a `` | Toggle staged alle |  |
| `` <enter> `` | Stage individuele hunks/lijnen |  |
| `` d `` | Bekijk 'veranderingen ongedaan maken' opties |  |
| `` g `` | Bekijk upstream reset opties |  |
| `` D `` | Bekijk reset opties |  |
| `` ` `` | Toggle bestandsboom weergave |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Open external merge tool (git mergetool) |  |
| `` f `` | Fetch |  |
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
| `` <space> `` | Uitchecken |  |
| `` n `` | Nieuwe branch |  |
| `` o `` | Maak een pull-request |  |
| `` O `` | Bekijk opties voor pull-aanvraag |  |
| `` <c-y> `` | Kopieer de URL van het pull-verzoek naar het klembord |  |
| `` c `` | Uitchecken bij naam |  |
| `` F `` | Forceer checkout |  |
| `` d `` | View delete options |  |
| `` r `` | Rebase branch |  |
| `` M `` | Merge in met huidige checked out branch |  |
| `` f `` | Fast-forward deze branch vanaf zijn upstream |  |
| `` T `` | Creëer tag |  |
| `` s `` | Sort order |  |
| `` g `` | Bekijk reset opties |  |
| `` R `` | Hernoem branch |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream |
| `` w `` | View worktree options |  |
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
| `` <c-o> `` | Kopieer de vastgelegde bestandsnaam naar het klembord |  |
| `` c `` | Bestand uitchecken |  |
| `` d `` | Uitsluit deze commit zijn veranderingen aan dit bestand |  |
| `` o `` | Open bestand |  |
| `` e `` | Verander bestand |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle bestand inbegrepen in patch |  |
| `` a `` | Toggle all files included in patch |  |
| `` <enter> `` | Enter bestand om geselecteerde regels toe te voegen aan de patch |  |
| `` ` `` | Toggle bestandsboom weergave |  |
| `` / `` | Start met zoeken |  |

## Commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer commit SHA naar klembord |  |
| `` <c-r> `` | Reset cherry-picked (gekopieerde) commits selectie |  |
| `` b `` | View bisect options |  |
| `` s `` | Squash beneden |  |
| `` f `` | Fixup commit |  |
| `` r `` | Hernoem commit |  |
| `` R `` | Hernoem commit met editor |  |
| `` d `` | Verwijder commit |  |
| `` e `` | Wijzig commit |  |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Kies commit (wanneer midden in rebase) |  |
| `` F `` | Creëer fixup commit |  |
| `` S `` | Squash bovenstaande commits |  |
| `` <c-j> `` | Verplaats commit 1 naar beneden |  |
| `` <c-k> `` | Verplaats commit 1 naar boven |  |
| `` V `` | Plak commits (cherry-pick) |  |
| `` B `` | Mark commit as base commit for rebase | Select a base commit for the next rebase; this will effectively perform a 'git rebase --onto'. |
| `` A `` | Wijzig commit met staged veranderingen |  |
| `` a `` | Set/Reset commit author |  |
| `` t `` | Commit ongedaan maken |  |
| `` T `` | Tag commit |  |
| `` <c-l> `` | Open log menu |  |
| `` w `` | View worktree options |  |
| `` <space> `` | Checkout commit |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | Open commit in browser |  |
| `` n `` | Creëer nieuwe branch van commit |  |
| `` g `` | Bekijk reset opties |  |
| `` C `` | Kopieer commit (cherry-pick) |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Bekijk gecommite bestanden |  |
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
| `` e `` | Verander bestand |  |
| `` o `` | Open bestand |  |
| `` <left> `` | Selecteer voorgaand conflict |  |
| `` <right> `` | Selecteer volgende conflict |  |
| `` <up> `` | Selecteer bovenste hunk |  |
| `` <down> `` | Selecteer onderste hunk |  |
| `` z `` | Ongedaan maken |  |
| `` M `` | Open external merge tool (git mergetool) |  |
| `` <space> `` | Kies stuk |  |
| `` b `` | Kies beide stukken |  |
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
| `` a `` | Toggle selecteer hunk |  |
| `` <c-o> `` | Copy the selected text to the clipboard |  |
| `` o `` | Open bestand |  |
| `` e `` | Verander bestand |  |
| `` <space> `` | Voeg toe/verwijder lijn(en) in patch |  |
| `` <esc> `` | Sluit lijn-bij-lijn modus |  |
| `` / `` | Start met zoeken |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer commit SHA naar klembord |  |
| `` w `` | View worktree options |  |
| `` <space> `` | Checkout commit |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | Open commit in browser |  |
| `` n `` | Creëer nieuwe branch van commit |  |
| `` g `` | Bekijk reset opties |  |
| `` C `` | Kopieer commit (cherry-pick) |  |
| `` <c-r> `` | Reset cherry-picked (gekopieerde) commits selectie |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Bekijk commits |  |
| `` / `` | Filter the current view by text |  |

## Remote branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer branch name naar klembord |  |
| `` <space> `` | Uitchecken |  |
| `` n `` | Nieuwe branch |  |
| `` M `` | Merge in met huidige checked out branch |  |
| `` r `` | Rebase branch |  |
| `` d `` | Delete remote tag |  |
| `` u `` | Stel in als upstream van uitgecheckte branch |  |
| `` s `` | Sort order |  |
| `` g `` | Bekijk reset opties |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | Bekijk commits |  |
| `` / `` | Filter the current view by text |  |

## Remotes

| Key | Action | Info |
|-----|--------|-------------|
| `` f `` | Fetch remote |  |
| `` n `` | Voeg een nieuwe remote toe |  |
| `` d `` | Verwijder remote |  |
| `` e `` | Wijzig remote |  |
| `` / `` | Filter the current view by text |  |

## Staging

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Selecteer de vorige hunk |  |
| `` <right> `` | Selecteer de volgende hunk |  |
| `` v `` | Toggle drag selecteer |  |
| `` a `` | Toggle selecteer hunk |  |
| `` <c-o> `` | Copy the selected text to the clipboard |  |
| `` o `` | Open bestand |  |
| `` e `` | Verander bestand |  |
| `` <esc> `` | Ga terug naar het bestanden paneel |  |
| `` <tab> `` | Ga naar een ander paneel |  |
| `` <space> `` | Toggle lijnen staged / unstaged |  |
| `` d `` | Verwijdert change (git reset) |  |
| `` E `` | Edit hunk |  |
| `` c `` | Commit veranderingen |  |
| `` w `` | Commit veranderingen zonder pre-commit hook |  |
| `` C `` | Commit veranderingen met de git editor |  |
| `` / `` | Start met zoeken |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Toepassen |  |
| `` g `` | Pop |  |
| `` d `` | Laten vallen |  |
| `` n `` | Nieuwe branch |  |
| `` r `` | Rename stash |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | Bekijk gecommite bestanden |  |
| `` / `` | Filter the current view by text |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Open config bestand |  |
| `` e `` | Verander config bestand |  |
| `` u `` | Check voor updates |  |
| `` <enter> `` | Wissel naar een recente repo |  |
| `` a `` | Alle logs van de branch laten zien |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer commit SHA naar klembord |  |
| `` w `` | View worktree options |  |
| `` <space> `` | Checkout commit |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | Open commit in browser |  |
| `` n `` | Creëer nieuwe branch van commit |  |
| `` g `` | Bekijk reset opties |  |
| `` C `` | Kopieer commit (cherry-pick) |  |
| `` <c-r> `` | Reset cherry-picked (gekopieerde) commits selectie |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Bekijk gecommite bestanden |  |
| `` / `` | Start met zoeken |  |

## Submodules

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopieer submodule naam naar klembord |  |
| `` <enter> `` | Enter submodule |  |
| `` <space> `` | Enter submodule |  |
| `` d `` | Remove submodule |  |
| `` u `` | Update submodule |  |
| `` n `` | Voeg nieuwe submodule toe |  |
| `` e `` | Update submodule URL |  |
| `` i `` | Initialiseer submodule |  |
| `` b `` | Bekijk bulk submodule opties |  |
| `` / `` | Filter the current view by text |  |

## Tags

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Uitchecken |  |
| `` d `` | View delete options |  |
| `` P `` | Push tag |  |
| `` n `` | Creëer tag |  |
| `` g `` | Bekijk reset opties |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | Bekijk commits |  |
| `` / `` | Filter the current view by text |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | Create worktree |  |
| `` <space> `` | Switch to worktree |  |
| `` <enter> `` | Switch to worktree |  |
| `` o `` | Open in editor |  |
| `` d `` | Remove worktree |  |
| `` / `` | Filter the current view by text |  |
