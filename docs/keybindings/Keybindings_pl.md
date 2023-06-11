_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Globalne

<pre>
  <kbd>&lt;c-r&gt;</kbd>: Switch to a recent repo
  <kbd>&lt;pgup&gt;</kbd>: Scroll up main panel (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: Scroll down main panel (fn+down/shift+j)
  <kbd>@</kbd>: Open command log menu
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: Wykonaj własną komendę
  <kbd>&lt;c-p&gt;</kbd>: View custom patch options
  <kbd>m</kbd>: Widok scalenia/opcje zmiany bazy
  <kbd>R</kbd>: Odśwież
  <kbd>+</kbd>: Next screen mode (normal/half/fullscreen)
  <kbd>_</kbd>: Prev screen mode
  <kbd>?</kbd>: Open menu
  <kbd>&lt;c-s&gt;</kbd>: View filter-by-path options
  <kbd>W</kbd>: Open diff menu
  <kbd>&lt;c-e&gt;</kbd>: Open diff menu
  <kbd>&lt;c-w&gt;</kbd>: Toggle whether or not whitespace changes are shown in the diff view
  <kbd>z</kbd>: Undo
  <kbd>&lt;c-z&gt;</kbd>: Redo
  <kbd>P</kbd>: Push
  <kbd>p</kbd>: Pull
</pre>

## List panel navigation

<pre>
  <kbd>,</kbd>: Previous page
  <kbd>.</kbd>: Next page
  <kbd>&lt;</kbd>: Scroll to top
  <kbd>&gt;</kbd>: Scroll to bottom
  <kbd>/</kbd>: Search the current view by text
  <kbd>H</kbd>: Scroll left
  <kbd>L</kbd>: Scroll right
  <kbd>]</kbd>: Next tab
  <kbd>[</kbd>: Previous tab
</pre>

## Commit summary

<pre>
  <kbd>&lt;enter&gt;</kbd>: Potwierdź
  <kbd>&lt;esc&gt;</kbd>: Zamknij
</pre>

## Commity

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy commit SHA to clipboard
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>b</kbd>: View bisect options
  <kbd>s</kbd>: Ściśnij
  <kbd>f</kbd>: Napraw commit
  <kbd>r</kbd>: Zmień nazwę commita
  <kbd>R</kbd>: Zmień nazwę commita w edytorze
  <kbd>d</kbd>: Usuń commit
  <kbd>e</kbd>: Edytuj commit
  <kbd>p</kbd>: Wybierz commit (podczas zmiany bazy)
  <kbd>F</kbd>: Utwórz commit naprawczy dla tego commita
  <kbd>S</kbd>: Spłaszcz wszystkie commity naprawcze powyżej zaznaczonych commitów (autosquash)
  <kbd>&lt;c-j&gt;</kbd>: Przenieś commit 1 w dół
  <kbd>&lt;c-k&gt;</kbd>: Przenieś commit 1 w górę
  <kbd>v</kbd>: Wklej commity (przebieranie)
  <kbd>B</kbd>: Mark commit as base commit for rebase
  <kbd>A</kbd>: Popraw commit zmianami z poczekalni
  <kbd>a</kbd>: Set/Reset commit author
  <kbd>t</kbd>: Odwróć commit
  <kbd>T</kbd>: Tag commit
  <kbd>&lt;c-l&gt;</kbd>: Open log menu
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Checkout commit
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: Open commit in browser
  <kbd>n</kbd>: Create new branch off of commit
  <kbd>g</kbd>: Wyświetl opcje resetu
  <kbd>c</kbd>: Kopiuj commit (przebieranie)
  <kbd>C</kbd>: Kopiuj zakres commitów (przebieranie)
  <kbd>&lt;enter&gt;</kbd>: Przeglądaj pliki commita
  <kbd>/</kbd>: Search the current view by text
</pre>

## Confirmation panel

<pre>
  <kbd>&lt;enter&gt;</kbd>: Potwierdź
  <kbd>&lt;esc&gt;</kbd>: Zamknij
</pre>

## Local branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy branch name to clipboard
  <kbd>i</kbd>: Show git-flow options
  <kbd>&lt;space&gt;</kbd>: Przełącz
  <kbd>n</kbd>: Nowa gałąź
  <kbd>o</kbd>: Utwórz żądanie pobrania
  <kbd>O</kbd>: Utwórz opcje żądania ściągnięcia
  <kbd>&lt;c-y&gt;</kbd>: Skopiuj adres URL żądania pobrania do schowka
  <kbd>c</kbd>: Przełącz używając nazwy
  <kbd>F</kbd>: Wymuś przełączenie
  <kbd>d</kbd>: Usuń gałąź
  <kbd>r</kbd>: Zmiana bazy gałęzi
  <kbd>M</kbd>: Scal do obecnej gałęzi
  <kbd>f</kbd>: Fast-forward this branch from its upstream
  <kbd>T</kbd>: Create tag
  <kbd>g</kbd>: Wyświetl opcje resetu
  <kbd>R</kbd>: Rename branch
  <kbd>u</kbd>: Set/Unset upstream
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: View commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Main panel (patch building)

<pre>
  <kbd>&lt;left&gt;</kbd>: Poprzedni kawałek
  <kbd>&lt;right&gt;</kbd>: Następny kawałek
  <kbd>v</kbd>: Toggle drag select
  <kbd>V</kbd>: Toggle drag select
  <kbd>a</kbd>: Toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: Copy the selected text to the clipboard
  <kbd>o</kbd>: Otwórz plik
  <kbd>e</kbd>: Edytuj plik
  <kbd>&lt;space&gt;</kbd>: Add/Remove line(s) to patch
  <kbd>&lt;esc&gt;</kbd>: Wyście z trybu "linia po linii"
  <kbd>/</kbd>: Search the current view by text
</pre>

## Menu

<pre>
  <kbd>&lt;enter&gt;</kbd>: Wykonaj
  <kbd>&lt;esc&gt;</kbd>: Zamknij
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Pliki

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy the file name to the clipboard
  <kbd>d</kbd>: Pokaż opcje porzucania zmian
  <kbd>&lt;space&gt;</kbd>: Przełącz stan poczekalni
  <kbd>&lt;c-b&gt;</kbd>: Filter files by status
  <kbd>c</kbd>: Zatwierdź zmiany
  <kbd>w</kbd>: Zatwierdź zmiany bez skryptu pre-commit
  <kbd>A</kbd>: Zmień ostatni commit
  <kbd>C</kbd>: Zatwierdź zmiany używając edytora
  <kbd>e</kbd>: Edytuj plik
  <kbd>o</kbd>: Otwórz plik
  <kbd>i</kbd>: Ignore or exclude file
  <kbd>r</kbd>: Odśwież pliki
  <kbd>s</kbd>: Przechowaj zmiany
  <kbd>S</kbd>: Wyświetl opcje schowka
  <kbd>a</kbd>: Przełącz stan poczekalni wszystkich
  <kbd>&lt;enter&gt;</kbd>: Zatwierdź pojedyncze linie
  <kbd>g</kbd>: View upstream reset options
  <kbd>D</kbd>: Wyświetl opcje resetu
  <kbd>`</kbd>: Toggle file tree view
  <kbd>M</kbd>: Open external merge tool (git mergetool)
  <kbd>f</kbd>: Pobierz
  <kbd>/</kbd>: Search the current view by text
</pre>

## Pliki commita

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy the committed file name to the clipboard
  <kbd>c</kbd>: Plik wybierania
  <kbd>d</kbd>: Porzuć zmiany commita dla tego pliku
  <kbd>o</kbd>: Otwórz plik
  <kbd>e</kbd>: Edytuj plik
  <kbd>&lt;space&gt;</kbd>: Toggle file included in patch
  <kbd>a</kbd>: Toggle all files included in patch
  <kbd>&lt;enter&gt;</kbd>: Enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: Toggle file tree view
  <kbd>/</kbd>: Search the current view by text
</pre>

## Poczekalnia

<pre>
  <kbd>&lt;left&gt;</kbd>: Poprzedni kawałek
  <kbd>&lt;right&gt;</kbd>: Następny kawałek
  <kbd>v</kbd>: Toggle drag select
  <kbd>V</kbd>: Toggle drag select
  <kbd>a</kbd>: Toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: Copy the selected text to the clipboard
  <kbd>o</kbd>: Otwórz plik
  <kbd>e</kbd>: Edytuj plik
  <kbd>&lt;esc&gt;</kbd>: Wróć do panelu plików
  <kbd>&lt;tab&gt;</kbd>: Switch to other panel (staged/unstaged changes)
  <kbd>&lt;space&gt;</kbd>: Toggle line staged / unstaged
  <kbd>d</kbd>: Discard change (git reset)
  <kbd>E</kbd>: Edit hunk
  <kbd>c</kbd>: Zatwierdź zmiany
  <kbd>w</kbd>: Zatwierdź zmiany bez skryptu pre-commit
  <kbd>C</kbd>: Zatwierdź zmiany używając edytora
  <kbd>/</kbd>: Search the current view by text
</pre>

## Reflog

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy commit SHA to clipboard
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Checkout commit
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: Open commit in browser
  <kbd>n</kbd>: Create new branch off of commit
  <kbd>g</kbd>: Wyświetl opcje resetu
  <kbd>c</kbd>: Kopiuj commit (przebieranie)
  <kbd>C</kbd>: Kopiuj zakres commitów (przebieranie)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: View commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Remote branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy branch name to clipboard
  <kbd>&lt;space&gt;</kbd>: Przełącz
  <kbd>n</kbd>: Nowa gałąź
  <kbd>M</kbd>: Scal do obecnej gałęzi
  <kbd>r</kbd>: Zmiana bazy gałęzi
  <kbd>d</kbd>: Usuń gałąź
  <kbd>u</kbd>: Set as upstream of checked-out branch
  <kbd>g</kbd>: Wyświetl opcje resetu
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: View commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Remotes

<pre>
  <kbd>f</kbd>: Fetch remote
  <kbd>n</kbd>: Add new remote
  <kbd>d</kbd>: Remove remote
  <kbd>e</kbd>: Edit remote
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Scalanie

<pre>
  <kbd>e</kbd>: Edytuj plik
  <kbd>o</kbd>: Otwórz plik
  <kbd>&lt;left&gt;</kbd>: Poprzedni konflikt
  <kbd>&lt;right&gt;</kbd>: Następny konflikt
  <kbd>&lt;up&gt;</kbd>: Wybierz poprzedni kawałek
  <kbd>&lt;down&gt;</kbd>: Wybierz następny kawałek
  <kbd>z</kbd>: Cofnij
  <kbd>M</kbd>: Open external merge tool (git mergetool)
  <kbd>&lt;space&gt;</kbd>: Wybierz kawałek
  <kbd>b</kbd>: Wybierz oba kawałki
  <kbd>&lt;esc&gt;</kbd>: Wróć do panelu plików
</pre>

## Schowek

<pre>
  <kbd>&lt;space&gt;</kbd>: Zastosuj
  <kbd>g</kbd>: Wyciągnij
  <kbd>d</kbd>: Porzuć
  <kbd>n</kbd>: Nowa gałąź
  <kbd>r</kbd>: Rename stash
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: Przeglądaj pliki commita
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Status

<pre>
  <kbd>o</kbd>: Otwórz konfigurację
  <kbd>e</kbd>: Edytuj konfigurację
  <kbd>u</kbd>: Sprawdź aktualizacje
  <kbd>&lt;enter&gt;</kbd>: Switch to a recent repo
  <kbd>a</kbd>: Pokaż wszystkie logi gałęzi
</pre>

## Sub-commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy commit SHA to clipboard
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Checkout commit
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: Open commit in browser
  <kbd>n</kbd>: Create new branch off of commit
  <kbd>g</kbd>: Wyświetl opcje resetu
  <kbd>c</kbd>: Kopiuj commit (przebieranie)
  <kbd>C</kbd>: Kopiuj zakres commitów (przebieranie)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: Przeglądaj pliki commita
  <kbd>/</kbd>: Search the current view by text
</pre>

## Submodules

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Copy submodule name to clipboard
  <kbd>&lt;enter&gt;</kbd>: Enter submodule
  <kbd>&lt;space&gt;</kbd>: Enter submodule
  <kbd>d</kbd>: Remove submodule
  <kbd>u</kbd>: Update submodule
  <kbd>n</kbd>: Add new submodule
  <kbd>e</kbd>: Update submodule URL
  <kbd>i</kbd>: Initialize submodule
  <kbd>b</kbd>: View bulk submodule options
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Tags

<pre>
  <kbd>&lt;space&gt;</kbd>: Przełącz
  <kbd>d</kbd>: Delete tag
  <kbd>P</kbd>: Push tag
  <kbd>n</kbd>: Create tag
  <kbd>g</kbd>: Wyświetl opcje resetu
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: View commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Worktrees

<pre>
  <kbd>n</kbd>: Create worktree
  <kbd>&lt;space&gt;</kbd>: Switch to worktree
  <kbd>&lt;enter&gt;</kbd>: Switch to worktree
  <kbd>o</kbd>: Open in editor
  <kbd>d</kbd>: Remove worktree
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Zwykłe

<pre>
  <kbd>mouse wheel down</kbd>: Przewiń w dół (fn+up)
  <kbd>mouse wheel up</kbd>: Przewiń w górę (fn+down)
</pre>
