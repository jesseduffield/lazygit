package flaggy

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Flag holds the base methods for all flag types
type Flag struct {
	ShortName     string
	LongName      string
	Description   string
	rawValue      string // the value as a string before being parsed
	Hidden        bool   // indicates this flag should be hidden from help and suggestions
	AssignmentVar interface{}
	defaultValue  string // the value (as a string), that was set by default before any parsing and assignment
	parsed        bool   // indicates that this flag has already been parsed
}

// HasName indicates that this flag's short or long name matches the
// supplied name string
func (f *Flag) HasName(name string) bool {
	name = strings.TrimSpace(name)
	if f.ShortName == name || f.LongName == name {
		return true
	}
	return false
}

// identifyAndAssignValue identifies the type of the incoming value
// and assigns it to the AssignmentVar pointer's target value.  If
// the value is a type that needs parsing, that is performed as well.
func (f *Flag) identifyAndAssignValue(value string) error {

	var err error

	// Only parse this flag default value once. This keeps us from
	// overwriting the default value in help output
	if !f.parsed {
		f.parsed = true
		// parse the default value as a string and remember it for help output
		f.defaultValue, err = f.returnAssignmentVarValueAsString()
		if err != nil {
			return err
		}
	}

	debugPrint("attempting to assign value", value, "to flag", f.LongName)
	f.rawValue = value // remember the raw value

	// depending on the type of the assignment variable, we convert the
	// incoming string and assign it.  We only use pointers to variables
	// in flagy.  No returning vars by value.
	switch f.AssignmentVar.(type) {
	case *string:
		v, _ := (f.AssignmentVar).(*string)
		*v = value
	case *[]string:
		v := f.AssignmentVar.(*[]string)
		splitString := strings.Split(value, ",")
		new := append(*v, splitString...)
		*v = new
	case *bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		a, _ := (f.AssignmentVar).(*bool)
		*a = v
	case *[]bool:
		// parse the incoming bool
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		// cast the assignment var
		existing := f.AssignmentVar.(*[]bool)
		// deref the assignment var and append to it
		v := append(*existing, b)
		// pointer the new value and assign it
		a, _ := (f.AssignmentVar).(*[]bool)
		*a = v
	case *time.Duration:
		v, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		a, _ := (f.AssignmentVar).(*time.Duration)
		*a = v
	case *[]time.Duration:
		t, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*[]time.Duration)
		// deref the assignment var and append to it
		v := append(*existing, t)
		// pointer the new value and assign it
		a, _ := (f.AssignmentVar).(*[]time.Duration)
		*a = v
	case *float32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		float := float32(v)
		a, _ := (f.AssignmentVar).(*float32)
		*a = float
	case *[]float32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		float := float32(v)
		existing := f.AssignmentVar.(*[]float32)
		new := append(*existing, float)
		*existing = new
	case *float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		a, _ := (f.AssignmentVar).(*float64)
		*a = v
	case *[]float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*[]float64)
		new := append(*existing, v)

		*existing = new
	case *int:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		e := f.AssignmentVar.(*int)
		*e = v
	case *[]int:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*[]int)
		new := append(*existing, v)
		*existing = new
	case *uint:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*uint)
		*existing = uint(v)
	case *[]uint:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*[]uint)
		new := append(*existing, uint(v))
		*existing = new
	case *uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*uint64)
		*existing = v
	case *[]uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*[]uint64)
		new := append(*existing, v)
		*existing = new
	case *uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*uint32)
		*existing = uint32(v)
	case *[]uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*[]uint32)
		new := append(*existing, uint32(v))
		*existing = new
	case *uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return err
		}
		val := uint16(v)
		existing := f.AssignmentVar.(*uint16)
		*existing = val
	case *[]uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*[]uint16)
		new := append(*existing, uint16(v))
		*existing = new
	case *uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return err
		}
		val := uint8(v)
		existing := f.AssignmentVar.(*uint8)
		*existing = val
	case *[]uint8:
		var newSlice []uint8

		v, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return err
		}
		newV := uint8(v)
		existing := f.AssignmentVar.(*[]uint8)
		newSlice = append(*existing, newV)
		*existing = newSlice
	case *int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*int64)
		*existing = v
	case *[]int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		existingSlice := f.AssignmentVar.(*[]int64)
		newSlice := append(*existingSlice, v)
		*existingSlice = newSlice
	case *int32:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		converted := int32(v)
		existing := f.AssignmentVar.(*int32)
		*existing = converted
	case *[]int32:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		existingSlice := f.AssignmentVar.(*[]int32)
		newSlice := append(*existingSlice, int32(v))
		*existingSlice = newSlice
	case *int16:
		v, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}
		converted := int16(v)
		existing := f.AssignmentVar.(*int16)
		*existing = converted
	case *[]int16:
		v, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}
		existingSlice := f.AssignmentVar.(*[]int16)
		newSlice := append(*existingSlice, int16(v))
		*existingSlice = newSlice
	case *int8:
		v, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return err
		}
		converted := int8(v)
		existing := f.AssignmentVar.(*int8)
		*existing = converted
	case *[]int8:
		v, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return err
		}
		existingSlice := f.AssignmentVar.(*[]int8)
		newSlice := append(*existingSlice, int8(v))
		*existingSlice = newSlice
	case *net.IP:
		v := net.ParseIP(value)
		existing := f.AssignmentVar.(*net.IP)
		*existing = v
	case *[]net.IP:
		v := net.ParseIP(value)
		existing := f.AssignmentVar.(*[]net.IP)
		new := append(*existing, v)
		*existing = new
	case *net.HardwareAddr:
		v, err := net.ParseMAC(value)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*net.HardwareAddr)
		*existing = v
	case *[]net.HardwareAddr:
		v, err := net.ParseMAC(value)
		if err != nil {
			return err
		}
		existing := f.AssignmentVar.(*[]net.HardwareAddr)
		new := append(*existing, v)
		*existing = new
	case *net.IPMask:
		v := net.IPMask(net.ParseIP(value).To4())
		existing := f.AssignmentVar.(*net.IPMask)
		*existing = v
	case *[]net.IPMask:
		v := net.IPMask(net.ParseIP(value).To4())
		existing := f.AssignmentVar.(*[]net.IPMask)
		new := append(*existing, v)
		*existing = new
	default:
		return errors.New("Unknown flag assignmentVar supplied in flag " + f.LongName + " " + f.ShortName)
	}

	return err
}

