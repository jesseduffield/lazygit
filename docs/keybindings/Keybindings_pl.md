_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Skróty klawiszowe

_Legenda: `<c-b>` oznacza ctrl+b, `<a-b>` oznacza alt+b, `B` oznacza shift+b_

## Globalne skróty klawiszowe

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Przełącz na ostatnie repozytorium |  |
| `` <pgup> (fn+up/shift+k) `` | Przewiń główne okno w górę |  |
| `` <pgdown> (fn+down/shift+j) `` | Przewiń główne okno w dół |  |
| `` @ `` | Pokaż opcje dziennika poleceń | Pokaż opcje dla dziennika poleceń, np. pokazywanie/ukrywanie dziennika poleceń i skupienie na dzienniku poleceń. |
| `` P `` | Wypchnij | Wypchnij bieżącą gałąź do jej gałęzi nadrzędnej. Jeśli nie skonfigurowano gałęzi nadrzędnej, zostaniesz poproszony o skonfigurowanie gałęzi nadrzędnej. |
| `` p `` | Pociągnij | Pociągnij zmiany z zdalnego dla bieżącej gałęzi. Jeśli nie skonfigurowano gałęzi nadrzędnej, zostaniesz poproszony o skonfigurowanie gałęzi nadrzędnej. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` } `` | Zwiększ rozmiar kontekstu w widoku różnic | Zwiększ ilość kontekstu pokazywanego wokół zmian w widoku różnic. |
| `` { `` | Zmniejsz rozmiar kontekstu w widoku różnic | Zmniejsz ilość kontekstu pokazywanego wokół zmian w widoku różnic. |
| `` : `` | Wykonaj polecenie niestandardowe | Wyświetl monit, w którym możesz wprowadzić polecenie powłoki do wykonania. Nie należy mylić z wcześniej skonfigurowanymi poleceniami niestandardowymi. |
| `` <c-p> `` | Wyświetl opcje niestandardowej łatki |  |
| `` m `` | Pokaż opcje scalania/rebase | Pokaż opcje do przerwania/kontynuowania/pominięcia bieżącego scalania/rebase. |
| `` R `` | Odśwież | Odśwież stan git (tj. uruchom `git status`, `git branch`, itp. w tle, aby zaktualizować zawartość paneli). To nie uruchamia `git fetch`. |
| `` + `` | Następny tryb ekranu (normalny/półpełny/pełnoekranowy) |  |
| `` _ `` | Poprzedni tryb ekranu |  |
| `` ? `` | Otwórz menu przypisań klawiszy |  |
| `` <c-s> `` | Pokaż opcje filtrowania | Pokaż opcje filtrowania dziennika commitów, tak aby pokazywane były tylko commity pasujące do filtra. |
| `` W `` | Pokaż opcje różnicowania | Pokaż opcje dotyczące różnicowania dwóch refów, np. różnicowanie względem wybranego refa, wprowadzanie refa do różnicowania i odwracanie kierunku różnic. |
| `` <c-e> `` | Pokaż opcje różnicowania | Pokaż opcje dotyczące różnicowania dwóch refów, np. różnicowanie względem wybranego refa, wprowadzanie refa do różnicowania i odwracanie kierunku różnic. |
| `` q `` | Wyjdź |  |
| `` <esc> `` | Anuluj |  |
| `` <c-w> `` | Przełącz białe znaki | Przełącz czy zmiany białych znaków są pokazywane w widoku różnic. |
| `` z `` | Cofnij | Dziennik reflog zostanie użyty do określenia, jakie polecenie git należy uruchomić, aby cofnąć ostatnie polecenie git. Nie obejmuje to zmian w drzewie roboczym; brane są pod uwagę tylko commity. |
| `` <c-z> `` | Ponów | Dziennik reflog zostanie użyty do określenia, jakie polecenie git należy uruchomić, aby ponowić ostatnie polecenie git. Nie obejmuje to zmian w drzewie roboczym; brane są pod uwagę tylko commity. |

## Nawigacja panelu listy

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Poprzednia strona |  |
| `` . `` | Następna strona |  |
| `` < `` | Przewiń do góry |  |
| `` > `` | Przewiń do dołu |  |
| `` v `` | Przełącz zaznaczenie zakresu |  |
| `` <s-down> `` | Zaznacz zakres w dół |  |
| `` <s-up> `` | Zaznacz zakres w górę |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |
| `` H `` | Przewiń w lewo |  |
| `` L `` | Przewiń w prawo |  |
| `` ] `` | Następna zakładka |  |
| `` [ `` | Poprzednia zakładka |  |

## Commity

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopiuj hash commita do schowka |  |
| `` <c-r> `` | Resetuj wybrane (cherry-picked) commity |  |
| `` b `` | Zobacz opcje bisect |  |
| `` s `` | Scal | Scal wybrany commit z commitami poniżej. Wiadomość wybranego commita zostanie dołączona do commita poniżej. |
| `` f `` | Poprawka | Włącz wybrany commit do commita poniżej. Podobnie do fixup, ale wiadomość wybranego commita zostanie odrzucona. |
| `` r `` | Przeformułuj | Przeformułuj wiadomość wybranego commita. |
| `` R `` | Przeformułuj za pomocą edytora |  |
| `` d `` | Usuń | Usuń wybrany commit. To usunie commit z gałęzi za pomocą rebazowania. Jeśli commit wprowadza zmiany, od których zależą późniejsze commity, być może będziesz musiał rozwiązać konflikty scalania. |
| `` e `` | Edytuj (rozpocznij interaktywne rebazowanie) | Edytuj wybrany commit. Użyj tego, aby rozpocząć interaktywne rebazowanie od wybranego commita. Podczas trwania rebazowania, to oznaczy wybrany commit do edycji, co oznacza, że po kontynuacji rebazowania, rebazowanie zostanie wstrzymane na wybranym commicie, aby umożliwić wprowadzenie zmian. |
| `` i `` | Rozpocznij interaktywny rebase | Rozpocznij interaktywny rebase dla commitów na twoim branchu. To będzie zawierać wszystkie commity od HEAD do pierwszego commita scalenia lub commita głównego brancha.
Jeśli chcesz zamiast tego rozpocząć interaktywny rebase od wybranego commita, naciśnij `e`. |
| `` p `` | Wybierz | Oznacz wybrany commit do wybrania (podczas rebazowania). Oznacza to, że commit zostanie zachowany po kontynuacji rebazowania. |
| `` F `` | Utwórz commit fixup | Utwórz commit 'fixup!' dla wybranego commita. Później możesz nacisnąć `S` na tym samym commicie, aby zastosować wszystkie powyższe commity fixup. |
| `` S `` | Zastosuj commity fixup | Scal wszystkie commity 'fixup!', albo powyżej wybranego commita, albo wszystkie w bieżącej gałęzi (autosquash). |
| `` <c-j> `` | Przesuń commit w dół |  |
| `` <c-k> `` | Przesuń commit w górę |  |
| `` V `` | Wklej (cherry-pick) |  |
| `` B `` | Oznacz jako bazowy commit dla rebase | Wybierz bazowy commit dla następnego rebase. Kiedy robisz rebase na branch, tylko commity powyżej bazowego commita zostaną przeniesione. Używa to polecenia `git rebase --onto`. |
| `` A `` | Popraw | Popraw commit ze zmianami zatwierdzonymi. Jeśli wybrany commit jest commit HEAD, to wykona `git commit --amend`. W przeciwnym razie commit zostanie poprawiony za pomocą rebazowania. |
| `` a `` | Popraw atrybut commita | Ustaw/Resetuj autora commita lub ustaw współautora. |
| `` t `` | Cofnij | Utwórz commit cofający dla wybranego commita, który stosuje zmiany wybranego commita w odwrotnej kolejności. |
| `` T `` | Otaguj commit | Utwórz nowy tag wskazujący na wybrany commit. Zostaniesz poproszony o wprowadzenie nazwy tagu i opcjonalnego opisu. |
| `` <c-l> `` | Zobacz opcje logów | Zobacz opcje dla logów commitów, np. zmiana kolejności sortowania, ukrywanie grafu gita, pokazywanie całego grafu gita. |
| `` <space> `` | Przełącz | Przełącz wybrany commit jako odłączoną HEAD. |
| `` y `` | Kopiuj atrybut commita do schowka | Kopiuj atrybut commita do schowka (np. hash, URL, różnice, wiadomość, autor). |
| `` o `` | Otwórz commit w przeglądarce |  |
| `` n `` | Utwórz nową gałąź z commita |  |
| `` g `` | Reset | Wyświetl opcje resetu (miękki/mieszany/twardy) do wybranego elementu. |
| `` C `` | Kopiuj (cherry-pick) | Oznacz commit jako skopiowany. Następnie, w widoku lokalnych commitów, możesz nacisnąć `V`, aby wkleić (cherry-pick) skopiowane commity do sprawdzonej gałęzi. W dowolnym momencie możesz nacisnąć `<esc>`, aby anulować zaznaczenie. |
| `` <c-t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` <enter> `` | Wyświetl pliki |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Drzewa pracy

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | Nowe drzewo pracy |  |
| `` <space> `` | Przełącz | Przełącz do wybranego drzewa pracy. |
| `` o `` | Otwórz w edytorze |  |
| `` d `` | Usuń | Usuń wybrane drzewo pracy. To usunie zarówno katalog drzewa pracy, jak i metadane o drzewie pracy w katalogu .git. |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Główny panel (budowanie łatki)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Idź do poprzedniego fragmentu |  |
| `` <right> `` | Idź do następnego fragmentu |  |
| `` v `` | Przełącz zaznaczenie zakresu |  |
| `` a `` | Zaznacz fragment | Przełącz tryb zaznaczania fragmentu. |
| `` <c-o> `` | Kopiuj zaznaczony tekst do schowka |  |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` e `` | Edytuj plik | Otwórz plik w zewnętrznym edytorze. |
| `` <space> `` | Przełącz linie w łatce |  |
| `` <esc> `` | Wyjdź z budowniczego niestandardowej łatki |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Lokalne gałęzie

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopiuj nazwę gałęzi do schowka |  |
| `` i `` | Pokaż opcje git-flow |  |
| `` <space> `` | Przełącz | Przełącz wybrany element. |
| `` n `` | Nowa gałąź |  |
| `` o `` | Utwórz żądanie ściągnięcia |  |
| `` O `` | Zobacz opcje tworzenia pull requesta |  |
| `` <c-y> `` | Kopiuj adres URL żądania ściągnięcia do schowka |  |
| `` c `` | Przełącz według nazwy | Przełącz według nazwy. W polu wprowadzania możesz wpisać '-' aby przełączyć się na ostatnią gałąź. |
| `` F `` | Wymuś przełączenie | Wymuś przełączenie wybranej gałęzi. To spowoduje odrzucenie wszystkich lokalnych zmian w drzewie roboczym przed przełączeniem na wybraną gałąź. |
| `` d `` | Usuń | Wyświetl opcje usuwania lokalnej/odległej gałęzi. |
| `` r `` | Przebazuj | Przebazuj przełączoną gałąź na wybraną gałąź. |
| `` M `` | Scal | Scal wybraną gałąź z aktualnie sprawdzoną gałęzią. |
| `` f `` | Szybkie przewijanie | Szybkie przewijanie wybranej gałęzi z jej źródła. |
| `` T `` | Nowy tag |  |
| `` s `` | Kolejność sortowania |  |
| `` g `` | Reset |  |
| `` R `` | Zmień nazwę gałęzi |  |
| `` u `` | Pokaż opcje upstream | Pokaż opcje dotyczące upstream gałęzi, np. ustawianie/usuwanie upstream i resetowanie do upstream. |
| `` <c-t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` <enter> `` | Pokaż commity |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Menu

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Wykonaj |  |
| `` <esc> `` | Zamknij |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Panel główny (normalny)

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | Przewiń w dół |  |
| `` mouse wheel up (fn+down) `` | Przewiń w górę |  |

