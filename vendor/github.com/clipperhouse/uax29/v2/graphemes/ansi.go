package graphemes

// ansiEscapeLength returns the byte length of a valid 7-bit ANSI escape
// sequence at the start of data, or 0 if none.
//
// Recognized forms (ECMA-48 / ISO 6429):
//   - CSI: ESC [ then parameter bytes (0x30-0x3F), intermediate (0x20-0x2F), final (0x40-0x7E)
//   - OSC: ESC ] then payload until BEL (0x07), 7-bit ST (ESC \), CAN (0x18), or SUB (0x1A)
//   - DCS, SOS, PM, APC: ESC P/X/^/_ then payload until 7-bit ST (ESC \), CAN, or SUB
//   - Two-byte: ESC + Fe/Fs (0x40-0x7E excluding above), or Fp (0x30-0x3F), or nF (0x20-0x2F then final)
func ansiEscapeLength[T ~string | ~[]byte](data T) int {
	n := len(data)
	if n < 2 || data[0] != esc {
		return 0
	}

	b1 := data[1]
	switch b1 {
	case '[': // CSI
		body := csiBodyLength(data[2:])
		if body == 0 {
			return 0
		}
		return 2 + body
	case ']': // OSC - allows BEL or 7-bit ST terminator
		body := oscLength(data[2:])
		if body < 0 {
			return 0
		}
		return 2 + body
	case 'P', 'X', '^', '_': // DCS, SOS, PM, APC
		body := stSequenceLength(data[2:])
		if body < 0 {
			return 0
		}
		return 2 + body
	}

	if b1 >= 0x40 && b1 <= 0x7E {
		// Fe/Fs two-byte; [ ] P X ^ _ handled above
		return 2
	}
	if b1 >= 0x30 && b1 <= 0x3F {
		// Fp (private) two-byte
		return 2
	}
	if b1 >= 0x20 && b1 <= 0x2F {
		// nF: intermediates then one final (0x30-0x7E)
		i := 2
		for i < n && data[i] >= 0x20 && data[i] <= 0x2F {
			i++
		}
		if i < n && data[i] >= 0x30 && data[i] <= 0x7E {
			return i + 1
		}
		return 0
	}

	return 0
}

// csiBodyLength returns the length of the CSI body (param/intermediate/final bytes).
// data is the slice after "ESC [".
// Per ECMA-48, the CSI body has the form:
//
//	parameters (0x30–0x3F)*, intermediates (0x20–0x2F)*, final (0x40–0x7E)
//
// Once an intermediate byte is seen, subsequent parameter bytes are invalid.
func csiBodyLength[T ~string | ~[]byte](data T) int {
	seenIntermediate := false
	for i := 0; i < len(data); i++ {
		b := data[i]
		if b >= 0x30 && b <= 0x3F {
			if seenIntermediate {
				return 0
			}
			continue
		}
		if b >= 0x20 && b <= 0x2F {
			seenIntermediate = true
			continue
		}
		if b >= 0x40 && b <= 0x7E {
			return i + 1
		}
		return 0
	}
	return 0
}

// oscLength returns the length of the OSC body.
// data is the slice after "ESC ]".
//
// Returns:
//   - n >= 0: consumed body length (includes BEL/ST terminator when present)
//   - -1: not terminated in the provided data
//
// OSC accepts BEL (0x07) or 7-bit ST (ESC \) as terminators by widespread convention.
// Per ECMA-48, CAN (0x18) and SUB (0x1A) cancel the control string; in that
// case they are not part of the OSC sequence length.
func oscLength[T ~string | ~[]byte](data T) int {
	for i := 0; i < len(data); i++ {
		b := data[i]
		if b == bel {
			return i + 1
		}
		if b == can || b == sub {
			return i
		}
		if b == esc && i+1 < len(data) && data[i+1] == '\\' {
			return i + 2
		}
	}
	return -1
}

// stSequenceLength returns the length of a control-string body.
// data is the slice after "ESC x".
//
// Returns:
//   - n >= 0: consumed body length (includes ST terminator when present)
//   - -1: not terminated in the provided data
//
// Used for DCS, SOS, PM, and APC, which per ECMA-48 terminate with ST.
// ST here is the 7-bit form (ESC \).
// CAN (0x18) and SUB (0x1A) cancel the control string; in that case they are
// not part of the sequence length.
func stSequenceLength[T ~string | ~[]byte](data T) int {
	for i := 0; i < len(data); i++ {
		if data[i] == can || data[i] == sub {
			return i
		}
		if data[i] == esc && i+1 < len(data) && data[i+1] == '\\' {
			return i + 2
		}
	}
	return -1
}
