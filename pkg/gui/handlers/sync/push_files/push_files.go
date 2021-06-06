package push_files

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

type Gui interface {
	PopupPanelFocused() bool
	CurrentBranch() *models.Branch
	GetUserConfig() *config.UserConfig
	UpstreamForBranchInConfig(string) (string, error)
	SurfaceError(error) error
	GetGitCommand() commands.IGitCommand
	Prompt(PromptOpts) error
	Ask(AskOpts) error
	GetTr() *i18n.TranslationSet
	CreateErrorPanel(string) error
	HandleCredentialsPopup(error)
	WithPopupWaitingStatus(string, func() error) error
	PromptUserForCredential(passOrUname string) string
	RefreshSidePanels(RefreshOptions) error
}

type PushFilesHandler struct {
	Gui
}

func New(gui Gui) *PushFilesHandler {
	return &PushFilesHandler{Gui: gui}
}

func (gui *PushFilesHandler) Run() error {
	if gui.PopupPanelFocused() {
		return nil
	}

	// if we have pullables we'll ask if the user wants to force push
	currentBranch := gui.CurrentBranch()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return nil
	}

	if currentBranch.IsTrackingRemote() {
		if currentBranch.HasCommitsToPull() {
			return gui.requestToForcePush()
		} else {
			return gui.pushWithForceFlag(false, "", "")
		}
	} else {
		// see if we have an upstream for this branch in our config
		upstream, err := gui.UpstreamForBranchInConfig(currentBranch.Name)
		if err != nil {
			return gui.SurfaceError(err)
		}

		if upstream != "" {
			return gui.pushWithForceFlag(false, "", upstream)
		}

		if gui.GetGitCommand().GetPushToCurrent() {
			return gui.pushWithForceFlag(false, "", "--set-upstream")
		} else {
			return gui.promptToSetUpstreamAndPush(currentBranch.Name)
		}
	}
}

func (gui *PushFilesHandler) promptToSetUpstreamAndPush(currentBranchName string) error {
	return gui.Prompt(PromptOpts{
		Title:          gui.GetTr().EnterUpstream,
		InitialContent: "origin " + currentBranchName,
		HandleConfirm: func(upstream string) error {
			return gui.pushWithForceFlag(false, upstream, "")
		},
	})
}

func (gui *PushFilesHandler) pushWithForceFlag(force bool, upstream string, args string) error {
	return gui.WithPopupWaitingStatus(gui.GetTr().PushWait, func() error {
		branchName := gui.CurrentBranch().Name

		err := gui.GetGitCommand().WithSpan(gui.GetTr().Spans.Push).Push(branchName, force, upstream, args, gui.PromptUserForCredential)
		if err != nil && !force && strings.Contains(err.Error(), "Updates were rejected") {
			forcePushDisabled := gui.GetUserConfig().Git.DisableForcePushing
			if forcePushDisabled {
				_ = gui.CreateErrorPanel(gui.GetTr().UpdatesRejectedAndForcePushDisabled)
				return nil
			}

			return gui.Ask(AskOpts{
				Title:  gui.GetTr().ForcePush,
				Prompt: gui.GetTr().ForcePushPrompt,
				HandleConfirm: func() error {
					return gui.pushWithForceFlag(true, upstream, args)
				},
			})
		}

		gui.HandleCredentialsPopup(err)
		_ = gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})

		return nil
	})
}

func (gui *PushFilesHandler) requestToForcePush() error {
	forcePushDisabled := gui.GetUserConfig().Git.DisableForcePushing
	if forcePushDisabled {
		return gui.CreateErrorPanel(gui.GetTr().ForcePushDisabled)
	}

	return gui.Ask(AskOpts{
		Title:  gui.GetTr().ForcePush,
		Prompt: gui.GetTr().ForcePushPrompt,
		HandleConfirm: func() error {
			return gui.pushWithForceFlag(true, "", "")
		},
	})
}
