package main

import (
	"io"
	"path/filepath"
	"sync"

	"github.com/cheggaaa/pb"
	getter "github.com/hashicorp/go-getter"
)

// defaultProgressBar is the default instance of a cheggaaa
// progress bar.
var defaultProgressBar getter.ProgressTracker = &ProgressBar{}

// ProgressBar wraps a github.com/cheggaaa/pb.Pool
// in order to display download progress for one or multiple
// downloads.
//
// If two different instance of ProgressBar try to
// display a progress only one will be displayed.
// It is therefore recommended to use DefaultProgressBar
type ProgressBar struct {
	// lock everything below
	lock sync.Mutex

	pool *pb.Pool

	pbs int
}

func ProgressBarConfig(bar *pb.ProgressBar, prefix string) {
	bar.SetUnits(pb.U_BYTES)
	bar.Prefix(prefix)
}

// TrackProgress instantiates a new progress bar that will
// display the progress of stream until closed.
// total can be 0.
func (cpb *ProgressBar) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	cpb.lock.Lock()
	defer cpb.lock.Unlock()

	newPb := pb.New64(totalSize)
	newPb.Set64(currentSize)
	ProgressBarConfig(newPb, filepath.Base(src))
	if cpb.pool == nil {
		cpb.pool = pb.NewPool()
		cpb.pool.Start()
	}
	cpb.pool.Add(newPb)
	reader := newPb.NewProxyReader(stream)

	cpb.pbs++
	return &readCloser{
		Reader: reader,
		close: func() error {
			cpb.lock.Lock()
			defer cpb.lock.Unlock()

			newPb.Finish()
			cpb.pbs--
			if cpb.pbs <= 0 {
				cpb.pool.Stop()
				cpb.pool = nil
			}
			return nil
		},
	}
}

type readCloser struct {
	io.Reader
	close func() error
}

func (c *readCloser) Close() error { return c.close() }
