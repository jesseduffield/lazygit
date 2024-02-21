package memoize

import (
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

// Memoizer allows you to memoize function calls. Memoizer is safe for concurrent use by multiple goroutines.
type Memoizer struct {

	// Storage exposes the underlying cache of memoized results to manipulate as desired - for example, to Flush().
	Storage *cache.Cache

	group singleflight.Group
}

// NewMemoizer creates a new Memoizer with the configured expiry and cleanup policies.
// If desired, use cache.NoExpiration to cache values forever.
func NewMemoizer(defaultExpiration, cleanupInterval time.Duration) *Memoizer {
	return &Memoizer{
		Storage: cache.New(defaultExpiration, cleanupInterval),
		group:   singleflight.Group{},
	}
}

// Memoize executes and returns the results of the given function, unless there was a cached value of the same key.
// Only one execution is in-flight for a given key at a time.
// The boolean return value indicates whether v was previously stored.
func (m *Memoizer) Memoize(key string, fn func() (interface{}, error)) (interface{}, error, bool) {
	// Check cache
	value, found := m.Storage.Get(key)
	if found {
		return value, nil, true
	}

	// Combine memoized function with a cache store
	value, err, _ := m.group.Do(key, func() (interface{}, error) {
		data, innerErr := fn()

		if innerErr == nil {
			m.Storage.Set(key, data, cache.DefaultExpiration)
		}

		return data, innerErr
	})
	return value, err, false
}
