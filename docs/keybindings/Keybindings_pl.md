_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit Keybindings

## Globalne

<pre>
  <kbd>ctrl+r</kbd>: switch to a recent repo
  <kbd>pgup</kbd>: scroll up main panel (fn+up/shift+k)
  <kbd>pgdown</kbd>: scroll down main panel (fn+down/shift+j)
  <kbd>m</kbd>: widok scalenia/opcje zmiany bazy
  <kbd>ctrl+p</kbd>: view custom patch options
  <kbd>R</kbd>: odśwież
  <kbd>x</kbd>: open menu
  <kbd>+</kbd>: next screen mode (normal/half/fullscreen)
  <kbd>_</kbd>: prev screen mode
  <kbd>ctrl+s</kbd>: view filter-by-path options
  <kbd>W</kbd>: open diff menu
  <kbd>ctrl+e</kbd>: open diff menu
  <kbd>@</kbd>: open command log menu
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: wykonaj własną komendę
  <kbd>z</kbd>: undo (via reflog) (experimental)
  <kbd>ctrl+z</kbd>: redo (via reflog) (experimental)
  <kbd>P</kbd>: push
  <kbd>p</kbd>: pull
</pre>

## List Panel Navigation

<pre>
  <kbd>.</kbd>: next page
  <kbd>,</kbd>: previous page
  <kbd><</kbd>: scroll to top
  <kbd>></kbd>: scroll to bottom
  <kbd>/</kbd>: start search
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
  <kbd>]</kbd>: next tab
  <kbd>[</kbd>: previous tab
</pre>

## Commity

<pre>
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>b</kbd>: view bisect options
  <kbd>s</kbd>: ściśnij
  <kbd>f</kbd>: napraw commit
  <kbd>r</kbd>: zmień nazwę commita
  <kbd>R</kbd>: zmień nazwę commita w edytorze
  <kbd>d</kbd>: usuń commit
  <kbd>e</kbd>: edytuj commit
  <kbd>p</kbd>: wybierz commit (podczas zmiany bazy)
  <kbd>F</kbd>: utwórz commit naprawczy dla tego commita
  <kbd>S</kbd>: spłaszcz wszystkie commity naprawcze powyżej zaznaczonych commitów (autosquash)
  <kbd>ctrl+j</kbd>: przenieś commit 1 w dół
  <kbd>ctrl+k</kbd>: przenieś commit 1 w górę
  <kbd>v</kbd>: wklej commity (przebieranie)
  <kbd>A</kbd>: popraw commit zmianami z poczekalni
  <kbd>a</kbd>: reset commit author
  <kbd>t</kbd>: odwróć commit
  <kbd>T</kbd>: tag commit
  <kbd>ctrl+l</kbd>: open log menu
  <kbd>space</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>c</kbd>: kopiuj commit (przebieranie)
  <kbd>C</kbd>: kopiuj zakres commitów (przebieranie)
  <kbd>enter</kbd>: przeglądaj pliki commita
</pre>

## Local Branches

<pre>
  <kbd>ctrl+o</kbd>: copy branch name to clipboard
  <kbd>i</kbd>: show git-flow options
  <kbd>space</kbd>: przełącz
  <kbd>n</kbd>: nowa gałąź
  <kbd>o</kbd>: utwórz żądanie pobrania
  <kbd>O</kbd>: utwórz opcje żądania ściągnięcia
  <kbd>ctrl+y</kbd>: skopiuj adres URL żądania pobrania do schowka
  <kbd>c</kbd>: przełącz używając nazwy
  <kbd>F</kbd>: wymuś przełączenie
  <kbd>d</kbd>: usuń gałąź
  <kbd>r</kbd>: zmiana bazy gałęzi
  <kbd>M</kbd>: scal do obecnej gałęzi
  <kbd>f</kbd>: fast-forward this branch from its upstream
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>R</kbd>: rename branch
  <kbd>u</kbd>: set/unset upstream
  <kbd>enter</kbd>: view commits
</pre>

## Main Panel (Patch Building)

<pre>
  <kbd>esc</kbd>: wyście z trybu "linia po linii"
  <kbd>o</kbd>: otwórz plik
  <kbd>▲</kbd>: poprzednia linia
  <kbd>▼</kbd>: następna linia
  <kbd>◄</kbd>: poprzedni kawałek
  <kbd>►</kbd>: następny kawałek
  <kbd>ctrl+o</kbd>: copy the selected text to the clipboard
  <kbd>space</kbd>: add/remove line(s) to patch
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
</pre>

## Pliki

<pre>
  <kbd>ctrl+o</kbd>: copy the file name to the clipboard
  <kbd>ctrl+w</kbd>: Toggle whether or not whitespace changes are shown in the diff view
  <kbd>d</kbd>: pokaż opcje porzucania zmian
  <kbd>space</kbd>: przełącz stan poczekalni
  <kbd>ctrl+b</kbd>: Filter files (staged/unstaged)
  <kbd>c</kbd>: Zatwierdź zmiany
  <kbd>w</kbd>: zatwierdź zmiany bez skryptu pre-commit
  <kbd>A</kbd>: Zmień ostatni commit
  <kbd>C</kbd>: Zatwierdź zmiany używając edytora
  <kbd>e</kbd>: edytuj plik
  <kbd>o</kbd>: otwórz plik
  <kbd>i</kbd>: Ignore or Exclude file
  <kbd>r</kbd>: odśwież pliki
  <kbd>s</kbd>: przechowaj zmiany
  <kbd>S</kbd>: wyświetl opcje schowka
  <kbd>a</kbd>: przełącz stan poczekalni wszystkich
  <kbd>enter</kbd>: zatwierdź pojedyncze linie
  <kbd>g</kbd>: view upstream reset options
  <kbd>D</kbd>: wyświetl opcje resetu
  <kbd>`</kbd>: toggle file tree view
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>f</kbd>: pobierz
</pre>

## Pliki commita

<pre>
  <kbd>ctrl+o</kbd>: copy the committed file name to the clipboard
  <kbd>c</kbd>: plik wybierania
  <kbd>d</kbd>: porzuć zmiany commita dla tego pliku
  <kbd>o</kbd>: otwórz plik
  <kbd>e</kbd>: edytuj plik
  <kbd>space</kbd>: toggle file included in patch
  <kbd>a</kbd>: toggle all files included in patch
  <kbd>enter</kbd>: enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: toggle file tree view
</pre>

## Poczekalnia

<pre>
  <kbd>esc</kbd>: wróć do panelu plików
  <kbd>space</kbd>: toggle line staged / unstaged
  <kbd>d</kbd>: delete change (git reset)
  <kbd>tab</kbd>: switch to other panel
  <kbd>o</kbd>: otwórz plik
  <kbd>▲</kbd>: poprzednia linia
  <kbd>▼</kbd>: następna linia
  <kbd>◄</kbd>: poprzedni kawałek
  <kbd>►</kbd>: następny kawałek
  <kbd>ctrl+o</kbd>: copy the selected text to the clipboard
  <kbd>e</kbd>: edytuj plik
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>E</kbd>: edit hunk
</pre>

## Reflog

<pre>
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
  <kbd>space</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>c</kbd>: kopiuj commit (przebieranie)
  <kbd>C</kbd>: kopiuj zakres commitów (przebieranie)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>enter</kbd>: view commits
</pre>

## Remote Branches

<pre>
  <kbd>space</kbd>: przełącz
  <kbd>n</kbd>: nowa gałąź
  <kbd>M</kbd>: scal do obecnej gałęzi
  <kbd>r</kbd>: zmiana bazy gałęzi
  <kbd>d</kbd>: usuń gałąź
  <kbd>u</kbd>: set as upstream of checked-out branch
  <kbd>esc</kbd>: wróć do listy repozytoriów zdalnych
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>enter</kbd>: view commits
</pre>

## Remotes

<pre>
  <kbd>f</kbd>: fetch remote
  <kbd>n</kbd>: add new remote
  <kbd>d</kbd>: remove remote
  <kbd>e</kbd>: edit remote
</pre>

## Scalanie

<pre>
  <kbd>esc</kbd>: wróć do panelu plików
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>space</kbd>: wybierz kawałek
  <kbd>b</kbd>: wybierz wszystkie kawałki
  <kbd>◄</kbd>: poprzedni konflikt
  <kbd>►</kbd>: następny konflikt
  <kbd>▲</kbd>: wybierz poprzedni kawałek
  <kbd>▼</kbd>: wybierz następny kawałek
  <kbd>e</kbd>: edytuj plik
  <kbd>o</kbd>: otwórz plik
  <kbd>z</kbd>: cofnij
</pre>

## Schowek

<pre>
  <kbd>space</kbd>: zastosuj
  <kbd>g</kbd>: wyciągnij
  <kbd>d</kbd>: porzuć
  <kbd>n</kbd>: nowa gałąź
  <kbd>enter</kbd>: przeglądaj pliki commita
</pre>

## Status

<pre>
  <kbd>e</kbd>: edytuj konfigurację
  <kbd>o</kbd>: otwórz konfigurację
  <kbd>u</kbd>: sprawdź aktualizacje
  <kbd>enter</kbd>: switch to a recent repo
  <kbd>a</kbd>: pokaż wszystkie logi gałęzi
</pre>

## Sub-commits

<pre>
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
  <kbd>space</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>c</kbd>: kopiuj commit (przebieranie)
  <kbd>C</kbd>: kopiuj zakres commitów (przebieranie)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>enter</kbd>: przeglądaj pliki commita
</pre>

## Submodules

<pre>
  <kbd>ctrl+o</kbd>: copy submodule name to clipboard
  <kbd>enter</kbd>: enter submodule
  <kbd>d</kbd>: remove submodule
  <kbd>u</kbd>: update submodule
  <kbd>n</kbd>: add new submodule
  <kbd>e</kbd>: update submodule URL
  <kbd>i</kbd>: initialize submodule
  <kbd>b</kbd>: view bulk submodule options
</pre>

## Tags

<pre>
  <kbd>space</kbd>: przełącz
  <kbd>d</kbd>: delete tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: create tag
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>enter</kbd>: view commits
</pre>

## Zwykłe

<pre>
  <kbd>mouse wheel down</kbd>: przewiń w dół (fn+up)
  <kbd>mouse wheel up</kbd>: przewiń w górę (fn+down)
</pre>
