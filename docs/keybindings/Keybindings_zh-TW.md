_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit 鍵盤快捷鍵

_說明：`<c-b>` 表示 Ctrl＋B、`<a-b>` 表示 Alt＋B，`B`表示 Shift＋B_

## 全域快捷鍵

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | 切換到最近使用的版本庫 |  |
| `` <pgup> (fn+up/shift+k) `` | 向上捲動主面板 |  |
| `` <pgdown> (fn+down/shift+j) `` | 向下捲動主面板 |  |
| `` @ `` | 開啟命令記錄選單 | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | 推送 | 推送到遠端。如果沒有設定遠端，會開啟設定視窗。 |
| `` p `` | 拉取 | 從遠端同步當前分支。如果沒有設定遠端，會開啟設定視窗。 |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` } `` | 增加差異檢視中顯示變更周圍上下文的大小 | Increase the amount of the context shown around changes in the diff view. |
| `` { `` | 減小差異檢視中顯示變更周圍上下文的大小 | Decrease the amount of the context shown around changes in the diff view. |
| `` : `` | Execute shell command | Bring up a prompt where you can enter a shell command to execute. |
| `` <c-p> `` | 檢視自訂補丁選項 |  |
| `` m `` | 查看合併/變基選項 | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | 重新整理 | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | 下一個螢幕模式（常規/半螢幕/全螢幕） |  |
| `` _ `` | 上一個螢幕模式 |  |
| `` ? `` | 開啟選單 |  |
| `` <c-s> `` | 檢視篩選路徑選項 | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | 開啟差異比較選單 | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | 開啟差異比較選單 | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | 結束 |  |
| `` <esc> `` | 取消 |  |
| `` <c-w> `` | 切換是否在差異檢視中顯示空格變更 | Toggle whether or not whitespace changes are shown in the diff view. |
| `` z `` | 復原 | 將使用 reflog 確任 git 指令以復原。這不包括工作區更改；只考慮提交。 |
| `` <c-z> `` | 取消復原 | 將使用 reflog 確任 git 指令以重作。這不包括工作區更改；只考慮提交。 |

## 移動

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | 上一頁 |  |
| `` . `` | 下一頁 |  |
| `` < `` | 捲動到頂部 |  |
| `` > `` | 捲動到底部 |  |
| `` v `` | 切換拖曳選擇 |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | 搜尋 |  |
| `` H `` | 向左捲動 |  |
| `` L `` | 向右捲動 |  |
| `` ] `` | 下一個索引標籤 |  |
| `` [ `` | 上一個索引標籤 |  |

## 主面板 (補丁生成)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 選擇上一段 |  |
| `` <right> `` | 選擇下一段 |  |
| `` v `` | 切換拖曳選擇 |  |
| `` a `` | 切換選擇程式碼塊 | Toggle hunk selection mode. |
| `` <c-o> `` | 複製所選文本至剪貼簿 |  |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` e `` | 編輯檔案 | 使用外部編輯器開啟 |
| `` <space> `` | 向 (或從) 補丁中添加/刪除行 |  |
| `` <esc> `` | 退出自訂補丁建立器 |  |
| `` / `` | 搜尋 |  |

## 主面板（一般）

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | 向下捲動 |  |
| `` mouse wheel up (fn+down) `` | 向上捲動 |  |

## 主面板（合併）

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 挑選程式碼片段 |  |
| `` b `` | 挑選所有程式碼片段 |  |
| `` <up> `` | 選擇上一段 |  |
| `` <down> `` | 選擇下一段 |  |
| `` <left> `` | 選擇上一個衝突 |  |
| `` <right> `` | 選擇下一個衝突 |  |
| `` z `` | 復原 | Undo last merge conflict resolution. |
| `` e `` | 編輯檔案 | 使用外部編輯器開啟 |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` M `` | 開啟外部合併工具 | 執行 `git mergetool`。 |
| `` <esc> `` | 返回檔案面板 |  |

## 主面板（預存）

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 選擇上一段 |  |
| `` <right> `` | 選擇下一段 |  |
| `` v `` | 切換拖曳選擇 |  |
| `` a `` | 切換選擇程式碼塊 | Toggle hunk selection mode. |
| `` <c-o> `` | 複製所選文本至剪貼簿 |  |
| `` <space> `` | 切換預存 | 切換現有行的狀態 (已預存/未預存) |
| `` d `` | 刪除變更 (git reset) | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` e `` | 編輯檔案 | 使用外部編輯器開啟 |
| `` <esc> `` | 返回檔案面板 |  |
| `` <tab> `` | 切換至另一個面板 (已預存/未預存更改) | Switch to other view (staged/unstaged changes). |
| `` E `` | 編輯程式碼塊 | Edit selected hunk in external editor. |
| `` c `` | 提交變更 | 提交暫存區變更 |
| `` w `` | 沒有預提交 hook 就提交更改 |  |
| `` C `` | 使用 git 編輯器提交變更 |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | 搜尋 |  |

