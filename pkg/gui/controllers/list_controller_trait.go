package controllers

import "github.com/jesseduffield/lazygit/pkg/gui/types"

// Embed this into your list controller to get some convenience methods for
// ensuring a single item is selected, etc.

type ListControllerTrait[T comparable] struct {
	c           *ControllerCommon
	context     types.IListContext
	getSelected func() T
}

func NewListControllerTrait[T comparable](
	c *ControllerCommon,
	context types.IListContext,
	getSelected func() T,
) *ListControllerTrait[T] {
	return &ListControllerTrait[T]{
		c:           c,
		context:     context,
		getSelected: getSelected,
	}
}

// Convenience function for combining multiple disabledReason callbacks.
// The first callback to return a disabled reason will be the one returned.
func (self *ListControllerTrait[T]) require(callbacks ...func() *types.DisabledReason) func() *types.DisabledReason {
	return func() *types.DisabledReason {
		for _, callback := range callbacks {
			if disabledReason := callback(); disabledReason != nil {
				return disabledReason
			}
		}

		return nil
	}
}

// Convenience function for enforcing that a single item is selected.
// Also takes callbacks for additional disabled reasons, and passes the selected
// item into each one.
func (self *ListControllerTrait[T]) singleItemSelected(callbacks ...func(T) *types.DisabledReason) func() *types.DisabledReason {
	return func() *types.DisabledReason {
		if self.context.GetList().AreMultipleItemsSelected() {
			return &types.DisabledReason{Text: self.c.Tr.RangeSelectNotSupported}
		}

		var zeroValue T
		item := self.getSelected()
		if item == zeroValue {
			return &types.DisabledReason{Text: self.c.Tr.NoItemSelected}
		}

		for _, callback := range callbacks {
			if reason := callback(item); reason != nil {
				return reason
			}
		}

		return nil
	}
}

// Passes the selected item to the callback. Used for handler functions.
func (self *ListControllerTrait[T]) withItem(callback func(T) error) func() error {
	return func() error {
		var zeroValue T
		commit := self.getSelected()
		if commit == zeroValue {
			return self.c.ErrorMsg(self.c.Tr.NoItemSelected)
		}

		return callback(commit)
	}
}

// Like withItem, but doesn't show an error message if no item is selected.
// Use this for click actions (it's a no-op to click empty space)
func (self *ListControllerTrait[T]) withItemGraceful(callback func(T) error) func() error {
	return func() error {
		var zeroValue T
		commit := self.getSelected()
		if commit == zeroValue {
			return nil
		}

		return callback(commit)
	}
}

// All controllers must implement this method so we're defining it here for convenience
func (self *ListControllerTrait[T]) Context() types.Context {
	return self.context
}
