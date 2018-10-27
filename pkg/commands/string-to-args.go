package commands

import "runtime"

// ToArgv converts string s into an argv for exec.
func ToArgv(s string) []string {
	const (
		InArg = iota
		InArgQuote
		OutOfArg
	)
	currentState := OutOfArg
	currentQuoteChar := "\x00" // to distinguish between ' and " quotations
	// this allows to use "foo'bar"
	currentArg := ""
	argv := []string{}

	isQuote := func(c string) bool {
		return c == `"` || c == `'`
	}

	isEscape := func(c string) bool {
		return c == `\`
	}

	isWhitespace := func(c string) bool {
		return c == " " || c == "\t"
	}

	L := len(s)
	for i := 0; i < L; i++ {
		c := s[i : i+1]

		//fmt.Printf("c %s state %v arg %s argv %v i %d\n", c, currentState, currentArg, args, i)
		if isQuote(c) {
			switch currentState {
			case OutOfArg:
				currentArg = ""
				fallthrough
			case InArg:
				currentState = InArgQuote
				currentQuoteChar = c

			case InArgQuote:
				if c == currentQuoteChar {
					currentState = InArg
				} else {
					currentArg += c
				}
			}

		} else if isWhitespace(c) {
			switch currentState {
			case InArg:
				argv = append(argv, currentArg)
				currentState = OutOfArg
			case InArgQuote:
				currentArg += c
			case OutOfArg:
				// nothing
			}

		} else if isEscape(c) {
			switch currentState {
			case OutOfArg:
				currentArg = ""
				currentState = InArg
				fallthrough
			case InArg:
				fallthrough
			case InArgQuote:
				if i == L-1 {
					if runtime.GOOS == "windows" {
						// just add \ to end for windows
						currentArg += c
					} else {
						panic("Escape character at end string")
					}
				} else {
					if runtime.GOOS == "windows" {
						peek := s[i+1 : i+2]
						if peek != `"` {
							currentArg += c
						}
					} else {
						i++
						c = s[i : i+1]
						currentArg += c
					}
				}
			}
		} else {
			switch currentState {
			case InArg, InArgQuote:
				currentArg += c

			case OutOfArg:
				currentArg = ""
				currentArg += c
				currentState = InArg
			}
		}
	}

	if currentState == InArg {
		argv = append(argv, currentArg)
	} else if currentState == InArgQuote {
		panic("Starting quote has no ending quote.")
	}

	return argv
}
