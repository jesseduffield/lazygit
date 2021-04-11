package gui

type ListContext struct {
	GetItemsLength      func() int
	GetDisplayStrings   func() [][]string
	OnFocus             func() error
	OnFocusLost         func() error
	OnClickSelectedItem func() error

	// the boolean here tells us whether the item is nil. This is needed because you can't work it out on the calling end once the pointer is wrapped in an interface (unless you want to use reflection)
	SelectedItem  func() (ListItem, bool)
	GetPanelState func() IListPanelState

	Gui                        *Gui
	ResetMainViewOriginOnFocus bool

	*BasicContext
}

type IListPanelState interface {
	SetSelectedLineIdx(int)
	GetSelectedLineIdx() int
}

type ListItem interface {
	// ID is a SHA when the item is a commit, a filename when the item is a file, 'stash@{4}' when it's a stash entry, 'my_branch' when it's a branch
	ID() string

	// Description is something we would show in a message e.g. '123as14: push blah' for a commit
	Description() string
}

func (lc *ListContext) GetSelectedItem() (ListItem, bool) {
	return lc.SelectedItem()
}

func (lc *ListContext) GetSelectedItemId() string {
	item, ok := lc.SelectedItem()

	if !ok {
		return ""
	}

	return item.ID()
}

// OnFocus assumes that the content of the context has already been rendered to the view. OnRender is the function which actually renders the content to the view
func (lc *ListContext) OnRender() error {
	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return nil
	}

	if lc.GetDisplayStrings != nil {
		lc.Gui.refreshSelectedLine(lc.GetPanelState(), lc.GetItemsLength())
		lc.Gui.renderDisplayStrings(view, lc.GetDisplayStrings())
	}

	return nil
}

func (lc *ListContext) HandleFocusLost() error {
	if lc.OnFocusLost != nil {
		return lc.OnFocusLost()
	}

	return nil
}

func (lc *ListContext) HandleFocus() error {
	if lc.Gui.popupPanelFocused() {
		return nil
	}

	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return nil
	}

	view.FocusPoint(0, lc.GetPanelState().GetSelectedLineIdx())

	if lc.ResetMainViewOriginOnFocus {
		if err := lc.Gui.resetOrigin(lc.Gui.Views.Main); err != nil {
			return err
		}
		if err := lc.Gui.resetOrigin(lc.Gui.Views.Secondary); err != nil {
			return err
		}
	}

	if lc.Gui.State.Modes.Diffing.Active() {
		return lc.Gui.renderDiff()
	}

	if lc.OnFocus != nil {
		return lc.OnFocus()
	}

	return nil
}

func (lc *ListContext) HandleRender() error {
	return lc.OnRender()
}

func (lc *ListContext) handlePrevLine() error {
	return lc.handleLineChange(-1)
}

func (lc *ListContext) handleNextLine() error {
	return lc.handleLineChange(1)
}

func (lc *ListContext) handleLineChange(change int) error {
	if !lc.Gui.isPopupPanel(lc.ViewName) && lc.Gui.popupPanelFocused() {
		return nil
	}

	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return err
	}

	selectedLineIdx := lc.GetPanelState().GetSelectedLineIdx()
	if (change < 0 && selectedLineIdx == 0) || (change > 0 && selectedLineIdx == lc.GetItemsLength()-1) {
		return nil
	}

	lc.Gui.changeSelectedLine(lc.GetPanelState(), lc.GetItemsLength(), change)
	view.FocusPoint(0, lc.GetPanelState().GetSelectedLineIdx())

	return lc.HandleFocus()
}

func (lc *ListContext) handleNextPage() error {
	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return nil
	}
	delta := lc.Gui.pageDelta(view)

	return lc.handleLineChange(delta)
}

func (lc *ListContext) handleGotoTop() error {
	return lc.handleLineChange(-lc.GetItemsLength())
}

func (lc *ListContext) handleGotoBottom() error {
	return lc.handleLineChange(lc.GetItemsLength())
}

func (lc *ListContext) handlePrevPage() error {
	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return nil
	}

	delta := lc.Gui.pageDelta(view)

	return lc.handleLineChange(-delta)
}

func (lc *ListContext) handleClick() error {
	if !lc.Gui.isPopupPanel(lc.ViewName) && lc.Gui.popupPanelFocused() {
		return nil
	}

	view, err := lc.Gui.g.View(lc.ViewName)
	if err != nil {
		return nil
	}

	prevSelectedLineIdx := lc.GetPanelState().GetSelectedLineIdx()
	newSelectedLineIdx := view.SelectedLineIdx()

	// we need to focus the view
	if err := lc.Gui.pushContext(lc); err != nil {
		return err
	}

	if newSelectedLineIdx > lc.GetItemsLength()-1 {
		return nil
	}

	lc.GetPanelState().SetSelectedLineIdx(newSelectedLineIdx)

	prevViewName := lc.Gui.currentViewName()
	if prevSelectedLineIdx == newSelectedLineIdx && prevViewName == lc.ViewName && lc.OnClickSelectedItem != nil {
		return lc.OnClickSelectedItem()
	}
	return lc.HandleFocus()
}

func (lc *ListContext) onSearchSelect(selectedLineIdx int) error {
	lc.GetPanelState().SetSelectedLineIdx(selectedLineIdx)
	return lc.HandleFocus()
}
