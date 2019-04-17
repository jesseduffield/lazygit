package getter

import (
	"io"
)

// WithProgress allows for a user to track
// the progress of a download.
// For example by displaying a progress bar with
// current download.
// Not all getters have progress support yet.
func WithProgress(pl ProgressTracker) func(*Client) error {
	return func(c *Client) error {
		c.ProgressListener = pl
		return nil
	}
}

// ProgressTracker allows to track the progress of downloads.
type ProgressTracker interface {
	// TrackProgress should be called when
	// a new object is being downloaded.
	// src is the location the file is
	// downloaded from.
	// currentSize is the current size of
	// the file in case it is a partial
	// download.
	// totalSize is the total size in bytes,
	// size can be zero if the file size
	// is not known.
	// stream is the file being downloaded, every
	// written byte will add up to processed size.
	//
	// TrackProgress returns a ReadCloser that wraps the
	// download in progress ( stream ).
	// When the download is finished, body shall be closed.
	TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) (body io.ReadCloser)
}
