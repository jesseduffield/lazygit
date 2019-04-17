package getter

import (
	"context"
	"io"
)

// readerFunc is syntactic sugar for read interface.
type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

// Copy is a io.Copy cancellable by context
func Copy(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	// Copy will call the Reader and Writer interface multiple time, in order
	// to copy by chunk (avoiding loading the whole file in memory).
	return io.Copy(dst, readerFunc(func(p []byte) (int, error) {

		select {
		case <-ctx.Done():
			// context has been canceled
			// stop process and propagate "context canceled" error
			return 0, ctx.Err()
		default:
			// otherwise just run default io.Reader implementation
			return src.Read(p)
		}
	}))
}
