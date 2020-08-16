# Lazygit Keybindings

## Globalne

<pre>
  <kbd>pgup</kbd>: scroll up main panel (fn+up)
  <kbd>pgdown</kbd>: scroll down main panel (fn+down)
  <kbd>m</kbd>: view merge/rebase options
  <kbd>ctrl+p</kbd>: view custom patch options
  <kbd>P</kbd>: push
  <kbd>p</kbd>: pull
  <kbd>R</kbd>: odśwież
  <kbd>x</kbd>: open menu
  <kbd>z</kbd>: undo (via reflog) (experimental)
  <kbd>ctrl+z</kbd>: redo (via reflog) (experimental)
  <kbd>+</kbd>: next screen mode (normal/half/fullscreen)
  <kbd>_</kbd>: prev screen mode
  <kbd>:</kbd>: execute custom command
  <kbd>|</kbd>: view scoping options
  <kbd>∂</kbd>: open diff menu
</pre>

## Gałęzie Panel

<pre>
  <kbd>]</kbd>: next tab
  <kbd>[</kbd>: previous tab
</pre>

## Gałęzie Panel (Branches Tab)

<pre>
  <kbd>space</kbd>: przełącz
  <kbd>o</kbd>: utwórz żądanie wyciągnięcia
  <kbd>c</kbd>: przełącz używając nazwy
  <kbd>F</kbd>: wymuś przełączenie
  <kbd>n</kbd>: nowa gałąź
  <kbd>d</kbd>: usuń gałąź
  <kbd>r</kbd>: rebase branch
  <kbd>M</kbd>: scal do obecnej gałęzi
  <kbd>i</kbd>: show git-flow options
  <kbd>f</kbd>: fast-forward this branch from its upstream
  <kbd>g</kbd>: view reset options
  <kbd>R</kbd>: rename branch
  <kbd>ctrl+o</kbd>: copy branch name to clipboard
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Gałęzie Panel (Remote Branches (in Remotes tab))

<pre>
  <kbd>esc</kbd>: return to remotes list
  <kbd>g</kbd>: view reset options
  <kbd>space</kbd>: przełącz
  <kbd>n</kbd>: nowa gałąź
  <kbd>M</kbd>: scal do obecnej gałęzi
  <kbd>d</kbd>: usuń gałąź
  <kbd>r</kbd>: rebase branch
  <kbd>u</kbd>: set as upstream of checked-out branch
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Gałęzie Panel (Remotes Tab)

<pre>
  <kbd>f</kbd>: fetch remote
  <kbd>n</kbd>: add new remote
  <kbd>d</kbd>: remove remote
  <kbd>e</kbd>: edit remote
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Gałęzie Panel (Tags Tab)

<pre>
  <kbd>space</kbd>: przełącz
  <kbd>d</kbd>: delete tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: create tag
  <kbd>g</kbd>: view reset options
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Commit files Panel

<pre>
  <kbd>esc</kbd>: go back
  <kbd>c</kbd>: checkout file
  <kbd>d</kbd>: discard this commit's changes to this file
  <kbd>o</kbd>: otwórz plik
  <kbd>e</kbd>: edytuj plik
  <kbd>space</kbd>: toggle file included in patch
  <kbd>enter</kbd>: enter file to add selected lines to the patch
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Commity Panel

<pre>
  <kbd>]</kbd>: next tab
  <kbd>[</kbd>: previous tab
</pre>

## Commity Panel (Commits Tab)

<pre>
  <kbd>s</kbd>: ściśnij w dół
  <kbd>r</kbd>: przemianuj commit
  <kbd>R</kbd>: przemianuj commit w edytorze
  <kbd>g</kbd>: zresetuj do tego commita
  <kbd>f</kbd>: napraw commit
  <kbd>F</kbd>: create fixup commit for this commit
  <kbd>S</kbd>: squash above commits
  <kbd>d</kbd>: delete commit
  <kbd>ctrl+j</kbd>: move commit down one
  <kbd>ctrl+k</kbd>: move commit up one
  <kbd>e</kbd>: edit commit
  <kbd>A</kbd>: amend commit with staged changes
  <kbd>p</kbd>: pick commit (when mid-rebase)
  <kbd>t</kbd>: revert commit
  <kbd>c</kbd>: copy commit (cherry-pick)
  <kbd>ctrl+o</kbd>: copy commit SHA to clipboard
  <kbd>C</kbd>: copy commit range (cherry-pick)
  <kbd>v</kbd>: paste commits (cherry-pick)
  <kbd>enter</kbd>: view commit's files
  <kbd>space</kbd>: checkout commit
  <kbd>n</kbd>: create new branch off of commit
  <kbd>T</kbd>: tag commit
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Commity Panel (Reflog Tab)

<pre>
  <kbd>space</kbd>: checkout commit
  <kbd>g</kbd>: view reset options
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Pliki Panel

<pre>
  <kbd>c</kbd>: commituj zmiany
  <kbd>w</kbd>: commit changes without pre-commit hook
  <kbd>A</kbd>: zmień ostatnie zatwierdzenie
  <kbd>C</kbd>: commituj zmiany używając edytora z gita
  <kbd>space</kbd>: przełącz zatwierdzenie
  <kbd>d</kbd>: view 'discard changes' options
  <kbd>e</kbd>: edytuj plik
  <kbd>o</kbd>: otwórz plik
  <kbd>i</kbd>: dodaj do .gitignore
  <kbd>r</kbd>: odśwież pliki
  <kbd>s</kbd>: przechowaj pliki
  <kbd>S</kbd>: view stash options
  <kbd>a</kbd>: przełącz wszystkie zatwierdzenia
  <kbd>D</kbd>: view reset options
  <kbd>enter</kbd>: zatwierdź pojedyncze linie
  <kbd>f</kbd>: fetch
  <kbd>g</kbd>: view upstream reset options
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Main Panel (Merging)

<pre>
  <kbd>esc</kbd>: wróć do panelu plików
  <kbd>space</kbd>: pick hunk
  <kbd>b</kbd>: pick both hunks
  <kbd>◄</kbd>: select previous conflict
  <kbd>►</kbd>: select next conflict
  <kbd>▲</kbd>: select top hunk
  <kbd>▼</kbd>: select bottom hunk
  <kbd>z</kbd>: cofnij
</pre>

## Main Panel (Normal)

<pre>
  <kbd>￣</kbd>: scroll down (fn+up)
  <kbd>￤</kbd>: scroll up (fn+down)
</pre>

## Main Panel (Patch Building)

<pre>
  <kbd>esc</kbd>: exit line-by-line mode
  <kbd>o</kbd>: otwórz plik
  <kbd>▲</kbd>: select previous line
  <kbd>▼</kbd>: select next line
  <kbd>◄</kbd>: select previous hunk
  <kbd>►</kbd>: select next hunk
  <kbd>space</kbd>: add/remove line(s) to patch
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
</pre>

## Main Panel (Zatwierdzanie)

<pre>
  <kbd>esc</kbd>: wróć do panelu plików
  <kbd>space</kbd>: toggle line staged / unstaged
  <kbd>d</kbd>: delete change (git reset)
  <kbd>tab</kbd>: switch to other panel
  <kbd>o</kbd>: otwórz plik
  <kbd>▲</kbd>: select previous line
  <kbd>▼</kbd>: select next line
  <kbd>◄</kbd>: select previous hunk
  <kbd>►</kbd>: select next hunk
  <kbd>e</kbd>: edytuj plik
  <kbd>o</kbd>: otwórz plik
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>c</kbd>: commituj zmiany
  <kbd>w</kbd>: commit changes without pre-commit hook
  <kbd>C</kbd>: commituj zmiany używając edytora z gita
</pre>

## Menu Panel

<pre>
  <kbd>esc</kbd>: close menu
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Schowek Panel

<pre>
  <kbd>space</kbd>: zastosuj
  <kbd>g</kbd>: wyciągnij
  <kbd>d</kbd>: porzuć
  <kbd>,</kbd>: previous page
  <kbd>.</kbd>: next page
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Status Panel

<pre>
  <kbd>e</kbd>: edytuj plik konfiguracyjny
  <kbd>o</kbd>: otwórz plik konfiguracyjny
  <kbd>u</kbd>: sprawdź aktualizacje
  <kbd>enter</kbd>: switch to a recent repo
</pre>
