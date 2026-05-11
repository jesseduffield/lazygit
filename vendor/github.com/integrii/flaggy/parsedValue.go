package flaggy

// parsedValue represents a flag or subcommand that was parsed.  Primarily used
// to account for all parsed values in order to determine if unknown values were
// passed to the root parser after all subcommands have been parsed.
type parsedValue struct {
	Key          string
	Value        string
	IsPositional bool // indicates that this value was positional and not a key/value
	ConsumesNext bool // indicates that parsing this value consumed the following CLI token
}

// newParsedValue creates and returns a new parsedValue struct with the
// supplied values set
func newParsedValue(key string, value string, isPositional bool, consumesNext bool) parsedValue {
	if len(key) == 0 && len(value) == 0 {
		panic("can't add parsed value with no key or value")
	}
	return parsedValue{
		Key:          key,
		Value:        value,
		IsPositional: isPositional,
		ConsumesNext: consumesNext,
	}
}