const argIsPositional = "positional"       // subcommand or positional value
const argIsFlagWithSpace = "flagWithSpace" // -f path or --file path
const argIsFlagWithValue = "flagWithValue" // -f=path or --file=path
const argIsFinal = "final"                 // the final argument only '--'

// determineArgType determines if the specified arg is a flag with space
// separated value, a flag with a connected value, or neither (positional)
func determineArgType(arg string) string {

	// if the arg is --, then its the final arg
	if arg == "--" {
		return argIsFinal
	}

	// if it has the prefix --, then its a long flag
	if strings.HasPrefix(arg, "--") {
		// if it contains an equals, it is a joined value
		if strings.Contains(arg, "=") {
			return argIsFlagWithValue
		}
		return argIsFlagWithSpace
	}

	// if it has the prefix -, then its a short flag
	if strings.HasPrefix(arg, "-") {
		// if it contains an equals, it is a joined value
		if strings.Contains(arg, "=") {
			return argIsFlagWithValue
		}
		return argIsFlagWithSpace
	}

	return argIsPositional
}

// parseArgWithValue parses a key=value concatenated argument into a key and
// value
func parseArgWithValue(arg string) (key string, value string) {

	// remove up to two minuses from start of flag
	arg = strings.TrimPrefix(arg, "-")
	arg = strings.TrimPrefix(arg, "-")

	// debugPrint("parseArgWithValue parsing", arg)

	// break at the equals
	args := strings.SplitN(arg, "=", 2)

	// if its a bool arg, with no explicit value, we return a blank
	if len(args) == 1 {
		return args[0], ""
	}

	// if its a key and value pair, we return those
	if len(args) == 2 {
		// debugPrint("parseArgWithValue parsed", args[0], args[1])
		return args[0], args[1]
	}

	fmt.Println("Warning: attempted to parseArgWithValue but did not have correct parameter count.", arg, "->", args)
	return "", ""
}

