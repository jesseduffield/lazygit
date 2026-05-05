package displaywidth

// Options allows you to specify the treatment of ambiguous East Asian
// characters and ANSI escape sequences.
type Options struct {
	// EastAsianWidth specifies whether to treat ambiguous East Asian characters
	// as width 1 or 2. When false (default), ambiguous East Asian characters
	// are treated as width 1. When true, they are width 2.
	EastAsianWidth bool

	// ControlSequences specifies whether to ignore 7-bit ECMA-48 escape sequences
	// when calculating the display width. When false (default), ANSI escape
	// sequences are treated as just a series of characters. When true, they are
	// treated as a single zero-width unit.
	ControlSequences bool
	// ControlSequences8Bit specifies whether to ignore 8-bit ECMA-48 escape sequences
	// when calculating the display width. When false (default), these are treated
	// as just a series of characters. When true, they are treated as a single
	// zero-width unit.
	ControlSequences8Bit bool
}

// DefaultOptions is the default options for the display width
// calculation, which is EastAsianWidth false, ControlSequences false, and
// ControlSequences8Bit false.
var DefaultOptions = Options{
	EastAsianWidth:       false,
	ControlSequences:     false,
	ControlSequences8Bit: false,
}
