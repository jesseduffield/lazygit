package flaggy

import (
	"errors"
	"fmt"
	"os"

	"text/template"
)

// Parser represents the set of vars and subcommands we are expecting
// from our input args, and the parser than handles them all.
type Parser struct {
	Subcommand
	Version                    string             // the optional version of the parser.
	ShowHelpWithHFlag          bool               // display help when -h or --help passed
	ShowVersionWithVersionFlag bool               // display the version when --version passed
	ShowHelpOnUnexpected       bool               // display help when an unexpected flag is passed
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
	// debugPrint("Kicking off parsing with args:", args)
	return p.parse(p, args, 0)
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

// Disable show version with --version. It is enabled by default.
func (p *Parser) DisableShowVersionWithVersion() {
	p.ShowVersionWithVersionFlag = false
}
