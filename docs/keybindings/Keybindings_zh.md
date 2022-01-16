_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit 按键绑定

## 全局键绑定

<pre>
  <kbd>ctrl+r</kbd>: 切换到最近的仓库 (<c-r>)
  <kbd>pgup</kbd>: 向上滚动主面板 (fn+up)
  <kbd>pgdown</kbd>: 向下滚动主面板 (fn+down)
  <kbd>m</kbd>: 查看 合并/变基 选项
  <kbd>ctrl+p</kbd>: 查看自定义补丁选项
  <kbd>P</kbd>: 推送
  <kbd>p</kbd>: 拉取
  <kbd>R</kbd>: 刷新
  <kbd>x</kbd>: 打开菜单
  <kbd>z</kbd>: （通过 reflog）撤销「实验功能」
  <kbd>ctrl+z</kbd>: （通过 reflog）重做「实验功能」
  <kbd>+</kbd>: 下一屏模式（正常/半屏/全屏）
  <kbd>_</kbd>: 上一屏模式
  <kbd>:</kbd>: 执行自定义命令
  <kbd>ctrl+s</kbd>: 查看按路径过滤选项
  <kbd>W</kbd>: 打开 diff 菜单
  <kbd>ctrl+e</kbd>: 打开 diff 菜单
  <kbd>@</kbd>: 打开命令日志菜单
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
</pre>

## 列表面板导航

