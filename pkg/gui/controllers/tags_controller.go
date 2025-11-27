package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type TagsController struct {
	baseController
	*ListControllerTrait[*models.Tag]
	c *ControllerCommon
}

var _ types.IController = &TagsController{}

func NewTagsController(
	c *ControllerCommon,
) *TagsController {
	return &TagsController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().Tags,
			c.Contexts().Tags.GetSelected,
			c.Contexts().Tags.GetSelectedItems,
		),
		c: c,
	}
}

func (self *TagsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.checkout),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			Tooltip:           self.c.Tr.TagCheckoutTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:             opts.GetKey(opts.Config.Universal.New),
			Handler:         self.create,
			Description:     self.c.Tr.NewTag,
			Tooltip:         self.c.Tr.NewTagTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItem(self.delete),
			Description:       self.c.Tr.Delete,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Tooltip:           self.c.Tr.TagDeleteTooltip,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.PushTag),
			Handler:           self.withItem(self.push),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.PushTag,
			Tooltip:           self.c.Tr.PushTagTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:           self.withItem(self.createResetMenu),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Reset,
			Tooltip:           self.c.Tr.ResetTooltip,
			DisplayOnScreen:   true,
			OpensMenu:         true,
		},
		{
			Key: opts.GetKey(opts.Config.Universal.OpenDiffTool),
			Handler: self.withItem(func(selectedTag *models.Tag) error {
				return self.c.Helpers().Diff.OpenDiffToolForRef(selectedTag)
			}),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
	}

	return bindings
}

func (self *TagsController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			var task types.UpdateTask
			tag := self.context().GetSelected()
			if tag == nil {
				task = types.NewRenderStringTask("No tags")
			} else {
				cmdObj := self.c.Git().Branch.GetGraphCmdObj(tag.FullRefName())
				prefix := self.getTagInfo(tag) + "\n\n---\n\n"
				task = types.NewRunCommandTaskWithPrefix(cmdObj.GetCmd(), prefix)
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: "Tag",
					Task:  task,
				},
			})
		})
	}
}

func (self *TagsController) getTagInfo(tag *models.Tag) string {
	tagIsAnnotated, err := self.c.Git().Tag.IsTagAnnotated(tag.Name)
	if err != nil {
		self.c.Log.Warnf("Error checking if tag is annotated: %v", err)
	}

	if tagIsAnnotated {
		info := fmt.Sprintf("%s: %s", self.c.Tr.AnnotatedTag, style.AttrBold.Sprint(style.FgYellow.Sprint(tag.Name)))
		output, err := self.c.Git().Tag.ShowAnnotationInfo(tag.Name)
		if err == nil {
			info += "\n\n" + strings.TrimRight(filterOutPgpSignature(output), "\n")
		}
		return info
	}

	return fmt.Sprintf("%s: %s", self.c.Tr.LightweightTag, style.AttrBold.Sprint(style.FgYellow.Sprint(tag.Name)))
}

func filterOutPgpSignature(output string) string {
	lines := strings.Split(output, "\n")
	inPgpSignature := false
	filteredLines := lo.Filter(lines, func(line string, _ int) bool {
		if line == "-----END PGP SIGNATURE-----" {
			inPgpSignature = false
			return false
		}
		if line == "-----BEGIN PGP SIGNATURE-----" {
			inPgpSignature = true
		}
		return !inPgpSignature
	})
	return strings.Join(filteredLines, "\n")
}

func (self *TagsController) checkout(tag *models.Tag) error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutTag)
	if err := self.c.Helpers().Refs.CheckoutRef(tag.FullRefName(), types.CheckoutRefOptions{}); err != nil {
		return err
	}
	return nil
}

func (self *TagsController) localDelete(tag *models.Tag) error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.DeleteLocalTag)
		err := self.c.Git().Tag.LocalDelete(tag.Name)
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS}})
		return err
	})
}

