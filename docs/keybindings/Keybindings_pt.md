_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit Keybindings

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## Combinações globais de teclas

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | Mudar para um repositório recente |  |
| `` <pgup> (fn+up/shift+k) `` | Rolar janela principal para cima |  |
| `` <pgdown> (fn+down/shift+j) `` | Rolar a janela principal para baixo |  |
| `` @ `` | View command log options | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | Empurre (Push) | Faça push do branch atual para o seu branch upstream. Se nenhum upstream estiver configurado, você será solicitado a configurar um branch a montante. |
| `` p `` | Puxar (Pull) | Puxe alterações do controle remoto para o ramo atual. Se nenhum upstream estiver configurado, será solicitado configurar um ramo a montante. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename.<br><br>The default can be changed in the config file with the key 'git.renameSimilarityThreshold'. |
| `` } `` | Increase diff context size | Increase the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` { `` | Decrease diff context size | Decrease the amount of the context shown around changes in the diff view.<br><br>The default can be changed in the config file with the key 'git.diffContextSize'. |
| `` : `` | Executar comando da shell | Traga um prompt onde você pode digitar um comando shell para executar. |
| `` <c-p> `` | Ver opções de patch personalizadas |  |
| `` m `` | Ver opções de mesclar/rebase | Ver opções para abortar/continuar/pular o merge/rebase atual. |
| `` R `` | Atualizar | Atualize o estado do git (ou seja, execute `git status`, `git branch`, etc em segundo plano para atualizar o conteúdo de painéis). Isso não executa `git fetch`. |
| `` + `` | Next screen mode (normal/half/fullscreen) |  |
| `` _ `` | Prev screen mode |  |
| `` \| `` | Cycle pagers | Choose the next pager in the list of configured pagers |
| `` <esc> `` | Cancelar |  |
| `` ? `` | Open keybindings menu |  |
| `` <c-s> `` | View filter options | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | View diffing options | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | View diffing options | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | Sair |  |
| `` <c-z> `` | Suspend the application |  |
| `` <c-w> `` | Toggle whitespace | Toggle whether or not whitespace changes are shown in the diff view.<br><br>The default can be changed in the config file with the key 'git.ignoreWhitespaceInDiffView'. |
| `` z `` | Desfazer | O reflog será usado para determinar qual comando git para executar para desfazer o último comando git. Isto não inclui mudanças na árvore de trabalho; apenas compromissos são tidos em consideração. |
| `` Z `` | Refazer | O reflog será usado para determinar qual comando git para executar para refazer o último comando git. Isto não inclui mudanças na árvore de trabalho; apenas compromissos são tidos em consideração. |

## List panel navigation

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | Previous page |  |
| `` . `` | Next page |  |
| `` < (<home>) `` | Scroll to top |  |
| `` > (<end>) `` | Scroll to bottom |  |
| `` v `` | Toggle range select |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | Search the current view by text |  |
| `` H `` | Rolar à esquerda |  |
| `` L `` | Scroll para a direita |  |
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
| `` w `` | Fazer commit de alterações sem pré-commit |  |
| `` A `` | Alterar último commit |  |
| `` C `` | Enviar alteração usando um editor Git |  |
| `` <c-f> `` | Encontrar commit da base para consertar | Encontre o commit em que as suas mudanças atuais estão se baseando, para alterar/consertar o commit. Isso poupa-te você de ter que olhar pelos commits da sua branch um por um para ver qual commit deve ser alterado/consertado<br>Veja a documentação:<br><https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Editar | Abrir arquivo no editor externo. |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` i `` | Ignore or exclude file |  |
| `` r `` | Atualizar arquivos |  |
| `` s `` | Stash | Stash todas as alterações. Para outras variações de armazenamento, use a fixação de teclas de armazenamento. |
| `` S `` | Ver opções de stash | Ver opções de stash (por exemplo, trash all, stash staged, stash unsttued). |
| `` a `` | Stage completo | Alternar para todos os arquivos na árvore de trabalho |
| `` <enter> `` | Stage lines / Colapso diretório | Se o item selecionado for um arquivo, o foco na exibição de preparo para o estágio de cenas/linhas individuais. Se o item selecionado for um diretório, recolher/expandi-lo. |
| `` d `` | Descartar | Exibir opções para descartar alterações para o arquivo selecionado. |
| `` g `` | View upstream reset options |  |
| `` D `` | Restaurar | Opções de redefinição de exibição para árvore de trabalho (por exemplo, nukando a árvore de trabalho). |
| `` ` `` | Alternar exibição de árvore de arquivo | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` M `` | View merge conflict options | View options for resolving merge conflicts. |
| `` f `` | Buscar | Buscar alterações do controle remoto. |
| `` - `` | Recolher todos os arquivos | Recolher todos os diretórios na árvore de arquivos |
| `` = `` | Expandir todos os arquivos | Expandir todos os diretórios na árvore do arquivo |
| `` 0 `` | Focus main view |  |
| `` / `` | Search the current view by text |  |

