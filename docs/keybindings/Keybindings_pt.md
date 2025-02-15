_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Combinações globais de teclas

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Mudar para um repositório recente |  |
| `` <pgup> (fn+up/shift+k) `` | Scroll up main window |  |
| `` <pgdown> (fn+down/shift+j) `` | Scroll down main window |  |
| `` @ `` | View command log options | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | Empurre (Push) | Faça push do branch atual para o seu branch upstream. Se nenhum upstream estiver configurado, você será solicitado a configurar um branch a montante. |
| `` p `` | Puxar (Pull) | Puxe alterações do controle remoto para o ramo atual. Se nenhum upstream estiver configurado, será solicitado configurar um ramo a montante. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` } `` | Increase diff context size | Increase the amount of the context shown around changes in the diff view. |
| `` { `` | Decrease diff context size | Decrease the amount of the context shown around changes in the diff view. |
| `` : `` | Execute shell command | Bring up a prompt where you can enter a shell command to execute. |
| `` <c-p> `` | View custom patch options |  |
| `` m `` | Ver opções de mesclar/rebase | Ver opções para abortar/continuar/pular o merge/rebase atual. |
| `` R `` | Atualizar | Atualize o estado do git (ou seja, execute `git status`, `git branch`, etc em segundo plano para atualizar o conteúdo de painéis). Isso não executa `git fetch`. |
| `` + `` | Next screen mode (normal/half/fullscreen) |  |
| `` _ `` | Prev screen mode |  |
| `` ? `` | Open keybindings menu |  |
| `` <c-s> `` | View filter options | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | View diffing options | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | View diffing options | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | Sair |  |
| `` <esc> `` | Cancel |  |
| `` <c-w> `` | Toggle whitespace | Toggle whether or not whitespace changes are shown in the diff view. |
| `` z `` | Desfazer | O reflog será usado para determinar qual comando git para executar para desfazer o último comando git. Isto não inclui mudanças na árvore de trabalho; apenas compromissos são tidos em consideração. |
| `` <c-z> `` | Refazer | O reflog será usado para determinar qual comando git para executar para refazer o último comando git. Isto não inclui mudanças na árvore de trabalho; apenas compromissos são tidos em consideração. |

## List panel navigation

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Previous page |  |
| `` . `` | Next page |  |
| `` < `` | Scroll to top |  |
| `` > `` | Scroll to bottom |  |
| `` v `` | Toggle range select |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | Search the current view by text |  |
| `` H `` | Scroll left |  |
| `` L `` | Scroll right |  |
| `` ] `` | Next tab |  |
| `` [ `` | Previous tab |  |

## Arquivos

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy path to clipboard |  |
| `` <space> `` | Etapa | Alternar para staging para o arquivo selecionado. |
| `` <c-b> `` | Filtrar arquivos por status |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Commit | Submeter mudanças em staging |
| `` w `` | Commit changes without pre-commit hook |  |
| `` A `` | Alterar último commit |  |
| `` C `` | Enviar alteração usando um editor Git |  |
| `` <c-f> `` | Encontrar commit da base para consertar | Encontre o commit em que as suas mudanças atuais estão se baseando, para alterar/consertar o commit. Isso poupa-te você de ter que olhar pelos commits da sua branch um por um para ver qual commit deve ser alterado/consertado
Veja a documentação:
<https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Editar | Abrir arquivo no editor externo. |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` i `` | Ignore or exclude file |  |
| `` r `` | Atualizar arquivos |  |
| `` s `` | Stash | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | View stash options | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | Stage completo | Alternar para todos os arquivos na árvore de trabalho |
| `` <enter> `` | Stage lines / Colapso diretório | Se o item selecionado for um arquivo, o foco na exibição de preparo para o estágio de cenas/linhas individuais. Se o item selecionado for um diretório, recolher/expandi-lo. |
| `` d `` | Discard | View options for discarding changes to the selected file. |
| `` g `` | View upstream reset options |  |
| `` D `` | Reset | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | Alternar exibição de árvore de arquivo | Alternar a visualização de arquivo entre layout plano e layout de árvore. Layout plano mostra todos os caminhos de arquivo em uma única lista, layout de árvore agrupa arquivos por diretório. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` M `` | Abrir ferramenta de merge externa | Execute `git mergetool`. |
| `` f `` | Buscar | Buscar alterações do controle remoto. |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Expand all files | Expand all directories in the file tree |
| `` / `` | Search the current view by text |  |

## Branches locais

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` i `` | Show git-flow options |  |
| `` <space> `` | Verificar | Checar item selecionado |
| `` n `` | Nova branch |  |
| `` o `` | Create pull request |  |
| `` O `` | View create pull request options |  |
| `` <c-y> `` | Copiar URL do pull request para área de transferência |  |
| `` c `` | Checar por nome | Checar por nome. Na caixa de entrada você pode inserir '-' para trocar para a última branch  |
| `` F `` | Forçar checagem | Forçar checagem da branch selecionada. Isso irá descartar todas as mudanças no seu diretório de trabalho antes cheque a branch selecionada   |
| `` d `` | Delete | View delete options for local/remote branch. |
| `` r `` | Refazer | Refazer a branch checada na branch selecionada |
| `` M `` | Mesclar | Ver opções para mesclar o item selecionado no branch atual (mesclar regularmente, mesclar squash) |
| `` f `` | Avanço rápido | Encaminhamento rápido de branch selecionada a partir do upstream. |
| `` T `` | New tag |  |
| `` s `` | Sort order |  |
| `` g `` | Reset |  |
| `` R `` | Rename branch |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Branches remotos

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` <space> `` | Verificar | Checar a nova branch baseada na brach remota selecionada, ou a branch remota como HEAD, desanexado |
| `` n `` | Nova branch |  |
| `` M `` | Mesclar | Ver opções para mesclar o item selecionado no branch atual (mesclar regularmente, mesclar squash) |
| `` r `` | Refazer | Refazer a branch checada na branch selecionada |
| `` d `` | Delete | Delete the remote branch from the remote. |
| `` u `` | Set as upstream | Set the selected remote branch as the upstream of the checked-out branch. |
| `` s `` | Sort order |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Commit files

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy path to clipboard |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Verificar | Checkout file. This replaces the file in your working tree with the version from the selected commit. |
| `` d `` | Remove | Discard this commit's changes to this file. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes this file. |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` e `` | Editar | Abrir arquivo no editor externo. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` <space> `` | Toggle file included in patch | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Toggle all files | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Enter file / Toggle directory collapsed | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | Alternar exibição de árvore de arquivo | Alternar a visualização de arquivo entre layout plano e layout de árvore. Layout plano mostra todos os caminhos de arquivo em uma única lista, layout de árvore agrupa arquivos por diretório. |
| `` - `` | Collapse all files | Collapse all directories in the files tree |
| `` = `` | Expand all files | Expand all directories in the file tree |
| `` / `` | Search the current view by text |  |

## Commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit hash to clipboard |  |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` b `` | View bisect options |  |
| `` s `` | Squash | Squash o commit selecionado no commit abaixo dele. A mensagem do commit selecionado será anexada ao commit abaixo dele. |
| `` f `` | Fixup | Faça o commit selecionado no commit abaixo dele. Semelhante para o squash, mas a mensagem do commit selecionado será descartada. |
| `` r `` | Reword | Repetir a mensagem de submissão selecionada. |
| `` R `` | Republicar com o editor |  |
| `` d `` | Descartar | Solte o commit selecionado. Isso irá remover o commit do branch através de uma rebase. Se o commit faz com que as alterações em commits posteriores dependem, você pode precisar resolver conflitos de merge. |
| `` e `` | Editar (iniciar rebase interativa) | Editar o commit selecionado. Use isto para iniciar uma rebase interativa a partir do commit selecionado. Quando já estiver no meio da reconstrução, isto irá marcar o commit selecionado para edição, o que significa que ao continuar com a reformulação. a rebase irá pausar no commit selecionado para permitir que você faça alterações. |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Escolher | Marque o commit selecionado para ser escolhido (quando meados da base). Isso significa que o commit será mantido ao continuar o rebase. |
| `` F `` | Create fixup commit | Create 'fixup!' commit for the selected commit. Later on, you can press `S` on this same commit to apply all above fixup commits. |
| `` S `` | Apply fixup commits | Squash all 'fixup!' commits, either above the selected commit, or all in current branch (autosquash). |
| `` <c-j> `` | Mover commit um para baixo |  |
| `` <c-k> `` | Mover o commit um para cima |  |
| `` V `` | Colar (cherry-pick) |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Modificar | Alterar o commit com mudanças em sted. Se o commit selecionado for o commit HEAD, ele executará o `git commit --amend`. Caso contrário, o compromisso será alterado por meio de uma base de apoio. |
| `` a `` | Alterar atributo de commit | Definir/Redefinir autor de submissão ou co-autor definido. |
| `` t `` | Reverter | Crie um commit reverter para o commit selecionado, que aplica as alterações do commit selecionado em reverso. |
| `` T `` | Tag commit | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <c-l> `` | View log options | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` <space> `` | Verificar | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Copiar (cherry-pick) | Marcar commit como copiado. Então, dentro da visualização local de commits, você pode pressionar `V` para colar (cherry-pick) o(s) commit(s) copiado(s) em seu branch de check-out. A qualquer momento você pode pressionar `<esc>` para cancelar a seleção. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` <enter> `` | View files |  |
| `` w `` | View worktree options |  |
| `` / `` | Search the current view by text |  |

