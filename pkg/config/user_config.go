package config

type CommitPrefixConfig struct {
	Pattern string `yaml:"pattern"`
	Replace string `yaml:"replace"`
}

type ThemeConfig struct {
	LightTheme           bool     `yaml:"lightTheme"`
	ActiveBorderColor    []string `yaml:"activeBorderColor"`
	InactiveBorderColor  []string `yaml:"inactiveBorderColor"`
	OptionsTextColor     []string `yaml:"optionsTextColor"`
	SelectedLineBgColor  []string `yaml:"selectedLineBgColor"`
	SelectedRangeBgColor []string `yaml:"selectedRangeBgColor"`
}

type CustomCommandMenuOption struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}

type CustomCommandPrompt struct {
	Type  string `yaml:"type"` // one of 'input' and 'menu'
	Title string `yaml:"title"`

	// this only apply to prompts
	InitialValue string `yaml:"initialValue"`

	// this only applies to menus
	Options []CustomCommandMenuOption
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

type UserConfig struct {
	Gui struct {
		ScrollHeight           int         `yaml:"scrollHeight"`
		ScrollPastBottom       bool        `yaml:"scrollPastBottom"`
		MouseEvents            bool        `yaml:"mouseEvents"`
		SkipUnstageLineWarning bool        `yaml:"skipUnstageLineWarning"`
		SkipStashWarning       bool        `yaml:"skipStashWarning"`
		SidePanelWidth         float64     `yaml:"sidePanelWidth"`
		ExpandFocusedSidePanel bool        `yaml:"expandFocusedSidePanel"`
		MainPanelSplitMode     string      `yaml:"mainPanelSplitMode"`
		Theme                  ThemeConfig `yaml:"theme"`
		CommitLength           struct {
			Show bool `yaml:"show"`
		} `yaml:"commitLength"`
	} `yaml:"gui"`
	Git struct {
		Paging struct {
			ColorArg  string `yaml:"colorArg"`
			Pager     string `yaml:"pager"`
			UseConfig bool   `yaml:"useConfig"`
		} `yaml:"paging"`
		Merging struct {
			ManualCommit bool   `yaml:"manualCommit"`
			Args         string `yaml:"args"`
		} `yaml:"merging"`
		Pull struct {
			Mode string `yaml:"mode"`
		} `yaml:"pull"`
		SkipHookPrefix      string                        `yaml:"skipHookPrefix"`
		AutoFetch           bool                          `yaml:"autoFetch"`
		BranchLogCmd        string                        `yaml:"branchLogCmd"`
		OverrideGpg         bool                          `yaml:"overrideGpg"`
		DisableForcePushing bool                          `yaml:"disableForcePushing"`
		CommitPrefixes      map[string]CommitPrefixConfig `yaml:"commitPrefixes"`
	} `yaml:"git"`
	Update struct {
		Method string `yaml:"method"`
		Days   int64  `yaml:"days"`
	} `yaml:"update"`
	Reporting            string `yaml:"reporting"`
	SplashUpdatesIndex   int    `yaml:"splashUpdatesIndex"`
	ConfirmOnQuit        bool   `yaml:"confirmOnQuit"`
	QuitOnTopLevelReturn bool   `yaml:"quitOnTopLevelReturn"`
	Keybinding           struct {
		Universal struct {
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
		} `yaml:"universal"`
		Status struct {
			CheckForUpdate string `yaml:"checkForUpdate"`
			RecentRepos    string `yaml:"recentRepos"`
		} `yaml:"status"`
		Files struct {
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
		} `yaml:"files"`
		Branches struct {
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
		} `yaml:"branches"`
		Commits struct {
			SquashDown             string `yaml:"squashDown"`
			RenameCommit           string `yaml:"renameCommit"`
			RenameCommitWithEditor string `yaml:"renameCommitWithEditor"`
			ViewResetOptions       string `yaml:"viewResetOptions"`
			MarkCommitAsFixup      string `yaml:"markCommitAsFixup"`
			CreateFixupCommit      string `yaml:"createFixupCommit"`
			SquashAboveCommits     string `yaml:"squashAboveCommits"`
			MoveDownCommit         string `yaml:"moveDownCommit"`
			MoveUpCommit           string `yaml:"moveUpCommit"`
			AmendToCommit          string `yaml:"amendToCommit"`
			PickCommit             string `yaml:"pickCommit"`
			RevertCommit           string `yaml:"revertCommit"`
			CherryPickCopy         string `yaml:"cherryPickCopy"`
			CherryPickCopyRange    string `yaml:"cherryPickCopyRange"`
			PasteCommits           string `yaml:"pasteCommits"`
			TagCommit              string `yaml:"tagCommit"`
			CheckoutCommit         string `yaml:"checkoutCommit"`
			ResetCherryPick        string `yaml:"resetCherryPick"`
		} `yaml:"commits"`
		Stash struct {
			PopStash string `yaml:"popStash"`
		} `yaml:"stash"`
		CommitFiles struct {
			CheckoutCommitFile string `yaml:"checkoutCommitFile"`
		} `yaml:"commitFiles"`
		Main struct {
			ToggleDragSelect    string `yaml:"toggleDragSelect"`
			ToggleDragSelectAlt string `yaml:"toggleDragSelect-alt"`
			ToggleSelectHunk    string `yaml:"toggleSelectHunk"`
			PickBothHunks       string `yaml:"pickBothHunks"`
		} `yaml:"main"`
		Submodules struct {
			Init     string `yaml:"init"`
			Update   string `yaml:"update"`
			BulkMenu string `yaml:"bulkMenu"`
		} `yaml:"submodules"`
	} `yaml:"keybinding"`
	// OS determines what defaults are set for opening files and links
	OS                   OSConfig          `yaml:"os,omitempty"`
	DisableStartupPopups bool              `yaml:"disableStartupPopups"`
	CustomCommands       []CustomCommand   `yaml:"customCommands"`
	Services             map[string]string `yaml:"services"`
}
