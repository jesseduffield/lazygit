package commands

// sometimes we need to do a sequence of things in a rebase but the user needs to
// fix merge conflicts along the way. When this happens we queue up the next step
// so that after the next successful rebase continue we can continue from where we left off.
// At the moment we've just got a single callback, but later on we could use a queue here
type RebaseWorkflow struct {
	onSuccessfulContinue func() error
}

func (c *RebaseWorkflow) Start(f func() error) {
	c.onSuccessfulContinue = f
}

func (c *RebaseWorkflow) InProgress() bool {
	return c.onSuccessfulContinue != nil
}

func (c *RebaseWorkflow) Abort() {
	c.onSuccessfulContinue = nil
}

// we could name this function 'Complete' but in future we may actually have a
// queue here rather than a single function so Continue is a more correct name
func (c *RebaseWorkflow) Continue() error {
	f := c.onSuccessfulContinue
	c.onSuccessfulContinue = nil
	return f()
}
