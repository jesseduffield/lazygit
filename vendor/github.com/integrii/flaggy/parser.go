package flaggy

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	initialSubcommandContext   *Subcommand        // points to the initial help context prior to parsing
	ShowCompletion             bool               // indicates that bash and zsh completion output is possible
	SortFlags                  bool               // when true, help output flags are sorted alphabetically
	SortFlagsReverse           bool               // when true with SortFlags, sort order is reversed (Z..A)
}

// supportedCompletionShells lists every shell that can receive generated completion output.
var supportedCompletionShells = []string{"bash", "zsh", "fish", "powershell", "nushell"}

// completionShellList joins the supported completion shell names into a space separated string.
func completionShellList() string {
	return strings.Join(supportedCompletionShells, " ")
}

// isSupportedCompletionShell reports whether the provided shell is eligible for generated completions.
func isSupportedCompletionShell(shell string) bool {
	for _, supported := range supportedCompletionShells {
		if shell == supported {
			return true
		}
	}
	return false
}

// TrailingSubcommand returns the last and most specific subcommand invoked.
func (p *Parser) TrailingSubcommand() *Subcommand {
	return p.subcommandContext
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
	p.ShowCompletion = true
	p.SortFlags = false
	p.SortFlagsReverse = false
	p.SetHelpTemplate(DefaultHelpTemplate)
	initialContext := &Subcommand{}
	p.subcommandContext = initialContext
	p.initialSubcommandContext = initialContext
	return p
}

// isTopLevelHelpContext returns true when help output should be shown for the top
// level parser instead of a specific subcommand.
func (p *Parser) isTopLevelHelpContext() bool {
	if p.subcommandContext == nil {
		return true
	}
	if p.subcommandContext == &p.Subcommand {
		return true
	}
	if p.initialSubcommandContext != nil && p.subcommandContext == p.initialSubcommandContext {
		return true
	}
	return false
}

// SortFlagsByLongName enables alphabetical sorting by long flag name
// (case-insensitive) for help output on this parser.
func (p *Parser) SortFlagsByLongName() {
	p.SortFlags = true
	p.SortFlagsReverse = false
}

// SortFlagsByLongNameReversed enables reverse alphabetical sorting by
// long flag name (case-insensitive) for help output on this parser.
func (p *Parser) SortFlagsByLongNameReversed() {
	p.SortFlags = true
	p.SortFlagsReverse = true
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

	// Handle shell completion before any parsing to avoid unknown-argument exits.
	if p.ShowCompletion {
		if len(args) >= 1 && strings.EqualFold(args[0], "completion") {
			// no shell provided
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "Please specify a shell for completion. Supported shells: %s\n", completionShellList())
				exitOrPanic(2)
			}

			shell := strings.ToLower(args[1])
			if isSupportedCompletionShell(shell) {
				p.Completion(shell)
				exitOrPanic(0)
			}
			fmt.Fprintf(os.Stderr, "Unsupported shell specified for completion: %s\nSupported shells: %s\n", args[1], completionShellList())
			exitOrPanic(2)
		}
	}

	debugPrint("Kicking off parsing with args:", args)
	err := p.parse(p, args)
	if err != nil {
		return err
	}

	// if we are set to exit on unexpected args, look for those here
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

// Completion takes in a shell type and outputs the completion script for
// that shell.
func (p *Parser) Completion(completionType string) {
	switch strings.ToLower(completionType) {
	case "bash":
		fmt.Print(GenerateBashCompletion(p))
	case "zsh":
		fmt.Print(GenerateZshCompletion(p))
	case "fish":
		fmt.Print(GenerateFishCompletion(p))
	case "powershell":
		fmt.Print(GeneratePowerShellCompletion(p))
	case "nushell":
		fmt.Print(GenerateNushellCompletion(p))
	default:
		fmt.Fprintf(os.Stderr, "Unsupported shell specified for completion: %s\nSupported shells: %s\n", completionType, completionShellList())
	}
}

// findArgsNotInParsedValues finds arguments not used in parsed values.  The
// incoming args should be in the order supplied by the user and should not
// include the invoked binary, which is normally the first thing in os.Args.
func findArgsNotInParsedValues(args []string, parsedValues []parsedValue) []string {
	// DebugMode = true
	// defer func() {
	// 	DebugMode = false
	// }()

	var argsNotUsed []string
	var skipNext bool

	for i := 0; i < len(args); i++ {
		a := args[i]

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

		// Determine token type and normalized key/value
		arg := parseFlagToName(a)
		isFlagToken := strings.HasPrefix(a, "-")

		// skip args that start with 'test.' because they are injected with go test
		debugPrint("flagsNotParsed: checking arg for test prefix:", arg)
		if strings.HasPrefix(arg, "test.") {
			debugPrint("skipping test. prefixed arg has test prefix:", arg)
			continue
		}
		debugPrint("flagsNotParsed: flag is not a test. flag:", arg)

		// indicates that we found this arg used in one of the parsed values. Used
		// to indicate which values should be added to argsNotUsed.
		var foundArgUsed bool

		// For flag tokens, only allow non-positional (flag) matches.
		if isFlagToken {
			for _, pv := range parsedValues {
				debugPrint(pv.Key + "==" + arg + " || (" + strconv.FormatBool(pv.IsPositional) + " && " + pv.Value + " == " + arg + ")")
				if !pv.IsPositional && pv.Key == arg {
					debugPrint("Found matching parsed flag for " + pv.Key)
					foundArgUsed = true
					if pv.ConsumesNext {
						skipNext = true
					} else if i+1 < len(args) && pv.Value == args[i+1] {
						skipNext = true
					}
					break
				}
			}
			if foundArgUsed {
				continue
			}
		}

		// For non-flag tokens, prefer positional matches first.
		if !isFlagToken {
			for _, pv := range parsedValues {
				debugPrint(pv.Key + "==" + arg + " || (" + strconv.FormatBool(pv.IsPositional) + " && " + pv.Value + " == " + arg + ")")
				if pv.IsPositional && pv.Value == arg {
					debugPrint("Found matching parsed positional for " + pv.Value)
					foundArgUsed = true
					break
				}
			}
			if foundArgUsed {
				continue
			}

			// Fallback for non-flag tokens: allow matching a non-positional flag by bare name.
			for _, pv := range parsedValues {
				debugPrint(pv.Key + "==" + arg + " || (" + strconv.FormatBool(pv.IsPositional) + " && " + pv.Value + " == " + arg + ")")
				if !pv.IsPositional && pv.Key == arg {
					debugPrint("Found matching parsed flag for " + pv.Key)
					foundArgUsed = true
					if pv.ConsumesNext {
						skipNext = true
					} else if i+1 < len(args) && pv.Value == args[i+1] {
						skipNext = true
					}
					break
				}
			}
			if foundArgUsed {
				continue
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
	return p.ParseArgs(os.Args[1:])
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
