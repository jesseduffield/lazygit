package i18n

func dutchTranslationSet() TranslationSet {
	return TranslationSet{
		NotEnoughSpace:                      "Niet genoeg ruimte om de panelen te renderen",
		DiffTitle:                           "Diff",
		LogTitle:                            "Log",
		FilesTitle:                          "Bestanden",
		BranchesTitle:                       "Branches",
		CommitsTitle:                        "Commits",
		StashTitle:                          "Stash",
		UnstagedChanges:                     `Unstaged wijzigingen`,
		StagedChanges:                       `Staged wijzigingen`,
		PatchBuildingMainTitle:              `Voeg lijnen/hunks toe aan patch`,
		MergingMainTitle:                    "Los merge conflicten op",
		MainTitle:                           "Hoofd",
		StagingTitle:                        "Staging",
		NormalTitle:                         "Normaal",
		CommitMessage:                       "Commitbericht",
		CredentialsUsername:                 "Gebruikersnaam",
		CredentialsPassword:                 "Wachtwoord",
		CredentialsPassphrase:               "Voer een wachtwoordzin in voor de SSH-sleutel",
		PassUnameWrong:                      "Wachtwoord en/of gebruikersnaam verkeerd",
		CommitChanges:                       "Commit veranderingen",
		AmendLastCommit:                     "wijzig laatste commit",
		SureToAmend:                         "Weet je zeker dat je de laatste commit wilt wijzigen? U kunt het commit-bericht wijzigen vanuit het commits-paneel.",
		NoCommitToAmend:                     "Er is geen commits om te wijzigen.",
		CommitChangesWithEditor:             "commit veranderingen met de git editor",
		StatusTitle:                         "Status",
		LcNavigate:                          "navigeer",
		LcMenu:                              "menu",
		LcExecute:                           "uitvoeren",
		LcOpen:                              "open",
		LcIgnore:                            "negeren",
		LcDelete:                            "verwijderen",
		LcToggleStaged:                      "toggle staged",
		LcToggleStagedAll:                   "toggle staged alle",
		LcRefresh:                           "verversen",
		LcPush:                              "push",
		LcPull:                              "pull",
		LcEdit:                              "bewerken",
		LcScroll:                            "scroll",
		LcAbortMerge:                        "samenvoegen afbreken",
		LcResolveMergeConflicts:             "los merge conflicten op",
		LcCommitFileFilter:                  "Commit dossiers filteren",
		FilterModifiedFiles:                 "Show only modified files",
		FilterUntrackedFiles:                "Show only untracked files",
		ResetCommitFilterState:              "Reset commit file state filter",
		MergeConflictsTitle:                 "Merge Conflicten",
		LcCheckout:                          "uitchecken",
		FileHasNoUnstagedChanges:            "Het bestand heeft geen unstaged veranderingen om toe te voegen",
		CannotGitAdd:                        "Kan commando niet uitvoeren git add --path untracked files",
		NoFilesDisplay:                      "Geen bestanden om te laten zien",
		NotAFile:                            "Dit is geen bestand",
		PullWait:                            "Pullen...",
		PushWait:                            "Pushen...",
		FetchWait:                           "Fetchen...",
		FileNoMergeCons:                     "Dit bestand heeft geen merge conflicten",
		LcSoftReset:                         "zacht reset",
		SureTo:                              "Weet je het zeker dat je {{.fileName}} wilt {{.deleteVerb}} (je veranderingen zullen worden verwijderd)",
		AlreadyCheckedOutBranch:             "Je hebt deze branch al uitgecheckt",
		SureForceCheckout:                   "Weet je zeker dat je het uitchecken wil forceren? Al je lokale verandering zullen worden verwijdert",
		ForceCheckoutBranch:                 "Forceer uitchecken op deze branch",
		BranchName:                          "Branch naam",
		NewBranchNameBranchOff:              "Nieuw branch naam (Branch is afgeleid van {{.branchName}})",
		CantDeleteCheckOutBranch:            "Je kan een uitgecheckte branch niet verwijderen!",
		DeleteBranch:                        "Verwijder branch",
		DeleteBranchMessage:                 "Weet je zeker dat je branch {{.selectedBranchName}} wilt verwijderen?",
		ForceDeleteBranchMessage:            "Weet je zeker dat je branch {{.selectedBranchName}} geforceerd wil verwijderen?",
		LcRebaseBranch:                      "rebase branch",
		CantRebaseOntoSelf:                  "Je kan niet een branch rebasen op zichzelf",
		CantMergeBranchIntoItself:           "Je kan niet een branch in zichzelf mergen",
		LcForceCheckout:                     "forceer checkout",
		LcMerge:                             "samenvoegen",
		LcCheckoutByName:                    "uitchecken bij naam",
		LcNewBranch:                         "nieuwe branch",
		LcDeleteBranch:                      "verwijder branch",
		LcForceDeleteBranch:                 "verwijder branch (forceer)",
		NoBranchesThisRepo:                  "Geen branches voor deze repo",
		NoTrackingThisBranch:                "deze branch wordt niet gevolgd",
		CommitMessageConfirm:                "{{.keyBindClose}}: Sluiten, {{.keyBindNewLine}}: Nieuwe lijn, {{.keyBindConfirm}}: Bevestig",
		CommitWithoutMessageErr:             "Je kan geen commit maken zonder commit bericht",
		CloseConfirm:                        "{{.keyBindClose}}: Sluiten, {{.keyBindConfirm}}: Bevestig",
		LcClose:                             "sluiten",
		LcQuit:                              "quit",
		SureResetThisCommit:                 "Weet je het zeker dat je wil resetten naar deze commit?",
		ResetToCommit:                       "Reset Naar Commit",
		LcSquashDown:                        "squash beneden",
		LcRename:                            "hernoemen",
		LcResetToThisCommit:                 "reset naar deze commit",
		LcFixupCommit:                       "Fixup commit",
		OnlySquashTopmostCommit:             "Kan alleen bovenste commit squashen",
		YouNoCommitsToSquash:                "Je hebt geen commits om mee te squashen",
		CantFixupWhileUnstagedChanges:       "Kan geen Fixup uitvoeren op unstaged veranderingen",
		Fixup:                               "Fixup",
		SureFixupThisCommit:                 "Weet je zeker dat je fixup wil uitvoeren op deze commit? De commit hieronder zol worden squashed in deze",
		SureSquashThisCommit:                "Weet je zeker dat je deze commit wil samenvoegen met de commit hieronder?",
		Squash:                              "Squash",
		LcPickCommit:                        "kies commit (wanneer midden in rebase)",
		LcRevertCommit:                      "commit ongedaan maken",
		OnlyRenameTopCommit:                 "Je kan alleen de bovenste commit hernoemen",
		LcRenameCommit:                      "hernoem commit",
		LcDeleteCommit:                      "verwijder commit",
		LcMoveDownCommit:                    "verplaats commit 1 naar beneden",
		LcMoveUpCommit:                      "verplaats commit 1 naar boven",
		LcEditCommit:                        "wijzig commit",
		LcAmendToCommit:                     "wijzig commit met staged veranderingen",
		LcRenameCommitEditor:                "hernoem commit met editor",
		PotentialErrInGetselectedCommit:     "Er is mogelijk een error in getSelected Commit (geen match tussen ui en state)",
		NoCommitsThisBranch:                 "Geen commits in deze branch",
		Error:                               "Foutmelding",
		RunningSubprocess:                   "subprocess lopend",
		LcSelectHunk:                        "selecteer stuk",
		LcNavigateConflicts:                 "navigeer conflicts",
		LcPickHunk:                          "kies stuk",
		LcPickAllHunks:                      "kies beide stukken",
		LcUndo:                              "ongedaan maken",
		LcUndoReflog:                        "ongedaan maken (via reflog) (experimenteel)",
		LcRedoReflog:                        "redo (via reflog) (experimenteel)",
		LcPop:                               "pop",
		LcDrop:                              "laten vallen",
		LcApply:                             "toepassen",
		NoStashEntries:                      "Geen stash items",
		StashDrop:                           "Stash laten vallen",
		SureDropStashEntry:                  "Weet je het zeker dat je deze stash entry wil laten vallen?",
		StashPop:                            "Stash pop",
		SurePopStashEntry:                   "Weet je zeker dat je deze stash entry wil poppen?",
		StashApply:                          "Stash toepassen",
		SureApplyStashEntry:                 "Weet je zeker dat je deze stash entry wil toepassen?",
		NoStashTo:                           "Geen stash voor {{.method}}",
		NoTrackedStagedFilesStash:           "Je hebt geen tracked/staged bestanden om te laten stashen",
		StashChanges:                        "Stash veranderingen",
		IssntListOfViews:                    "{{.name}} is niet in de lijst van weergaves",
		LcNewFocusedViewIs:                  "nieuw gefocussed weergave is {{.newFocusedView}}",
		NoChangedFiles:                      "Geen veranderde bestanden",
		MergeAborted:                        "Merge afgebroken",
		OpenConfig:                          "open config bestand",
		EditConfig:                          "verander config bestand",
		ForcePush:                           "Forceer push",
		ForcePushPrompt:                     "Jouw branch is afgeweken van de remote branch. Druk 'esc' om te annuleren, of 'enter' om geforceert te pushen.",
		ForcePushDisabled:                   "Your branch has diverged from the remote branch and you've disabled force pushing",
		UpdatesRejectedAndForcePushDisabled: "Updates were rejected and you have disabled force pushing",
		LcCheckForUpdate:                    "check voor updates",
		CheckingForUpdates:                  "zoeken naar updates...",
		OnLatestVersionErr:                  "Je hebt al de laatste versie",
		MajorVersionErr:                     "Nieuwe versie ({{.newVersion}}) is niet backwards compatibele vergeleken met de huidige versie ({{.currentVersion}})",
		CouldNotFindBinaryErr:               "Kon geen binary vinden op {{.url}}",
		AnonymousReportingTitle:             "Help lazygit te verbeteren",
		AnonymousReportingPrompt:            "Zou je anonieme data rapportage willen aanzetten om lazygit beter te kunnen maken? (enter/esc)",
		IntroPopupMessage:                   "Bedankt voor het gebruik maken van lazygit! 2 dingen die je moet weten:\n\n1) Als je meer van lazygit zijn features wilt leren bekijk dan deze video:\n   https://youtu.be/CPLdltN7wgE\n\n2) Als je git gebruikt, ben je een programmeur! Met jouw hulp kunnen we lazygit verbeteren, dus overweeg om een ​​donateur te worden en mee te doen aan het plezier op\n   https://github.com/jesseduffield/lazygit",
		GitconfigParseErr:                   `Gogit kon je gitconfig bestand niet goed parsen door de aanwezigheid van losstaande '\' tekens. Het weghalen van deze tekens zou het probleem moeten oplossen. `,
		LcEditFile:                          `verander bestand`,
		LcOpenFile:                          `open bestand`,
		LcIgnoreFile:                        `voeg toe aan .gitignore`,
		LcRefreshFiles:                      `refresh bestanden`,
		LcMergeIntoCurrentBranch:            `merge in met huidige checked out branch`,
		ConfirmQuit:                         `Weet je zeker dat je dit programma wil sluiten?`,
		SwitchRepo:                          "wissel naar een recente repo",
		LcAllBranchesLogGraph:               `alle logs van de branch laten zien`,
		UnsupportedGitService:               `Niet-ondersteunde git-service`,
		LcCreatePullRequest:                 `maak een pull-request`,
		LcCopyPullRequestURL:                `kopieer de URL van het pull-verzoek naar het klembord`,
		NoBranchOnRemote:                    `Deze branch bestaat niet op de remote. U moet het eerst naar de remote pushen.`,
		LcFetch:                             `fetch`,
		NoAutomaticGitFetchTitle:            `Geen automatische git fetch`,
		NoAutomaticGitFetchBody:             `Lazygit kan niet "git fetch" uitvoeren in een privé repository, gebruik f in het branches paneel om "git fetch" manueel uit te voeren`,
		FileEnter:                           `stage individuele hunks/lijnen`,
		FileStagingRequirements:             `Kan alleen individuele lijnen stagen van getrackte bestanden met onstaged veranderingen`,
		SelectHunk:                          `selecteer hunk`,
		StageSelection:                      `toggle lijnen staged / unstaged`,
		ResetSelection:                      `verwijdert change (git reset)`,
		ToggleDragSelect:                    `toggle drag selecteer`,
		ToggleSelectHunk:                    `toggle selecteer hunk`,
		ToggleSelectionForPatch:             `voeg toe/verwijder lijn(en) in patch`,
		TogglePanel:                         `ga naar een ander paneel`,
		CantStageStaged:                     `Je kan niet al gestaged verandering stagen!`,
		ReturnToFilesPanel:                  `ga terug naar het bestanden paneel`,
		CantFindHunks:                       `Kan geen hunks vinden in deze patch`,
		CantFindHunk:                        `Kan geen hunk vinden`,
		FastForward:                         `fast-forward deze branch vanaf zijn upstream`,
		Fetching:                            "fetching en fast-forwarding {{.from}} -> {{.to}} ...",
		FoundConflicts:                      "Conflicten!, Om af te breken druk 'esc', anders druk op 'enter'",
		FoundConflictsTitle:                 "Auto-merge mislukt",
		Undo:                                "ongedaan maken",
		PickHunk:                            "kies hunk",
		PickAllHunks:                        "kies bijde hunks",
		ViewMergeRebaseOptions:              "bekijk merge/rebase opties",
		NotMergingOrRebasing:                "Je bent momenteel niet aan het rebasen of mergen",
		RecentRepos:                         "recente repositories",
		MergeOptionsTitle:                   "Merge Opties",
		RebaseOptionsTitle:                  "Rebase Opties",
		CommitMessageTitle:                  "Commit Bericht",
		LocalBranchesTitle:                  "Branches Tabblad",
		SearchTitle:                         "Zoek",
		TagsTitle:                           "Tags Tabblad",
		BranchCommitsTitle:                  "Commits Tabblad",
		MenuTitle:                           "Menu",
		RemotesTitle:                        "Remotes Tabblad",
		CredentialsTitle:                    "Credentials",
		RemoteBranchesTitle:                 "Remote Branches (in Remotes tabblad)",
		PatchBuildingTitle:                  "Patch Bouwen",
		InformationTitle:                    "Informatie",
		SecondaryTitle:                      "Secondary",
		ReflogCommitsTitle:                  "Reflog Tabblad",
		Title:                               "Title",
		GlobalTitle:                         "Globale Sneltoetsen",
		ConflictsResolved:                   "alle merge conflicten zijn opgelost. Wilt je verder gaan?",
		RebasingTitle:                       "Rebasen",
		MergingTitle:                        "Mergen",
		ConfirmRebase:                       "Weet je zeker dat je {{.checkedOutBranch}} op {{.selectedBranch}} wil rebasen?",
		ConfirmMerge:                        "Weet je zeker dat je {{.selectedBranch}} in {{.checkedOutBranch}} wil mergen?",
		FwdNoUpstream:                       "Kan niet de branch vooruitspoelen zonder upstream",
		FwdCommitsToPush:                    "Je kan niet vooruitspoelen als de branch geen nieuwe commits heeft",
		ErrorOccurred:                       "Er is iets fout gegaan! Zou je hier een issue aan willen maken",
		NoRoom:                              "Niet genoeg ruimte",
		YouAreHere:                          "JE BENT HIER",
		LcRewordNotSupported:                "herformatteren van commits in interactief rebasen is nog niet ondersteund",
		LcCherryPickCopy:                    "kopieer commit (cherry-pick)",
		LcCherryPickCopyRange:               "kopieer commit reeks (cherry-pick)",
		LcPasteCommits:                      "plak commits (cherry-pick)",
		SureCherryPick:                      "Weet je zeker dat je de gekopieerde commits naar deze branch wil cherry-picken?",
		CherryPick:                          "Cherry-Pick",
		CannotRebaseOntoFirstCommit:         "Je kan niet interactief rebasen naar de eerste commit",
		CannotSquashOntoSecondCommit:        "Je kan niet een squash/fixup doen naar de 2de commit",
		Donate:                              "Doneer",
		PrevLine:                            "selecteer de vorige lijn",
		NextLine:                            "selecteer de volgende lijn",
		PrevHunk:                            "selecteer de vorige hunk",
		NextHunk:                            "selecteer de volgende hunk",
		PrevConflict:                        "selecteer voorgaand conflict",
		NextConflict:                        "selecteer volgende conflict",
		SelectPrevHunk:                      "selecteer bovenste hunk",
		SelectNextHunk:                      "selecteer onderste hunk",
		ScrollDown:                          "scroll omlaag",
		ScrollUp:                            "scroll omhoog",
		LcScrollUpMainPanel:                 "scroll naar beneden vanaf hoofdpaneel",
		LcScrollDownMainPanel:               "scroll naar beneden vanaf hoofdpaneel",
		AmendCommitTitle:                    "Commit wijzigen",
		AmendCommitPrompt:                   "Weet je zeker dat je deze commit wil wijzigen met de vorige staged bestanden?",
		DeleteCommitTitle:                   "Verwijder Commit",
		DeleteCommitPrompt:                  "Weet je zeker dat je deze commit wil verwijderen?",
		SquashingStatus:                     "squashen",
		FixingStatus:                        "fixing up",
		DeletingStatus:                      "verwijderen",
		MovingStatus:                        "verplaatsen",
		RebasingStatus:                      "rebasen",
		AmendingStatus:                      "wijzigen",
		CherryPickingStatus:                 "cherry-picken",
		UndoingStatus:                       "ongedaan maken",
		RedoingStatus:                       "redoing",
		CheckingOutStatus:                   "uitchecken",
		CommitFiles:                         "Commit bestanden",
		LcViewCommitFiles:                   "bekijk gecommite bestanden",
		CommitFilesTitle:                    "Commit bestanden",
		LcGoBack:                            "ga terug",
		NoCommiteFiles:                      "Geen bestanden voor deze commit",
		LcCheckoutCommitFile:                "bestand uitchecken",
		LcDiscardOldFileChange:              "uitsluit deze commit zijn veranderingen aan dit bestand",
		DiscardFileChangesTitle:             "uitsluit bestand zijn veranderingen",
		DiscardFileChangesPrompt:            "Weet je zeker dat je de wijzigingen van deze commit in dit bestand wilt weggooien? Als dit bestand is gecreëerd in deze commit dan zal dit bestand worden verwijdert",
		DisabledForGPG:                      "Onderdelen niet beschikbaar voor gebruikers die GPG gebruiken",
		CreateRepo:                          "Niet in een git repository. Creëer een nieuwe git repository? (y/n): ",
		AutoStashTitle:                      "Autostash?",
		AutoStashPrompt:                     "Je moet je veranderingen stashen en poppen om ze over te brengen. Dit automatisch doen? (enter/esc)",
		StashPrefix:                         "Auto-stashing veranderingen voor ",
		LcViewDiscardOptions:                "bekijk 'veranderingen ongedaan maken' opties",
		LcCancel:                            "annuleren",
		LcDiscardAllChanges:                 "negeer alle wijzigingen",
		LcDiscardUnstagedChanges:            "negeer unstaged wijzigingen",
		LcDiscardAllChangesToAllFiles:       "verwijder werkende tree",
		LcDiscardAnyUnstagedChanges:         "gooi unstaged wijzigingen weg",
		LcDiscardUntrackedFiles:             "negeer niet-gevonden bestanden",
		LcViewResetOptions:                  `bekijk reset opties`,
		LcHardReset:                         "harde reset",
		LcHardResetUpstream:                 "harde naar upstream branch",
		LcCreateFixupCommit:                 `creëer fixup commit voor deze commit`,
		LcSquashAboveCommits:                `squash bovenstaande commits`,
		SquashAboveCommits:                  `Squash bovenstaande commits`,
		SureSquashAboveCommits:              `Weet je zeker dat je alles wil squash/fixup! voor de bovenstaand commits {{.commit}}?`,
		CreateFixupCommit:                   `Creëer fixup commit`,
		SureCreateFixupCommit:               `Weet je zeker dat je een fixup wil maken! commit voor commit {{.commit}}?`,
		LcExecuteCustomCommand:              "voor aangepaste commando uit",
		CustomCommand:                       "Aangepaste commando:",
		LcCommitChangesWithoutHook:          "commit veranderingen zonder pre-commit hook",
		SkipHookPrefixNotConfigured:         "Je hebt nog niet een commit bericht voorvoegsel ingesteld voor het overslaan van hooks. Set `git.skipHookPrefix = 'WIP'` in je config",
		LcResetTo:                           `reset naar`,
		PressEnterToReturn:                  "Press om terug te gaan naar lazygit",
		LcViewStashOptions:                  "bekijk stash opties",
		LcStashAllChanges:                   "stash-bestanden",
		LcStashStagedChanges:                "stash staged wijzigingen",
		LcStashOptions:                      "Stash opties",
		NotARepository:                      "Fout: moet in een git repository uitgevoerd worden",
		LcJump:                              "ga naar paneel",
		DiscardPatch:                        "Patch weg gooien",
		DiscardPatchConfirm:                 "Je kan alleen maar een patch bouwen van 1 commit. Huidige patch weggooien?",
		CantPatchWhileRebasingError:         "Je kan geen patch bouwen of patch commando uitvoeren wanneer je in een merging of rebasing state zit",
		LcToggleAddToPatch:                  "toggle bestand inbegrepen in patch",
		ViewPatchOptions:                    "bekijk aangepaste patch opties",
		PatchOptionsTitle:                   "Patch Opties",
		NoPatchError:                        "Nog geen patch gecreëerd. Om een patch te bouwen gebruik 'space' op een commit bestand of 'enter' om een spesiefieke lijnen toe te voegen",
		LcEnterFile:                         "enter bestand om geselecteerde regels toe te voegen aan de patch",
		ExitLineByLineMode:                  `sluit lijn-bij-lijn modus`,
		EnterUpstream:                       `Enter upstream als '<remote> <branchnaam>'`,
		EnterUpstreamWithSlash:              `Enter upstream als '<remote>/<branchnaam>'`,
		LcNotTrackingRemote:                 "(nog geen remote aan het volgen)",
		ReturnToRemotesList:                 `Ga terug naar remotes lijst`,
		LcAddNewRemote:                      `voeg een nieuwe remote toe`,
		LcNewRemoteName:                     `Nieuwe remote name:`,
		LcNewRemoteUrl:                      `Nieuwe remote url:`,
		LcEditRemoteName:                    `Enter updated remote naam voor {{.remoteName}}:`,
		LcEditRemoteUrl:                     `Enter updated remote url voor {{.remoteName}}:`,
		LcRemoveRemote:                      `verwijder remote`,
		LcRemoveRemotePrompt:                "Weet je zeker dat je deze remote wilt verwijderen",
		DeleteRemoteBranch:                  "Verwijder Remote Branch",
		DeleteRemoteBranchMessage:           "Weet je zeker dat je deze remote branch wilt verwijderen",
		LcSetUpstream:                       "stel in als upstream van uitgecheckte branch",
		SetUpstreamTitle:                    "Stel in als upstream branch",
		SetUpstreamMessage:                  "Weet je zeker dat je de upstream branch van '{{.checkedOut}}' naar '{{.selected}}' wilt zetten",
		LcEditRemote:                        "wijzig remote",
		LcTagCommit:                         "tag commit",
		TagNameTitle:                        "Tag naam:",
		LcDeleteTag:                         "verwijder tag",
		DeleteTagTitle:                      "Verwijder tag",
		DeleteTagPrompt:                     "Weet je zeker dat je '{{.tagName}}' wil verwijderen?",
		PushTagTitle:                        "remote om tag '{{.tagName}}' te pushen naar:",
		LcPushTag:                           "push tag",
		LcCreateTag:                         "creëer tag",
		CreateTagTitle:                      "Tag naam:",
		LcFetchRemote:                       "fetch remote",
		FetchingRemoteStatus:                "remote fetchen",
		LcCheckoutCommit:                    "checkout commit",
		SureCheckoutThisCommit:              "Weet je zeker dat je deze commit wil uitchecken?",
		LcGitFlowOptions:                    "laat git-flow opties zien",
		NotAGitFlowBranch:                   "Dit lijkt geen git flow branch te zijn",
		NewGitFlowBranchPrompt:              "nieuwe {{.branchType}} naam:",
		IgnoreTracked:                       "Negeer tracked bestand",
		IgnoreTrackedPrompt:                 "weet je zeker dat je een getracked bestand wil negeren?",
		LcViewResetToUpstreamOptions:        "bekijk upstream reset opties",
		LcNextScreenMode:                    "volgende scherm modus (normaal/half/groot)",
		LcPrevScreenMode:                    "vorige scherm modus",
		LcStartSearch:                       "start met zoeken",
		Panel:                               "Paneel",
		Keybindings:                         "Sneltoetsen",
		LcRenameBranch:                      "hernoem branch",
		NewBranchNamePrompt:                 "Noem een nieuwe branch naam",
		RenameBranchWarning:                 "Deze branch volgt een remote. Deze actie zal alleen de locale branch name wijzigen niet de naam van de remote branch. Verder gaan?",
		LcOpenMenu:                          "open menu",
		LcCloseMenu:                         "sluit menu",
		LcResetCherryPick:                   "reset cherry-picked (gekopieerde) commits selectie",
		LcNextTab:                           "volgende tabblad",
		LcPrevTab:                           "vorige tabblad",
		LcCantUndoWhileRebasing:             "Kan niet ongedaan maken terwijl je aan het rebasen bent",
		LcCantRedoWhileRebasing:             "Kan niet opnieuw doen (redo) terwijl je aan het rebasen bent",
		MustStashWarning:                    "Een patch in de index stoppen vereist stashen en onstashen van je wijzigingen. Als er iets verkeert gaat kan je je bestanden terug vinden in de stash. Verder gaan?",
		MustStashTitle:                      "Moet stashen",
		ConfirmationTitle:                   "Bevestigingspaneel",
		LcPrevPage:                          "vorige pagina",
		LcNextPage:                          "volgende pagina",
		LcGotoTop:                           "scroll naar boven",
		LcGotoBottom:                        "scroll naar beneden",
		LcFilteringBy:                       "filteren bij",
		ResetInParentheses:                  "(reset)",
		LcOpenFilteringMenu:                 "bekijk scoping opties",
		LcFilterBy:                          "filter bij",
		LcExitFilterMode:                    "stop met filteren bij pad",
		LcFilterPathOption:                  "vulin pad om op te filteren",
		LcEnterFileName:                     "vulin path:",
		FilteringMenuTitle:                  "Filteren",
		MustExitFilterModeTitle:             "Command niet beschikbaar",
		MustExitFilterModePrompt:            "Command niet beschikbaar in filter modus. Sluit filter modus?",
		LcDiff:                              "diff",
		LcEnterRefToDiff:                    "vul in ref naar diff",
		LcEnteRefName:                       "vul in ref:",
		LcExitDiffMode:                      "sluit diff mode",
		DiffingMenuTitle:                    "Diffen",
		LcSwapDiff:                          "keer diff richting om",
		LcOpenDiffingMenu:                   "open diff menu",
		LcShowingGitDiff:                    "laat output zien voor:",
		LcCopyCommitShaToClipboard:          "kopieer commit SHA naar klembord",
		LcCopyCommitMessageToClipboard:      "kopieer commit bericht naar klembord",
		LcCopyBranchNameToClipboard:         "kopieer branch name naar klembord",
		LcCopyFileNameToClipboard:           "kopieer de bestandsnaam naar het klembord",
		LcCopyCommitFileNameToClipboard:     "kopieer de vastgelegde bestandsnaam naar het klembord",
		LcCommitPrefixPatternError:          "Fout in commitPrefix patroon",
		NoFilesStagedTitle:                  "geen bestanden gestaged",
		NoFilesStagedPrompt:                 "Je hebt geen bestanden gestaged. Commit alle bestanden?",
		BranchNotFoundTitle:                 "Branch niet gevonden",
		BranchNotFoundPrompt:                "Branch niet gevonden. Creëer een nieuwe branch genaamd",
		PullRequestURLCopiedToClipboard:     "Pull-aanvraag-URL gekopieerd naar klembord",
		CommitMessageCopiedToClipboard:      "Commit message gekopieerd naar klembord",
		LcCopiedToClipboard:                 "gekopieerd naar klembord",
		NavigationTitle:                     "Lijstpaneel Navigatie",
		LcViewCommits:                       "bekijk commits",
		LcToggleTreeView:                    "toggle bestandsboom weergave",
		LcCreateNewBranchFromCommit:         "creëer nieuwe branch van commit",
		LcCopySubmoduleNameToClipboard:      "kopieer submodule naam naar klembord",
		LcEnterSubmodule:                    "enter submodule",
		LcViewResetAndRemoveOptions:         "bekijk reset en verwijder submodule opties",
		LcAddSubmodule:                      "voeg nieuwe submodule toe",
		LcInitSubmodule:                     "initialiseer submodule",
		LcViewBulkSubmoduleOptions:          "bekijk bulk submodule opties",
		LcViewStashFiles:                    "bekijk bestanden van stash entry",
		CreatePullRequestOptions:            "Bekijk opties voor pull-aanvraag",
		LcCreatePullRequestOptions:          "bekijk opties voor pull-aanvraag",
	}
}
