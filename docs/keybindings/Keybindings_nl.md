# Lazygit Keybindings

## Global

<pre>
  <kbd>pgup</kbd>: scroll up main panel (fn+up)
  <kbd>pgdown</kbd>: scroll down main panel (fn+down)
  <kbd>m</kbd>: bekijk merge/rebase opties
  <kbd>ctrl+p</kbd>: view custom patch options
  <kbd>P</kbd>: push
  <kbd>p</kbd>: pull
  <kbd>R</kbd>: verversen
  <kbd>x</kbd>: open menu
  <kbd>z</kbd>: undo (via reflog) (experimental)
  <kbd>ctrl+z</kbd>: redo (via reflog) (experimental)
  <kbd>+</kbd>: next screen mode (normal/half/fullscreen)
  <kbd>_</kbd>: prev screen mode
  <kbd>:</kbd>: voor aangepast commando uit
</pre>

## Branches Panel

<pre>
  <kbd>]</kbd>: next tab
  <kbd>[</kbd>: previous tab
</pre>

## Branches Panel (Branches Tab)

<pre>
  <kbd>space</kbd>: uitchecken
  <kbd>o</kbd>: maak een pull-aanvraag
  <kbd>c</kbd>: uitchecken bij naam
  <kbd>F</kbd>: forceer checkout
  <kbd>n</kbd>: nieuwe branch
  <kbd>d</kbd>: verwijder branch
  <kbd>r</kbd>: rebase branch
  <kbd>M</kbd>: merge in met huidige checked out branch
  <kbd>i</kbd>: show git-flow options
  <kbd>f</kbd>: fast-forward this branch from its upstream
  <kbd>g</kbd>: bekijk reset opties
  <kbd>R</kbd>: rename branch
  <kbd>/</kbd>: start search
</pre>

## Branches Panel (Remote Branches (in Remotes tab))

<pre>
  <kbd>esc</kbd>: return to remotes list
  <kbd>g</kbd>: bekijk reset opties
  <kbd>space</kbd>: uitchecken
  <kbd>M</kbd>: merge in met huidige checked out branch
  <kbd>d</kbd>: verwijder branch
  <kbd>r</kbd>: rebase branch
  <kbd>u</kbd>: set as upstream of checked-out branch
  <kbd>/</kbd>: start search
</pre>

## Branches Panel (Remotes Tab)

<pre>
  <kbd>f</kbd>: fetch remote
  <kbd>n</kbd>: add new remote
  <kbd>d</kbd>: remove remote
  <kbd>e</kbd>: edit remote
  <kbd>/</kbd>: start search
</pre>

## Branches Panel (Tags Tab)

<pre>
  <kbd>space</kbd>: uitchecken
  <kbd>d</kbd>: delete tag
  <kbd>P</kbd>: push tag
  <kbd>n</kbd>: create tag
  <kbd>g</kbd>: bekijk reset opties
  <kbd>/</kbd>: start search
</pre>

## Commit bestanden Panel

<pre>
  <kbd>esc</kbd>: ga terug
  <kbd>c</kbd>: bestand uitchecken
  <kbd>d</kbd>: uitsluit deze commit zijn veranderingen aan dit bestand
  <kbd>o</kbd>: open bestand
  <kbd>space</kbd>: toggle file included in patch
  <kbd>enter</kbd>: enter file to add selected lines to the patch
  <kbd>/</kbd>: start search
</pre>

## Commits Panel

<pre>
  <kbd>]</kbd>: next tab
  <kbd>[</kbd>: previous tab
  <kbd>/</kbd>: start search
</pre>

## Commits Panel (Commits Tab)

<pre>
  <kbd>s</kbd>: squash beneden
  <kbd>r</kbd>: hernoem commit
  <kbd>R</kbd>: rename commit with editor
  <kbd>g</kbd>: reset naar deze commit
  <kbd>f</kbd>: Fixup commit
  <kbd>F</kbd>: creëer fixup commit voor deze commit
  <kbd>S</kbd>: squash bovenstaande commits
  <kbd>d</kbd>: verwijder commit
  <kbd>ctrl+j</kbd>: verplaats commit 1 omlaag
  <kbd>ctrl+k</kbd>: verplaats commit 1 omhoog
  <kbd>e</kbd>: verander commit
  <kbd>A</kbd>: wijzig commit met staged veranderingen
  <kbd>p</kbd>: pick commit (when mid-rebase)
  <kbd>t</kbd>: commit omgedaan maken
  <kbd>c</kbd>: kopiëer commit (cherry-pick)
  <kbd>C</kbd>: kopiëer commit reeks (cherry-pick)
  <kbd>v</kbd>: plak commits (cherry-pick)
  <kbd>enter</kbd>: bekijk gecommite bestanden
  <kbd>space</kbd>: checkout commit
  <kbd>i</kbd>: select commit to diff with another commit
  <kbd>T</kbd>: tag commit
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
</pre>

## Commits Panel (Reflog Tab)

<pre>
  <kbd>space</kbd>: checkout commit
  <kbd>g</kbd>: bekijk reset opties
</pre>

## Bestanden Panel

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
  <kbd>S</kbd>: view stash options
  <kbd>a</kbd>: toggle staged alle
  <kbd>D</kbd>: bekijk reset opties
  <kbd>enter</kbd>: stage individuele hunks/lijnen
  <kbd>f</kbd>: fetch
  <kbd>g</kbd>: view upstream reset options
  <kbd>/</kbd>: start search
</pre>

## Hoofd Panel (Merging)

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

## Hoofd Panel (Normaal)

<pre>
  <kbd>￣</kbd>: scroll omlaag (fn+up)
  <kbd>￤</kbd>: scroll omhoog (fn+down)
</pre>

## Hoofd Panel (Patch Building)

<pre>
  <kbd>esc</kbd>: exit line-by-line mode
  <kbd>▲</kbd>: selecteer de vorige lijn
  <kbd>▼</kbd>: selecteer de volgende lijn
  <kbd>◄</kbd>: selecteer de vorige hunk
  <kbd>►</kbd>: selecteer de volgende hunk
  <kbd>space</kbd>: add/remove line(s) to patch
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
</pre>

## Hoofd Panel (Stage Lines/Hunks)

<pre>
  <kbd>esc</kbd>: ga terug naar het bestanden paneel
  <kbd>space</kbd>: toggle line staged / unstaged
  <kbd>d</kbd>: delete change (git reset)
  <kbd>tab</kbd>: switch to other panel
  <kbd>▲</kbd>: selecteer de vorige lijn
  <kbd>▼</kbd>: selecteer de volgende lijn
  <kbd>◄</kbd>: selecteer de vorige hunk
  <kbd>►</kbd>: selecteer de volgende hunk
  <kbd>e</kbd>: verander bestand
  <kbd>o</kbd>: open bestand
  <kbd>v</kbd>: toggle drag select
  <kbd>V</kbd>: toggle drag select
  <kbd>a</kbd>: toggle select hunk
  <kbd>c</kbd>: Commit veranderingen
  <kbd>w</kbd>: commit veranderingen zonder pre-commit hook
  <kbd>C</kbd>: commit veranderingen met de git editor
</pre>

## Menu Panel

<pre>
  <kbd>esc</kbd>: close menu
  <kbd>q</kbd>: close menu
  <kbd>/</kbd>: start search
</pre>

## Stash Panel

<pre>
  <kbd>space</kbd>: toepassen
  <kbd>g</kbd>: pop
  <kbd>d</kbd>: drop
  <kbd>/</kbd>: start search
</pre>

## Status Panel

<pre>
  <kbd>e</kbd>: verander config file
  <kbd>o</kbd>: open config file
  <kbd>u</kbd>: check voor updates
  <kbd>enter</kbd>: wissel naar een recente repo
</pre>
