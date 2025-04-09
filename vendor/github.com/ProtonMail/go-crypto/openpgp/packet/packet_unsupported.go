package packet

import (
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
)

// UnsupportedPackage represents a OpenPGP packet with a known packet type
// but with unsupported content.
type UnsupportedPacket struct {
	IncompletePacket Packet
	Error            errors.UnsupportedError
}

// Implements the Packet interface
func (up *UnsupportedPacket) parse(read io.Reader) error {
	err := up.IncompletePacket.parse(read)
	if castedErr, ok := err.(errors.UnsupportedError); ok {
		up.Error = castedErr
		return nil
	}
	return err
}
