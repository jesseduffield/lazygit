_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit 按键绑定

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## 全局键绑定

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | 切换到最近的仓库 |  |
| `` <pgup> (fn+up/shift+k) `` | 向上滚动主面板 |  |
| `` <pgdown> (fn+down/shift+j) `` | 向下滚动主面板 |  |
| `` @ `` | 打开命令日志菜单 |  |
| `` } `` | 扩大差异视图中显示的上下文范围 |  |
| `` { `` | 缩小差异视图中显示的上下文范围 |  |
| `` : `` | 执行自定义命令 |  |
| `` <c-p> `` | 查看自定义补丁选项 |  |
| `` m `` | 查看 合并/变基 选项 |  |
| `` R `` | 刷新 |  |
| `` + `` | 下一屏模式（正常/半屏/全屏） |  |
| `` _ `` | 上一屏模式 |  |
| `` ? `` | 打开菜单 |  |
| `` <c-s> `` | 查看按路径过滤选项 |  |
| `` W `` | 打开 diff 菜单 |  |
| `` <c-e> `` | 打开 diff 菜单 |  |
| `` <c-w> `` | 切换是否在差异视图中显示空白字符差异 |  |
| `` z `` | （通过 reflog）撤销「实验功能」 | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | （通过 reflog）重做「实验功能」 | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` P `` | 推送 |  |
| `` p `` | 拉取 |  |

## 列表面板导航

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | 上一页 |  |
| `` . `` | 下一页 |  |
| `` < `` | 滚动到顶部 |  |
| `` > `` | 滚动到底部 |  |
| `` v `` | 切换拖动选择 |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | 开始搜索 |  |
| `` H `` | 向左滚动 |  |
| `` L `` | 向右滚动 |  |
| `` ] `` | 下一个标签 |  |
| `` [ `` | 上一个标签 |  |

## Reflog 页面

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将提交的 SHA 复制到剪贴板 |  |
| `` w `` | View worktree options |  |
| `` <space> `` | 检出提交 |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | 在浏览器中打开提交 |  |
| `` n `` | 从提交创建新分支 |  |
| `` g `` | 查看重置选项 |  |
| `` C `` | 复制提交（拣选） |  |
| `` <c-r> `` | 重置已拣选（复制）的提交 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 查看提交 |  |
| `` / `` | Filter the current view by text |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | Create worktree |  |
| `` <space> `` | Switch to worktree |  |
| `` <enter> `` | Switch to worktree |  |
| `` o `` | Open in editor |  |
| `` d `` | Remove worktree |  |
| `` / `` | Filter the current view by text |  |

## 分支页面

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将分支名称复制到剪贴板 |  |
| `` i `` | 显示 git-flow 选项 |  |
| `` <space> `` | 检出 |  |
| `` n `` | 新分支 |  |
| `` o `` | 创建抓取请求 |  |
| `` O `` | 创建抓取请求选项 |  |
| `` <c-y> `` | 将抓取请求 URL 复制到剪贴板 |  |
| `` c `` | 按名称检出 |  |
| `` F `` | 强制检出 |  |
| `` d `` | View delete options |  |
| `` r `` | 将已检出的分支变基到该分支 |  |
| `` M `` | 合并到当前检出的分支 |  |
| `` f `` | 从上游快进此分支 |  |
| `` T `` | 创建标签 |  |
| `` s `` | Sort order |  |
| `` g `` | 查看重置选项 |  |
| `` R `` | 重命名分支 |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream |
| `` w `` | View worktree options |  |
| `` <enter> `` | 查看提交 |  |
| `` / `` | Filter the current view by text |  |

## 子提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将提交的 SHA 复制到剪贴板 |  |
| `` w `` | View worktree options |  |
| `` <space> `` | 检出提交 |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | 在浏览器中打开提交 |  |
| `` n `` | 从提交创建新分支 |  |
| `` g `` | 查看重置选项 |  |
| `` C `` | 复制提交（拣选） |  |
| `` <c-r> `` | 重置已拣选（复制）的提交 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 查看提交的文件 |  |
| `` / `` | 开始搜索 |  |

## 子模块

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将子模块名称复制到剪贴板 |  |
| `` <enter> `` | 输入子模块 |  |
| `` <space> `` | 输入子模块 |  |
| `` d `` | 删除子模块 |  |
| `` u `` | 更新子模块 |  |
| `` n `` | 添加新的子模块 |  |
| `` e `` | 更新子模块 URL |  |
| `` i `` | 初始化子模块 |  |
| `` b `` | 查看批量子模块选项 |  |
| `` / `` | Filter the current view by text |  |

## 提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将提交的 SHA 复制到剪贴板 |  |
| `` <c-r> `` | 重置已拣选（复制）的提交 |  |
| `` b `` | 查看二分查找选项 |  |
| `` s `` | 向下压缩 |  |
| `` f `` | 修正提交（fixup） |  |
| `` r `` | 改写提交 |  |
| `` R `` | 使用编辑器重命名提交 |  |
| `` d `` | 删除提交 |  |
| `` e `` | 编辑提交 |  |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | 选择提交（变基过程中） |  |
| `` F `` | 创建修正提交 |  |
| `` S `` | 压缩在所选提交之上的所有“fixup!”提交（自动压缩） |  |
| `` <c-j> `` | 下移提交 |  |
| `` <c-k> `` | 上移提交 |  |
| `` V `` | 粘贴提交（拣选） |  |
| `` B `` | Mark commit as base commit for rebase | Select a base commit for the next rebase; this will effectively perform a 'git rebase --onto'. |
| `` A `` | 用已暂存的更改来修补提交 |  |
| `` a `` | Set/Reset commit author |  |
| `` t `` | 还原提交 |  |
| `` T `` | 标签提交 |  |
| `` <c-l> `` | 打开日志菜单 |  |
| `` w `` | View worktree options |  |
| `` <space> `` | 检出提交 |  |
| `` y `` | Copy commit attribute |  |
| `` o `` | 在浏览器中打开提交 |  |
| `` n `` | 从提交创建新分支 |  |
| `` g `` | 查看重置选项 |  |
| `` C `` | 复制提交（拣选） |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 查看提交的文件 |  |
| `` / `` | 开始搜索 |  |

## 提交文件

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将提交的文件名复制到剪贴板 |  |
| `` c `` | 检出文件 |  |
| `` d `` | 放弃对此文件的提交更改 |  |
| `` o `` | 打开文件 |  |
| `` e `` | 编辑文件 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | 补丁中包含的切换文件 |  |
| `` a `` | Toggle all files included in patch |  |
| `` <enter> `` | 输入文件以将所选行添加到补丁中（或切换目录折叠） |  |
| `` ` `` | 切换文件树视图 |  |
| `` / `` | 开始搜索 |  |

## 提交讯息

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 确认 |  |
| `` <esc> `` | 关闭 |  |

## 文件

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将文件名复制到剪贴板 |  |
| `` <space> `` | 切换暂存状态 |  |
| `` <c-b> `` | Filter files by status |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | 提交更改 |  |
| `` w `` | 提交更改而无需预先提交钩子 |  |
| `` A `` | 修补最后一次提交 |  |
| `` C `` | 提交更改（使用编辑器编辑提交信息） |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | 编辑文件 |  |
| `` o `` | 打开文件 |  |
| `` i `` | 忽略文件 |  |
| `` r `` | 刷新文件 |  |
| `` s `` | 将所有更改加入贮藏 |  |
| `` S `` | 查看贮藏选项 |  |
| `` a `` | 切换所有文件的暂存状态 |  |
| `` <enter> `` | 暂存单个 块/行 用于文件, 或 折叠/展开 目录 |  |
| `` d `` | 查看'放弃更改'选项 |  |
| `` g `` | 查看上游重置选项 |  |
| `` D `` | 查看重置选项 |  |
| `` ` `` | 切换文件树视图 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | 打开外部合并工具 (git mergetool) |  |
| `` f `` | 抓取 |  |
| `` / `` | 开始搜索 |  |

## 构建补丁中

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 选择上一个区块 |  |
| `` <right> `` | 选择下一个区块 |  |
| `` v `` | 切换拖动选择 |  |
| `` a `` | 切换选择区块 |  |
| `` <c-o> `` | 将选中文本复制到剪贴板 |  |
| `` o `` | 打开文件 |  |
| `` e `` | 编辑文件 |  |
| `` <space> `` | 添加/移除 行到补丁 |  |
| `` <esc> `` | 退出逐行模式 |  |
| `` / `` | 开始搜索 |  |

## 标签页面

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 检出 |  |
| `` d `` | View delete options |  |
| `` P `` | 推送标签 |  |
| `` n `` | 创建标签 |  |
| `` g `` | 查看重置选项 |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | 查看提交 |  |
| `` / `` | Filter the current view by text |  |

## 正在合并

| Key | Action | Info |
|-----|--------|-------------|
| `` e `` | 编辑文件 |  |
| `` o `` | 打开文件 |  |
| `` <left> `` | 选择上一个冲突 |  |
| `` <right> `` | 选择下一个冲突 |  |
| `` <up> `` | 选择顶部块 |  |
| `` <down> `` | 选择底部块 |  |
| `` z `` | 撤销 |  |
| `` M `` | 打开外部合并工具 (git mergetool) |  |
| `` <space> `` | 选中区块 |  |
| `` b `` | 选中所有区块 |  |
| `` <esc> `` | 返回文件面板 |  |

## 正在暂存

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 选择上一个区块 |  |
| `` <right> `` | 选择下一个区块 |  |
| `` v `` | 切换拖动选择 |  |
| `` a `` | 切换选择区块 |  |
| `` <c-o> `` | 将选中文本复制到剪贴板 |  |
| `` o `` | 打开文件 |  |
| `` e `` | 编辑文件 |  |
| `` <esc> `` | 返回文件面板 |  |
| `` <tab> `` | 切换到其他面板 |  |
| `` <space> `` | 切换行暂存状态 |  |
| `` d `` | 取消变更 (git reset) |  |
| `` E `` | Edit hunk |  |
| `` c `` | 提交更改 |  |
| `` w `` | 提交更改而无需预先提交钩子 |  |
| `` C `` | 提交更改（使用编辑器编辑提交信息） |  |
| `` / `` | 开始搜索 |  |

## 正常

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | 向下滚动 |  |
| `` mouse wheel up (fn+down) `` | 向上滚动 |  |

## 状态

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 打开配置文件 |  |
| `` e `` | 编辑配置文件 |  |
| `` u `` | 检查更新 |  |
| `` <enter> `` | 切换到最近的仓库 |  |
| `` a `` | 显示所有分支的日志 |  |

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
| `` / `` | Filter the current view by text |  |

## 贮藏

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 应用 |  |
| `` g `` | 应用并删除 |  |
| `` d `` | 删除 |  |
| `` n `` | 新分支 |  |
| `` r `` | Rename stash |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | 查看提交的文件 |  |
| `` / `` | Filter the current view by text |  |

## 远程分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将分支名称复制到剪贴板 |  |
| `` <space> `` | 检出 |  |
| `` n `` | 新分支 |  |
| `` M `` | 合并到当前检出的分支 |  |
| `` r `` | 将已检出的分支变基到该分支 |  |
| `` d `` | Delete remote tag |  |
| `` u `` | 设置为检出分支的上游 |  |
| `` s `` | Sort order |  |
| `` g `` | 查看重置选项 |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | 查看提交 |  |
| `` / `` | Filter the current view by text |  |

## 远程页面

| Key | Action | Info |
|-----|--------|-------------|
| `` f `` | 抓取远程仓库 |  |
| `` n `` | 添加新的远程仓库 |  |
| `` d `` | 删除远程 |  |
| `` e `` | 编辑远程仓库 |  |
| `` / `` | Filter the current view by text |  |
