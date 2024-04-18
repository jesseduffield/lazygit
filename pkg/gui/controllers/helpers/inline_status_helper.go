package helpers

import (
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sasha-s/go-deadlock"
)

type InlineStatusHelper struct {
	c *HelperCommon

	windowHelper             *WindowHelper
	contextsWithInlineStatus map[types.ContextKey]*inlineStatusInfo
	mutex                    *deadlock.Mutex
}

func NewInlineStatusHelper(c *HelperCommon, windowHelper *WindowHelper) *InlineStatusHelper {
	return &InlineStatusHelper{
		c:                        c,
		windowHelper:             windowHelper,
		contextsWithInlineStatus: make(map[types.ContextKey]*inlineStatusInfo),
		mutex:                    &deadlock.Mutex{},
	}
}

type InlineStatusOpts struct {
	Item       types.HasUrn
	Operation  types.ItemOperation
	ContextKey types.ContextKey
}

type inlineStatusInfo struct {
	refCount int
	stop     chan struct{}
}

// A custom task for WithInlineStatus calls; it wraps the original one and
// hides the status whenever the task is paused, and shows it again when
// continued.
type inlineStatusHelperTask struct {
	gocui.Task

	inlineStatusHelper *InlineStatusHelper
	opts               InlineStatusOpts
}

// poor man's version of explicitly saying that struct X implements interface Y
var _ gocui.Task = inlineStatusHelperTask{}

func (self inlineStatusHelperTask) Pause() {
	self.inlineStatusHelper.stop(self.opts)
	self.Task.Pause()

	self.inlineStatusHelper.renderContext(self.opts.ContextKey)
}

func (self inlineStatusHelperTask) Continue() {
	self.Task.Continue()
	self.inlineStatusHelper.start(self.opts)
}

func (self *InlineStatusHelper) WithInlineStatus(opts InlineStatusOpts, f func(gocui.Task) error) {
	context := self.c.ContextForKey(opts.ContextKey).(types.IListContext)
	view := context.GetView()
	visible := view.Visible && self.windowHelper.TopViewInWindow(context.GetWindowName(), false) == view
	if visible && context.IsItemVisible(opts.Item) {
		self.c.OnWorker(func(task gocui.Task) error {
			self.start(opts)
			defer self.stop(opts)

			return f(inlineStatusHelperTask{task, self, opts})
		})
	} else {
		message := presentation.ItemOperationToString(opts.Operation, self.c.Tr)
		_ = self.c.WithWaitingStatus(message, func(t gocui.Task) error {
			// We still need to set the item operation, because it might be used
			// for other (non-presentation) purposes
			self.c.State().SetItemOperation(opts.Item, opts.Operation)
			defer self.c.State().ClearItemOperation(opts.Item)

			return f(t)
		})
	}
}

func (self *InlineStatusHelper) start(opts InlineStatusOpts) {
	self.c.State().SetItemOperation(opts.Item, opts.Operation)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	info := self.contextsWithInlineStatus[opts.ContextKey]
	if info == nil {
		info = &inlineStatusInfo{refCount: 0, stop: make(chan struct{})}
		self.contextsWithInlineStatus[opts.ContextKey] = info

		go utils.Safe(func() {
			ticker := time.NewTicker(time.Millisecond * time.Duration(self.c.UserConfig.Gui.Spinner.Rate))
			defer ticker.Stop()
		outer:
			for {
				select {
				case <-ticker.C:
					self.renderContext(opts.ContextKey)
				case <-info.stop:
					break outer
				}
			}
		})
	}

	info.refCount++
}

func (self *InlineStatusHelper) stop(opts InlineStatusOpts) {
	self.mutex.Lock()

	if info := self.contextsWithInlineStatus[opts.ContextKey]; info != nil {
		info.refCount--
		if info.refCount <= 0 {
			info.stop <- struct{}{}
			delete(self.contextsWithInlineStatus, opts.ContextKey)
		}
	}

	self.mutex.Unlock()

	self.c.State().ClearItemOperation(opts.Item)

	// When recording a demo we need to re-render the context again here to
	// remove the inline status. In normal usage we don't want to do this
	// because in the case of pushing a branch this would first reveal the ↑3↓7
	// status from before the push for a brief moment, to be replaced by a green
	// checkmark a moment later when the async refresh is done. This looks
	// jarring, so normally we rely on the async refresh to redraw with the
	// status removed. (In some rare cases, where there's no refresh at all, we
	// need to redraw manually in the controller; see TagsController.push() for
	// an example.)
	//
	// In demos, however, we turn all async refreshes into sync ones, because
	// this looks better in demos. In this case the refresh happens while the
	// status is still set, so we need to render again after removing it.
	if self.c.InDemo() {
		self.renderContext(opts.ContextKey)
	}
}

func (self *InlineStatusHelper) renderContext(contextKey types.ContextKey) {
	self.c.OnUIThread(func() error {
		_ = self.c.ContextForKey(contextKey).HandleRender()
		return nil
	})
}
