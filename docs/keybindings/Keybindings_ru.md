_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Связки клавиш

_Связки клавиш_

## Глобальные сочетания клавиш

<pre>
  <kbd>&lt;c-r&gt;</kbd>: Переключиться на последний репозиторий
  <kbd>&lt;pgup&gt;</kbd>: Прокрутить вверх главную панель (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: Прокрутить вниз главную панель (fn+down/shift+j)
  <kbd>@</kbd>: Открыть меню журнала команд
  <kbd>}</kbd>: Увеличить размер контекста, отображаемого вокруг изменений в просмотрщике сравнении
  <kbd>{</kbd>: Уменьшите размер контекста, отображаемого вокруг изменений в просмотрщике сравнении
  <kbd>:</kbd>: Выполнить пользовательскую команду
  <kbd>&lt;c-p&gt;</kbd>: Просмотреть пользовательские параметры патча
  <kbd>m</kbd>: Просмотреть параметры слияния/перебазирования
  <kbd>R</kbd>: Обновить
  <kbd>+</kbd>: Следующий режим экрана (нормальный/полуэкранный/полноэкранный)
  <kbd>_</kbd>: Предыдущий режим экрана
  <kbd>?</kbd>: Открыть меню
  <kbd>&lt;c-s&gt;</kbd>: Просмотреть параметры фильтрации по пути
  <kbd>W</kbd>: Открыть меню сравнении
  <kbd>&lt;c-e&gt;</kbd>: Открыть меню сравнении
  <kbd>&lt;c-w&gt;</kbd>: Переключить отображение изменении пробелов в просмотрщике сравнении
  <kbd>z</kbd>: Отменить (через reflog) (экспериментальный)
  <kbd>&lt;c-z&gt;</kbd>: Повторить (через reflog) (экспериментальный)
  <kbd>P</kbd>: Отправить изменения
  <kbd>p</kbd>: Получить и слить изменения
</pre>

## Навигация по панели списка

<pre>
  <kbd>,</kbd>: Предыдущая страница
  <kbd>.</kbd>: Следующая страница
  <kbd>&lt;</kbd>: Пролистать наверх
  <kbd>&gt;</kbd>: Прокрутить вниз
  <kbd>v</kbd>: Переключить выборку перетаскивания
  <kbd>&lt;s-down&gt;</kbd>: Range select down
  <kbd>&lt;s-up&gt;</kbd>: Range select up
  <kbd>/</kbd>: Найти
  <kbd>H</kbd>: Прокрутить влево
  <kbd>L</kbd>: Прокрутить вправо
  <kbd>]</kbd>: Следующая вкладка
  <kbd>[</kbd>: Предыдущая вкладка
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

## Главная панель (Индексирование)

<pre>
  <kbd>&lt;left&gt;</kbd>: Выбрать предыдущую часть
  <kbd>&lt;right&gt;</kbd>: Выбрать следующую часть
  <kbd>v</kbd>: Переключить выборку перетаскивания
  <kbd>a</kbd>: Переключить выборку частей
  <kbd>&lt;c-o&gt;</kbd>: Скопировать выделенный текст в буфер обмена
  <kbd>o</kbd>: Открыть файл
  <kbd>e</kbd>: Редактировать файл
  <kbd>&lt;esc&gt;</kbd>: Вернуться к панели файлов
  <kbd>&lt;tab&gt;</kbd>: Переключиться на другую панель (проиндексированные/непроиндексированные изменения)
  <kbd>&lt;space&gt;</kbd>: Переключить строку в проиндексированные / непроиндексированные
  <kbd>d</kbd>: Отменить изменение (git reset)
  <kbd>E</kbd>: Изменить эту часть
  <kbd>c</kbd>: Сохранить изменения
  <kbd>w</kbd>: Закоммитить изменения без предварительного хука коммита
  <kbd>C</kbd>: Сохранить изменения с помощью редактора git
  <kbd>/</kbd>: Найти
</pre>

## Главная панель (Обычный)

<pre>
  <kbd>mouse wheel down</kbd>: Прокрутить вниз (fn+up)
  <kbd>mouse wheel up</kbd>: Прокрутить вверх (fn+down)
</pre>

## Главная панель (Слияние)

<pre>
  <kbd>e</kbd>: Редактировать файл
  <kbd>o</kbd>: Открыть файл
  <kbd>&lt;left&gt;</kbd>: Выбрать предыдущий конфликт
  <kbd>&lt;right&gt;</kbd>: Выбрать следующий конфликт
  <kbd>&lt;up&gt;</kbd>: Выбрать предыдущую часть
  <kbd>&lt;down&gt;</kbd>: Выбрать следующую часть
  <kbd>z</kbd>: Отменить
  <kbd>M</kbd>: Открыть внешний инструмент слияния (git mergetool)
  <kbd>&lt;space&gt;</kbd>: Выбрать эту часть
  <kbd>b</kbd>: Выбрать все части
  <kbd>&lt;esc&gt;</kbd>: Вернуться к панели файлов
</pre>

## Главная панель (сборка патчей)

<pre>
  <kbd>&lt;left&gt;</kbd>: Выбрать предыдущую часть
  <kbd>&lt;right&gt;</kbd>: Выбрать следующую часть
  <kbd>v</kbd>: Переключить выборку перетаскивания
  <kbd>a</kbd>: Переключить выборку частей
  <kbd>&lt;c-o&gt;</kbd>: Скопировать выделенный текст в буфер обмена
  <kbd>o</kbd>: Открыть файл
  <kbd>e</kbd>: Редактировать файл
  <kbd>&lt;space&gt;</kbd>: Добавить/удалить строку(и) для патча
  <kbd>&lt;esc&gt;</kbd>: Выйти из сборщика пользовательских патчей
  <kbd>/</kbd>: Найти
</pre>

## Журнал ссылок (Reflog)

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Скопировать SHA коммита в буфер обмена
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Переключить коммит
  <kbd>y</kbd>: Скопировать атрибут коммита
  <kbd>o</kbd>: Открыть коммит в браузере
  <kbd>n</kbd>: Создать новую ветку с этого коммита
  <kbd>g</kbd>: Просмотреть параметры сброса
  <kbd>C</kbd>: Скопировать отобранные коммит (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Сбросить отобранную (скопированную | cherry-picked) выборку коммитов
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>&lt;enter&gt;</kbd>: Просмотреть коммиты
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Коммиты

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Скопировать SHA коммита в буфер обмена
  <kbd>&lt;c-r&gt;</kbd>: Сбросить отобранную (скопированную | cherry-picked) выборку коммитов
  <kbd>b</kbd>: Просмотреть параметры бинарного поиска
  <kbd>s</kbd>: Объединить несколько коммитов в один нижний
  <kbd>f</kbd>: Объединить несколько коммитов в один отбросив сообщение коммита
  <kbd>r</kbd>: Перефразировать коммит
  <kbd>R</kbd>: Переписать коммит с помощью редактора
  <kbd>d</kbd>: Удалить коммит
  <kbd>e</kbd>: Изменить коммит
  <kbd>i</kbd>: Start interactive rebase
  <kbd>p</kbd>: Выбрать коммит (в середине перебазирования)
  <kbd>F</kbd>: Создать fixup коммит для этого коммита
  <kbd>S</kbd>: Объединить все 'fixup!' коммиты выше в выбранный коммит (автосохранение)
  <kbd>&lt;c-j&gt;</kbd>: Переместить коммит вниз на один
  <kbd>&lt;c-k&gt;</kbd>: Переместить коммит вверх на один
  <kbd>V</kbd>: Вставить отобранные коммиты (cherry-pick)
  <kbd>B</kbd>: Mark commit as base commit for rebase
  <kbd>A</kbd>: Править последний коммит с проиндексированными изменениями
  <kbd>a</kbd>: Установить/убрать автора коммита
  <kbd>t</kbd>: Отменить коммит
  <kbd>T</kbd>: Пометить коммит тегом
  <kbd>&lt;c-l&gt;</kbd>: Открыть меню журнала
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Переключить коммит
  <kbd>y</kbd>: Скопировать атрибут коммита
  <kbd>o</kbd>: Открыть коммит в браузере
  <kbd>n</kbd>: Создать новую ветку с этого коммита
  <kbd>g</kbd>: Просмотреть параметры сброса
  <kbd>C</kbd>: Скопировать отобранные коммит (cherry-pick)
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>&lt;enter&gt;</kbd>: Просмотреть файлы выбранного элемента
  <kbd>/</kbd>: Найти
</pre>

## Локальные Ветки

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Скопировать название ветки в буфер обмена
  <kbd>i</kbd>: Показать параметры git-flow
  <kbd>&lt;space&gt;</kbd>: Переключить
  <kbd>n</kbd>: Новая ветка
  <kbd>o</kbd>: Создать запрос на принятие изменений
  <kbd>O</kbd>: Создать параметры запроса принятие изменений
  <kbd>&lt;c-y&gt;</kbd>: Скопировать URL запроса на принятие изменений в буфер обмена
  <kbd>c</kbd>: Переключить по названию
  <kbd>F</kbd>: Принудительное переключение
  <kbd>d</kbd>: View delete options
  <kbd>r</kbd>: Перебазировать переключённую ветку на эту ветку
  <kbd>M</kbd>: Слияние с текущей переключённой веткой
  <kbd>f</kbd>: Перемотать эту ветку вперёд из её upstream-ветки
  <kbd>T</kbd>: Создать тег
  <kbd>s</kbd>: Порядок сортировки
  <kbd>g</kbd>: Просмотреть параметры сброса
  <kbd>R</kbd>: Переименовать ветку
  <kbd>u</kbd>: View upstream options
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: Просмотреть коммиты
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Меню

<pre>
  <kbd>&lt;enter&gt;</kbd>: Выполнить
  <kbd>&lt;esc&gt;</kbd>: Закрыть
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Панель Подтверждения

<pre>
  <kbd>&lt;enter&gt;</kbd>: Подтвердить
  <kbd>&lt;esc&gt;</kbd>: Закрыть/отменить
</pre>

## Подкоммиты

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Скопировать SHA коммита в буфер обмена
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: Переключить коммит
  <kbd>y</kbd>: Скопировать атрибут коммита
  <kbd>o</kbd>: Открыть коммит в браузере
  <kbd>n</kbd>: Создать новую ветку с этого коммита
  <kbd>g</kbd>: Просмотреть параметры сброса
  <kbd>C</kbd>: Скопировать отобранные коммит (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Сбросить отобранную (скопированную | cherry-picked) выборку коммитов
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>&lt;enter&gt;</kbd>: Просмотреть файлы выбранного элемента
  <kbd>/</kbd>: Найти
</pre>

## Подмодули

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Скопировать название подмодуля в буфер обмена
  <kbd>&lt;enter&gt;</kbd>: Ввести подмодуль
  <kbd>&lt;space&gt;</kbd>: Ввести подмодуль
  <kbd>d</kbd>: Удалить подмодуль
  <kbd>u</kbd>: Обновить подмодуль
  <kbd>n</kbd>: Добавить новый подмодуль
  <kbd>e</kbd>: Обновить URL подмодуля
  <kbd>i</kbd>: Инициализировать подмодуль
  <kbd>b</kbd>: Просмотреть параметры массового подмодуля
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Сводка коммита

<pre>
  <kbd>&lt;enter&gt;</kbd>: Подтвердить
  <kbd>&lt;esc&gt;</kbd>: Закрыть
</pre>

## Сохранить Изменения Файлов

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Скопировать закомиченное имя файла в буфер обмена
  <kbd>c</kbd>: Переключить файл
  <kbd>d</kbd>: Отменить изменения коммита в этом файле
  <kbd>o</kbd>: Открыть файл
  <kbd>e</kbd>: Редактировать файл
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>&lt;space&gt;</kbd>: Переключить файлы включённые в патч
  <kbd>a</kbd>: Переключить все файлы, включённые в патч
  <kbd>&lt;enter&gt;</kbd>: Введите файл, чтобы добавить выбранные строки в патч (или свернуть каталог переключения)
  <kbd>`</kbd>: Переключить вид дерева файлов
  <kbd>/</kbd>: Найти
</pre>

## Статус

<pre>
  <kbd>o</kbd>: Открыть файл конфигурации
  <kbd>e</kbd>: Редактировать файл конфигурации
  <kbd>u</kbd>: Проверить обновления
  <kbd>&lt;enter&gt;</kbd>: Переключиться на последний репозиторий
  <kbd>a</kbd>: Показать все логи ветки
</pre>

## Теги

<pre>
  <kbd>&lt;space&gt;</kbd>: Переключить
  <kbd>d</kbd>: View delete options
  <kbd>P</kbd>: Отправить тег
  <kbd>n</kbd>: Создать тег
  <kbd>g</kbd>: Просмотреть параметры сброса
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: Просмотреть коммиты
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Удалённые ветки

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Скопировать название ветки в буфер обмена
  <kbd>&lt;space&gt;</kbd>: Переключить
  <kbd>n</kbd>: Новая ветка
  <kbd>M</kbd>: Слияние с текущей переключённой веткой
  <kbd>r</kbd>: Перебазировать переключённую ветку на эту ветку
  <kbd>d</kbd>: Delete remote tag
  <kbd>u</kbd>: Установить как upstream-ветку переключённую ветку
  <kbd>s</kbd>: Порядок сортировки
  <kbd>g</kbd>: Просмотреть параметры сброса
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: Просмотреть коммиты
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Удалённые репозитории

<pre>
  <kbd>f</kbd>: Получение изменения из удалённого репозитория
  <kbd>n</kbd>: Добавить новую удалённую ветку
  <kbd>d</kbd>: Удалить удалённую ветку
  <kbd>e</kbd>: Редактировать удалённый репозитории
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Файлы

<pre>
  <kbd>&lt;c-o&gt;</kbd>: Скопировать название файла в буфер обмена
  <kbd>&lt;space&gt;</kbd>: Переключить индекс
  <kbd>&lt;c-b&gt;</kbd>: Фильтровать файлы (проиндексированные/непроиндексированные)
  <kbd>y</kbd>: Copy to clipboard
  <kbd>c</kbd>: Сохранить изменения
  <kbd>w</kbd>: Закоммитить изменения без предварительного хука коммита
  <kbd>A</kbd>: Правка последнего коммита
  <kbd>C</kbd>: Сохранить изменения с помощью редактора git
  <kbd>&lt;c-f&gt;</kbd>: Find base commit for fixup
  <kbd>e</kbd>: Редактировать файл
  <kbd>o</kbd>: Открыть файл
  <kbd>i</kbd>: Игнорировать или исключить файл
  <kbd>r</kbd>: Обновить файлы
  <kbd>s</kbd>: Припрятать все изменения
  <kbd>S</kbd>: Просмотреть параметры хранилища
  <kbd>a</kbd>: Все проиндексированные/непроиндексированные
  <kbd>&lt;enter&gt;</kbd>: Проиндексировать отдельные части/строки для файла или свернуть/развернуть для каталога
  <kbd>d</kbd>: Просмотреть параметры «отмены изменении»
  <kbd>g</kbd>: Просмотреть параметры сброса upstream-ветки
  <kbd>D</kbd>: Просмотреть параметры сброса
  <kbd>`</kbd>: Переключить вид дерева файлов
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>M</kbd>: Открыть внешний инструмент слияния (git mergetool)
  <kbd>f</kbd>: Получить изменения
  <kbd>/</kbd>: Найти
</pre>

## Хранилище

<pre>
  <kbd>&lt;space&gt;</kbd>: Применить припрятанные изменения
  <kbd>g</kbd>: Применить припрятанные изменения и тут же удалить их из хранилища
  <kbd>d</kbd>: Удалить припрятанные изменения из хранилища
  <kbd>n</kbd>: Новая ветка
  <kbd>r</kbd>: Переименовать хранилище
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: Просмотреть файлы выбранного элемента
  <kbd>/</kbd>: Filter the current view by text
</pre>
