package flaggy

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Subcommand represents a subcommand which contains a set of child
// subcommands along with a set of flags relevant to it.  Parsing
// runs until a subcommand is detected by matching its name and
// position.  Once a matching subcommand is found, the next set
// of parsing occurs within that matched subcommand.
type Subcommand struct {
	Name                  string
	ShortName             string
	Description           string
	Position              int // the position of this subcommand, not including flags
	Subcommands           []*Subcommand
	Flags                 []*Flag
	PositionalFlags       []*PositionalValue
	ParsedValues          []parsedValue // a list of values and positionals parsed
	AdditionalHelpPrepend string        // additional prepended message when Help is displayed
	AdditionalHelpAppend  string        // additional appended message when Help is displayed
	Used                  bool          // indicates this subcommand was found and parsed
	Hidden                bool          // indicates this subcommand should be hidden from help
}

// NewSubcommand creates a new subcommand that can have flags or PositionalFlags
// added to it.  The position starts with 1, not 0
func NewSubcommand(name string) *Subcommand {
	if len(name) == 0 {
		fmt.Fprintln(os.Stderr, "Error creating subcommand (NewSubcommand()).  No subcommand name was specified.")
		exitOrPanic(2)
	}
	newSC := &Subcommand{
		Name: name,
	}
	return newSC
}

// parseAllFlagsFromArgs parses the non-positional flags such as -f or -v=value
// out of the supplied args and returns the resulting positional items in order,
// all the flag names found (without values), a bool to indicate if help was
// requested, and any errors found during parsing
func (sc *Subcommand) parseAllFlagsFromArgs(p *Parser, args []string) ([]string, bool, error) {

	var positionalOnlyArguments []string
	var helpRequested bool // indicates the user has supplied -h and we
	// should render help if we are the last subcommand

	// indicates we should skip the next argument, like when parsing a flag
	// that separates key and value by space
	var skipNext bool

	// endArgfound indicates that a -- was found and everything
	// remaining should be added to the trailing arguments slices
	var endArgFound bool

	// find all the normal flags (not positional) and parse them out
	for i, a := range args {

		debugPrint("parsing arg:", a)

		// evaluate if there is a following arg to avoid panics
		var nextArgExists bool
		var nextArg string
		if len(args)-1 >= i+1 {
			nextArgExists = true
			nextArg = args[i+1]
		}

		// if end arg -- has been found, just add everything to TrailingArguments
		if endArgFound {
			if !p.trailingArgumentsExtracted {
				p.TrailingArguments = append(p.TrailingArguments, a)
			}
			continue
		}

		// skip this run if specified
		if skipNext {
			skipNext = false
			debugPrint("skipping flag because it is an arg:", a)
			continue
		}

		// parse the flag into its name for consideration without dashes
		flagName := parseFlagToName(a)

		// if the flag being passed is version or v and the option to display
		// version with version flags, then display version
		if p.ShowVersionWithVersionFlag {
			if flagName == versionFlagLongName {
				p.ShowVersionAndExit()
			}
		}

		// if the show Help on h flag option is set, then show Help when h or Help
		// is passed as an option
		if p.ShowHelpWithHFlag {
			if flagName == helpFlagShortName || flagName == helpFlagLongName {
				// Ensure this is the last subcommand passed so we give the correct
				// help output
				helpRequested = true
				continue
			}
		}

		// determine what kind of flag this is
		argType := determineArgType(a)

		// strip flags from arg
		// debugPrint("Parsing flag named", a, "of type", argType)

		// depending on the flag type, parse the key and value out, then apply it
		switch argType {
		case argIsFinal:
			// debugPrint("Arg", i, "is final:", a)
			endArgFound = true
		case argIsPositional:
			// debugPrint("Arg is positional or subcommand:", a)
			// this positional argument into a slice of their own, so that
			// we can determine if its a subcommand or positional value later
			positionalOnlyArguments = append(positionalOnlyArguments, a)
			// track this as a parsed value with the subcommand
			sc.addParsedPositionalValue(a)
		case argIsFlagWithSpace: // a flag with a space. ex) -k v or --key value
			a = parseFlagToName(a)

			// debugPrint("Arg", i, "is flag with space:", a)
			// parse next arg as value to this flag and apply to subcommand flags
			// if the flag is a bool flag, then we check for a following positional
			// and skip it if necessary
			if flagIsBool(sc, p, a) {
				debugPrint(sc.Name, "bool flag", a, "next var is:", nextArg)
				// set the value in this subcommand and its root parser
				valueSet, err := setValueForParsers(a, "true", p, sc)

				// if an error occurs, just return it and quit parsing
				if err != nil {
					return []string{}, false, err
				}

				// log all values parsed by this subcommand.  We leave the value blank
				// because the bool value had no explicit true or false supplied
				if valueSet {
					sc.addParsedFlag(a, "")
				}

				// we've found and set a standalone bool flag, so we move on to the next
				// argument in the list of arguments
				continue
			}

			skipNext = true
			// debugPrint(sc.Name, "NOT bool flag", a)

			// if the next arg was not found, then show a Help message
			if !nextArgExists {
				p.ShowHelpWithMessage("Expected a following arg for flag " + a + ", but it did not exist.")
				exitOrPanic(2)
			}
			valueSet, err := setValueForParsers(a, nextArg, p, sc)
			if err != nil {
				return []string{}, false, err
			}

			// log all parsed values in the subcommand
			if valueSet {
				sc.addParsedFlag(a, nextArg)
			}
		case argIsFlagWithValue: // a flag with an equals sign. ex) -k=v or --key=value
			// debugPrint("Arg", i, "is flag with value:", a)
			a = parseFlagToName(a)

			// parse flag into key and value and apply to subcommand flags
			key, val := parseArgWithValue(a)

			// set the value in this subcommand and its root parser
			valueSet, err := setValueForParsers(key, val, p, sc)
			if err != nil {
				return []string{}, false, err
			}

			// log all values parsed by the subcommand
			if valueSet {
				sc.addParsedFlag(a, val)
			}
		}
	}

	return positionalOnlyArguments, helpRequested, nil
}

