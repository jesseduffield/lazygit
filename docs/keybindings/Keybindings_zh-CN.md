_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit 按键绑定

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## 全局键绑定

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | 切换到最近的仓库 |  |
| `` <pgup> (fn+up/shift+k) `` | 向上滚动主面板 |  |
| `` <pgdown> (fn+down/shift+j) `` | 向下滚动主面板 |  |
| `` @ `` | 打开命令日志菜单 | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | 推送 | Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` p `` | 拉取 | Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` } `` | 扩大差异视图中显示的上下文范围 | Increase the amount of the context shown around changes in the diff view. |
| `` { `` | 缩小差异视图中显示的上下文范围 | Decrease the amount of the context shown around changes in the diff view. |
| `` : `` | 执行自定义命令 | Bring up a prompt where you can enter a shell command to execute. Not to be confused with pre-configured custom commands. |
| `` <c-p> `` | 查看自定义补丁选项 |  |
| `` m `` | 查看 合并/变基 选项 | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | 刷新 | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | 下一屏模式（正常/半屏/全屏） |  |
| `` _ `` | 上一屏模式 |  |
| `` ? `` | 打开菜单 |  |
| `` <c-s> `` | 查看按路径过滤选项 | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | 打开 diff 菜单 | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | 打开 diff 菜单 | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | 退出 |  |
| `` <esc> `` | 取消 |  |
| `` <c-w> `` | 切换是否在差异视图中显示空白字符差异 | Toggle whether or not whitespace changes are shown in the diff view. |
| `` z `` | （通过 reflog）撤销「实验功能」 | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | （通过 reflog）重做「实验功能」 | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |

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
| `` <c-o> `` | 将提交的 hash 复制到剪贴板 |  |
| `` <space> `` | 检出 | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 在浏览器中打开提交 |  |
| `` n `` | 从提交创建新分支 |  |
| `` g `` | 查看重置选项 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | 复制提交（拣选） | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | 重置已拣选（复制）的提交 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 查看提交 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | Open in editor |  |
| `` d `` | Remove | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | Filter the current view by text |  |

## 分支页面

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将分支名称复制到剪贴板 |  |
| `` i `` | 显示 git-flow 选项 |  |
| `` <space> `` | 检出 | Checkout selected item. |
| `` n `` | 新分支 |  |
| `` o `` | 创建抓取请求 |  |
| `` O `` | 创建抓取请求选项 |  |
| `` <c-y> `` | 将抓取请求 URL 复制到剪贴板 |  |
| `` c `` | 按名称检出 | Checkout by name. In the input box you can enter '-' to switch to the last branch. |
| `` F `` | 强制检出 | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | Delete | View delete options for local/remote branch. |
| `` r `` | 将已检出的分支变基到该分支 | Rebase the checked-out branch onto the selected branch. |
| `` M `` | 合并到当前检出的分支 | Merge selected branch into currently checked out branch. |
| `` f `` | 从上游快进此分支 | Fast-forward selected branch from its upstream. |
| `` T `` | 创建标签 |  |
| `` s `` | Sort order |  |
| `` g `` | 查看重置选项 |  |
| `` R `` | 重命名分支 |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <enter> `` | 查看提交 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## 子提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将提交的 hash 复制到剪贴板 |  |
| `` <space> `` | 检出 | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 在浏览器中打开提交 |  |
| `` n `` | 从提交创建新分支 |  |
| `` g `` | 查看重置选项 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | 复制提交（拣选） | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | 重置已拣选（复制）的提交 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 查看提交的文件 |  |
| `` w `` | View worktree options |  |
| `` / `` | 开始搜索 |  |

## 子模块

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将子模块名称复制到剪贴板 |  |
| `` <enter> `` | Enter | 输入子模块 |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | 更新子模块 |
| `` n `` | 添加新的子模块 |  |
| `` e `` | 更新子模块 URL |  |
| `` i `` | Initialize | 初始化子模块 |
| `` b `` | 查看批量子模块选项 |  |
| `` / `` | Filter the current view by text |  |

## 提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将提交的 hash 复制到剪贴板 |  |
| `` <c-r> `` | 重置已拣选（复制）的提交 |  |
| `` b `` | 查看二分查找选项 |  |
| `` s `` | 压缩 | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | 修正（fixup） | Meld the selected commit into the commit below it. Similar to fixup, but the selected commit's message will be discarded. |
| `` r `` | 改写提交 | Reword the selected commit's message. |
| `` R `` | 使用编辑器重命名提交 |  |
| `` d `` | 删除提交 | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | Edit (start interactive rebase) | 编辑提交 |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick | 选择提交（变基过程中） |
| `` F `` | 为此提交创建修正 | 创建修正提交 |
| `` S `` | Apply fixup commits | 压缩在所选提交之上的所有“fixup!”提交（自动压缩） |
| `` <c-j> `` | 下移提交 |  |
| `` <c-k> `` | 上移提交 |  |
| `` V `` | 粘贴提交（拣选） |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Amend | 用已暂存的更改来修补提交 |
| `` a `` | Amend commit attribute | Set/Reset commit author or set co-author. |
| `` t `` | Revert | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | 标签提交 | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <c-l> `` | 打开日志菜单 | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` <space> `` | 检出 | Checkout the selected commit as a detached HEAD. |
| `` y `` | Copy commit attribute to clipboard | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 在浏览器中打开提交 |  |
| `` n `` | 从提交创建新分支 |  |
| `` g `` | 查看重置选项 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | 复制提交（拣选） | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 查看提交的文件 |  |
| `` w `` | View worktree options |  |
| `` / `` | 开始搜索 |  |

## 提交文件

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将文件名复制到剪贴板 |  |
| `` c `` | 检出 | 检出文件 |
| `` d `` | Remove | 放弃对此文件的提交更改 |
| `` o `` | 打开文件 | Open file in default application. |
| `` e `` | Edit | Open file in external editor. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | 补丁中包含的切换文件 | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Toggle all files | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | 输入文件以将所选行添加到补丁中（或切换目录折叠） | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | 切换文件树视图 | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
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
| `` <space> `` | 切换暂存状态 | Toggle staged for selected file. |
| `` <c-b> `` | Filter files by status |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | 提交更改 | Commit staged changes. |
| `` w `` | 提交更改而无需预先提交钩子 |  |
| `` A `` | 修补最后一次提交 |  |
| `` C `` | 提交更改（使用编辑器编辑提交信息） |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit | Open file in external editor. |
| `` o `` | 打开文件 | Open file in default application. |
| `` i `` | 忽略文件 |  |
| `` r `` | 刷新文件 |  |
| `` s `` | Stash | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | 查看贮藏选项 | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | 切换所有文件的暂存状态 | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | 暂存单个 块/行 用于文件, 或 折叠/展开 目录 | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | 查看'放弃更改'选项 | View options for discarding changes to the selected file. |
| `` g `` | 查看上游重置选项 |  |
| `` D `` | Reset | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | 切换文件树视图 | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | 打开外部合并工具 (git mergetool) | Run `git mergetool`. |
| `` f `` | 抓取 | Fetch changes from remote. |
| `` / `` | 开始搜索 |  |

## 构建补丁中

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 选择上一个区块 |  |
| `` <right> `` | 选择下一个区块 |  |
| `` v `` | 切换拖动选择 |  |
| `` a `` | 切换选择区块 | Toggle hunk selection mode. |
| `` <c-o> `` | 将选中文本复制到剪贴板 |  |
| `` o `` | 打开文件 | Open file in default application. |
| `` e `` | 编辑文件 | Open file in external editor. |
| `` <space> `` | 添加/移除 行到补丁 |  |
| `` <esc> `` | 退出逐行模式 |  |
| `` / `` | 开始搜索 |  |

## 标签页面

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 检出 | Checkout the selected tag tag as a detached HEAD. |
| `` n `` | 创建标签 | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Delete | View delete options for local/remote tag. |
| `` P `` | 推送标签 | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <enter> `` | 查看提交 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## 正在合并

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 选中区块 |  |
| `` b `` | 选中所有区块 |  |
| `` <up> `` | 选择顶部块 |  |
| `` <down> `` | 选择底部块 |  |
| `` <left> `` | 选择上一个冲突 |  |
| `` <right> `` | 选择下一个冲突 |  |
| `` z `` | 撤销 | Undo last merge conflict resolution. |
| `` e `` | 编辑文件 | Open file in external editor. |
| `` o `` | 打开文件 | Open file in default application. |
| `` M `` | 打开外部合并工具 (git mergetool) | Run `git mergetool`. |
| `` <esc> `` | 返回文件面板 |  |

## 正在暂存

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 选择上一个区块 |  |
| `` <right> `` | 选择下一个区块 |  |
| `` v `` | 切换拖动选择 |  |
| `` a `` | 切换选择区块 | Toggle hunk selection mode. |
| `` <c-o> `` | 将选中文本复制到剪贴板 |  |
| `` <space> `` | 切换暂存状态 | 切换行暂存状态 |
| `` d `` | 取消变更 (git reset) | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | 打开文件 | Open file in default application. |
| `` e `` | 编辑文件 | Open file in external editor. |
| `` <esc> `` | 返回文件面板 |  |
| `` <tab> `` | 切换到其他面板 | Switch to other view (staged/unstaged changes). |
| `` E `` | Edit hunk | Edit selected hunk in external editor. |
| `` c `` | 提交更改 | Commit staged changes. |
| `` w `` | 提交更改而无需预先提交钩子 |  |
| `` C `` | 提交更改（使用编辑器编辑提交信息） |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | 开始搜索 |  |

## 正常

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | 向下滚动 |  |
| `` mouse wheel up (fn+down) `` | 向上滚动 |  |

## 状态

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 打开配置文件 | Open file in default application. |
| `` e `` | 编辑配置文件 | Open file in external editor. |
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
| `` <space> `` | 应用 | Apply the stash entry to your working directory. |
| `` g `` | 应用并删除 | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | 删除 | Remove the stash entry from the stash list. |
| `` n `` | 新分支 | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | Rename stash |  |
| `` <enter> `` | 查看提交的文件 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## 远程分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 将分支名称复制到剪贴板 |  |
| `` <space> `` | 检出 | Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head. |
| `` n `` | 新分支 |  |
| `` M `` | 合并到当前检出的分支 | Merge selected branch into currently checked out branch. |
| `` r `` | 将已检出的分支变基到该分支 | Rebase the checked-out branch onto the selected branch. |
| `` d `` | Delete | Delete the remote branch from the remote. |
| `` u `` | Set as upstream | 设置为检出分支的上游 |
| `` s `` | Sort order |  |
| `` g `` | 查看重置选项 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <enter> `` | 查看提交 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## 远程页面

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | 添加新的远程仓库 |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Edit | 编辑远程仓库 |
| `` f `` | 抓取 | 抓取远程仓库 |
| `` / `` | Filter the current view by text |  |
