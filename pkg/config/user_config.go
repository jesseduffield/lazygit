package config

import (
	"time"
)

type UserConfig struct {
	Gui                  GuiConfig        `yaml:"gui"`
	Git                  GitConfig        `yaml:"git"`
	Update               UpdateConfig     `yaml:"update"`
	Refresher            RefresherConfig  `yaml:"refresher"`
	Reporting            string           `yaml:"reporting"`
	SplashUpdatesIndex   int              `yaml:"splashUpdatesIndex"`
	ConfirmOnQuit        bool             `yaml:"confirmOnQuit"`
	QuitOnTopLevelReturn bool             `yaml:"quitOnTopLevelReturn"`
	Keybinding           KeybindingConfig `yaml:"keybinding"`
	// OS determines what defaults are set for opening files and links
	OS                           OSConfig          `yaml:"os,omitempty"`
	DisableStartupPopups         bool              `yaml:"disableStartupPopups"`
	CustomCommands               []CustomCommand   `yaml:"customCommands"`
	Services                     map[string]string `yaml:"services"`
	NotARepository               string            `yaml:"notARepository"`
	PromptToReturnFromSubprocess bool              `yaml:"promptToReturnFromSubprocess"`
}

type RefresherConfig struct {
	RefreshInterval int `yaml:"refreshInterval"`
	FetchInterval   int `yaml:"fetchInterval"`
}

type GuiConfig struct {
	AuthorColors             map[string]string  `yaml:"authorColors"`
	BranchColors             map[string]string  `yaml:"branchColors"`
	ScrollHeight             int                `yaml:"scrollHeight"`
	ScrollPastBottom         bool               `yaml:"scrollPastBottom"`
	MouseEvents              bool               `yaml:"mouseEvents"`
	SkipUnstageLineWarning   bool               `yaml:"skipUnstageLineWarning"`
	SkipStashWarning         bool               `yaml:"skipStashWarning"`
	SidePanelWidth           float64            `yaml:"sidePanelWidth"`
	ExpandFocusedSidePanel   bool               `yaml:"expandFocusedSidePanel"`
	MainPanelSplitMode       string             `yaml:"mainPanelSplitMode"`
	Language                 string             `yaml:"language"`
	TimeFormat               string             `yaml:"timeFormat"`
	Theme                    ThemeConfig        `yaml:"theme"`
	CommitLength             CommitLengthConfig `yaml:"commitLength"`
	SkipNoStagedFilesWarning bool               `yaml:"skipNoStagedFilesWarning"`
	ShowListFooter           bool               `yaml:"showListFooter"`
	ShowFileTree             bool               `yaml:"showFileTree"`
	ShowRandomTip            bool               `yaml:"showRandomTip"`
	ShowCommandLog           bool               `yaml:"showCommandLog"`
	ShowBottomLine           bool               `yaml:"showBottomLine"`
	ShowIcons                bool               `yaml:"showIcons"`
	CommandLogSize           int                `yaml:"commandLogSize"`
	SplitDiff                string             `yaml:"splitDiff"`
}

type ThemeConfig struct {
	LightTheme                bool     `yaml:"lightTheme"`
	ActiveBorderColor         []string `yaml:"activeBorderColor"`
	InactiveBorderColor       []string `yaml:"inactiveBorderColor"`
	OptionsTextColor          []string `yaml:"optionsTextColor"`
	SelectedLineBgColor       []string `yaml:"selectedLineBgColor"`
	SelectedRangeBgColor      []string `yaml:"selectedRangeBgColor"`
	CherryPickedCommitBgColor []string `yaml:"cherryPickedCommitBgColor"`
	CherryPickedCommitFgColor []string `yaml:"cherryPickedCommitFgColor"`
	UnstagedChangesColor      []string `yaml:"unstagedChangesColor"`
}

type CommitLengthConfig struct {
	Show bool `yaml:"show"`
}