## Confirmation panel

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Confirmar |  |
| `` <esc> `` | Fechar/Cancelar |  |

## Etiquetas

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy tag to clipboard |  |
| `` <space> `` | Verificar | Checar a tag selecionada como um HEAD, desanexado |
| `` n `` | New tag | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Delete | View delete options for local/remote tag. |
| `` P `` | Push tag | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Menu

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Executar |  |
| `` <esc> `` | Fechar |  |
| `` / `` | Filter the current view by text |  |

## Painel Principal (Normal)

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | Scroll down |  |
| `` mouse wheel up (fn+down) `` | Scroll up |  |

## Painel Principal (preparação)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Go to previous hunk |  |
| `` <right> `` | Go to next hunk |  |
| `` v `` | Toggle range select |  |
| `` a `` | Selecione o local | Ativa/desativa modo seleção de hunk  |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` <space> `` | Etapa | Ativar/desativar seleção em staged/unstaged |
| `` d `` | Descartar | Quando a mudança não desejada for selecionada, descarte a mudança usando `git reset`. Quando a mudança em fase é selecionada, despare a mudança. |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` e `` | Editar arquivo | Abrir arquivo no editor externo. |
| `` <esc> `` | Retornar ao painel de arquivos |  |
| `` <tab> `` | Mudar de visão | Alternar para outra visão (staged/não processadas alterações). |
| `` E `` | Editar hunk | Editar o local selecionado no editor externo. |
| `` c `` | Commit | Submeter mudanças em staging |
| `` w `` | Commit changes without pre-commit hook |  |
| `` C `` | Enviar alteração usando um editor Git |  |
| `` <c-f> `` | Encontrar commit da base para consertar | Encontre o commit em que as suas mudanças atuais estão se baseando, para alterar/consertar o commit. Isso poupa-te você de ter que olhar pelos commits da sua branch um por um para ver qual commit deve ser alterado/consertado
Veja a documentação:
<https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | Search the current view by text |  |

