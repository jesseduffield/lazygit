// Package implements the decoding of logfmt key-value pairs.
//
// Example logfmt message:
//
//	foo=bar a=14 baz="hello kitty" cool%story=bro f %^asdf
//
// Example result in JSON:
//
//	{ "foo": "bar", "a": 14, "baz": "hello kitty", "cool%story": "bro", "f": true, "%^asdf": true }
//
// EBNFish:
//
// 	ident_byte = any byte greater than ' ', excluding '=' and '"'
// 	string_byte = any byte excluding '"' and '\'
// 	garbage = !ident_byte
// 	ident = ident_byte, { ident byte }
// 	key = ident
// 	value = ident | '"', { string_byte | '\', '"' }, '"'
// 	pair = key, '=', value | key, '=' | key
// 	message = { garbage, pair }, garbage
package logfmt

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Handler is the interface implemented by objects that accept logfmt
// key-value pairs. HandleLogfmt must copy the logfmt data if it
// wishes to retain the data after returning.
type Handler interface {
	HandleLogfmt(key, val []byte) error
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as
// logfmt handlers. If f is a function with the appropriate signature,
// HandlerFunc(f) is a Handler object that calls f.
type HandlerFunc func(key, val []byte) error

func (f HandlerFunc) HandleLogfmt(key, val []byte) error {
	return f(key, val)
}

// Unmarshal parses the logfmt encoding data and stores the result in the value
// pointed to by v. If v is an Handler, HandleLogfmt will be called for each
// key-value pair.
//
// If v is not a Handler, it will pass v to NewStructHandler and use the
// returned StructHandler for decoding.
func Unmarshal(data []byte, v interface{}) (err error) {
	h, ok := v.(Handler)
	if !ok {
		h, err = NewStructHandler(v)
		if err != nil {
			return err
		}
	}
	return gotoScanner(data, h)
}

// StructHandler unmarshals logfmt into a struct. It matches incoming keys to
// the the struct's fields (either the struct field name or its tag, preferring
// an exact match but also accepting a case-insensitive match.
//
// Field types supported by StructHandler are:
//
// 	all numeric types (e.g. float32, int, etc.)
// 	[]byte
// 	string
// 	bool - true if key is present, false otherwise (the value is ignored).
//	time.Duration - uses time.ParseDuration
//
// If a field is a pointer to an above type, and a matching key is not present
// in the logfmt data, the pointer will be untouched.
//
// If v is not a pointer to an Handler or struct, Unmarshal will return an
// error.
type StructHandler struct {
	rv reflect.Value
}

func NewStructHandler(v interface{}) (Handler, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	return &StructHandler{rv: rv}, nil
}

func (h *StructHandler) HandleLogfmt(key, val []byte) error {
	el := h.rv.Elem()
	skey := string(key)
	for i := 0; i < el.NumField(); i++ {
		fv := el.Field(i)
		ft := el.Type().Field(i)
		switch {
		case ft.Name == skey:
		case ft.Tag.Get("logfmt") == skey:
		case strings.EqualFold(ft.Name, skey):
		default:
			continue
		}
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				t := fv.Type().Elem()
				v := reflect.New(t)
				fv.Set(v)
				fv = v
			}
			fv = fv.Elem()
		}
		switch fv.Interface().(type) {
		case time.Duration:
			d, err := time.ParseDuration(string(val))
			if err != nil {
				return &UnmarshalTypeError{string(val), fv.Type()}
			}
			fv.Set(reflect.ValueOf(d))
		case string:
			fv.SetString(string(val))
		case []byte:
			b := make([]byte, len(val))
			copy(b, val)
			fv.SetBytes(b)
		case bool:
			fv.SetBool(true)
		default:
			switch {
			case reflect.Int <= fv.Kind() && fv.Kind() <= reflect.Int64:
				v, err := strconv.ParseInt(string(val), 10, 64)
				if err != nil {
					return err
				}
				fv.SetInt(v)
			case reflect.Uint32 <= fv.Kind() && fv.Kind() <= reflect.Uint64:
				v, err := strconv.ParseUint(string(val), 10, 64)
				if err != nil {
					return err
				}
				fv.SetUint(v)
			case reflect.Float32 <= fv.Kind() && fv.Kind() <= reflect.Float64:
				v, err := strconv.ParseFloat(string(val), 10)
				if err != nil {
					return err
				}
				fv.SetFloat(v)
			default:
				return &UnmarshalTypeError{string(val), fv.Type()}
			}
		}

	}
	return nil
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "logfmt: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "logfmt: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "logfmt: Unmarshal(nil " + e.Type.String() + ")"
}

// An UnmarshalTypeError describes a logfmt value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value string       // the logfmt value
	Type  reflect.Type // type of Go value it could not be assigned to
}

func (e *UnmarshalTypeError) Error() string {
	return "logfmt: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}
