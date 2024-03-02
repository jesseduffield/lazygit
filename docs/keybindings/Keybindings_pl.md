_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Globalne

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Switch to a recent repo |  |
| `` <pgup> (fn+up/shift+k) `` | Scroll up main window |  |
| `` <pgdown> (fn+down/shift+j) `` | Scroll down main window |  |
| `` @ `` | View command log options | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | Push | Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` p `` | Pull | Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` } `` | Increase diff context size | Increase the amount of the context shown around changes in the diff view. |
| `` { `` | Decrease diff context size | Decrease the amount of the context shown around changes in the diff view. |
| `` : `` | Wykonaj własną komendę | Bring up a prompt where you can enter a shell command to execute. Not to be confused with pre-configured custom commands. |
| `` <c-p> `` | View custom patch options |  |
| `` m `` | Widok scalenia/opcje zmiany bazy | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | Odśwież | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | Next screen mode (normal/half/fullscreen) |  |
| `` _ `` | Prev screen mode |  |
| `` ? `` | Open keybindings menu |  |
| `` <c-s> `` | View filter options | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | View diffing options | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | View diffing options | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | Quit |  |
| `` <esc> `` | Anuluj |  |
| `` <c-w> `` | Toggle whitespace | Toggle whether or not whitespace changes are shown in the diff view. |
| `` z `` | Undo | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | Redo | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |

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
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` b `` | View bisect options |  |
| `` s `` | Spłaszcz | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | Napraw | Meld the selected commit into the commit below it. Similar to fixup, but the selected commit's message will be discarded. |
| `` r `` | Zmień nazwę commita | Reword the selected commit's message. |
| `` R `` | Zmień nazwę commita w edytorze |  |
| `` d `` | Usuń commit | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | Edit (start interactive rebase) | Edytuj commit |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick | Wybierz commit (podczas zmiany bazy) |
| `` F `` | Create fixup commit | Utwórz commit naprawczy dla tego commita |
| `` S `` | Apply fixup commits | Spłaszcz wszystkie commity naprawcze powyżej zaznaczonych commitów (autosquash) |
| `` <c-j> `` | Przenieś commit 1 w dół |  |
| `` <c-k> `` | Przenieś commit 1 w górę |  |
| `` V `` | Wklej commity (przebieranie) |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Amend | Popraw commit zmianami z poczekalni |
| `` a `` | Amend commit attribute | Set/Reset commit author or set co-author. |
| `` t `` | Revert | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | Tag commit | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <c-l> `` | View log options | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` <space> `` | Przełącz | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Wyświetl opcje resetu | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Kopiuj commit (przebieranie) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Przeglądaj pliki commita |  |
| `` w `` | View worktree options |  |
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
| `` <space> `` | Przełącz | Checkout selected item. |
| `` n `` | Nowa gałąź |  |
| `` o `` | Utwórz żądanie pobrania |  |
| `` O `` | Utwórz opcje żądania ściągnięcia |  |
| `` <c-y> `` | Skopiuj adres URL żądania pobrania do schowka |  |
| `` c `` | Przełącz używając nazwy | Checkout by name. In the input box you can enter '-' to switch to the last branch. |
| `` F `` | Wymuś przełączenie | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | Delete | View delete options for local/remote branch. |
| `` r `` | Zmiana bazy gałęzi | Rebase the checked-out branch onto the selected branch. |
| `` M `` | Scal do obecnej gałęzi | Merge selected branch into currently checked out branch. |
| `` f `` | Fast-forward | Fast-forward selected branch from its upstream. |
| `` T `` | New tag |  |
| `` s `` | Sort order |  |
| `` g `` | Wyświetl opcje resetu |  |
| `` R `` | Rename branch |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Main panel (patch building)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Poprzedni kawałek |  |
| `` <right> `` | Następny kawałek |  |
| `` v `` | Toggle range select |  |
| `` a `` | Select hunk | Toggle hunk selection mode. |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` o `` | Otwórz plik | Open file in default application. |
| `` e `` | Edytuj plik | Open file in external editor. |
| `` <space> `` | Toggle lines in patch |  |
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
| `` <c-o> `` | Copy path to clipboard |  |
| `` <space> `` | Przełącz stan poczekalni | Toggle staged for selected file. |
| `` <c-b> `` | Filter files by status |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Zatwierdź zmiany | Commit staged changes. |
| `` w `` | Zatwierdź zmiany bez skryptu pre-commit |  |
| `` A `` | Zmień ostatni commit |  |
| `` C `` | Zatwierdź zmiany używając edytora |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit | Open file in external editor. |
| `` o `` | Otwórz plik | Open file in default application. |
| `` i `` | Ignore or exclude file |  |
| `` r `` | Odśwież pliki |  |
| `` s `` | Stash | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | Wyświetl opcje schowka | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | Przełącz stan poczekalni wszystkich | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | Zatwierdź pojedyncze linie | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | Pokaż opcje porzucania zmian | View options for discarding changes to the selected file. |
| `` g `` | View upstream reset options |  |
| `` D `` | Reset | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | Toggle file tree view | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Open external merge tool | Run `git mergetool`. |
| `` f `` | Pobierz | Fetch changes from remote. |
| `` / `` | Search the current view by text |  |

## Pliki commita

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy path to clipboard |  |
| `` c `` | Przełącz | Plik wybierania |
| `` d `` | Remove | Porzuć zmiany commita dla tego pliku |
| `` o `` | Otwórz plik | Open file in default application. |
| `` e `` | Edit | Open file in external editor. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle file included in patch | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Toggle all files | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Enter file / Toggle directory collapsed | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | Toggle file tree view | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` / `` | Search the current view by text |  |

## Poczekalnia

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Poprzedni kawałek |  |
| `` <right> `` | Następny kawałek |  |
| `` v `` | Toggle range select |  |
| `` a `` | Select hunk | Toggle hunk selection mode. |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` <space> `` | Przełącz stan poczekalni | Toggle selection staged / unstaged. |
| `` d `` | Discard | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | Otwórz plik | Open file in default application. |
| `` e `` | Edytuj plik | Open file in external editor. |
| `` <esc> `` | Wróć do panelu plików |  |
| `` <tab> `` | Switch view | Switch to other view (staged/unstaged changes). |
| `` E `` | Edit hunk | Edit selected hunk in external editor. |
| `` c `` | Zatwierdź zmiany | Commit staged changes. |
| `` w `` | Zatwierdź zmiany bez skryptu pre-commit |  |
| `` C `` | Zatwierdź zmiany używając edytora |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | Search the current view by text |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit SHA to clipboard |  |
| `` <space> `` | Przełącz | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Wyświetl opcje resetu | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Kopiuj commit (przebieranie) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Remote branches

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` <space> `` | Przełącz | Checkout a new local branch based on the selected remote branch. The new branch will track the remote branch. |
| `` n `` | Nowa gałąź |  |
| `` M `` | Scal do obecnej gałęzi | Merge selected branch into currently checked out branch. |
| `` r `` | Zmiana bazy gałęzi | Rebase the checked-out branch onto the selected branch. |
| `` d `` | Delete | Delete the remote branch from the remote. |
| `` u `` | Set as upstream | Set the selected remote branch as the upstream of the checked-out branch. |
| `` s `` | Sort order |  |
| `` g `` | Wyświetl opcje resetu | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Remotes

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | New remote |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Edit | Edit the selected remote's name or URL. |
| `` f `` | Pobierz | Fetch updates from the remote repository. This retrieves new commits and branches without merging them into your local branches. |
| `` / `` | Filter the current view by text |  |

## Scalanie

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Wybierz kawałek |  |
| `` b `` | Wybierz oba kawałki |  |
| `` <up> `` | Wybierz poprzedni kawałek |  |
| `` <down> `` | Wybierz następny kawałek |  |
| `` <left> `` | Poprzedni konflikt |  |
| `` <right> `` | Następny konflikt |  |
| `` z `` | Cofnij | Undo last merge conflict resolution. |
| `` e `` | Edytuj plik | Open file in external editor. |
| `` o `` | Otwórz plik | Open file in default application. |
| `` M `` | Open external merge tool | Run `git mergetool`. |
| `` <esc> `` | Wróć do panelu plików |  |

## Schowek

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Zastosuj | Apply the stash entry to your working directory. |
| `` g `` | Wyciągnij | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | Porzuć | Remove the stash entry from the stash list. |
| `` n `` | Nowa gałąź | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | Rename stash |  |
| `` <enter> `` | Przeglądaj pliki commita |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Otwórz konfigurację | Open file in default application. |
| `` e `` | Edytuj konfigurację | Open file in external editor. |
| `` u `` | Sprawdź aktualizacje |  |
| `` <enter> `` | Switch to a recent repo |  |
| `` a `` | Pokaż wszystkie logi gałęzi |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit SHA to clipboard |  |
| `` <space> `` | Przełącz | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Wyświetl opcje resetu | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Kopiuj commit (przebieranie) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Przeglądaj pliki commita |  |
| `` w `` | View worktree options |  |
| `` / `` | Search the current view by text |  |

## Submodules

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy submodule name to clipboard |  |
| `` <enter> `` | Enter | Enter submodule. After entering the submodule, you can press `<esc>` to escape back to the parent repo. |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | Update selected submodule. |
| `` n `` | New submodule |  |
| `` e `` | Update submodule URL |  |
| `` i `` | Initialize | Initialize the selected submodule to prepare for fetching. You probably want to follow this up by invoking the 'update' action to fetch the submodule. |
| `` b `` | View bulk submodule options |  |
| `` / `` | Filter the current view by text |  |

## Tags

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Przełącz | Checkout the selected tag tag as a detached HEAD. |
| `` n `` | New tag | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Delete | View delete options for local/remote tag. |
| `` P `` | Push tag | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <enter> `` | View commits |  |
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

## Zwykłe

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | Przewiń w dół |  |
| `` mouse wheel up (fn+down) `` | Przewiń w górę |  |
