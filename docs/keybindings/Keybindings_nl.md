# Lazygit Sneltoetsen

## Globaale Sneltoetsen

<pre>
  <kbd>pgup</kbd>: scroll naar beneden vanaf hooft paneel (fn+up)
  <kbd>pgdown</kbd>: scroll naar beneden vabaf hooft paneel (fn+down)
  <kbd>m</kbd>: bekijk merge/rebase opties
  <kbd>ctrl+p</kbd>: bekijk aangepaste patch opties
  <kbd>P</kbd>: push
  <kbd>p</kbd>: pull
  <kbd>R</kbd>: verversen
  <kbd>x</kbd>: open menu
  <kbd>z</kbd>: ongedaan maken (via reflog) (experimenteel)
  <kbd>ctrl+z</kbd>: redo (via reflog) (experimenteel)
  <kbd>+</kbd>: volgende schermmode (normaal/half/groot )
  <kbd>_</kbd>: vorige schermmode
  <kbd>:</kbd>: voor aangepast commando uit
  <kbd>|</kbd>: bekijk scoping opties
  <kbd>W</kbd>: open diff menu
  <kbd>ctrl+e</kbd>: open diff menu
</pre>

## List Panel Navigation

<pre>
  <kbd>.</kbd>: volgende pagina
  <kbd>,</kbd>: vorige pagina
  <kbd><</kbd>: scroll naar boven
  <kbd>></kbd>: scroll naar beneden
  <kbd>/</kbd>: start met zoekken
  <kbd>]</kbd>: volgende tab
  <kbd>[</kbd>: vorige tab
</pre>

## Branches Paneel (Branches Tab)

<pre>
  <kbd>space</kbd>: uitchecken
  <kbd>o</kbd>: maak een pull-aanvraag
  <kbd>ctrl+y</kbd>: kopieer de URL van het pull-verzoek naar het klembord
  <kbd>c</kbd>: uitchecken bij naam
  <kbd>F</kbd>: forceer checkout
  <kbd>n</kbd>: nieuwe branch
  <kbd>d</kbd>: verwijder branch
  <kbd>r</kbd>: rebase branch
  <kbd>M</kbd>: merge in met huidige checked out branch
  <kbd>i</kbd>: laat git-flow opties zien
  <kbd>f</kbd>: fast-forward deze branch vanaf zijn upstream
  <kbd>g</kbd>: bekijk reset opties
  <kbd>R</kbd>: hernoem branch
  <kbd>ctrl+o</kbd>: copieer branch name naar clipboard
  <kbd>enter</kbd>: view commits
</pre>

## Branches Paneel (Remote Branches (in Remotes tab))

<pre>
  <kbd>esc</kbd>: Ga terug naar remotes lijst
  <kbd>g</kbd>: bekijk reset opties
  <kbd>enter</kbd>: view commits
  <kbd>space</kbd>: uitchecken
  <kbd>n</kbd>: nieuwe branch
  <kbd>M</kbd>: merge in met huidige checked out branch
  <kbd>d</kbd>: verwijder branch
  <kbd>r</kbd>: rebase branch
  <kbd>u</kbd>: stel in als upstream van uitgecheckte branch
</pre>

## Branches Paneel (Remotes Tab)

<pre>
  <kbd>f</kbd>: fetch remote
  <kbd>n</kbd>: voeg een nieuwe remote toe
  <kbd>d</kbd>: verwijder remote
  <kbd>e</kbd>: wijzig remote
</pre>

## Branches Paneel (Sub-commits)

<pre>
  <kbd>enter</kbd>: bekijk gecommite bestanden
  <kbd>space</kbd>: checkout commit
  <kbd>g</kbd>: bekijk reset opties
  <kbd>n</kbd>: nieuwe branch
  <kbd>c</kbd>: kopiëer commit (cherry-pick)
  <kbd>C</kbd>: kopiëer commit reeks (cherry-pick)
  <kbd>ctrl+r</kbd>: reset cherry-picked (gecopieerde) commits selectie
  <kbd>ctrl+o</kbd>: copieer commit SHA naar clipboard
</pre>

## Branches Paneel (Tags Tab)

<pre>
  <kbd>space</kbd>: uitchecken
  <kbd>d</kbd>: verwijder tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: creëer tag
  <kbd>g</kbd>: bekijk reset opties
  <kbd>enter</kbd>: view commits
</pre>

## Commit bestanden Paneel

<pre>
  <kbd>ctrl+o</kbd>: kopieer de vastgelegde bestandsnaam naar het klembord
  <kbd>c</kbd>: bestand uitchecken
  <kbd>d</kbd>: uitsluit deze commit zijn veranderingen aan dit bestand
  <kbd>o</kbd>: open bestand
  <kbd>e</kbd>: verander bestand
  <kbd>space</kbd>: toggle bestand inbegrepen in patch
  <kbd>enter</kbd>: enter bestand to add selecteered lines to the patch
</pre>

## Commits Paneel (Commits)

<pre>
  <kbd>s</kbd>: squash beneden
  <kbd>r</kbd>: hernoem commit
  <kbd>R</kbd>: hernoem commit met editor
  <kbd>g</kbd>: reset naar deze commit
  <kbd>f</kbd>: Fixup commit
  <kbd>F</kbd>: creëer fixup commit voor deze commit
  <kbd>S</kbd>: squash bovenstaande commits
  <kbd>d</kbd>: verwijder commit
  <kbd>ctrl+j</kbd>: verplaats commit 1 naar beneden
  <kbd>ctrl+k</kbd>: verplaats commit 1 naar boven
  <kbd>e</kbd>: wijzig commit
  <kbd>A</kbd>: wijzig commit met staged veranderingen
  <kbd>p</kbd>: kies commit (wanneer midden in rebase)
  <kbd>t</kbd>: commit ongedaan maken
  <kbd>c</kbd>: kopiëer commit (cherry-pick)
  <kbd>ctrl+o</kbd>: copieer commit SHA naar clipboard
  <kbd>C</kbd>: kopiëer commit reeks (cherry-pick)
  <kbd>v</kbd>: plak commits (cherry-pick)
  <kbd>enter</kbd>: bekijk gecommite bestanden
  <kbd>space</kbd>: checkout commit
  <kbd>n</kbd>: create new branch off of commit
  <kbd>T</kbd>: tag commit
  <kbd>ctrl+r</kbd>: reset cherry-picked (gecopieerde) commits selectie
  <kbd>ctrl+y</kbd>: copieer commit bericht naar clipboard
</pre>

## Commits Paneel (Reflog Tab)

<pre>
  <kbd>enter</kbd>: bekijk gecommite bestanden
  <kbd>space</kbd>: checkout commit
  <kbd>g</kbd>: bekijk reset opties
  <kbd>c</kbd>: kopiëer commit (cherry-pick)
  <kbd>C</kbd>: kopiëer commit reeks (cherry-pick)
  <kbd>ctrl+r</kbd>: reset cherry-picked (gecopieerde) commits selectie
  <kbd>ctrl+o</kbd>: copieer commit SHA naar clipboard
</pre>

## Bestanden Paneel (Bestanden)

<pre>
  <kbd>c</kbd>: Commit veranderingen
  <kbd>w</kbd>: commit veranderingen zonder pre-commit hook
  <kbd>A</kbd>: wijzig laatste commit
  <kbd>C</kbd>: commit veranderingen met de git editor
  <kbd>space</kbd>: toggle staged
  <kbd>d</kbd>: bekijk 'veranderingen ongedaan maken' opties
  <kbd>e</kbd>: verander bestand
  <kbd>o</kbd>: open bestand
  <kbd>i</kbd>: voeg toe aan .gitignore
  <kbd>r</kbd>: refresh bestanden
  <kbd>s</kbd>: stash-bestanden
  <kbd>S</kbd>: bekijk stash opties
  <kbd>a</kbd>: toggle staged alle
  <kbd>D</kbd>: bekijk reset opties
  <kbd>enter</kbd>: stage individuele hunks/lijnen
  <kbd>f</kbd>: fetch
  <kbd>ctrl+o</kbd>: kopieer de bestandsnaam naar het klembord
  <kbd>g</kbd>: bekijk upstream reset opties
</pre>

## Bestanden Paneel (Submodules)

<pre>
  <kbd>ctrl+o</kbd>: copy submodule name to clipboard
  <kbd>enter</kbd>: enter submodule
  <kbd>d</kbd>: view reset and remove submodule options
  <kbd>u</kbd>: update submodule
  <kbd>n</kbd>: add new submodule
  <kbd>e</kbd>: update submodule URL
  <kbd>i</kbd>: initialize submodule
  <kbd>b</kbd>: view bulk submodule options
</pre>

## Hooft Paneel (Merggen)

<pre>
  <kbd>esc</kbd>: ga terug naar het bestanden paneel
  <kbd>space</kbd>: kies hunk
  <kbd>b</kbd>: kies bijde hunks
  <kbd>◄</kbd>: selecteer voorgaand conflict
  <kbd>►</kbd>: selecteer volgende conflict
  <kbd>▲</kbd>: selecteer bovenste hunk
  <kbd>▼</kbd>: selecteer onderste hunk
  <kbd>z</kbd>: ongedaan maken
</pre>

## Hooft Paneel (Normaal)

<pre>
  <kbd>￣</kbd>: scroll omlaag (fn+up)
  <kbd>￤</kbd>: scroll omhoog (fn+down)
</pre>

## Hooft Paneel (Patch Bouwen)

<pre>
  <kbd>esc</kbd>: sluit lijn-bij-lijn mode
  <kbd>o</kbd>: open bestand
  <kbd>▲</kbd>: selecteer de vorige lijn
  <kbd>▼</kbd>: selecteer de volgende lijn
  <kbd>◄</kbd>: selecteer de vorige hunk
  <kbd>►</kbd>: selecteer de volgende hunk
  <kbd>space</kbd>: voeg toe/verwijder lijn(en) in patch
  <kbd>v</kbd>: toggle drag  selecteer
  <kbd>V</kbd>: toggle drag  selecteer
  <kbd>a</kbd>: toggle  selecteer hunk
</pre>

## Hooft Paneel (Staging)

<pre>
  <kbd>esc</kbd>: ga terug naar het bestanden paneel
  <kbd>space</kbd>: toggle lijnen staged / unstaged
  <kbd>d</kbd>: verwijdert change (git reset)
  <kbd>tab</kbd>: ga naar een ander paneel
  <kbd>o</kbd>: open bestand
  <kbd>▲</kbd>: selecteer de vorige lijn
  <kbd>▼</kbd>: selecteer de volgende lijn
  <kbd>◄</kbd>: selecteer de vorige hunk
  <kbd>►</kbd>: selecteer de volgende hunk
  <kbd>e</kbd>: verander bestand
  <kbd>o</kbd>: open bestand
  <kbd>v</kbd>: toggle drag  selecteer
  <kbd>V</kbd>: toggle drag  selecteer
  <kbd>a</kbd>: toggle  selecteer hunk
  <kbd>c</kbd>: Commit veranderingen
  <kbd>w</kbd>: commit veranderingen zonder pre-commit hook
  <kbd>C</kbd>: commit veranderingen met de git editor
</pre>

## Menu Paneel

<pre>
  <kbd>esc</kbd>: sluit menu
</pre>

## Stash Paneel

<pre>
  <kbd>enter</kbd>: view stash entry's files
  <kbd>space</kbd>: toepassen
  <kbd>g</kbd>: pop
  <kbd>d</kbd>: laten vallen
  <kbd>n</kbd>: nieuwe branch
</pre>

## Status Paneel

<pre>
  <kbd>e</kbd>: verander config bestand
  <kbd>o</kbd>: open config bestand
  <kbd>u</kbd>: check voor updates
  <kbd>enter</kbd>: wissel naar een recente repo
  <kbd>a</kbd>: alle takken van het houtblok laten zien
</pre>
