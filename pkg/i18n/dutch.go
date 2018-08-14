package lang

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func addDutch(i18nObject *i18n.Bundle) *i18n.Bundle {
	i18nObject.AddMessages(language.English,
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
		},
	)
	return i18nObject
}
