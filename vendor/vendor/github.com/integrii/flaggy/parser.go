package flaggy

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"text/template"
)

// Parser represents the set of flags and subcommands we are expecting
// from our input arguments.  Parser is the top level struct responsible for
// parsing an entire set of subcommands and flags.
type Parser struct {
	Subcommand
	Version                    string             // the optional version of the parser.
	ShowHelpWithHFlag          bool               // display help when -h or --help passed
	ShowVersionWithVersionFlag bool               // display the version when --version passed
	ShowHelpOnUnexpected       bool               // display help when an unexpected flag or subcommand is passed
	TrailingArguments          []string           // everything after a -- is placed here
	HelpTemplate               *template.Template // template for Help output
	trailingArgumentsExtracted bool               // indicates that trailing args have been parsed and should not be appended again
	parsed                     bool               // indicates this parser has parsed
	subcommandContext          *Subcommand        // points to the most specific subcommand being used
}

// NewParser creates a new ArgumentParser ready to parse inputs
func NewParser(name string) *Parser {
	// this can not be done inline because of struct embedding
	p := &Parser{}
	p.Name = name
	p.Version = defaultVersion
	p.ShowHelpOnUnexpected = true
	p.ShowHelpWithHFlag = true
	p.ShowVersionWithVersionFlag = true
	p.SetHelpTemplate(DefaultHelpTemplate)
	p.subcommandContext = &Subcommand{}
	return p
}

// ParseArgs parses as if the passed args were the os.Args, but without the
// binary at the 0 position in the array.  An error is returned if there
// is a low level issue converting flags to their proper type.  No error
// is returned for invalid arguments or missing require subcommands.
func (p *Parser) ParseArgs(args []string) error {
	if p.parsed {
		return errors.New("Parser.Parse() called twice on parser with name: " + " " + p.Name + " " + p.ShortName)
	}
	p.parsed = true

	debugPrint("Kicking off parsing with args:", args)
	err := p.parse(p, args, 0)
	if err != nil {
		return err
	}

	// if we are set to crash on unexpected args, look for those here TODO
	if p.ShowHelpOnUnexpected {
		parsedValues := p.findAllParsedValues()
		debugPrint("parsedValues:", parsedValues)
		argsNotParsed := findArgsNotInParsedValues(args, parsedValues)
		if len(argsNotParsed) > 0 {
			// flatten out unused args for our error message
			var argsNotParsedFlat string
			for _, a := range argsNotParsed {
				argsNotParsedFlat = argsNotParsedFlat + " " + a
			}
			p.ShowHelpAndExit("Unknown arguments supplied: " + argsNotParsedFlat)
		}
	}

	return nil
}

// findArgsNotInParsedValues finds arguments not used in parsed values.  The
// incoming args should be in the order supplied by the user and should not
// include the invoked binary, which is normally the first thing in os.Args.
func findArgsNotInParsedValues(args []string, parsedValues []parsedValue) []string {
	var argsNotUsed []string
	var skipNext bool
	for _, a := range args {

		// if the final argument (--) is seen, then we stop checking because all
		// further values are trailing arguments.
		if determineArgType(a) == argIsFinal {
			return argsNotUsed
		}

		// allow for skipping the next arg when needed
		if skipNext {
			skipNext = false
			continue
		}

		// strip flag slashes from incoming arguments so they match up with the
		// keys from parsedValues.
		arg := parseFlagToName(a)

		// indicates that we found this arg used in one of the parsed values. Used
		// to indicate which values should be added to argsNotUsed.
		var foundArgUsed bool

		// search all args for a corresponding parsed value
		for _, pv := range parsedValues {
			// this argumenet was a key
			// debugPrint(pv.Key, "==", arg)
			debugPrint(pv.Key + "==" + arg + " || (" + strconv.FormatBool(pv.IsPositional) + " && " + pv.Value + " == " + arg + ")")
			if pv.Key == arg || (pv.IsPositional && pv.Value == arg) {
				debugPrint("Found matching parsed arg for " + pv.Key)
				foundArgUsed = true // the arg was used in this parsedValues set
				// if the value is not a positional value and the parsed value had a
				// value that was not blank, we skip the next value in the argument list
				if !pv.IsPositional && len(pv.Value) > 0 {
					skipNext = true
					break
				}
			}
			// this prevents excessive parsed values from being checked after we find
			// the arg used for the first time
			if foundArgUsed {
				break
			}
		}

		// if the arg was not used in any parsed values, then we add it to the slice
		// of arguments not used
		if !foundArgUsed {
			argsNotUsed = append(argsNotUsed, arg)
		}
	}

	return argsNotUsed
}

// ShowVersionAndExit shows the version of this parser
func (p *Parser) ShowVersionAndExit() {
	fmt.Println("Version:", p.Version)
	exitOrPanic(0)
}

// SetHelpTemplate sets the go template this parser will use when rendering
// Help.
func (p *Parser) SetHelpTemplate(tmpl string) error {
	var err error
	p.HelpTemplate = template.New(helpFlagLongName)
	p.HelpTemplate, err = p.HelpTemplate.Parse(tmpl)
	if err != nil {
		return err
	}
	return nil
}

// Parse calculates all flags and subcommands
func (p *Parser) Parse() error {

	err := p.ParseArgs(os.Args[1:])
	if err != nil {
		return err
	}
	return nil

}

// ShowHelp shows Help without an error message
func (p *Parser) ShowHelp() {
	debugPrint("showing help for", p.subcommandContext.Name)
	p.ShowHelpWithMessage("")
}

// ShowHelpAndExit shows parser help and exits with status code 2
func (p *Parser) ShowHelpAndExit(message string) {
	p.ShowHelpWithMessage(message)
	exitOrPanic(2)
}

// ShowHelpWithMessage shows the Help for this parser with an optional string error
// message as a header.  The supplied subcommand will be the context of Help
// displayed to the user.
func (p *Parser) ShowHelpWithMessage(message string) {

	// create a new Help values template and extract values into it
	help := Help{}
	help.ExtractValues(p, message)
	err := p.HelpTemplate.Execute(os.Stderr, help)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error rendering Help template:", err)
	}
}

// DisableShowVersionWithVersion disables the showing of version information
// with --version. It is enabled by default.
func (p *Parser) DisableShowVersionWithVersion() {
	p.ShowVersionWithVersionFlag = false
}
