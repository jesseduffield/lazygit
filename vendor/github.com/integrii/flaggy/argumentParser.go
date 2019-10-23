package flaggy

// setValueForParsers sets the value for a specified key in the
// specified parsers (which normally include a Parser and Subcommand).
// The return values represent the key being set, and any errors
// returned when setting the key, such as failures to convert the string
// into the appropriate flag value.  We stop assigning values as soon
// as we find a parser that accepts it.
func setValueForParsers(key string, value string, parsers ...ArgumentParser) (bool, error) {

	var valueWasSet bool

	for _, p := range parsers {
		valueWasSet, err := p.SetValueForKey(key, value)
		if err != nil {
			return valueWasSet, err
		}
		if valueWasSet {
			break
		}
	}

	return valueWasSet, nil
}

// ArgumentParser represents a parser or subcommand
type ArgumentParser interface {
	SetValueForKey(key string, value string) (bool, error)
}
