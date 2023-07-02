_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit 按键绑定

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## 全局键绑定

<pre>
  <kbd>&lt;c-r&gt;</kbd>: 切换到最近的仓库
  <kbd>&lt;pgup&gt;</kbd>: 向上滚动主面板 (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: 向下滚动主面板 (fn+down/shift+j)
  <kbd>@</kbd>: 打开命令日志菜单
  <kbd>}</kbd>: 扩大差异视图中显示的上下文范围
  <kbd>{</kbd>: 缩小差异视图中显示的上下文范围
  <kbd>:</kbd>: 执行自定义命令
  <kbd>&lt;c-p&gt;</kbd>: 查看自定义补丁选项
  <kbd>m</kbd>: 查看 合并/变基 选项
  <kbd>R</kbd>: 刷新
  <kbd>+</kbd>: 下一屏模式（正常/半屏/全屏）
  <kbd>_</kbd>: 上一屏模式
  <kbd>?</kbd>: 打开菜单
  <kbd>&lt;c-s&gt;</kbd>: 查看按路径过滤选项
  <kbd>W</kbd>: 打开 diff 菜单
  <kbd>&lt;c-e&gt;</kbd>: 打开 diff 菜单
  <kbd>&lt;c-w&gt;</kbd>: 切换是否在差异视图中显示空白字符差异
  <kbd>z</kbd>: （通过 reflog）撤销「实验功能」
  <kbd>&lt;c-z&gt;</kbd>: （通过 reflog）重做「实验功能」
  <kbd>P</kbd>: 推送
  <kbd>p</kbd>: 拉取
</pre>

## 列表面板导航

<pre>
  <kbd>,</kbd>: 上一页
  <kbd>.</kbd>: 下一页
  <kbd>&lt;</kbd>: 滚动到顶部
  <kbd>&gt;</kbd>: 滚动到底部
  <kbd>/</kbd>: 开始搜索
  <kbd>H</kbd>: 向左滚动
  <kbd>L</kbd>: 向右滚动
  <kbd>]</kbd>: 下一个标签
  <kbd>[</kbd>: 上一个标签
</pre>

## Reflog 页面

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 将提交的 SHA 复制到剪贴板
  <kbd>&lt;space&gt;</kbd>: 检出提交
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: 在浏览器中打开提交
  <kbd>n</kbd>: 从提交创建新分支
  <kbd>g</kbd>: 查看重置选项
  <kbd>c</kbd>: 复制提交（拣选）
  <kbd>C</kbd>: 复制提交范围（拣选）
  <kbd>&lt;c-r&gt;</kbd>: 重置已拣选（复制）的提交
  <kbd>&lt;enter&gt;</kbd>: 查看提交
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 分支页面

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 将分支名称复制到剪贴板
  <kbd>i</kbd>: 显示 git-flow 选项
  <kbd>&lt;space&gt;</kbd>: 检出
  <kbd>n</kbd>: 新分支
  <kbd>o</kbd>: 创建抓取请求
  <kbd>O</kbd>: 创建抓取请求选项
  <kbd>&lt;c-y&gt;</kbd>: 将抓取请求 URL 复制到剪贴板
  <kbd>c</kbd>: 按名称检出
  <kbd>F</kbd>: 强制检出
  <kbd>d</kbd>: 删除分支
  <kbd>r</kbd>: 将已检出的分支变基到该分支
  <kbd>M</kbd>: 合并到当前检出的分支
  <kbd>f</kbd>: 从上游快进此分支
  <kbd>T</kbd>: 创建标签
  <kbd>g</kbd>: 查看重置选项
  <kbd>R</kbd>: 重命名分支
  <kbd>u</kbd>: Set/Unset upstream
  <kbd>&lt;enter&gt;</kbd>: 查看提交
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 子提交

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 将提交的 SHA 复制到剪贴板
  <kbd>&lt;space&gt;</kbd>: 检出提交
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: 在浏览器中打开提交
  <kbd>n</kbd>: 从提交创建新分支
  <kbd>g</kbd>: 查看重置选项
  <kbd>c</kbd>: 复制提交（拣选）
  <kbd>C</kbd>: 复制提交范围（拣选）
  <kbd>&lt;c-r&gt;</kbd>: 重置已拣选（复制）的提交
  <kbd>&lt;enter&gt;</kbd>: 查看提交的文件
  <kbd>/</kbd>: 开始搜索
</pre>

## 子模块

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 将子模块名称复制到剪贴板
  <kbd>&lt;enter&gt;</kbd>: 输入子模块
  <kbd>d</kbd>: 删除子模块
  <kbd>u</kbd>: 更新子模块
  <kbd>n</kbd>: 添加新的子模块
  <kbd>e</kbd>: 更新子模块 URL
  <kbd>i</kbd>: 初始化子模块
  <kbd>b</kbd>: 查看批量子模块选项
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 提交

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 将提交的 SHA 复制到剪贴板
  <kbd>&lt;c-r&gt;</kbd>: 重置已拣选（复制）的提交
  <kbd>b</kbd>: 查看二分查找选项
  <kbd>s</kbd>: 向下压缩
  <kbd>f</kbd>: 修正提交（fixup）
  <kbd>r</kbd>: 改写提交
  <kbd>R</kbd>: 使用编辑器重命名提交
  <kbd>d</kbd>: 删除提交
  <kbd>e</kbd>: 编辑提交
  <kbd>p</kbd>: 选择提交（变基过程中）
  <kbd>F</kbd>: 创建修正提交
  <kbd>S</kbd>: 压缩在所选提交之上的所有“fixup!”提交（自动压缩）
  <kbd>&lt;c-j&gt;</kbd>: 下移提交
  <kbd>&lt;c-k&gt;</kbd>: 上移提交
  <kbd>v</kbd>: 粘贴提交（拣选）
  <kbd>A</kbd>: 用已暂存的更改来修补提交
  <kbd>a</kbd>: Set/Reset commit author
  <kbd>t</kbd>: 还原提交
  <kbd>T</kbd>: 标签提交
  <kbd>&lt;c-l&gt;</kbd>: 打开日志菜单
  <kbd>&lt;space&gt;</kbd>: 检出提交
  <kbd>y</kbd>: Copy commit attribute
  <kbd>o</kbd>: 在浏览器中打开提交
  <kbd>n</kbd>: 从提交创建新分支
  <kbd>g</kbd>: 查看重置选项
  <kbd>c</kbd>: 复制提交（拣选）
  <kbd>C</kbd>: 复制提交范围（拣选）
  <kbd>&lt;enter&gt;</kbd>: 查看提交的文件
  <kbd>/</kbd>: 开始搜索
</pre>

## 提交文件

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 将提交的文件名复制到剪贴板
  <kbd>c</kbd>: 检出文件
  <kbd>d</kbd>: 放弃对此文件的提交更改
  <kbd>o</kbd>: 打开文件
  <kbd>e</kbd>: 编辑文件
  <kbd>&lt;space&gt;</kbd>: 补丁中包含的切换文件
  <kbd>a</kbd>: Toggle all files included in patch
  <kbd>&lt;enter&gt;</kbd>: 输入文件以将所选行添加到补丁中（或切换目录折叠）
  <kbd>`</kbd>: 切换文件树视图
  <kbd>/</kbd>: 开始搜索
</pre>

## 提交讯息

<pre>
  <kbd>&lt;enter&gt;</kbd>: 确认
  <kbd>&lt;esc&gt;</kbd>: 关闭
</pre>

## 文件

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 将文件名复制到剪贴板
  <kbd>d</kbd>: 查看'放弃更改'选项
  <kbd>&lt;space&gt;</kbd>: 切换暂存状态
  <kbd>&lt;c-b&gt;</kbd>: Filter files by status
  <kbd>c</kbd>: 提交更改
  <kbd>w</kbd>: 提交更改而无需预先提交钩子
  <kbd>A</kbd>: 修补最后一次提交
  <kbd>C</kbd>: 提交更改（使用编辑器编辑提交信息）
  <kbd>e</kbd>: 编辑文件
  <kbd>o</kbd>: 打开文件
  <kbd>i</kbd>: 忽略文件
  <kbd>r</kbd>: 刷新文件
  <kbd>s</kbd>: 将所有更改加入贮藏
  <kbd>S</kbd>: 查看贮藏选项
  <kbd>a</kbd>: 切换所有文件的暂存状态
  <kbd>&lt;enter&gt;</kbd>: 暂存单个 块/行 用于文件, 或 折叠/展开 目录
  <kbd>g</kbd>: 查看上游重置选项
  <kbd>D</kbd>: 查看重置选项
  <kbd>`</kbd>: 切换文件树视图
  <kbd>M</kbd>: 打开外部合并工具 (git mergetool)
  <kbd>f</kbd>: 抓取
  <kbd>/</kbd>: 开始搜索
</pre>

## 构建补丁中

<pre>
  <kbd>&lt;left&gt;</kbd>: 选择上一个区块
  <kbd>&lt;right&gt;</kbd>: 选择下一个区块
  <kbd>v</kbd>: 切换拖动选择
  <kbd>V</kbd>: 切换拖动选择
  <kbd>a</kbd>: 切换选择区块
  <kbd>&lt;c-o&gt;</kbd>: 将选中文本复制到剪贴板
  <kbd>o</kbd>: 打开文件
  <kbd>e</kbd>: 编辑文件
  <kbd>&lt;space&gt;</kbd>: 添加/移除 行到补丁
  <kbd>&lt;esc&gt;</kbd>: 退出逐行模式
  <kbd>/</kbd>: 开始搜索
</pre>

## 标签页面

<pre>
  <kbd>&lt;space&gt;</kbd>: 检出
  <kbd>d</kbd>: 删除标签
  <kbd>P</kbd>: 推送标签
  <kbd>n</kbd>: 创建标签
  <kbd>g</kbd>: 查看重置选项
  <kbd>&lt;enter&gt;</kbd>: 查看提交
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 正在合并

<pre>
  <kbd>e</kbd>: 编辑文件
  <kbd>o</kbd>: 打开文件
  <kbd>&lt;left&gt;</kbd>: 选择上一个冲突
  <kbd>&lt;right&gt;</kbd>: 选择下一个冲突
  <kbd>&lt;up&gt;</kbd>: 选择顶部块
  <kbd>&lt;down&gt;</kbd>: 选择底部块
  <kbd>z</kbd>: 撤销
  <kbd>M</kbd>: 打开外部合并工具 (git mergetool)
  <kbd>&lt;space&gt;</kbd>: 选中区块
  <kbd>b</kbd>: 选中所有区块
  <kbd>&lt;esc&gt;</kbd>: 返回文件面板
</pre>

## 正在暂存

<pre>
  <kbd>&lt;left&gt;</kbd>: 选择上一个区块
  <kbd>&lt;right&gt;</kbd>: 选择下一个区块
  <kbd>v</kbd>: 切换拖动选择
  <kbd>V</kbd>: 切换拖动选择
  <kbd>a</kbd>: 切换选择区块
  <kbd>&lt;c-o&gt;</kbd>: 将选中文本复制到剪贴板
  <kbd>o</kbd>: 打开文件
  <kbd>e</kbd>: 编辑文件
  <kbd>&lt;esc&gt;</kbd>: 返回文件面板
  <kbd>&lt;tab&gt;</kbd>: 切换到其他面板
  <kbd>&lt;space&gt;</kbd>: 切换行暂存状态
  <kbd>d</kbd>: 取消变更 (git reset)
  <kbd>E</kbd>: Edit hunk
  <kbd>c</kbd>: 提交更改
  <kbd>w</kbd>: 提交更改而无需预先提交钩子
  <kbd>C</kbd>: 提交更改（使用编辑器编辑提交信息）
  <kbd>/</kbd>: 开始搜索
</pre>

## 正常

<pre>
  <kbd>mouse wheel down</kbd>: 向下滚动 (fn+up)
  <kbd>mouse wheel up</kbd>: 向上滚动 (fn+down)
</pre>

## 状态

<pre>
  <kbd>o</kbd>: 打开配置文件
  <kbd>e</kbd>: 编辑配置文件
  <kbd>u</kbd>: 检查更新
  <kbd>&lt;enter&gt;</kbd>: 切换到最近的仓库
  <kbd>a</kbd>: 显示所有分支的日志
</pre>

## 确认面板

<pre>
  <kbd>&lt;enter&gt;</kbd>: 确认
  <kbd>&lt;esc&gt;</kbd>: 关闭
</pre>

## 菜单

<pre>
  <kbd>&lt;enter&gt;</kbd>: 执行
  <kbd>&lt;esc&gt;</kbd>: 关闭
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 贮藏

<pre>
  <kbd>&lt;space&gt;</kbd>: 应用
  <kbd>g</kbd>: 应用并删除
  <kbd>d</kbd>: 删除
  <kbd>n</kbd>: 新分支
  <kbd>r</kbd>: Rename stash
  <kbd>&lt;enter&gt;</kbd>: 查看提交的文件
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 远程分支

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 将分支名称复制到剪贴板
  <kbd>&lt;space&gt;</kbd>: 检出
  <kbd>n</kbd>: 新分支
  <kbd>M</kbd>: 合并到当前检出的分支
  <kbd>r</kbd>: 将已检出的分支变基到该分支
  <kbd>d</kbd>: 删除分支
  <kbd>u</kbd>: 设置为检出分支的上游
  <kbd>g</kbd>: 查看重置选项
  <kbd>&lt;enter&gt;</kbd>: 查看提交
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 远程页面

<pre>
  <kbd>f</kbd>: 抓取远程仓库
  <kbd>n</kbd>: 添加新的远程仓库
  <kbd>d</kbd>: 删除远程
  <kbd>e</kbd>: 编辑远程仓库
  <kbd>/</kbd>: Filter the current view by text
</pre>
