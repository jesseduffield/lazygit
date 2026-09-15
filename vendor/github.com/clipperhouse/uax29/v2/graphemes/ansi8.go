package graphemes

// ansiEscapeLength8Bit returns the byte length of a valid 8-bit C1 ANSI
// sequence at the start of data, or 0 if none.
//
// Recognized forms (ECMA-48 / ISO 6429):
//   - C1 CSI (0x9B) body as parameter/intermediate/final bytes
//   - C1 OSC (0x9D) body terminated by BEL, C1 ST, CAN, or SUB
//   - C1 DCS/SOS/PM/APC (0x90/0x98/0x9E/0x9F) body terminated by C1 ST, CAN, or SUB
//   - Standalone C1 controls (0x80..0x9F not listed above): single byte
func ansiEscapeLength8Bit[T ~string | ~[]byte](data T) int {
	if len(data) == 0 {
		return 0
	}

	switch data[0] {
	case 0x9B: // C1 CSI
		body := csiBodyLength(data[1:])
		if body == 0 {
			return 0
		}
		return 1 + body
	case 0x9D: // C1 OSC
		body := oscLengthC1(data[1:])
		if body < 0 {
			return 0
		}
		return 1 + body
	case 0x90, 0x98, 0x9E, 0x9F: // C1 DCS, SOS, PM, APC
		body := stSequenceLengthC1(data[1:])
		if body < 0 {
			return 0
		}
		return 1 + body
	default:
		if data[0] >= 0x80 && data[0] <= 0x9F {
			return 1
		}
	}

	return 0
}

// oscLengthC1 returns the length of a C1 OSC body.
// data is the slice after the C1 OSC initiator (0x9D).
//
// Returns:
//   - n >= 0: consumed body length (includes BEL/ST terminator when present)
//   - -1: not terminated in the provided data
//
// Terminators: BEL (0x07) or C1 ST (0x9C).
// CAN (0x18) and SUB (0x1A) cancel the control string.
func oscLengthC1[T ~string | ~[]byte](data T) int {
	for i := 0; i < len(data); i++ {
		b := data[i]
		if b == bel || b == st {
			return i + 1
		}
		if b == can || b == sub {
			return i
		}
	}
	return -1
}

// stSequenceLengthC1 parses DCS/SOS/PM/APC bodies that terminate with C1 ST
// (0x9C), or are canceled by CAN/SUB.
func stSequenceLengthC1[T ~string | ~[]byte](data T) int {
	for i := 0; i < len(data); i++ {
		b := data[i]
		if b == can || b == sub {
			return i
		}
		if b == st {
			return i + 1
		}
	}
	return -1
}
