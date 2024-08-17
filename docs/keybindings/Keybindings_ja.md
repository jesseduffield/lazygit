_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit キーバインド

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## グローバルキーバインド

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | 最近使用したリポジトリに切り替え |  |
| `` <pgup> (fn+up/shift+k) `` | メインパネルを上にスクロール |  |
| `` <pgdown> (fn+down/shift+j) `` | メインパネルを下にスクロール |  |
| `` @ `` | コマンドログメニューを開く | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | Push | Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` p `` | Pull | Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` } `` | Increase diff context size | Increase the amount of the context shown around changes in the diff view. |
| `` { `` | Decrease diff context size | Decrease the amount of the context shown around changes in the diff view. |
| `` : `` | Execute shell command | Bring up a prompt where you can enter a shell command to execute. |
| `` <c-p> `` | View custom patch options |  |
| `` m `` | View merge/rebase options | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | リフレッシュ | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | 次のスクリーンモード (normal/half/fullscreen) |  |
| `` _ `` | 前のスクリーンモード |  |
| `` ? `` | メニューを開く |  |
| `` <c-s> `` | View filter options | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | 差分メニューを開く | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | 差分メニューを開く | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | 終了 |  |
| `` <esc> `` | キャンセル |  |
| `` <c-w> `` | 空白文字の差分の表示有無を切り替え | Toggle whether or not whitespace changes are shown in the diff view. |
| `` z `` | アンドゥ (via reflog) (experimental) | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | リドゥ (via reflog) (experimental) | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |

## 一覧パネルの操作

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | 前のページ |  |
| `` . `` | 次のページ |  |
| `` < `` | 最上部までスクロール |  |
| `` > `` | 最下部までスクロール |  |
| `` v `` | 範囲選択を切り替え |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | 検索を開始 |  |
| `` H `` | 左スクロール |  |
| `` L `` | 右スクロール |  |
| `` ] `` | 次のタブ |  |
| `` [ `` | 前のタブ |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 適用 | Apply the stash entry to your working directory. |
| `` g `` | Pop | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | Drop | Remove the stash entry from the stash list. |
| `` n `` | 新しいブランチを作成 | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | Stashを変更 |  |
| `` <enter> `` | View files |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | コミットのhashをクリップボードにコピー |  |
| `` <space> `` | チェックアウト | Checkout the selected commit as a detached HEAD. |
| `` y `` | コミットの情報をコピー | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | ブラウザでコミットを開く |  |
| `` n `` | コミットにブランチを作成 |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | コミットをコピー (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View files |  |
| `` w `` | View worktree options |  |
| `` / `` | 検索を開始 |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | Open in editor |  |
| `` d `` | Remove | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
| `` / `` | Filter the current view by text |  |

## コミット

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | コミットのhashをクリップボードにコピー |  |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` b `` | View bisect options |  |
| `` s `` | Squash | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | Fixup | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` r `` | コミットメッセージを変更 | Reword the selected commit's message. |
| `` R `` | エディタでコミットメッセージを編集 |  |
| `` d `` | コミットを削除 | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | Edit (start interactive rebase) | コミットを編集 |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick | Mark the selected commit to be picked (when mid-rebase). This means that the commit will be retained upon continuing the rebase. |
| `` F `` | Fixupコミットを作成 | このコミットに対するfixupコミットを作成 |
| `` S `` | Apply fixup commits | Squash all 'fixup!' commits, either above the selected commit, or all in current branch (autosquash). |
| `` <c-j> `` | コミットを1つ下に移動 |  |
| `` <c-k> `` | コミットを1つ上に移動 |  |
| `` V `` | コミットを貼り付け (cherry-pick) |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Amend | ステージされた変更でamendコミット |
| `` a `` | Amend commit attribute | Set/Reset commit author or set co-author. |
| `` t `` | Revert | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | タグを作成 | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <c-l> `` | ログメニューを開く | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` <space> `` | チェックアウト | Checkout the selected commit as a detached HEAD. |
| `` y `` | コミットの情報をコピー | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | ブラウザでコミットを開く |  |
| `` n `` | コミットにブランチを作成 |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | コミットをコピー (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View files |  |
| `` w `` | View worktree options |  |
| `` / `` | 検索を開始 |  |

## コミットファイル

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | ファイル名をクリップボードにコピー |  |
| `` c `` | チェックアウト | Checkout file. This replaces the file in your working tree with the version from the selected commit. |
| `` d `` | Remove | Discard this commit's changes to this file. This runs an interactive rebase in the background, so you may get a merge conflict if a later commit also changes this file. |
| `` o `` | ファイルを開く | Open file in default application. |
| `` e `` | Edit | Open file in external editor. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle file included in patch | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Toggle all files | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Enter file / Toggle directory collapsed | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | ファイルツリーの表示を切り替え | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` / `` | 検索を開始 |  |

## コミットメッセージ

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 確認 |  |
| `` <esc> `` | 閉じる |  |

## サブモジュール

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | サブモジュール名をクリップボードにコピー |  |
| `` <enter> `` | Enter | サブモジュールを開く |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | サブモジュールを更新 |
| `` n `` | サブモジュールを新規追加 |  |
| `` e `` | サブモジュールのURLを更新 |  |
| `` i `` | Initialize | サブモジュールを初期化 |
| `` b `` | View bulk submodule options |  |
| `` / `` | Filter the current view by text |  |

## ステータス

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 設定ファイルを開く | Open file in default application. |
| `` e `` | 設定ファイルを編集 | Open file in external editor. |
| `` u `` | 更新を確認 |  |
| `` <enter> `` | 最近使用したリポジトリに切り替え |  |
| `` a `` | すべてのブランチログを表示 |  |

## タグ

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | チェックアウト | Checkout the selected tag tag as a detached HEAD. |
| `` n `` | タグを作成 | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Delete | View delete options for local/remote tag. |
| `` P `` | タグをpush | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | コミットを閲覧 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## ファイル

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | ファイル名をクリップボードにコピー |  |
| `` <space> `` | ステージ/アンステージ | Toggle staged for selected file. |
| `` <c-b> `` | ファイルをフィルタ (ステージ/アンステージ) |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | 変更をコミット | Commit staged changes. |
| `` w `` | pre-commitフックを実行せずに変更をコミット |  |
| `` A `` | 最新のコミットにamend |  |
| `` C `` | gitエディタを使用して変更をコミット |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit | Open file in external editor. |
| `` o `` | ファイルを開く | Open file in default application. |
| `` i `` | ファイルをignore |  |
| `` r `` | ファイルをリフレッシュ |  |
| `` s `` | Stash | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | View stash options | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | すべての変更をステージ/アンステージ | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | Stage lines / Collapse directory | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | Discard | View options for discarding changes to the selected file. |
| `` g `` | View upstream reset options |  |
| `` D `` | Reset | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | ファイルツリーの表示を切り替え | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Git mergetoolを開く | Run `git mergetool`. |
| `` f `` | Fetch | Fetch changes from remote. |
| `` / `` | 検索を開始 |  |

## ブランチ

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | ブランチ名をクリップボードにコピー |  |
| `` i `` | Show git-flow options |  |
| `` <space> `` | チェックアウト | Checkout selected item. |
| `` n `` | 新しいブランチを作成 |  |
| `` o `` | Pull Requestを作成 |  |
| `` O `` | View create pull request options |  |
| `` <c-y> `` | Pull RequestのURLをクリップボードにコピー |  |
| `` c `` | Checkout by name | Checkout by name. In the input box you can enter '-' to switch to the last branch. |
| `` F `` | Force checkout | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | Delete | View delete options for local/remote branch. |
| `` r `` | Rebase | Rebase the checked-out branch onto the selected branch. |
| `` M `` | 現在のブランチにマージ | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` f `` | Fast-forward | Fast-forward selected branch from its upstream. |
| `` T `` | タグを作成 |  |
| `` s `` | 並び替え |  |
| `` g `` | Reset |  |
| `` R `` | ブランチ名を変更 |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | コミットを閲覧 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## メインパネル (Merging)

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | Pick hunk |  |
| `` b `` | Pick all hunks |  |
| `` <up> `` | 前のhunkを選択 |  |
| `` <down> `` | 次のhunkを選択 |  |
| `` <left> `` | 前のコンフリクトを選択 |  |
| `` <right> `` | 次のコンフリクトを選択 |  |
| `` z `` | アンドゥ | Undo last merge conflict resolution. |
| `` e `` | ファイルを編集 | Open file in external editor. |
| `` o `` | ファイルを開く | Open file in default application. |
| `` M `` | Git mergetoolを開く | Run `git mergetool`. |
| `` <esc> `` | ファイル一覧に戻る |  |

## メインパネル (Normal)

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | 下にスクロール |  |
| `` mouse wheel up (fn+down) `` | 上にスクロール |  |

## メインパネル (Patch Building)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 前のhunkを選択 |  |
| `` <right> `` | 次のhunkを選択 |  |
| `` v `` | 範囲選択を切り替え |  |
| `` a `` | Hunk選択を切り替え | Toggle hunk selection mode. |
| `` <c-o> `` | 選択されたテキストをクリップボードにコピー |  |
| `` o `` | ファイルを開く | Open file in default application. |
| `` e `` | ファイルを編集 | Open file in external editor. |
| `` <space> `` | 行をパッチに追加/削除 |  |
| `` <esc> `` | Exit custom patch builder |  |
| `` / `` | 検索を開始 |  |

## メインパネル (Staging)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 前のhunkを選択 |  |
| `` <right> `` | 次のhunkを選択 |  |
| `` v `` | 範囲選択を切り替え |  |
| `` a `` | Hunk選択を切り替え | Toggle hunk selection mode. |
| `` <c-o> `` | 選択されたテキストをクリップボードにコピー |  |
| `` <space> `` | ステージ/アンステージ | 選択行をステージ/アンステージ |
| `` d `` | 変更を削除 (git reset) | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | ファイルを開く | Open file in default application. |
| `` e `` | ファイルを編集 | Open file in external editor. |
| `` <esc> `` | ファイル一覧に戻る |  |
| `` <tab> `` | パネルを切り替え | Switch to other view (staged/unstaged changes). |
| `` E `` | Edit hunk | Edit selected hunk in external editor. |
| `` c `` | 変更をコミット | Commit staged changes. |
| `` w `` | pre-commitフックを実行せずに変更をコミット |  |
| `` C `` | gitエディタを使用して変更をコミット |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | 検索を開始 |  |

## メニュー

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 実行 |  |
| `` <esc> `` | 閉じる |  |
| `` / `` | Filter the current view by text |  |

## リモート

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | リモートを新規追加 |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Edit | リモートを編集 |
| `` f `` | Fetch | リモートをfetch |
| `` / `` | Filter the current view by text |  |

## リモートブランチ

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | ブランチ名をクリップボードにコピー |  |
| `` <space> `` | チェックアウト | Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head. |
| `` n `` | 新しいブランチを作成 |  |
| `` M `` | 現在のブランチにマージ | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` r `` | Rebase | Rebase the checked-out branch onto the selected branch. |
| `` d `` | Delete | Delete the remote branch from the remote. |
| `` u `` | Set as upstream | Set the selected remote branch as the upstream of the checked-out branch. |
| `` s `` | 並び替え |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | コミットを閲覧 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## 参照ログ

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | コミットのhashをクリップボードにコピー |  |
| `` <space> `` | チェックアウト | Checkout the selected commit as a detached HEAD. |
| `` y `` | コミットの情報をコピー | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | ブラウザでコミットを開く |  |
| `` n `` | コミットにブランチを作成 |  |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | コミットをコピー (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset copied (cherry-picked) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | コミットを閲覧 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## 確認パネル

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 確認 |  |
| `` <esc> `` | 閉じる/キャンセル |  |
