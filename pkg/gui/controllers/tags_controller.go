package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
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
			Keys:              opts.GetKeys(opts.Config.Universal.Select),
			Handler:           self.withItem(self.checkout),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			Tooltip:           self.c.Tr.TagCheckoutTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.New),
			Handler:         self.create,
			Description:     self.c.Tr.NewTag,
			Tooltip:         self.c.Tr.NewTagTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.NewWorktree),
			Handler:     self.withItem(self.c.Helpers().Worktree.NewWorktreeMenuForTag),
			Description: self.c.Tr.NewWorktree,
			OpensMenu:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Remove),
			Handler:           self.withItems(self.delete),
			Description:       self.c.Tr.Delete,
			GetDisabledReason: self.require(self.itemsSelected()),
			Tooltip:           self.c.Tr.TagDeleteTooltip,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Branches.PushTag),
			Handler:           self.withItems(self.push),
			GetDisabledReason: self.require(self.itemsSelected()),
			Description:       self.c.Tr.PushTag,
			Tooltip:           self.c.Tr.PushTagTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Commits.ViewResetOptions),
			Handler:           self.withItem(self.createResetMenu),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Reset,
			Tooltip:           self.c.Tr.ResetTooltip,
			DisplayOnScreen:   true,
			OpensMenu:         true,
		},
		{
			Keys: opts.GetKeys(opts.Config.Universal.OpenDiffTool),
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
				task = types.NewRunCommandTaskWithPrefix(cmdObj.GetCmd(), types.StaticPrefix(prefix))
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

func (self *TagsController) localDelete(tags []*models.Tag) error {
	return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.DeleteLocalTag)
		err := self.c.Git().Tag.LocalDelete(tagNames(tags))
		self.refreshAfterLocalDelete()
		return err
	})
}

func (self *TagsController) remoteDelete(tags []*models.Tag) error {
	confirmPromptTemplate := lo.Ternary(len(tags) > 1, self.c.Tr.DeleteRemoteTagsPrompt, self.c.Tr.DeleteRemoteTagPrompt)
	return self.confirmRemoteDelete(tags, confirmPromptTemplate, func(task gocui.Task, upstream string) error {
		self.c.LogAction(self.c.Tr.Actions.DeleteRemoteTag)
		if err := self.c.Git().Remote.DeleteRemoteTag(task, upstream, tagNames(tags)); err != nil {
			return err
		}
		self.c.Toast(lo.Ternary(len(tags) > 1, self.c.Tr.RemoteTagsDeletedMessage, self.c.Tr.RemoteTagDeletedMessage))
		self.c.RefreshFromWorker(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS, types.TAGS}})
		return nil
	})
}

func (self *TagsController) localAndRemoteDelete(tags []*models.Tag) error {
	confirmPromptTemplate := lo.Ternary(len(tags) > 1, self.c.Tr.DeleteLocalAndRemoteTagsPrompt, self.c.Tr.DeleteLocalAndRemoteTagPrompt)
	return self.confirmRemoteDelete(tags, confirmPromptTemplate, func(task gocui.Task, upstream string) error {
		self.c.LogAction(self.c.Tr.Actions.DeleteRemoteTag)
		if err := self.c.Git().Remote.DeleteRemoteTag(task, upstream, tagNames(tags)); err != nil {
			return err
		}

		self.c.LogAction(self.c.Tr.Actions.DeleteLocalTag)
		if err := self.c.Git().Tag.LocalDelete(tagNames(tags)); err != nil {
			return err
		}
		self.refreshAfterLocalDelete()
		return nil
	})
}

// Refreshes after deleting local tags. If the tags were a range selection,
// this also collapses it to its first line; otherwise it would select the tags
// that moved up into the place of the deleted ones. Collapsing it in the Then
// of a batched refresh draws the shorter list and the new selection in the
// same frame.
func (self *TagsController) refreshAfterLocalDelete() {
	self.c.RefreshFromWorker(types.RefreshOptions{
		Scope:          []types.RefreshableView{types.COMMITS, types.TAGS},
		BatchUIUpdates: true,
		Then: func() error {
			if self.context().IsSelectingRange() {
				self.context().CollapseRangeSelectionToTop()
				self.c.PostRefreshUpdate(self.context())
			}
			return nil
		},
	})
}