type GitConfig struct {
	Paging              PagingConfig                  `yaml:"paging"`
	Commit              CommitConfig                  `yaml:"commit"`
	Merging             MergingConfig                 `yaml:"merging"`
	SkipHookPrefix      string                        `yaml:"skipHookPrefix"`
	AutoFetch           bool                          `yaml:"autoFetch"`
	AutoRefresh         bool                          `yaml:"autoRefresh"`
	BranchLogCmd        string                        `yaml:"branchLogCmd"`
	AllBranchesLogCmd   string                        `yaml:"allBranchesLogCmd"`
	OverrideGpg         bool                          `yaml:"overrideGpg"`
	DisableForcePushing bool                          `yaml:"disableForcePushing"`
	CommitPrefixes      map[string]CommitPrefixConfig `yaml:"commitPrefixes"`
	// this should really be under 'gui', not 'git'
	ParseEmoji      bool      `yaml:"parseEmoji"`
	Log             LogConfig `yaml:"log"`
	DiffContextSize int       `yaml:"diffContextSize"`
}

type PagingConfig struct {
	ColorArg  string `yaml:"colorArg"`
	Pager     string `yaml:"pager"`
	UseConfig bool   `yaml:"useConfig"`
}

type CommitConfig struct {
	SignOff bool `yaml:"signOff"`
}

type MergingConfig struct {
	ManualCommit bool   `yaml:"manualCommit"`
	Args         string `yaml:"args"`
}

