package lo

import (
	"context"
	"sync"
	"time"

	"github.com/samber/lo/internal/xrand"
)

// DispatchingStrategy is a function that distributes messages to channels.
type DispatchingStrategy[T any] func(msg T, index uint64, channels []<-chan T) int

// ChannelDispatcher distributes messages from input channels into N child channels.
// Close events are propagated to children.
// Underlying channels can have a fixed buffer capacity or be unbuffered when cap is 0.
// Play: https://go.dev/play/p/UZGu2wVg3J2
func ChannelDispatcher[T any](stream <-chan T, count, channelBufferCap int, strategy DispatchingStrategy[T]) []<-chan T {
	children := createChannels[T](count, channelBufferCap)

	roChildren := channelsToReadOnly(children)

	go func() {
		// propagate channel closing to children
		defer closeChannels(children)

		var i uint64

		for msg := range stream {
			destination := strategy(msg, i, roChildren) % count
			children[destination] <- msg

			i++
		}
	}()

	return roChildren
}

func createChannels[T any](count, channelBufferCap int) []chan T {
	children := make([]chan T, 0, count)

	for i := 0; i < count; i++ {
		children = append(children, make(chan T, channelBufferCap))
	}

	return children
}

func channelsToReadOnly[T any](children []chan T) []<-chan T {
	roChildren := make([]<-chan T, 0, len(children))

	for i := range children {
		roChildren = append(roChildren, children[i])
	}

	return roChildren
}

func closeChannels[T any](children []chan T) {
	for i := 0; i < len(children); i++ {
		close(children[i])
	}
}

func channelIsNotFull[T any](ch <-chan T) bool {
	return cap(ch) == 0 || len(ch) < cap(ch)
}

// DispatchingStrategyRoundRobin distributes messages in a rotating sequential manner.
// If the channel capacity is exceeded, the next channel will be selected and so on.
// Play: https://go.dev/play/p/UZGu2wVg3J2
func DispatchingStrategyRoundRobin[T any](msg T, index uint64, channels []<-chan T) int {
	for {
		i := int(index % uint64(len(channels)))
		if channelIsNotFull(channels[i]) {
			return i
		}

		index++
		time.Sleep(10 * time.Microsecond) // prevent CPU from burning 🔥
	}
}

// DispatchingStrategyRandom distributes messages in a random manner.
// If the channel capacity is exceeded, another random channel will be selected and so on.
// Play: https://go.dev/play/p/GEyGn3TdGk4
func DispatchingStrategyRandom[T any](msg T, index uint64, channels []<-chan T) int {
	for {
		i := xrand.IntN(len(channels))
		if channelIsNotFull(channels[i]) {
			return i
		}

		time.Sleep(10 * time.Microsecond) // prevent CPU from burning 🔥
	}
}

// DispatchingStrategyWeightedRandom distributes messages in a weighted manner.
// If the channel capacity is exceeded, another random channel will be selected and so on.
// Play: https://go.dev/play/p/v0eMh8NZG2L
func DispatchingStrategyWeightedRandom[T any](weights []int) DispatchingStrategy[T] {
	seq := []int{}

	for i, weight := range weights {
		for j := 0; j < weight; j++ {
			seq = append(seq, i)
		}
	}

	return func(msg T, index uint64, channels []<-chan T) int {
		for {
			i := seq[xrand.IntN(len(seq))]
			if channelIsNotFull(channels[i]) {
				return i
			}

			time.Sleep(10 * time.Microsecond) // prevent CPU from burning 🔥
		}
	}
}

// DispatchingStrategyFirst distributes messages in the first non-full channel.
// If the capacity of the first channel is exceeded, the second channel will be selected and so on.
// Play: https://go.dev/play/p/OrJCvOmk42f
func DispatchingStrategyFirst[T any](msg T, index uint64, channels []<-chan T) int {
	for {
		for i := range channels {
			if channelIsNotFull(channels[i]) {
				return i
			}
		}

		time.Sleep(10 * time.Microsecond) // prevent CPU from burning 🔥
	}
}

// DispatchingStrategyLeast distributes messages in the emptiest channel.
// Play: https://go.dev/play/p/ypy0jrRcEe7
func DispatchingStrategyLeast[T any](msg T, index uint64, channels []<-chan T) int {
	_, i := MinIndexBy(channels, func(a, b <-chan T) bool {
		return len(a) < len(b)
	})

	return i
}

