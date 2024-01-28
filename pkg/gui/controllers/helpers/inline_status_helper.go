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
		self.c.OnWorker(func(task gocui.Task) {
			self.start(opts)

			err := f(inlineStatusHelperTask{task, self, opts})
			if err != nil {
				self.c.OnUIThread(func() error {
					return self.c.Error(err)
				})
			}

			self.stop(opts)
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
			ticker := time.NewTicker(time.Millisecond * utils.LoaderAnimationInterval)
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
}

func (self *InlineStatusHelper) renderContext(contextKey types.ContextKey) {
	self.c.OnUIThread(func() error {
		_ = self.c.ContextForKey(contextKey).HandleRender()
		return nil
	})
}
