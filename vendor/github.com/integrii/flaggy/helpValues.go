package flaggy

import (
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Help represents the values needed to render a Help page
type Help struct {
	Subcommands    []HelpSubcommand
	Positionals    []HelpPositional
	Flags          []HelpFlag
	GlobalFlags    []HelpFlag
	UsageString    string
	CommandName    string
	PrependMessage string
	AppendMessage  string
	ShowCompletion bool
	Message        string
	Description    string
	Lines          []string
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
	ShortDisplay string
	LongDisplay  string
}

// ExtractValues extracts Help template values from a subcommand and its parent
// parser. The parser is required in order to detect default flag settings
// for help and version output.
func (h *Help) ExtractValues(p *Parser, message string) {
	// accept message string for output
	h.Message = message

	ctx := p.subcommandContext
	if ctx == nil || ctx == p.initialSubcommandContext {
		ctx = &p.Subcommand
	}
	isRootContext := ctx == &p.Subcommand

	// extract Help values from the current subcommand in context
	// prependMessage string
	h.PrependMessage = ctx.AdditionalHelpPrepend
	// appendMessage  string
	h.AppendMessage = ctx.AdditionalHelpAppend
	// command name
	h.CommandName = ctx.Name
	// description
	h.Description = ctx.Description
	// shell completion
	showCompletion := p.ShowCompletion && p.isTopLevelHelpContext()
	h.ShowCompletion = showCompletion

	// determine the max length of subcommand names for spacer calculation.
	maxLength := getLongestNameLength(ctx.Subcommands, 0)
	// include the synthetic completion subcommand in spacer calculation
	if showCompletion {
		if l := len("completion"); l > maxLength {
			maxLength = l
		}
	}

	// subcommands    []HelpSubcommand
	for _, cmd := range ctx.Subcommands {
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

	// Append a synthetic completion subcommand at the end when enabled.
	// This shows users the correct invocation: "./appName completion [bash|zsh]".
	if showCompletion {
		completionHelp := HelpSubcommand{
			ShortName:   "",
			LongName:    "completion",
			Description: "Generate shell completion script for bash or zsh.",
			Position:    0,
			Spacer:      makeSpacer("completion", maxLength),
		}
		h.Subcommands = append(h.Subcommands, completionHelp)
	}

	maxLength = getLongestNameLength(ctx.PositionalFlags, 0)

	// parse positional flags into help output structs
	for _, pos := range ctx.PositionalFlags {
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

	// if the built-in version flag is enabled, then add it to the appropriate help collection
	if p.ShowVersionWithVersionFlag {
		defaultVersionFlag := HelpFlag{
			ShortName:    "",
			LongName:     versionFlagLongName,
			Description:  "Displays the program version string.",
			DefaultValue: "",
		}
		if isRootContext {
			h.addFlagToSlice(&h.Flags, defaultVersionFlag)
		} else {
			h.addFlagToSlice(&h.GlobalFlags, defaultVersionFlag)
		}
	}

	// if the built-in help flag exists, then add it as a help flag
	if p.ShowHelpWithHFlag {
		defaultHelpFlag := HelpFlag{
			ShortName:    helpFlagShortName,
			LongName:     helpFlagLongName,
			Description:  "Displays help with available flag, subcommand, and positional value parameters.",
			DefaultValue: "",
		}
		if isRootContext {
			h.addFlagToSlice(&h.Flags, defaultHelpFlag)
		} else {
			h.addFlagToSlice(&h.GlobalFlags, defaultHelpFlag)
		}
	}

	// go through every flag in the subcommand and add it to help output
	h.parseFlagsToHelpFlags(ctx.Flags, &h.Flags)

	// go through every flag in the parent parser and add it to help output
	if isRootContext {
		h.parseFlagsToHelpFlags(p.Flags, &h.Flags)
	} else {
		h.parseFlagsToHelpFlags(p.Flags, &h.GlobalFlags)
	}

	// Optionally sort flags alphabetically by long name (fallback to short name)
	if p.SortFlags {
		sort.SliceStable(h.Flags, func(i, j int) bool {
			a := h.Flags[i]
			b := h.Flags[j]
			aName := strings.ToLower(strings.TrimSpace(a.LongName))
			bName := strings.ToLower(strings.TrimSpace(b.LongName))
			if aName == "" {
				aName = strings.ToLower(strings.TrimSpace(a.ShortName))
			}
			if bName == "" {
				bName = strings.ToLower(strings.TrimSpace(b.ShortName))
			}
			if p.SortFlagsReverse {
				return aName > bName
			}
			return aName < bName
		})
		sort.SliceStable(h.GlobalFlags, func(i, j int) bool {
			a := h.GlobalFlags[i]
			b := h.GlobalFlags[j]
			aName := strings.ToLower(strings.TrimSpace(a.LongName))
			bName := strings.ToLower(strings.TrimSpace(b.LongName))
			if aName == "" {
				aName = strings.ToLower(strings.TrimSpace(a.ShortName))
			}
			if bName == "" {
				bName = strings.ToLower(strings.TrimSpace(b.ShortName))
			}
			if p.SortFlagsReverse {
				return aName > bName
			}
			return aName < bName
		})
	}

	// formulate the usage string
	// first, we capture all the command and positional names by position
	commandsByPosition := make(map[int]string)
	for _, pos := range ctx.PositionalFlags {
		if pos.Hidden {
			continue
		}
		if len(commandsByPosition[pos.Position]) > 0 {
			commandsByPosition[pos.Position] = commandsByPosition[pos.Position] + "|" + pos.Name
		} else {
			commandsByPosition[pos.Position] = pos.Name
		}
	}
	for _, cmd := range ctx.Subcommands {
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
		usageString = ctx.Name
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

	alignHelpFlags(h.Flags)
	alignHelpFlags(h.GlobalFlags)
	h.composeLines()
}

// parseFlagsToHelpFlags parses the specified slice of flags into
// help flags on the the calling help command
func (h *Help) parseFlagsToHelpFlags(flags []*Flag, dest *[]HelpFlag) {
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
			if !*b {
				defaultValue = ""
			}
		}

		newHelpFlag := HelpFlag{
			ShortName:    f.ShortName,
			LongName:     f.LongName,
			Description:  f.Description,
			DefaultValue: defaultValue,
		}
		h.addFlagToSlice(dest, newHelpFlag)
	}
}

// addFlagToSlice adds a flag to the provided slice if it does not exist already.
func (h *Help) addFlagToSlice(dest *[]HelpFlag, f HelpFlag) {
	for _, existingFlag := range *dest {
		if len(existingFlag.ShortName) > 0 && existingFlag.ShortName == f.ShortName {
			return
		}
		if len(existingFlag.LongName) > 0 && existingFlag.LongName == f.LongName {
			return
		}
	}
	*dest = append(*dest, f)
}

// getLongestNameLength takes a slice of any supported flag and returns the length of the longest of their names
func getLongestNameLength(slice interface{}, min int) int {
	var maxLength = min

	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		log.Panicf("Parameter given to getLongestNameLength() is of type %s. Expected slice", s.Kind())
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

func alignHelpFlags(flags []HelpFlag) {
	if len(flags) == 0 {
		return
	}

	shortWidth := 0
	longWidth := 0

	for _, flag := range flags {
		shortCol := flagShortColumn(flag.ShortName)
		longCol := flagLongColumn(flag.LongName)
		if l := utf8.RuneCountInString(shortCol); l > shortWidth {
			shortWidth = l
		}
		if l := utf8.RuneCountInString(longCol); l > longWidth {
			longWidth = l
		}
	}

	const shortGap = "  "
	const descGap = "   "

	for i := range flags {
		shortCol := flagShortColumn(flags[i].ShortName)
		longCol := flagLongColumn(flags[i].LongName)

		if shortWidth > 0 {
			flags[i].ShortDisplay = padRight(shortCol, shortWidth) + shortGap
		} else {
			flags[i].ShortDisplay = shortGap
		}

		if longWidth > 0 {
			flags[i].LongDisplay = padRight(longCol, longWidth) + descGap
		} else {
			flags[i].LongDisplay = descGap
		}
	}
}

func flagShortColumn(shortName string) string {
	if shortName == "" {
		return ""
	}
	return "-" + shortName
}

func flagLongColumn(longName string) string {
	if longName == "" {
		return ""
	}
	return "--" + longName
}

func padRight(input string, width int) string {
	delta := width - utf8.RuneCountInString(input)
	if delta <= 0 {
		return input
	}
	return input + strings.Repeat(" ", delta)
}

func (h *Help) composeLines() {
	lines := make([]string, 0, 16)

	appendBlank := func() {
		if len(lines) > 0 && lines[len(lines)-1] != "" {
			lines = append(lines, "")
		}
	}

	if h.CommandName != "" || h.Description != "" {
		header := h.CommandName
		if h.Description != "" {
			if header != "" {
				header += " - "
			}
			header += h.Description
		}
		lines = append(lines, header)
	}

	if h.PrependMessage != "" {
		lines = append(lines, splitLines(h.PrependMessage)...)
	}

	appendSection := func(section []string) {
		if len(section) == 0 {
			return
		}
		appendBlank()
		lines = append(lines, section...)
	}

	if h.UsageString != "" {
		section := []string{
			"  Usage:",
			"    " + h.UsageString,
		}
		appendSection(section)
	}

	if len(h.Positionals) > 0 {
		section := []string{"  Positional Variables:"}
		for _, pos := range h.Positionals {
			line := "    " + pos.Name + "  " + pos.Spacer
			if pos.Description != "" {
				line += " " + pos.Description
			}
			if pos.DefaultValue != "" {
				line += " (default: " + pos.DefaultValue + ")"
			} else if pos.Required {
				line += " (Required)"
			}
			section = append(section, line)
		}
		appendSection(section)
	}

	if len(h.Subcommands) > 0 {
		section := []string{"  Subcommands:"}
		for _, sub := range h.Subcommands {
			line := "    " + sub.LongName
			if sub.ShortName != "" {
				line += " (" + sub.ShortName + ")"
			}
			if sub.Position > 1 {
				line += "  (position " + strconv.Itoa(sub.Position) + ")"
			}
			if sub.Description != "" {
				line += "   " + sub.Spacer + sub.Description
			}
			section = append(section, line)
		}
		appendSection(section)
	}

	if len(h.Flags) > 0 {
		section := []string{"  Flags:"}
		for _, flag := range h.Flags {
			line := "    " + flag.ShortDisplay + flag.LongDisplay
			descAdded := false
			if flag.Description != "" {
				line += flag.Description
				descAdded = true
			}
			if flag.DefaultValue != "" {
				if descAdded {
					line += " (default: " + flag.DefaultValue + ")"
				} else {
					line += "(default: " + flag.DefaultValue + ")"
				}
			}
			section = append(section, line)
		}
		appendSection(section)
	}

	if len(h.GlobalFlags) > 0 {
		section := []string{"  Global Flags:"}
		for _, flag := range h.GlobalFlags {
			line := "    " + flag.ShortDisplay + flag.LongDisplay
			descAdded := false
			if flag.Description != "" {
				line += flag.Description
				descAdded = true
			}
			if flag.DefaultValue != "" {
				if descAdded {
					line += " (default: " + flag.DefaultValue + ")"
				} else {
					line += "(default: " + flag.DefaultValue + ")"
				}
			}
			section = append(section, line)
		}
		appendSection(section)
	}

	appendText := func(text string) {
		if text == "" {
			return
		}
		appendBlank()
		lines = append(lines, splitLines(text)...)
	}

	appendText(h.AppendMessage)
	appendText(h.Message)

	if len(lines) == 0 {
		lines = append(lines, "")
	} else {
		if lines[0] != "" {
			lines = append([]string{""}, lines...)
		}
		if lines[len(lines)-1] != "" {
			lines = append(lines, "")
		}
	}

	h.Lines = lines
}

func splitLines(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, "\n")
}
