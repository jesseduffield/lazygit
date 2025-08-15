// Package flaggy is a input flag parsing package that supports recursive
// subcommands, positional values, and any-position flags without
// unnecessary complexeties.
//
// For a getting started tutorial and full feature list, check out the
// readme at https://github.com/integrii/flaggy.
package flaggy // import "github.com/integrii/flaggy"

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// strings used for builtin help and version flags both short and long
const versionFlagLongName = "version"
const helpFlagLongName = "help"
const helpFlagShortName = "h"

// defaultVersion is applied to parsers when they are created
const defaultVersion = "0.0.0"

// DebugMode indicates that debug output should be enabled
var DebugMode bool

// DefaultHelpTemplate is the help template that will be used
// for newly created subcommands and commands
var DefaultHelpTemplate = defaultHelpTemplate

// DefaultParser is the default parser that is used with the package-level public
// functions
var DefaultParser *Parser

// TrailingArguments holds trailing arguments in the main parser after parsing
// has been run.
var TrailingArguments []string

func init() {

	// set the default help template
	// allow usage like flaggy.StringVar by enabling a default Parser
	ResetParser()
}

// ResetParser resets the default parser to a fresh instance.  Uses the
// name of the binary executing as the program name by default.
func ResetParser() {
	if len(os.Args) > 0 {
		chunks := strings.Split(os.Args[0], "/")
		DefaultParser = NewParser(chunks[len(chunks)-1])
	} else {
		DefaultParser = NewParser("default")
	}
}

// Parse parses flags as requested in the default package parser
func Parse() {
	err := DefaultParser.Parse()
	TrailingArguments = DefaultParser.TrailingArguments
	if err != nil {
		log.Panicln("Error from argument parser:", err)
	}
}

// ParseArgs parses the passed args as if they were the arguments to the
// running binary.  Targets the default main parser for the package.
func ParseArgs(args []string) {
	err := DefaultParser.ParseArgs(args)
	TrailingArguments = DefaultParser.TrailingArguments
	if err != nil {
		log.Panicln("Error from argument parser:", err)
	}
}

