package http

import (
	"net/http"
	"net/url"
)

// transportOptions contains transport specific configuration.
type transportOptions struct {
	insecureSkipTLS bool
	// []byte is not comparable.
	caBundle string
	proxyURL url.URL
}

func (c *client) addTransport(opts transportOptions, transport *http.Transport) {
	c.mutex.Lock()
	c.transports.Add(opts, transport)
	c.mutex.Unlock()
}

func (c *client) removeTransport(opts transportOptions) {
	c.mutex.Lock()
	c.transports.Remove(opts)
	c.mutex.Unlock()
}

func (c *client) fetchTransport(opts transportOptions) (*http.Transport, bool) {
	c.mutex.RLock()
	t, ok := c.transports.Get(opts)
	c.mutex.RUnlock()
	if !ok {
		return nil, false
	}
	transport, ok := t.(*http.Transport)
	if !ok {
		return nil, false
	}
	return transport, true
}