## Panel główny (scalanie)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Wybierz fragment |  |
| `` b `` | Wybierz wszystkie fragmenty |  |
| `` <up> `` | Poprzedni fragment |  |
| `` <down> `` | Następny fragment |  |
| `` <left> `` | Poprzedni konflikt |  |
| `` <right> `` | Następny konflikt |  |
| `` z `` | Cofnij | Cofnij ostatnie rozwiązanie konfliktu scalania. |
| `` e `` | Edytuj plik | Otwórz plik w zewnętrznym edytorze. |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` M `` | Otwórz zewnętrzne narzędzie scalania | Uruchom `git mergetool`. |
| `` <esc> `` | Wróć do panelu plików |  |

## Panel główny (zatwierdzanie)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Idź do poprzedniego fragmentu |  |
| `` <right> `` | Idź do następnego fragmentu |  |
| `` v `` | Przełącz zaznaczenie zakresu |  |
| `` a `` | Zaznacz fragment | Przełącz tryb zaznaczania fragmentu. |
| `` <c-o> `` | Kopiuj zaznaczony tekst do schowka |  |
| `` <space> `` | Zatwierdź | Przełącz zaznaczenie zatwierdzone/niezatwierdzone. |
| `` d `` | Odrzuć | Gdy zaznaczona jest niezatwierdzona zmiana, odrzuć ją używając `git reset`. Gdy zaznaczona jest zatwierdzona zmiana, cofnij zatwierdzenie. |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` e `` | Edytuj plik | Otwórz plik w zewnętrznym edytorze. |
| `` <esc> `` | Wróć do panelu plików |  |
| `` <tab> `` | Przełącz widok | Przełącz na inny widok (zatwierdzone/niezatwierdzone zmiany). |
| `` E `` | Edytuj fragment | Edytuj wybrany fragment w zewnętrznym edytorze. |
| `` c `` | Commit | Zatwierdź zmiany zatwierdzone. |
| `` w `` | Zatwierdź zmiany bez hooka pre-commit |  |
| `` C `` | Zatwierdź zmiany używając edytora git |  |
| `` <c-f> `` | Znajdź bazowy commit do poprawki | Znajdź commit, na którym opierają się Twoje obecne zmiany, w celu poprawienia/zmiany commita. To pozwala Ci uniknąć przeglądania commitów w Twojej gałęzi jeden po drugim, aby zobaczyć, który commit powinien być poprawiony/zmieniony. Zobacz dokumentację: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Panel potwierdzenia

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Potwierdź |  |
| `` <esc> `` | Zamknij/Anuluj |  |

