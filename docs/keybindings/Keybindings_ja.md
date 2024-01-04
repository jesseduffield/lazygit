_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit キーバインド

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## グローバルキーバインド

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | 最近使用したリポジトリに切り替え |  |
| `` <pgup> (fn+up/shift+k) `` | メインパネルを上にスクロール |  |
| `` <pgdown> (fn+down/shift+j) `` | メインパネルを下にスクロール |  |
| `` @ `` | コマンドログメニューを開く |  |
| `` } `` | Increase the size of the context shown around changes in the diff view |  |
| `` { `` | Decrease the size of the context shown around changes in the diff view |  |
| `` : `` | カスタムコマンドを実行 |  |
| `` <c-p> `` | View custom patch options |  |
| `` m `` | View merge/rebase options |  |
| `` R `` | リフレッシュ |  |
| `` + `` | 次のスクリーンモード (normal/half/fullscreen) |  |
| `` _ `` | 前のスクリーンモード |  |
| `` ? `` | メニューを開く |  |
| `` <c-s> `` | View filter-by-path options |  |
| `` W `` | 差分メニューを開く |  |
| `` <c-e> `` | 差分メニューを開く |  |
| `` <c-w> `` | 空白文字の差分の表示有無を切り替え |  |
| `` z `` | アンドゥ (via reflog) (experimental) | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | リドゥ (via reflog) (experimental) | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` P `` | Push |  |
| `` p `` | Pull |  |

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
| `` <space> `` | 適用 |  |
| `` g `` | Pop |  |
| `` d `` | Drop |  |
| `` n `` | 新しいブランチを作成 |  |
| `` r `` | Stashを変更 |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | View selected item's files |  |
| `` / `` | Filter the current view by text |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | コミットのSHAをクリップボードにコピー |  |
| `` w `` | View worktree options |  |
| `` <space> `` | コミットをチェックアウト |  |
| `` y `` | コミットの情報をコピー |  |
| `` o `` | ブラウザでコミットを開く |  |
| `` n `` | コミットにブランチを作成 |  |
| `` g `` | View reset options |  |
| `` C `` | コミットをコピー (cherry-pick) |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View selected item's files |  |
| `` / `` | 検索を開始 |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | Create worktree |  |
| `` <space> `` | Switch to worktree |  |
| `` <enter> `` | Switch to worktree |  |
| `` o `` | Open in editor |  |
| `` d `` | Remove worktree |  |
| `` / `` | Filter the current view by text |  |

## コミット

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | コミットのSHAをクリップボードにコピー |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` b `` | View bisect options |  |
| `` s `` | Squash down |  |
| `` f `` | Fixup commit |  |
| `` r `` | コミットメッセージを変更 |  |
| `` R `` | エディタでコミットメッセージを編集 |  |
| `` d `` | コミットを削除 |  |
| `` e `` | コミットを編集 |  |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick commit (when mid-rebase) |  |
| `` F `` | このコミットに対するfixupコミットを作成 |  |
| `` S `` | Squash all 'fixup!' commits above selected commit (autosquash) |  |
| `` <c-j> `` | コミットを1つ下に移動 |  |
| `` <c-k> `` | コミットを1つ上に移動 |  |
| `` V `` | コミットを貼り付け (cherry-pick) |  |
| `` B `` | Mark commit as base commit for rebase | Select a base commit for the next rebase; this will effectively perform a 'git rebase --onto'. |
| `` A `` | ステージされた変更でamendコミット |  |
| `` a `` | Set/Reset commit author |  |
| `` t `` | コミットをrevert |  |
| `` T `` | タグを作成 |  |
| `` <c-l> `` | ログメニューを開く |  |
| `` w `` | View worktree options |  |
| `` <space> `` | コミットをチェックアウト |  |
| `` y `` | コミットの情報をコピー |  |
| `` o `` | ブラウザでコミットを開く |  |
| `` n `` | コミットにブランチを作成 |  |
| `` g `` | View reset options |  |
| `` C `` | コミットをコピー (cherry-pick) |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View selected item's files |  |
| `` / `` | 検索を開始 |  |

## コミットファイル

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | コミットされたファイル名をクリップボードにコピー |  |
| `` c `` | Checkout file |  |
| `` d `` | Discard this commit's changes to this file |  |
| `` o `` | ファイルを開く |  |
| `` e `` | ファイルを編集 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle file included in patch |  |
| `` a `` | Toggle all files included in patch |  |
| `` <enter> `` | Enter file to add selected lines to the patch (or toggle directory collapsed) |  |
| `` ` `` | ファイルツリーの表示を切り替え |  |
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
| `` <enter> `` | サブモジュールを開く |  |
| `` <space> `` | サブモジュールを開く |  |
| `` d `` | サブモジュールを削除 |  |
| `` u `` | サブモジュールを更新 |  |
| `` n `` | サブモジュールを新規追加 |  |
| `` e `` | サブモジュールのURLを更新 |  |
| `` i `` | サブモジュールを初期化 |  |
| `` b `` | View bulk submodule options |  |
| `` / `` | Filter the current view by text |  |

## ステータス

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 設定ファイルを開く |  |
| `` e `` | 設定ファイルを編集 |  |
| `` u `` | 更新を確認 |  |
| `` <enter> `` | 最近使用したリポジトリに切り替え |  |
| `` a `` | すべてのブランチログを表示 |  |

## タグ

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | チェックアウト |  |
| `` d `` | View delete options |  |
| `` P `` | タグをpush |  |
| `` n `` | タグを作成 |  |
| `` g `` | View reset options |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | コミットを閲覧 |  |
| `` / `` | Filter the current view by text |  |

## ファイル

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | ファイル名をクリップボードにコピー |  |
| `` <space> `` | ステージ/アンステージ |  |
| `` <c-b> `` | ファイルをフィルタ (ステージ/アンステージ) |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | 変更をコミット |  |
| `` w `` | pre-commitフックを実行せずに変更をコミット |  |
| `` A `` | 最新のコミットにamend |  |
| `` C `` | gitエディタを使用して変更をコミット |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | ファイルを編集 |  |
| `` o `` | ファイルを開く |  |
| `` i `` | ファイルをignore |  |
| `` r `` | ファイルをリフレッシュ |  |
| `` s `` | 変更をstash |  |
| `` S `` | View stash options |  |
| `` a `` | すべての変更をステージ/アンステージ |  |
| `` <enter> `` | Stage individual hunks/lines for file, or collapse/expand for directory |  |
| `` d `` | View 'discard changes' options |  |
| `` g `` | View upstream reset options |  |
| `` D `` | View reset options |  |
| `` ` `` | ファイルツリーの表示を切り替え |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Git mergetoolを開く |  |
| `` f `` | Fetch |  |
| `` / `` | 検索を開始 |  |

## ブランチ

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | ブランチ名をクリップボードにコピー |  |
| `` i `` | Show git-flow options |  |
| `` <space> `` | チェックアウト |  |
| `` n `` | 新しいブランチを作成 |  |
| `` o `` | Pull Requestを作成 |  |
| `` O `` | Create pull request options |  |
| `` <c-y> `` | Pull RequestのURLをクリップボードにコピー |  |
| `` c `` | Checkout by name, enter '-' to switch to last |  |
| `` F `` | Force checkout |  |
| `` d `` | View delete options |  |
| `` r `` | Rebase checked-out branch onto this branch |  |
| `` M `` | 現在のブランチにマージ |  |
| `` f `` | Fast-forward this branch from its upstream |  |
| `` T `` | タグを作成 |  |
| `` s `` | 並び替え |  |
| `` g `` | View reset options |  |
| `` R `` | ブランチ名を変更 |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream |
| `` w `` | View worktree options |  |
| `` <enter> `` | コミットを閲覧 |  |
| `` / `` | Filter the current view by text |  |

## メインパネル (Merging)

| Key | Action | Info |
|-----|--------|-------------|
| `` e `` | ファイルを編集 |  |
| `` o `` | ファイルを開く |  |
| `` <left> `` | 前のコンフリクトを選択 |  |
| `` <right> `` | 次のコンフリクトを選択 |  |
| `` <up> `` | 前のhunkを選択 |  |
| `` <down> `` | 次のhunkを選択 |  |
| `` z `` | アンドゥ |  |
| `` M `` | Git mergetoolを開く |  |
| `` <space> `` | Pick hunk |  |
| `` b `` | Pick all hunks |  |
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
| `` a `` | Hunk選択を切り替え |  |
| `` <c-o> `` | 選択されたテキストをクリップボードにコピー |  |
| `` o `` | ファイルを開く |  |
| `` e `` | ファイルを編集 |  |
| `` <space> `` | 行をパッチに追加/削除 |  |
| `` <esc> `` | Exit custom patch builder |  |
| `` / `` | 検索を開始 |  |

## メインパネル (Staging)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 前のhunkを選択 |  |
| `` <right> `` | 次のhunkを選択 |  |
| `` v `` | 範囲選択を切り替え |  |
| `` a `` | Hunk選択を切り替え |  |
| `` <c-o> `` | 選択されたテキストをクリップボードにコピー |  |
| `` o `` | ファイルを開く |  |
| `` e `` | ファイルを編集 |  |
| `` <esc> `` | ファイル一覧に戻る |  |
| `` <tab> `` | パネルを切り替え |  |
| `` <space> `` | 選択行をステージ/アンステージ |  |
| `` d `` | 変更を削除 (git reset) |  |
| `` E `` | Edit hunk |  |
| `` c `` | 変更をコミット |  |
| `` w `` | pre-commitフックを実行せずに変更をコミット |  |
| `` C `` | gitエディタを使用して変更をコミット |  |
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
| `` f `` | リモートをfetch |  |
| `` n `` | リモートを新規追加 |  |
| `` d `` | リモートを削除 |  |
| `` e `` | リモートを編集 |  |
| `` / `` | Filter the current view by text |  |

## リモートブランチ

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | ブランチ名をクリップボードにコピー |  |
| `` <space> `` | チェックアウト |  |
| `` n `` | 新しいブランチを作成 |  |
| `` M `` | 現在のブランチにマージ |  |
| `` r `` | Rebase checked-out branch onto this branch |  |
| `` d `` | Delete remote tag |  |
| `` u `` | Set as upstream of checked-out branch |  |
| `` s `` | 並び替え |  |
| `` g `` | View reset options |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | コミットを閲覧 |  |
| `` / `` | Filter the current view by text |  |

## 参照ログ

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | コミットのSHAをクリップボードにコピー |  |
| `` w `` | View worktree options |  |
| `` <space> `` | コミットをチェックアウト |  |
| `` y `` | コミットの情報をコピー |  |
| `` o `` | ブラウザでコミットを開く |  |
| `` n `` | コミットにブランチを作成 |  |
| `` g `` | View reset options |  |
| `` C `` | コミットをコピー (cherry-pick) |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | コミットを閲覧 |  |
| `` / `` | Filter the current view by text |  |

## 確認パネル

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 確認 |  |
| `` <esc> `` | 閉じる/キャンセル |  |