## Branches locais

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy branch name to clipboard |  |
| `` i `` | Exibir opções do git-flow |  |
| `` <space> `` | Verificar | Checar item selecionado |
| `` n `` | Nova branch |  |
| `` N `` | Mover commits para uma nova branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` o `` | Create pull request |  |
| `` O `` | View create pull request options |  |
| `` <c-y> `` | Copiar URL do pull request para área de transferência |  |
| `` c `` | Checar por nome | Checar por nome. Na caixa de entrada você pode inserir '-' para trocar para a última branch  |
| `` - `` | Checkout da branch anterior |  |
| `` F `` | Forçar checagem | Forçar checagem da branch selecionada. Isso irá descartar todas as mudanças no seu diretório de trabalho antes cheque a branch selecionada   |
| `` d `` | Apagar | Ver opções de exclusão para a branch local/remoto. |
| `` r `` | Refazer | Refazer a branch checada na branch selecionada |
| `` M `` | Mesclar | Ver opções para mesclar o item selecionado no branch atual (mesclar regularmente, mesclar squash) |
| `` f `` | Avanço rápido | Encaminhamento rápido de branch selecionada a partir do upstream. |
| `` T `` | New tag |  |
| `` s `` | Sort order |  |
| `` g `` | Restaurar |  |
| `` R `` | Rename branch |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` 0 `` | Focus main view |  |
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
| `` d `` | Apagar | Excluir o branch remoto do controle remoto. |
| `` u `` | Definir como upstream | Definir o ramo remoto selecionado como fluxo do branch check-out. |
| `` s `` | Sort order |  |
| `` g `` | Restaurar | Ver opções de redefinição (soft/mixed/hard) para redefinir para o item selecionado. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Commit arquivos

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy path to clipboard |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | Verificar | Arquivo de check-out. Isso substitui o arquivo em sua árvore de trabalho com a versão do commit selecionado. |
| `` d `` | Remover | Descartar as alterações desse commit para este arquivo. Isso executa uma rebase interativa em segundo plano, então você pode ter um conflito de merge se um commit posterior também alterar este arquivo. |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` e `` | Editar | Abrir arquivo no editor externo. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` <space> `` | Alternar entre o arquivo incluído no patch | Alternar se o arquivo está incluído no patch personalizado. Veja https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Alternar todos os arquivos | Adicionar/remover todos os arquivos de commit para atualização personalizada. Consulte https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Insira o arquivo / Alternar diretório recolhido | Se um arquivo estiver selecionado, insira o arquivo para que você possa adicionar/remover linhas individuais no patch personalizado. Se um diretório for selecionado, ative o diretório. |
| `` ` `` | Alternar exibição de árvore de arquivo | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory.<br><br>The default can be changed in the config file with the key 'gui.showFileTree'. |
| `` - `` | Recolher todos os arquivos | Recolher todos os diretórios na árvore de arquivos |
| `` = `` | Expandir todos os arquivos | Expandir todos os diretórios na árvore do arquivo |
| `` 0 `` | Focus main view |  |
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
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.<br>If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Escolher | Marque o commit selecionado para ser escolhido (quando meados da base). Isso significa que o commit será mantido ao continuar o rebase. |
| `` F `` | Criar commit de correção | Crie o commit 'correção!' para o commit selecionado. Mais tarde, você pode pressionar `S` neste mesmo commit para aplicar todas os commits de correção acima. |
| `` S `` | Aplicar commits de correções | Aplicar Squash all 'correção!', seja acima do commit selecionado, ou tudo no branch atual (autosquash). |
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
| `` N `` | Mover commits para uma nova branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Restaurar | Ver opções de redefinição (soft/mixed/hard) para redefinir para o item selecionado. |
| `` C `` | Copiar (cherry-pick) | Marcar commit como copiado. Então, dentro da visualização local de commits, você pode pressionar `V` para colar (cherry-pick) o(s) commit(s) copiado(s) em seu branch de check-out. A qualquer momento você pode pressionar `<esc>` para cancelar a seleção. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Ver arquivos |  |
| `` w `` | View worktree options |  |
| `` / `` | Search the current view by text |  |

## Confirmation panel

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Confirmar |  |
| `` <esc> `` | Fechar/Cancelar |  |
| `` <c-o> `` | Copy to clipboard |  |

## Etiquetas

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy tag to clipboard |  |
| `` <space> `` | Verificar | Checar a tag selecionada como um HEAD, desanexado |
| `` n `` | New tag | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Apagar | Ver opções de exclusão para tag local/remoto. |
| `` P `` | Push tag | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Restaurar | Ver opções de redefinição (soft/mixed/hard) para redefinir para o item selecionado. |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Input prompt

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Confirmar |  |
| `` <esc> `` | Fechar/Cancelar |  |

## Menu

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Executar |  |
| `` <esc> `` | Fechar/Cancelar |  |
| `` / `` | Filter the current view by text |  |

## Painel Principal (Normal)

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | Rolar para baixo |  |
| `` mouse wheel up (fn+down) `` | Rolar para cima |  |
| `` <tab> `` | Mudar de visão | Alternar para outra visão (staged/não processadas alterações). |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | Search the current view by text |  |

## Painel Principal (preparação)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Ir para o local anterior |  |
| `` <right> `` | Ir para o próximo trecho |  |
| `` v `` | Toggle range select |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` <space> `` | Etapa | Ativar/desativar seleção em staged/unstaged |
| `` d `` | Descartar | Quando a mudança não desejada for selecionada, descarte a mudança usando `git reset`. Quando a mudança em fase é selecionada, despare a mudança. |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` e `` | Editar arquivo | Abrir arquivo no editor externo. |
| `` <esc> `` | Retornar ao painel de arquivos |  |
| `` <tab> `` | Mudar de visão | Alternar para outra visão (staged/não processadas alterações). |
| `` E `` | Editar hunk | Editar o local selecionado no editor externo. |
| `` c `` | Commit | Submeter mudanças em staging |
| `` w `` | Fazer commit de alterações sem pré-commit |  |
| `` C `` | Enviar alteração usando um editor Git |  |
| `` <c-f> `` | Encontrar commit da base para consertar | Encontre o commit em que as suas mudanças atuais estão se baseando, para alterar/consertar o commit. Isso poupa-te você de ter que olhar pelos commits da sua branch um por um para ver qual commit deve ser alterado/consertado<br>Veja a documentação:<br><https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | Search the current view by text |  |

## Painel principal (mesclagem)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Escolha o local |  |
| `` b `` | Pegar todos os pedaços |  |
| `` <up> `` | Trecho anterior |  |
| `` <down> `` | Próximo trecho |  |
| `` <left> `` | Conflito anterior |  |
| `` <right> `` | Próximo conflito |  |
| `` z `` | Desfazer | Desfazer resolução de conflitos de última mesclagem. |
| `` e `` | Editar arquivo | Abrir arquivo no editor externo. |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` M `` | View merge conflict options | View options for resolving merge conflicts. |
| `` <esc> `` | Retornar ao painel de arquivos |  |

