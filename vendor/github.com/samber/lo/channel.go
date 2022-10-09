package lo

import (
	"math/rand"
	"time"
)

type DispatchingStrategy[T any] func(msg T, index uint64, channels []<-chan T) int

// ChannelDispatcher distributes messages from input channels into N child channels.
// Close events are propagated to children.
// Underlying channels can have a fixed buffer capacity or be unbuffered when cap is 0.
func ChannelDispatcher[T any](stream <-chan T, count int, channelBufferCap int, strategy DispatchingStrategy[T]) []<-chan T {
	children := createChannels[T](count, channelBufferCap)

	roChildren := channelsToReadOnly(children)

	go func() {
		// propagate channel closing to children
		defer closeChannels(children)

		var i uint64 = 0

		for {
			msg, ok := <-stream
			if !ok {
				return
			}

			destination := strategy(msg, i, roChildren) % count
			children[destination] <- msg

			i++
		}
	}()

	return roChildren
}

func createChannels[T any](count int, channelBufferCap int) []chan T {
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
func DispatchingStrategyRandom[T any](msg T, index uint64, channels []<-chan T) int {
	for {
		i := rand.Intn(len(channels))
		if channelIsNotFull(channels[i]) {
			return i
		}

		time.Sleep(10 * time.Microsecond) // prevent CPU from burning 🔥
	}
}

// DispatchingStrategyRandom distributes messages in a weighted manner.
// If the channel capacity is exceeded, another random channel will be selected and so on.
func DispatchingStrategyWeightedRandom[T any](weights []int) DispatchingStrategy[T] {
	seq := []int{}

	for i := 0; i < len(weights); i++ {
		for j := 0; j < weights[i]; j++ {
			seq = append(seq, i)
		}
	}

	return func(msg T, index uint64, channels []<-chan T) int {
		for {
			i := seq[rand.Intn(len(seq))]
			if channelIsNotFull(channels[i]) {
				return i
			}

			time.Sleep(10 * time.Microsecond) // prevent CPU from burning 🔥
		}
	}
}

// DispatchingStrategyFirst distributes messages in the first non-full channel.
// If the capacity of the first channel is exceeded, the second channel will be selected and so on.
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
func DispatchingStrategyLeast[T any](msg T, index uint64, channels []<-chan T) int {
	seq := Range(len(channels))

	return MinBy(seq, func(item int, min int) bool {
		return len(channels[item]) < len(channels[min])
	})
}

// DispatchingStrategyMost distributes messages in the fulliest channel.
// If the channel capacity is exceeded, the next channel will be selected and so on.
func DispatchingStrategyMost[T any](msg T, index uint64, channels []<-chan T) int {
	seq := Range(len(channels))

	return MaxBy(seq, func(item int, max int) bool {
		return len(channels[item]) > len(channels[max]) && channelIsNotFull(channels[item])
	})
}

// SliceToChannel returns a read-only channels of collection elements.
func SliceToChannel[T any](bufferSize int, collection []T) <-chan T {
	ch := make(chan T, bufferSize)

	go func() {
		for _, item := range collection {
			ch <- item
		}

		close(ch)
	}()

	return ch
}

// Generator implements the generator design pattern.
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

// Batch creates a slice of n elements from a channel. Returns the slice and the slice length.
// @TODO: we should probaby provide an helper that reuse the same buffer.
func Batch[T any](ch <-chan T, size int) (collection []T, length int, readTime time.Duration, ok bool) {
	buffer := make([]T, 0, size)
	index := 0
	now := time.Now()

	for ; index < size; index++ {
		item, ok := <-ch
		if !ok {
			return buffer, index, time.Since(now), false
		}

		buffer = append(buffer, item)
	}

	return buffer, index, time.Since(now), true
}

// BatchWithTimeout creates a slice of n elements from a channel, with timeout. Returns the slice and the slice length.
// @TODO: we should probaby provide an helper that reuse the same buffer.
func BatchWithTimeout[T any](ch <-chan T, size int, timeout time.Duration) (collection []T, length int, readTime time.Duration, ok bool) {
	expire := time.NewTimer(timeout)
	defer expire.Stop()

	buffer := make([]T, 0, size)
	index := 0
	now := time.Now()

	for ; index < size; index++ {
		select {
		case item, ok := <-ch:
			if !ok {
				return buffer, index, time.Since(now), false
			}

			buffer = append(buffer, item)

		case <-expire.C:
			return buffer, index, time.Since(now), true
		}
	}

	return buffer, index, time.Since(now), true
}
