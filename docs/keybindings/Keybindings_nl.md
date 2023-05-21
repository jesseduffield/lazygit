_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit Sneltoetsen

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Globale Sneltoetsen

<pre>
  <kbd>&lt;c-r&gt;</kbd>: wissel naar een recente repo
  <kbd>&lt;pgup&gt;</kbd>: scroll naar beneden vanaf hoofdpaneel (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: scroll naar beneden vanaf hoofdpaneel (fn+down/shift+j)
  <kbd>@</kbd>: open command log menu
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: voer aangepaste commando uit
  <kbd>&lt;c-p&gt;</kbd>: bekijk aangepaste patch opties
  <kbd>m</kbd>: bekijk merge/rebase opties
  <kbd>R</kbd>: verversen
  <kbd>+</kbd>: volgende scherm modus (normaal/half/groot)
  <kbd>_</kbd>: vorige scherm modus
  <kbd>?</kbd>: open menu
  <kbd>&lt;c-s&gt;</kbd>: bekijk scoping opties
  <kbd>W</kbd>: open diff menu
  <kbd>&lt;c-e&gt;</kbd>: open diff menu
  <kbd>&lt;c-w&gt;</kbd>: Toggle whether or not whitespace changes are shown in the diff view
  <kbd>z</kbd>: ongedaan maken (via reflog) (experimenteel)
  <kbd>&lt;c-z&gt;</kbd>: redo (via reflog) (experimenteel)
  <kbd>P</kbd>: push
  <kbd>p</kbd>: pull
</pre>

## Lijstpaneel Navigatie

<pre>
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd>&lt;</kbd>: scroll naar boven
  <kbd>/</kbd>: start met zoeken
  <kbd>&gt;</kbd>: scroll naar beneden
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
  <kbd>]</kbd>: volgende tabblad
  <kbd>[</kbd>: vorige tabblad
</pre>

## Bestanden

<pre>
  <kbd>&lt;c-o&gt;</kbd>: kopieer de bestandsnaam naar het klembord
  <kbd>d</kbd>: bekijk 'veranderingen ongedaan maken' opties
  <kbd>&lt;space&gt;</kbd>: toggle staged
  <kbd>&lt;c-b&gt;</kbd>: Filter files (staged/unstaged)
  <kbd>c</kbd>: commit veranderingen
  <kbd>w</kbd>: commit veranderingen zonder pre-commit hook
  <kbd>A</kbd>: wijzig laatste commit
  <kbd>C</kbd>: commit veranderingen met de git editor
  <kbd>e</kbd>: verander bestand
  <kbd>o</kbd>: open bestand
  <kbd>i</kbd>: ignore or exclude file
  <kbd>r</kbd>: refresh bestanden
  <kbd>s</kbd>: stash-bestanden
  <kbd>S</kbd>: bekijk stash opties
  <kbd>a</kbd>: toggle staged alle
  <kbd>&lt;enter&gt;</kbd>: stage individuele hunks/lijnen
  <kbd>g</kbd>: bekijk upstream reset opties
  <kbd>D</kbd>: bekijk reset opties
  <kbd>`</kbd>: toggle bestandsboom weergave
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>f</kbd>: fetch
</pre>

## Bevestigingspaneel

<pre>
  <kbd>&lt;enter&gt;</kbd>: bevestig
  <kbd>&lt;esc&gt;</kbd>: sluiten
</pre>

## Branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: kopieer branch name naar klembord
  <kbd>i</kbd>: laat git-flow opties zien
  <kbd>&lt;space&gt;</kbd>: uitchecken
  <kbd>n</kbd>: nieuwe branch
  <kbd>o</kbd>: maak een pull-request
  <kbd>O</kbd>: bekijk opties voor pull-aanvraag
  <kbd>&lt;c-y&gt;</kbd>: kopieer de URL van het pull-verzoek naar het klembord
  <kbd>c</kbd>: uitchecken bij naam
  <kbd>F</kbd>: forceer checkout
  <kbd>d</kbd>: verwijder branch
  <kbd>r</kbd>: rebase branch
  <kbd>M</kbd>: merge in met huidige checked out branch
  <kbd>f</kbd>: fast-forward deze branch vanaf zijn upstream
  <kbd>T</kbd>: creëer tag
  <kbd>g</kbd>: bekijk reset opties
  <kbd>R</kbd>: hernoem branch
  <kbd>u</kbd>: set/unset upstream
  <kbd>&lt;enter&gt;</kbd>: bekijk commits
</pre>

## Commit Bericht

<pre>
  <kbd>&lt;enter&gt;</kbd>: bevestig
  <kbd>&lt;esc&gt;</kbd>: sluiten
</pre>

## Commit bestanden

<pre>
  <kbd>&lt;c-o&gt;</kbd>: kopieer de vastgelegde bestandsnaam naar het klembord
  <kbd>c</kbd>: bestand uitchecken
  <kbd>d</kbd>: uitsluit deze commit zijn veranderingen aan dit bestand
  <kbd>o</kbd>: open bestand
  <kbd>e</kbd>: verander bestand
  <kbd>&lt;space&gt;</kbd>: toggle bestand inbegrepen in patch
  <kbd>a</kbd>: toggle all files included in patch
  <kbd>&lt;enter&gt;</kbd>: enter bestand om geselecteerde regels toe te voegen aan de patch
  <kbd>`</kbd>: toggle bestandsboom weergave
</pre>

## Commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: kopieer commit SHA naar klembord
  <kbd>&lt;c-r&gt;</kbd>: reset cherry-picked (gekopieerde) commits selectie
  <kbd>b</kbd>: view bisect options
  <kbd>s</kbd>: squash beneden
  <kbd>f</kbd>: Fixup commit
  <kbd>r</kbd>: hernoem commit
  <kbd>R</kbd>: hernoem commit met editor
  <kbd>d</kbd>: verwijder commit
  <kbd>e</kbd>: wijzig commit
  <kbd>p</kbd>: kies commit (wanneer midden in rebase)
  <kbd>F</kbd>: creëer fixup commit voor deze commit
  <kbd>S</kbd>: squash bovenstaande commits
  <kbd>&lt;c-j&gt;</kbd>: verplaats commit 1 naar beneden
  <kbd>&lt;c-k&gt;</kbd>: verplaats commit 1 naar boven
  <kbd>v</kbd>: plak commits (cherry-pick)
  <kbd>A</kbd>: wijzig commit met staged veranderingen
  <kbd>a</kbd>: reset commit author
  <kbd>t</kbd>: commit ongedaan maken
  <kbd>T</kbd>: tag commit
  <kbd>&lt;c-l&gt;</kbd>: open log menu
  <kbd>&lt;space&gt;</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: creëer nieuwe branch van commit
  <kbd>g</kbd>: bekijk reset opties
  <kbd>c</kbd>: kopieer commit (cherry-pick)
  <kbd>C</kbd>: kopieer commit reeks (cherry-pick)
  <kbd>&lt;enter&gt;</kbd>: bekijk gecommite bestanden
</pre>

## Menu

<pre>
  <kbd>&lt;enter&gt;</kbd>: uitvoeren
  <kbd>&lt;esc&gt;</kbd>: sluiten
</pre>

## Mergen

<pre>
  <kbd>e</kbd>: verander bestand
  <kbd>o</kbd>: open bestand
  <kbd>&lt;left&gt;</kbd>: selecteer voorgaand conflict
  <kbd>&lt;right&gt;</kbd>: selecteer volgende conflict
  <kbd>&lt;up&gt;</kbd>: selecteer bovenste hunk
  <kbd>&lt;down&gt;</kbd>: selecteer onderste hunk
  <kbd>z</kbd>: ongedaan maken
  <kbd>M</kbd>: open external merge tool (git mergetool)
  <kbd>&lt;space&gt;</kbd>: kies hunk
  <kbd>b</kbd>: kies bijde hunks
  <kbd>&lt;esc&gt;</kbd>: ga terug naar het bestanden paneel
</pre>

## Normaal

<pre>
  <kbd>mouse wheel down</kbd>: scroll omlaag (fn+up)
  <kbd>mouse wheel up</kbd>: scroll omhoog (fn+down)
</pre>

## Patch Bouwen

<pre>
  <kbd>&lt;left&gt;</kbd>: selecteer de vorige hunk
  <kbd>&lt;right&gt;</kbd>: selecteer de volgende hunk
  <kbd>v</kbd>: toggle drag selecteer
  <kbd>V</kbd>: toggle drag selecteer
  <kbd>a</kbd>: toggle selecteer hunk
  <kbd>&lt;c-o&gt;</kbd>: copy the selected text to the clipboard
  <kbd>o</kbd>: open bestand
  <kbd>e</kbd>: verander bestand
  <kbd>&lt;space&gt;</kbd>: voeg toe/verwijder lijn(en) in patch
  <kbd>&lt;esc&gt;</kbd>: sluit lijn-bij-lijn modus
</pre>

## Reflog

<pre>
  <kbd>&lt;c-o&gt;</kbd>: kopieer commit SHA naar klembord
  <kbd>&lt;space&gt;</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: creëer nieuwe branch van commit
  <kbd>g</kbd>: bekijk reset opties
  <kbd>c</kbd>: kopieer commit (cherry-pick)
  <kbd>C</kbd>: kopieer commit reeks (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: reset cherry-picked (gekopieerde) commits selectie
  <kbd>&lt;enter&gt;</kbd>: bekijk commits
</pre>

## Remote Branches

<pre>
  <kbd>&lt;c-o&gt;</kbd>: kopieer branch name naar klembord
  <kbd>&lt;space&gt;</kbd>: uitchecken
  <kbd>n</kbd>: nieuwe branch
  <kbd>M</kbd>: merge in met huidige checked out branch
  <kbd>r</kbd>: rebase branch
  <kbd>d</kbd>: verwijder branch
  <kbd>u</kbd>: stel in als upstream van uitgecheckte branch
  <kbd>&lt;esc&gt;</kbd>: ga terug naar remotes lijst
  <kbd>g</kbd>: bekijk reset opties
  <kbd>&lt;enter&gt;</kbd>: bekijk commits
</pre>

## Remotes

<pre>
  <kbd>f</kbd>: fetch remote
  <kbd>n</kbd>: voeg een nieuwe remote toe
  <kbd>d</kbd>: verwijder remote
  <kbd>e</kbd>: wijzig remote
</pre>

## Staging

<pre>
  <kbd>&lt;left&gt;</kbd>: selecteer de vorige hunk
  <kbd>&lt;right&gt;</kbd>: selecteer de volgende hunk
  <kbd>v</kbd>: toggle drag selecteer
  <kbd>V</kbd>: toggle drag selecteer
  <kbd>a</kbd>: toggle selecteer hunk
  <kbd>&lt;c-o&gt;</kbd>: copy the selected text to the clipboard
  <kbd>o</kbd>: open bestand
  <kbd>e</kbd>: verander bestand
  <kbd>&lt;esc&gt;</kbd>: ga terug naar het bestanden paneel
  <kbd>&lt;tab&gt;</kbd>: ga naar een ander paneel
  <kbd>&lt;space&gt;</kbd>: toggle lijnen staged / unstaged
  <kbd>d</kbd>: verwijdert change (git reset)
  <kbd>E</kbd>: edit hunk
  <kbd>c</kbd>: commit veranderingen
  <kbd>w</kbd>: commit veranderingen zonder pre-commit hook
  <kbd>C</kbd>: commit veranderingen met de git editor
</pre>

## Stash

<pre>
  <kbd>&lt;space&gt;</kbd>: toepassen
  <kbd>g</kbd>: pop
  <kbd>d</kbd>: laten vallen
  <kbd>n</kbd>: nieuwe branch
  <kbd>r</kbd>: rename stash
  <kbd>&lt;enter&gt;</kbd>: bekijk gecommite bestanden
</pre>

## Status

<pre>
  <kbd>o</kbd>: open config bestand
  <kbd>e</kbd>: verander config bestand
  <kbd>u</kbd>: check voor updates
  <kbd>&lt;enter&gt;</kbd>: wissel naar een recente repo
  <kbd>a</kbd>: alle logs van de branch laten zien
</pre>

## Sub-commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: kopieer commit SHA naar klembord
  <kbd>&lt;space&gt;</kbd>: checkout commit
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: open commit in browser
  <kbd>n</kbd>: creëer nieuwe branch van commit
  <kbd>g</kbd>: bekijk reset opties
  <kbd>c</kbd>: kopieer commit (cherry-pick)
  <kbd>C</kbd>: kopieer commit reeks (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: reset cherry-picked (gekopieerde) commits selectie
  <kbd>&lt;enter&gt;</kbd>: bekijk gecommite bestanden
</pre>

## Submodules

<pre>
  <kbd>&lt;c-o&gt;</kbd>: kopieer submodule naam naar klembord
  <kbd>&lt;enter&gt;</kbd>: enter submodule
  <kbd>d</kbd>: remove submodule
  <kbd>u</kbd>: update submodule
  <kbd>n</kbd>: voeg nieuwe submodule toe
  <kbd>e</kbd>: update submodule URL
  <kbd>i</kbd>: initialiseer submodule
  <kbd>b</kbd>: bekijk bulk submodule opties
</pre>

## Tags

<pre>
  <kbd>&lt;space&gt;</kbd>: uitchecken
  <kbd>d</kbd>: verwijder tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: creëer tag
  <kbd>g</kbd>: bekijk reset opties
  <kbd>&lt;enter&gt;</kbd>: bekijk commits
</pre>