## Painel principal (patch build)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | Ir para o local anterior |  |
| `` <right> `` | Ir para o próximo trecho |  |
| `` v `` | Toggle range select |  |
| `` a `` | Toggle hunk selection | Toggle line-by-line vs. hunk selection mode. |
| `` <c-o> `` | Copy selected text to clipboard |  |
| `` o `` | Abrir arquivo | Abrir arquivo no aplicativo padrão. |
| `` e `` | Editar arquivo | Abrir arquivo no editor externo. |
| `` <space> `` | Alternar linhas no caminho |  |
| `` <esc> `` | Sair do construtor de patch personalizado |  |
| `` / `` | Search the current view by text |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit hash to clipboard |  |
| `` <space> `` | Verificar | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` N `` | Mover commits para uma nova branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Restaurar | Ver opções de redefinição (soft/mixed/hard) para redefinir para o item selecionado. |
| `` C `` | Copiar (cherry-pick) | Marcar commit como copiado. Então, dentro da visualização local de commits, você pode pressionar `V` para colar (cherry-pick) o(s) commit(s) copiado(s) em seu branch de check-out. A qualquer momento você pode pressionar `<esc>` para cancelar a seleção. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | View commits |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Remotes

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | Ver branches |  |
| `` n `` | Novo controle |  |
| `` d `` | Remover | Remover o controle remoto. Quaisquer ramificações locais de rastreamento de um ramo remoto do controle não serão afetadas. |
| `` e `` | Editar | Edit the selected remote's name or URL. |
| `` f `` | Buscar | Fetch updates from the remote repository. This retrieves new commits and branches without merging them into your local branches. |
| `` / `` | Filter the current view by text |  |

## Secundário

| Key | Action | Info |
|-----|--------|-------------|
| `` <tab> `` | Mudar de visão | Alternar para outra visão (staged/não processadas alterações). |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | Search the current view by text |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Aplicar | Aplique o stash no seu diretório de trabalho. |
| `` g `` | Pop | Aplique a entrada de stash no seu diretório de trabalho e remova a entrada de stash. |
| `` d `` | Descartar | Remova a entrada do stash da lista de armazenamento. |
| `` n `` | Nova branch | Criar um novo ramo a partir da entrada de lixo selecionada. Isso funciona verificando o commit do qual a entrada de lixo foi criada, criar um novo branch a partir desse commit e, em seguida, aplicar a entrada de lixo ao novo branch como um commit adicional. |
| `` r `` | Renomear o stasj |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Ver arquivos |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Status

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | Abrir o ficheiro de config | Abrir arquivo no aplicativo padrão. |
| `` e `` | Editar arquivo de configuração | Abrir arquivo no editor externo. |
| `` u `` | Verificar atualização |  |
| `` <enter> `` | Mudar para um repositório recente |  |
| `` a `` | Mostrar/ciclo todos os logs de filiais |  |
| `` 0 `` | Focus main view |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy commit hash to clipboard |  |
| `` <space> `` | Verificar | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | Open commit in browser |  |
| `` n `` | Create new branch off of commit |  |
| `` N `` | Mover commits para uma nova branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.<br><br>Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | Restaurar | Ver opções de redefinição (soft/mixed/hard) para redefinir para o item selecionado. |
| `` C `` | Copiar (cherry-pick) | Marcar commit como copiado. Então, dentro da visualização local de commits, você pode pressionar `V` para colar (cherry-pick) o(s) commit(s) copiado(s) em seu branch de check-out. A qualquer momento você pode pressionar `<esc>` para cancelar a seleção. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Abrir ferramenta de diff externa (git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | Ver arquivos |  |
| `` w `` | View worktree options |  |
| `` / `` | Search the current view by text |  |

## Submodules

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | Copy submodule name to clipboard |  |
| `` <enter> `` | Enter | Enter submodule. After entering the submodule, you can press `<esc>` to escape back to the parent repo. |
| `` d `` | Remover | Remove the selected submodule and its corresponding directory. |
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
| `` d `` | Remover | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | Filter the current view by text |  |
