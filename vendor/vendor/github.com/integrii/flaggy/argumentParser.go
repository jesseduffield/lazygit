package flaggy

// setValueForParsers sets the value for a specified key in the
// specified parsers (which normally include a Parser and Subcommand).
// The return values represent the key being set, and any errors
// returned when setting the key, such as failures to convert the string
// into the appropriate flag value.  We stop assigning values as soon
// as we find a any parser that accepts it.
func setValueForParsers(key string, value string, parsers ...ArgumentParser) (bool, error) {

	for _, p := range parsers {
		valueWasSet, err := p.SetValueForKey(key, value)
		if err != nil {
			return valueWasSet, err
		}
		if valueWasSet {
			return true, nil
		}
	}

	return false, nil
}

// ArgumentParser represents a parser or subcommand
type ArgumentParser interface {
	SetValueForKey(key string, value string) (bool, error)
}
