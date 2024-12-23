_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Связки клавиш

_Связки клавиш_

## Глобальные сочетания клавиш

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Переключиться на последний репозиторий |  |
| `` <pgup> (fn+up/shift+k) `` | Прокрутить вверх главную панель |  |
| `` <pgdown> (fn+down/shift+j) `` | Прокрутить вниз главную панель |  |
| `` @ `` | Открыть меню журнала команд | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | Отправить изменения | Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` p `` | Получить и слить изменения | Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` } `` | Увеличить размер контекста, отображаемого вокруг изменений в просмотрщике сравнении | Increase the amount of the context shown around changes in the diff view. |
| `` { `` | Уменьшите размер контекста, отображаемого вокруг изменений в просмотрщике сравнении | Decrease the amount of the context shown around changes in the diff view. |
| `` : `` | Execute shell command | Bring up a prompt where you can enter a shell command to execute. |
| `` <c-p> `` | Просмотреть пользовательские параметры патча |  |
| `` m `` | Просмотреть параметры слияния/перебазирования | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | Обновить | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | Следующий режим экрана (нормальный/полуэкранный/полноэкранный) |  |
| `` _ `` | Предыдущий режим экрана |  |
| `` ? `` | Открыть меню |  |
| `` <c-s> `` | Просмотреть параметры фильтрации по пути | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | Открыть меню сравнении | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | Открыть меню сравнении | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | Выйти |  |
| `` <esc> `` | Отменить |  |
| `` <c-w> `` | Переключить отображение изменении пробелов в просмотрщике сравнении | Toggle whether or not whitespace changes are shown in the diff view. |
| `` z `` | Отменить (через reflog) (экспериментальный) | Журнал ссылок (reflog) будет использоваться для определения того, какую команду git запустить, чтобы отменить последнюю команду git. Сюда не входят изменения в рабочем дереве; учитываются только коммиты. |
| `` <c-z> `` | Повторить (через reflog) (экспериментальный) | Журнал ссылок (reflog) будет использоваться для определения того, какую команду git нужно запустить, чтобы повторить последнюю команду git. Сюда не входят изменения в рабочем дереве; учитываются только коммиты. |

## Навигация по панели списка

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Предыдущая страница |  |
| `` . `` | Следующая страница |  |
| `` < `` | Пролистать наверх |  |
| `` > `` | Прокрутить вниз |  |
| `` v `` | Переключить выборку перетаскивания |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | Найти |  |
| `` H `` | Прокрутить влево |  |
| `` L `` | Прокрутить вправо |  |
| `` ] `` | Следующая вкладка |  |
| `` [ `` | Предыдущая вкладка |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | Open in editor |  |
| `` d `` | Remove | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | Filter the current view by text |  |

## Главная панель (Индексирование)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Выбрать предыдущую часть |  |
| `` <right> `` | Выбрать следующую часть |  |
| `` v `` | Переключить выборку перетаскивания |  |
| `` a `` | Переключить выборку частей | Toggle hunk selection mode. |
| `` <c-o> `` | Скопировать выделенный текст в буфер обмена |  |
| `` <space> `` | Переключить индекс | Переключить строку в проиндексированные / непроиндексированные |
| `` d `` | Отменить изменение (git reset) | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | Открыть файл | Open file in default application. |
| `` e `` | Редактировать файл | Open file in external editor. |
| `` <esc> `` | Вернуться к панели файлов |  |
| `` <tab> `` | Переключиться на другую панель (проиндексированные/непроиндексированные изменения) | Switch to other view (staged/unstaged changes). |
| `` E `` | Изменить эту часть | Edit selected hunk in external editor. |
| `` c `` | Сохранить изменения | Commit staged changes. |
| `` w `` | Закоммитить изменения без предварительного хука коммита |  |
| `` C `` | Сохранить изменения с помощью редактора git |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | Найти |  |

## Главная панель (Обычный)

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | Прокрутить вниз |  |
| `` mouse wheel up (fn+down) `` | Прокрутить вверх |  |

## Главная панель (Слияние)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Выбрать эту часть |  |
| `` b `` | Выбрать все части |  |
| `` <up> `` | Выбрать предыдущую часть |  |
| `` <down> `` | Выбрать следующую часть |  |
| `` <left> `` | Выбрать предыдущий конфликт |  |
| `` <right> `` | Выбрать следующий конфликт |  |
| `` z `` | Отменить | Undo last merge conflict resolution. |
| `` e `` | Редактировать файл | Open file in external editor. |
| `` o `` | Открыть файл | Open file in default application. |
| `` M `` | Открыть внешний инструмент слияния (git mergetool) | Run `git mergetool`. |
| `` <esc> `` | Вернуться к панели файлов |  |

## Главная панель (сборка патчей)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Выбрать предыдущую часть |  |
| `` <right> `` | Выбрать следующую часть |  |
| `` v `` | Переключить выборку перетаскивания |  |
| `` a `` | Переключить выборку частей | Toggle hunk selection mode. |
| `` <c-o> `` | Скопировать выделенный текст в буфер обмена |  |
| `` o `` | Открыть файл | Open file in default application. |
| `` e `` | Редактировать файл | Open file in external editor. |
| `` <space> `` | Добавить/удалить строку(и) для патча |  |
| `` <esc> `` | Выйти из сборщика пользовательских патчей |  |
| `` / `` | Найти |  |

## Журнал ссылок (Reflog)

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Скопировать hash коммита в буфер обмена |  |
| `` <space> `` | Переключить | Checkout the selected commit as a detached HEAD. |
| `` y `` | Скопировать атрибут коммита | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Открыть коммит в браузере |  |
| `` n `` | Создать новую ветку с этого коммита |  |
| `` g `` | Просмотреть параметры сброса | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Скопировать отобранные коммит (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Сбросить отобранную (скопированную | cherry-picked) выборку коммитов |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Просмотреть коммиты |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Коммиты

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Скопировать hash коммита в буфер обмена |  |
| `` <c-r> `` | Сбросить отобранную (скопированную | cherry-picked) выборку коммитов |  |
| `` b `` | Просмотреть параметры бинарного поиска |  |
| `` s `` | Объединить коммиты (Squash) | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | Объединить несколько коммитов в один отбросив сообщение коммита (Fixup)  | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` r `` | Перефразировать коммит | Reword the selected commit's message. |
| `` R `` | Переписать коммит с помощью редактора |  |
| `` d `` | Удалить коммит | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | Edit (start interactive rebase) | Изменить коммит |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick | Выбрать коммит (в середине перебазирования) |
| `` F `` | Создать fixup коммит | Создать fixup коммит для этого коммита |
| `` S `` | Apply fixup commits | Объединить все 'fixup!' коммиты выше в выбранный коммит (автосохранение) |
| `` <c-j> `` | Переместить коммит вниз на один |  |
| `` <c-k> `` | Переместить коммит вверх на один |  |
| `` V `` | Вставить отобранные коммиты (cherry-pick) |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Amend | Править последний коммит с проиндексированными изменениями |
| `` a `` | Установить/убрать автора коммита | Set/Reset commit author or set co-author. |
| `` t `` | Revert | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | Пометить коммит тегом | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <c-l> `` | Открыть меню журнала | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` <space> `` | Переключить | Checkout the selected commit as a detached HEAD. |
| `` y `` | Скопировать атрибут коммита | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Открыть коммит в браузере |  |
| `` n `` | Создать новую ветку с этого коммита |  |
| `` g `` | Просмотреть параметры сброса | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Скопировать отобранные коммит (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Просмотреть файлы выбранного элемента |  |
| `` w `` | View worktree options |  |
| `` / `` | Найти |  |

## Локальные Ветки

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Скопировать название ветки в буфер обмена |  |
| `` i `` | Показать параметры git-flow |  |
| `` <space> `` | Переключить | Checkout selected item. |
| `` n `` | Новая ветка |  |
| `` o `` | Создать запрос на принятие изменений |  |
| `` O `` | Создать параметры запроса принятие изменений |  |
| `` <c-y> `` | Скопировать URL запроса на принятие изменений в буфер обмена |  |
| `` c `` | Переключить по названию | Checkout by name. In the input box you can enter '-' to switch to the last branch. |
| `` F `` | Принудительное переключение | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | Delete | View delete options for local/remote branch. |
| `` r `` | Перебазировать переключённую ветку на эту ветку | Rebase the checked-out branch onto the selected branch. |
| `` M `` | Слияние с текущей переключённой веткой | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` f `` | Перемотать эту ветку вперёд из её upstream-ветки | Fast-forward selected branch from its upstream. |
| `` T `` | Создать тег |  |
| `` s `` | Порядок сортировки |  |
| `` g `` | Просмотреть параметры сброса |  |
| `` R `` | Переименовать ветку |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Просмотреть коммиты |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Меню

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Выполнить |  |
| `` <esc> `` | Закрыть |  |
| `` / `` | Filter the current view by text |  |

## Панель Подтверждения

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Подтвердить |  |
| `` <esc> `` | Закрыть/отменить |  |

## Подкоммиты

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Скопировать hash коммита в буфер обмена |  |
| `` <space> `` | Переключить | Checkout the selected commit as a detached HEAD. |
| `` y `` | Скопировать атрибут коммита | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Открыть коммит в браузере |  |
| `` n `` | Создать новую ветку с этого коммита |  |
| `` g `` | Просмотреть параметры сброса | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Скопировать отобранные коммит (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Сбросить отобранную (скопированную | cherry-picked) выборку коммитов |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Просмотреть файлы выбранного элемента |  |
| `` w `` | View worktree options |  |
| `` / `` | Найти |  |

## Подмодули

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Скопировать название подмодуля в буфер обмена |  |
| `` <enter> `` | Enter | Ввести подмодуль |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | Обновить подмодуль |
| `` n `` | Добавить новый подмодуль |  |
| `` e `` | Обновить URL подмодуля |  |
| `` i `` | Initialize | Инициализировать подмодуль |
| `` b `` | Просмотреть параметры массового подмодуля |  |
| `` / `` | Filter the current view by text |  |

## Сводка коммита

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Подтвердить |  |
| `` <esc> `` | Закрыть |  |

## Сохранить Изменения Файлов

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Скопировать название файла в буфер обмена |  |
| `` c `` | Переключить | Переключить файл |
| `` d `` | Remove | Отменить изменения коммита в этом файле |
| `` o `` | Открыть файл | Open file in default application. |
| `` e `` | Edit | Open file in external editor. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Переключить файлы включённые в патч | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Переключить все файлы, включённые в патч | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Введите файл, чтобы добавить выбранные строки в патч (или свернуть каталог переключения) | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | Переключить вид дерева файлов | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` / `` | Найти |  |

## Статус

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Открыть файл конфигурации | Open file in default application. |
| `` e `` | Редактировать файл конфигурации | Open file in external editor. |
| `` u `` | Проверить обновления |  |
| `` <enter> `` | Переключиться на последний репозиторий |  |
| `` a `` | Показать все логи ветки |  |

## Теги

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Переключить | Checkout the selected tag as a detached HEAD. |
| `` n `` | Создать тег | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Delete | View delete options for local/remote tag. |
| `` P `` | Отправить тег | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Просмотреть коммиты |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Удалённые ветки

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Скопировать название ветки в буфер обмена |  |
| `` <space> `` | Переключить | Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head. |
| `` n `` | Новая ветка |  |
| `` M `` | Слияние с текущей переключённой веткой | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` r `` | Перебазировать переключённую ветку на эту ветку | Rebase the checked-out branch onto the selected branch. |
| `` d `` | Delete | Delete the remote branch from the remote. |
| `` u `` | Set as upstream | Установить как upstream-ветку переключённую ветку |
| `` s `` | Порядок сортировки |  |
| `` g `` | Просмотреть параметры сброса | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | Просмотреть коммиты |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Удалённые репозитории

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | Добавить новую удалённую ветку |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Edit | Редактировать удалённый репозитории |
| `` f `` | Получить изменения | Получение изменения из удалённого репозитория |
| `` / `` | Filter the current view by text |  |

## Файлы

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Скопировать название файла в буфер обмена |  |
| `` <space> `` | Переключить индекс | Toggle staged for selected file. |
| `` <c-b> `` | Фильтровать файлы (проиндексированные/непроиндексированные) |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Сохранить изменения | Commit staged changes. |
| `` w `` | Закоммитить изменения без предварительного хука коммита |  |
| `` A `` | Правка последнего коммита |  |
| `` C `` | Сохранить изменения с помощью редактора git |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit | Open file in external editor. |
| `` o `` | Открыть файл | Open file in default application. |
| `` i `` | Игнорировать или исключить файл |  |
| `` r `` | Обновить файлы |  |
| `` s `` | Stash | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | Просмотреть параметры хранилища | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | Все проиндексированные/непроиндексированные | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | Проиндексировать отдельные части/строки для файла или свернуть/развернуть для каталога | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | Просмотреть параметры «отмены изменении» | View options for discarding changes to the selected file. |
| `` g `` | Просмотреть параметры сброса upstream-ветки |  |
| `` D `` | Reset | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | Переключить вид дерева файлов | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Открыть внешний инструмент слияния (git mergetool) | Run `git mergetool`. |
| `` f `` | Получить изменения | Fetch changes from remote. |
| `` / `` | Найти |  |

## Хранилище

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Применить припрятанные изменения | Apply the stash entry to your working directory. |
| `` g `` | Применить припрятанные изменения и тут же удалить их из хранилища | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | Удалить припрятанные изменения из хранилища | Remove the stash entry from the stash list. |
| `` n `` | Новая ветка | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | Переименовать хранилище |  |
| `` <enter> `` | Просмотреть файлы выбранного элемента |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |
