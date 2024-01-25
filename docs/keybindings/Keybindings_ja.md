_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit キーバインド

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## グローバルキーバインド

<pre>
  <kbd>&lt;c-r&gt;</kbd>: 最近使用したリポジトリに切り替え
  <kbd>&lt;pgup&gt;</kbd>: メインパネルを上にスクロール (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: メインパネルを下にスクロール (fn+down/shift+j)
  <kbd>@</kbd>: コマンドログメニューを開く
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: カスタムコマンドを実行
  <kbd>&lt;c-p&gt;</kbd>: View custom patch options
  <kbd>m</kbd>: View merge/rebase options
  <kbd>R</kbd>: リフレッシュ
  <kbd>+</kbd>: 次のスクリーンモード (normal/half/fullscreen)
  <kbd>_</kbd>: 前のスクリーンモード
  <kbd>?</kbd>: メニューを開く
  <kbd>&lt;c-s&gt;</kbd>: View filter-by-path options
  <kbd>W</kbd>: 差分メニューを開く
  <kbd>&lt;c-e&gt;</kbd>: 差分メニューを開く
  <kbd>&lt;c-w&gt;</kbd>: 空白文字の差分の表示有無を切り替え
  <kbd>z</kbd>: アンドゥ (via reflog) (experimental)
  <kbd>&lt;c-z&gt;</kbd>: リドゥ (via reflog) (experimental)
  <kbd>P</kbd>: Push
  <kbd>p</kbd>: Pull
</pre>

## 一覧パネルの操作

<pre>
  <kbd>,</kbd>: 前のページ
  <kbd>.</kbd>: 次のページ
  <kbd>&lt;</kbd>: 最上部までスクロール
  <kbd>&gt;</kbd>: 最下部までスクロール
  <kbd>v</kbd>: 範囲選択を切り替え
  <kbd>&lt;s-down&gt;</kbd>: Range select down
  <kbd>&lt;s-up&gt;</kbd>: Range select up
  <kbd>/</kbd>: 検索を開始
  <kbd>H</kbd>: 左スクロール
  <kbd>L</kbd>: 右スクロール
  <kbd>]</kbd>: 次のタブ
  <kbd>[</kbd>: 前のタブ
</pre>

## Stash

<pre>
  <kbd>&lt;space&gt;</kbd>: 適用
  <kbd>g</kbd>: Pop
  <kbd>d</kbd>: Drop
  <kbd>n</kbd>: 新しいブランチを作成
  <kbd>r</kbd>: Stashを変更
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: View selected item's files
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Sub-commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: コミットのSHAをクリップボードにコピー
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: コミットをチェックアウト
  <kbd>y</kbd>: コミットの情報をコピー
  <kbd>o</kbd>: ブラウザでコミットを開く
  <kbd>n</kbd>: コミットにブランチを作成
  <kbd>g</kbd>: View reset options
  <kbd>C</kbd>: コミットをコピー (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>&lt;enter&gt;</kbd>: View selected item's files
  <kbd>/</kbd>: 検索を開始
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

## コミット

<pre>
  <kbd>&lt;c-o&gt;</kbd>: コミットのSHAをクリップボードにコピー
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>b</kbd>: View bisect options
  <kbd>s</kbd>: Squash down
  <kbd>f</kbd>: Fixup commit
  <kbd>r</kbd>: コミットメッセージを変更
  <kbd>R</kbd>: エディタでコミットメッセージを編集
  <kbd>d</kbd>: コミットを削除
  <kbd>e</kbd>: コミットを編集
  <kbd>i</kbd>: Start interactive rebase
  <kbd>p</kbd>: Pick commit (when mid-rebase)
  <kbd>F</kbd>: このコミットに対するfixupコミットを作成
  <kbd>S</kbd>: Squash all 'fixup!' commits above selected commit (autosquash)
  <kbd>&lt;c-j&gt;</kbd>: コミットを1つ下に移動
  <kbd>&lt;c-k&gt;</kbd>: コミットを1つ上に移動
  <kbd>V</kbd>: コミットを貼り付け (cherry-pick)
  <kbd>B</kbd>: Mark commit as base commit for rebase
  <kbd>A</kbd>: ステージされた変更でamendコミット
  <kbd>a</kbd>: Set/Reset commit author
  <kbd>t</kbd>: コミットをrevert
  <kbd>T</kbd>: タグを作成
  <kbd>&lt;c-l&gt;</kbd>: ログメニューを開く
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: コミットをチェックアウト
  <kbd>y</kbd>: コミットの情報をコピー
  <kbd>o</kbd>: ブラウザでコミットを開く
  <kbd>n</kbd>: コミットにブランチを作成
  <kbd>g</kbd>: View reset options
  <kbd>C</kbd>: コミットをコピー (cherry-pick)
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>&lt;enter&gt;</kbd>: View selected item's files
  <kbd>/</kbd>: 検索を開始
</pre>

## コミットファイル

<pre>
  <kbd>&lt;c-o&gt;</kbd>: コミットされたファイル名をクリップボードにコピー
  <kbd>c</kbd>: Checkout file
  <kbd>d</kbd>: Discard this commit's changes to this file
  <kbd>o</kbd>: ファイルを開く
  <kbd>e</kbd>: ファイルを編集
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>&lt;space&gt;</kbd>: Toggle file included in patch
  <kbd>a</kbd>: Toggle all files included in patch
  <kbd>&lt;enter&gt;</kbd>: Enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: ファイルツリーの表示を切り替え
  <kbd>/</kbd>: 検索を開始
</pre>

## コミットメッセージ

<pre>
  <kbd>&lt;enter&gt;</kbd>: 確認
  <kbd>&lt;esc&gt;</kbd>: 閉じる
</pre>

## サブモジュール

<pre>
  <kbd>&lt;c-o&gt;</kbd>: サブモジュール名をクリップボードにコピー
  <kbd>&lt;enter&gt;</kbd>: サブモジュールを開く
  <kbd>&lt;space&gt;</kbd>: サブモジュールを開く
  <kbd>d</kbd>: サブモジュールを削除
  <kbd>u</kbd>: サブモジュールを更新
  <kbd>n</kbd>: サブモジュールを新規追加
  <kbd>e</kbd>: サブモジュールのURLを更新
  <kbd>i</kbd>: サブモジュールを初期化
  <kbd>b</kbd>: View bulk submodule options
  <kbd>/</kbd>: Filter the current view by text
</pre>

## ステータス

<pre>
  <kbd>o</kbd>: 設定ファイルを開く
  <kbd>e</kbd>: 設定ファイルを編集
  <kbd>u</kbd>: 更新を確認
  <kbd>&lt;enter&gt;</kbd>: 最近使用したリポジトリに切り替え
  <kbd>a</kbd>: すべてのブランチログを表示
</pre>

## タグ

<pre>
  <kbd>&lt;space&gt;</kbd>: チェックアウト
  <kbd>d</kbd>: View delete options
  <kbd>P</kbd>: タグをpush
  <kbd>n</kbd>: タグを作成
  <kbd>g</kbd>: View reset options
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: コミットを閲覧
  <kbd>/</kbd>: Filter the current view by text
</pre>

## ファイル

<pre>
  <kbd>&lt;c-o&gt;</kbd>: ファイル名をクリップボードにコピー
  <kbd>&lt;space&gt;</kbd>: ステージ/アンステージ
  <kbd>&lt;c-b&gt;</kbd>: ファイルをフィルタ (ステージ/アンステージ)
  <kbd>y</kbd>: Copy to clipboard
  <kbd>c</kbd>: 変更をコミット
  <kbd>w</kbd>: pre-commitフックを実行せずに変更をコミット
  <kbd>A</kbd>: 最新のコミットにamend
  <kbd>C</kbd>: gitエディタを使用して変更をコミット
  <kbd>&lt;c-f&gt;</kbd>: Find base commit for fixup
  <kbd>e</kbd>: ファイルを編集
  <kbd>o</kbd>: ファイルを開く
  <kbd>i</kbd>: ファイルをignore
  <kbd>r</kbd>: ファイルをリフレッシュ
  <kbd>s</kbd>: 変更をstash
  <kbd>S</kbd>: View stash options
  <kbd>a</kbd>: すべての変更をステージ/アンステージ
  <kbd>&lt;enter&gt;</kbd>: Stage individual hunks/lines for file, or collapse/expand for directory
  <kbd>d</kbd>: View 'discard changes' options
  <kbd>g</kbd>: View upstream reset options
  <kbd>D</kbd>: View reset options
  <kbd>`</kbd>: ファイルツリーの表示を切り替え
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>M</kbd>: Git mergetoolを開く
  <kbd>f</kbd>: Fetch
  <kbd>/</kbd>: 検索を開始
</pre>

## ブランチ

<pre>
  <kbd>&lt;c-o&gt;</kbd>: ブランチ名をクリップボードにコピー
  <kbd>i</kbd>: Show git-flow options
  <kbd>&lt;space&gt;</kbd>: チェックアウト
  <kbd>n</kbd>: 新しいブランチを作成
  <kbd>o</kbd>: Pull Requestを作成
  <kbd>O</kbd>: Create pull request options
  <kbd>&lt;c-y&gt;</kbd>: Pull RequestのURLをクリップボードにコピー
  <kbd>c</kbd>: Checkout by name, enter '-' to switch to last
  <kbd>F</kbd>: Force checkout
  <kbd>d</kbd>: View delete options
  <kbd>r</kbd>: Rebase checked-out branch onto this branch
  <kbd>M</kbd>: 現在のブランチにマージ
  <kbd>f</kbd>: Fast-forward this branch from its upstream
  <kbd>T</kbd>: タグを作成
  <kbd>s</kbd>: 並び替え
  <kbd>g</kbd>: View reset options
  <kbd>R</kbd>: ブランチ名を変更
  <kbd>u</kbd>: View upstream options
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: コミットを閲覧
  <kbd>/</kbd>: Filter the current view by text
</pre>

## メインパネル (Merging)

<pre>
  <kbd>e</kbd>: ファイルを編集
  <kbd>o</kbd>: ファイルを開く
  <kbd>&lt;left&gt;</kbd>: 前のコンフリクトを選択
  <kbd>&lt;right&gt;</kbd>: 次のコンフリクトを選択
  <kbd>&lt;up&gt;</kbd>: 前のhunkを選択
  <kbd>&lt;down&gt;</kbd>: 次のhunkを選択
  <kbd>z</kbd>: アンドゥ
  <kbd>M</kbd>: Git mergetoolを開く
  <kbd>&lt;space&gt;</kbd>: Pick hunk
  <kbd>b</kbd>: Pick all hunks
  <kbd>&lt;esc&gt;</kbd>: ファイル一覧に戻る
</pre>

## メインパネル (Normal)

<pre>
  <kbd>mouse wheel down</kbd>: 下にスクロール (fn+up)
  <kbd>mouse wheel up</kbd>: 上にスクロール (fn+down)
</pre>

## メインパネル (Patch Building)

<pre>
  <kbd>&lt;left&gt;</kbd>: 前のhunkを選択
  <kbd>&lt;right&gt;</kbd>: 次のhunkを選択
  <kbd>v</kbd>: 範囲選択を切り替え
  <kbd>a</kbd>: Hunk選択を切り替え
  <kbd>&lt;c-o&gt;</kbd>: 選択されたテキストをクリップボードにコピー
  <kbd>o</kbd>: ファイルを開く
  <kbd>e</kbd>: ファイルを編集
  <kbd>&lt;space&gt;</kbd>: 行をパッチに追加/削除
  <kbd>&lt;esc&gt;</kbd>: Exit custom patch builder
  <kbd>/</kbd>: 検索を開始
</pre>

## メインパネル (Staging)

<pre>
  <kbd>&lt;left&gt;</kbd>: 前のhunkを選択
  <kbd>&lt;right&gt;</kbd>: 次のhunkを選択
  <kbd>v</kbd>: 範囲選択を切り替え
  <kbd>a</kbd>: Hunk選択を切り替え
  <kbd>&lt;c-o&gt;</kbd>: 選択されたテキストをクリップボードにコピー
  <kbd>o</kbd>: ファイルを開く
  <kbd>e</kbd>: ファイルを編集
  <kbd>&lt;esc&gt;</kbd>: ファイル一覧に戻る
  <kbd>&lt;tab&gt;</kbd>: パネルを切り替え
  <kbd>&lt;space&gt;</kbd>: 選択行をステージ/アンステージ
  <kbd>d</kbd>: 変更を削除 (git reset)
  <kbd>E</kbd>: Edit hunk
  <kbd>c</kbd>: 変更をコミット
  <kbd>w</kbd>: pre-commitフックを実行せずに変更をコミット
  <kbd>C</kbd>: gitエディタを使用して変更をコミット
  <kbd>/</kbd>: 検索を開始
</pre>

## メニュー

<pre>
  <kbd>&lt;enter&gt;</kbd>: 実行
  <kbd>&lt;esc&gt;</kbd>: 閉じる
  <kbd>/</kbd>: Filter the current view by text
</pre>

## リモート

<pre>
  <kbd>f</kbd>: リモートをfetch
  <kbd>n</kbd>: リモートを新規追加
  <kbd>d</kbd>: リモートを削除
  <kbd>e</kbd>: リモートを編集
  <kbd>/</kbd>: Filter the current view by text
</pre>

## リモートブランチ

<pre>
  <kbd>&lt;c-o&gt;</kbd>: ブランチ名をクリップボードにコピー
  <kbd>&lt;space&gt;</kbd>: チェックアウト
  <kbd>n</kbd>: 新しいブランチを作成
  <kbd>M</kbd>: 現在のブランチにマージ
  <kbd>r</kbd>: Rebase checked-out branch onto this branch
  <kbd>d</kbd>: Delete remote tag
  <kbd>u</kbd>: Set as upstream of checked-out branch
  <kbd>s</kbd>: 並び替え
  <kbd>g</kbd>: View reset options
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: コミットを閲覧
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 参照ログ

<pre>
  <kbd>&lt;c-o&gt;</kbd>: コミットのSHAをクリップボードにコピー
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: コミットをチェックアウト
  <kbd>y</kbd>: コミットの情報をコピー
  <kbd>o</kbd>: ブラウザでコミットを開く
  <kbd>n</kbd>: コミットにブランチを作成
  <kbd>g</kbd>: View reset options
  <kbd>C</kbd>: コミットをコピー (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>&lt;c-t&gt;</kbd>: Open external diff tool (git difftool)
  <kbd>&lt;enter&gt;</kbd>: コミットを閲覧
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 確認パネル

<pre>
  <kbd>&lt;enter&gt;</kbd>: 確認
  <kbd>&lt;esc&gt;</kbd>: 閉じる/キャンセル
</pre>
