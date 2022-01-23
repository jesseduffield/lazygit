_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit Keybindings

## Globalne

<pre>
  <kbd>ctrl+r</kbd>: switch to a recent repo (<c-r>)
  <kbd>pgup</kbd>: scroll up main panel (fn+up)
  <kbd>pgdown</kbd>: scroll down main panel (fn+down)
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
  <kbd>]</kbd>: next tab
  <kbd>[</kbd>: previous tab
</pre>

## Gałęzie Panel (Branches Tab)

<pre>
  <kbd>space</kbd>: przełącz
  <kbd>o</kbd>: utwórz żądanie pobrania
  <kbd>O</kbd>: utwórz opcje żądania ściągnięcia
  <kbd>ctrl+y</kbd>: skopiuj adres URL żądania pobrania do schowka
  <kbd>c</kbd>: przełącz używając nazwy
  <kbd>F</kbd>: wymuś przełączenie
  <kbd>n</kbd>: nowa gałąź
  <kbd>d</kbd>: usuń gałąź
  <kbd>r</kbd>: zmiana bazy gałęzi
  <kbd>M</kbd>: scal do obecnej gałęzi
  <kbd>i</kbd>: show git-flow options
  <kbd>f</kbd>: fast-forward this branch from its upstream
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>R</kbd>: rename branch
  <kbd>ctrl+o</kbd>: copy branch name to clipboard
  <kbd>enter</kbd>: view commits
</pre>

## Gałęzie Panel (Remote Branches (in Remotes tab))

<pre>
  <kbd>esc</kbd>: wróć do listy repozytoriów zdalnych
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>enter</kbd>: view commits
  <kbd>space</kbd>: przełącz
  <kbd>n</kbd>: nowa gałąź
  <kbd>M</kbd>: scal do obecnej gałęzi
  <kbd>d</kbd>: usuń gałąź
  <kbd>r</kbd>: zmiana bazy gałęzi
  <kbd>u</kbd>: set as upstream of checked-out branch
</pre>

## Gałęzie Panel (Remotes Tab)

<pre>
  <kbd>f</kbd>: fetch remote
  <kbd>n</kbd>: add new remote
  <kbd>d</kbd>: remove remote
  <kbd>e</kbd>: edit remote
</pre>

## Gałęzie Panel (Sub-commits)

<pre>
  <kbd>enter</kbd>: przeglądaj pliki commita
  <kbd>space</kbd>: checkout commit
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>n</kbd>: nowa gałąź
  <kbd>c</kbd>: kopiuj commit (przebieranie)
  <kbd>C</kbd>: kopiuj zakres commitów (przebieranie)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
</pre>

## Gałęzie Panel (Tags Tab)

<pre>
  <kbd>space</kbd>: przełącz
  <kbd>d</kbd>: delete tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: create tag
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>enter</kbd>: view commits
</pre>

## Pliki commita Panel

<pre>
  <kbd>ctrl+o</kbd>: copy the committed file name to the clipboard
  <kbd>c</kbd>: plik wybierania
  <kbd>d</kbd>: porzuć zmiany commita dla tego pliku
  <kbd>o</kbd>: otwórz plik
  <kbd>e</kbd>: edytuj plik
  <kbd>space</kbd>: toggle file included in patch
  <kbd>enter</kbd>: enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: toggle file tree view
</pre>

## Commity Panel (Commity)

<pre>
  <kbd>c</kbd>: kopiuj commit (przebieranie)
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
  <kbd>C</kbd>: kopiuj zakres commitów (przebieranie)
  <kbd>v</kbd>: wklej commity (przebieranie)
  <kbd>n</kbd>: create new branch off of commit
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
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
  <kbd>A</kbd>: popraw commit zmianami z poczekalni
  <kbd>t</kbd>: odwróć commit
  <kbd>ctrl+l</kbd>: open log menu
  <kbd>g</kbd>: zresetuj do tego commita
  <kbd>enter</kbd>: przeglądaj pliki commita
  <kbd>space</kbd>: checkout commit
  <kbd>T</kbd>: tag commit
  <kbd>ctrl+y</kbd>: copy commit message to clipboard
  <kbd>o</kbd>: open commit in browser
  <kbd>b</kbd>: view bisect options
</pre>

## Commity Panel (Reflog Tab)

<pre>
  <kbd>enter</kbd>: przeglądaj pliki commita
  <kbd>space</kbd>: checkout commit
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>c</kbd>: kopiuj commit (przebieranie)
  <kbd>C</kbd>: kopiuj zakres commitów (przebieranie)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
</pre>

## Extras Panel

<pre>
  <kbd>@</kbd>: open command log menu
</pre>

## Pliki Panel (Pliki)

<pre>
  <kbd>d</kbd>: pokaż opcje porzucania zmian
  <kbd>D</kbd>: wyświetl opcje resetu
  <kbd>f</kbd>: pobierz
  <kbd>ctrl+o</kbd>: copy the file name to the clipboard
  <kbd>ctrl+w</kbd>: Toggle whether or not whitespace changes are shown in the diff view
  <kbd>space</kbd>: przełącz stan poczekalni
  <kbd>ctrl+b</kbd>: Filter files (staged/unstaged)
  <kbd>c</kbd>: Zatwierdź zmiany
  <kbd>w</kbd>: zatwierdź zmiany bez skryptu pre-commit
  <kbd>A</kbd>: Zmień ostatni commit
  <kbd>C</kbd>: Zatwierdź zmiany używając edytora
  <kbd>e</kbd>: edytuj plik
  <kbd>o</kbd>: otwórz plik
  <kbd>i</kbd>: dodaj do .gitignore
  <kbd>r</kbd>: odśwież pliki
  <kbd>s</kbd>: przechowaj zmiany
  <kbd>S</kbd>: wyświetl opcje schowka
  <kbd>a</kbd>: przełącz stan poczekalni wszystkich
  <kbd>enter</kbd>: zatwierdź pojedyncze linie
  <kbd>:</kbd>: wykonaj własną komendę
  <kbd>g</kbd>: view upstream reset options
  <kbd>`</kbd>: toggle file tree view
  <kbd>M</kbd>: open external merge tool (git mergetool)
</pre>

## Pliki Panel (Submodules)

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

## Główne Panel (Scalanie)

<pre>
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
  <kbd>esc</kbd>: wróć do panelu plików
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>space</kbd>: wybierz kawałek
  <kbd>b</kbd>: wybierz wszystkie kawałki
  <kbd>◄</kbd>: poprzedni konflikt
  <kbd>►</kbd>: następny konflikt
  <kbd>▲</kbd>: wybierz poprzedni kawałek
  <kbd>▼</kbd>: wybierz następny kawałek
  <kbd>z</kbd>: cofnij
</pre>

## Główne Panel (Zwykłe)

<pre>
  <kbd>Ő</kbd>: przewiń w dół (fn+up)
  <kbd>ő</kbd>: przewiń w górę (fn+down)
</pre>

## Główne Panel (Patch Building)

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
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
</pre>

## Główne Panel (Poczekalnia)

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
  <kbd>o</kbd>: otwórz plik
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
  <kbd>c</kbd>: Zatwierdź zmiany
  <kbd>w</kbd>: zatwierdź zmiany bez skryptu pre-commit
  <kbd>C</kbd>: Zatwierdź zmiany używając edytora
</pre>

## Menu Panel

<pre>
  <kbd>esc</kbd>: close menu
</pre>

## Schowek Panel

<pre>
  <kbd>enter</kbd>: view stash entry's files
  <kbd>space</kbd>: zastosuj
  <kbd>g</kbd>: wyciągnij
  <kbd>d</kbd>: porzuć
  <kbd>n</kbd>: nowa gałąź
</pre>

## Status Panel

<pre>
  <kbd>e</kbd>: edytuj konfigurację
  <kbd>o</kbd>: otwórz konfigurację
  <kbd>u</kbd>: sprawdź aktualizacje
  <kbd>enter</kbd>: switch to a recent repo
  <kbd>a</kbd>: pokaż wszystkie logi gałęzi
</pre>