// Asks for the remote to delete the tags from and for a confirmation, and then
// runs deleteTags on a worker with the tags shown as being deleted.
// confirmPromptTemplate can use the placeholder upstream, and tagName if there
// is only one tag.
func (self *TagsController) confirmRemoteDelete(
	tags []*models.Tag,
	confirmPromptTemplate string,
	deleteTags func(task gocui.Task, upstream string) error,
) error {
	var title string
	if len(tags) == 1 {
		title = utils.ResolvePlaceholderString(
			self.c.Tr.SelectRemoteTagUpstream,
			map[string]string{
				"tagName": tags[0].Name,
			},
		)
	} else {
		title = self.c.Tr.SelectRemoteTagsUpstream
	}

	self.c.Prompt(types.PromptOpts{
		Title:               title,
		InitialContent:      "origin",
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteSuggestionsFunc(),
		HandleConfirm: func(upstream string) error {
			confirmPrompt := utils.ResolvePlaceholderString(
				confirmPromptTemplate,
				map[string]string{
					"tagName":  tags[0].Name,
					"upstream": upstream,
				},
			)

			self.c.Confirm(types.ConfirmOpts{
				Title:  self.deleteTagsTitle(tags),
				Prompt: confirmPrompt,
				HandleConfirm: func() error {
					return helpers.WithInlineStatusOnItems(self.c.HelperCommon, tags, types.ItemOperationDeleting, context.TAGS_CONTEXT_KEY, func(task gocui.Task) error {
						return deleteTags(task, upstream)
					})
				},
			})

			return nil
		},
	})

	return nil
}

func (self *TagsController) delete(tags []*models.Tag) error {
	menuItems := []*types.MenuItem{
		{
			Label: lo.Ternary(len(tags) > 1, self.c.Tr.DeleteLocalTags, self.c.Tr.DeleteLocalTag),
			Keys:  menuKey('c'),
			OnPress: func() error {
				return self.localDelete(tags)
			},
		},
		{
			Label:     lo.Ternary(len(tags) > 1, self.c.Tr.DeleteRemoteTags, self.c.Tr.DeleteRemoteTag),
			Keys:      menuKey('r'),
			OpensMenu: true,
			OnPress: func() error {
				return self.remoteDelete(tags)
			},
		},
		{
			Label:     lo.Ternary(len(tags) > 1, self.c.Tr.DeleteLocalAndRemoteTags, self.c.Tr.DeleteLocalAndRemoteTag),
			Keys:      menuKey('b'),
			OpensMenu: true,
			OnPress: func() error {
				return self.localAndRemoteDelete(tags)
			},
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.deleteTagsTitle(tags),
		Items: menuItems,
	})
}

func (self *TagsController) deleteTagsTitle(tags []*models.Tag) string {
	if len(tags) > 1 {
		return self.c.Tr.DeleteTagsTitle
	}

	return utils.ResolvePlaceholderString(
		self.c.Tr.DeleteTagTitle,
		map[string]string{
			"tagName": tags[0].Name,
		},
	)
}

func (self *TagsController) push(tags []*models.Tag) error {
	var title string
	if len(tags) == 1 {
		title = utils.ResolvePlaceholderString(
			self.c.Tr.PushTagTitle,
			map[string]string{
				"tagName": tags[0].Name,
			},
		)
	} else {
		title = self.c.Tr.PushTagsTitle
	}

	self.c.Prompt(types.PromptOpts{
		Title:               title,
		InitialContent:      "origin",
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			return helpers.WithInlineStatusOnItems(self.c.HelperCommon, tags, types.ItemOperationPushing, context.TAGS_CONTEXT_KEY, func(task gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.PushTag)
				return self.c.Git().Tag.Push(task, response, tagNames(tags))
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

func tagNames(tags []*models.Tag) []string {
	return lo.Map(tags, func(tag *models.Tag, _ int) string { return tag.Name })
}
