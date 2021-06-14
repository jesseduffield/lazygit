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
	SurfaceError(error) error
	GetGit() commands.IGit
	Prompt(PromptOpts) error
	Ask(AskOpts) error
	GetTr() *i18n.TranslationSet
	CreateErrorPanel(string) error
	WithPopupWaitingStatus(string, func() error) error
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

	opts := commands.PushOpts{}

	if currentBranch.IsTrackingRemote() {
		if currentBranch.HasCommitsToPull() {
			return gui.requestToForcePush(opts, gui.GetTr().ForcePushDisabled)
		} else {
			return gui.attemptToPush(opts)
		}
	} else {
		opts.SetUpstream = true
		// see if we have an upstream for this branch in our config
		remoteName, err := gui.GetGit().FindRemoteForBranchInConfig(currentBranch.Name)
		if err != nil {
			return gui.SurfaceError(err)
		}
		if remoteName != "" {
			opts.DestinationRemote = remoteName
			opts.DestinationBranch = currentBranch.Name

			return gui.attemptToPush(opts)
		}

		if gui.GetGit().GetPushToCurrent() {
			return gui.attemptToPush(opts)
		} else {
			return gui.promptToSetDestinationAndPush(opts, currentBranch.Name)
		}
	}
}

func (gui *PushFilesHandler) promptToSetDestinationAndPush(opts commands.PushOpts, currentBranchName string) error {
	return gui.Prompt(PromptOpts{
		Title:          gui.GetTr().EnterUpstream,
		InitialContent: "origin " + currentBranchName,
		HandleConfirm: func(upstream string) error {
			split := strings.Split(upstream, " ")
			remote, branch := split[0], split[1]
			opts.DestinationRemote = remote
			opts.DestinationBranch = branch
			return gui.attemptToPush(opts)
		},
	})
}

func (gui *PushFilesHandler) attemptToPush(opts commands.PushOpts) error {
	return gui.WithPopupWaitingStatus(gui.GetTr().PushWait, func() error {
		rejected, err := gui.GetGit().WithSpan(gui.GetTr().Spans.Push).Push(opts)
		if err != nil {
			return gui.SurfaceError(err)
		}
		if rejected {
			if opts.Force {
				// this should really never happen
				return gui.CreateErrorPanel("Force push rejected")
			} else {
				return gui.requestToForcePush(opts, gui.GetTr().UpdatesRejectedAndForcePushDisabled)
			}
		}

		_ = gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})

		return nil
	})
}

func (gui *PushFilesHandler) isRejectionErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Updates were rejected")
}

func (gui *PushFilesHandler) requestToForcePush(opts commands.PushOpts, messageOnError string) error {
	opts.Force = true

	forcePushDisabled := gui.GetUserConfig().Git.DisableForcePushing
	if forcePushDisabled {
		_ = gui.CreateErrorPanel(messageOnError)
		return nil
	}

	return gui.Ask(AskOpts{
		Title:  gui.GetTr().ForcePush,
		Prompt: gui.GetTr().ForcePushPrompt,
		HandleConfirm: func() error {
			// return request to /blah with these options: {opts}
			return gui.attemptToPush(opts)
		},
	})
}