// findAllParsedValues finds all values parsed by all subcommands and this
// subcommand and its child subcommands
func (sc *Subcommand) findAllParsedValues() []parsedValue {
	parsedValues := sc.ParsedValues
	for _, sc := range sc.Subcommands {
		// skip unused subcommands
		if !sc.Used {
			continue
		}
		parsedValues = append(parsedValues, sc.findAllParsedValues()...)
	}
	return parsedValues
}

// parse causes the argument parser to parse based on the supplied []string.
// depth specifies the non-flag subcommand positional depth.  A slice of flags
// and subcommands parsed is returned so that the parser can ultimately decide
// if there were any unexpected values supplied by the user
func (sc *Subcommand) parse(p *Parser, args []string, depth int) error {

	debugPrint("- Parsing subcommand", sc.Name, "with depth of", depth, "and args", args)

	// if a command is parsed, its used
	sc.Used = true
	debugPrint("used subcommand", sc.Name, sc.ShortName)
	if len(sc.Name) > 0 {
		sc.addParsedPositionalValue(sc.Name)
	}
	if len(sc.ShortName) > 0 {
		sc.addParsedPositionalValue(sc.ShortName)
	}

	// as subcommands are used, they become the context of the parser.  This helps
	// us understand how to display help based on which subcommand is being used
	p.subcommandContext = sc

	// ensure that help and version flags are not used if the parser has the
	// built-in help and version flags enabled
	if p.ShowHelpWithHFlag {
		sc.ensureNoConflictWithBuiltinHelp()
	}
	if p.ShowVersionWithVersionFlag {
		sc.ensureNoConflictWithBuiltinVersion()
	}

	// Parse the normal flags out of the argument list and return the positionals
	// (subcommands and positional values), along with the flags used.
	// Then the flag values are applied to the parent parser and the current
	// subcommand being parsed.
	positionalOnlyArguments, helpRequested, err := sc.parseAllFlagsFromArgs(p, args)
	if err != nil {
		return err
	}

	// indicate that trailing arguments have been extracted, so that they aren't
	// appended a second time
	p.trailingArgumentsExtracted = true

	// loop over positional values and look for their matching positional
	// parameter, or their positional command.  If neither are found, then
	// we throw an error
	var parsedArgCount int
	for pos, v := range positionalOnlyArguments {

		// the first relative positional argument will be human natural at position 1
		// but offset for the depth of relative commands being parsed for currently.
		relativeDepth := pos - depth + 1
		// debugPrint("Parsing positional only position", relativeDepth, "with value", v)

		if relativeDepth < 1 {
			// debugPrint(sc.Name, "skipped value:", v)
			continue
		}
		parsedArgCount++

		// determine subcommands and parse them by positional value and name
		for _, cmd := range sc.Subcommands {
			// debugPrint("Subcommand being compared", relativeDepth, "==", cmd.Position, "and", v, "==", cmd.Name, "==", cmd.ShortName)
			if relativeDepth == cmd.Position && (v == cmd.Name || v == cmd.ShortName) {
				debugPrint("Decending into positional subcommand", cmd.Name, "at relativeDepth", relativeDepth, "and absolute depth", depth+1)
				return cmd.parse(p, args, depth+parsedArgCount) // continue recursive positional parsing
			}
		}

		// determine positional args and parse them by positional value and name
		var foundPositional bool
		for _, val := range sc.PositionalFlags {
			if relativeDepth == val.Position {
				debugPrint("Found a positional value at relativePos:", relativeDepth, "value:", v)

				// set original value for help output
				val.defaultValue = *val.AssignmentVar

				// defrerence the struct pointer, then set the pointer property within it
				*val.AssignmentVar = v
				// debugPrint("set positional to value", *val.AssignmentVar)
				foundPositional = true
				val.Found = true
				break
			}
		}

		// if there aren't any positional flags but there are subcommands that
		// were not used, display a useful message with subcommand options.
		if !foundPositional && p.ShowHelpOnUnexpected {
			debugPrint("No positional at position", relativeDepth)
			var foundSubcommandAtDepth bool
			for _, cmd := range sc.Subcommands {
				if cmd.Position == relativeDepth {
					foundSubcommandAtDepth = true
				}
			}

			// if there is a subcommand here but it was not specified, display them all
			// as a suggestion to the user before exiting.
			if foundSubcommandAtDepth {
				// determine which name to use in upcoming help output
				fmt.Fprintln(os.Stderr, sc.Name+":", "No subcommand or positional value found at position", strconv.Itoa(relativeDepth)+".")
				var output string
				for _, cmd := range sc.Subcommands {
					if cmd.Hidden {
						continue
					}
					output = output + " " + cmd.Name
				}
				// if there are available subcommands, let the user know
				if len(output) > 0 {
					output = strings.TrimLeft(output, " ")
					fmt.Println("Available subcommands:", output)
				}
				exitOrPanic(2)
			}

			// if there were not any flags or subcommands at this position at all, then
			// throw an error (display Help if necessary)
			p.ShowHelpWithMessage("Unexpected argument: " + v)
			exitOrPanic(2)
		}
	}

	// if help was requested and we should show help when h is passed,
	if helpRequested && p.ShowHelpWithHFlag {
		p.ShowHelp()
		exitOrPanic(0)
	}

	// find any positionals that were not used on subcommands that were
	// found and throw help (unknown argument) in the global parse or subcommand
	for _, pv := range p.PositionalFlags {
		if pv.Required && !pv.Found {
			p.ShowHelpWithMessage("Required global positional variable " + pv.Name + " not found at position " + strconv.Itoa(pv.Position))
			exitOrPanic(2)
		}
	}
	for _, pv := range sc.PositionalFlags {
		if pv.Required && !pv.Found {
			p.ShowHelpWithMessage("Required positional of subcommand " + sc.Name + " named " + pv.Name + " not found at position " + strconv.Itoa(pv.Position))
			exitOrPanic(2)
		}
	}

	return nil
}