## Pliki

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopiuj ścieżkę do schowka |  |
| `` <space> `` | Zatwierdź | Przełącz zatwierdzenie dla wybranego pliku. |
| `` <c-b> `` | Filtruj pliki według statusu |  |
| `` y `` | Kopiuj do schowka |  |
| `` c `` | Commit | Zatwierdź zmiany zatwierdzone. |
| `` w `` | Zatwierdź zmiany bez hooka pre-commit |  |
| `` A `` | Popraw ostatni commit |  |
| `` C `` | Zatwierdź zmiany używając edytora git |  |
| `` <c-f> `` | Znajdź bazowy commit do poprawki | Znajdź commit, na którym opierają się Twoje obecne zmiany, w celu poprawienia/zmiany commita. To pozwala Ci uniknąć przeglądania commitów w Twojej gałęzi jeden po drugim, aby zobaczyć, który commit powinien być poprawiony/zmieniony. Zobacz dokumentację: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edytuj | Otwórz plik w zewnętrznym edytorze. |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` i `` | Ignoruj lub wyklucz plik |  |
| `` r `` | Odśwież pliki |  |
| `` s `` | Schowaj | Schowaj wszystkie zmiany. Dla innych wariantów schowania, użyj klawisza wyświetlania opcji schowka. |
| `` S `` | Wyświetl opcje schowka | Wyświetl opcje schowka (np. schowaj wszystko, schowaj zatwierdzone, schowaj niezatwierdzone). |
| `` a `` | Zatwierdź wszystko | Przełącz zatwierdzenie/odznaczenie dla wszystkich plików w drzewie roboczym. |
| `` <enter> `` | Zatwierdź linie / Zwiń katalog | Jeśli wybrany element jest plikiem, skup się na widoku zatwierdzania, aby móc zatwierdzać poszczególne fragmenty/linie. Jeśli wybrany element jest katalogiem, zwiń/rozwiń go. |
| `` d `` | Odrzuć | Wyświetl opcje odrzucania zmian w wybranym pliku. |
| `` g `` | Pokaż opcje resetowania do upstream |  |
| `` D `` | Reset | Wyświetl opcje resetu dla drzewa roboczego (np. zniszczenie drzewa roboczego). |
| `` ` `` | Przełącz widok drzewa plików | Przełącz widok plików między płaskim a drzewem. Płaski układ pokazuje wszystkie ścieżki plików na jednej liście, układ drzewa grupuje pliki według katalogów. |
| `` <c-t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` M `` | Otwórz zewnętrzne narzędzie scalania | Uruchom `git mergetool`. |
| `` f `` | Pobierz | Pobierz zmiany ze zdalnego serwera. |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Pliki commita

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopiuj ścieżkę do schowka |  |
| `` c `` | Przełącz | Przełącz plik. Zastępuje plik w twoim drzewie roboczym wersją z wybranego commita. |
| `` d `` | Usuń | Odrzuć zmiany w tym pliku z tego commita. Uruchamia interaktywny rebase w tle, więc możesz otrzymać konflikt scalania, jeśli późniejszy commit również zmienia ten plik. |
| `` o `` | Otwórz plik | Otwórz plik w domyślnej aplikacji. |
| `` e `` | Edytuj | Otwórz plik w zewnętrznym edytorze. |
| `` <c-t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` <space> `` | Przełącz plik włączony w łatkę | Przełącz, czy plik jest włączony w niestandardową łatkę. Zobacz https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Przełącz wszystkie pliki | Dodaj/usuń wszystkie pliki commita do niestandardowej łatki. Zobacz https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Wejdź do pliku / Przełącz zwiń katalog | Jeśli plik jest wybrany, wejdź do pliku, aby móc dodawać/usuwać poszczególne linie do niestandardowej łatki. Jeśli wybrany jest katalog, przełącz katalog. |
| `` ` `` | Przełącz widok drzewa plików | Przełącz widok plików między płaskim a drzewem. Płaski układ pokazuje wszystkie ścieżki plików na jednej liście, układ drzewa grupuje pliki według katalogów. |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Podsumowanie commita

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Potwierdź |  |
| `` <esc> `` | Zamknij |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopiuj hash commita do schowka |  |
| `` <space> `` | Przełącz | Przełącz wybrany commit jako odłączoną HEAD. |
| `` y `` | Kopiuj atrybut commita do schowka | Kopiuj atrybut commita do schowka (np. hash, URL, różnice, wiadomość, autor). |
| `` o `` | Otwórz commit w przeglądarce |  |
| `` n `` | Utwórz nową gałąź z commita |  |
| `` g `` | Reset | Wyświetl opcje resetu (miękki/mieszany/twardy) do wybranego elementu. |
| `` C `` | Kopiuj (cherry-pick) | Oznacz commit jako skopiowany. Następnie, w widoku lokalnych commitów, możesz nacisnąć `V`, aby wkleić (cherry-pick) skopiowane commity do sprawdzonej gałęzi. W dowolnym momencie możesz nacisnąć `<esc>`, aby anulować zaznaczenie. |
| `` <c-r> `` | Resetuj wybrane (cherry-picked) commity |  |
| `` <c-t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` <enter> `` | Pokaż commity |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Schowek

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Zastosuj | Zastosuj wpis schowka do katalogu roboczego. |
| `` g `` | Wyciągnij | Zastosuj wpis schowka do katalogu roboczego i usuń wpis schowka. |
| `` d `` | Usuń | Usuń wpis schowka z listy schowka. |
| `` n `` | Nowa gałąź | Utwórz nową gałąź z wybranego wpisu schowka. Działa poprzez przełączenie git na commit, na którym wpis schowka został utworzony, tworzenie nowej gałęzi z tego commita, a następnie zastosowanie wpisu schowka do nowej gałęzi jako dodatkowego commita. |
| `` r `` | Zmień nazwę schowka |  |
| `` <enter> `` | Wyświetl pliki |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Otwórz plik konfiguracyjny | Otwórz plik w domyślnej aplikacji. |
| `` e `` | Edytuj plik konfiguracyjny | Otwórz plik w zewnętrznym edytorze. |
| `` u `` | Sprawdź aktualizacje |  |
| `` <enter> `` | Przełącz na ostatnie repozytorium |  |
| `` a `` | Pokaż wszystkie gałęzie w logach |  |

