_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit 按键绑定

## 全局键绑定

<pre>
  <kbd>ctrl+r</kbd>: 切换到最近的仓库
  <kbd>pgup</kbd>: 向上滚动主面板 (fn+up/shift+k)
  <kbd>pgdown</kbd>: 向下滚动主面板 (fn+down/shift+j)
  <kbd>m</kbd>: 查看 合并/变基 选项
  <kbd>ctrl+p</kbd>: 查看自定义补丁选项
  <kbd>R</kbd>: 刷新
  <kbd>x</kbd>: 打开菜单
  <kbd>+</kbd>: 下一屏模式（正常/半屏/全屏）
  <kbd>_</kbd>: 上一屏模式
  <kbd>ctrl+s</kbd>: 查看按路径过滤选项
  <kbd>W</kbd>: 打开 diff 菜单
  <kbd>ctrl+e</kbd>: 打开 diff 菜单
  <kbd>@</kbd>: 打开命令日志菜单
  <kbd>}</kbd>: 扩大差异视图中显示的上下文范围
  <kbd>{</kbd>: 缩小差异视图中显示的上下文范围
  <kbd>:</kbd>: 执行自定义命令
  <kbd>z</kbd>: （通过 reflog）撤销「实验功能」
  <kbd>ctrl+z</kbd>: （通过 reflog）重做「实验功能」
  <kbd>P</kbd>: 推送
  <kbd>p</kbd>: 拉取
</pre>

## 列表面板导航

<pre>
  <kbd>,</kbd>: 上一页
  <kbd>.</kbd>: 下一页
  <kbd><</kbd>: 滚动到顶部
  <kbd>/</kbd>: 开始搜索
  <kbd>></kbd>: 滚动到底部
  <kbd>H</kbd>: 向左滚动
  <kbd>L</kbd>: 向右滚动
  <kbd>]</kbd>: 下一个标签
  <kbd>[</kbd>: 上一个标签
</pre>

## Reflog 页面

<pre>
  <kbd>ctrl+o</kbd>: 将提交的 SHA 复制到剪贴板
  <kbd>space</kbd>: 检出提交
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: 在浏览器中打开提交
  <kbd>n</kbd>: 从提交创建新分支
  <kbd>g</kbd>: 查看重置选项
  <kbd>c</kbd>: 复制提交（拣选）
  <kbd>C</kbd>: 复制提交范围（拣选）
  <kbd>ctrl+r</kbd>: 重置已拣选（复制）的提交
  <kbd>enter</kbd>: 查看提交
</pre>

## 分支页面

<pre>
  <kbd>ctrl+o</kbd>: 将分支名称复制到剪贴板
  <kbd>i</kbd>: 显示 git-flow 选项
  <kbd>space</kbd>: 检出
  <kbd>n</kbd>: 新分支
  <kbd>o</kbd>: 创建抓取请求
  <kbd>O</kbd>: 创建抓取请求选项
  <kbd>ctrl+y</kbd>: 将抓取请求 URL 复制到剪贴板
  <kbd>c</kbd>: 按名称检出
  <kbd>F</kbd>: 强制检出
  <kbd>d</kbd>: 删除分支
  <kbd>r</kbd>: 将已检出的分支变基到该分支
  <kbd>M</kbd>: 合并到当前检出的分支
  <kbd>f</kbd>: 从上游快进此分支
  <kbd>g</kbd>: 查看重置选项
  <kbd>R</kbd>: 重命名分支
  <kbd>u</kbd>: set/unset upstream
  <kbd>enter</kbd>: 查看提交
</pre>

## 子提交

<pre>
  <kbd>ctrl+o</kbd>: 将提交的 SHA 复制到剪贴板
  <kbd>space</kbd>: 检出提交
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: 在浏览器中打开提交
  <kbd>n</kbd>: 从提交创建新分支
  <kbd>g</kbd>: 查看重置选项
  <kbd>c</kbd>: 复制提交（拣选）
  <kbd>C</kbd>: 复制提交范围（拣选）
  <kbd>ctrl+r</kbd>: 重置已拣选（复制）的提交
  <kbd>enter</kbd>: 查看提交的文件
</pre>

## 子模块

<pre>
  <kbd>ctrl+o</kbd>: 将子模块名称复制到剪贴板
  <kbd>enter</kbd>: 输入子模块
  <kbd>d</kbd>: 删除子模块
  <kbd>u</kbd>: 更新子模块
  <kbd>n</kbd>: 添加新的子模块
  <kbd>e</kbd>: 更新子模块 URL
  <kbd>i</kbd>: 初始化子模块
  <kbd>b</kbd>: 查看批量子模块选项
</pre>

## 提交

<pre>
  <kbd>ctrl+o</kbd>: 将提交的 SHA 复制到剪贴板
  <kbd>ctrl+r</kbd>: 重置已拣选（复制）的提交
  <kbd>b</kbd>: 查看二分查找选项
  <kbd>s</kbd>: 向下压缩
  <kbd>f</kbd>: 修正提交（fixup）
  <kbd>r</kbd>: 改写提交
  <kbd>R</kbd>: 使用编辑器重命名提交
  <kbd>d</kbd>: 删除提交
  <kbd>e</kbd>: 编辑提交
  <kbd>p</kbd>: 选择提交（变基过程中）
  <kbd>F</kbd>: 为此提交创建修正
  <kbd>S</kbd>: 压缩在所选提交之上的所有“fixup!”提交（自动压缩）
  <kbd>ctrl+j</kbd>: 下移提交
  <kbd>ctrl+k</kbd>: 上移提交
  <kbd>v</kbd>: 粘贴提交（拣选）
  <kbd>A</kbd>: 用已暂存的更改来修补提交
  <kbd>a</kbd>: reset commit author
  <kbd>t</kbd>: 还原提交
  <kbd>T</kbd>: 标签提交
  <kbd>ctrl+l</kbd>: 打开日志菜单
  <kbd>space</kbd>: 检出提交
  <kbd>y</kbd>: copy commit attribute
  <kbd>o</kbd>: 在浏览器中打开提交
  <kbd>n</kbd>: 从提交创建新分支
  <kbd>g</kbd>: 查看重置选项
  <kbd>c</kbd>: 复制提交（拣选）
  <kbd>C</kbd>: 复制提交范围（拣选）
  <kbd>enter</kbd>: 查看提交的文件
</pre>

## 提交文件

<pre>
  <kbd>ctrl+o</kbd>: 将提交的文件名复制到剪贴板
  <kbd>c</kbd>: 检出文件
  <kbd>d</kbd>: 放弃对此文件的提交更改
  <kbd>o</kbd>: 打开文件
  <kbd>e</kbd>: 编辑文件
  <kbd>space</kbd>: 补丁中包含的切换文件
  <kbd>a</kbd>: toggle all files included in patch
  <kbd>enter</kbd>: 输入文件以将所选行添加到补丁中（或切换目录折叠）
  <kbd>`</kbd>: 切换文件树视图
</pre>

## 文件

<pre>
  <kbd>ctrl+o</kbd>: 将文件名复制到剪贴板
  <kbd>ctrl+w</kbd>: 切换是否在差异视图中显示空白字符差异
  <kbd>d</kbd>: 查看'放弃更改'选项
  <kbd>space</kbd>: 切换暂存状态
  <kbd>ctrl+b</kbd>: Filter files (staged/unstaged)
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
  <kbd>enter</kbd>: 暂存单个 块/行 用于文件, 或 折叠/展开 目录
  <kbd>g</kbd>: 查看上游重置选项
  <kbd>D</kbd>: 查看重置选项
  <kbd>`</kbd>: 切换文件树视图
  <kbd>M</kbd>: 打开外部合并工具 (git mergetool)
  <kbd>f</kbd>: 抓取
</pre>

## 构建补丁中

<pre>
  <kbd>◄</kbd>: 选择上一个区块
  <kbd>►</kbd>: 选择下一个区块
  <kbd>v</kbd>: 切换拖动选择
  <kbd>V</kbd>: 切换拖动选择
  <kbd>a</kbd>: 切换选择区块
  <kbd>ctrl+o</kbd>: 将选中文本复制到剪贴板
  <kbd>o</kbd>: 打开文件
  <kbd>e</kbd>: 编辑文件
  <kbd>space</kbd>: 添加/移除 行到补丁
  <kbd>esc</kbd>: 退出逐行模式
</pre>

## 标签页面

<pre>
  <kbd>space</kbd>: 检出
  <kbd>d</kbd>: 删除标签
  <kbd>P</kbd>: 推送标签
  <kbd>n</kbd>: 创建标签
  <kbd>g</kbd>: 查看重置选项
  <kbd>enter</kbd>: 查看提交
</pre>

## 正在合并

<pre>
  <kbd>e</kbd>: 编辑文件
  <kbd>o</kbd>: 打开文件
  <kbd>◄</kbd>: 选择上一个冲突
  <kbd>►</kbd>: 选择下一个冲突
  <kbd>▲</kbd>: 选择顶部块
  <kbd>▼</kbd>: 选择底部块
  <kbd>z</kbd>: 撤销
  <kbd>M</kbd>: 打开外部合并工具 (git mergetool)
  <kbd>space</kbd>: 选中区块
  <kbd>b</kbd>: 选中所有区块
  <kbd>esc</kbd>: 返回文件面板
</pre>

## 正在暂存

<pre>
  <kbd>◄</kbd>: 选择上一个区块
  <kbd>►</kbd>: 选择下一个区块
  <kbd>v</kbd>: 切换拖动选择
  <kbd>V</kbd>: 切换拖动选择
  <kbd>a</kbd>: 切换选择区块
  <kbd>ctrl+o</kbd>: 将选中文本复制到剪贴板
  <kbd>o</kbd>: 打开文件
  <kbd>e</kbd>: 编辑文件
  <kbd>esc</kbd>: 返回文件面板
  <kbd>tab</kbd>: 切换到其他面板
  <kbd>space</kbd>: 切换行暂存状态
  <kbd>d</kbd>: 取消变更 (git reset)
  <kbd>E</kbd>: edit hunk
  <kbd>c</kbd>: 提交更改
  <kbd>w</kbd>: 提交更改而无需预先提交钩子
  <kbd>C</kbd>: 提交更改（使用编辑器编辑提交信息）
</pre>

## 正常

<pre>
  <kbd>mouse wheel ▼</kbd>: 向下滚动 (fn+up)
  <kbd>mouse wheel ▲</kbd>: 向上滚动 (fn+down)
</pre>

## 状态

<pre>
  <kbd>e</kbd>: 编辑配置文件
  <kbd>o</kbd>: 打开配置文件
  <kbd>u</kbd>: 检查更新
  <kbd>enter</kbd>: 切换到最近的仓库
  <kbd>a</kbd>: 显示所有分支的日志
</pre>

## 贮藏

<pre>
  <kbd>space</kbd>: 应用
  <kbd>g</kbd>: 应用并删除
  <kbd>d</kbd>: 删除
  <kbd>n</kbd>: 新分支
  <kbd>r</kbd>: rename stash
  <kbd>enter</kbd>: 查看提交的文件
</pre>

## 远程分支

<pre>
  <kbd>space</kbd>: 检出
  <kbd>n</kbd>: 新分支
  <kbd>M</kbd>: 合并到当前检出的分支
  <kbd>r</kbd>: 将已检出的分支变基到该分支
  <kbd>d</kbd>: 删除分支
  <kbd>u</kbd>: 设置为检出分支的上游
  <kbd>esc</kbd>: 返回远程仓库列表
  <kbd>g</kbd>: 查看重置选项
  <kbd>enter</kbd>: 查看提交
</pre>

## 远程页面

<pre>
  <kbd>f</kbd>: 抓取远程仓库
  <kbd>n</kbd>: 添加新的远程仓库
  <kbd>d</kbd>: 删除远程
  <kbd>e</kbd>: 编辑远程仓库
</pre>
