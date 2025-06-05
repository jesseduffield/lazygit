_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit 按键绑定

_图例：`<c-b>` 意味着ctrl+b, `<a-b>意味着Alt+b, `B` 意味着shift+b_

## 全局键绑定

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | 切换到最近的仓库 |  |
| `` <pgup> (fn+up/shift+k) `` | 向上滚动主面板 |  |
| `` <pgdown> (fn+down/shift+j) `` | 向下滚动主面板 |  |
| `` @ `` | 打开命令日志菜单 | 查看命令日志的选项，例如显示/隐藏命令日志以及聚焦命令日志 |
| `` P `` | 推送 | 推送当前分支到它的上游。如果上游未配置，你可以在弹窗中配置上游分支。 |
| `` p `` | 拉取 | 从当前分支的远程分支获取改动。如果上游未配置，你可以在弹窗中配置上游分支。 |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` } `` | 扩大差异视图中显示的上下文范围 | 增加diff视图中围绕更改显示的上下文数量 |
| `` { `` | 缩小差异视图中显示的上下文范围 | 减少diff视图中围绕更改显示的上下文数量 |
| `` : `` | 执行 Shell 命令 | Bring up a prompt where you can enter a shell command to execute. |
| `` <c-p> `` | 查看自定义补丁选项 |  |
| `` m `` | 查看 合并/变基 选项 | 查看当前合并或变基的中止、继续、跳过选项 |
| `` R `` | 刷新 | 刷新git状态(即在后台上运行`git status`,`git branch`等命令以更新面板内容) 不会运行`git fetch` |
| `` + `` | 下一屏模式(正常/半屏/全屏) |  |
| `` _ `` | 上一屏模式 |  |
| `` ? `` | 打开菜单 |  |
| `` <c-s> `` | 查看按路径过滤选项 | 查看用于过滤提交日志的选项，以便仅显示与过滤器匹配的提交。 |
| `` W `` | 打开 diff 菜单 | 查看与比较两个引用相关的选项，例如与选定的 ref 进行比较，输入要比较的 ref，然后反转比较方向。 |
| `` <c-e> `` | 打开 diff 菜单 | 查看与比较两个引用相关的选项，例如与选定的 ref 进行比较，输入要比较的 ref，然后反转比较方向。 |
| `` q `` | 退出 |  |
| `` <esc> `` | 取消 |  |
| `` <c-w> `` | 切换是否在差异视图中显示空白字符差异 | 切换是否在diff视图中显示空白更改 |
| `` z `` | 撤销 | Reflog将用于确定运行哪个git命令来撤消最后一个git命令。这并不包括对工作树的更改，只考虑提交。 |
| `` <c-z> `` | 重做 | Reflog将用于确定运行哪个git命令来重做上一个git命令。这并不包括对工作树的更改，只考虑提交。 |

## 列表面板导航

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | 上一页 |  |
| `` . `` | 下一页 |  |
| `` < (<home>) `` | 滚动到顶部 |  |
| `` > (<end>) `` | 滚动到底部 |  |
| `` v `` | 切换拖动选择 |  |
| `` <s-down> `` | 向下扩展选择范围 |  |
| `` <s-up> `` | 向上扩展选择范围 |  |
| `` / `` | 开始搜索 |  |
| `` H `` | 向左滚动 |  |
| `` L `` | 向右滚动 |  |
| `` ] `` | 下一个标签 |  |
| `` [ `` | 上一个标签 |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将提交的 hash 复制到剪贴板 |  |
| `` <space> `` | 检出 | 检出所选择的提交作为分离HEAD。 |
| `` y `` | 复制提交属性到剪贴板 | 复制提交属性到剪贴板(例如，hash、URL、diff、消息、作者)。 |
| `` o `` | 在浏览器中打开提交 |  |
| `` n `` | 从提交创建新分支 |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.

Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | 查看重置选项 | 查看重置选项 (soft/mixed/hard) 用于重置到选择项 |
| `` C `` | 复制提交(拣选) | 标记提交为已复制。然后，在本地提交视图中，你可以按 `V` (Cherry-Pick) 将已复制的提交粘贴到已检出的分支中。任何时候都可以按 `<esc>` 来取消选择。 |
| `` <c-r> `` | 重置已拣选(复制)的提交 |  |
| `` <c-t> `` | 使用外部差异比较工具(git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 查看提交 |  |
| `` w `` | 查看工作区选项 |  |
| `` / `` | 通过文本过滤当前视图 |  |

## 子提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将提交的 hash 复制到剪贴板 |  |
| `` <space> `` | 检出 | 检出所选择的提交作为分离HEAD。 |
| `` y `` | 复制提交属性到剪贴板 | 复制提交属性到剪贴板(例如，hash、URL、diff、消息、作者)。 |
| `` o `` | 在浏览器中打开提交 |  |
| `` n `` | 从提交创建新分支 |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.

Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | 查看重置选项 | 查看重置选项 (soft/mixed/hard) 用于重置到选择项 |
| `` C `` | 复制提交(拣选) | 标记提交为已复制。然后，在本地提交视图中，你可以按 `V` (Cherry-Pick) 将已复制的提交粘贴到已检出的分支中。任何时候都可以按 `<esc>` 来取消选择。 |
| `` <c-r> `` | 重置已拣选(复制)的提交 |  |
| `` <c-t> `` | 使用外部差异比较工具(git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 查看提交的文件 |  |
| `` w `` | 查看工作区选项 |  |
| `` / `` | 开始搜索 |  |

## 子模块

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将子模块名称复制到剪贴板 |  |
| `` <enter> `` | 进入 | 输入子模块 |
| `` d `` | 删除 | 删除选定的子模块及其相应的目录 |
| `` u `` | 更新 | 更新子模块 |
| `` n `` | 添加新的子模块 |  |
| `` e `` | 更新子模块 URL |  |
| `` i `` | 初始化 | 初始化子模块 |
| `` b `` | 查看批量子模块选项 |  |
| `` / `` | 通过文本过滤当前视图 |  |

## 工作区

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | 新建工作树 |  |
| `` <space> `` | 切换 | 切换到选中的工作树 |
| `` o `` | 在编辑器中编写 |  |
| `` d `` | 删除 | 删除选定的工作树。这将删除工作树的目录以及 .git 目录中有关工作树的元数据。 |
| `` / `` | 通过文本过滤当前视图 |  |

## 提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将提交的 hash 复制到剪贴板 |  |
| `` <c-r> `` | 重置已拣选(复制)的提交 |  |
| `` b `` | 查看二分查找选项 |  |
| `` s `` | 压缩(Squash) | 将已选提交压缩到该提交之下。这些选定的提交的消息会附加到该提交的消息之下。 |
| `` f `` | 修正 （fixup） | 将选定的提交合并到其下面的提交中。与压缩类似，但所选提交的消息将被丢弃。 |
| `` r `` | 改写提交 | 重写所选提交的消息。 |
| `` R `` | 使用编辑器重命名提交 |  |
| `` d `` | 删除提交 | 删除选中的提交。这将通过变基从分支中删除该提交，如果该提交修改的内容依赖于后续的提交，则需要解决合并冲突。 |
| `` e `` | 编辑(开始交互式变基) | 编辑提交 |
| `` i `` | 开始交互式变基 | 为分支上的提交启动交互式变基。这将包括从 HEAD 提交到第一个合并提交或主分支提交的所有提交。
如果您想从所选提交启动交互式变基，请按 `e`。 |
| `` p `` | 拣选(Pick) | 标记选中的提交为 picked（变基过程中）。这意味该提交将在后续的变基中保留。 |
| `` F `` | 为此提交创建修正 | 创建修正提交 |
| `` S `` | 应用该修复提交 | 压缩所选提交之上或当前分支的所有 “fixup!” 提交（自动压缩）。 |
| `` <c-j> `` | 下移提交 |  |
| `` <c-k> `` | 上移提交 |  |
| `` V `` | 粘贴提交(拣选) |  |
| `` B `` | 标记一个主提交用于变基 | 选择下一次变基的主提交。当您变基到一个分支时，只有高于主提交的提交才会被引入。这使用“git rebase --onto”命令。 |
| `` A `` | 修补(Amend) | 用已暂存的变更来修补提交 |
| `` a `` | 修补提交属性 | 设置或重置提交的作者，或添加其他作者。 |
| `` t `` | 撤销(Revert) | 为所选提交创建还原提交，这会反向应用所选提交的更改。 |
| `` T `` | 标签提交 | 创建一个新标签指向所选提交。你可以在弹窗中输入标签名称和描述(可选)。 |
| `` <c-l> `` | 打开日志菜单 | 查看提交日志的选项，例如更改排序顺序、隐藏 git graph、显示整个 git graph。 |
| `` <space> `` | 检出 | 检出所选择的提交作为分离HEAD。 |
| `` y `` | 复制提交属性到剪贴板 | 复制提交属性到剪贴板(例如，hash、URL、diff、消息、作者)。 |
| `` o `` | 在浏览器中打开提交 |  |
| `` n `` | 从提交创建新分支 |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.

Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` g `` | 查看重置选项 | 查看重置选项 (soft/mixed/hard) 用于重置到选择项 |
| `` C `` | 复制提交(拣选) | 标记提交为已复制。然后，在本地提交视图中，你可以按 `V` (Cherry-Pick) 将已复制的提交粘贴到已检出的分支中。任何时候都可以按 `<esc>` 来取消选择。 |
| `` <c-t> `` | 使用外部差异比较工具(git difftool) |  |
| `` * `` | Select commits of current branch |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 查看提交的文件 |  |
| `` w `` | 查看工作区选项 |  |
| `` / `` | 开始搜索 |  |

## 提交信息

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 确认 |  |
| `` <esc> `` | 关闭 |  |

## 提交文件

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将文件名复制到剪贴板 |  |
| `` y `` | 复制到剪贴板 |  |
| `` c `` | 检出 | 检出文件 |
| `` d `` | 删除 | 放弃对此文件的提交变更 |
| `` o `` | 打开文件 | 使用默认程序打开该文件 |
| `` e `` | 编辑 | 使用外部编辑器打开文件 |
| `` <c-t> `` | 使用外部差异比较工具(git difftool) |  |
| `` <space> `` | 补丁中包含的切换文件 | 切换文件是否包含在自定义补丁中。请参阅 https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches。 |
| `` a `` | 操作所有文件 | 添加或删除所有提交中的文件到自定义的补丁中。请参阅 https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches。 |
| `` <enter> `` | 输入文件以将所选行添加到补丁中(或切换目录折叠) | 如果已选择一个文件，则Enter进入该文件，以便您可以向自定义补丁添加/删除单独的行。如果选择了目录，则切换目录。 |
| `` ` `` | 切换文件树视图 | 在平铺部署与树布局之间切换文件视图。平铺布局在一个列表中展示所有文件路径，树布局则根据目录分组展示。 |
| `` - `` | 折叠全部文件 | 折叠文件树中的全部目录 |
| `` = `` | 展开全部文件 | 展开文件树中的全部目录 |
| `` 0 `` | Focus main view |  |
| `` / `` | 开始搜索 |  |

## 文件

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将文件名复制到剪贴板 |  |
| `` <space> `` | 切换暂存状态 | 为选定的文件切换暂存状态 |
| `` <c-b> `` | 通过状态过滤文件 |  |
| `` y `` | 复制到剪贴板 |  |
| `` c `` | 提交变更 | 提交暂存文件 |
| `` w `` | 提交变更而无需预先提交钩子 |  |
| `` A `` | 修补最后一次提交 |  |
| `` C `` | 使用 Git 编辑器提交变更 |  |
| `` <c-f> `` | 找到用于修复的基准提交 | 找到您当前变更所基于的提交，以便于修正/改进该提交。这样做可以省去您逐一查看分支提交来确定应该修正/改进哪个提交的麻烦。请参阅文档: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | 编辑 | 使用外部编辑器打开文件 |
| `` o `` | 打开文件 | 使用默认程序打开该文件 |
| `` i `` | 忽略文件 |  |
| `` r `` | 刷新文件 |  |
| `` s `` | 贮藏 | 贮藏所有变更.若要使用其他贮藏变体,请使用查看贮藏选项快捷键 |
| `` S `` | 查看贮藏选项 | 查看贮藏选项（例如：贮藏所有、贮藏已暂存变更、贮藏未暂存变更） |
| `` a `` | 切换所有文件的暂存状态 | 切换工作区中所有文件的已暂存/未暂存状态 |
| `` <enter> `` | 暂存单个 块/行 用于文件, 或 折叠/展开 目录 | 如果选中的是一个文件，则会进入到暂存视图，以便可以暂存单个代码块/行。如果选中的是一个目录，则会折叠/展开这个目录 |
| `` d `` | 查看'放弃变更'选项 | 查看选中文件的放弃变更选项 |
| `` g `` | 查看上游重置选项 |  |
| `` D `` | 重置 | 查看工作树的重置选项（例如：清除工作树）。 |
| `` ` `` | 切换文件树视图 | 在平铺部署与树布局之间切换文件视图。平铺布局在一个列表中展示所有文件路径，树布局则根据目录分组展示。 |
| `` <c-t> `` | 使用外部差异比较工具(git difftool) |  |
| `` M `` | 打开外部合并工具(git mergetool) | 执行 `git mergetool`. |
| `` f `` | 抓取 | 从远程获取变更 |
| `` - `` | 折叠全部文件 | 折叠文件树中的全部目录 |
| `` = `` | 展开全部文件 | 展开文件树中的全部目录 |
| `` 0 `` | Focus main view |  |
| `` / `` | 开始搜索 |  |

## 本地分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将分支名称复制到剪贴板 |  |
| `` i `` | 显示 git-flow 选项 |  |
| `` <space> `` | 检出 | 检出选中的项目 |
| `` n `` | 新分支 |  |
| `` N `` | Move commits to new branch | Create a new branch and move the unpushed commits of the current branch to it. Useful if you meant to start new work and forgot to create a new branch first.

Note that this disregards the selection, the new branch is always created either from the main branch or stacked on top of the current branch (you get to choose which). |
| `` o `` | 创建拉取请求 |  |
| `` O `` | 创建拉取请求选项 |  |
| `` <c-y> `` | 将拉取请求 URL 复制到剪贴板 |  |
| `` c `` | 按名称检出 | 按名称检出。在输入框中，您可以输入'-' 来切换到最后一个分支。 |
| `` F `` | 强制检出 | 强制检出所选分支。这将在检出所选分支之前放弃工作目录中的所有本地更改。 |
| `` d `` | 删除 | 查看本地/远程分支的删除选项 |
| `` r `` | 将已检出的分支变基到该分支 | 将检出的分支变基到所选的分支上。 |
| `` M `` | 合并到当前检出的分支 | Merge selected branch into currently checked out branch. |
| `` f `` | 从上游快进此分支 | 将当前分支直接移动到远程追踪分支的最新提交 |
| `` T `` | 创建标签 |  |
| `` s `` | 排序 |  |
| `` g `` | 查看重置选项 |  |
| `` R `` | 重命名分支 |  |
| `` u `` | 查看上游选项 | 查看与分支上游相关的选项，例如设置/取消设置上游和重置为上游。 |
| `` <c-t> `` | 使用外部差异比较工具(git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 查看提交 |  |
| `` w `` | 查看工作区选项 |  |
| `` / `` | 通过文本过滤当前视图 |  |

## 构建补丁中

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 选择上一个区块 |  |
| `` <right> `` | 选择下一个区块 |  |
| `` v `` | 切换拖动选择 |  |
| `` a `` | 切换选择代码块 | 切换代码块选择模式 |
| `` <c-o> `` | 将选中文本复制到剪贴板 |  |
| `` o `` | 打开文件 | 使用默认程序打开该文件 |
| `` e `` | 编辑文件 | 使用外部编辑器打开文件 |
| `` <space> `` | 添加/移除 行到补丁 |  |
| `` <esc> `` | 退出逐行模式 |  |
| `` / `` | 开始搜索 |  |

## 标签

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将标签复制到剪贴板 |  |
| `` <space> `` | 检出 | 检出选择的标签作为分离的HEAD |
| `` n `` | 创建标签 | 基于当前提交创建一个新标签。你将在弹窗中输入标签名称和描述(可选)。 |
| `` d `` | 删除 | 查看本地/远程标签的删除选项 |
| `` P `` | 推送标签 | 推送选择的标签到远端。你将在弹窗中选择一个远端。 |
| `` g `` | 重置 | 查看重置选项 (soft/mixed/hard) 用于重置到选择项 |
| `` <c-t> `` | 使用外部差异比较工具(git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 查看提交 |  |
| `` w `` | 查看工作区选项 |  |
| `` / `` | 通过文本过滤当前视图 |  |

## 次要

| Key | Action | Info |
|-----|--------|-------------|
| `` <tab> `` | 切换到其他面板 | 切换到其他视图（已暂存/未暂存的变更） |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | 开始搜索 |  |

## 正在合并

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 选中区块 |  |
| `` b `` | 选中所有区块 |  |
| `` <up> `` | 选择顶部块 |  |
| `` <down> `` | 选择底部块 |  |
| `` <left> `` | 选择上一个冲突 |  |
| `` <right> `` | 选择下一个冲突 |  |
| `` z `` | 撤销 | 撤消上次合并冲突解决 |
| `` e `` | 编辑文件 | 使用外部编辑器打开文件 |
| `` o `` | 打开文件 | 使用默认程序打开该文件 |
| `` M `` | 打开外部合并工具(git mergetool) | 执行 `git mergetool`. |
| `` <esc> `` | 返回文件面板 |  |

## 正在暂存

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 选择上一个区块 |  |
| `` <right> `` | 选择下一个区块 |  |
| `` v `` | 切换拖动选择 |  |
| `` a `` | 切换选择代码块 | 切换代码块选择模式 |
| `` <c-o> `` | 将选中文本复制到剪贴板 |  |
| `` <space> `` | 切换暂存状态 | 切换行暂存状态 |
| `` d `` | 取消变更(git reset) | 当选择未暂存的变更时，使用git reset丢弃该变更。当选择已暂存的变更时，取消暂存该变更 |
| `` o `` | 打开文件 | 使用默认程序打开该文件 |
| `` e `` | 编辑文件 | 使用外部编辑器打开文件 |
| `` <esc> `` | 返回文件面板 |  |
| `` <tab> `` | 切换到其他面板 | 切换到其他视图（已暂存/未暂存的变更） |
| `` E `` | 编辑代码块 | 在外部编辑器中编辑选中的代码块 |
| `` c `` | 提交变更 | 提交暂存文件 |
| `` w `` | 提交变更而无需预先提交钩子 |  |
| `` C `` | 使用 Git 编辑器提交变更 |  |
| `` <c-f> `` | 找到用于修复的基准提交 | 找到您当前变更所基于的提交，以便于修正/改进该提交。这样做可以省去您逐一查看分支提交来确定应该修正/改进哪个提交的麻烦。请参阅文档: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | 开始搜索 |  |

## 正常

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | 向下滚动 |  |
| `` mouse wheel up (fn+down) `` | 向上滚动 |  |
| `` <tab> `` | 切换到其他面板 | 切换到其他视图（已暂存/未暂存的变更） |
| `` <esc> `` | Exit back to side panel |  |
| `` / `` | 开始搜索 |  |

## 状态

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 打开配置文件 | 使用默认程序打开该文件 |
| `` e `` | 编辑配置文件 | 使用外部编辑器打开文件 |
| `` u `` | 检查更新 |  |
| `` <enter> `` | 切换到最近的仓库 |  |
| `` a `` | Show/cycle all branch logs |  |

## 确认面板

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 确认 |  |
| `` <esc> `` | 关闭 |  |

## 菜单

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 执行 |  |
| `` <esc> `` | 关闭 |  |
| `` / `` | 通过文本过滤当前视图 |  |

## 贮藏

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 应用 | 将贮藏项应用到您的工作目录。 |
| `` g `` | 应用并删除 | 将存储项应用到工作目录并删除存储项。 |
| `` d `` | 删除 | 从贮藏列表中删除该贮藏项 |
| `` n `` | 新分支 | 从选定的贮藏项创建一个新分支。这是通过 git 检查创建贮藏项的提交，从该提交创建一个新分支，然后将贮藏项作为附加提交应用到新分支来实现的。 |
| `` r `` | 重命名贮藏 |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 查看提交的文件 |  |
| `` w `` | 查看工作区选项 |  |
| `` / `` | 通过文本过滤当前视图 |  |

## 远程

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 查看分支 |  |
| `` n `` | 添加新的远程仓库 |  |
| `` d `` | 删除 | 删除选中的远程。从远程跟踪远程分支的任何本地分支都不会受到影响。 |
| `` e `` | 编辑 | 编辑远程仓库 |
| `` f `` | 抓取 | 抓取远程仓库 |
| `` / `` | 通过文本过滤当前视图 |  |

## 远程分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将分支名称复制到剪贴板 |  |
| `` <space> `` | 检出 | 基于当前选中的远程分支检出一个新的本地分支，或者将远程分支作分离的HEAD。 |
| `` n `` | 新分支 |  |
| `` M `` | 合并到当前检出的分支 | Merge selected branch into currently checked out branch. |
| `` r `` | 将已检出的分支变基到该分支 | 将检出的分支变基到所选的分支上。 |
| `` d `` | 删除 | 从远程删除远程分支。 |
| `` u `` | 设置为上游 | 设置为检出分支的上游 |
| `` s `` | 排序 |  |
| `` g `` | 查看重置选项 | 查看重置选项 (soft/mixed/hard) 用于重置到选择项 |
| `` <c-t> `` | 使用外部差异比较工具(git difftool) |  |
| `` 0 `` | Focus main view |  |
| `` <enter> `` | 查看提交 |  |
| `` w `` | 查看工作区选项 |  |
| `` / `` | 通过文本过滤当前视图 |  |
