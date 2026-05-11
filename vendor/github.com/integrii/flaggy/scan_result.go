package flaggy

// flagScanResult summarizes the outcome of scanning arguments for a parser.
type flagScanResult struct {
	// Positionals lists positional tokens (subcommands or positional args) in
	// the order they were encountered, along with their indexes in the source
	// argument slice.
	Positionals []positionalToken
	// ForwardArgs contains arguments that were intentionally left untouched so
	// that downstream parsers can process them. These tokens maintain their
	// original order.
	ForwardArgs []string
	// HelpRequested reports whether a help flag (-h/--help) was encountered
	// while scanning this parser.
	HelpRequested bool
	// Subcommand holds the first subcommand encountered while scanning. When
	// non-nil, scanning stops and the remaining arguments are handed off to the
	// referenced parser.
	Subcommand *subcommandMatch
}

// positionalToken tracks a positional argument's value and the index it was
// read from in the source slice.
type positionalToken struct {
	Value string
	Index int
}

// subcommandMatch captures the metadata necessary to hand control over to a
// downstream subcommand parser.
type subcommandMatch struct {
	// Command references the subcommand that matched the positional token.
	Command *Subcommand
	// Token points to the positional token that triggered the match.
	Token positionalToken
	// RelativeDepth tracks the positional depth (1-based) where the match was
	// found. This mirrors how subcommand positions are configured.
	RelativeDepth int
}
