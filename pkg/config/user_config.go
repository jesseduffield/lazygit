package config

type UserConfig struct {
	Gui                  GuiConfig        `yaml:"gui"`
	Git                  GitConfig        `yaml:"git"`
	Update               UpdateConfig     `yaml:"update"`
	Reporting            string           `yaml:"reporting"`
	SplashUpdatesIndex   int              `yaml:"splashUpdatesIndex"`
	ConfirmOnQuit        bool             `yaml:"confirmOnQuit"`
	QuitOnTopLevelReturn bool             `yaml:"quitOnTopLevelReturn"`
	Keybinding           KeybindingConfig `yaml:"keybinding"`
	// OS determines what defaults are set for opening files and links
	OS                   OSConfig          `yaml:"os,omitempty"`
	DisableStartupPopups bool              `yaml:"disableStartupPopups"`
	CustomCommands       []CustomCommand   `yaml:"customCommands"`
	Services             map[string]string `yaml:"services"`
}

type GuiConfig struct {
	ScrollHeight           int                `yaml:"scrollHeight"`
	ScrollPastBottom       bool               `yaml:"scrollPastBottom"`
	MouseEvents            bool               `yaml:"mouseEvents"`
	SkipUnstageLineWarning bool               `yaml:"skipUnstageLineWarning"`
	SkipStashWarning       bool               `yaml:"skipStashWarning"`
	SidePanelWidth         float64            `yaml:"sidePanelWidth"`
	ExpandFocusedSidePanel bool               `yaml:"expandFocusedSidePanel"`
	MainPanelSplitMode     string             `yaml:"mainPanelSplitMode"`
	Theme                  ThemeConfig        `yaml:"theme"`
	CommitLength           CommitLengthConfig `yaml:"commitLength"`
}

type ThemeConfig struct {
	LightTheme           bool     `yaml:"lightTheme"`
	ActiveBorderColor    []string `yaml:"activeBorderColor"`
	InactiveBorderColor  []string `yaml:"inactiveBorderColor"`
	OptionsTextColor     []string `yaml:"optionsTextColor"`
	SelectedLineBgColor  []string `yaml:"selectedLineBgColor"`
	SelectedRangeBgColor []string `yaml:"selectedRangeBgColor"`
}

type CommitLengthConfig struct {
	Show bool `yaml:"show"`
}

type GitConfig struct {
	Paging              PagingConfig                  `yaml:"paging"`
	Merging             MergingConfig                 `yaml:"merging"`
	Pull                PullConfig                    `yaml:"pull"`
	SkipHookPrefix      string                        `yaml:"skipHookPrefix"`
	AutoFetch           bool                          `yaml:"autoFetch"`
	BranchLogCmd        string                        `yaml:"branchLogCmd"`
	OverrideGpg         bool                          `yaml:"overrideGpg"`
	DisableForcePushing bool                          `yaml:"disableForcePushing"`
	CommitPrefixes      map[string]CommitPrefixConfig `yaml:"commitPrefixes"`
}

type PagingConfig struct {
	ColorArg  string `yaml:"colorArg"`
	Pager     string `yaml:"pager"`
	UseConfig bool   `yaml:"useConfig"`
}

type MergingConfig struct {
	ManualCommit bool   `yaml:"manualCommit"`
	Args         string `yaml:"args"`
}

type PullConfig struct {
	Mode string `yaml:"mode"`
}

type CommitPrefixConfig struct {
	Pattern string `yaml:"pattern"`
	Replace string `yaml:"replace"`
}

type UpdateConfig struct {
	Method string `yaml:"method"`
	Days   int64  `yaml:"days"`
}

type KeybindingConfig struct {
	Universal   KeybindingUniversalConfig   `yaml:"universal"`
	Status      KeybindingStatusConfig      `yaml:"status"`
	Files       KeybindingFilesConfig       `yaml:"files"`
	Branches    KeybindingBranchesConfig    `yaml:"branches"`
	Commits     KeybindingCommitsConfig     `yaml:"commits"`
	Stash       KeybindingStashConfig       `yaml:"stash"`
	CommitFiles KeybindingCommitFilesConfig `yaml:"commitFiles"`
	Main        KeybindingMainConfig        `yaml:"main"`
	Submodules  KeybindingSubmodulesConfig  `yaml:"submodules"`
}

