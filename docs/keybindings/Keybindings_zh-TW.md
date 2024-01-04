_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit 鍵盤快捷鍵

_說明：`<c-b>` 表示 Ctrl+B、`<a-b>` 表示 Alt+B，`B`表示 Shift+B_

## 全局快捷鍵

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | 切換到最近使用的版本庫 |  |
| `` <pgup> (fn+up/shift+k) `` | 向上捲動主面板 |  |
| `` <pgdown> (fn+down/shift+j) `` | 向下捲動主面板 |  |
| `` @ `` | 開啟命令記錄選單 |  |
| `` } `` | 增加差異檢視中顯示變更周圍上下文的大小 |  |
| `` { `` | 減小差異檢視中顯示變更周圍上下文的大小 |  |
| `` : `` | 執行自訂命令 |  |
| `` <c-p> `` | 檢視自訂補丁選項 |  |
| `` m `` | 查看合併/變基選項 |  |
| `` R `` | 重新整理 |  |
| `` + `` | 下一個螢幕模式（常規/半螢幕/全螢幕） |  |
| `` _ `` | 上一個螢幕模式 |  |
| `` ? `` | 開啟選單 |  |
| `` <c-s> `` | 檢視篩選路徑選項 |  |
| `` W `` | 開啟差異比較選單 |  |
| `` <c-e> `` | 開啟差異比較選單 |  |
| `` <c-w> `` | 切換是否在差異檢視中顯示空格變更 |  |
| `` z `` | 復原 | 將使用 reflog 確定要運行哪個 git 命令以復原上一個 git 命令。這不包括工作區的更改；只考慮提交。 |
| `` <c-z> `` | 取消復原 | 將使用 reflog 確定要運行哪個 git 命令以重作上一個 git 命令。這不包括工作區的更改；只考慮提交。 |
| `` P `` | 推送 |  |
| `` p `` | 拉取 |  |

## 列表面板導航

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | 上一頁 |  |
| `` . `` | 下一頁 |  |
| `` < `` | 捲動到頂部 |  |
| `` > `` | 捲動到底部 |  |
| `` v `` | 切換拖曳選擇 |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | 開始搜尋 |  |
| `` H `` | 向左捲動 |  |
| `` L `` | 向右捲動 |  |
| `` ] `` | 下一個索引標籤 |  |
| `` [ `` | 上一個索引標籤 |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製提交 SHA 到剪貼簿 |  |
| `` w `` | View worktree options |  |
| `` <space> `` | 檢出提交 |  |
| `` y `` | 複製提交屬性 |  |
| `` o `` | 在瀏覽器中開啟提交 |  |
| `` n `` | 從提交建立新分支 |  |
| `` g `` | 檢視重設選項 |  |
| `` C `` | 複製提交 (揀選) |  |
| `` <c-r> `` | 重設選定的揀選 (複製) 提交 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 檢視提交 |  |
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

## 主視窗 (一般)

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | 向下捲動 |  |
| `` mouse wheel up (fn+down) `` | 向上捲動 |  |

## 主視窗 (合併中)

| Key | Action | Info |
|-----|--------|-------------|
| `` e `` | 編輯檔案 |  |
| `` o `` | 開啟檔案 |  |
| `` <left> `` | 選擇上一個衝突 |  |
| `` <right> `` | 選擇下一個衝突 |  |
| `` <up> `` | 選擇上一段 |  |
| `` <down> `` | 選擇下一段 |  |
| `` z `` | 復原 |  |
| `` M `` | 開啟外部合併工具 (git mergetool) |  |
| `` <space> `` | 挑選程式碼片段 |  |
| `` b `` | 挑選所有程式碼片段 |  |
| `` <esc> `` | 返回檔案面板 |  |

## 主視窗 (預存中)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 選擇上一段 |  |
| `` <right> `` | 選擇下一段 |  |
| `` v `` | 切換拖曳選擇 |  |
| `` a `` | 切換選擇程式碼塊 |  |
| `` <c-o> `` | 複製所選文本至剪貼簿 |  |
| `` o `` | 開啟檔案 |  |
| `` e `` | 編輯檔案 |  |
| `` <esc> `` | 返回檔案面板 |  |
| `` <tab> `` | 切換至另一個面板 (已預存/未預存更改) |  |
| `` <space> `` | 切換現有行的狀態 (已預存/未預存) |  |
| `` d `` | 刪除變更 (git reset) |  |
| `` E `` | 編輯程式碼塊 |  |
| `` c `` | 提交變更 |  |
| `` w `` | 沒有預提交 hook 就提交更改 |  |
| `` C `` | 使用 git 編輯器提交變更 |  |
| `` / `` | 開始搜尋 |  |

## 主面板 (補丁生成)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 選擇上一段 |  |
| `` <right> `` | 選擇下一段 |  |
| `` v `` | 切換拖曳選擇 |  |
| `` a `` | 切換選擇程式碼塊 |  |
| `` <c-o> `` | 複製所選文本至剪貼簿 |  |
| `` o `` | 開啟檔案 |  |
| `` e `` | 編輯檔案 |  |
| `` <space> `` | 向 (或從) 補丁中添加/刪除行 |  |
| `` <esc> `` | 退出自訂補丁建立器 |  |
| `` / `` | 開始搜尋 |  |

## 功能表

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 執行 |  |
| `` <esc> `` | 關閉 |  |
| `` / `` | Filter the current view by text |  |

## 子提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製提交 SHA 到剪貼簿 |  |
| `` w `` | View worktree options |  |
| `` <space> `` | 檢出提交 |  |
| `` y `` | 複製提交屬性 |  |
| `` o `` | 在瀏覽器中開啟提交 |  |
| `` n `` | 從提交建立新分支 |  |
| `` g `` | 檢視重設選項 |  |
| `` C `` | 複製提交 (揀選) |  |
| `` <c-r> `` | 重設選定的揀選 (複製) 提交 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 檢視所選項目的檔案 |  |
| `` / `` | 開始搜尋 |  |

## 子模組

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製子模組名稱到剪貼簿 |  |
| `` <enter> `` | 進入子模組 |  |
| `` <space> `` | 進入子模組 |  |
| `` d `` | 移除子模組 |  |
| `` u `` | 更新子模組 |  |
| `` n `` | 新增子模組 |  |
| `` e `` | 更新子模組 URL |  |
| `` i `` | 初始化子模組 |  |
| `` b `` | 查看批量子模組選項 |  |
| `` / `` | Filter the current view by text |  |

## 提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製提交 SHA 到剪貼簿 |  |
| `` <c-r> `` | 重設選定的揀選 (複製) 提交 |  |
| `` b `` | 查看二分選項 |  |
| `` s `` | 向下壓縮 |  |
| `` f `` | 修復提交 (Fixup) |  |
| `` r `` | 改寫提交 |  |
| `` R `` | 使用編輯器改寫提交 |  |
| `` d `` | 刪除提交 |  |
| `` e `` | 編輯提交 |  |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | 挑選提交 (於變基過程中) |  |
| `` F `` | 為此提交建立修復提交 |  |
| `` S `` | 壓縮上方所有的“fixup!”提交 (自動壓縮) |  |
| `` <c-j> `` | 向下移動提交 |  |
| `` <c-k> `` | 向上移動提交 |  |
| `` V `` | 貼上提交 (揀選) |  |
| `` B `` | Mark commit as base commit for rebase | Select a base commit for the next rebase; this will effectively perform a 'git rebase --onto'. |
| `` A `` | 使用已預存的更改修正提交 |  |
| `` a `` | 設置/重設提交作者 |  |
| `` t `` | 還原提交 |  |
| `` T `` | 打標籤到提交 |  |
| `` <c-l> `` | 開啟記錄選單 |  |
| `` w `` | View worktree options |  |
| `` <space> `` | 檢出提交 |  |
| `` y `` | 複製提交屬性 |  |
| `` o `` | 在瀏覽器中開啟提交 |  |
| `` n `` | 從提交建立新分支 |  |
| `` g `` | 檢視重設選項 |  |
| `` C `` | 複製提交 (揀選) |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 檢視所選項目的檔案 |  |
| `` / `` | 開始搜尋 |  |

## 提交摘要

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 確認 |  |
| `` <esc> `` | 關閉 |  |

## 提交檔案

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製提交的檔案名稱到剪貼簿 |  |
| `` c `` | 檢出檔案 |  |
| `` d `` | 捨棄此提交對此檔案的更改 |  |
| `` o `` | 開啟檔案 |  |
| `` e `` | 編輯檔案 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | 切換檔案是否包含在補丁中 |  |
| `` a `` | 切換所有檔案是否包含在補丁中 |  |
| `` <enter> `` | 輸入檔案以將選定的行添加至補丁（或切換目錄折疊） |  |
| `` ` `` | 切換檔案樹狀視圖 |  |
| `` / `` | 開始搜尋 |  |

## 收藏 (Stash)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 套用 |  |
| `` g `` | 還原 |  |
| `` d `` | 捨棄 |  |
| `` n `` | 新分支 |  |
| `` r `` | 重新命名收藏 |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | 檢視所選項目的檔案 |  |
| `` / `` | Filter the current view by text |  |

## 本地分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製分支名稱到剪貼簿 |  |
| `` i `` | 顯示 git-flow 選項 |  |
| `` <space> `` | 檢出 |  |
| `` n `` | 新分支 |  |
| `` o `` | 建立拉取請求 |  |
| `` O `` | 建立拉取請求選項 |  |
| `` <c-y> `` | 複製拉取請求的 URL 到剪貼板 |  |
| `` c `` | 根據名稱檢出 |  |
| `` F `` | 強制檢出 |  |
| `` d `` | View delete options |  |
| `` r `` | 將已檢出的分支變基至此分支 |  |
| `` M `` | 合併到當前檢出的分支 |  |
| `` f `` | 從上游快進此分支 |  |
| `` T `` | 建立標籤 |  |
| `` s `` | Sort order |  |
| `` g `` | 檢視重設選項 |  |
| `` R `` | 重新命名分支 |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream |
| `` w `` | View worktree options |  |
| `` <enter> `` | 檢視提交 |  |
| `` / `` | Filter the current view by text |  |

## 標籤

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 檢出 |  |
| `` d `` | View delete options |  |
| `` P `` | 推送標籤 |  |
| `` n `` | 建立標籤 |  |
| `` g `` | 檢視重設選項 |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | 檢視提交 |  |
| `` / `` | Filter the current view by text |  |

## 檔案

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製檔案名稱到剪貼簿 |  |
| `` <space> `` | 切換預存 |  |
| `` <c-b> `` | 篩選檔案 (預存/未預存) |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | 提交變更 |  |
| `` w `` | 沒有預提交 hook 就提交更改 |  |
| `` A `` | 修正上次提交 |  |
| `` C `` | 使用 git 編輯器提交變更 |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | 編輯檔案 |  |
| `` o `` | 開啟檔案 |  |
| `` i `` | 忽略或排除檔案 |  |
| `` r `` | 重新整理檔案 |  |
| `` s `` | 收藏所有變更 |  |
| `` S `` | 檢視收藏選項 |  |
| `` a `` | 全部預存/取消預存 |  |
| `` <enter> `` | 選擇檔案中的單個程式碼塊/行，或展開/折疊目錄 |  |
| `` d `` | 檢視“捨棄更改”的選項 |  |
| `` g `` | 檢視上游重設選項 |  |
| `` D `` | 檢視重設選項 |  |
| `` ` `` | 切換檔案樹狀視圖 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | 開啟外部合併工具 (git mergetool) |  |
| `` f `` | 擷取 |  |
| `` / `` | 開始搜尋 |  |

## 狀態

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 開啟設定檔案 |  |
| `` e `` | 編輯設定檔案 |  |
| `` u `` | 檢查更新 |  |
| `` <enter> `` | 切換到最近使用的版本庫 |  |
| `` a `` | 顯示所有分支日誌 |  |

## 確認面板

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 確認 |  |
| `` <esc> `` | 關閉/取消 |  |

## 遠端

| Key | Action | Info |
|-----|--------|-------------|
| `` f `` | 擷取遠端 |  |
| `` n `` | 新增遠端 |  |
| `` d `` | 移除遠端 |  |
| `` e `` | 編輯遠端 |  |
| `` / `` | Filter the current view by text |  |

## 遠端分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製分支名稱到剪貼簿 |  |
| `` <space> `` | 檢出 |  |
| `` n `` | 新分支 |  |
| `` M `` | 合併到當前檢出的分支 |  |
| `` r `` | 將已檢出的分支變基至此分支 |  |
| `` d `` | Delete remote tag |  |
| `` u `` | 將此分支設為當前分支之上游 |  |
| `` s `` | Sort order |  |
| `` g `` | 檢視重設選項 |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | 檢視提交 |  |
| `` / `` | Filter the current view by text |  |
