_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit キーバインド

## グローバルキーバインド

<pre>
  <kbd>ctrl+r</kbd>: 最近使用したリポジトリに切り替え
  <kbd>pgup</kbd>: メインパネルを上にスクロール (fn+up/shift+k)
  <kbd>pgdown</kbd>: メインパネルを下にスクロール (fn+down/shift+j)
  <kbd>m</kbd>: view merge/rebase options
  <kbd>ctrl+p</kbd>: view custom patch options
  <kbd>R</kbd>: リフレッシュ
  <kbd>x</kbd>: メニューを開く
  <kbd>+</kbd>: 次のスクリーンモード (normal/half/fullscreen)
  <kbd>_</kbd>: 前のスクリーンモード
  <kbd>ctrl+s</kbd>: view filter-by-path options
  <kbd>W</kbd>: 差分メニューを開く
  <kbd>ctrl+e</kbd>: 差分メニューを開く
  <kbd>@</kbd>: コマンドログメニューを開く
  <kbd>}</kbd>: Increase the size of the context shown around changes in the diff view
  <kbd>{</kbd>: Decrease the size of the context shown around changes in the diff view
  <kbd>:</kbd>: カスタムコマンドを実行
  <kbd>z</kbd>: アンドゥ (via reflog) (experimental)
  <kbd>ctrl+z</kbd>: リドゥ (via reflog) (experimental)
  <kbd>P</kbd>: push
  <kbd>p</kbd>: pull
</pre>

## 一覧パネルの操作

<pre>
  <kbd>,</kbd>: 前のページ
  <kbd>.</kbd>: 次のページ
  <kbd><</kbd>: 最上部までスクロール
  <kbd>/</kbd>: 検索を開始
  <kbd>></kbd>: 最下部までスクロール
  <kbd>H</kbd>: 左スクロール
  <kbd>L</kbd>: 右スクロール
  <kbd>]</kbd>: 次のタブ
  <kbd>[</kbd>: 前のタブ
</pre>

## Stash

<pre>
  <kbd>space</kbd>: 適用
  <kbd>g</kbd>: pop
  <kbd>d</kbd>: drop
  <kbd>n</kbd>: 新しいブランチを作成
  <kbd>r</kbd>: Stashを変更
  <kbd>enter</kbd>: view selected item's files
</pre>

## Sub-commits

<pre>
  <kbd>ctrl+o</kbd>: コミットのSHAをクリップボードにコピー
  <kbd>space</kbd>: コミットをチェックアウト
  <kbd>y</kbd>: コミットの情報をコピー
  <kbd>o</kbd>: ブラウザでコミットを開く
  <kbd>n</kbd>: コミットにブランチを作成
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: コミットをコピー (cherry-pick)
  <kbd>C</kbd>: コミットを範囲コピー (cherry-pick)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>enter</kbd>: view selected item's files
</pre>

## コミット

<pre>
  <kbd>ctrl+o</kbd>: コミットのSHAをクリップボードにコピー
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>b</kbd>: view bisect options
  <kbd>s</kbd>: squash down
  <kbd>f</kbd>: fixup commit
  <kbd>r</kbd>: コミットメッセージを変更
  <kbd>R</kbd>: エディタでコミットメッセージを編集
  <kbd>d</kbd>: コミットを削除
  <kbd>e</kbd>: コミットを編集
  <kbd>p</kbd>: pick commit (when mid-rebase)
  <kbd>F</kbd>: このコミットに対するfixupコミットを作成
  <kbd>S</kbd>: squash all 'fixup!' commits above selected commit (autosquash)
  <kbd>ctrl+j</kbd>: コミットを1つ下に移動
  <kbd>ctrl+k</kbd>: コミットを1つ上に移動
  <kbd>v</kbd>: コミットを貼り付け (cherry-pick)
  <kbd>A</kbd>: ステージされた変更でamendコミット
  <kbd>a</kbd>: reset commit author
  <kbd>t</kbd>: コミットをrevert
  <kbd>T</kbd>: タグを作成
  <kbd>ctrl+l</kbd>: ログメニューを開く
  <kbd>space</kbd>: コミットをチェックアウト
  <kbd>y</kbd>: コミットの情報をコピー
  <kbd>o</kbd>: ブラウザでコミットを開く
  <kbd>n</kbd>: コミットにブランチを作成
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: コミットをコピー (cherry-pick)
  <kbd>C</kbd>: コミットを範囲コピー (cherry-pick)
  <kbd>enter</kbd>: view selected item's files
</pre>

## コミットファイル

<pre>
  <kbd>ctrl+o</kbd>: コミットされたファイル名をクリップボードにコピー
  <kbd>c</kbd>: checkout file
  <kbd>d</kbd>: discard this commit's changes to this file
  <kbd>o</kbd>: ファイルを開く
  <kbd>e</kbd>: ファイルを編集
  <kbd>space</kbd>: toggle file included in patch
  <kbd>a</kbd>: toggle all files included in patch
  <kbd>enter</kbd>: enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: ファイルツリーの表示を切り替え
</pre>

## サブモジュール

<pre>
  <kbd>ctrl+o</kbd>: サブモジュール名をクリップボードにコピー
  <kbd>enter</kbd>: サブモジュールを開く
  <kbd>d</kbd>: サブモジュールを削除
  <kbd>u</kbd>: サブモジュールを更新
  <kbd>n</kbd>: サブモジュールを新規追加
  <kbd>e</kbd>: サブモジュールのURLを更新
  <kbd>i</kbd>: サブモジュールを初期化
  <kbd>b</kbd>: view bulk submodule options
</pre>

## ステータス

<pre>
  <kbd>e</kbd>: 設定ファイルを編集
  <kbd>o</kbd>: 設定ファイルを開く
  <kbd>u</kbd>: 更新を確認
  <kbd>enter</kbd>: 最近使用したリポジトリに切り替え
  <kbd>a</kbd>: すべてのブランチログを表示
</pre>

## タグ

<pre>
  <kbd>space</kbd>: チェックアウト
  <kbd>d</kbd>: タグを削除
  <kbd>P</kbd>: タグをpush
  <kbd>n</kbd>: タグを作成
  <kbd>g</kbd>: view reset options
  <kbd>enter</kbd>: コミットを閲覧
</pre>

## ファイル

<pre>
  <kbd>ctrl+o</kbd>: ファイル名をクリップボードにコピー
  <kbd>ctrl+w</kbd>: 空白文字の差分の表示有無を切り替え
  <kbd>d</kbd>: view 'discard changes' options
  <kbd>space</kbd>: ステージ/アンステージ
  <kbd>ctrl+b</kbd>: ファイルをフィルタ (ステージ/アンステージ)
  <kbd>c</kbd>: 変更をコミット
  <kbd>w</kbd>: pre-commitフックを実行せずに変更をコミット
  <kbd>A</kbd>: 最新のコミットにamend
  <kbd>C</kbd>: gitエディタを使用して変更をコミット
  <kbd>e</kbd>: ファイルを編集
  <kbd>o</kbd>: ファイルを開く
  <kbd>i</kbd>: ファイルをignore
  <kbd>r</kbd>: ファイルをリフレッシュ
  <kbd>s</kbd>: 変更をstash
  <kbd>S</kbd>: view stash options
  <kbd>a</kbd>: すべての変更をステージ/アンステージ
  <kbd>enter</kbd>: stage individual hunks/lines for file, or collapse/expand for directory
  <kbd>g</kbd>: view upstream reset options
  <kbd>D</kbd>: view reset options
  <kbd>`</kbd>: ファイルツリーの表示を切り替え
  <kbd>M</kbd>: git mergetoolを開く
  <kbd>f</kbd>: fetch
</pre>

## ブランチ

<pre>
  <kbd>ctrl+o</kbd>: ブランチ名をクリップボードにコピー
  <kbd>i</kbd>: show git-flow options
  <kbd>space</kbd>: チェックアウト
  <kbd>n</kbd>: 新しいブランチを作成
  <kbd>o</kbd>: Pull Requestを作成
  <kbd>O</kbd>: create pull request options
  <kbd>ctrl+y</kbd>: Pull RequestのURLをクリップボードにコピー
  <kbd>c</kbd>: checkout by name
  <kbd>F</kbd>: force checkout
  <kbd>d</kbd>: ブランチを削除
  <kbd>r</kbd>: rebase checked-out branch onto this branch
  <kbd>M</kbd>: 現在のブランチにマージ
  <kbd>f</kbd>: fast-forward this branch from its upstream
  <kbd>g</kbd>: view reset options
  <kbd>R</kbd>: ブランチ名を変更
  <kbd>u</kbd>: set/unset upstream
  <kbd>enter</kbd>: コミットを閲覧
</pre>

## メインパネル (Merging)

<pre>
  <kbd>e</kbd>: ファイルを編集
  <kbd>o</kbd>: ファイルを開く
  <kbd>◄</kbd>: 前のコンフリクトを選択
  <kbd>►</kbd>: 次のコンフリクトを選択
  <kbd>▲</kbd>: 前のhunkを選択
  <kbd>▼</kbd>: 次のhunkを選択
  <kbd>z</kbd>: アンドゥ
  <kbd>M</kbd>: git mergetoolを開く
  <kbd>space</kbd>: pick hunk
  <kbd>b</kbd>: pick all hunks
  <kbd>esc</kbd>: ファイル一覧に戻る
</pre>

## メインパネル (Normal)

<pre>
  <kbd>mouse wheel ▼</kbd>: 下にスクロール (fn+up)
  <kbd>mouse wheel ▲</kbd>: 上にスクロール (fn+down)
</pre>

## メインパネル (Patch Building)

<pre>
  <kbd>◄</kbd>: 前のhunkを選択
  <kbd>►</kbd>: 次のhunkを選択
  <kbd>v</kbd>: 範囲選択を切り替え
  <kbd>V</kbd>: 範囲選択を切り替え
  <kbd>a</kbd>: hunk選択を切り替え
  <kbd>ctrl+o</kbd>: 選択されたテキストをクリップボードにコピー
  <kbd>o</kbd>: ファイルを開く
  <kbd>e</kbd>: ファイルを編集
  <kbd>space</kbd>: 行をパッチに追加/削除
  <kbd>esc</kbd>: exit custom patch builder
</pre>

## メインパネル (Staging)

<pre>
  <kbd>◄</kbd>: 前のhunkを選択
  <kbd>►</kbd>: 次のhunkを選択
  <kbd>v</kbd>: 範囲選択を切り替え
  <kbd>V</kbd>: 範囲選択を切り替え
  <kbd>a</kbd>: hunk選択を切り替え
  <kbd>ctrl+o</kbd>: 選択されたテキストをクリップボードにコピー
  <kbd>o</kbd>: ファイルを開く
  <kbd>e</kbd>: ファイルを編集
  <kbd>esc</kbd>: ファイル一覧に戻る
  <kbd>tab</kbd>: パネルを切り替え
  <kbd>space</kbd>: 選択行をステージ/アンステージ
  <kbd>d</kbd>: 変更を削除 (git reset)
  <kbd>E</kbd>: edit hunk
  <kbd>c</kbd>: 変更をコミット
  <kbd>w</kbd>: pre-commitフックを実行せずに変更をコミット
  <kbd>C</kbd>: gitエディタを使用して変更をコミット
</pre>

## リモート

<pre>
  <kbd>f</kbd>: リモートをfetch
  <kbd>n</kbd>: リモートを新規追加
  <kbd>d</kbd>: リモートを削除
  <kbd>e</kbd>: リモートを編集
</pre>

## リモートブランチ

<pre>
  <kbd>space</kbd>: チェックアウト
  <kbd>n</kbd>: 新しいブランチを作成
  <kbd>M</kbd>: 現在のブランチにマージ
  <kbd>r</kbd>: rebase checked-out branch onto this branch
  <kbd>d</kbd>: ブランチを削除
  <kbd>u</kbd>: set as upstream of checked-out branch
  <kbd>esc</kbd>: リモート一覧に戻る
  <kbd>g</kbd>: view reset options
  <kbd>enter</kbd>: コミットを閲覧
</pre>

## 参照ログ

<pre>
  <kbd>ctrl+o</kbd>: コミットのSHAをクリップボードにコピー
  <kbd>space</kbd>: コミットをチェックアウト
  <kbd>y</kbd>: コミットの情報をコピー
  <kbd>o</kbd>: ブラウザでコミットを開く
  <kbd>n</kbd>: コミットにブランチを作成
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: コミットをコピー (cherry-pick)
  <kbd>C</kbd>: コミットを範囲コピー (cherry-pick)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>enter</kbd>: コミットを閲覧
</pre>