## Painel principal (mesclagem)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Escolha o local |  |
| `` b `` | Pegar todos os pedaços |  |
| `` <up> `` | Previous hunk |  |
| `` <down> `` | Next hunk |  |
| `` <left> `` | Previous conflict |  |
| `` <right> `` | Next conflict |  |
| `` z `` | Desfazer | Desfazer resolução de conflitos de última mesclagem. |
| `` e `` | Editar arquivo | Abrir arquivo no editor externo. |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` M `` | Abrir ferramenta de merge externa | Execute `git mergetool`. |
| `` <esc> `` | Retornar ao painel de arquivos |  |

## Painel principal (patch build)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Go to previous hunk |  |
| `` <right> `` | Go to next hunk |  |
| `` v `` | Toggle range select |  |
| `` a `` | Selecione o local | Ativa/desativa modo seleção de hunk  |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` e `` | Editar arquivo | Abrir arquivo no editor externo. |
| `` <space> `` | Alternar linhas no caminho |  |
| `` <esc> `` | Exit custom patch builder |  |
| `` / `` | Search the current view by text |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit hash to clipboard |  |
| `` <space> `` | Verificar | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Copiar (cherry-pick) | Marcar commit como copiado. Então, dentro da visualização local de commits, você pode pressionar `V` para colar (cherry-pick) o(s) commit(s) copiado(s) em seu branch de check-out. A qualquer momento você pode pressionar `<esc>` para cancelar a seleção. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Remotes

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | New remote |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Editar | Edit the selected remote's name or URL. |
| `` f `` | Buscar | Fetch updates from the remote repository. This retrieves new commits and branches without merging them into your local branches. |
| `` / `` | Filter the current view by text |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Aplicar | Aplique o stash no seu diretório de trabalho. |
| `` g `` | Pop | Aplique a entrada de stash no seu diretório de trabalho e remova a entrada de stash. |
| `` d `` | Descartar | Remova a entrada do stash da lista de armazenamento. |
| `` n `` | Nova branch | Criar um novo ramo a partir da entrada de lixo selecionada. Isso funciona verificando o commit do qual a entrada de lixo foi criada, criar um novo branch a partir desse commit e, em seguida, aplicar a entrada de lixo ao novo branch como um commit adicional. |
| `` r `` | Renomear o stasj |  |
| `` <enter> `` | View files |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Abrir o ficheiro de config | Abrir arquivo no aplicativo padrão. |
| `` e `` | Editar arquivo de configuração | Abrir arquivo no editor externo. |
| `` u `` | Verificar atualização |  |
| `` <enter> `` | Mudar para um repositório recente |  |
| `` a `` | Mostrar todos os logs da branch |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit hash to clipboard |  |
| `` <space> `` | Verificar | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | Copiar (cherry-pick) | Marcar commit como copiado. Então, dentro da visualização local de commits, você pode pressionar `V` para colar (cherry-pick) o(s) commit(s) copiado(s) em seu branch de check-out. A qualquer momento você pode pressionar `<esc>` para cancelar a seleção. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` <enter> `` | View files |  |
| `` w `` | View worktree options |  |
| `` / `` | Search the current view by text |  |

## Submodules

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy submodule name to clipboard |  |
| `` <enter> `` | Enter | Enter submodule. After entering the submodule, you can press `<esc>` to escape back to the parent repo. |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | Update selected submodule. |
| `` n `` | New submodule |  |
| `` e `` | Update submodule URL |  |
| `` i `` | Initialize | Initialize the selected submodule to prepare for fetching. You probably want to follow this up by invoking the 'update' action to fetch the submodule. |
| `` b `` | View bulk submodule options |  |
| `` / `` | Filter the current view by text |  |

## Sumário do commit

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Confirmar |  |
| `` <esc> `` | Fechar |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | Abrir no editor |  |
| `` d `` | Remove | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | Filter the current view by text |  |
