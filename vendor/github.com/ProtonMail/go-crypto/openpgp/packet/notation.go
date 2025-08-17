package packet

// Notation type represents a Notation Data subpacket
// see https://tools.ietf.org/html/rfc4880#section-5.2.3.16
type Notation struct {
	Name            string
	Value           []byte
	IsCritical      bool
	IsHumanReadable bool
}

func (notation *Notation) getData() []byte {
	nameData := []byte(notation.Name)
	nameLen := len(nameData)
	valueLen := len(notation.Value)

	data := make([]byte, 8+nameLen+valueLen)
	if notation.IsHumanReadable {
		data[0] = 0x80
	}

	data[4] = byte(nameLen >> 8)
	data[5] = byte(nameLen)
	data[6] = byte(valueLen >> 8)
	data[7] = byte(valueLen)
	copy(data[8:8+nameLen], nameData)
	copy(data[8+nameLen:], notation.Value)
	return data
}