## 功能表

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 執行 |  |
| `` <esc> `` | 關閉 |  |
| `` / `` | 搜尋 |  |

## 子提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製提交 hash 到剪貼簿 |  |
| `` <space> `` | 檢出 | Checkout the selected commit as a detached HEAD. |
| `` y `` | 複製提交屬性 | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 在瀏覽器中開啟提交 |  |
| `` n `` | 從提交建立新分支 |  |
| `` g `` | 檢視重設選項 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | 複製提交 (揀選) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | 重設選定的揀選 (複製) 提交 |  |
| `` <c-t> `` | 開啟外部差異工具 (git difftool) |  |
| `` <enter> `` | 檢視所選項目的檔案 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 子模組

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製子模組名稱到剪貼簿 |  |
| `` <enter> `` | Enter | 進入子模組 |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | 更新子模組 |
| `` n `` | 新增子模組 |  |
| `` e `` | 更新子模組 URL |  |
| `` i `` | Initialize | 初始化子模組 |
| `` b `` | 查看批量子模組選項 |  |
| `` / `` | 搜尋 |  |

## 工作目錄

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | 在編輯器中開啟 |  |
| `` d `` | Remove | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | 搜尋 |  |

## 提交

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製提交 hash 到剪貼簿 |  |
| `` <c-r> `` | 重設選定的揀選 (複製) 提交 |  |
| `` b `` | 查看二分選項 |  |
| `` s `` | 壓縮 (Squash) | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | 修復 (Fixup) | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` r `` | 改寫提交 | 改寫選中的提交訊息 |
| `` R `` | 使用編輯器改寫提交 |  |
| `` d `` | 刪除提交 | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | 編輯(開始互動變基) | 編輯提交 |
| `` i `` | 開始互動變基 | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | 挑選 | 挑選提交 (於變基過程中) |
| `` F `` | 建立修復提交 | 為此提交建立修復提交 |
| `` S `` | 壓縮上方所有「fixup」提交（自動壓縮） | 是否壓縮上方 {{.commit}} 所有「fixup」提交？ |
| `` <c-j> `` | 向下移動提交 |  |
| `` <c-k> `` | 向上移動提交 |  |
| `` V `` | 貼上提交 (揀選) |  |
| `` B `` | 為了變基已標注提交為基準提交 | 請為了下一次變基選擇一項基準提交；此將執行 `git rebase --onto`。 |
| `` A `` | 修改 | 使用已預存的更改修正提交 |
| `` a `` | 設定/重設提交作者 | Set/Reset commit author or set co-author. |
| `` t `` | 還原 | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | 打標籤到提交 | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <c-l> `` | 開啟記錄選單 | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` <space> `` | 檢出 | Checkout the selected commit as a detached HEAD. |
| `` y `` | 複製提交屬性 | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 在瀏覽器中開啟提交 |  |
| `` n `` | 從提交建立新分支 |  |
| `` g `` | 檢視重設選項 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | 複製提交 (揀選) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-t> `` | 開啟外部差異工具 (git difftool) |  |
| `` <enter> `` | 檢視所選項目的檔案 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 提交摘要

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 確認 |  |
| `` <esc> `` | 關閉 |  |

