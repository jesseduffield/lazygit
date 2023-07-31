_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit 鍵盤快捷鍵

_說明：`<c-b>` 表示 Ctrl+B、`<a-b>` 表示 Alt+B，`B`表示 Shift+B_

## 全局快捷鍵

<pre>
  <kbd>&lt;c-r&gt;</kbd>: 切換到最近使用的版本庫
  <kbd>&lt;pgup&gt;</kbd>: 向上捲動主面板 (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: 向下捲動主面板 (fn+down/shift+j)
  <kbd>@</kbd>: 開啟命令記錄選單
  <kbd>}</kbd>: 增加差異檢視中顯示變更周圍上下文的大小
  <kbd>{</kbd>: 減小差異檢視中顯示變更周圍上下文的大小
  <kbd>:</kbd>: 執行自訂命令
  <kbd>&lt;c-p&gt;</kbd>: 檢視自訂補丁選項
  <kbd>m</kbd>: 查看合併/變基選項
  <kbd>R</kbd>: 重新整理
  <kbd>+</kbd>: 下一個螢幕模式（常規/半螢幕/全螢幕）
  <kbd>_</kbd>: 上一個螢幕模式
  <kbd>?</kbd>: 開啟選單
  <kbd>&lt;c-s&gt;</kbd>: 檢視篩選路徑選項
  <kbd>W</kbd>: 開啟差異比較選單
  <kbd>&lt;c-e&gt;</kbd>: 開啟差異比較選單
  <kbd>&lt;c-w&gt;</kbd>: 切換是否在差異檢視中顯示空格變更
  <kbd>z</kbd>: 復原
  <kbd>&lt;c-z&gt;</kbd>: 取消復原
  <kbd>P</kbd>: 推送
  <kbd>p</kbd>: 拉取
</pre>

## 列表面板導航

<pre>
  <kbd>,</kbd>: 上一頁
  <kbd>.</kbd>: 下一頁
  <kbd>&lt;</kbd>: 捲動到頂部
  <kbd>&gt;</kbd>: 捲動到底部
  <kbd>/</kbd>: 開始搜尋
  <kbd>H</kbd>: 向左捲動
  <kbd>L</kbd>: 向右捲動
  <kbd>]</kbd>: 下一個索引標籤
  <kbd>[</kbd>: 上一個索引標籤
</pre>

## Reflog

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 複製提交 SHA 到剪貼簿
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: 檢出提交
  <kbd>y</kbd>: 複製提交屬性
  <kbd>o</kbd>: 在瀏覽器中開啟提交
  <kbd>n</kbd>: 從提交建立新分支
  <kbd>g</kbd>: 檢視重設選項
  <kbd>c</kbd>: 複製提交 (揀選)
  <kbd>C</kbd>: 複製提交範圍 (揀選)
  <kbd>&lt;c-r&gt;</kbd>: 重設選定的揀選 (複製) 提交
  <kbd>&lt;enter&gt;</kbd>: 檢視提交
  <kbd>/</kbd>: Filter the current view by text
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

## 主視窗 (一般)

<pre>
  <kbd>mouse wheel down</kbd>: 向下捲動 (fn+up)
  <kbd>mouse wheel up</kbd>: 向上捲動 (fn+down)
</pre>

## 主視窗 (合併中)

<pre>
  <kbd>e</kbd>: 編輯檔案
  <kbd>o</kbd>: 開啟檔案
  <kbd>&lt;left&gt;</kbd>: 選擇上一個衝突
  <kbd>&lt;right&gt;</kbd>: 選擇下一個衝突
  <kbd>&lt;up&gt;</kbd>: 選擇上一段
  <kbd>&lt;down&gt;</kbd>: 選擇下一段
  <kbd>z</kbd>: 復原
  <kbd>M</kbd>: 開啟外部合併工具 (git mergetool)
  <kbd>&lt;space&gt;</kbd>: 挑選程式碼片段
  <kbd>b</kbd>: 挑選所有程式碼片段
  <kbd>&lt;esc&gt;</kbd>: 返回檔案面板
</pre>

## 主視窗 (預存中)

<pre>
  <kbd>&lt;left&gt;</kbd>: 選擇上一段
  <kbd>&lt;right&gt;</kbd>: 選擇下一段
  <kbd>v</kbd>: 切換拖曳選擇
  <kbd>V</kbd>: 切換拖曳選擇
  <kbd>a</kbd>: 切換選擇程式碼塊
  <kbd>&lt;c-o&gt;</kbd>: 複製所選文本至剪貼簿
  <kbd>o</kbd>: 開啟檔案
  <kbd>e</kbd>: 編輯檔案
  <kbd>&lt;esc&gt;</kbd>: 返回檔案面板
  <kbd>&lt;tab&gt;</kbd>: 切換至另一個面板 (已預存/未預存更改)
  <kbd>&lt;space&gt;</kbd>: 切換現有行的狀態 (已預存/未預存)
  <kbd>d</kbd>: 刪除變更 (git reset)
  <kbd>E</kbd>: 編輯程式碼塊
  <kbd>c</kbd>: 提交變更
  <kbd>w</kbd>: 沒有預提交 hook 就提交更改
  <kbd>C</kbd>: 使用 git 編輯器提交變更
  <kbd>/</kbd>: 開始搜尋
</pre>

## 主面板 (補丁生成)

<pre>
  <kbd>&lt;left&gt;</kbd>: 選擇上一段
  <kbd>&lt;right&gt;</kbd>: 選擇下一段
  <kbd>v</kbd>: 切換拖曳選擇
  <kbd>V</kbd>: 切換拖曳選擇
  <kbd>a</kbd>: 切換選擇程式碼塊
  <kbd>&lt;c-o&gt;</kbd>: 複製所選文本至剪貼簿
  <kbd>o</kbd>: 開啟檔案
  <kbd>e</kbd>: 編輯檔案
  <kbd>&lt;space&gt;</kbd>: 向 (或從) 補丁中添加/刪除行
  <kbd>&lt;esc&gt;</kbd>: 退出自訂補丁建立器
  <kbd>/</kbd>: 開始搜尋
</pre>

## 功能表

<pre>
  <kbd>&lt;enter&gt;</kbd>: 執行
  <kbd>&lt;esc&gt;</kbd>: 關閉
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 子提交

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 複製提交 SHA 到剪貼簿
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: 檢出提交
  <kbd>y</kbd>: 複製提交屬性
  <kbd>o</kbd>: 在瀏覽器中開啟提交
  <kbd>n</kbd>: 從提交建立新分支
  <kbd>g</kbd>: 檢視重設選項
  <kbd>c</kbd>: 複製提交 (揀選)
  <kbd>C</kbd>: 複製提交範圍 (揀選)
  <kbd>&lt;c-r&gt;</kbd>: 重設選定的揀選 (複製) 提交
  <kbd>&lt;enter&gt;</kbd>: 檢視所選項目的檔案
  <kbd>/</kbd>: 開始搜尋
</pre>

## 子模組

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 複製子模組名稱到剪貼簿
  <kbd>&lt;enter&gt;</kbd>: 進入子模組
  <kbd>&lt;space&gt;</kbd>: 進入子模組
  <kbd>d</kbd>: 移除子模組
  <kbd>u</kbd>: 更新子模組
  <kbd>n</kbd>: 新增子模組
  <kbd>e</kbd>: 更新子模組 URL
  <kbd>i</kbd>: 初始化子模組
  <kbd>b</kbd>: 查看批量子模組選項
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 提交

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 複製提交 SHA 到剪貼簿
  <kbd>&lt;c-r&gt;</kbd>: 重設選定的揀選 (複製) 提交
  <kbd>b</kbd>: 查看二分選項
  <kbd>s</kbd>: 向下壓縮
  <kbd>f</kbd>: 修復提交 (Fixup)
  <kbd>r</kbd>: 改寫提交
  <kbd>R</kbd>: 使用編輯器改寫提交
  <kbd>d</kbd>: 刪除提交
  <kbd>e</kbd>: 編輯提交
  <kbd>p</kbd>: 挑選提交 (於變基過程中)
  <kbd>F</kbd>: 為此提交建立修復提交
  <kbd>S</kbd>: 壓縮上方所有的“fixup!”提交 (自動壓縮)
  <kbd>&lt;c-j&gt;</kbd>: 向下移動提交
  <kbd>&lt;c-k&gt;</kbd>: 向上移動提交
  <kbd>v</kbd>: 貼上提交 (揀選)
  <kbd>B</kbd>: Mark commit as base commit for rebase
  <kbd>A</kbd>: 使用已預存的更改修正提交
  <kbd>a</kbd>: 設置/重設提交作者
  <kbd>t</kbd>: 還原提交
  <kbd>T</kbd>: 打標籤到提交
  <kbd>&lt;c-l&gt;</kbd>: 開啟記錄選單
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: 檢出提交
  <kbd>y</kbd>: 複製提交屬性
  <kbd>o</kbd>: 在瀏覽器中開啟提交
  <kbd>n</kbd>: 從提交建立新分支
  <kbd>g</kbd>: 檢視重設選項
  <kbd>c</kbd>: 複製提交 (揀選)
  <kbd>C</kbd>: 複製提交範圍 (揀選)
  <kbd>&lt;enter&gt;</kbd>: 檢視所選項目的檔案
  <kbd>/</kbd>: 開始搜尋
</pre>

## 提交摘要

<pre>
  <kbd>&lt;enter&gt;</kbd>: 確認
  <kbd>&lt;esc&gt;</kbd>: 關閉
</pre>

## 提交檔案

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 複製提交的檔案名稱到剪貼簿
  <kbd>c</kbd>: 檢出檔案
  <kbd>d</kbd>: 捨棄此提交對此檔案的更改
  <kbd>o</kbd>: 開啟檔案
  <kbd>e</kbd>: 編輯檔案
  <kbd>&lt;space&gt;</kbd>: 切換檔案是否包含在補丁中
  <kbd>a</kbd>: 切換所有檔案是否包含在補丁中
  <kbd>&lt;enter&gt;</kbd>: 輸入檔案以將選定的行添加至補丁（或切換目錄折疊）
  <kbd>`</kbd>: 切換檔案樹狀視圖
  <kbd>/</kbd>: 開始搜尋
</pre>

## 收藏 (Stash)

<pre>
  <kbd>&lt;space&gt;</kbd>: 套用
  <kbd>g</kbd>: 還原
  <kbd>d</kbd>: 捨棄
  <kbd>n</kbd>: 新分支
  <kbd>r</kbd>: 重新命名收藏
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: 檢視所選項目的檔案
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 本地分支

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 複製分支名稱到剪貼簿
  <kbd>i</kbd>: 顯示 git-flow 選項
  <kbd>&lt;space&gt;</kbd>: 檢出
  <kbd>n</kbd>: 新分支
  <kbd>o</kbd>: 建立拉取請求
  <kbd>O</kbd>: 建立拉取請求選項
  <kbd>&lt;c-y&gt;</kbd>: 複製拉取請求的 URL 到剪貼板
  <kbd>c</kbd>: 根據名稱檢出
  <kbd>F</kbd>: 強制檢出
  <kbd>d</kbd>: 刪除分支
  <kbd>r</kbd>: 將已檢出的分支變基至此分支
  <kbd>M</kbd>: 合併到當前檢出的分支
  <kbd>f</kbd>: 從上游快進此分支
  <kbd>T</kbd>: 建立標籤
  <kbd>g</kbd>: 檢視重設選項
  <kbd>R</kbd>: 重新命名分支
  <kbd>u</kbd>: 設定/取消設定上游
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: 檢視提交
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 標籤

<pre>
  <kbd>&lt;space&gt;</kbd>: 檢出
  <kbd>d</kbd>: 刪除標籤
  <kbd>P</kbd>: 推送標籤
  <kbd>n</kbd>: 建立標籤
  <kbd>g</kbd>: 檢視重設選項
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: 檢視提交
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 檔案

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 複製檔案名稱到剪貼簿
  <kbd>d</kbd>: 檢視“捨棄更改”的選項
  <kbd>&lt;space&gt;</kbd>: 切換預存
  <kbd>&lt;c-b&gt;</kbd>: 篩選檔案 (預存/未預存)
  <kbd>c</kbd>: 提交變更
  <kbd>w</kbd>: 沒有預提交 hook 就提交更改
  <kbd>A</kbd>: 修正上次提交
  <kbd>C</kbd>: 使用 git 編輯器提交變更
  <kbd>e</kbd>: 編輯檔案
  <kbd>o</kbd>: 開啟檔案
  <kbd>i</kbd>: 忽略或排除檔案
  <kbd>r</kbd>: 重新整理檔案
  <kbd>s</kbd>: 收藏所有變更
  <kbd>S</kbd>: 檢視收藏選項
  <kbd>a</kbd>: 全部預存/取消預存
  <kbd>&lt;enter&gt;</kbd>: 選擇檔案中的單個程式碼塊/行，或展開/折疊目錄
  <kbd>g</kbd>: 檢視上游重設選項
  <kbd>D</kbd>: 檢視重設選項
  <kbd>`</kbd>: 切換檔案樹狀視圖
  <kbd>M</kbd>: 開啟外部合併工具 (git mergetool)
  <kbd>f</kbd>: 擷取
  <kbd>/</kbd>: 開始搜尋
</pre>

## 狀態

<pre>
  <kbd>o</kbd>: 開啟設定檔案
  <kbd>e</kbd>: 編輯設定檔案
  <kbd>u</kbd>: 檢查更新
  <kbd>&lt;enter&gt;</kbd>: 切換到最近使用的版本庫
  <kbd>a</kbd>: 顯示所有分支日誌
</pre>

## 確認面板

<pre>
  <kbd>&lt;enter&gt;</kbd>: 確認
  <kbd>&lt;esc&gt;</kbd>: 關閉/取消
</pre>

## 遠端

<pre>
  <kbd>f</kbd>: 擷取遠端
  <kbd>n</kbd>: 新增遠端
  <kbd>d</kbd>: 移除遠端
  <kbd>e</kbd>: 編輯遠端
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 遠端分支

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 複製分支名稱到剪貼簿
  <kbd>&lt;space&gt;</kbd>: 檢出
  <kbd>n</kbd>: 新分支
  <kbd>M</kbd>: 合併到當前檢出的分支
  <kbd>r</kbd>: 將已檢出的分支變基至此分支
  <kbd>d</kbd>: 刪除分支
  <kbd>u</kbd>: 將此分支設為當前分支之上游
  <kbd>g</kbd>: 檢視重設選項
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: 檢視提交
  <kbd>/</kbd>: Filter the current view by text
</pre>
