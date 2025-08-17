package packet

import (
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
)

type Marker struct{}

const markerString = "PGP"

// parse just checks if the packet contains "PGP".
func (m *Marker) parse(reader io.Reader) error {
	var buffer [3]byte
	if _, err := io.ReadFull(reader, buffer[:]); err != nil {
		return err
	}
	if string(buffer[:]) != markerString {
		return errors.StructuralError("invalid marker packet")
	}
	return nil
}

// SerializeMarker writes a marker packet to writer.
func SerializeMarker(writer io.Writer) error {
	err := serializeHeader(writer, packetTypeMarker, len(markerString))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(markerString))
	return err
}
