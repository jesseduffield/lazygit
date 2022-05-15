package context

import (
	"log"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func shouldShowGraph(state GuiContextState, userConfig *config.UserConfig) bool {
	if state.Modes().Filtering.Active() {
		return false
	}

	if state.Needle() != "" {
		return false
	}

	value := userConfig.Git.Log.ShowGraph
	switch value {
	case "always":
		return true
	case "never":
		return false
	case "when-maximised":
		return state.ScreenMode() != types.SCREEN_NORMAL
	}

	log.Fatalf("Unknown value for git.log.showGraph: %s. Expected one of: 'always', 'never', 'when-maximised'", value)
	return false
}

func getCommitDisplayStrings(
	selectedCommit *models.Commit,
	commits []*models.Commit,
	guiContextState GuiContextState,
	userConfig *config.UserConfig,
	startIdx int,
	length int,
) [][]string {
	selectedCommitSha := ""
	if guiContextState.IsFocused() && selectedCommit != nil {
		selectedCommitSha = selectedCommit.Sha
	}
	return presentation.GetCommitListDisplayStrings(
		commits,
		guiContextState.ScreenMode() != types.SCREEN_NORMAL,
		cherryPickedCommitShaSet(guiContextState),
		guiContextState.Modes().Diffing.Ref,
		userConfig.Git.ParseEmoji,
		userConfig.Gui.TimeFormat,
		selectedCommitSha,
		startIdx,
		length,
		shouldShowGraph(guiContextState, userConfig),
		guiContextState.BisectInfo(),
	)
}