<pre>
  <kbd>.</kbd>: 下一页
  <kbd>,</kbd>: 上一页
  <kbd><</kbd>: 滚动到顶部
  <kbd>></kbd>: 滚动到底部
  <kbd>/</kbd>: 开始搜索
  <kbd>]</kbd>: 下一个标签
  <kbd>[</kbd>: 上一个标签
</pre>

## 分支 面板 (分支标签)

<pre>
  <kbd>space</kbd>: 检出
  <kbd>o</kbd>: 创建抓取请求
  <kbd>O</kbd>: 创建抓取请求选项
  <kbd>ctrl+y</kbd>: 将抓取请求 URL 复制到剪贴板
  <kbd>c</kbd>: 按名称检出
  <kbd>F</kbd>: 强制检出
  <kbd>n</kbd>: 新分支
  <kbd>d</kbd>: 删除分支
  <kbd>r</kbd>: 将已检出的分支变基到该分支
  <kbd>M</kbd>: 合并到当前检出的分支
  <kbd>i</kbd>: 显示 git-flow 选项
  <kbd>f</kbd>: 从上游快进此分支
  <kbd>g</kbd>: 查看重置选项
  <kbd>R</kbd>: 重命名分支
  <kbd>ctrl+o</kbd>: 将分支名称复制到剪贴板
  <kbd>enter</kbd>: 查看提交
</pre>

## 分支 面板 (远程分支（在远程页面中）)

<pre>
  <kbd>esc</kbd>: 返回远程仓库列表
  <kbd>g</kbd>: 查看重置选项
  <kbd>enter</kbd>: 查看提交
  <kbd>space</kbd>: 检出
  <kbd>n</kbd>: 新分支
  <kbd>M</kbd>: 合并到当前检出的分支
  <kbd>d</kbd>: 删除分支
  <kbd>r</kbd>: 将已检出的分支变基到该分支
  <kbd>u</kbd>: 设置为检出分支的上游
</pre>

## 分支 面板 (远程页面)

<pre>
  <kbd>f</kbd>: 抓取远程仓库
  <kbd>n</kbd>: 添加新的远程仓库
  <kbd>d</kbd>: 删除远程
  <kbd>e</kbd>: 编辑远程仓库
</pre>

## 分支 面板 (子提交)

<pre>
  <kbd>enter</kbd>: 查看提交的文件
  <kbd>space</kbd>: 检出提交
  <kbd>g</kbd>: 查看重置选项
  <kbd>n</kbd>: 新分支
  <kbd>c</kbd>: 复制提交（拣选）
  <kbd>C</kbd>: 复制提交范围（拣选）
  <kbd>ctrl+r</kbd>: 重置已拣选（复制）的提交
  <kbd>ctrl+o</kbd>: 将提交的 SHA 复制到剪贴板
</pre>

## 分支 面板 (标签页面)

<pre>
  <kbd>space</kbd>: 检出
  <kbd>d</kbd>: 删除标签
  <kbd>P</kbd>: 推送标签
  <kbd>n</kbd>: 创建标签
  <kbd>g</kbd>: 查看重置选项
  <kbd>enter</kbd>: 查看提交
</pre>

## 提交文件 面板

<pre>
  <kbd>ctrl+o</kbd>: 将提交的文件名复制到剪贴板
  <kbd>c</kbd>: 检出文件
  <kbd>d</kbd>: 放弃对此文件的提交更改
  <kbd>o</kbd>: 打开文件
  <kbd>e</kbd>: 编辑文件
  <kbd>space</kbd>: 补丁中包含的切换文件
  <kbd>enter</kbd>: 输入文件以将所选行添加到补丁中（或切换目录折叠）
  <kbd>`</kbd>: 切换文件树视图
</pre>

## 提交 面板 (提交)

<pre>
  <kbd>c</kbd>: 复制提交（拣选）
  <kbd>ctrl+o</kbd>: 将提交的 SHA 复制到剪贴板
  <kbd>C</kbd>: 复制提交范围（拣选）
  <kbd>v</kbd>: 粘贴提交（拣选）
  <kbd>n</kbd>: 从提交创建新分支
  <kbd>ctrl+r</kbd>: 重置已拣选（复制）的提交
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
  <kbd>A</kbd>: 用已暂存的更改来修补提交
  <kbd>t</kbd>: 还原提交
  <kbd>ctrl+l</kbd>: open log menu
  <kbd>g</kbd>: 重置为此提交
  <kbd>enter</kbd>: 查看提交的文件
  <kbd>space</kbd>: 检出提交
  <kbd>T</kbd>: 标签提交
  <kbd>ctrl+y</kbd>: 将提交消息复制到剪贴板
  <kbd>o</kbd>: open commit in browser
  <kbd>b</kbd>: view bisect options
</pre>

## 提交 面板 (Reflog)

<pre>
  <kbd>enter</kbd>: 查看提交的文件
  <kbd>space</kbd>: 检出提交
  <kbd>g</kbd>: 查看重置选项
  <kbd>c</kbd>: 复制提交（拣选）
  <kbd>C</kbd>: 复制提交范围（拣选）
  <kbd>ctrl+r</kbd>: 重置已拣选（复制）的提交
  <kbd>ctrl+o</kbd>: 将提交的 SHA 复制到剪贴板
</pre>

## Extras 面板

<pre>
  <kbd>@</kbd>: 打开命令日志菜单
</pre>

## 文件 面板

<pre>
  <kbd>ctrl+b</kbd>: 过滤提交文件
</pre>

## 文件 面板 (文件)

<pre>
  <kbd>c</kbd>: 提交更改
  <kbd>w</kbd>: 提交更改而无需预先提交钩子
  <kbd>A</kbd>: 修补最后一次提交
  <kbd>C</kbd>: 提交更改（使用编辑器编辑提交信息）
  <kbd>d</kbd>: 查看'放弃更改‘选项
  <kbd>e</kbd>: 编辑文件
  <kbd>o</kbd>: 打开文件
  <kbd>i</kbd>: 添加到 .gitignore
  <kbd>r</kbd>: 刷新文件
  <kbd>s</kbd>: 将所有更改加入贮藏
  <kbd>S</kbd>: 查看隐藏选项
  <kbd>a</kbd>: 切换所有文件的暂存状态
  <kbd>D</kbd>: 查看重置选项
  <kbd>enter</kbd>: 暂存单个 块/行 用于文件, 或 折叠/展开 目录
  <kbd>f</kbd>: 抓取
  <kbd>ctrl+o</kbd>: 将文件名复制到剪贴板
  <kbd>g</kbd>: 查看上游重置选项
  <kbd>`</kbd>: 切换文件树视图
  <kbd>M</kbd>: 打开合并工具
  <kbd>ctrl+w</kbd>: 切换是否在差异视图中显示空白更改
  <kbd>space</kbd>: 切换暂存状态
</pre>

## 文件 面板 (子模块)

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

## 主要 面板 (合并中)

<pre>
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
  <kbd>esc</kbd>: 返回文件面板
  <kbd>M</kbd>: 打开合并工具
  <kbd>space</kbd>: 选中区块
  <kbd>b</kbd>: 选中所有区块
  <kbd>◄</kbd>: 选择上一个冲突
  <kbd>►</kbd>: 选择下一个冲突
  <kbd>▲</kbd>: 选择顶部块
  <kbd>▼</kbd>: 选择底部块
  <kbd>z</kbd>: 撤销
</pre>

## 主要 面板 (正常)

<pre>
  <kbd>Ő</kbd>: 向下滚动 (fn+up)
  <kbd>ő</kbd>: 向上滚动 (fn+down)
</pre>

## 主要 面板 (构建补丁中)

<pre>
  <kbd>esc</kbd>: 退出逐行模式
  <kbd>o</kbd>: 打开文件
  <kbd>▲</kbd>: 选择上一行
  <kbd>▼</kbd>: 选择下一行
  <kbd>◄</kbd>: 选择上一个区块
  <kbd>►</kbd>: 选择下一个区块
  <kbd>ctrl+o</kbd>: copy the selected text to the clipboard
  <kbd>space</kbd>: 添加/移除 行到补丁
  <kbd>v</kbd>: 切换拖动选择
  <kbd>V</kbd>: 切换拖动选择
  <kbd>a</kbd>: 切换选择区块
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
</pre>

## 主要 面板 (正在暂存)

<pre>
  <kbd>esc</kbd>: 返回文件面板
  <kbd>space</kbd>: 切换行暂存状态
  <kbd>d</kbd>: 取消变更 (git reset)
  <kbd>tab</kbd>: 切换到其他面板
  <kbd>o</kbd>: 打开文件
  <kbd>▲</kbd>: 选择上一行
  <kbd>▼</kbd>: 选择下一行
  <kbd>◄</kbd>: 选择上一个区块
  <kbd>►</kbd>: 选择下一个区块
  <kbd>ctrl+o</kbd>: copy the selected text to the clipboard
  <kbd>e</kbd>: 编辑文件
  <kbd>o</kbd>: 打开文件
  <kbd>v</kbd>: 切换拖动选择
  <kbd>V</kbd>: 切换拖动选择
  <kbd>a</kbd>: 切换选择区块
  <kbd>H</kbd>: scroll left
  <kbd>L</kbd>: scroll right
  <kbd>c</kbd>: 提交更改
  <kbd>w</kbd>: 提交更改而无需预先提交钩子
  <kbd>C</kbd>: 提交更改（使用编辑器编辑提交信息）
</pre>

## 菜单 面板

<pre>
  <kbd>esc</kbd>: 关闭菜单
</pre>

## 贮藏 面板

<pre>
  <kbd>enter</kbd>: 查看贮藏条目中的文件
  <kbd>space</kbd>: 应用
  <kbd>g</kbd>: 应用并删除
  <kbd>d</kbd>: 删除
  <kbd>n</kbd>: 新分支
</pre>

## 状态 面板

<pre>
  <kbd>e</kbd>: 编辑配置文件
  <kbd>o</kbd>: 打开配置文件
  <kbd>u</kbd>: 检查更新
  <kbd>enter</kbd>: 切换到最近的仓库
  <kbd>a</kbd>: 显示所有分支的日志
</pre>