type LogConfig struct {
	Order          string `yaml:"order"`     // one of date-order, author-date-order, topo-order
	ShowGraph      string `yaml:"showGraph"` // one of always, never, when-maximised
	ShowWholeGraph bool   `yaml:"showWholeGraph"`
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

// damn looks like we have some inconsistencies here with -alt and -alt1
type KeybindingUniversalConfig struct {
	Quit                         string   `yaml:"quit"`
	QuitAlt1                     string   `yaml:"quit-alt1"`
	Return                       string   `yaml:"return"`
	QuitWithoutChangingDirectory string   `yaml:"quitWithoutChangingDirectory"`
	TogglePanel                  string   `yaml:"togglePanel"`
	PrevItem                     string   `yaml:"prevItem"`
	NextItem                     string   `yaml:"nextItem"`
	PrevItemAlt                  string   `yaml:"prevItem-alt"`
	NextItemAlt                  string   `yaml:"nextItem-alt"`
	PrevPage                     string   `yaml:"prevPage"`
	NextPage                     string   `yaml:"nextPage"`
	ScrollLeft                   string   `yaml:"scrollLeft"`
	ScrollRight                  string   `yaml:"scrollRight"`
	GotoTop                      string   `yaml:"gotoTop"`
	GotoBottom                   string   `yaml:"gotoBottom"`
	PrevBlock                    string   `yaml:"prevBlock"`
	NextBlock                    string   `yaml:"nextBlock"`
	PrevBlockAlt                 string   `yaml:"prevBlock-alt"`
	NextBlockAlt                 string   `yaml:"nextBlock-alt"`
	NextBlockAlt2                string   `yaml:"nextBlock-alt2"`
	PrevBlockAlt2                string   `yaml:"prevBlock-alt2"`
	JumpToBlock                  []string `yaml:"jumpToBlock"`
	NextMatch                    string   `yaml:"nextMatch"`
	PrevMatch                    string   `yaml:"prevMatch"`
	StartSearch                  string   `yaml:"startSearch"`
	OptionMenu                   string   `yaml:"optionMenu"`
	OptionMenuAlt1               string   `yaml:"optionMenu-alt1"`
	Select                       string   `yaml:"select"`
	GoInto                       string   `yaml:"goInto"`
	Confirm                      string   `yaml:"confirm"`
	ConfirmAlt1                  string   `yaml:"confirm-alt1"`
	Remove                       string   `yaml:"remove"`
	New                          string   `yaml:"new"`
	Edit                         string   `yaml:"edit"`
	OpenFile                     string   `yaml:"openFile"`
	ScrollUpMain                 string   `yaml:"scrollUpMain"`
	ScrollDownMain               string   `yaml:"scrollDownMain"`
	ScrollUpMainAlt1             string   `yaml:"scrollUpMain-alt1"`
	ScrollDownMainAlt1           string   `yaml:"scrollDownMain-alt1"`
	ScrollUpMainAlt2             string   `yaml:"scrollUpMain-alt2"`
	ScrollDownMainAlt2           string   `yaml:"scrollDownMain-alt2"`
	ExecuteCustomCommand         string   `yaml:"executeCustomCommand"`
	CreateRebaseOptionsMenu      string   `yaml:"createRebaseOptionsMenu"`
	PushFiles                    string   `yaml:"pushFiles"`
	PullFiles                    string   `yaml:"pullFiles"`
	Refresh                      string   `yaml:"refresh"`
	CreatePatchOptionsMenu       string   `yaml:"createPatchOptionsMenu"`
	NextTab                      string   `yaml:"nextTab"`
	PrevTab                      string   `yaml:"prevTab"`
	NextScreenMode               string   `yaml:"nextScreenMode"`
	PrevScreenMode               string   `yaml:"prevScreenMode"`
	Undo                         string   `yaml:"undo"`
	Redo                         string   `yaml:"redo"`
	FilteringMenu                string   `yaml:"filteringMenu"`
	DiffingMenu                  string   `yaml:"diffingMenu"`
	DiffingMenuAlt               string   `yaml:"diffingMenu-alt"`
	CopyToClipboard              string   `yaml:"copyToClipboard"`
	OpenRecentRepos              string   `yaml:"openRecentRepos"`
	SubmitEditorText             string   `yaml:"submitEditorText"`
	AppendNewline                string   `yaml:"appendNewline"`
	ExtrasMenu                   string   `yaml:"extrasMenu"`
	ToggleWhitespaceInDiffView   string   `yaml:"toggleWhitespaceInDiffView"`
	IncreaseContextInDiffView    string   `yaml:"increaseContextInDiffView"`
	DecreaseContextInDiffView    string   `yaml:"decreaseContextInDiffView"`
}

type KeybindingStatusConfig struct {
	CheckForUpdate      string `yaml:"checkForUpdate"`
	RecentRepos         string `yaml:"recentRepos"`
	AllBranchesLogGraph string `yaml:"allBranchesLogGraph"`
}

type KeybindingFilesConfig struct {
	CommitChanges            string `yaml:"commitChanges"`
	CommitChangesWithoutHook string `yaml:"commitChangesWithoutHook"`
	AmendLastCommit          string `yaml:"amendLastCommit"`
	CommitChangesWithEditor  string `yaml:"commitChangesWithEditor"`
	IgnoreOrExcludeFile      string `yaml:"IgnoreOrExcludeFile"`
	RefreshFiles             string `yaml:"refreshFiles"`
	StashAllChanges          string `yaml:"stashAllChanges"`
	ViewStashOptions         string `yaml:"viewStashOptions"`
	ToggleStagedAll          string `yaml:"toggleStagedAll"`
	ViewResetOptions         string `yaml:"viewResetOptions"`
	Fetch                    string `yaml:"fetch"`
	ToggleTreeView           string `yaml:"toggleTreeView"`
	OpenMergeTool            string `yaml:"openMergeTool"`
	OpenStatusFilter         string `yaml:"openStatusFilter"`
}

type KeybindingBranchesConfig struct {
	CreatePullRequest      string `yaml:"createPullRequest"`
	ViewPullRequestOptions string `yaml:"viewPullRequestOptions"`
	CopyPullRequestURL     string `yaml:"copyPullRequestURL"`
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
	SquashDown                     string `yaml:"squashDown"`
	RenameCommit                   string `yaml:"renameCommit"`
	RenameCommitWithEditor         string `yaml:"renameCommitWithEditor"`
	ViewResetOptions               string `yaml:"viewResetOptions"`
	MarkCommitAsFixup              string `yaml:"markCommitAsFixup"`
	CreateFixupCommit              string `yaml:"createFixupCommit"`
	SquashAboveCommits             string `yaml:"squashAboveCommits"`
	MoveDownCommit                 string `yaml:"moveDownCommit"`
	MoveUpCommit                   string `yaml:"moveUpCommit"`
	AmendToCommit                  string `yaml:"amendToCommit"`
	ResetCommitAuthor              string `yaml:"resetCommitAuthor"`
	PickCommit                     string `yaml:"pickCommit"`
	RevertCommit                   string `yaml:"revertCommit"`
	CherryPickCopy                 string `yaml:"cherryPickCopy"`
	CherryPickCopyRange            string `yaml:"cherryPickCopyRange"`
	PasteCommits                   string `yaml:"pasteCommits"`
	TagCommit                      string `yaml:"tagCommit"`
	CheckoutCommit                 string `yaml:"checkoutCommit"`
	ResetCherryPick                string `yaml:"resetCherryPick"`
	CopyCommitAttributeToClipboard string `yaml:"copyCommitAttributeToClipboard"`
	OpenLogMenu                    string `yaml:"openLogMenu"`
	OpenInBrowser                  string `yaml:"openInBrowser"`
	ViewBisectOptions              string `yaml:"viewBisectOptions"`
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
	EditSelectHunk      string `yaml:"editSelectHunk"`
}

type KeybindingSubmodulesConfig struct {
	Init     string `yaml:"init"`
	Update   string `yaml:"update"`
	BulkMenu string `yaml:"bulkMenu"`
}

// OSConfig contains config on the level of the os
type OSConfig struct {
	// EditCommand is the command for editing a file
	EditCommand string `yaml:"editCommand,omitempty"`

	// EditCommandTemplate is the command template for editing a file
	EditCommandTemplate string `yaml:"editCommandTemplate,omitempty"`

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
	Stream      bool                  `yaml:"stream"`
}

type CustomCommandPrompt struct {
	Type  string `yaml:"type"` // one of 'input', 'menu', or 'confirm'
	Title string `yaml:"title"`

	// this only apply to input prompts
	InitialValue string `yaml:"initialValue"`

	// this only applies to confirm prompts
	Body string `yaml:"body"`

	// this only applies to menus
	Options []CustomCommandMenuOption

	// this only applies to menuFromCommand
	Command     string `yaml:"command"`
	Filter      string `yaml:"filter"`
	ValueFormat string `yaml:"valueFormat"`
	LabelFormat string `yaml:"labelFormat"`
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
			SkipStashWarning:       false,
			SidePanelWidth:         0.3333,
			ExpandFocusedSidePanel: false,
			MainPanelSplitMode:     "flexible",
			Language:               "auto",
			TimeFormat:             time.RFC822,
			Theme: ThemeConfig{
				LightTheme:                false,
				ActiveBorderColor:         []string{"green", "bold"},
				InactiveBorderColor:       []string{"white"},
				OptionsTextColor:          []string{"blue"},
				SelectedLineBgColor:       []string{"blue"},
				SelectedRangeBgColor:      []string{"blue"},
				CherryPickedCommitBgColor: []string{"cyan"},
				CherryPickedCommitFgColor: []string{"blue"},
				UnstagedChangesColor:      []string{"red"},
			},
			CommitLength:             CommitLengthConfig{Show: true},
			SkipNoStagedFilesWarning: false,
			ShowListFooter:           true,
			ShowCommandLog:           true,
			ShowBottomLine:           true,
			ShowFileTree:             true,
			ShowRandomTip:            true,
			ShowIcons:                false,
			CommandLogSize:           8,
			SplitDiff:                "auto",
		},
		Git: GitConfig{
			Paging: PagingConfig{
				ColorArg:  "always",
				Pager:     "",
				UseConfig: false,
			},
			Commit: CommitConfig{
				SignOff: false,
			},
			Merging: MergingConfig{
				ManualCommit: false,
				Args:         "",
			},
			Log: LogConfig{
				Order:          "topo-order",
				ShowGraph:      "when-maximised",
				ShowWholeGraph: false,
			},
			SkipHookPrefix:      "WIP",
			AutoFetch:           true,
			AutoRefresh:         true,
			BranchLogCmd:        "git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --",
			AllBranchesLogCmd:   "git log --graph --all --color=always --abbrev-commit --decorate --date=relative  --pretty=medium",
			DisableForcePushing: false,
			CommitPrefixes:      map[string]CommitPrefixConfig(nil),
			ParseEmoji:          false,
			DiffContextSize:     3,
		},
		Refresher: RefresherConfig{
			RefreshInterval: 10,
			FetchInterval:   60,
		},
		Update: UpdateConfig{
			Method: "prompt",
			Days:   14,
		},
		Reporting:            "undetermined",
		SplashUpdatesIndex:   0,
		ConfirmOnQuit:        false,
		QuitOnTopLevelReturn: false,
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
				ScrollLeft:                   "H",
				ScrollRight:                  "L",
				GotoTop:                      "<",
				GotoBottom:                   ">",
				PrevBlock:                    "<left>",
				NextBlock:                    "<right>",
				PrevBlockAlt:                 "h",
				NextBlockAlt:                 "l",
				PrevBlockAlt2:                "<backtab>",
				NextBlockAlt2:                "<tab>",
				JumpToBlock:                  []string{"1", "2", "3", "4", "5"},
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
				OpenRecentRepos:              "<c-r>",
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
				AppendNewline:                "<a-enter>",
				ExtrasMenu:                   "@",
				ToggleWhitespaceInDiffView:   "<c-w>",
				IncreaseContextInDiffView:    "}",
				DecreaseContextInDiffView:    "{",
			},
			Status: KeybindingStatusConfig{
				CheckForUpdate:      "u",
				RecentRepos:         "<enter>",
				AllBranchesLogGraph: "a",
			},
			Files: KeybindingFilesConfig{
				CommitChanges:            "c",
				CommitChangesWithoutHook: "w",
				AmendLastCommit:          "A",
				CommitChangesWithEditor:  "C",
				IgnoreOrExcludeFile:      "i",
				RefreshFiles:             "r",
				StashAllChanges:          "s",
				ViewStashOptions:         "S",
				ToggleStagedAll:          "a",
				ViewResetOptions:         "D",
				Fetch:                    "f",
				ToggleTreeView:           "`",
				OpenMergeTool:            "M",
				OpenStatusFilter:         "<c-b>",
			},
			Branches: KeybindingBranchesConfig{
				CopyPullRequestURL:     "<c-y>",
				CreatePullRequest:      "o",
				ViewPullRequestOptions: "O",
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
				SquashDown:                     "s",
				RenameCommit:                   "r",
				RenameCommitWithEditor:         "R",
				ViewResetOptions:               "g",
				MarkCommitAsFixup:              "f",
				CreateFixupCommit:              "F",
				SquashAboveCommits:             "S",
				MoveDownCommit:                 "<c-j>",
				MoveUpCommit:                   "<c-k>",
				AmendToCommit:                  "A",
				ResetCommitAuthor:              "a",
				PickCommit:                     "p",
				RevertCommit:                   "t",
				CherryPickCopy:                 "c",
				CherryPickCopyRange:            "C",
				PasteCommits:                   "v",
				TagCommit:                      "T",
				CheckoutCommit:                 "<space>",
				ResetCherryPick:                "<c-R>",
				CopyCommitAttributeToClipboard: "y",
				OpenLogMenu:                    "<c-l>",
				OpenInBrowser:                  "o",
				ViewBisectOptions:              "b",
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
				EditSelectHunk:      "E",
			},
			Submodules: KeybindingSubmodulesConfig{
				Init:     "i",
				Update:   "u",
				BulkMenu: "b",
			},
		},
		OS:                           GetPlatformDefaultConfig(),
		DisableStartupPopups:         false,
		CustomCommands:               []CustomCommand(nil),
		Services:                     map[string]string(nil),
		NotARepository:               "prompt",
		PromptToReturnFromSubprocess: true,
	}
}
