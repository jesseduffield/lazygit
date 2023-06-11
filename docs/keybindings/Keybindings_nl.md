_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit Sneltoetsen

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Globale sneltoetsen

<pre>
  <kbd>&lt;c-r&gt;</kbd>: Wissel naar een recente repo
  <kbd>&lt;pgup&gt;</kbd>: Scroll naar beneden vanaf hoofdpaneel (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: Scroll naar beneden vanaf hoofdpaneel (fn+down/shift+j)
  <kbd>@</kbd>: Open command log menu
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: Voer aangepaste commando uit
  <kbd>&lt;c-p&gt;</kbd>: Bekijk aangepaste patch opties
  <kbd>m</kbd>: Bekijk merge/rebase opties
  <kbd>R</kbd>: Verversen
  <kbd>+</kbd>: Volgende scherm modus (normaal/half/groot)
  <kbd>_</kbd>: Vorige scherm modus
  <kbd>?</kbd>: Open menu
  <kbd>&lt;c-s&gt;</kbd>: Bekijk scoping opties
  <kbd>W</kbd>: Open diff menu
  <kbd>&lt;c-e&gt;</kbd>: Open diff menu
  <kbd>&lt;c-w&gt;</kbd>: Toggle whether or not whitespace changes are shown in the diff view
  <kbd>z</kbd>: Ongedaan maken (via reflog) (experimenteel)
  <kbd>&lt;c-z&gt;</kbd>: Redo (via reflog) (experimenteel)
  <kbd>P</kbd>: Push
  <kbd>p</kbd>: Pull
</pre>

## Lijstpaneel navigatie

<pre>
  <kbd>,</kbd>: Vorige pagina
  <kbd>.</kbd>: Volgende pagina
  <kbd>&lt;</kbd>: Scroll naar boven
  <kbd>&gt;</kbd>: Scroll naar beneden
  <kbd>/</kbd>: Start met zoeken
  <kbd>H</kbd>: Scroll left
  <kbd>L</kbd>: Scroll right
  <kbd>]</kbd>: Volgende tabblad
  <kbd>[</kbd>: Vorige tabblad
</pre>

## Bestanden

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Kopieer de bestandsnaam naar het klembord
  <kbd>d</kbd>: Bekijk 'veranderingen ongedaan maken' opties
  <kbd>&lt;space&gt;</kbd>: Toggle staged
  <kbd>&lt;c-b&gt;</kbd>: Filter files by status
  <kbd>c</kbd>: Commit veranderingen
  <kbd>w</kbd>: Commit veranderingen zonder pre-commit hook
  <kbd>A</kbd>: Wijzig laatste commit
  <kbd>C</kbd>: Commit veranderingen met de git editor
  <kbd>e</kbd>: Verander bestand
  <kbd>o</kbd>: Open bestand
  <kbd>i</kbd>: Ignore or exclude file
  <kbd>r</kbd>: Refresh bestanden
  <kbd>s</kbd>: Stash-bestanden
  <kbd>S</kbd>: Bekijk stash opties
  <kbd>a</kbd>: Toggle staged alle
  <kbd>&lt;enter&gt;</kbd>: Stage individuele hunks/lijnen
  <kbd>g</kbd>: Bekijk upstream reset opties
  <kbd>D</kbd>: Bekijk reset opties
  <kbd>`</kbd>: Toggle bestandsboom weergave
  <kbd>M</kbd>: Open external merge tool (git mergetool)
  <kbd>f</kbd>: Fetch
  <kbd>/</kbd>: Start met zoeken
</pre>

## Bevestigingspaneel

<pre>
  <kbd>&lt;enter&gt;</kbd>: Bevestig
  <kbd>&lt;esc&gt;</kbd>: Sluiten
</pre>

## Branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Kopieer branch name naar klembord
  <kbd>i</kbd>: Laat git-flow opties zien
  <kbd>&lt;space&gt;</kbd>: Uitchecken
  <kbd>n</kbd>: Nieuwe branch
  <kbd>o</kbd>: Maak een pull-request
  <kbd>O</kbd>: Bekijk opties voor pull-aanvraag
  <kbd>&lt;c-y&gt;</kbd>: Kopieer de URL van het pull-verzoek naar het klembord
  <kbd>c</kbd>: Uitchecken bij naam
  <kbd>F</kbd>: Forceer checkout
  <kbd>d</kbd>: Verwijder branch
  <kbd>r</kbd>: Rebase branch
  <kbd>M</kbd>: Merge in met huidige checked out branch
  <kbd>f</kbd>: Fast-forward deze branch vanaf zijn upstream
  <kbd>T</kbd>: Creëer tag
  <kbd>g</kbd>: Bekijk reset opties
  <kbd>R</kbd>: Hernoem branch
  <kbd>u</kbd>: Set/Unset upstream
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: Bekijk commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Commit bericht

<pre>
  <kbd>&lt;enter&gt;</kbd>: Bevestig
  <kbd>&lt;esc&gt;</kbd>: Sluiten
</pre>

## Commit bestanden

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Kopieer de vastgelegde bestandsnaam naar het klembord
  <kbd>c</kbd>: Bestand uitchecken
  <kbd>d</kbd>: Uitsluit deze commit zijn veranderingen aan dit bestand
  <kbd>o</kbd>: Open bestand
  <kbd>e</kbd>: Verander bestand
  <kbd>&lt;space&gt;</kbd>: Toggle bestand inbegrepen in patch
  <kbd>a</kbd>: Toggle all files included in patch
  <kbd>&lt;enter&gt;</kbd>: Enter bestand om geselecteerde regels toe te voegen aan de patch
  <kbd>`</kbd>: Toggle bestandsboom weergave
  <kbd>/</kbd>: Start met zoeken
</pre>

## Commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Kopieer commit SHA naar klembord
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (gekopieerde) commits selectie
  <kbd>b</kbd>: View bisect options
  <kbd>s</kbd>: Squash beneden
  <kbd>f</kbd>: Fixup commit
  <kbd>r</kbd>: Hernoem commit
  <kbd>R</kbd>: Hernoem commit met editor
  <kbd>d</kbd>: Verwijder commit
  <kbd>e</kbd>: Wijzig commit
  <kbd>p</kbd>: Kies commit (wanneer midden in rebase)
  <kbd>F</kbd>: Creëer fixup commit
  <kbd>S</kbd>: Squash bovenstaande commits
  <kbd>&lt;c-j&gt;</kbd>: Verplaats commit 1 naar beneden
  <kbd>&lt;c-k&gt;</kbd>: Verplaats commit 1 naar boven
  <kbd>v</kbd>: Plak commits (cherry-pick)
  <kbd>B</kbd>: Mark commit as base commit for rebase
  <kbd>A</kbd>: Wijzig commit met staged veranderingen
  <kbd>a</kbd>: Set/Reset commit author
  <kbd>t</kbd>: Commit ongedaan maken
  <kbd>T</kbd>: Tag commit
  <kbd>&lt;c-l&gt;</kbd>: Open log menu
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Checkout commit
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: Open commit in browser
  <kbd>n</kbd>: Creëer nieuwe branch van commit
  <kbd>g</kbd>: Bekijk reset opties
  <kbd>c</kbd>: Kopieer commit (cherry-pick)
  <kbd>C</kbd>: Kopieer commit reeks (cherry-pick)
  <kbd>&lt;enter&gt;</kbd>: Bekijk gecommite bestanden
  <kbd>/</kbd>: Start met zoeken
</pre>

## Menu

<pre>
  <kbd>&lt;enter&gt;</kbd>: Uitvoeren
  <kbd>&lt;esc&gt;</kbd>: Sluiten
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Mergen

<pre>
  <kbd>e</kbd>: Verander bestand
  <kbd>o</kbd>: Open bestand
  <kbd>&lt;left&gt;</kbd>: Selecteer voorgaand conflict
  <kbd>&lt;right&gt;</kbd>: Selecteer volgende conflict
  <kbd>&lt;up&gt;</kbd>: Selecteer bovenste hunk
  <kbd>&lt;down&gt;</kbd>: Selecteer onderste hunk
  <kbd>z</kbd>: Ongedaan maken
  <kbd>M</kbd>: Open external merge tool (git mergetool)
  <kbd>&lt;space&gt;</kbd>: Kies stuk
  <kbd>b</kbd>: Kies beide stukken
  <kbd>&lt;esc&gt;</kbd>: Ga terug naar het bestanden paneel
</pre>

## Normaal

<pre>
  <kbd>mouse wheel down</kbd>: Scroll omlaag (fn+up)
  <kbd>mouse wheel up</kbd>: Scroll omhoog (fn+down)
</pre>

## Patch bouwen

<pre>
  <kbd>&lt;left&gt;</kbd>: Selecteer de vorige hunk
  <kbd>&lt;right&gt;</kbd>: Selecteer de volgende hunk
  <kbd>v</kbd>: Toggle drag selecteer
  <kbd>V</kbd>: Toggle drag selecteer
  <kbd>a</kbd>: Toggle selecteer hunk
  <kbd>&lt;c-o&gt;</kbd>: Copy the selected text to the clipboard
  <kbd>o</kbd>: Open bestand
  <kbd>e</kbd>: Verander bestand
  <kbd>&lt;space&gt;</kbd>: Voeg toe/verwijder lijn(en) in patch
  <kbd>&lt;esc&gt;</kbd>: Sluit lijn-bij-lijn modus
  <kbd>/</kbd>: Start met zoeken
</pre>

## Reflog

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Kopieer commit SHA naar klembord
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Checkout commit
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: Open commit in browser
  <kbd>n</kbd>: Creëer nieuwe branch van commit
  <kbd>g</kbd>: Bekijk reset opties
  <kbd>c</kbd>: Kopieer commit (cherry-pick)
  <kbd>C</kbd>: Kopieer commit reeks (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (gekopieerde) commits selectie
  <kbd>&lt;enter&gt;</kbd>: Bekijk commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Remote branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Kopieer branch name naar klembord
  <kbd>&lt;space&gt;</kbd>: Uitchecken
  <kbd>n</kbd>: Nieuwe branch
  <kbd>M</kbd>: Merge in met huidige checked out branch
  <kbd>r</kbd>: Rebase branch
  <kbd>d</kbd>: Verwijder branch
  <kbd>u</kbd>: Stel in als upstream van uitgecheckte branch
  <kbd>g</kbd>: Bekijk reset opties
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: Bekijk commits
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Remotes

<pre>
  <kbd>f</kbd>: Fetch remote
  <kbd>n</kbd>: Voeg een nieuwe remote toe
  <kbd>d</kbd>: Verwijder remote
  <kbd>e</kbd>: Wijzig remote
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Staging

<pre>
  <kbd>&lt;left&gt;</kbd>: Selecteer de vorige hunk
  <kbd>&lt;right&gt;</kbd>: Selecteer de volgende hunk
  <kbd>v</kbd>: Toggle drag selecteer
  <kbd>V</kbd>: Toggle drag selecteer
  <kbd>a</kbd>: Toggle selecteer hunk
  <kbd>&lt;c-o&gt;</kbd>: Copy the selected text to the clipboard
  <kbd>o</kbd>: Open bestand
  <kbd>e</kbd>: Verander bestand
  <kbd>&lt;esc&gt;</kbd>: Ga terug naar het bestanden paneel
  <kbd>&lt;tab&gt;</kbd>: Ga naar een ander paneel
  <kbd>&lt;space&gt;</kbd>: Toggle lijnen staged / unstaged
  <kbd>d</kbd>: Verwijdert change (git reset)
  <kbd>E</kbd>: Edit hunk
  <kbd>c</kbd>: Commit veranderingen
  <kbd>w</kbd>: Commit veranderingen zonder pre-commit hook
  <kbd>C</kbd>: Commit veranderingen met de git editor
  <kbd>/</kbd>: Start met zoeken
</pre>

## Stash

<pre>
  <kbd>&lt;space&gt;</kbd>: Toepassen
  <kbd>g</kbd>: Pop
  <kbd>d</kbd>: Laten vallen
  <kbd>n</kbd>: Nieuwe branch
  <kbd>r</kbd>: Rename stash
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: Bekijk gecommite bestanden
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Status

<pre>
  <kbd>o</kbd>: Open config bestand
  <kbd>e</kbd>: Verander config bestand
  <kbd>u</kbd>: Check voor updates
  <kbd>&lt;enter&gt;</kbd>: Wissel naar een recente repo
  <kbd>a</kbd>: Alle logs van de branch laten zien
</pre>

## Sub-commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Kopieer commit SHA naar klembord
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Checkout commit
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: Open commit in browser
  <kbd>n</kbd>: Creëer nieuwe branch van commit
  <kbd>g</kbd>: Bekijk reset opties
  <kbd>c</kbd>: Kopieer commit (cherry-pick)
  <kbd>C</kbd>: Kopieer commit reeks (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (gekopieerde) commits selectie
  <kbd>&lt;enter&gt;</kbd>: Bekijk gecommite bestanden
  <kbd>/</kbd>: Start met zoeken
</pre>

## Submodules

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Kopieer submodule naam naar klembord
  <kbd>&lt;enter&gt;</kbd>: Enter submodule
  <kbd>&lt;space&gt;</kbd>: Enter submodule
  <kbd>d</kbd>: Remove submodule
  <kbd>u</kbd>: Update submodule
  <kbd>n</kbd>: Voeg nieuwe submodule toe
  <kbd>e</kbd>: Update submodule URL
  <kbd>i</kbd>: Initialiseer submodule
  <kbd>b</kbd>: Bekijk bulk submodule opties
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Tags

<pre>
  <kbd>&lt;space&gt;</kbd>: Uitchecken
  <kbd>d</kbd>: Verwijder tag
  <kbd>P</kbd>: Push tag
  <kbd>n</kbd>: Creëer tag
  <kbd>g</kbd>: Bekijk reset opties
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: Bekijk commits
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
