package gui

type BasicContext struct {
	OnFocus         func() error
	OnFocusLost     func() error
	OnRender        func() error
	Kind            ContextKind
	Key             ContextKey
	ViewName        string
	WindowName      string
	OnGetOptionsMap func() map[string]string

	ParentContext Context
	// we can't know on the calling end whether a Context is actually a nil value without reflection, so we're storing this flag here to tell us. There has got to be a better way around this
	hasParent bool
}

func (c *BasicContext) GetOptionsMap() map[string]string {
	if c.OnGetOptionsMap != nil {
		return c.OnGetOptionsMap()
	}
	return nil
}

func (c *BasicContext) SetParentContext(context Context) {
	c.ParentContext = context
	c.hasParent = true
}

func (c *BasicContext) GetParentContext() (Context, bool) {
	return c.ParentContext, c.hasParent
}

func (c *BasicContext) SetWindowName(windowName string) {
	c.WindowName = windowName
}

func (c *BasicContext) GetWindowName() string {
	windowName := c.WindowName

	if windowName != "" {
		return windowName
	}

	// TODO: actually set this for everything so we don't default to the view name
	return c.ViewName
}

func (c *BasicContext) HandleRender() error {
	if c.OnRender != nil {
		return c.OnRender()
	}
	return nil
}

func (c *BasicContext) GetViewName() string {
	return c.ViewName
}

func (c *BasicContext) HandleFocus() error {
	return c.OnFocus()
}

func (c *BasicContext) HandleFocusLost() error {
	if c.OnFocusLost != nil {
		return c.OnFocusLost()
	}
	return nil
}

func (c *BasicContext) GetKind() ContextKind {
	return c.Kind
}

func (c *BasicContext) GetKey() ContextKey {
	return c.Key
}