## 提交檔案

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製檔案名稱到剪貼簿 |  |
| `` c `` | 檢出 | 檢出檔案 |
| `` d `` | Remove | Discard this commit's changes to this file. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes this file. |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` e `` | 編輯 | 使用外部編輯器開啟 |
| `` <c-t> `` | 開啟外部差異工具 (git difftool) |  |
| `` <space> `` | 切換檔案是否包含在補丁中 | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | 切換所有檔案是否包含在補丁中 | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | 輸入檔案以將選定的行添加至補丁（或切換目錄折疊） | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | 顯示檔案樹狀視圖 | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` / `` | 搜尋 |  |

## 收藏 (Stash)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 套用 | Apply the stash entry to your working directory. |
| `` g `` | 還原 | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | 捨棄 | Remove the stash entry from the stash list. |
| `` n `` | 新分支 | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | 重新命名收藏 |  |
| `` <enter> `` | 檢視所選項目的檔案 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 日誌

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製提交 hash 到剪貼簿 |  |
| `` <space> `` | 檢出 | Checkout the selected commit as a detached HEAD. |
| `` y `` | 複製提交屬性 | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 在瀏覽器中開啟提交 |  |
| `` n `` | 從提交建立新分支 |  |
| `` g `` | 檢視重設選項 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | 複製提交 (揀選) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | 重設選定的揀選 (複製) 提交 |  |
| `` <c-t> `` | 開啟外部差異工具 (git difftool) |  |
| `` <enter> `` | 檢視提交 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 本地分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製分支名稱到剪貼簿 |  |
| `` i `` | 顯示 git-flow 選項 |  |
| `` <space> `` | 檢出 | 檢出選定的項目。 |
| `` n `` | 新分支 |  |
| `` o `` | 建立拉取請求 |  |
| `` O `` | 建立拉取請求選項 |  |
| `` <c-y> `` | 複製拉取請求的 URL 到剪貼板 |  |
| `` c `` | 根據名稱檢出 | Checkout by name. In the input box you can enter '-' to switch to the last branch. |
| `` F `` | 強制檢出 | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | 刪除 | View delete options for local/remote branch. |
| `` r `` | 將已檢出的分支變基至此分支 | Rebase the checked-out branch onto the selected branch. |
| `` M `` | 合併到當前檢出的分支 | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` f `` | 從上游快進此分支 | 從遠端快進所選的分支 |
| `` T `` | 建立標籤 |  |
| `` s `` | 排序規則 |  |
| `` g `` | 檢視重設選項 |  |
| `` R `` | 重新命名分支 |  |
| `` u `` | 檢視遠端設定 | 檢視有關遠端分支的設定（例如重設至遠端） |
| `` <c-t> `` | 開啟外部差異工具 (git difftool) |  |
| `` <enter> `` | 檢視提交 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 標籤

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 檢出 | Checkout the selected tag as a detached HEAD. |
| `` n `` | 建立標籤 | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | 刪除 | View delete options for local/remote tag. |
| `` P `` | 推送標籤 | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | 重設 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | 開啟外部差異工具 (git difftool) |  |
| `` <enter> `` | 檢視提交 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |

## 檔案

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製檔案名稱到剪貼簿 |  |
| `` <space> `` | 切換預存 | Toggle staged for selected file. |
| `` <c-b> `` | 篩選檔案 (預存/未預存) |  |
| `` y `` | 複製到剪貼簿 |  |
| `` c `` | 提交變更 | 提交暫存區變更 |
| `` w `` | 沒有預提交 hook 就提交更改 |  |
| `` A `` | 修改上次提交 |  |
| `` C `` | 使用 git 編輯器提交變更 |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | 編輯 | 使用外部編輯器開啟 |
| `` o `` | 開啟檔案 | 使用預設軟體開啟 |
| `` i `` | 忽略或排除檔案 |  |
| `` r `` | 重新整理檔案 |  |
| `` s `` | 收藏 | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | 檢視收藏選項 | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | 全部預存/取消預存 | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | 選擇檔案中的單個程式碼塊/行，或展開/折疊目錄 | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | 捨棄 | 檢視選中變動進行捨棄復原 |
| `` g `` | 檢視遠端重設選項 |  |
| `` D `` | 重設 | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | 顯示檔案樹狀視圖 | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` <c-t> `` | 開啟外部差異工具 (git difftool) |  |
| `` M `` | 開啟外部合併工具 | 執行 `git mergetool`。 |
| `` f `` | 擷取 | 同步遠端異動 |
| `` / `` | 搜尋 |  |

## 狀態

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 開啟設定檔案 | 使用預設軟體開啟 |
| `` e `` | 編輯設定檔案 | 使用外部編輯器開啟 |
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
| `` <enter> `` | View branches |  |
| `` n `` | 新增遠端 |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | 編輯 | 編輯遠端 |
| `` f `` | 擷取 | 擷取遠端 |
| `` / `` | 搜尋 |  |

## 遠端分支

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 複製分支名稱到剪貼簿 |  |
| `` <space> `` | 檢出 | Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head. |
| `` n `` | 新分支 |  |
| `` M `` | 合併到當前檢出的分支 | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` r `` | 將已檢出的分支變基至此分支 | Rebase the checked-out branch onto the selected branch. |
| `` d `` | 刪除 | Delete the remote branch from the remote. |
| `` u `` | 設置為遠端 | 將此分支設為當前分支之遠端 |
| `` s `` | 排序規則 |  |
| `` g `` | 檢視重設選項 | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | 開啟外部差異工具 (git difftool) |  |
| `` <enter> `` | 檢視提交 |  |
| `` w `` | 檢視工作目錄選項 |  |
| `` / `` | 搜尋 |  |