type KeybindingUniversalConfig struct {
	Quit                         string `yaml:"quit"`
	QuitAlt1                     string `yaml:"quit-alt1"`
	Return                       string `yaml:"return"`
	QuitWithoutChangingDirectory string `yaml:"quitWithoutChangingDirectory"`
	TogglePanel                  string `yaml:"togglePanel"`
	PrevItem                     string `yaml:"prevItem"`
	NextItem                     string `yaml:"nextItem"`
	PrevItemAlt                  string `yaml:"prevItem-alt"`
	NextItemAlt                  string `yaml:"nextItem-alt"`
	PrevPage                     string `yaml:"prevPage"`
	NextPage                     string `yaml:"nextPage"`
	GotoTop                      string `yaml:"gotoTop"`
	GotoBottom                   string `yaml:"gotoBottom"`
	PrevBlock                    string `yaml:"prevBlock"`
	NextBlock                    string `yaml:"nextBlock"`
	PrevBlockAlt                 string `yaml:"prevBlock-alt"`
	NextBlockAlt                 string `yaml:"nextBlock-alt"`
	NextMatch                    string `yaml:"nextMatch"`
	PrevMatch                    string `yaml:"prevMatch"`
	StartSearch                  string `yaml:"startSearch"`
	OptionMenu                   string `yaml:"optionMenu"`
	OptionMenuAlt1               string `yaml:"optionMenu-alt1"`
	Select                       string `yaml:"select"`
	GoInto                       string `yaml:"goInto"`
	Confirm                      string `yaml:"confirm"`
	ConfirmAlt1                  string `yaml:"confirm-alt1"`
	Remove                       string `yaml:"remove"`
	New                          string `yaml:"new"`
	Edit                         string `yaml:"edit"`
	OpenFile                     string `yaml:"openFile"`
	ScrollUpMain                 string `yaml:"scrollUpMain"`
	ScrollDownMain               string `yaml:"scrollDownMain"`
	ScrollUpMainAlt1             string `yaml:"scrollUpMain-alt1"`
	ScrollDownMainAlt1           string `yaml:"scrollDownMain-alt1"`
	ScrollUpMainAlt2             string `yaml:"scrollUpMain-alt2"`
	ScrollDownMainAlt2           string `yaml:"scrollDownMain-alt2"`
	ExecuteCustomCommand         string `yaml:"executeCustomCommand"`
	CreateRebaseOptionsMenu      string `yaml:"createRebaseOptionsMenu"`
	PushFiles                    string `yaml:"pushFiles"`
	PullFiles                    string `yaml:"pullFiles"`
	Refresh                      string `yaml:"refresh"`
	CreatePatchOptionsMenu       string `yaml:"createPatchOptionsMenu"`
	NextTab                      string `yaml:"nextTab"`
	PrevTab                      string `yaml:"prevTab"`
	NextScreenMode               string `yaml:"nextScreenMode"`
	PrevScreenMode               string `yaml:"prevScreenMode"`
	Undo                         string `yaml:"undo"`
	Redo                         string `yaml:"redo"`
	FilteringMenu                string `yaml:"filteringMenu"`
	DiffingMenu                  string `yaml:"diffingMenu"`
	DiffingMenuAlt               string `yaml:"diffingMenu-alt"`
	CopyToClipboard              string `yaml:"copyToClipboard"`
	SubmitEditorText             string `yaml:"submitEditorText"`
	AppendNewline                string `yaml:"appendNewline"`
}

type KeybindingStatusConfig struct {
	CheckForUpdate string `yaml:"checkForUpdate"`
	RecentRepos    string `yaml:"recentRepos"`
}

type KeybindingFilesConfig struct {
	CommitChanges            string `yaml:"commitChanges"`
	CommitChangesWithoutHook string `yaml:"commitChangesWithoutHook"`
	AmendLastCommit          string `yaml:"amendLastCommit"`
	CommitChangesWithEditor  string `yaml:"commitChangesWithEditor"`
	IgnoreFile               string `yaml:"ignoreFile"`
	RefreshFiles             string `yaml:"refreshFiles"`
	StashAllChanges          string `yaml:"stashAllChanges"`
	ViewStashOptions         string `yaml:"viewStashOptions"`
	ToggleStagedAll          string `yaml:"toggleStagedAll"`
	ViewResetOptions         string `yaml:"viewResetOptions"`
	Fetch                    string `yaml:"fetch"`
}

type KeybindingBranchesConfig struct {
	CreatePullRequest      string `yaml:"createPullRequest"`
	CheckoutBranchByName   string `yaml:"checkoutBranchByName"`
	ForceCheckoutBranch    string `yaml:"forceCheckoutBranch"`
	RebaseBranch           string `yaml:"rebaseBranch"`
	RenameBranch           string `yaml:"renameBranch"`
	MergeIntoCurrentBranch string `yaml:"mergeIntoCurrentBranch"`
	ViewGitFlowOptions     string `yaml:"viewGitFlowOptions"`
	FastForward            string `yaml:"fastForward"`
	PushTag                string `yaml:"pushTag"`
	SetUpstream            string `yaml:"setUpstream"`
	FetchRemote            string `yaml:"fetchRemote"`
}

