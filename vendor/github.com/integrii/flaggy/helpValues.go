package flaggy

import (
	"log"
	"reflect"
	"strings"
	"unicode/utf8"
)

// Help represents the values needed to render a Help page
type Help struct {
	Subcommands    []HelpSubcommand
	Positionals    []HelpPositional
	Flags          []HelpFlag
	UsageString    string
	CommandName    string
	PrependMessage string
	AppendMessage  string
	Message        string
	Description    string
}

// HelpSubcommand is used to template subcommand Help output
type HelpSubcommand struct {
	ShortName   string
	LongName    string
	Description string
	Position    int
	Spacer      string
}

// HelpPositional is used to template positional Help output
type HelpPositional struct {
	Name         string
	Description  string
	Required     bool
	Position     int
	DefaultValue string
	Spacer       string
}

// HelpFlag is used to template string flag Help output
type HelpFlag struct {
	ShortName    string
	LongName     string
	Description  string
	DefaultValue string
	Spacer       string
}

// ExtractValues extracts Help template values from a subcommand and its parent
// parser. The parser is required in order to detect default flag settings
// for help and version output.
func (h *Help) ExtractValues(p *Parser, message string) {

	// accept message string for output
	h.Message = message

	// extract Help values from the current subcommand in context
	// prependMessage string
	h.PrependMessage = p.subcommandContext.AdditionalHelpPrepend
	// appendMessage  string
	h.AppendMessage = p.subcommandContext.AdditionalHelpAppend
	// command name
	h.CommandName = p.subcommandContext.Name
	// description
	h.Description = p.subcommandContext.Description

	maxLength := getLongestNameLength(p.subcommandContext.Subcommands, 0)

	// subcommands    []HelpSubcommand
	for _, cmd := range p.subcommandContext.Subcommands {
		if cmd.Hidden {
			continue
		}
		newHelpSubcommand := HelpSubcommand{
			ShortName:   cmd.ShortName,
			LongName:    cmd.Name,
			Description: cmd.Description,
			Position:    cmd.Position,
			Spacer:      makeSpacer(cmd.Name, maxLength),
		}
		h.Subcommands = append(h.Subcommands, newHelpSubcommand)
	}

	maxLength = getLongestNameLength(p.subcommandContext.PositionalFlags, 0)

	// parse positional flags into help output structs
	for _, pos := range p.subcommandContext.PositionalFlags {
		if pos.Hidden {
			continue
		}
		newHelpPositional := HelpPositional{
			Name:         pos.Name,
			Position:     pos.Position,
			Description:  pos.Description,
			Required:     pos.Required,
			DefaultValue: pos.defaultValue,
			Spacer:       makeSpacer(pos.Name, maxLength),
		}
		h.Positionals = append(h.Positionals, newHelpPositional)
	}

	maxLength = len(versionFlagLongName)
	if len(helpFlagLongName) > maxLength {
		maxLength = len(helpFlagLongName)
	}
	maxLength = getLongestNameLength(p.subcommandContext.Flags, maxLength)
	maxLength = getLongestNameLength(p.Flags, maxLength)

	// if the built-in version flag is enabled, then add it as a help flag
	if p.ShowVersionWithVersionFlag {
		defaultVersionFlag := HelpFlag{
			ShortName:    "",
			LongName:     versionFlagLongName,
			Description:  "Displays the program version string.",
			DefaultValue: "",
			Spacer:       makeSpacer(versionFlagLongName, maxLength),
		}
		h.Flags = append(h.Flags, defaultVersionFlag)
	}

	// if the built-in help flag exists, then add it as a help flag
	if p.ShowHelpWithHFlag {
		defaultHelpFlag := HelpFlag{
			ShortName:    helpFlagShortName,
			LongName:     helpFlagLongName,
			Description:  "Displays help with available flag, subcommand, and positional value parameters.",
			DefaultValue: "",
			Spacer:       makeSpacer(helpFlagLongName, maxLength),
		}
		h.Flags = append(h.Flags, defaultHelpFlag)
	}

	// go through every flag in the subcommand and add it to help output
	h.parseFlagsToHelpFlags(p.subcommandContext.Flags, maxLength)

	// go through every flag in the parent parser and add it to help output
	h.parseFlagsToHelpFlags(p.Flags, maxLength)

	// formulate the usage string
	// first, we capture all the command and positional names by position
	commandsByPosition := make(map[int]string)
	for _, pos := range p.subcommandContext.PositionalFlags {
		if pos.Hidden {
			continue
		}
		if len(commandsByPosition[pos.Position]) > 0 {
			commandsByPosition[pos.Position] = commandsByPosition[pos.Position] + "|" + pos.Name
		} else {
			commandsByPosition[pos.Position] = pos.Name
		}
	}
	for _, cmd := range p.subcommandContext.Subcommands {
		if cmd.Hidden {
			continue
		}
		if len(commandsByPosition[cmd.Position]) > 0 {
			commandsByPosition[cmd.Position] = commandsByPosition[cmd.Position] + "|" + cmd.Name
		} else {
			commandsByPosition[cmd.Position] = cmd.Name
		}
	}

	// find the highest position count in the map
	var highestPosition int
	for i := range commandsByPosition {
		if i > highestPosition {
			highestPosition = i
		}
	}

	// only have a usage string if there are positional items
	var usageString string
	if highestPosition > 0 {
		// find each positional value and make our final string
		usageString = p.subcommandContext.Name
		for i := 1; i <= highestPosition; i++ {
			if len(commandsByPosition[i]) > 0 {
				usageString = usageString + " [" + commandsByPosition[i] + "]"
			} else {
				// dont keep listing after the first position without any properties
				// it will be impossible to reach anything beyond here anyway
				break
			}
		}
	}

	h.UsageString = usageString

}

