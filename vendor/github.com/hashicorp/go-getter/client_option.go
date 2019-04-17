package getter

import "context"

// A ClientOption allows to configure a client
type ClientOption func(*Client) error

// Configure configures a client with options.
func (c *Client) Configure(opts ...ClientOption) error {
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}
	c.Options = opts
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return err
		}
	}
	// Default decompressor values
	if c.Decompressors == nil {
		c.Decompressors = Decompressors
	}
	// Default detector values
	if c.Detectors == nil {
		c.Detectors = Detectors
	}
	// Default getter values
	if c.Getters == nil {
		c.Getters = Getters
	}

	for _, getter := range c.Getters {
		getter.SetClient(c)
	}
	return nil
}

// WithContext allows to pass a context to operation
// in order to be able to cancel a download in progress.
func WithContext(ctx context.Context) func(*Client) error {
	return func(c *Client) error {
		c.Ctx = ctx
		return nil
	}
}