type KeybindingCommitsConfig struct {
	SquashDown                   string `yaml:"squashDown"`
	RenameCommit                 string `yaml:"renameCommit"`
	RenameCommitWithEditor       string `yaml:"renameCommitWithEditor"`
	ViewResetOptions             string `yaml:"viewResetOptions"`
	MarkCommitAsFixup            string `yaml:"markCommitAsFixup"`
	CreateFixupCommit            string `yaml:"createFixupCommit"`
	SquashAboveCommits           string `yaml:"squashAboveCommits"`
	MoveDownCommit               string `yaml:"moveDownCommit"`
	MoveUpCommit                 string `yaml:"moveUpCommit"`
	AmendToCommit                string `yaml:"amendToCommit"`
	PickCommit                   string `yaml:"pickCommit"`
	RevertCommit                 string `yaml:"revertCommit"`
	CherryPickCopy               string `yaml:"cherryPickCopy"`
	CherryPickCopyRange          string `yaml:"cherryPickCopyRange"`
	PasteCommits                 string `yaml:"pasteCommits"`
	TagCommit                    string `yaml:"tagCommit"`
	CheckoutCommit               string `yaml:"checkoutCommit"`
	ResetCherryPick              string `yaml:"resetCherryPick"`
	CopyCommitMessageToClipboard string `yaml:"copyCommitMessageToClipboard"`
}

type KeybindingStashConfig struct {
	PopStash string `yaml:"popStash"`
}

type KeybindingCommitFilesConfig struct {
	CheckoutCommitFile string `yaml:"checkoutCommitFile"`
}

type KeybindingMainConfig struct {
	ToggleDragSelect    string `yaml:"toggleDragSelect"`
	ToggleDragSelectAlt string `yaml:"toggleDragSelect-alt"`
	ToggleSelectHunk    string `yaml:"toggleSelectHunk"`
	PickBothHunks       string `yaml:"pickBothHunks"`
}

type KeybindingSubmodulesConfig struct {
	Init     string `yaml:"init"`
	Update   string `yaml:"update"`
	BulkMenu string `yaml:"bulkMenu"`
}

// OSConfig contains config on the level of the os
type OSConfig struct {
	// OpenCommand is the command for opening a file
	OpenCommand string `yaml:"openCommand,omitempty"`

	// OpenCommand is the command for opening a link
	OpenLinkCommand string `yaml:"openLinkCommand,omitempty"`
}

type CustomCommand struct {
	Key         string                `yaml:"key"`
	Context     string                `yaml:"context"`
	Command     string                `yaml:"command"`
	Subprocess  bool                  `yaml:"subprocess"`
	Prompts     []CustomCommandPrompt `yaml:"prompts"`
	LoadingText string                `yaml:"loadingText"`
	Description string                `yaml:"description"`
}

type CustomCommandPrompt struct {
	Type  string `yaml:"type"` // one of 'input' and 'menu'
	Title string `yaml:"title"`

	// this only apply to prompts
	InitialValue string `yaml:"initialValue"`

	// this only applies to menus
	Options []CustomCommandMenuOption
}

type CustomCommandMenuOption struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}

