package controllers

import "github.com/jesseduffield/lazygit/pkg/gui/types"

// Embed this into your list controller to get some convenience methods for
// ensuring a single item is selected, etc.

type ListControllerTrait[T comparable] struct {
	c                *ControllerCommon
	context          types.IListContext
	getSelectedItem  func() T
	getSelectedItems func() ([]T, int, int)
}

func NewListControllerTrait[T comparable](
	c *ControllerCommon,
	context types.IListContext,
	getSelected func() T,
	getSelectedItems func() ([]T, int, int),
) *ListControllerTrait[T] {
	return &ListControllerTrait[T]{
		c:                c,
		context:          context,
		getSelectedItem:  getSelected,
		getSelectedItems: getSelectedItems,
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
		item := self.getSelectedItem()
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

// Ensures that at least one item is selected.
func (self *ListControllerTrait[T]) itemRangeSelected(callbacks ...func([]T, int, int) *types.DisabledReason) func() *types.DisabledReason {
	return func() *types.DisabledReason {
		items, startIdx, endIdx := self.getSelectedItems()
		if len(items) == 0 {
			return &types.DisabledReason{Text: self.c.Tr.NoItemSelected}
		}

		for _, callback := range callbacks {
			if reason := callback(items, startIdx, endIdx); reason != nil {
				return reason
			}
		}

		return nil
	}
}

func (self *ListControllerTrait[T]) itemsSelected(callbacks ...func([]T) *types.DisabledReason) func() *types.DisabledReason { //nolint:unused
	return func() *types.DisabledReason {
		items, _, _ := self.getSelectedItems()
		if len(items) == 0 {
			return &types.DisabledReason{Text: self.c.Tr.NoItemSelected}
		}

		for _, callback := range callbacks {
			if reason := callback(items); reason != nil {
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
		commit := self.getSelectedItem()
		if commit == zeroValue {
			return self.c.ErrorMsg(self.c.Tr.NoItemSelected)
		}

		return callback(commit)
	}
}

func (self *ListControllerTrait[T]) withItems(callback func([]T) error) func() error {
	return func() error {
		items, _, _ := self.getSelectedItems()
		if len(items) == 0 {
			return self.c.ErrorMsg(self.c.Tr.NoItemSelected)
		}

		return callback(items)
	}
}

// like withItems but also passes the start and end index of the selection
func (self *ListControllerTrait[T]) withItemsRange(callback func([]T, int, int) error) func() error {
	return func() error {
		items, startIdx, endIdx := self.getSelectedItems()
		if len(items) == 0 {
			return self.c.ErrorMsg(self.c.Tr.NoItemSelected)
		}

		return callback(items, startIdx, endIdx)
	}
}

// Like withItem, but doesn't show an error message if no item is selected.
// Use this for click actions (it's a no-op to click empty space)
func (self *ListControllerTrait[T]) withItemGraceful(callback func(T) error) func() error {
	return func() error {
		var zeroValue T
		commit := self.getSelectedItem()
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
