package icons

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

var (
	BRANCH_ICON         = map[int]string{2: "\ufb2b", 3: "\U000f062c"} // שׂ 󰘬
	DETACHED_HEAD_ICON  = "\ue729"                                     // 
	TAG_ICON            = "\uf02b"                                     // 
	COMMIT_ICON         = map[int]string{2: "\ufc16", 3: "\U000f0718"} // ﰖ  󰜘
	MERGE_COMMIT_ICON   = map[int]string{2: "\ufb2c", 3: "\U000f062d"} // שּׁ 󰘭
	DEFAULT_REMOTE_ICON = map[int]string{2: "\uf7a1", 3: "\U000f02a2"} //   󰊢
	STASH_ICON          = "\uf01c"                                     // 
)

type remoteIcon struct {
	domain string
	icons  map[int]string
}

var remoteIcons = []remoteIcon{
	{domain: "github.com", icons: map[int]string{2: "\ue709"}},                     // 
	{domain: "bitbucket.org", icons: map[int]string{2: "\ue703"}},                  // 
	{domain: "gitlab.com", icons: map[int]string{2: "\uf296"}},                     // 
	{domain: "dev.azure.com", icons: map[int]string{2: "\ufd03", 3: "\U000f0805"}}, // ﴃ  or  󰠅
}

func IconForBranch(branch *models.Branch) string {
	if branch.DetachedHead {
		return DETACHED_HEAD_ICON
	}
	return BRANCH_ICON[GetNerdFontsVersion()]
}

func IconForRemoteBranch(branch *models.RemoteBranch) string {
	return BRANCH_ICON[GetNerdFontsVersion()]
}

func IconForTag(tag *models.Tag) string {
	return TAG_ICON
}

func IconForCommit(commit *models.Commit) string {
	if len(commit.Parents) > 1 {
		return MERGE_COMMIT_ICON[GetNerdFontsVersion()]
	}
	return COMMIT_ICON[GetNerdFontsVersion()]
}

func IconForRemote(remote *models.Remote) string {
	for _, r := range remoteIcons {
		for _, url := range remote.Urls {
			if strings.Contains(url, r.domain) {
				if icon, ok := r.icons[GetNerdFontsVersion()]; ok {
					return icon
				}
				return r.icons[2]
			}
		}
	}
	return DEFAULT_REMOTE_ICON[GetNerdFontsVersion()]
}

func IconForStash(stash *models.StashEntry) string {
	return STASH_ICON
}
