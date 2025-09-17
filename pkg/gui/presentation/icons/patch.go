package icons

import "github.com/jesseduffield/lazygit/pkg/config"

func PatchHardcodedIcons(customIcons config.CustomIconsConfig) {
	patchFileIconFromConfig(customIcons.FileIcons)
	patchVCSIconFromConfig(customIcons.VCSIcons)
}

func patchVCSIconFromConfig(vcsIcons map[string]string) {
	for name, icon := range vcsIcons {
		// replace Variables(e.x: SOME_UPPER_SNAKECASE_VARIABLE)
		switch name {
		case "branch":
			BRANCH_ICON = icon
		case "detached-head":
			DETACHED_HEAD_ICON = icon
		case "tag":
			TAG_ICON = icon
		case "commit":
			COMMIT_ICON = icon
		case "merge-commit":
			MERGE_COMMIT_ICON = icon
		case "remote":
			DEFAULT_REMOTE_ICON = icon
		case "stash":
			STASH_ICON = icon
		case "linked-worktree":
			LINKED_WORKTREE_ICON = icon
		case "missing-linked-worktree":
			MISSING_LINKED_WORKTREE_ICON = icon
		}
		// replace value in remoteIconsMap if the key exists.
		if _, ok := remoteIcons[name]; ok {
			remoteIcons[name] = icon
		}
	}
}

func changeStruct(p config.IconProperties) IconProperties {
	return IconProperties{Icon: p.Icon, Color: p.Color}
}

func patchFileIconFromConfig(fileIcons map[string]config.IconProperties) {
	for name, icon := range fileIcons {
		switch name {
		case "file":
			DEFAULT_FILE_ICON = changeStruct(icon)
		case "submodule":
			DEFAULT_SUBMODULE_ICON = changeStruct(icon)
		case "directory":
			DEFAULT_DIRECTORY_ICON = changeStruct(icon)
		}
	}
}
