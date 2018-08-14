package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// addDutch will add all dutch translations
func addDutch(i18nObject *i18n.Bundle) *i18n.Bundle {

	// add the translations
	i18nObject.AddMessages(language.Dutch,
		&i18n.Message{
			ID:    "NotEnoughSpace",
			Other: "Niet genoeg ruimte om de panelen te renderen",
		}, &i18n.Message{
			ID:    "DiffTitle",
			Other: "Diff",
		}, &i18n.Message{
			ID:    "FilesTitle",
			Other: "Bestanden",
		}, &i18n.Message{
			ID:    "BranchesTitle",
			Other: "Branches",
		}, &i18n.Message{
			ID:    "CommitsTitle",
			Other: "Commits",
		}, &i18n.Message{
			ID:    "StashTitle",
			Other: "Stash",
		}, &i18n.Message{
			ID:    "CommitMessage",
			Other: "Commit Bericht",
		}, &i18n.Message{
			ID:    "CommitChanges",
			Other: "Commit Veranderingen",
		}, &i18n.Message{
			ID:    "StatusTitle",
			Other: "Status",
		}, &i18n.Message{
			ID:    "navigate",
			Other: "navigeer",
		}, &i18n.Message{
			ID:    "stashFiles",
			Other: "stash-bestanden",
		}, &i18n.Message{
			ID:    "open",
			Other: "open",
		}, &i18n.Message{
			ID:    "ignore",
			Other: "negeren",
		}, &i18n.Message{
			ID:    "delete",
			Other: "verwijderen",
		}, &i18n.Message{
			ID:    "toggleStaged",
			Other: "toggle staged",
		}, &i18n.Message{
			ID:    "refresh",
			Other: "verversen",
		}, &i18n.Message{
			ID:    "addPatch",
			Other: "verandering toevoegen",
		}, &i18n.Message{
			ID:    "edit",
			Other: "veranderen",
		}, &i18n.Message{
			ID:    "scroll",
			Other: "scroll",
		}, &i18n.Message{
			ID:    "abortMerge",
			Other: "samenvoegen afbreken",
		}, &i18n.Message{
			ID:    "resolveMergeConflicts",
			Other: "verhelp samenvoegen fouten",
		}, &i18n.Message{
			ID:    "checkout",
			Other: "uitchecken",
		}, &i18n.Message{
			ID:    "NoChangedFiles",
			Other: "Geen Bestanden verandert",
		}, &i18n.Message{
			ID:    "NoUsernameSetErr",
			Other: `Geen gebruikersnaam ingesteld. Doe: git config --global user.name "Jou Naam"`,
		}, &i18n.Message{
			ID:    "FileHasNoUnstagedChanges",
			Other: "Het bestand heeft geen unstaged veranderingen om toe te voegen",
		}, &i18n.Message{
			ID:    "CannotGitAdd",
			Other: "Kan commando niet uitvoeren git add --path untracked files",
		}, &i18n.Message{
			ID:    "CantIgnoreTrackFiles",
			Other: "Kan gevolgde bestanden niet negeren",
		},
	)

	// return the new i18nObject
	return i18nObject
}