func GetDefaultConfig() *UserConfig {
	return &UserConfig{
		Gui: GuiConfig{
			ScrollHeight:           2,
			ScrollPastBottom:       true,
			MouseEvents:            true,
			SkipUnstageLineWarning: false,
			SkipStashWarning:       true,
			SidePanelWidth:         0.3333,
			ExpandFocusedSidePanel: false,
			MainPanelSplitMode:     "flexible",
			Theme: ThemeConfig{
				LightTheme:           false,
				ActiveBorderColor:    []string{"green", "bold"},
				InactiveBorderColor:  []string{"white"},
				OptionsTextColor:     []string{"blue"},
				SelectedLineBgColor:  []string{"default"},
				SelectedRangeBgColor: []string{"blue"},
			},
			CommitLength: CommitLengthConfig{Show: true},
		},
		Git: GitConfig{
			Paging: PagingConfig{
				ColorArg:  "always",
				Pager:     "",
				UseConfig: false},
			Merging: MergingConfig{
				ManualCommit: false,
				Args:         "",
			},
			Pull: PullConfig{
				Mode: "merge",
			},
			SkipHookPrefix:      "WIP",
			AutoFetch:           true,
			BranchLogCmd:        "git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --",
			OverrideGpg:         false,
			DisableForcePushing: false,
			CommitPrefixes:      map[string]CommitPrefixConfig(nil),
		},
		Update: UpdateConfig{
			Method: "prompt",
			Days:   14,
		},
		Reporting:            "undetermined",
		SplashUpdatesIndex:   0,
		ConfirmOnQuit:        false,
		QuitOnTopLevelReturn: true,
		Keybinding: KeybindingConfig{
			Universal: KeybindingUniversalConfig{
				Quit:                         "q",
				QuitAlt1:                     "<c-c>",
				Return:                       "<esc>",
				QuitWithoutChangingDirectory: "Q",
				TogglePanel:                  "<tab>",
				PrevItem:                     "<up>",
				NextItem:                     "<down>",
				PrevItemAlt:                  "k",
				NextItemAlt:                  "j",
				PrevPage:                     ",",
				NextPage:                     ".",
				GotoTop:                      "<",
				GotoBottom:                   ">",
				PrevBlock:                    "<left>",
				NextBlock:                    "<right>",
				PrevBlockAlt:                 "h",
				NextBlockAlt:                 "l",
				NextMatch:                    "n",
				PrevMatch:                    "N",
				StartSearch:                  "/",
				OptionMenu:                   "x",
				OptionMenuAlt1:               "?",
				Select:                       "<space>",
				GoInto:                       "<enter>",
				Confirm:                      "<enter>",
				ConfirmAlt1:                  "y",
				Remove:                       "d",
				New:                          "n",
				Edit:                         "e",
				OpenFile:                     "o",
				ScrollUpMain:                 "<pgup>",
				ScrollDownMain:               "<pgdown>",
				ScrollUpMainAlt1:             "K",
				ScrollDownMainAlt1:           "J",
				ScrollUpMainAlt2:             "<c-u>",
				ScrollDownMainAlt2:           "<c-d>",
				ExecuteCustomCommand:         ":",
				CreateRebaseOptionsMenu:      "m",
				PushFiles:                    "P",
				PullFiles:                    "p",
				Refresh:                      "R",
				CreatePatchOptionsMenu:       "<c-p>",
				NextTab:                      "]",
				PrevTab:                      "[",
				NextScreenMode:               "+",
				PrevScreenMode:               "_",
				Undo:                         "z",
				Redo:                         "<c-z>",
				FilteringMenu:                "<c-s>",
				DiffingMenu:                  "W",
				DiffingMenuAlt:               "<c-e>",
				CopyToClipboard:              "<c-o>",
				SubmitEditorText:             "<enter>",
				AppendNewline:                "<tab>",
			},
			Status: KeybindingStatusConfig{
				CheckForUpdate: "u",
				RecentRepos:    "<enter>",
			},
			Files: KeybindingFilesConfig{
				CommitChanges:            "c",
				CommitChangesWithoutHook: "w",
				AmendLastCommit:          "A",
				CommitChangesWithEditor:  "C",
				IgnoreFile:               "i",
				RefreshFiles:             "r",
				StashAllChanges:          "s",
				ViewStashOptions:         "S",
				ToggleStagedAll:          "a",
				ViewResetOptions:         "D",
				Fetch:                    "f",
			},
			Branches: KeybindingBranchesConfig{
				CreatePullRequest:      "o",
				CheckoutBranchByName:   "c",
				ForceCheckoutBranch:    "F",
				RebaseBranch:           "r",
				RenameBranch:           "R",
				MergeIntoCurrentBranch: "M",
				ViewGitFlowOptions:     "i",
				FastForward:            "f",
				PushTag:                "P",
				SetUpstream:            "u",
				FetchRemote:            "f",
			},
			Commits: KeybindingCommitsConfig{
				SquashDown:                   "s",
				RenameCommit:                 "r",
				RenameCommitWithEditor:       "R",
				ViewResetOptions:             "g",
				MarkCommitAsFixup:            "f",
				CreateFixupCommit:            "F",
				SquashAboveCommits:           "S",
				MoveDownCommit:               "<c-j>",
				MoveUpCommit:                 "<c-k>",
				AmendToCommit:                "A",
				PickCommit:                   "p",
				RevertCommit:                 "t",
				CherryPickCopy:               "c",
				CherryPickCopyRange:          "C",
				PasteCommits:                 "v",
				TagCommit:                    "T",
				CheckoutCommit:               "<space>",
				ResetCherryPick:              "<c-R>",
				CopyCommitMessageToClipboard: "<c-y>",
			},
			Stash: KeybindingStashConfig{
				PopStash: "g",
			},
			CommitFiles: KeybindingCommitFilesConfig{
				CheckoutCommitFile: "c",
			},
			Main: KeybindingMainConfig{
				ToggleDragSelect:    "v",
				ToggleDragSelectAlt: "V",
				ToggleSelectHunk:    "a",
				PickBothHunks:       "b",
			},
			Submodules: KeybindingSubmodulesConfig{
				Init:     "i",
				Update:   "u",
				BulkMenu: "b",
			},
		},
		OS:                   GetPlatformDefaultConfig(),
		DisableStartupPopups: false,
		CustomCommands:       []CustomCommand(nil),
		Services:             map[string]string(nil),
	}
}
