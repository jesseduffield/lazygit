# Lazygit Sneltoetsen

## Global

<pre>
  <kbd>pgup</kbd>: scroll omhoog naar hooft paneel (fn+up)
  <kbd>pgdown</kbd>: scroll beneden naar hooft paneel (fn+down)
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
  <kbd>|</kbd>: view scoping options
  <kbd>ctrl+e</kbd>: open diff menu
</pre>

## Branches Paneel

<pre>
  <kbd>]</kbd>: volgende tab
  <kbd>[</kbd>: vorige tab
</pre>

## Branches Paneel (Branches Tab)

<pre>
  <kbd>space</kbd>: uitchecken
  <kbd>o</kbd>: maak een pull-aanvraag
  <kbd>c</kbd>: uitchecken bij naam
  <kbd>F</kbd>: forceer checkout
  <kbd>n</kbd>: nieuwe branch
  <kbd>d</kbd>: verwijder branch
  <kbd>r</kbd>: rebase branch
  <kbd>M</kbd>: merge in met huidige uitgecheckte branch
  <kbd>i</kbd>: laat git-flow opties zien
  <kbd>f</kbd>: fast-forward deze branch van zijn upstream
  <kbd>g</kbd>: bekijk reset opties
  <kbd>R</kbd>: hernoem branch
  <kbd>ctrl+o</kbd>: copieer branch name clipboard
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll naar bovenkant
  <kbd>/</kbd>: start met zoekken
  <kbd>></kbd>: scroll naar bodem
</pre>

## Branches Paneel (Remote Branches (in Remotes tab))

<pre>
  <kbd>esc</kbd>: ga terug naar remotes lijst
  <kbd>g</kbd>: bekijk reset opties
  <kbd>space</kbd>: uitchecken
  <kbd>n</kbd>: nieuwe branch
  <kbd>M</kbd>: merge in met huidige uitgecheckte branch
  <kbd>d</kbd>: verwijder branch
  <kbd>r</kbd>: rebase branch
  <kbd>u</kbd>: set as upstream of checked-out branch
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll naar bovenkant
  <kbd>/</kbd>: start met zoekken
  <kbd>></kbd>: scroll naar bodem
</pre>

## Branches Paneel (Remotes Tab)

<pre>
  <kbd>f</kbd>: remote ophalen
  <kbd>n</kbd>: nieuwe remote toevoegen
  <kbd>d</kbd>: verwijder remote
  <kbd>e</kbd>: wijzig remote
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll naar bovenkant
  <kbd>/</kbd>: start met zoekken
  <kbd>></kbd>: scroll naar bodem
</pre>

## Branches Paneel (Tags Tab)

<pre>
  <kbd>space</kbd>: uitchecken
  <kbd>d</kbd>: verwijdert tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: nieuwe tag
  <kbd>g</kbd>: bekijk reset opties
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll naar bovenkant
  <kbd>/</kbd>: start met zoekken
  <kbd>></kbd>: scroll naar bodem
</pre>

## Commit bestanden Paneel

<pre>
  <kbd>esc</kbd>: ga terug
  <kbd>c</kbd>: bestand uitchecken
  <kbd>d</kbd>: uitsluit deze commit zijn wijzigingen aan dit bestand
  <kbd>o</kbd>: open bestand
  <kbd>e</kbd>: verander bestand
  <kbd>space</kbd>: wissel bestand opgenomen in patch
  <kbd>enter</kbd>: open bestand om specifieke lijnen toe te voegen aan patch
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll naar bovenkant
  <kbd>/</kbd>: start met zoekken
  <kbd>></kbd>: scroll naar bodem
</pre>

## Commits Paneel

<pre>
  <kbd>]</kbd>: volgende tab
  <kbd>[</kbd>: vorige tab
</pre>

## Commits Paneel (Commits Tab)

<pre>
  <kbd>s</kbd>: squash beneden
  <kbd>r</kbd>: hernoem commit
  <kbd>R</kbd>: hernoem commit met editor
  <kbd>g</kbd>: reset naar deze commit
  <kbd>f</kbd>: Fixup commit
  <kbd>F</kbd>: creëer fixup commit voor deze commit
  <kbd>S</kbd>: squash bovenstaande commits
  <kbd>d</kbd>: verwijder commit
  <kbd>ctrl+j</kbd>: verplaats commit 1 omlaag
  <kbd>ctrl+k</kbd>: verplaats commit 1 omhoog
  <kbd>e</kbd>: verander commit
  <kbd>A</kbd>: wijzig commit met staged wijzigingen
  <kbd>p</kbd>: pick commit (wanneer midden in rebase)
  <kbd>t</kbd>: commit omgedaan maken
  <kbd>c</kbd>: kopiëer commit (cherry-pick)
  <kbd>ctrl+o</kbd>: kopiëer commit SHA naar clipboard
  <kbd>C</kbd>: kopiëer commit reeks (cherry-pick)
  <kbd>v</kbd>: plak commits (cherry-pick)
  <kbd>enter</kbd>: bekijk gecommite bestanden
  <kbd>space</kbd>: checkout commit
  <kbd>T</kbd>: tag commit
  <kbd>ctrl+r</kbd>: reset cherry-picked (gecopieerde) commits selectie
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll naar bovenkant
  <kbd>/</kbd>: start met zoekken
  <kbd>></kbd>: scroll naar bodem
</pre>

## Commits Paneel (Reflog Tab)

<pre>
  <kbd>space</kbd>: checkout commit
  <kbd>g</kbd>: bekijk reset opties
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll naar bovenkant
  <kbd>/</kbd>: start met zoekken
  <kbd>></kbd>: scroll naar bodem
</pre>

## Bestanden Paneel

<pre>
  <kbd>c</kbd>: commit wijzigingen
  <kbd>w</kbd>: commit wijzigingen zonder pre-commit hook
  <kbd>A</kbd>: wijzig laatste commit
  <kbd>C</kbd>: commit wijzigingen met de git editor
  <kbd>space</kbd>: wissel staged
  <kbd>d</kbd>: bekijk 'wijzigingen ongedaan maken' opties
  <kbd>e</kbd>: verander bestand
  <kbd>o</kbd>: open bestand
  <kbd>i</kbd>: voeg toe aan .gitignore
  <kbd>r</kbd>: bestanden vernieuwen
  <kbd>s</kbd>: stash-bestanden
  <kbd>S</kbd>: bekijk stash opties
  <kbd>a</kbd>: wissel staged alle
  <kbd>D</kbd>: bekijk reset opties
  <kbd>enter</kbd>: stage individuele hunks/lijnen
  <kbd>f</kbd>: fetch
  <kbd>g</kbd>: bekijk upstream reset optie
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll naar bovenkant
  <kbd>/</kbd>: start met zoekken
  <kbd>></kbd>: scroll naar bodem
</pre>

## Hoofd Paneel (Merging)

<pre>
  <kbd>esc</kbd>: ga terug naar het bestanden paneel
  <kbd>space</kbd>: pick hunk
  <kbd>b</kbd>: pick beide hunks
  <kbd>◄</kbd>: selecteer voorgaand conflict
  <kbd>►</kbd>: selecteer volgende conflict
  <kbd>▲</kbd>: selecteer bovenste hunk
  <kbd>▼</kbd>: selecteer onderste hunk
  <kbd>z</kbd>: ongedaan maken
</pre>

## Hoofd Paneel (Normaal)

<pre>
  <kbd>￣</kbd>: scroll omlaag (fn+up)
  <kbd>￤</kbd>: scroll omhoog (fn+down)
</pre>

## Hoofd Paneel (Patch Building)

<pre>
  <kbd>esc</kbd>: sluit lijn-bij-lijn mode
  <kbd>▲</kbd>: selecteer de vorige lijn
  <kbd>▼</kbd>: selecteer de volgende lijn
  <kbd>◄</kbd>: selecteer de vorige hunk
  <kbd>►</kbd>: selecteer de volgende hunk
  <kbd>space</kbd>: voeg toe/verwijdert lijn(en) in patch
  <kbd>v</kbd>: wissel drag selectie
  <kbd>V</kbd>: wissel drag selectie
  <kbd>a</kbd>: wissel selectie hunk
</pre>

## Hoofd Paneel (Stage Lines/Hunks)

<pre>
  <kbd>esc</kbd>: ga terug naar het bestanden paneel
  <kbd>space</kbd>: wissel lijn staged / unstaged
  <kbd>d</kbd>: verwijdert wijziging (git reset)
  <kbd>tab</kbd>: ga naar ander paneel
  <kbd>▲</kbd>: selecteer de vorige lijn
  <kbd>▼</kbd>: selecteer de volgende lijn
  <kbd>◄</kbd>: selecteer de vorige hunk
  <kbd>►</kbd>: selecteer de volgende hunk
  <kbd>e</kbd>: verander bestand
  <kbd>o</kbd>: open bestand
  <kbd>v</kbd>: wissel drag selectie
  <kbd>V</kbd>: wissel drag selectie
  <kbd>a</kbd>: wissel selectie hunk
  <kbd>c</kbd>: commit wijzigingen
  <kbd>w</kbd>: commit wijzigingen zonder pre-commit hook
  <kbd>C</kbd>: commit wijzigingen met de git editor
</pre>

## Menu Paneel

<pre>
  <kbd>esc</kbd>: sluit menu
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll to top
  <kbd>/</kbd>: start search
  <kbd>></kbd>: scroll to bottom
</pre>

## Stash Paneel

<pre>
  <kbd>space</kbd>: toepassen
  <kbd>g</kbd>: pop
  <kbd>d</kbd>: drop
  <kbd>,</kbd>: vorige pagina
  <kbd>.</kbd>: volgende pagina
  <kbd><</kbd>: scroll naar bovenkant
  <kbd>/</kbd>: start met zoekken
  <kbd>></kbd>: scroll naar bodem
</pre>

## Status Paneel

<pre>
  <kbd>e</kbd>: verander config file
  <kbd>o</kbd>: open config bestand
  <kbd>u</kbd>: check voor updates
  <kbd>enter</kbd>: wissel naar een recente repo
</pre>
