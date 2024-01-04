_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Globalne

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Switch to a recent repo |  |
| `` <pgup> (fn+up/shift+k) `` | Scroll up main panel |  |
| `` <pgdown> (fn+down/shift+j) `` | Scroll down main panel |  |
| `` @ `` | Open command log menu |  |
| `` } `` | Increase the size of the context shown around changes in the diff view |  |
| `` { `` | Decrease the size of the context shown around changes in the diff view |  |
| `` : `` | Wykonaj własną komendę |  |
| `` <c-p> `` | View custom patch options |  |
| `` m `` | Widok scalenia/opcje zmiany bazy |  |
| `` R `` | Odśwież |  |
| `` + `` | Next screen mode (normal/half/fullscreen) |  |
| `` _ `` | Prev screen mode |  |
| `` ? `` | Open menu |  |
| `` <c-s> `` | View filter-by-path options |  |
| `` W `` | Open diff menu |  |
| `` <c-e> `` | Open diff menu |  |
| `` <c-w> `` | Toggle whether or not whitespace changes are shown in the diff view |  |
| `` z `` | Undo | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | Redo | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` P `` | Push |  |
| `` p `` | Pull |  |

## List panel navigation

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Previous page |  |
| `` . `` | Next page |  |
| `` < `` | Scroll to top |  |
| `` > `` | Scroll to bottom |  |
| `` v `` | Toggle range select |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | Search the current view by text |  |
| `` H `` | Scroll left |  |
| `` L `` | Scroll right |  |
| `` ] `` | Next tab |  |
| `` [ `` | Previous tab |  |

## Commit summary

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Potwierdź |  |
| `` <esc> `` | Zamknij |  |

## Commity

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit SHA to clipboard |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` b `` | View bisect options |  |
| `` s `` | Ściśnij |  |
| `` f `` | Napraw commit |  |
| `` r `` | Zmień nazwę commita |  |
| `` R `` | Zmień nazwę commita w edytorze |  |
| `` d `` | Usuń commit |  |
| `` e `` | Edytuj commit |  |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Wybierz commit (podczas zmiany bazy) |  |
| `` F `` | Utwórz commit naprawczy dla tego commita |  |
| `` S `` | Spłaszcz wszystkie commity naprawcze powyżej zaznaczonych commitów (autosquash) |  |
| `` <c-j> `` | Przenieś commit 1 w dół |  |
| `` <c-k> `` | Przenieś commit 1 w górę |  |
| `` V `` | Wklej commity (przebieranie) |  |
| `` B `` | Mark commit as base commit for rebase | Select a base commit for the next rebase; this will effectively perform a 'git rebase --onto'. |
| `` A `` | Popraw commit zmianami z poczekalni |  |
| `` a `` | Set/Reset commit author |  |
| `` t `` | Odwróć commit |  |
| `` T `` | Tag commit |  |
| `` <c-l> `` | Open log menu |  |
| `` w `` | View worktree options |  |
| `` <space> `` | Checkout commit |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Wyświetl opcje resetu |  |
| `` C `` | Kopiuj commit (przebieranie) |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Przeglądaj pliki commita |  |
| `` / `` | Search the current view by text |  |

## Confirmation panel

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Potwierdź |  |
| `` <esc> `` | Zamknij |  |

## Local branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` i `` | Show git-flow options |  |
| `` <space> `` | Przełącz |  |
| `` n `` | Nowa gałąź |  |
| `` o `` | Utwórz żądanie pobrania |  |
| `` O `` | Utwórz opcje żądania ściągnięcia |  |
| `` <c-y> `` | Skopiuj adres URL żądania pobrania do schowka |  |
| `` c `` | Przełącz używając nazwy |  |
| `` F `` | Wymuś przełączenie |  |
| `` d `` | View delete options |  |
| `` r `` | Zmiana bazy gałęzi |  |
| `` M `` | Scal do obecnej gałęzi |  |
| `` f `` | Fast-forward this branch from its upstream |  |
| `` T `` | Create tag |  |
| `` s `` | Sort order |  |
| `` g `` | Wyświetl opcje resetu |  |
| `` R `` | Rename branch |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream |
| `` w `` | View worktree options |  |
| `` <enter> `` | View commits |  |
| `` / `` | Filter the current view by text |  |

## Main panel (patch building)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Poprzedni kawałek |  |
| `` <right> `` | Następny kawałek |  |
| `` v `` | Toggle range select |  |
| `` a `` | Toggle select hunk |  |
| `` <c-o> `` | Copy the selected text to the clipboard |  |
| `` o `` | Otwórz plik |  |
| `` e `` | Edytuj plik |  |
| `` <space> `` | Add/Remove line(s) to patch |  |
| `` <esc> `` | Wyście z trybu "linia po linii" |  |
| `` / `` | Search the current view by text |  |

## Menu

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Wykonaj |  |
| `` <esc> `` | Zamknij |  |
| `` / `` | Filter the current view by text |  |

## Pliki

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy the file name to the clipboard |  |
| `` <space> `` | Przełącz stan poczekalni |  |
| `` <c-b> `` | Filter files by status |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Zatwierdź zmiany |  |
| `` w `` | Zatwierdź zmiany bez skryptu pre-commit |  |
| `` A `` | Zmień ostatni commit |  |
| `` C `` | Zatwierdź zmiany używając edytora |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edytuj plik |  |
| `` o `` | Otwórz plik |  |
| `` i `` | Ignore or exclude file |  |
| `` r `` | Odśwież pliki |  |
| `` s `` | Przechowaj zmiany |  |
| `` S `` | Wyświetl opcje schowka |  |
| `` a `` | Przełącz stan poczekalni wszystkich |  |
| `` <enter> `` | Zatwierdź pojedyncze linie |  |
| `` d `` | Pokaż opcje porzucania zmian |  |
| `` g `` | View upstream reset options |  |
| `` D `` | Wyświetl opcje resetu |  |
| `` ` `` | Toggle file tree view |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Open external merge tool (git mergetool) |  |
| `` f `` | Pobierz |  |
| `` / `` | Search the current view by text |  |

## Pliki commita

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy the committed file name to the clipboard |  |
| `` c `` | Plik wybierania |  |
| `` d `` | Porzuć zmiany commita dla tego pliku |  |
| `` o `` | Otwórz plik |  |
| `` e `` | Edytuj plik |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle file included in patch |  |
| `` a `` | Toggle all files included in patch |  |
| `` <enter> `` | Enter file to add selected lines to the patch (or toggle directory collapsed) |  |
| `` ` `` | Toggle file tree view |  |
| `` / `` | Search the current view by text |  |

## Poczekalnia

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Poprzedni kawałek |  |
| `` <right> `` | Następny kawałek |  |
| `` v `` | Toggle range select |  |
| `` a `` | Toggle select hunk |  |
| `` <c-o> `` | Copy the selected text to the clipboard |  |
| `` o `` | Otwórz plik |  |
| `` e `` | Edytuj plik |  |
| `` <esc> `` | Wróć do panelu plików |  |
| `` <tab> `` | Switch to other panel (staged/unstaged changes) |  |
| `` <space> `` | Toggle line staged / unstaged |  |
| `` d `` | Discard change (git reset) |  |
| `` E `` | Edit hunk |  |
| `` c `` | Zatwierdź zmiany |  |
| `` w `` | Zatwierdź zmiany bez skryptu pre-commit |  |
| `` C `` | Zatwierdź zmiany używając edytora |  |
| `` / `` | Search the current view by text |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit SHA to clipboard |  |
| `` w `` | View worktree options |  |
| `` <space> `` | Checkout commit |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Wyświetl opcje resetu |  |
| `` C `` | Kopiuj commit (przebieranie) |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` / `` | Filter the current view by text |  |

## Remote branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` <space> `` | Przełącz |  |
| `` n `` | Nowa gałąź |  |
| `` M `` | Scal do obecnej gałęzi |  |
| `` r `` | Zmiana bazy gałęzi |  |
| `` d `` | Delete remote tag |  |
| `` u `` | Set as upstream of checked-out branch |  |
| `` s `` | Sort order |  |
| `` g `` | Wyświetl opcje resetu |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | View commits |  |
| `` / `` | Filter the current view by text |  |

## Remotes

| Key | Action | Info |
|-----|--------|-------------|
| `` f `` | Fetch remote |  |
| `` n `` | Add new remote |  |
| `` d `` | Remove remote |  |
| `` e `` | Edit remote |  |
| `` / `` | Filter the current view by text |  |

## Scalanie

| Key | Action | Info |
|-----|--------|-------------|
| `` e `` | Edytuj plik |  |
| `` o `` | Otwórz plik |  |
| `` <left> `` | Poprzedni konflikt |  |
| `` <right> `` | Następny konflikt |  |
| `` <up> `` | Wybierz poprzedni kawałek |  |
| `` <down> `` | Wybierz następny kawałek |  |
| `` z `` | Cofnij |  |
| `` M `` | Open external merge tool (git mergetool) |  |
| `` <space> `` | Wybierz kawałek |  |
| `` b `` | Wybierz oba kawałki |  |
| `` <esc> `` | Wróć do panelu plików |  |

## Schowek

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Zastosuj |  |
| `` g `` | Wyciągnij |  |
| `` d `` | Porzuć |  |
| `` n `` | Nowa gałąź |  |
| `` r `` | Rename stash |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | Przeglądaj pliki commita |  |
| `` / `` | Filter the current view by text |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Otwórz konfigurację |  |
| `` e `` | Edytuj konfigurację |  |
| `` u `` | Sprawdź aktualizacje |  |
| `` <enter> `` | Switch to a recent repo |  |
| `` a `` | Pokaż wszystkie logi gałęzi |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit SHA to clipboard |  |
| `` w `` | View worktree options |  |
| `` <space> `` | Checkout commit |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Wyświetl opcje resetu |  |
| `` C `` | Kopiuj commit (przebieranie) |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Przeglądaj pliki commita |  |
| `` / `` | Search the current view by text |  |

## Submodules

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy submodule name to clipboard |  |
| `` <enter> `` | Enter submodule |  |
| `` <space> `` | Enter submodule |  |
| `` d `` | Remove submodule |  |
| `` u `` | Update submodule |  |
| `` n `` | Add new submodule |  |
| `` e `` | Update submodule URL |  |
| `` i `` | Initialize submodule |  |
| `` b `` | View bulk submodule options |  |
| `` / `` | Filter the current view by text |  |

## Tags

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Przełącz |  |
| `` d `` | View delete options |  |
| `` P `` | Push tag |  |
| `` n `` | Create tag |  |
| `` g `` | Wyświetl opcje resetu |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | View commits |  |
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

## Zwykłe

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | Przewiń w dół |  |
| `` mouse wheel up (fn+down) `` | Przewiń w górę |  |
