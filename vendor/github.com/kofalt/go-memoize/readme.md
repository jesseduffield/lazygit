# go-memoize

There wasn't a decent [memoizer](https://wikipedia.org/wiki/Memoization) for Golang out there, so I lashed two nice libraries together and made one.

Dead-simple. Safe for concurrent use.

[![Reference](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/kofalt/go-memoize)
[![Linter](https://goreportcard.com/badge/github.com/kofalt/go-memoize?style=flat-square)](https://goreportcard.com/report/github.com/kofalt/go-memoize)
[![Build status](https://github.com/kofalt/go-memoize/workflows/Build/badge.svg)](https://github.com/kofalt/go-memoize/actions)

## Project status

**Complete.** Latest commit timestamp might be old - that's okay.

Go-memoize has been in production for a few years, and has yet to burn the house down.

## Usage

Cache expensive function calls in memory, with a configurable timeout and purge interval:

```golang
import (
	"time"

	"github.com/kofalt/go-memoize"
)

// Any expensive call that you wish to cache
expensive := func() (interface{}, error) {
	time.Sleep(3 * time.Second)
	return "some data", nil
}

// Cache expensive calls in memory for 90 seconds, purging old entries every 10 minutes.
cache := memoize.NewMemoizer(90*time.Second, 10*time.Minute)

// This will call the expensive func
result, err, cached := cache.Memoize("key1", expensive)

// This will be cached
result, err, cached = cache.Memoize("key1", expensive)

// This uses a new cache key, so expensive is called again
result, err, cached = cache.Memoize("key2", expensive)
```

In the example above, `result` is:
1. the return value from your function if `cached` is false, or
1. a previously stored value if `cached` is true.

All the hard stuff is punted to patrickmn's [go-cache](https://github.com/patrickmn/go-cache) and the Go team's [x/sync/singleflight](https://github.com/golang/sync), I just lashed them together.

Also note that `cache.Storage` is exported, so you can use the underlying cache features - such as [Flush](https://godoc.org/github.com/patrickmn/go-cache#Cache.Flush) or [SaveFile](https://godoc.org/github.com/patrickmn/go-cache#Cache.SaveFile).