// parseFlagToName parses a flag with space value down to a key name:
//     --path -> path
//     -p -> p
func parseFlagToName(arg string) string {
	// remove minus from start
	arg = strings.TrimLeft(arg, "-")
	arg = strings.TrimLeft(arg, "-")
	return arg
}

// flagIsBool determines if the flag is a bool within the specified parser
// and subcommand's context
func flagIsBool(sc *Subcommand, p *Parser, key string) bool {
	for _, f := range append(sc.Flags, p.Flags...) {
		if f.HasName(key) {
			_, isBool := f.AssignmentVar.(*bool)
			_, isBoolSlice := f.AssignmentVar.(*[]bool)
			if isBool || isBoolSlice {
				return true
			}
		}
	}

	// by default, the answer is false
	return false
}

// returnAssignmentVarValueAsString returns the value of the flag's
// assignment variable as a string.  This is used to display the
// default value of flags before they are assigned (like when help is output).
func (f *Flag) returnAssignmentVarValueAsString() (string, error) {

	debugPrint("returning current value of assignment var of flag", f.LongName)

	var err error

	// depending on the type of the assignment variable, we convert the
	// incoming string and assign it.  We only use pointers to variables
	// in flagy.  No returning vars by value.
	switch f.AssignmentVar.(type) {
	case *string:
		v, _ := (f.AssignmentVar).(*string)
		return *v, err
	case *[]string:
		v := f.AssignmentVar.(*[]string)
		return strings.Join(*v, ","), err
	case *bool:
		a, _ := (f.AssignmentVar).(*bool)
		return strconv.FormatBool(*a), err
	case *[]bool:
		value := f.AssignmentVar.(*[]bool)
		var ss []string
		for _, b := range *value {
			ss = append(ss, strconv.FormatBool(b))
		}
		return strings.Join(ss, ","), err
	case *time.Duration:
		a := f.AssignmentVar.(*time.Duration)
		return (*a).String(), err
	case *[]time.Duration:
		tds := f.AssignmentVar.(*[]time.Duration)
		var asSlice []string
		for _, td := range *tds {
			asSlice = append(asSlice, td.String())
		}
		return strings.Join(asSlice, ","), err
	case *float32:
		a := f.AssignmentVar.(*float32)
		return strconv.FormatFloat(float64(*a), 'f', 2, 32), err
	case *[]float32:
		v := f.AssignmentVar.(*[]float32)
		var strSlice []string
		for _, f := range *v {
			formatted := strconv.FormatFloat(float64(f), 'f', 2, 32)
			strSlice = append(strSlice, formatted)
		}
		return strings.Join(strSlice, ","), err
	case *float64:
		a := f.AssignmentVar.(*float64)
		return strconv.FormatFloat(float64(*a), 'f', 2, 64), err
	case *[]float64:
		v := f.AssignmentVar.(*[]float64)
		var strSlice []string
		for _, f := range *v {
			formatted := strconv.FormatFloat(float64(f), 'f', 2, 64)
			strSlice = append(strSlice, formatted)
		}
		return strings.Join(strSlice, ","), err
	case *int:
		a := f.AssignmentVar.(*int)
		return strconv.Itoa(*a), err
	case *[]int:
		val := f.AssignmentVar.(*[]int)
		var strSlice []string
		for _, i := range *val {
			str := strconv.Itoa(i)
			strSlice = append(strSlice, str)
		}
		return strings.Join(strSlice, ","), err
	case *uint:
		v := f.AssignmentVar.(*uint)
		return strconv.FormatUint(uint64(*v), 10), err
	case *[]uint:
		values := f.AssignmentVar.(*[]uint)
		var strVars []string
		for _, i := range *values {
			strVars = append(strVars, strconv.FormatUint(uint64(i), 10))
		}
		return strings.Join(strVars, ","), err
	case *uint64:
		v := f.AssignmentVar.(*uint64)
		return strconv.FormatUint(*v, 10), err
	case *[]uint64:
		values := f.AssignmentVar.(*[]uint64)
		var strVars []string
		for _, i := range *values {
			strVars = append(strVars, strconv.FormatUint(i, 10))
		}
		return strings.Join(strVars, ","), err
	case *uint32:
		v := f.AssignmentVar.(*uint32)
		return strconv.FormatUint(uint64(*v), 10), err
	case *[]uint32:
		values := f.AssignmentVar.(*[]uint32)
		var strVars []string
		for _, i := range *values {
			strVars = append(strVars, strconv.FormatUint(uint64(i), 10))
		}
		return strings.Join(strVars, ","), err
	case *uint16:
		v := f.AssignmentVar.(*uint16)
		return strconv.FormatUint(uint64(*v), 10), err
	case *[]uint16:
		values := f.AssignmentVar.(*[]uint16)
		var strVars []string
		for _, i := range *values {
			strVars = append(strVars, strconv.FormatUint(uint64(i), 10))
		}
		return strings.Join(strVars, ","), err
	case *uint8:
		v := f.AssignmentVar.(*uint8)
		return strconv.FormatUint(uint64(*v), 10), err
	case *[]uint8:
		values := f.AssignmentVar.(*[]uint8)
		var strVars []string
		for _, i := range *values {
			strVars = append(strVars, strconv.FormatUint(uint64(i), 10))
		}
		return strings.Join(strVars, ","), err
	case *int64:
		v := f.AssignmentVar.(*int64)
		return strconv.FormatInt(int64(*v), 10), err
	case *[]int64:
		values := f.AssignmentVar.(*[]int64)
		var strVars []string
		for _, i := range *values {
			strVars = append(strVars, strconv.FormatInt(i, 10))
		}
		return strings.Join(strVars, ","), err
	case *int32:
		v := f.AssignmentVar.(*int32)
		return strconv.FormatInt(int64(*v), 10), err
	case *[]int32:
		values := f.AssignmentVar.(*[]int32)
		var strVars []string
		for _, i := range *values {
			strVars = append(strVars, strconv.FormatInt(int64(i), 10))
		}
		return strings.Join(strVars, ","), err
	case *int16:
		v := f.AssignmentVar.(*int16)
		return strconv.FormatInt(int64(*v), 10), err
	case *[]int16:
		values := f.AssignmentVar.(*[]int16)
		var strVars []string
		for _, i := range *values {
			strVars = append(strVars, strconv.FormatInt(int64(i), 10))
		}
		return strings.Join(strVars, ","), err
	case *int8:
		v := f.AssignmentVar.(*int8)
		return strconv.FormatInt(int64(*v), 10), err
	case *[]int8:
		values := f.AssignmentVar.(*[]int8)
		var strVars []string
		for _, i := range *values {
			strVars = append(strVars, strconv.FormatInt(int64(i), 10))
		}
		return strings.Join(strVars, ","), err
	case *net.IP:
		val := f.AssignmentVar.(*net.IP)
		return val.String(), err
	case *[]net.IP:
		val := f.AssignmentVar.(*[]net.IP)
		var strSlice []string
		for _, ip := range *val {
			strSlice = append(strSlice, ip.String())
		}
		return strings.Join(strSlice, ","), err
	case *net.HardwareAddr:
		val := f.AssignmentVar.(*net.HardwareAddr)
		return val.String(), err
	case *[]net.HardwareAddr:
		val := f.AssignmentVar.(*[]net.HardwareAddr)
		var strSlice []string
		for _, mac := range *val {
			strSlice = append(strSlice, mac.String())
		}
		return strings.Join(strSlice, ","), err
	case *net.IPMask:
		val := f.AssignmentVar.(*net.IPMask)
		return val.String(), err
	case *[]net.IPMask:
		val := f.AssignmentVar.(*[]net.IPMask)
		var strSlice []string
		for _, m := range *val {
			strSlice = append(strSlice, m.String())
		}
		return strings.Join(strSlice, ","), err
	default:
		return "", errors.New("Unknown flag assignmentVar found in flag " + f.LongName + " " + f.ShortName + ". Type not supported: " + reflect.TypeOf(f.AssignmentVar).String())
	}
}
