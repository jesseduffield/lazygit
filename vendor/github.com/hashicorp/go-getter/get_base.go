package getter

import "context"

// getter is our base getter; it regroups
// fields all getters have in common.
type getter struct {
	client *Client
}

func (g *getter) SetClient(c *Client) { g.client = c }

// Context tries to returns the Contex from the getter's
// client. otherwise context.Background() is returned.
func (g *getter) Context() context.Context {
	if g == nil || g.client == nil {
		return context.Background()
	}
	return g.client.Ctx
}
