package i18n

const frenchIntroPopupMessage = `
Merci d'utiliser lazygit! 3 choses que je souhaite te partager:

 1) Pour en apprendre plus sur les fonctionnalitées de lazygit, regarde
			cette vidéo (en anglais):
      https://youtu.be/CPLdltN7wgE

 2) Les notes de mise à jour sont disponible en suivant ce lien:
      https://github.com/jesseduffield/lazygit/releases

 3) Si tu utilises git, alors tu es un développeur! Avec ton aide, nous pouvons
		améliorer laygit, donc n'hésite pas à nous en rejoindre en devenant un
		contributeur:
      https://github.com/jesseduffield/lazygit
		Tu peux aussi me sponsoriser et me dire ce sur quoi tu souhaiterais que je
		travaille en suivant ce lien:
			https://github.com/sponsors/jesseduffield
		Ou même donner une étoile au dépôt git pour montrer ton soutiens!
`

// exporting this so we can use it in tests
func frenchTranslationSet() TranslationSet {
	return TranslationSet{
		NotEnoughSpace:                      "Pas assez de place pour afficher les panels",
		DiffTitle:                           "Diff",
		FilesTitle:                          "Fichier",
		BranchesTitle:                       "Branches",
		CommitsTitle:                        "Commits",
		StashTitle:                          "Stash",
		UnstagedChanges:                     `Changements Non-Indexées`,
		StagedChanges:                       `Changement Indexées`,
		MainTitle:                           "Principal",
		MergeConfirmTitle:                   "Fusion",
		StagingTitle:                        "Panel principal (Indexation)",
		MergingTitle:                        "Panel principal (Fusion)",
		NormalTitle:                         "Panel principal (Normal)",
		CommitMessage:                       "Message de Commit",
		CredentialsUsername:                 "Nom d'utilisateur",
		CredentialsPassword:                 "Mot de passe",
		CredentialsPassphrase:               "Entrer la phrase secrète de votre clé SSH",
		PassUnameWrong:                      "Mauvais mot de passe, phrase secrète et/ou nom d'utilisateur",
		CommitChanges:                       "commit les changements",
		AmendLastCommit:                     "rectifier le dernier commit",
		SureToAmend:                         "Êtes vous sûr de vouloir rectifier le dernier commit? Vous pourrez ensuite modifier le message de commit depuis le panel des commits.",
		NoCommitToAmend:                     "Il n'y a aucun commit à rectifier",
		CommitChangesWithEditor:             "commit les changements en utilisant l'éditeur git par défaut",
		StatusTitle:                         "Status",
		LcNavigate:                          "naviguer",
		LcMenu:                              "menu",
		LcExecute:                           "executer",
		LcToggleStaged:                      "alterner indexé/non-indexé",
		LcToggleStagedAll:                   "tout indexés/non-indexés",
		LcToggleTreeView:                    "afficher/cacher l'arbre des fichiers",
		LcOpenMergeTool:                     "Ouvrir l'outil de fusion externe (git mergetool)",
		LcRefresh:                           "rafraîchir",
		LcPush:                              "push",
		LcPull:                              "pull",
		LcScroll:                            "défiler",
		MergeConflictsTitle:                 "Conflits de Fusion",
		LcCheckout:                          "changer de branche",
		LcFileFilter:                        "Filtrer les fichiers (indexés/non-indexés)",
		FilterStagedFiles:                   "Uniquement afficher les fichiers indexés",
		FilterUnstagedFiles:                 "Uniquement afficher les fichiers non-indexés",
		ResetCommitFilterState:              "Réinitialiser les filtres",
	}
}
