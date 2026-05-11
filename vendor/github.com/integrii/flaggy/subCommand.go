package flaggy

import (
	"fmt"
	"log"
	"math/big"
	"net"
	netip "net/netip"
	"net/url"
	"os"
	"regexp"
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
func (sc *Subcommand) parseAllFlagsFromArgs(p *Parser, args []string) (flagScanResult, error) {

	result := flagScanResult{}
	positionalCount := 0

	for i := 0; i < len(args); i++ {
		a := args[i]

		debugPrint("parsing arg:", a)

		argType := determineArgType(a)

		if argType == argIsFinal {
			if !p.trailingArgumentsExtracted {
				p.TrailingArguments = append(p.TrailingArguments, args[i+1:]...)
			}
			break
		}

		flagName := parseFlagToName(a)

		if p.ShowVersionWithVersionFlag && flagName == versionFlagLongName {
			p.ShowVersionAndExit()
		}

		if p.ShowHelpWithHFlag && (flagName == helpFlagShortName || flagName == helpFlagLongName) {
			result.HelpRequested = true
			continue
		}

		debugPrint("Parsing flag named", a, "of type", argType)

		switch argType {
		case argIsPositional:
			positionalCount++
			token := positionalToken{Value: a, Index: i}
			result.Positionals = append(result.Positionals, token)
			sc.addParsedPositionalValue(a)

			// Detect subcommands early so we avoid parsing child flags at this level.
			var matched *Subcommand
			for _, cmd := range sc.Subcommands {
				if a == cmd.Name || a == cmd.ShortName {
					// Prefer an exact positional match when available.
					if cmd.Position == positionalCount {
						matched = cmd
						break
					}
					if matched == nil {
						matched = cmd
					}
				}
			}
			if matched != nil {
				// Ignore the tentative match when a positional value is already defined at this depth.
				if matched.Position != positionalCount && hasPositionalAtDepth(sc, positionalCount) {
					matched = nil
				}
			}
			if matched != nil {
				// Drop the provisional positional bookkeeping because the token actually belongs to the child.
				if len(result.Positionals) > 0 {
					result.Positionals = result.Positionals[:len(result.Positionals)-1]
				}
				if len(sc.ParsedValues) > 0 {
					lastIdx := len(sc.ParsedValues) - 1
					if sc.ParsedValues[lastIdx].IsPositional {
						sc.ParsedValues = sc.ParsedValues[:lastIdx]
					}
				}
				// Record which subcommand will own the remainder of the arguments.
				result.Subcommand = &subcommandMatch{
					Command:       matched,
					Token:         token,
					RelativeDepth: matched.Position,
				}
				// Stop scanning so the child can handle the remainder.
				return result, nil
			}
		case argIsFlagWithSpace:
			key := flagName

			if flagIsBool(sc, p, key) {
				valueSet, err := setValueForParsers(key, "true", p, sc)
				if err != nil {
					return result, err
				}
				if valueSet {
					sc.addParsedFlag(key, "", false)
				}
				continue
			}

			if !flagIsDefined(sc, p, key) {
				result.ForwardArgs = append(result.ForwardArgs, args[i])
				if i+1 < len(args) && shouldReserveNextArgForChild(sc, positionalCount, args[i+1]) {
					result.ForwardArgs = append(result.ForwardArgs, args[i+1])
					i++
				}
				continue
			}

			if i+1 >= len(args) {
				p.ShowHelpWithMessage("Expected a following arg for flag " + key + ", but it did not exist.")
				exitOrPanic(2)
			}

			nextArg := args[i+1]
			valueSet, err := setValueForParsers(key, nextArg, p, sc)
			if err != nil {
				return result, err
			}
			if valueSet {
				sc.addParsedFlag(key, nextArg, true)
			}
			i++
		case argIsFlagWithValue:
			keyWithValue := flagName
			key, val := parseArgWithValue(keyWithValue)

			if !flagIsDefined(sc, p, key) {
				result.ForwardArgs = append(result.ForwardArgs, args[i])
				continue
			}

			valueSet, err := setValueForParsers(key, val, p, sc)
			if err != nil {
				return result, err
			}
			if valueSet {
				sc.addParsedFlag(keyWithValue, val, false)
			}
		}
	}

	return result, nil
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

// shouldReserveNextArgForChild determines if the following argument should be
// left untouched so that a downstream subcommand can parse it as its own flag
// or positional value.
func shouldReserveNextArgForChild(sc *Subcommand, positionalCount int, nextArg string) bool {
	if determineArgType(nextArg) == argIsFinal {
		return false
	}

	position := positionalCount + 1
	for _, cmd := range sc.Subcommands {
		if cmd.Position == position && (cmd.Name == nextArg || cmd.ShortName == nextArg) {
			return false
		}
	}

	for _, pos := range sc.PositionalFlags {
		if pos.Position == position {
			return false
		}
	}

	return true
}

func hasPositionalAtDepth(sc *Subcommand, depth int) bool {
	for _, pos := range sc.PositionalFlags {
		if pos.Position == depth {
			return true
		}
	}
	return false
}

// parse causes the argument parser to parse based on the supplied []string.
// The args slice should contain only values that have not already been
// consumed by parent parsers. The parser records any values it parses so that
// the root parser can detect unexpected arguments after parsing is complete.
func (sc *Subcommand) parse(p *Parser, args []string) error {

	debugPrint("- Parsing subcommand", sc.Name, "with args", args)

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

	scan, err := sc.parseAllFlagsFromArgs(p, args)
	if err != nil {
		return err
	}

	for idx, token := range scan.Positionals {
		relativeDepth := idx + 1
		value := token.Value

		if scan.Subcommand != nil && token.Index == scan.Subcommand.Token.Index {
			debugPrint("Descending into positional subcommand", scan.Subcommand.Command.Name, "at relativeDepth", scan.Subcommand.RelativeDepth)
			childArgs := append([]string{}, scan.ForwardArgs...)
			childArgs = append(childArgs, args[token.Index+1:]...)
			return scan.Subcommand.Command.parse(p, childArgs)
		}

		var foundPositional bool
		for _, val := range sc.PositionalFlags {
			if relativeDepth == val.Position {
				debugPrint("Found a positional value at relativePos:", relativeDepth, "value:", value)

				val.defaultValue = *val.AssignmentVar
				*val.AssignmentVar = value
				foundPositional = true
				val.Found = true
				break
			}
		}

		if !foundPositional {
			if p.ShowHelpOnUnexpected {
				debugPrint("No positional at position", relativeDepth)
				var foundSubcommandAtDepth bool
				for _, cmd := range sc.Subcommands {
					if cmd.Position == relativeDepth {
						foundSubcommandAtDepth = true
					}
				}

				if foundSubcommandAtDepth {
					fmt.Fprintln(os.Stderr, sc.Name+":", "No subcommand or positional value found at position", strconv.Itoa(relativeDepth)+".")
					var output string
					for _, cmd := range sc.Subcommands {
						if cmd.Hidden {
							continue
						}
						output = output + " " + cmd.Name
					}
					if len(output) > 0 {
						output = strings.TrimLeft(output, " ")
						fmt.Println("Available subcommands:", output)
					}
					exitOrPanic(2)
				}

				p.ShowHelpWithMessage("Unexpected argument: " + value)
				exitOrPanic(2)
			} else {
				p.TrailingArguments = append(p.TrailingArguments, value)
			}
		}
	}

	if scan.Subcommand != nil {
		// If we recorded a subcommand but didn't descend, ensure the remaining
		// arguments are handed off now.
		debugPrint("Descending into positional subcommand", scan.Subcommand.Command.Name, "at relativeDepth", scan.Subcommand.RelativeDepth)
		childArgs := append([]string{}, scan.ForwardArgs...)
		childArgs = append(childArgs, args[scan.Subcommand.Token.Index+1:]...)
		return scan.Subcommand.Command.parse(p, childArgs)
	}

	if scan.HelpRequested && p.ShowHelpWithHFlag {
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

	// indicate that trailing arguments have been extracted, so that they aren't
	// appended a second time by parent parsers.
	p.trailingArgumentsExtracted = true

	return nil
}

// addParsedFlag makes it easy to append flag values parsed by the subcommand
func (sc *Subcommand) addParsedFlag(key string, value string, consumesNext bool) {
	sc.ParsedValues = append(sc.ParsedValues, newParsedValue(key, value, false, consumesNext))
}

// addParsedPositionalValue makes it easy to append positionals parsed by the
// subcommand
func (sc *Subcommand) addParsedPositionalValue(value string) {
	sc.ParsedValues = append(sc.ParsedValues, newParsedValue("", value, true, false))
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

// BytesBase64 adds a new []byte flag parsed from base64 input.
func (sc *Subcommand) BytesBase64(assignmentVar *Base64Bytes, shortName string, longName string, description string) {
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

// Time adds a new time.Time flag. Supports RFC3339/RFC3339Nano, RFC1123, and unix seconds.
func (sc *Subcommand) Time(assignmentVar *time.Time, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// URL adds a new url.URL flag.
func (sc *Subcommand) URL(assignmentVar *url.URL, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// IPNet adds a new net.IPNet flag parsed from CIDR.
func (sc *Subcommand) IPNet(assignmentVar *net.IPNet, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// TCPAddr adds a new net.TCPAddr flag parsed from host:port.
func (sc *Subcommand) TCPAddr(assignmentVar *net.TCPAddr, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// UDPAddr adds a new net.UDPAddr flag parsed from host:port.
func (sc *Subcommand) UDPAddr(assignmentVar *net.UDPAddr, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// FileMode adds a new os.FileMode flag parsed from octal/decimal (base auto-detected).
func (sc *Subcommand) FileMode(assignmentVar *os.FileMode, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Regexp adds a new regexp.Regexp flag.
func (sc *Subcommand) Regexp(assignmentVar *regexp.Regexp, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Location adds a new time.Location flag.
func (sc *Subcommand) Location(assignmentVar *time.Location, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Month adds a new time.Month flag.
func (sc *Subcommand) Month(assignmentVar *time.Month, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// Weekday adds a new time.Weekday flag.
func (sc *Subcommand) Weekday(assignmentVar *time.Weekday, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// BigInt adds a new big.Int flag.
func (sc *Subcommand) BigInt(assignmentVar *big.Int, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// BigRat adds a new big.Rat flag.
func (sc *Subcommand) BigRat(assignmentVar *big.Rat, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// NetipAddr adds a new netip.Addr flag.
func (sc *Subcommand) NetipAddr(assignmentVar *netip.Addr, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// NetipPrefix adds a new netip.Prefix flag.
func (sc *Subcommand) NetipPrefix(assignmentVar *netip.Prefix, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// NetipAddrPort adds a new netip.AddrPort flag.
func (sc *Subcommand) NetipAddrPort(assignmentVar *netip.AddrPort, shortName string, longName string, description string) {
	sc.add(assignmentVar, shortName, longName, description)
}

// AddPositionalValue adds a positional value to the subcommand.  the
// relativePosition starts at 1 and is relative to the subcommand it belongs to
func (sc *Subcommand) AddPositionalValue(assignmentVar *string, name string, relativePosition int, required bool, description string) {

	// ensure no other positionals are at this depth
	for _, other := range sc.PositionalFlags {
		if relativePosition == other.Position {
			log.Panicln("Unable to add positional value " + name + " because " + other.Name + " already exists at position: " + strconv.Itoa(relativePosition))
		}
	}

	// ensure no subcommands at this depth
	for _, other := range sc.Subcommands {
		if relativePosition == other.Position {
			log.Panicln("Unable to add positional value " + name + "because a subcommand, " + other.Name + ", already exists at position: " + strconv.Itoa(relativePosition))
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
			if err := f.identifyAndAssignValue(value); err != nil {
				return false, err
			}
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