// addParsedFlag makes it easy to append flag values parsed by the subcommand
func (sc *Subcommand) addParsedFlag(key string, value string) {
	sc.ParsedValues = append(sc.ParsedValues, newParsedValue(key, value, false))
}

// addParsedPositionalValue makes it easy to append positionals parsed by the
// subcommand
func (sc *Subcommand) addParsedPositionalValue(value string) {
	sc.ParsedValues = append(sc.ParsedValues, newParsedValue("", value, true))
}

// FlagExists lets you know if the flag name exists as either a short or long
// name in the (sub)command
func (sc *Subcommand) FlagExists(name string) bool {

	for _, f := range sc.Flags {
		if f.HasName(name) {
			return true
		}
	}

	return false
}

// AttachSubcommand adds a possible subcommand to the Parser.
func (sc *Subcommand) AttachSubcommand(newSC *Subcommand, relativePosition int) {

	// assign the depth of the subcommand when its attached
	newSC.Position = relativePosition

	// ensure no subcommands at this depth with this name
	for _, other := range sc.Subcommands {
		if newSC.Position == other.Position {
			if newSC.Name != "" {
				if newSC.Name == other.Name {
					log.Panicln("Unable to add subcommand because one already exists at position" + strconv.Itoa(newSC.Position) + " with name " + other.Name)
				}
			}
			if newSC.ShortName != "" {
				if newSC.ShortName == other.ShortName {
					log.Panicln("Unable to add subcommand because one already exists at position" + strconv.Itoa(newSC.Position) + " with name " + other.ShortName)
				}
			}
		}
	}

	// ensure no positionals at this depth
	for _, other := range sc.PositionalFlags {
		if newSC.Position == other.Position {
			log.Panicln("Unable to add subcommand because a positional value already exists at position " + strconv.Itoa(newSC.Position) + ": " + other.Name)
		}
	}

	sc.Subcommands = append(sc.Subcommands, newSC)
}

