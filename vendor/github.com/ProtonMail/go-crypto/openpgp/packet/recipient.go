package packet

// Recipient type represents a Intended Recipient Fingerprint subpacket
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh#name-intended-recipient-fingerpr
type Recipient struct {
	KeyVersion  int
	Fingerprint []byte
}

func (r *Recipient) Serialize() []byte {
	packet := make([]byte, len(r.Fingerprint)+1)
	packet[0] = byte(r.KeyVersion)
	copy(packet[1:], r.Fingerprint)
	return packet
}