func (self *TagsController) remoteDelete(tag *models.Tag) error {
	title := utils.ResolvePlaceholderString(
		self.c.Tr.SelectRemoteTagUpstream,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	self.c.Prompt(types.PromptOpts{
		Title:               title,
		InitialContent:      "origin",
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteSuggestionsFunc(),
		HandleConfirm: func(upstream string) error {
			confirmTitle := utils.ResolvePlaceholderString(
				self.c.Tr.DeleteTagTitle,
				map[string]string{
					"tagName": tag.Name,
				},
			)
			confirmPrompt := utils.ResolvePlaceholderString(
				self.c.Tr.DeleteRemoteTagPrompt,
				map[string]string{
					"tagName":  tag.Name,
					"upstream": upstream,
				},
			)

			self.c.Confirm(types.ConfirmOpts{
				Title:  confirmTitle,
				Prompt: confirmPrompt,
				HandleConfirm: func() error {
					return self.c.WithInlineStatus(tag, types.ItemOperationDeleting, context.TAGS_CONTEXT_KEY, func(task gocui.Task) error {
						self.c.LogAction(self.c.Tr.Actions.DeleteRemoteTag)
						if err := self.c.Git().Remote.DeleteRemoteTag(task, upstream, tag.Name); err != nil {
							return err
						}
						self.c.Toast(self.c.Tr.RemoteTagDeletedMessage)
						self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS}})
						return nil
					})
				},
			})

			return nil
		},
	})

	return nil
}

func (self *TagsController) localAndRemoteDelete(tag *models.Tag) error {
	title := utils.ResolvePlaceholderString(
		self.c.Tr.SelectRemoteTagUpstream,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	self.c.Prompt(types.PromptOpts{
		Title:               title,
		InitialContent:      "origin",
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteSuggestionsFunc(),
		HandleConfirm: func(upstream string) error {
			confirmTitle := utils.ResolvePlaceholderString(
				self.c.Tr.DeleteTagTitle,
				map[string]string{
					"tagName": tag.Name,
				},
			)
			confirmPrompt := utils.ResolvePlaceholderString(
				self.c.Tr.DeleteLocalAndRemoteTagPrompt,
				map[string]string{
					"tagName":  tag.Name,
					"upstream": upstream,
				},
			)

			self.c.Confirm(types.ConfirmOpts{
				Title:  confirmTitle,
				Prompt: confirmPrompt,
				HandleConfirm: func() error {
					return self.c.WithInlineStatus(tag, types.ItemOperationDeleting, context.TAGS_CONTEXT_KEY, func(task gocui.Task) error {
						self.c.LogAction(self.c.Tr.Actions.DeleteRemoteTag)
						if err := self.c.Git().Remote.DeleteRemoteTag(task, upstream, tag.Name); err != nil {
							return err
						}

						self.c.LogAction(self.c.Tr.Actions.DeleteLocalTag)
						if err := self.c.Git().Tag.LocalDelete(tag.Name); err != nil {
							return err
						}
						self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS}})
						return nil
					})
				},
			})

			return nil
		},
	})

	return nil
}

func (self *TagsController) delete(tag *models.Tag) error {
	menuTitle := utils.ResolvePlaceholderString(
		self.c.Tr.DeleteTagTitle,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	menuItems := []*types.MenuItem{
		{
			Label: self.c.Tr.DeleteLocalTag,
			Key:   'c',
			OnPress: func() error {
				return self.localDelete(tag)
			},
		},
		{
			Label:     self.c.Tr.DeleteRemoteTag,
			Key:       'r',
			OpensMenu: true,
			OnPress: func() error {
				return self.remoteDelete(tag)
			},
		},
		{
			Label:     self.c.Tr.DeleteLocalAndRemoteTag,
			Key:       'b',
			OpensMenu: true,
			OnPress: func() error {
				return self.localAndRemoteDelete(tag)
			},
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: menuTitle,
		Items: menuItems,
	})
}

func (self *TagsController) push(tag *models.Tag) error {
	title := utils.ResolvePlaceholderString(
		self.c.Tr.PushTagTitle,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	self.c.Prompt(types.PromptOpts{
		Title:               title,
		InitialContent:      "origin",
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			return self.c.WithInlineStatus(tag, types.ItemOperationPushing, context.TAGS_CONTEXT_KEY, func(task gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.PushTag)
				err := self.c.Git().Tag.Push(task, response, tag.Name)

				// Render again to remove the inline status:
				self.c.OnUIThread(func() error {
					self.c.Contexts().Tags.HandleRender()
					return nil
				})

				return err
			})
		},
	})

	return nil
}

func (self *TagsController) createResetMenu(tag *models.Tag) error {
	return self.c.Helpers().Refs.CreateGitResetMenu(tag.Name, tag.FullRefName())
}

func (self *TagsController) create() error {
	// leaving commit hash blank so that we're just creating the tag for the current commit
	return self.c.Helpers().Tags.OpenCreateTagPrompt("", func() {
		self.context().SetSelection(0)
	})
}

func (self *TagsController) context() *context.TagsContext {
	return self.c.Contexts().Tags
}