// add is a "generic" to add flags of any type. Checks the supplied parent
// parser to ensure that the user isn't setting version or help flags that
// conflict with the built-in help and version flag behavior.
func (sc *Subcommand) add(assignmentVar interface{}, shortName string, longName string, description string) {

	// if the flag is already used, throw an error
	for _, existingFlag := range sc.Flags {
		if longName != "" && existingFlag.LongName == longName {
			log.Panicln("Flag " + longName + " added to subcommand " + sc.Name + " but the name is already assigned.")
		}
		if shortName != "" && existingFlag.ShortName == shortName {
			log.Panicln("Flag " + shortName + " added to subcommand " + sc.Name + " but the short name is already assigned.")
		}
	}

	newFlag := Flag{
		AssignmentVar: assignmentVar,
		ShortName:     shortName,
		LongName:      longName,
		Description:   description,
	}
	sc.Flags = append(sc.Flags, &newFlag)
}

// String adds a new string flag
func (sc *Subcommand) String(assignmentVar *string, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// StringSlice adds a new slice of strings flag
// Specify the flag multiple times to fill the slice
func (sc *Subcommand) StringSlice(assignmentVar *[]string, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Bool adds a new bool flag
func (sc *Subcommand) Bool(assignmentVar *bool, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// BoolSlice adds a new slice of bools flag
// Specify the flag multiple times to fill the slice
func (sc *Subcommand) BoolSlice(assignmentVar *[]bool, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// ByteSlice adds a new slice of bytes flag
// Specify the flag multiple times to fill the slice.  Takes hex as input.
func (sc *Subcommand) ByteSlice(assignmentVar *[]byte, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Duration adds a new time.Duration flag.
// Input format is described in time.ParseDuration().
// Example values: 1h, 1h50m, 32s
func (sc *Subcommand) Duration(assignmentVar *time.Duration, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// DurationSlice adds a new time.Duration flag.
// Input format is described in time.ParseDuration().
// Example values: 1h, 1h50m, 32s
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) DurationSlice(assignmentVar *[]time.Duration, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Float32 adds a new float32 flag.
func (sc *Subcommand) Float32(assignmentVar *float32, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Float32Slice adds a new float32 flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) Float32Slice(assignmentVar *[]float32, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Float64 adds a new float64 flag.
func (sc *Subcommand) Float64(assignmentVar *float64, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Float64Slice adds a new float64 flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) Float64Slice(assignmentVar *[]float64, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Int adds a new int flag
func (sc *Subcommand) Int(assignmentVar *int, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// IntSlice adds a new int slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) IntSlice(assignmentVar *[]int, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UInt adds a new uint flag
func (sc *Subcommand) UInt(assignmentVar *uint, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UIntSlice adds a new uint slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) UIntSlice(assignmentVar *[]uint, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UInt64 adds a new uint64 flag
func (sc *Subcommand) UInt64(assignmentVar *uint64, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UInt64Slice adds a new uint64 slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) UInt64Slice(assignmentVar *[]uint64, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UInt32 adds a new uint32 flag
func (sc *Subcommand) UInt32(assignmentVar *uint32, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UInt32Slice adds a new uint32 slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) UInt32Slice(assignmentVar *[]uint32, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UInt16 adds a new uint16 flag
func (sc *Subcommand) UInt16(assignmentVar *uint16, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UInt16Slice adds a new uint16 slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) UInt16Slice(assignmentVar *[]uint16, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UInt8 adds a new uint8 flag
func (sc *Subcommand) UInt8(assignmentVar *uint8, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UInt8Slice adds a new uint8 slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) UInt8Slice(assignmentVar *[]uint8, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Int64 adds a new int64 flag.
func (sc *Subcommand) Int64(assignmentVar *int64, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Int64Slice adds a new int64 slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) Int64Slice(assignmentVar *[]int64, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Int32 adds a new int32 flag
func (sc *Subcommand) Int32(assignmentVar *int32, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Int32Slice adds a new int32 slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) Int32Slice(assignmentVar *[]int32, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Int16 adds a new int16 flag
func (sc *Subcommand) Int16(assignmentVar *int16, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Int16Slice adds a new int16 slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) Int16Slice(assignmentVar *[]int16, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Int8 adds a new int8 flag
func (sc *Subcommand) Int8(assignmentVar *int8, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Int8Slice adds a new int8 slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) Int8Slice(assignmentVar *[]int8, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// IP adds a new net.IP flag.
func (sc *Subcommand) IP(assignmentVar *net.IP, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// IPSlice adds a new int8 slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) IPSlice(assignmentVar *[]net.IP, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// HardwareAddr adds a new net.HardwareAddr flag.
func (sc *Subcommand) HardwareAddr(assignmentVar *net.HardwareAddr, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// HardwareAddrSlice adds a new net.HardwareAddr slice flag.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) HardwareAddrSlice(assignmentVar *[]net.HardwareAddr, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// IPMask adds a new net.IPMask flag. IPv4 Only.
func (sc *Subcommand) IPMask(assignmentVar *net.IPMask, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// IPMaskSlice adds a new net.HardwareAddr slice flag. IPv4 only.
// Specify the flag multiple times to fill the slice.
func (sc *Subcommand) IPMaskSlice(assignmentVar *[]net.IPMask, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// AddPositionalValue adds a positional value to the subcommand.  the
// relativePosition starts at 1 and is relative to the subcommand it belongs to
func (sc *Subcommand) AddPositionalValue(assignmentVar *string, name string, relativePosition int, required bool, description string) {

	// ensure no other positionals are at this depth
	for _, other := range sc.PositionalFlags {
		if relativePosition == other.Position {
			log.Panicln("Unable to add positional value because one already exists at position: " + strconv.Itoa(relativePosition))
		}
	}

	// ensure no subcommands at this depth
	for _, other := range sc.Subcommands {
		if relativePosition == other.Position {
			log.Panicln("Unable to add positional value a subcommand already exists at position: " + strconv.Itoa(relativePosition))
		}
	}

	newPositionalValue := PositionalValue{
		Name:          name,
		Position:      relativePosition,
		AssignmentVar: assignmentVar,
		Required:      required,
		Description:   description,
		defaultValue:  *assignmentVar,
	}
	sc.PositionalFlags = append(sc.PositionalFlags, &newPositionalValue)
}

// SetValueForKey sets the value for the specified key. If setting a bool
// value, then send "true" or "false" as strings.  The returned bool indicates
// that a value was set.
func (sc *Subcommand) SetValueForKey(key string, value string) (bool, error) {

	// debugPrint("Looking to set key", key, "to value", value)
	// check for and assign flags that match the key
	for _, f := range sc.Flags {
		// debugPrint("Evaluating string flag", f.ShortName, "==", key, "||", f.LongName, "==", key)
		if f.ShortName == key || f.LongName == key {
			// debugPrint("Setting string value for", key, "to", value)
			f.identifyAndAssignValue(value)
			return true, nil
		}
	}

	// debugPrint(sc.Name, "was unable to find a key named", key, "to set to value", value)
	return false, nil
}

// ensureNoConflictWithBuiltinHelp ensures that the flags on this subcommand do
// not conflict with the builtin help flags (-h or --help). Exits the program
// if a conflict is found.
func (sc *Subcommand) ensureNoConflictWithBuiltinHelp() {
	for _, f := range sc.Flags {
		if f.LongName == helpFlagLongName {
			sc.exitBecauseOfHelpFlagConflict(f.LongName)
		}
		if f.LongName == helpFlagShortName {
			sc.exitBecauseOfHelpFlagConflict(f.LongName)
		}
		if f.ShortName == helpFlagLongName {
			sc.exitBecauseOfHelpFlagConflict(f.ShortName)
		}
		if f.ShortName == helpFlagShortName {
			sc.exitBecauseOfHelpFlagConflict(f.ShortName)
		}
	}
}

// ensureNoConflictWithBuiltinVersion ensures that the flags on this subcommand do
// not conflict with the builtin version flag (--version). Exits the program
// if a conflict is found.
func (sc *Subcommand) ensureNoConflictWithBuiltinVersion() {
	for _, f := range sc.Flags {
		if f.LongName == versionFlagLongName {
			sc.exitBecauseOfVersionFlagConflict(f.LongName)
		}
		if f.ShortName == versionFlagLongName {
			sc.exitBecauseOfVersionFlagConflict(f.ShortName)
		}
	}
}

// exitBecauseOfVersionFlagConflict exits the program with a message about how to prevent
// flags being defined from conflicting with the builtin flags.
func (sc *Subcommand) exitBecauseOfVersionFlagConflict(flagName string) {
	fmt.Println(`Flag with name '` + flagName + `' conflicts with the internal --version flag in flaggy.

You must either change the flag's name, or disable flaggy's internal version
flag with 'flaggy.DefaultParser.ShowVersionWithVersionFlag = false'.  If you are using
a custom parser, you must instead set '.ShowVersionWithVersionFlag = false' on it.`)
	exitOrPanic(1)
}

// exitBecauseOfHelpFlagConflict exits the program with a message about how to prevent
// flags being defined from conflicting with the builtin flags.
func (sc *Subcommand) exitBecauseOfHelpFlagConflict(flagName string) {
	fmt.Println(`Flag with name '` + flagName + `' conflicts with the internal --help or -h flag in flaggy.

You must either change the flag's name, or disable flaggy's internal help
flag with 'flaggy.DefaultParser.ShowHelpWithHFlag = false'.  If you are using
a custom parser, you must instead set '.ShowHelpWithHFlag = false' on it.`)
	exitOrPanic(1)
}
