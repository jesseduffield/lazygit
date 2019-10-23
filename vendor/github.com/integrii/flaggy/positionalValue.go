package flaggy

// PositionalValue represents a value which is determined by its position
// relative to where a subcommand was detected.
type PositionalValue struct {
	Name          string // used in documentation only
	Description   string
	AssignmentVar *string // the var that will get this variable
	Position      int     // the position, not including switches, of this variable
	Required      bool    // this subcommand must always be specified
	Found         bool    // was this positional found during parsing?
	Hidden        bool    // indicates this positional value should be hidden from help
	defaultValue  string  // used for help output
}