## Sub-commity

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopiuj hash commita do schowka |  |
| `` <space> `` | Przełącz | Przełącz wybrany commit jako odłączoną HEAD. |
| `` y `` | Kopiuj atrybut commita do schowka | Kopiuj atrybut commita do schowka (np. hash, URL, różnice, wiadomość, autor). |
| `` o `` | Otwórz commit w przeglądarce |  |
| `` n `` | Utwórz nową gałąź z commita |  |
| `` g `` | Reset | Wyświetl opcje resetu (miękki/mieszany/twardy) do wybranego elementu. |
| `` C `` | Kopiuj (cherry-pick) | Oznacz commit jako skopiowany. Następnie, w widoku lokalnych commitów, możesz nacisnąć `V`, aby wkleić (cherry-pick) skopiowane commity do sprawdzonej gałęzi. W dowolnym momencie możesz nacisnąć `<esc>`, aby anulować zaznaczenie. |
| `` <c-r> `` | Resetuj wybrane (cherry-picked) commity |  |
| `` <c-t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` <enter> `` | Wyświetl pliki |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Szukaj w bieżącym widoku po tekście |  |

## Submoduły

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopiuj nazwę submodułu do schowka |  |
| `` <enter> `` | Wejdź | Wejdź do submodułu. Po wejściu do submodułu możesz nacisnąć `<esc>`, aby wrócić do repozytorium nadrzędnego. |
| `` d `` | Usuń | Usuń wybrany submoduł i odpowiadający mu katalog. |
| `` u `` | Aktualizuj | Aktualizuj wybrany submoduł. |
| `` n `` | Nowy submoduł |  |
| `` e `` | Zaktualizuj URL submodułu |  |
| `` i `` | Zainicjuj | Zainicjuj wybrany submoduł, aby przygotować do pobrania. Prawdopodobnie chcesz to kontynuować, wywołując akcję 'update', aby pobrać submoduł. |
| `` b `` | Pokaż opcje masowych operacji na submodułach |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Tagi

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Przełącz | Przełącz wybrany tag jako odłączoną głowę (detached HEAD). |
| `` n `` | Nowy tag | Utwórz nowy tag z bieżącego commita. Zostaniesz poproszony o wprowadzenie nazwy tagu i opcjonalnego opisu. |
| `` d `` | Usuń | Wyświetl opcje usuwania lokalnego/odległego tagu. |
| `` P `` | Wyślij tag | Wyślij wybrany tag do zdalnego. Zostaniesz poproszony o wybranie zdalnego. |
| `` g `` | Reset | Wyświetl opcje resetu (miękki/mieszany/twardy) do wybranego elementu. |
| `` <c-t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` <enter> `` | Pokaż commity |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Zdalne

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Wyświetl gałęzie |  |
| `` n `` | Nowy zdalny |  |
| `` d `` | Usuń | Usuń wybrany zdalny. Wszelkie lokalne gałęzie śledzące gałąź zdalną z tego zdalnego nie zostaną dotknięte. |
| `` e `` | Edytuj | Edytuj nazwę lub URL wybranego zdalnego. |
| `` f `` | Pobierz | Pobierz aktualizacje z zdalnego repozytorium. Pobiera nowe commity i gałęzie bez scalania ich z lokalnymi gałęziami. |
| `` / `` | Filtruj bieżący widok po tekście |  |

## Zdalne gałęzie

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Kopiuj nazwę gałęzi do schowka |  |
| `` <space> `` | Przełącz | Przełącz na nową lokalną gałąź na podstawie wybranej gałęzi zdalnej. Nowa gałąź będzie śledzić gałąź zdalną. |
| `` n `` | Nowa gałąź |  |
| `` M `` | Scal | Scal wybraną gałąź z aktualnie sprawdzoną gałęzią. |
| `` r `` | Przebazuj | Przebazuj przełączoną gałąź na wybraną gałąź. |
| `` d `` | Usuń | Usuń gałąź zdalną ze zdalnego. |
| `` u `` | Ustaw jako upstream | Ustaw wybraną gałąź zdalną jako upstream sprawdzonej gałęzi. |
| `` s `` | Kolejność sortowania |  |
| `` g `` | Reset | Wyświetl opcje resetu (miękki/mieszany/twardy) do wybranego elementu. |
| `` <c-t> `` | Otwórz zewnętrzne narzędzie różnic (git difftool) |  |
| `` <enter> `` | Pokaż commity |  |
| `` w `` | Zobacz opcje drzewa pracy |  |
| `` / `` | Filtruj bieżący widok po tekście |  |
