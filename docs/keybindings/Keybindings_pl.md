_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Globalne

<pre>
  <kbd>&lt;c-r&gt;</kbd>: switch to a recent repo
  <kbd>&lt;pgup&gt;</kbd>: scroll up main panel (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: scroll down main panel (fn+down/shift+j)
  <kbd>@</kbd>: open command log menu
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: wykonaj własną komendę
  <kbd>&lt;c-p&gt;</kbd>: view custom patch options
  <kbd>m</kbd>: widok scalenia/opcje zmiany bazy
  <kbd>R</kbd>: odśwież
  <kbd>+</kbd>: next screen mode (normal/half/fullscreen)
  <kbd>_</kbd>: prev screen mode
  <kbd>?</kbd>: open menu
  <kbd>&lt;c-s&gt;</kbd>: view filter-by-path options
  <kbd>W</kbd>: open diff menu
  <kbd>&lt;c-e&gt;</kbd>: open diff menu
  <kbd>&lt;c-w&gt;</kbd>: Toggle whether or not whitespace changes are shown in the diff view
  <kbd>z</kbd>: undo (via reflog) (experimental)
  <kbd>&lt;c-z&gt;</kbd>: redo (via reflog) (experimental)
  <kbd>P</kbd>: push
  <kbd>p</kbd>: pull
</pre>

## List Panel Navigation

<pre>
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd>&lt;</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>&gt;</kbd>: scroll to bottom
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
  <kbd>]</kbd>: next tab
  <kbd>[</kbd>: previous tab
</pre>

## Commit Summary

<pre>
  <kbd>&lt;enter&gt;</kbd>: potwierdź
  <kbd>&lt;esc&gt;</kbd>: zamknij
</pre>

## Commity

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy commit SHA to clipboard
  <kbd>&lt;c-r&gt;</kbd>: reset cherry-picked (copied) commits selection
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
  <kbd>&lt;c-j&gt;</kbd>: przenieś commit 1 w dół
  <kbd>&lt;c-k&gt;</kbd>: przenieś commit 1 w górę
  <kbd>v</kbd>: wklej commity (przebieranie)
  <kbd>A</kbd>: popraw commit zmianami z poczekalni
  <kbd>a</kbd>: reset commit author
  <kbd>t</kbd>: odwróć commit
  <kbd>T</kbd>: tag commit
  <kbd>&lt;c-l&gt;</kbd>: open log menu
  <kbd>&lt;space&gt;</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>c</kbd>: kopiuj commit (przebieranie)
  <kbd>C</kbd>: kopiuj zakres commitów (przebieranie)
  <kbd>&lt;enter&gt;</kbd>: przeglądaj pliki commita
</pre>

## Confirmation Panel

<pre>
  <kbd>&lt;enter&gt;</kbd>: potwierdź
  <kbd>&lt;esc&gt;</kbd>: zamknij
</pre>

## Local Branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy branch name to clipboard
  <kbd>i</kbd>: show git-flow options
  <kbd>&lt;space&gt;</kbd>: przełącz
  <kbd>n</kbd>: nowa gałąź
  <kbd>o</kbd>: utwórz żądanie pobrania
  <kbd>O</kbd>: utwórz opcje żądania ściągnięcia
  <kbd>&lt;c-y&gt;</kbd>: skopiuj adres URL żądania pobrania do schowka
  <kbd>c</kbd>: przełącz używając nazwy
  <kbd>F</kbd>: wymuś przełączenie
  <kbd>d</kbd>: usuń gałąź
  <kbd>r</kbd>: zmiana bazy gałęzi
  <kbd>M</kbd>: scal do obecnej gałęzi
  <kbd>f</kbd>: fast-forward this branch from its upstream
  <kbd>T</kbd>: create tag
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>R</kbd>: rename branch
  <kbd>u</kbd>: set/unset upstream
  <kbd>&lt;enter&gt;</kbd>: view commits
</pre>

## Main Panel (Patch Building)

<pre>
  <kbd>&lt;left&gt;</kbd>: poprzedni kawałek
  <kbd>&lt;right&gt;</kbd>: następny kawałek
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: copy the selected text to the clipboard
  <kbd>o</kbd>: otwórz plik
  <kbd>e</kbd>: edytuj plik
  <kbd>&lt;space&gt;</kbd>: add/remove line(s) to patch
  <kbd>&lt;esc&gt;</kbd>: wyście z trybu "linia po linii"
</pre>

## Menu

<pre>
  <kbd>&lt;enter&gt;</kbd>: wykonaj
  <kbd>&lt;esc&gt;</kbd>: zamknij
</pre>

## Pliki

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy the file name to the clipboard
  <kbd>d</kbd>: pokaż opcje porzucania zmian
  <kbd>&lt;space&gt;</kbd>: przełącz stan poczekalni
  <kbd>&lt;c-b&gt;</kbd>: Filter files (staged/unstaged)
  <kbd>c</kbd>: Zatwierdź zmiany
  <kbd>w</kbd>: zatwierdź zmiany bez skryptu pre-commit
  <kbd>A</kbd>: Zmień ostatni commit
  <kbd>C</kbd>: Zatwierdź zmiany używając edytora
  <kbd>e</kbd>: edytuj plik
  <kbd>o</kbd>: otwórz plik
  <kbd>i</kbd>: ignore or exclude file
  <kbd>r</kbd>: odśwież pliki
  <kbd>s</kbd>: przechowaj zmiany
  <kbd>S</kbd>: wyświetl opcje schowka
  <kbd>a</kbd>: przełącz stan poczekalni wszystkich
  <kbd>&lt;enter&gt;</kbd>: zatwierdź pojedyncze linie
  <kbd>g</kbd>: view upstream reset options
  <kbd>D</kbd>: wyświetl opcje resetu
  <kbd>`</kbd>: toggle file tree view
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>f</kbd>: pobierz
</pre>

## Pliki commita

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy the committed file name to the clipboard
  <kbd>c</kbd>: plik wybierania
  <kbd>d</kbd>: porzuć zmiany commita dla tego pliku
  <kbd>o</kbd>: otwórz plik
  <kbd>e</kbd>: edytuj plik
  <kbd>&lt;space&gt;</kbd>: toggle file included in patch
  <kbd>a</kbd>: toggle all files included in patch
  <kbd>&lt;enter&gt;</kbd>: enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: toggle file tree view
</pre>

## Poczekalnia

<pre>
  <kbd>&lt;left&gt;</kbd>: poprzedni kawałek
  <kbd>&lt;right&gt;</kbd>: następny kawałek
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: copy the selected text to the clipboard
  <kbd>o</kbd>: otwórz plik
  <kbd>e</kbd>: edytuj plik
  <kbd>&lt;esc&gt;</kbd>: wróć do panelu plików
  <kbd>&lt;tab&gt;</kbd>: switch to other panel (staged/unstaged changes)
  <kbd>&lt;space&gt;</kbd>: toggle line staged / unstaged
  <kbd>d</kbd>: delete change (git reset)
  <kbd>E</kbd>: edit hunk
  <kbd>c</kbd>: Zatwierdź zmiany
  <kbd>w</kbd>: zatwierdź zmiany bez skryptu pre-commit
  <kbd>C</kbd>: Zatwierdź zmiany używając edytora
</pre>

## Reflog

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy commit SHA to clipboard
  <kbd>&lt;space&gt;</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>c</kbd>: kopiuj commit (przebieranie)
  <kbd>C</kbd>: kopiuj zakres commitów (przebieranie)
  <kbd>&lt;c-r&gt;</kbd>: reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: view commits
</pre>

## Remote Branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy branch name to clipboard
  <kbd>&lt;space&gt;</kbd>: przełącz
  <kbd>n</kbd>: nowa gałąź
  <kbd>M</kbd>: scal do obecnej gałęzi
  <kbd>r</kbd>: zmiana bazy gałęzi
  <kbd>d</kbd>: usuń gałąź
  <kbd>u</kbd>: set as upstream of checked-out branch
  <kbd>&lt;esc&gt;</kbd>: wróć do listy repozytoriów zdalnych
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>&lt;enter&gt;</kbd>: view commits
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
  <kbd>e</kbd>: edytuj plik
  <kbd>o</kbd>: otwórz plik
  <kbd>&lt;left&gt;</kbd>: poprzedni konflikt
  <kbd>&lt;right&gt;</kbd>: następny konflikt
  <kbd>&lt;up&gt;</kbd>: wybierz poprzedni kawałek
  <kbd>&lt;down&gt;</kbd>: wybierz następny kawałek
  <kbd>z</kbd>: cofnij
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>&lt;space&gt;</kbd>: wybierz kawałek
  <kbd>b</kbd>: wybierz wszystkie kawałki
  <kbd>&lt;esc&gt;</kbd>: wróć do panelu plików
</pre>

## Schowek

<pre>
  <kbd>&lt;space&gt;</kbd>: zastosuj
  <kbd>g</kbd>: wyciągnij
  <kbd>d</kbd>: porzuć
  <kbd>n</kbd>: nowa gałąź
  <kbd>r</kbd>: rename stash
  <kbd>&lt;enter&gt;</kbd>: przeglądaj pliki commita
</pre>

## Status

<pre>
  <kbd>o</kbd>: otwórz konfigurację
  <kbd>e</kbd>: edytuj konfigurację
  <kbd>u</kbd>: sprawdź aktualizacje
  <kbd>&lt;enter&gt;</kbd>: switch to a recent repo
  <kbd>a</kbd>: pokaż wszystkie logi gałęzi
</pre>

## Sub-commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy commit SHA to clipboard
  <kbd>&lt;space&gt;</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: create new branch off of commit
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>c</kbd>: kopiuj commit (przebieranie)
  <kbd>C</kbd>: kopiuj zakres commitów (przebieranie)
  <kbd>&lt;c-r&gt;</kbd>: reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: przeglądaj pliki commita
</pre>

## Submodules

<pre>
  <kbd>&lt;c-o&gt;</kbd>: copy submodule name to clipboard
  <kbd>&lt;enter&gt;</kbd>: enter submodule
  <kbd>d</kbd>: remove submodule
  <kbd>u</kbd>: update submodule
  <kbd>n</kbd>: add new submodule
  <kbd>e</kbd>: update submodule URL
  <kbd>i</kbd>: initialize submodule
  <kbd>b</kbd>: view bulk submodule options
</pre>

## Tags

<pre>
  <kbd>&lt;space&gt;</kbd>: przełącz
  <kbd>d</kbd>: delete tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: create tag
  <kbd>g</kbd>: wyświetl opcje resetu
  <kbd>&lt;enter&gt;</kbd>: view commits
</pre>

## Zwykłe

<pre>
  <kbd>mouse wheel down</kbd>: przewiń w dół (fn+up)
  <kbd>mouse wheel up</kbd>: przewiń w górę (fn+down)
</pre>
