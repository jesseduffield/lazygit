package packp

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
)

// UploadPackRequest represents a upload-pack request.
// Zero-value is not safe, use NewUploadPackRequest instead.
type UploadPackRequest struct {
	UploadRequest
	UploadHaves
}

// NewUploadPackRequest creates a new UploadPackRequest and returns a pointer.
func NewUploadPackRequest() *UploadPackRequest {
	ur := NewUploadRequest()
	return &UploadPackRequest{
		UploadHaves:   UploadHaves{},
		UploadRequest: *ur,
	}
}

// NewUploadPackRequestFromCapabilities creates a new UploadPackRequest and
// returns a pointer. The request capabilities are filled with the most optimal
// ones, based on the adv value (advertised capabilities), the UploadPackRequest
// it has no wants, haves or shallows and an infinite depth
func NewUploadPackRequestFromCapabilities(adv *capability.List) *UploadPackRequest {
	ur := NewUploadRequestFromCapabilities(adv)
	return &UploadPackRequest{
		UploadHaves:   UploadHaves{},
		UploadRequest: *ur,
	}
}

// IsEmpty returns whether a request is empty - it is empty if Haves are contained
// in the Wants, or if Wants length is zero, and we don't have any shallows
func (r *UploadPackRequest) IsEmpty() bool {
	return isSubset(r.Wants, r.Haves) && len(r.Shallows) == 0
}

func isSubset(needle []plumbing.Hash, haystack []plumbing.Hash) bool {
	for _, h := range needle {
		found := false
		for _, oh := range haystack {
			if h == oh {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// UploadHaves is a message to signal the references that a client has in a
// upload-pack. Do not use this directly. Use UploadPackRequest request instead.
type UploadHaves struct {
	Haves []plumbing.Hash
}

// Encode encodes the UploadHaves into the Writer. If flush is true, a flush
// command will be encoded at the end of the writer content.
func (u *UploadHaves) Encode(w io.Writer, flush bool) error {
	e := pktline.NewEncoder(w)

	plumbing.HashesSort(u.Haves)

	var last plumbing.Hash
	for _, have := range u.Haves {
		if bytes.Equal(last[:], have[:]) {
			continue
		}

		if err := e.Encodef("have %s\n", have); err != nil {
			return fmt.Errorf("sending haves for %q: %s", have, err)
		}

		last = have
	}

	if flush && len(u.Haves) != 0 {
		if err := e.Flush(); err != nil {
			return fmt.Errorf("sending flush-pkt after haves: %s", err)
		}
	}

	return nil
}
