package plural

import (
	"fmt"
	"strconv"
	"strings"
)

// Operands is a representation of http://unicode.org/reports/tr35/tr35-numbers.html#Operands
type Operands struct {
	N float64 // absolute value of the source number (integer and decimals)
	I int64   // integer digits of n
	V int64   // number of visible fraction digits in n, with trailing zeros
	W int64   // number of visible fraction digits in n, without trailing zeros
	F int64   // visible fractional digits in n, with trailing zeros
	T int64   // visible fractional digits in n, without trailing zeros
}

// NEqualsAny returns true if o represents an integer equal to any of the arguments.
func (o *Operands) NEqualsAny(any ...int64) bool {
	for _, i := range any {
		if o.I == i && o.T == 0 {
			return true
		}
	}
	return false
}

// NModEqualsAny returns true if o represents an integer equal to any of the arguments modulo mod.
func (o *Operands) NModEqualsAny(mod int64, any ...int64) bool {
	modI := o.I % mod
	for _, i := range any {
		if modI == i && o.T == 0 {
			return true
		}
	}
	return false
}

// NInRange returns true if o represents an integer in the closed interval [from, to].
func (o *Operands) NInRange(from, to int64) bool {
	return o.T == 0 && from <= o.I && o.I <= to
}

// NModInRange returns true if o represents an integer in the closed interval [from, to] modulo mod.
func (o *Operands) NModInRange(mod, from, to int64) bool {
	modI := o.I % mod
	return o.T == 0 && from <= modI && modI <= to
}

// NewOperands returns the operands for number.
func NewOperands(number interface{}) (*Operands, error) {
	switch number := number.(type) {
	case int:
		return newOperandsInt64(int64(number)), nil
	case int8:
		return newOperandsInt64(int64(number)), nil
	case int16:
		return newOperandsInt64(int64(number)), nil
	case int32:
		return newOperandsInt64(int64(number)), nil
	case int64:
		return newOperandsInt64(number), nil
	case string:
		return newOperandsString(number)
	case float32, float64:
		return nil, fmt.Errorf("floats should be formatted into a string")
	default:
		return nil, fmt.Errorf("invalid type %T; expected integer or string", number)
	}
}

func newOperandsInt64(i int64) *Operands {
	if i < 0 {
		i = -i
	}
	return &Operands{float64(i), i, 0, 0, 0, 0}
}

func newOperandsString(s string) (*Operands, error) {
	if s[0] == '-' {
		s = s[1:]
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	ops := &Operands{N: n}
	parts := strings.SplitN(s, ".", 2)
	ops.I, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	if len(parts) == 1 {
		return ops, nil
	}
	fraction := parts[1]
	ops.V = int64(len(fraction))
	for i := ops.V - 1; i >= 0; i-- {
		if fraction[i] != '0' {
			ops.W = i + 1
			break
		}
	}
	if ops.V > 0 {
		f, err := strconv.ParseInt(fraction, 10, 0)
		if err != nil {
			return nil, err
		}
		ops.F = f
	}
	if ops.W > 0 {
		t, err := strconv.ParseInt(fraction[:ops.W], 10, 0)
		if err != nil {
			return nil, err
		}
		ops.T = t
	}
	return ops, nil
}
