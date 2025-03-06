package git

import (
	"io"

	"github.com/go-git/go-git/v5/plumbing"
)

// signableObject is an object which can be signed.
type signableObject interface {
	EncodeWithoutSignature(o plumbing.EncodedObject) error
}

// Signer is an interface for signing git objects.
// message is a reader containing the encoded object to be signed.
// Implementors should return the encoded signature and an error if any.
// See https://git-scm.com/docs/gitformat-signature for more information.
type Signer interface {
	Sign(message io.Reader) ([]byte, error)
}

func signObject(signer Signer, obj signableObject) ([]byte, error) {
	encoded := &plumbing.MemoryObject{}
	if err := obj.EncodeWithoutSignature(encoded); err != nil {
		return nil, err
	}
	r, err := encoded.Reader()
	if err != nil {
		return nil, err
	}

	return signer.Sign(r)
}
