package i18n

const japaneseIntroPopupMessage = `
Thanks for using lazygit! Seriously you rock. Three things to share with you:

 1) If you want to learn about lazygit's features, watch this vid:
      https://youtu.be/CPLdltN7wgE

 2) Be sure to read the latest release notes at:
      https://github.com/jesseduffield/lazygit/releases

 3) If you're using git, that makes you a programmer! With your help we can make
    lazygit better, so consider becoming a contributor and joining the fun at
      https://github.com/jesseduffield/lazygit
    You can also sponsor me and tell me what to work on by clicking the donate
    button at the bottom right.
    Or even just star the repo to share the love!
`

// exporting this so we can use it in tests
func japaneseTranslationSet() TranslationSet {
	return TranslationSet{
		NotEnoughSpace:          "パネルの描画に十分な空間がありません",
		DiffTitle:               "差分",
		FilesTitle:              "ファイル",
		BranchesTitle:           "ブランチ",
		CommitsTitle:            "コミット",
		StashTitle:              "Stash",
		UnstagedChanges:         `ステージされていない変更`,
		StagedChanges:           `ステージされた変更`,
		MainTitle:               "メイン",
		MergeConfirmTitle:       "マージ",
		StagingTitle:            "メインパネル (Staging)",
		MergingTitle:            "メインパネル (Merging)",
		NormalTitle:             "メインパネル (Normal)",
		LogTitle:                "ログ",
		CommitSummary:           "コミットメッセージ",
		CredentialsUsername:     "ユーザ名",
		CredentialsPassword:     "パスワード",
		CredentialsPassphrase:   "SSH鍵のパスフレーズを入力",
		PassUnameWrong:          "パスワード, パスフレーズまたはユーザ名が間違っています。",
		CommitChanges:           "変更をコミット",
		AmendLastCommit:         "最新のコミットにamend",
		AmendLastCommitTitle:    "最新のコミットにamend",
		SureToAmend:             "最新のコミットに変更をamendします。よろしいですか? コミットメッセージはコミットパネルから変更できます。",
		NoCommitToAmend:         "Amend可能なコミットが存在しません。",
		CommitChangesWithEditor: "gitエディタを使用して変更をコミット",
		StatusTitle:             "ステータス",
		Menu:                    "メニュー",
		Execute:                 "実行",
		ToggleStaged:            "ステージ/アンステージ",
		ToggleStagedAll:         "すべての変更をステージ/アンステージ",
		ToggleTreeView:          "ファイルツリーの表示を切り替え",
		OpenMergeTool:           "Git mergetoolを開く",
		Refresh:                 "リフレッシュ",
		Push:                    "Push",
		Pull:                    "Pull",
		Scroll:                  "スクロール",
		MergeConflictsTitle:     "マージコンフリクト",
		Checkout:                "チェックアウト",
		FileFilter:              "ファイルをフィルタ (ステージ/アンステージ)",
		FilterStagedFiles:       "ステージされたファイルのみを表示",
		FilterUnstagedFiles:     "ステージされていないファイルのみを表示",
		ResetFilter:             "フィルタをリセット",
		// NoChangedFiles:                      "No changed files",
		PullWait:                "Pull中...",
		PushWait:                "Push中...",
		FetchWait:               "Fetch中...",
		SoftReset:               "Softリセット",
		AlreadyCheckedOutBranch: "ブランチはすでにチェックアウトされています。",
		// SureForceCheckout:                   "Are you sure you want force checkout? You will lose all local changes",
		// ForceCheckoutBranch:                 "Force Checkout Branch",
		BranchName:               "ブランチ名",
		NewBranchNameBranchOff:   "新規ブランチ名 ('{{.branchName}}' に作成)",
		CantDeleteCheckOutBranch: "チェックアウト中のブランチは削除できません!",
		DeleteBranch:             "ブランチを削除",
		DeleteBranchMessage:      "ブランチ '{{.selectedBranchName}}' を削除します。よろしいですか?",
		ForceDeleteBranchMessage: "'{{.selectedBranchName}}' はマージされていません。本当に削除しますか?",
		// LcRebaseBranch:                      "Rebase checked-out branch onto this branch",
		CantRebaseOntoSelf:        "ブランチを自分自身にリベースすることはできません。",
		CantMergeBranchIntoItself: "ブランチを自分自身にマージすることはできません。",
		// LcForceCheckout:                     "Force checkout",
		// LcCheckoutByName:                    "Checkout by name",
		NewBranch:               "新しいブランチを作成",
		NoBranchesThisRepo:      "リポジトリにブランチが存在しません",
		CommitWithoutMessageErr: "コミットメッセージを入力してください",
		CloseCancel:             "閉じる/キャンセル",
		Confirm:                 "確認",
		Close:                   "閉じる",
		Quit:                    "終了",
		// LcSquashDown:                        "Squash down",
		// LcFixupCommit:                       "Fixup commit",
		// NoCommitsThisBranch:                 "No commits for this branch",
		// CannotSquashOrFixupFirstCommit:      "There's no commit below to squash into",
		// Fixup:                               "Fixup",
		// SureFixupThisCommit:                 "Are you sure you want to 'fixup' this commit? It will be merged into the commit below",
		// SureSquashThisCommit:                "Are you sure you want to squash this commit into the commit below?",
		// Squash:                              "Squash",
		// LcPickCommit:                        "Pick commit (when mid-rebase)",
		RevertCommit:       "コミットをrevert",
		RewordCommit:       "コミットメッセージを変更",
		DeleteCommit:       "コミットを削除",
		MoveDownCommit:     "コミットを1つ下に移動",
		MoveUpCommit:       "コミットを1つ上に移動",
		EditCommit:         "コミットを編集",
		AmendToCommit:      "ステージされた変更でamendコミット",
		RenameCommitEditor: "エディタでコミットメッセージを編集",
		Error:              "エラー",
		// LcPickHunk:                          "Pick hunk",
		// LcPickAllHunks:                      "Pick all hunks",
		Undo:                "アンドゥ",
		UndoReflog:          "アンドゥ (via reflog) (experimental)",
		RedoReflog:          "リドゥ (via reflog) (experimental)",
		Pop:                 "Pop",
		Drop:                "Drop",
		Apply:               "適用",
		NoStashEntries:      "Stashが存在しません",
		StashDrop:           "Stashを削除",
		SureDropStashEntry:  "Stashを削除します。よろしいですか?",
		StashPop:            "Stashをpop",
		SurePopStashEntry:   "Stashをpopします。よろしいですか?",
		StashApply:          "Stashを適用",
		SureApplyStashEntry: "Stashを適用します。よろしいですか?",
		// NoTrackedStagedFilesStash:           "You have no tracked/staged files to stash",
		StashChanges:      "変更をStash",
		RenameStash:       "Stashを変更",
		RenameStashPrompt: "Stash名を変更: {{.stashName}}",
		OpenConfig:        "設定ファイルを開く",
		EditConfig:        "設定ファイルを編集",
		ForcePush:         "Force push",
		ForcePushPrompt:   "ブランチがリモートブランチから分岐しています。'esc'でキャンセル, または'enter'でforce pushします。",
		ForcePushDisabled: "ブランチがリモートブランチから分岐しています。force pushは無効化されています。",
		// UpdatesRejectedAndForcePushDisabled: "Updates were rejected and you have disabled force pushing",
		CheckForUpdate:                   "更新を確認",
		CheckingForUpdates:               "更新を確認中...",
		UpdateAvailableTitle:             "最新リリース!",
		UpdateAvailable:                  "バージョン {{.newVersion}} をインストールしますか?",
		UpdateInProgressWaitingStatus:    "更新中",
		UpdateCompletedTitle:             "更新完了!",
		UpdateCompleted:                  "更新のインストールに成功しました。lazygitを再起動してください。",
		FailedToRetrieveLatestVersionErr: "バージョン情報の取得に失敗しました",
		OnLatestVersionErr:               "使用中のバージョンは最新です",
		MajorVersionErr:                  "新バージョン ({{.newVersion}}) は現在のバージョン ({{.currentVersion}}) と後方互換性がありません。",
		CouldNotFindBinaryErr:            "{{.url}} にバイナリが存在しませんでした。",
		UpdateFailedErr:                  "更新失敗: {{.errMessage}}",
		ConfirmQuitDuringUpdateTitle:     "現在更新中",
		ConfirmQuitDuringUpdate:          "現在更新を実行中です。終了しますか?",
		MergeToolTitle:                   "マージツール",
		MergeToolPrompt:                  "`git mergetool`を開きます。よろしいですか?",
		IntroPopupMessage:                japaneseIntroPopupMessage,
		// GitconfigParseErr:                   `Gogit failed to parse your gitconfig file due to the presence of unquoted '\' characters. Removing these should fix the issue.`,
		EditFile:               `ファイルを編集`,
		OpenFile:               `ファイルを開く`,
		IgnoreFile:             `.gitignoreに追加`,
		RefreshFiles:           `ファイルをリフレッシュ`,
		MergeIntoCurrentBranch: `現在のブランチにマージ`,
		ConfirmQuit:            `終了します。よろしいですか?`,
		SwitchRepo:             `最近使用したリポジトリに切り替え`,
		AllBranchesLogGraph:    `すべてのブランチログを表示`,
		UnsupportedGitService:  `サポートされていないGitサービスです。`,
		CreatePullRequest:      `Pull Requestを作成`,
		CopyPullRequestURL:     `Pull RequestのURLをクリップボードにコピー`,
		NoBranchOnRemote:       `ブランチがリモートに存在しません。リモートにpushしてください。`,
		Fetch:                  `Fetch`,
		// NoAutomaticGitFetchTitle:            `No automatic git fetch`,
		// NoAutomaticGitFetchBody:             `Lazygit can't use "git fetch" in a private repo; use 'f' in the files panel to run "git fetch" manually`,
		// FileEnter:                           `stage individual hunks/lines for file, or collapse/expand for directory`,
		// FileStagingRequirements:             `Can only stage individual lines for tracked files`,
		StageSelection:          `選択行をステージ/アンステージ`,
		DiscardSelection:        `変更を削除 (git reset)`,
		ToggleDragSelect:        `範囲選択を切り替え`,
		ToggleSelectHunk:        `Hunk選択を切り替え`,
		ToggleSelectionForPatch: `行をパッチに追加/削除`,
		ToggleStagingPanel:      `パネルを切り替え`,
		ReturnToFilesPanel:      `ファイル一覧に戻る`,
		// FastForward:                         `fast-forward this branch from its upstream`,
		// Fetching:                            "Fetching and fast-forwarding {{.from}} -> {{.to}} ...",
		// FoundConflicts:                      "Conflicts! To abort press 'esc', otherwise press 'enter'",
		// FoundConflictsTitle:                 "Auto-merge failed",
		// PickHunk:                            "Pick hunk",
		// PickAllHunks:                        "Pick all hunks",
		// ViewMergeRebaseOptions:              "View merge/rebase options",
		// NotMergingOrRebasing:                "You are currently neither rebasing nor merging",
		RecentRepos: "最近使用したリポジトリ",
		// MergeOptionsTitle:                   "Merge Options",
		// RebaseOptionsTitle:                  "Rebase Options",
		CommitSummaryTitle:  "コミットメッセージ",
		LocalBranchesTitle:  "ブランチ",
		SearchTitle:         "検索",
		TagsTitle:           "タグ",
		MenuTitle:           "メニュー",
		RemotesTitle:        "リモート",
		RemoteBranchesTitle: "リモートブランチ",
		PatchBuildingTitle:  "メインパネル (Patch Building)",
		InformationTitle:    "Information",
		SecondaryTitle:      "Secondary",
		ReflogCommitsTitle:  "参照ログ",
		GlobalTitle:         "グローバルキーバインド",
		// ConflictsResolved:                   "All merge conflicts resolved. Continue?",
		// RebasingTitle:                       "Rebasing",
		// ConfirmRebase:                       "Are you sure you want to rebase '{{.checkedOutBranch}}' onto '{{.selectedBranch}}'?",
		// ConfirmMerge:                        "Are you sure you want to merge '{{.selectedBranch}}' into '{{.checkedOutBranch}}'?",
		// FwdNoUpstream:                       "Cannot fast-forward a branch with no upstream",
		// FwdNoLocalUpstream:                  "Cannot fast-forward a branch whose remote is not registered locally",
		// FwdCommitsToPush:                    "Cannot fast-forward a branch with commits to push",
		ErrorOccurred: "エラーが発生しました! issueを作成してください: ",
		// NoRoom:                              "Not enough room",
		YouAreHere: "現在位置",
		// LcRewordNotSupported:                "Rewording commits while interactively rebasing is not currently supported",
		CherryPickCopy:      "コミットをコピー (cherry-pick)",
		CherryPickCopyRange: "コミットを範囲コピー (cherry-pick)",
		PasteCommits:        "コミットを貼り付け (cherry-pick)",
		// SureCherryPick:                      "Are you sure you want to cherry-pick the copied commits onto this branch?",
		CherryPick:          "Cherry-Pick",
		Donate:              "支援",
		AskQuestion:         "質問",
		PrevLine:            "前の行を選択",
		NextLine:            "次の行を選択",
		PrevHunk:            "前のhunkを選択",
		NextHunk:            "次のhunkを選択",
		PrevConflict:        "前のコンフリクトを選択",
		NextConflict:        "次のコンフリクトを選択",
		SelectPrevHunk:      "前のhunkを選択",
		SelectNextHunk:      "次のhunkを選択",
		ScrollDown:          "下にスクロール",
		ScrollUp:            "上にスクロール",
		ScrollUpMainPanel:   "メインパネルを上にスクロール",
		ScrollDownMainPanel: "メインパネルを下にスクロール",
		AmendCommitTitle:    "Amendコミット",
		AmendCommitPrompt:   "ステージされたファイルで現在のコミットをamendします。よろしいですか?",
		DeleteCommitTitle:   "コミットを削除",
		DeleteCommitPrompt:  "選択されたコミットを削除します。よろしいですか?",
		// SquashingStatus:                     "Squashing",
		// FixingStatus:                        "Fixing up",
		// DeletingStatus:                      "Deleting",
		// MovingStatus:                        "Moving",
		// RebasingStatus:                      "Rebasing",
		// AmendingStatus:                      "Amending",
		// CherryPickingStatus:                 "Cherry-picking",
		// UndoingStatus:                       "Undoing",
		// RedoingStatus:                       "Redoing",
		// CheckingOutStatus:                   "Checking out",
		// CommittingStatus:                    "Committing",
		CommitFiles:                "Commit files",
		SubCommitsDynamicTitle:     "コミット (%s)",
		CommitFilesDynamicTitle:    "Diff files (%s)",
		RemoteBranchesDynamicTitle: "リモートブランチ (%s)",
		// LcViewItemFiles:                     "View selected item's files",
		CommitFilesTitle: "コミットファイル",
		// LcCheckoutCommitFile:                "Checkout file",
		// LcDiscardOldFileChange:              "Discard this commit's changes to this file",
		DiscardFileChangesTitle: "ファイルの変更を破棄",
		// DiscardFileChangesPrompt:            "Are you sure you want to discard this commit's changes to this file? If this file was created in this commit, it will be deleted",
		// DisabledForGPG:                      "Feature not available for users using GPG",
		CreateRepo: "Gitリポジトリではありません。リポジトリを作成しますか? (y/n): ",
		// AutoStashTitle:                      "Autostash?",
		// AutoStashPrompt:                     "You must stash and pop your changes to bring them across. Do this automatically? (enter/esc)",
		// StashPrefix:                         "Auto-stashing changes for ",
		// LcViewDiscardOptions:                "View 'discard changes' options",
		Cancel:            "キャンセル",
		DiscardAllChanges: "すべての変更を破棄",
		// LcDiscardUnstagedChanges:            "Discard unstaged changes",
		// LcDiscardAllChangesToAllFiles:       "Nuke working tree",
		// LcDiscardAnyUnstagedChanges:         "Discard unstaged changes",
		// LcDiscardUntrackedFiles:             "Discard untracked files",
		HardReset: "hardリセット",
		// LcViewResetOptions:                  `view reset options`,
		CreateFixupCommitDescription: `このコミットに対するfixupコミットを作成`,
		// LcSquashAboveCommits:                `squash all 'fixup!' commits above selected commit (autosquash)`,
		// SquashAboveCommits:                  `Squash all 'fixup!' commits above selected commit (autosquash)`,
		SureSquashAboveCommits:   `{{.commit}}に対するすべての fixup! コミットをsquashします。よろしいですか?`,
		CreateFixupCommit:        `Fixupコミットを作成`,
		SureCreateFixupCommit:    `{{.commit}} に対する fixup! コミットを作成します。よろしいですか?`,
		ExecuteCustomCommand:     "カスタムコマンドを実行",
		CustomCommand:            "カスタムコマンド:",
		CommitChangesWithoutHook: "pre-commitフックを実行せずに変更をコミット",
		// SkipHookPrefixNotConfigured:         "You have not configured a commit message prefix for skipping hooks. Set `git.skipHookPrefix = 'WIP'` in your config",
		// LcResetTo:                           `reset to`,
		PressEnterToReturn: "Enterを入力してください",
		// LcViewStashOptions:                  "View stash options",
		StashAllChanges: "変更をstash",
		// LcStashStagedChanges:                "Stash staged changes",
		// LcStashOptions:                      "Stash options",
		// NotARepository:                      "Error: must be run inside a git repository",
		Jump:            "パネルに移動",
		ScrollLeftRight: "左右にスクロール",
		ScrollLeft:      "左スクロール",
		ScrollRight:     "右スクロール",
		DiscardPatch:    "パッチを破棄",
		// DiscardPatchConfirm:                 "You can only build a patch from one commit/stash-entry at a time. Discard current patch?",
		// CantPatchWhileRebasingError:         "You cannot build a patch or run patch commands while in a merging or rebasing state",
		// LcToggleAddToPatch:                  "Toggle file included in patch",
		// LcToggleAllInPatch:                  "Toggle all files included in patch",
		// LcUpdatingPatch:                     "Updating patch",
		// ViewPatchOptions:                    "View custom patch options",
		// PatchOptionsTitle:                   "Patch Options",
		// NoPatchError:                        "No patch created yet. To start building a patch, use 'space' on a commit file or enter to add specific lines",
		// LcEnterFile:                         "Enter file to add selected lines to the patch (or toggle directory collapsed)",
		// ExitCustomPatchBuilder:    ``,
		EnterUpstream:             `'<remote> <branchname>' の形式でupstreamを入力`,
		InvalidUpstream:           "Upstreamの形式が正しくありません。'<remote> <branchname>' の形式で入力してください。",
		ReturnToRemotesList:       `リモート一覧に戻る`,
		AddNewRemote:              `リモートを新規追加`,
		NewRemoteName:             `新規リモート名:`,
		NewRemoteUrl:              `新規リモートURL:`,
		EditRemoteName:            `{{.remoteName}} の新しいリモート名を入力:`,
		EditRemoteUrl:             `{{.remoteName}} の新しいリモートURLを入力:`,
		RemoveRemote:              `リモートを削除`,
		RemoveRemotePrompt:        "リモートを削除します。よろしいですか?",
		DeleteRemoteBranch:        "リモートブランチを削除",
		DeleteRemoteBranchMessage: "リモートブランチを削除します。よろしいですか",
		// LcSetUpstream:                       "Set as upstream of checked-out branch",
		// SetUpstreamTitle:                    "Set upstream branch",
		// SetUpstreamMessage:                  "Are you sure you want to set the upstream branch of '{{.checkedOut}}' to '{{.selected}}'",
		EditRemote:             "リモートを編集",
		TagCommit:              "タグを作成",
		TagMenuTitle:           "タグを作成",
		TagNameTitle:           "タグ名",
		TagMessageTitle:        "タグメッセージ",
		AnnotatedTag:           "注釈付きタグ",
		LightweightTag:         "軽量タグ",
		DeleteTag:              "タグを削除",
		DeleteTagTitle:         "タグを削除",
		PushTagTitle:           "リモートにタグ '{{.tagName}}' をpush",
		PushTag:                "タグをpush",
		CreateTag:              "タグを作成",
		FetchRemote:            "リモートをfetch",
		FetchingRemoteStatus:   "リモートをfetch",
		CheckoutCommit:         "コミットをチェックアウト",
		SureCheckoutThisCommit: "選択されたコミットをチェックアウトします。よろしいですか?",
		// LcGitFlowOptions:                    "Show git-flow options",
		// NotAGitFlowBranch:                   "This does not seem to be a git flow branch",
		// NewGitFlowBranchPrompt:              "New {{.branchType}} name:",
		// IgnoreTracked:                       "Ignore tracked file",
		// IgnoreTrackedPrompt:                 "Are you sure you want to ignore a tracked file?",
		// LcViewResetToUpstreamOptions:        "View upstream reset options",
		NextScreenMode:      "次のスクリーンモード (normal/half/fullscreen)",
		PrevScreenMode:      "前のスクリーンモード",
		StartSearch:         "検索を開始",
		Panel:               "パネル",
		Keybindings:         "キーバインド",
		RenameBranch:        "ブランチ名を変更",
		NewBranchNamePrompt: "新しいブランチ名を入力",
		// RenameBranchWarning:                 "This branch is tracking a remote. This action will only rename the local branch name, not the name of the remote branch. Continue?",
		OpenMenu: "メニューを開く",
		// LcResetCherryPick:                   "Reset cherry-picked (copied) commits selection",
		NextTab:               "次のタブ",
		PrevTab:               "前のタブ",
		CantUndoWhileRebasing: "リベース中はアンドゥできません。",
		CantRedoWhileRebasing: "リベース中はリドゥできません。",
		// MustStashWarning:                    "Pulling a patch out into the index requires stashing and unstashing your changes. If something goes wrong, you'll be able to access your files from the stash. Continue?",
		// MustStashTitle:                      "Must stash",
		ConfirmationTitle: "確認パネル",
		PrevPage:          "前のページ",
		NextPage:          "次のページ",
		GotoTop:           "最上部までスクロール",
		GotoBottom:        "最下部までスクロール",
		// LcFilteringBy:                       "Filtering by",
		// ResetInParentheses:                  "(reset)",
		// LcOpenFilteringMenu:                 "View filter-by-path options",
		// LcFilterBy:                          "Filter by",
		// LcExitFilterMode:                    "Stop filtering by path",
		// LcFilterPathOption:                  "Enter path to filter by",
		// EnterFileName:                       "Enter path:",
		// FilteringMenuTitle:                  "Filtering",
		// MustExitFilterModeTitle:             "Command not available",
		// MustExitFilterModePrompt:            "Command not available in filtered mode. Exit filtered mode?",
		Diff: "差分",
		// LcEnterRefToDiff:                    "Enter ref to diff",
		EnterRefName:     "参照を入力:",
		ExitDiffMode:     "差分モードを終了",
		DiffingMenuTitle: "差分",
		// LcSwapDiff:                          "Reverse diff direction",
		OpenDiffingMenu: "差分メニューを開く",
		// // the actual view is the extras view which I intend to give more tabs in future but for now we'll only mention the command log part
		OpenExtrasMenu: "コマンドログメニューを開く",
		// LcShowingGitDiff:                    "Showing output for:",
		CommitDiff:                     "コミットの差分",
		CopyCommitShaToClipboard:       "コミットのSHAをクリップボードにコピー",
		CommitSha:                      "コミットのSHA",
		CommitURL:                      "コミットのURL",
		CopyCommitMessageToClipboard:   "コミットメッセージをクリップボードにコピー",
		CommitMessage:                  "コミットメッセージ",
		CommitAuthor:                   "コミットの作成者名",
		CopyCommitAttributeToClipboard: "コミットの情報をコピー",
		CopyBranchNameToClipboard:      "ブランチ名をクリップボードにコピー",
		CopyFileNameToClipboard:        "ファイル名をクリップボードにコピー",
		CopyCommitFileNameToClipboard:  "コミットされたファイル名をクリップボードにコピー",
		CopySelectedTexToClipboard:     "選択されたテキストをクリップボードにコピー",
		// LcCommitPrefixPatternError:          "Error in commitPrefix pattern",
		NoFilesStagedTitle:         "ファイルがステージされていません",
		NoFilesStagedPrompt:        "ファイルがステージされていません。すべての変更をコミットしますか?",
		BranchNotFoundTitle:        "ブランチが見つかりませんでした。",
		BranchNotFoundPrompt:       "ブランチが見つかりませんでした。新しくブランチを作成します ",
		DiscardChangeTitle:         "選択行をアンステージ",
		DiscardChangePrompt:        "選択された行を削除 (git reset) します。よろしいですか? この操作は取り消せません。\nこの警告を無効化するには設定ファイルの 'gui.skipDiscardChangeWarning' を true に設定してください。",
		CreateNewBranchFromCommit:  "コミットにブランチを作成",
		BuildingPatch:              "パッチを構築",
		ViewCommits:                "コミットを閲覧",
		MinGitVersionError:         "Lazygitの実行にはGit 2.20以降のバージョンが必要です。Gitを更新してください。もしくは、lazygitの後方互換性を改善するために https://github.com/jesseduffield/lazygit/issues にissueを作成してください。",
		RunningCustomCommandStatus: "カスタムコマンドを実行",
		// LcSubmoduleStashAndReset:            "Stash uncommitted submodule changes and update",
		// LcAndResetSubmodules:                "And reset submodules",
		EnterSubmodule:               "サブモジュールを開く",
		CopySubmoduleNameToClipboard: "サブモジュール名をクリップボードにコピー",
		RemoveSubmodule:              "サブモジュールを削除",
		RemoveSubmodulePrompt:        "サブモジュール '%s' とそのディレクトリを削除します。よろしいですか? この操作は取り消せません。",
		ResettingSubmoduleStatus:     "サブモジュールをリセット",
		NewSubmoduleName:             "新規サブモジュール名:",
		NewSubmoduleUrl:              "新規サブモジュールのURL:",
		NewSubmodulePath:             "新規サブモジュールのパス:",
		AddSubmodule:                 "サブモジュールを新規追加",
		AddingSubmoduleStatus:        "サブモジュールを新規追加",
		UpdateSubmoduleUrl:           "サブモジュール '%s' のURLを更新",
		UpdatingSubmoduleUrlStatus:   "URLを更新",
		EditSubmoduleUrl:             "サブモジュールのURLを更新",
		InitializingSubmoduleStatus:  "サブモジュールを初期化",
		InitSubmodule:                "サブモジュールを初期化",
		SubmoduleUpdate:              "サブモジュールを更新",
		UpdatingSubmoduleStatus:      "サブモジュールを更新",
		BulkInitSubmodules:           "サブモジュールを一括初期化",
		BulkUpdateSubmodules:         "サブモジュールを一括更新",
		// LcBulkDeinitSubmodules:              "Bulk deinit submodules",
		// LcViewBulkSubmoduleOptions:          "View bulk submodule options",
		// LcBulkSubmoduleOptions:              "Bulk submodule options",
		// LcRunningCommand:                    "Running command",
		// SubCommitsTitle:                     "Sub-commits",
		SubmodulesTitle: "サブモジュール",
		NavigationTitle: "一覧パネルの操作",
		// SuggestionsCheatsheetTitle:          "Suggestions",
		// SuggestionsTitle:                    "Suggestions (press %s to focus)",
		ExtrasTitle: "コマンドログ",
		// PushingTagStatus:                    "Pushing tag",
		PullRequestURLCopiedToClipboard:     "Pull requestのURLがクリップボードにコピーされました",
		CommitDiffCopiedToClipboard:         "コミットの差分がクリップボードにコピーされました",
		CommitSHACopiedToClipboard:          "コミットのSHAがクリップボードにコピーされました",
		CommitURLCopiedToClipboard:          "コミットのURLがクリップボードにコピーされました",
		CommitMessageCopiedToClipboard:      "コミットメッセージがクリップボードにコピーされました",
		CommitAuthorCopiedToClipboard:       "コミットの作成者名がクリップボードにコピーされました",
		CopiedToClipboard:                   "クリップボードにコピーされました",
		ErrCannotEditDirectory:              "ディレクトリは編集できません。",
		ErrStageDirWithInlineMergeConflicts: "マージコンフリクトの発生したファイルを含むディレクトリはステージ/アンステージできません。マージコンフリクトを解決してください。",
		ErrRepositoryMovedOrDeleted:         "リポジトリが見つかりません。すでに削除されたか、移動された可能性があります ¯\\_(ツ)_/¯",
		CommandLog:                          "コマンドログ",
		ToggleShowCommandLog:                "コマンドログの表示/非表示を切り替え",
		FocusCommandLog:                     "コマンドログにフォーカス",
		CommandLogHeader:                    "コマンドログの表示/非表示は '%s' で切り替えられます。\n",
		RandomTip:                           "ランダムTips",
		// SelectParentCommitForMerge:          "Select parent commit for merge",
		ToggleWhitespaceInDiffView: "空白文字の差分の表示有無を切り替え",
		// IncreaseContextInDiffView:           "Increase the size of the context shown around changes in the diff view",
		// DecreaseContextInDiffView:           "Decrease the size of the context shown around changes in the diff view",
		// CreatePullRequestOptions:            "Create pull request options",
		// LcCreatePullRequestOptions:          "Create pull request options",
		DefaultBranch:        "デフォルトブランチ",
		SelectBranch:         "ブランチを選択",
		SelectConfigFile:     "設定ファイルを選択",
		NoConfigFileFoundErr: "設定ファイルが見つかりませんでした。",
		// LcLoadingFileSuggestions:            "Loading file suggestions",
		// LcLoadingCommits:                    "Loading commits",
		// MustSpecifyOriginError:              "Must specify a remote if specifying a branch",
		// GitOutput:                           "Git output:",
		// GitCommandFailed:                    "Git command failed. Check command log for details (open with %s)",
		AbortTitle:   "%sを中止",
		AbortPrompt:  "実施中の%sを中止します。よろしいですか?",
		OpenLogMenu:  "ログメニューを開く",
		LogMenuTitle: "コミットログオプション",
		// ToggleShowGitGraphAll:               "Toggle show whole git graph (pass the `--all` flag to `git log`)",
		ShowGitGraph: "コミットグラフの表示",
		SortCommits:  "コミットの表示順",
		// CantChangeContextSizeError:          "Cannot change context while in patch building mode because we were too lazy to support it when releasing the feature. If you really want it, please let us know!",
		OpenCommitInBrowser: "ブラウザでコミットを開く",
		// LcViewBisectOptions:                 "View bisect options",
		// ConfirmRevertCommit:                 "Are you sure you want to revert {{.selectedCommit}}?",
		RewordInEditorTitle: "コミットメッセージをエディタで編集",
		// RewordInEditorPrompt:                "Are you sure you want to reword this commit in your editor?",
		// HardResetAutostashPrompt:            "Are you sure you want to hard reset to '%s'? An auto-stash will be performed if necessary.",
		// CheckoutPrompt:                      "Are you sure you want to checkout '%s'?",
		// UpstreamGone:                        "(upstream gone)",
		Actions: Actions{
			// TODO: combine this with the original keybinding descriptions (those are all in lowercase atm)
			CheckoutCommit:      "コミットをチェックアウト",
			CheckoutTag:         "タグをチェックアウト",
			CheckoutBranch:      "ブランチをチェックアウト",
			ForceCheckoutBranch: "ブランチを強制的にチェックアウト",
			DeleteBranch:        "ブランチを削除",
			Merge:               "マージ",
			// RebaseBranch:                      "Rebase branch",
			RenameBranch: "ブランチ名を変更",
			CreateBranch: "ブランチを作成",
			// CherryPick:                        "(Cherry-pick) Paste commits",
			CheckoutFile: "ファイルをチェックアウトs",
			// DiscardOldFileChange:              "Discard old file change",
			// SquashCommitDown:                  "Squash commit down",
			FixupCommit:       "Fixupコミット",
			RewordCommit:      "コミットメッセージを変更",
			DropCommit:        "コミットを削除",
			EditCommit:        "コミットを編集",
			AmendCommit:       "Amendコミット",
			RevertCommit:      "コミットをrevert",
			CreateFixupCommit: "fixupコミットを作成",
			// SquashAllAboveFixupCommits:        "Squash all above fixup commits",
			CreateLightweightTag:              "軽量タグを作成",
			CreateAnnotatedTag:                "注釈付きタグを作成",
			CopyCommitMessageToClipboard:      "コミットメッセージをクリップボードにコピー",
			CopyCommitDiffToClipboard:         "コミットの差分をクリップボードにコピー",
			CopyCommitSHAToClipboard:          "コミットSHAをクリップボードにコピー",
			CopyCommitURLToClipboard:          "コミットのURLをクリップボードにコピー",
			CopyCommitAuthorToClipboard:       "コミットの作成者名をクリップボードにコピー",
			CopyCommitAttributeToClipboard:    "クリップボードにコピー",
			MoveCommitUp:                      "コミットを上に移動",
			MoveCommitDown:                    "コミットを下に移動",
			CustomCommand:                     "カスタムコマンド",
			DiscardAllChangesInDirectory:      "ディレクトリ内のすべての変更を破棄",
			DiscardUnstagedChangesInDirectory: "ディレクトリ内のすべてのステージされていない変更を破棄",
			DiscardAllChangesInFile:           "ファイル内のすべての変更を破棄",
			DiscardAllUnstagedChangesInFile:   "ファイル内のすべてのステージされていない変更を破棄",
			StageFile:                         "ファイルをステージ",
			StageResolvedFiles:                "マージコンフリクトが解決されたすべてのファイルをステージ",
			UnstageFile:                       "ファイルをアンステージ",
			UnstageAllFiles:                   "すべてのファイルをアンステージ",
			StageAllFiles:                     "すべてのファイルをステージ",
			IgnoreExcludeFile:                 "ファイルをignore",
			Commit:                            "コミット",
			EditFile:                          "ファイルを編集",
			Push:                              "Push",
			Pull:                              "Pull",
			OpenFile:                          "ファイルを開く",
			StashAllChanges:                   "すべての変更をStash",
			StashStagedChanges:                "ステージされた変更をStash",
			GitFlowFinish:                     "Git flow finish",
			GitFlowStart:                      "Git Flow start",
			CopyToClipboard:                   "クリップボードにコピー",
			CopySelectedTextToClipboard:       "選択されたテキストをクリップボードにコピー",
			RemovePatchFromCommit:             "パッチをコミットから削除",
			MovePatchToSelectedCommit:         "パッチを選択したコミットに移動",
			MovePatchIntoIndex:                "パッチをindexに移動",
			MovePatchIntoNewCommit:            "パッチを次のコミットに移動",
			DeleteRemoteBranch:                "リモートブランチを削除",
			SetBranchUpstream:                 "Upstreamブランチを設定",
			AddRemote:                         "リモートを追加",
			RemoveRemote:                      "リモートを削除",
			UpdateRemote:                      "リモートを更新",
			ApplyPatch:                        "パッチを適用",
			Stash:                             "Stash",
			RenameStash:                       "Stash名を変更",
			RemoveSubmodule:                   "サブモジュールを削除",
			ResetSubmodule:                    "サブモジュールをリセット",
			AddSubmodule:                      "サブモジュールを追加",
			UpdateSubmoduleUrl:                "サブモジュールのURLを更新",
			InitialiseSubmodule:               "サブモジュールを初期化",
			BulkInitialiseSubmodules:          "サブモジュールを一括初期化",
			BulkUpdateSubmodules:              "サブモジュールを一括更新",
			// BulkDeinitialiseSubmodules:        "Bulk deinitialise submodules",
			UpdateSubmodule: "サブモジュールを更新",
			DeleteTag:       "タグを削除",
			PushTag:         "タグをpush",
			// NukeWorkingTree:                   "Nuke working tree",
			// DiscardUnstagedFileChanges:        "Discard unstaged file changes",
			// RemoveUntrackedFiles:              "Remove untracked files",
			SoftReset:           "Softリセット",
			MixedReset:          "Mixedリセット",
			HardReset:           "Hardリセット",
			FastForwardBranch:   "ブランチをfast forward",
			Undo:                "アンドゥ",
			Redo:                "リドゥ",
			CopyPullRequestURL:  "Pull requestのURLをコピー",
			OpenMergeTool:       "マージツールを開く",
			OpenCommitInBrowser: "コミットをブラウザで開く",
			OpenPullRequest:     "Pull requestをブラウザで開く",
			StartBisect:         "Bisectを開始",
			ResetBisect:         "Bisectをリセット",
			BisectSkip:          "Bisectをスキップ",
			BisectMark:          "Bisectをマーク",
		},
		Bisect: Bisect{
			// Mark:                        "Mark %s as %s",
			// MarkStart:                   "Mark %s as %s (start bisect)",
			SkipCurrent:     "%s をスキップする",
			ResetTitle:      "'git bisect' をリセット",
			ResetPrompt:     "'git bisect' をリセットします。よろしいですか?",
			ResetOption:     "Bisectをリセット",
			BisectMenuTitle: "bisect",
			CompleteTitle:   "Bisect完了",
			// CompletePrompt:              "Bisect complete! The following commit introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
			// CompletePromptIndeterminate: "Bisect complete! Some commits were skipped, so any of the following commits may have introduced the change:\n\n%s\n\nDo you want to reset 'git bisect' now?",
		},
	}
}