// DispatchingStrategyMost distributes messages in the fullest channel.
// If the channel capacity is exceeded, the next channel will be selected and so on.
// Play: https://go.dev/play/p/erHHone7rF9
func DispatchingStrategyMost[T any](msg T, index uint64, channels []<-chan T) int {
	_, i := MaxIndexBy(channels, func(a, b <-chan T) bool {
		return len(a) > len(b) && channelIsNotFull(a)
	})

	return i
}

// SliceToChannel returns a read-only channel of collection elements.
// Play: https://go.dev/play/p/lIbSY3QmiEg
func SliceToChannel[T any](bufferSize int, collection []T) <-chan T {
	ch := make(chan T, bufferSize)

	go func() {
		for i := range collection {
			ch <- collection[i]
		}

		close(ch)
	}()

	return ch
}

// ChannelToSlice returns a slice built from channel items. Blocks until channel closes.
// Play: https://go.dev/play/p/lIbSY3QmiEg
func ChannelToSlice[T any](ch <-chan T) []T {
	collection := []T{}

	for item := range ch {
		collection = append(collection, item)
	}

	return collection
}

// Generator implements the generator design pattern.
// Play: https://go.dev/play/p/lIbSY3QmiEg
//
// Deprecated: use "iter" package instead (Go >= 1.23).
func Generator[T any](bufferSize int, generator func(yield func(T))) <-chan T {
	ch := make(chan T, bufferSize)

	go func() {
		// WARNING: infinite loop
		generator(func(t T) {
			ch <- t
		})

		close(ch)
	}()

	return ch
}

// Buffer creates a slice of n elements from a channel. Returns the slice and the slice length.
// @TODO: we should probably provide a helper that reuses the same buffer.
// Play: https://go.dev/play/p/gPQ-6xmcKQI
func Buffer[T any](ch <-chan T, size int) (collection []T, length int, readTime time.Duration, ok bool) {
	buffer := make([]T, 0, size)
	now := time.Now()

	for index := 0; index < size; index++ {
		item, ok := <-ch
		if !ok {
			return buffer, index, time.Since(now), false
		}

		buffer = append(buffer, item)
	}

	return buffer, size, time.Since(now), true
}

// BufferWithContext creates a slice of n elements from a channel, with context. Returns the slice and the slice length.
// @TODO: we should probably provide a helper that reuses the same buffer.
// Play: https://go.dev/play/p/oRfOyJWK9YF
func BufferWithContext[T any](ctx context.Context, ch <-chan T, size int) (collection []T, length int, readTime time.Duration, ok bool) {
	buffer := make([]T, 0, size)
	now := time.Now()

	for index := 0; index < size; index++ {
		select {
		case item, ok := <-ch:
			if !ok {
				return buffer, index, time.Since(now), false
			}

			buffer = append(buffer, item)

		case <-ctx.Done():
			return buffer, index, time.Since(now), true
		}
	}

	return buffer, size, time.Since(now), true
}

// BufferWithTimeout creates a slice of n elements from a channel, with timeout. Returns the slice and the slice length.
// Play: https://go.dev/play/p/sxyEM3koo4n
func BufferWithTimeout[T any](ch <-chan T, size int, timeout time.Duration) (collection []T, length int, readTime time.Duration, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return BufferWithContext(ctx, ch, size)
}

// FanIn collects messages from multiple input channels into a single buffered channel.
// Output messages have no priority. When all upstream channels reach EOF, downstream channel closes.
// Play: https://go.dev/play/p/FH8Wq-T04Jb
func FanIn[T any](channelBufferCap int, upstreams ...<-chan T) <-chan T {
	out := make(chan T, channelBufferCap)
	var wg sync.WaitGroup

	// Start an output goroutine for each input channel in upstreams.
	wg.Add(len(upstreams))
	for i := range upstreams {
		go func(index int) {
			for n := range upstreams[index] {
				out <- n
			}
			wg.Done()
		}(i)
	}

	// Start a goroutine to close out once all the output goroutines are done.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// FanOut broadcasts all the upstream messages to multiple downstream channels.
// When upstream channel reaches EOF, downstream channels close. If any downstream
// channels is full, broadcasting is paused.
// Play: https://go.dev/play/p/2LHxcjKX23L
func FanOut[T any](count, channelsBufferCap int, upstream <-chan T) []<-chan T {
	downstreams := createChannels[T](count, channelsBufferCap)

	go func() {
		for msg := range upstream {
			for i := range downstreams {
				downstreams[i] <- msg
			}
		}

		// Close out once all the output goroutines are done.
		for i := range downstreams {
			close(downstreams[i])
		}
	}()

	return channelsToReadOnly(downstreams)
}