// parseFlagsToHelpFlags parses the specified slice of flags into
// help flags on the the calling help command
func (h *Help) parseFlagsToHelpFlags(flags []*Flag, maxLength int) {

	for _, f := range flags {
		if f.Hidden {
			continue
		}

		// parse help values out if the flag hasn't been parsed yet
		if !f.parsed {
			f.parsed = true
			// parse the default value as a string and remember it for help output
			f.defaultValue, _ = f.returnAssignmentVarValueAsString()
		}

		// determine the default value based on the assignment variable
		defaultValue := f.defaultValue

		// dont show nils
		if defaultValue == "<nil>" {
			defaultValue = ""
		}

		// for bools, dont show a default of false
		_, isBool := f.AssignmentVar.(*bool)
		if isBool {
			b := f.AssignmentVar.(*bool)
			if *b == false {
				defaultValue = ""
			}
		}

		newHelpFlag := HelpFlag{
			ShortName:    f.ShortName,
			LongName:     f.LongName,
			Description:  f.Description,
			DefaultValue: defaultValue,
			Spacer:       makeSpacer(f.LongName, maxLength),
		}
		h.AddFlagToHelp(newHelpFlag)
	}
}

// AddFlagToHelp adds a flag to help output if it does not exist
func (h *Help) AddFlagToHelp(f HelpFlag) {
	for _, existingFlag := range h.Flags {
		if len(existingFlag.ShortName) > 0 && existingFlag.ShortName == f.ShortName {
			return
		}
		if len(existingFlag.LongName) > 0 && existingFlag.LongName == f.LongName {
			return
		}
	}
	h.Flags = append(h.Flags, f)
}

// getLongestNameLength takes a slice of any supported flag and returns the length of the longest of their names
func getLongestNameLength(slice interface{}, min int) int {
	var maxLength = min

	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		log.Panicf("Paremeter given to getLongestNameLength() is of type %s. Expected slice", s.Kind())
	}

	for i := 0; i < s.Len(); i++ {
		option := s.Index(i).Interface()
		var name string
		switch t := option.(type) {
		case *Subcommand:
			name = t.Name
		case *Flag:
			name = t.LongName
		case *PositionalValue:
			name = t.Name
		default:
			log.Panicf("Unexpected type %T found in slice passed to getLongestNameLength(). Possible types: *Subcommand, *Flag, *PositionalValue", t)
		}
		length := len(name)
		if length > maxLength {
			maxLength = length
		}
	}

	return maxLength
}

// makeSpacer creates a string of whitespaces, with a length of the given
// maxLength minus the length of the given name
func makeSpacer(name string, maxLength int) string {
	length := maxLength - utf8.RuneCountInString(name)
	if length < 0 {
		length = 0
	}
	return strings.Repeat(" ", length)
}