// String adds a new string flag
func String(assignmentVar *string, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// StringSlice adds a new slice of strings flag
// Specify the flag multiple times to fill the slice
func StringSlice(assignmentVar *[]string, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Bool adds a new bool flag
func Bool(assignmentVar *bool, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// BoolSlice adds a new slice of bools flag
// Specify the flag multiple times to fill the slice
func BoolSlice(assignmentVar *[]bool, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// ByteSlice adds a new slice of bytes flag
// Specify the flag multiple times to fill the slice.  Takes hex as input.
func ByteSlice(assignmentVar *[]byte, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Duration adds a new time.Duration flag.
// Input format is described in time.ParseDuration().
// Example values: 1h, 1h50m, 32s
func Duration(assignmentVar *time.Duration, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// DurationSlice adds a new time.Duration flag.
// Input format is described in time.ParseDuration().
// Example values: 1h, 1h50m, 32s
// Specify the flag multiple times to fill the slice.
func DurationSlice(assignmentVar *[]time.Duration, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Float32 adds a new float32 flag.
func Float32(assignmentVar *float32, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Float32Slice adds a new float32 flag.
// Specify the flag multiple times to fill the slice.
func Float32Slice(assignmentVar *[]float32, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Float64 adds a new float64 flag.
func Float64(assignmentVar *float64, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Float64Slice adds a new float64 flag.
// Specify the flag multiple times to fill the slice.
func Float64Slice(assignmentVar *[]float64, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Int adds a new int flag
func Int(assignmentVar *int, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// IntSlice adds a new int slice flag.
// Specify the flag multiple times to fill the slice.
func IntSlice(assignmentVar *[]int, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UInt adds a new uint flag
func UInt(assignmentVar *uint, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UIntSlice adds a new uint slice flag.
// Specify the flag multiple times to fill the slice.
func UIntSlice(assignmentVar *[]uint, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UInt64 adds a new uint64 flag
func UInt64(assignmentVar *uint64, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UInt64Slice adds a new uint64 slice flag.
// Specify the flag multiple times to fill the slice.
func UInt64Slice(assignmentVar *[]uint64, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UInt32 adds a new uint32 flag
func UInt32(assignmentVar *uint32, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UInt32Slice adds a new uint32 slice flag.
// Specify the flag multiple times to fill the slice.
func UInt32Slice(assignmentVar *[]uint32, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UInt16 adds a new uint16 flag
func UInt16(assignmentVar *uint16, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UInt16Slice adds a new uint16 slice flag.
// Specify the flag multiple times to fill the slice.
func UInt16Slice(assignmentVar *[]uint16, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UInt8 adds a new uint8 flag
func UInt8(assignmentVar *uint8, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// UInt8Slice adds a new uint8 slice flag.
// Specify the flag multiple times to fill the slice.
func UInt8Slice(assignmentVar *[]uint8, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Int64 adds a new int64 flag
func Int64(assignmentVar *int64, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Int64Slice adds a new int64 slice flag.
// Specify the flag multiple times to fill the slice.
func Int64Slice(assignmentVar *[]int64, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Int32 adds a new int32 flag
func Int32(assignmentVar *int32, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Int32Slice adds a new int32 slice flag.
// Specify the flag multiple times to fill the slice.
func Int32Slice(assignmentVar *[]int32, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Int16 adds a new int16 flag
func Int16(assignmentVar *int16, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Int16Slice adds a new int16 slice flag.
// Specify the flag multiple times to fill the slice.
func Int16Slice(assignmentVar *[]int16, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Int8 adds a new int8 flag
func Int8(assignmentVar *int8, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// Int8Slice adds a new int8 slice flag.
// Specify the flag multiple times to fill the slice.
func Int8Slice(assignmentVar *[]int8, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// IP adds a new net.IP flag.
func IP(assignmentVar *net.IP, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// IPSlice adds a new int8 slice flag.
// Specify the flag multiple times to fill the slice.
func IPSlice(assignmentVar *[]net.IP, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// HardwareAddr adds a new net.HardwareAddr flag.
func HardwareAddr(assignmentVar *net.HardwareAddr, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// HardwareAddrSlice adds a new net.HardwareAddr slice flag.
// Specify the flag multiple times to fill the slice.
func HardwareAddrSlice(assignmentVar *[]net.HardwareAddr, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// IPMask adds a new net.IPMask flag. IPv4 Only.
func IPMask(assignmentVar *net.IPMask, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// IPMaskSlice adds a new net.HardwareAddr slice flag. IPv4 only.
// Specify the flag multiple times to fill the slice.
func IPMaskSlice(assignmentVar *[]net.IPMask, shortName string, longName string, description string) {
	DefaultParser.add(assignmentVar, shortName, longName, description)
}

// AttachSubcommand adds a subcommand for parsing
func AttachSubcommand(subcommand *Subcommand, relativePosition int) {
	DefaultParser.AttachSubcommand(subcommand, relativePosition)
}

// ShowHelp shows parser help
func ShowHelp(message string) {
	DefaultParser.ShowHelpWithMessage(message)
}

// SetDescription sets the description of the default package command parser
func SetDescription(description string) {
	DefaultParser.Description = description
}

// SetVersion sets the version of the default package command parser
func SetVersion(version string) {
	DefaultParser.Version = version
}

// SetName sets the name of the default package command parser
func SetName(name string) {
	DefaultParser.Name = name
}

// ShowHelpAndExit shows parser help and exits with status code 2
func ShowHelpAndExit(message string) {
	ShowHelp(message)
	exitOrPanic(2)
}

// PanicInsteadOfExit is used when running tests
var PanicInsteadOfExit bool

// exitOrPanic panics instead of calling os.Exit so that tests can catch
// more failures
func exitOrPanic(code int) {
	if PanicInsteadOfExit {
		panic("Panic instead of exit with code: " + strconv.Itoa(code))
	}
	os.Exit(code)
}

// ShowHelpOnUnexpectedEnable enables the ShowHelpOnUnexpected behavior on the
// default parser.  This causes unknown inputs to error out.
func ShowHelpOnUnexpectedEnable() {
	DefaultParser.ShowHelpOnUnexpected = true
}

// ShowHelpOnUnexpectedDisable disables the ShowHelpOnUnexpected behavior on the
// default parser.  This causes unknown inputs to error out.
func ShowHelpOnUnexpectedDisable() {
	DefaultParser.ShowHelpOnUnexpected = false
}

// AddPositionalValue adds a positional value to the main parser at the global
// context
func AddPositionalValue(assignmentVar *string, name string, relativePosition int, required bool, description string) {
	DefaultParser.AddPositionalValue(assignmentVar, name, relativePosition, required, description)
}

// debugPrint prints if debugging is enabled
func debugPrint(i ...interface{}) {
	if DebugMode {
		fmt.Println(i...)
	}
}
