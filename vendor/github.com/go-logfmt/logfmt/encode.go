package logfmt

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode/utf8"
)

// MarshalKeyvals returns the logfmt encoding of keyvals, a variadic sequence
// of alternating keys and values.
func MarshalKeyvals(keyvals ...interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := NewEncoder(buf).EncodeKeyvals(keyvals...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// An Encoder writes logfmt data to an output stream.
type Encoder struct {
	w       io.Writer
	scratch bytes.Buffer
	needSep bool
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

var (
	space   = []byte(" ")
	equals  = []byte("=")
	newline = []byte("\n")
	null    = []byte("null")
)

// EncodeKeyval writes the logfmt encoding of key and value to the stream. A
// single space is written before the second and subsequent keys in a record.
// Nothing is written if a non-nil error is returned.
func (enc *Encoder) EncodeKeyval(key, value interface{}) error {
	enc.scratch.Reset()
	if enc.needSep {
		if _, err := enc.scratch.Write(space); err != nil {
			return err
		}
	}
	if err := writeKey(&enc.scratch, key); err != nil {
		return err
	}
	if _, err := enc.scratch.Write(equals); err != nil {
		return err
	}
	if err := writeValue(&enc.scratch, value); err != nil {
		return err
	}
	_, err := enc.w.Write(enc.scratch.Bytes())
	enc.needSep = true
	return err
}

// EncodeKeyvals writes the logfmt encoding of keyvals to the stream. Keyvals
// is a variadic sequence of alternating keys and values. Keys of unsupported
// type are skipped along with their corresponding value. Values of
// unsupported type or that cause a MarshalerError are replaced by their error
// but do not cause EncodeKeyvals to return an error. If a non-nil error is
// returned some key/value pairs may not have be written.
func (enc *Encoder) EncodeKeyvals(keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 == 1 {
		keyvals = append(keyvals, nil)
	}
	for i := 0; i < len(keyvals); i += 2 {
		k, v := keyvals[i], keyvals[i+1]
		err := enc.EncodeKeyval(k, v)
		if err == ErrUnsupportedKeyType {
			continue
		}
		if _, ok := err.(*MarshalerError); ok || err == ErrUnsupportedValueType {
			v = err
			err = enc.EncodeKeyval(k, v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalerError represents an error encountered while marshaling a value.
type MarshalerError struct {
	Type reflect.Type
	Err  error
}

func (e *MarshalerError) Error() string {
	return "error marshaling value of type " + e.Type.String() + ": " + e.Err.Error()
}

// ErrNilKey is returned by Marshal functions and Encoder methods if a key is
// a nil interface or pointer value.
var ErrNilKey = errors.New("nil key")

// ErrInvalidKey is returned by Marshal functions and Encoder methods if, after
// dropping invalid runes, a key is empty.
var ErrInvalidKey = errors.New("invalid key")

// ErrUnsupportedKeyType is returned by Encoder methods if a key has an
// unsupported type.
var ErrUnsupportedKeyType = errors.New("unsupported key type")

// ErrUnsupportedValueType is returned by Encoder methods if a value has an
// unsupported type.
var ErrUnsupportedValueType = errors.New("unsupported value type")

func writeKey(w io.Writer, key interface{}) error {
	if key == nil {
		return ErrNilKey
	}

	switch k := key.(type) {
	case string:
		return writeStringKey(w, k)
	case []byte:
		if k == nil {
			return ErrNilKey
		}
		return writeBytesKey(w, k)
	case encoding.TextMarshaler:
		kb, err := safeMarshal(k)
		if err != nil {
			return err
		}
		if kb == nil {
			return ErrNilKey
		}
		return writeBytesKey(w, kb)
	case fmt.Stringer:
		ks, ok := safeString(k)
		if !ok {
			return ErrNilKey
		}
		return writeStringKey(w, ks)
	default:
		rkey := reflect.ValueOf(key)
		switch rkey.Kind() {
		case reflect.Array, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Struct:
			return ErrUnsupportedKeyType
		case reflect.Ptr:
			if rkey.IsNil() {
				return ErrNilKey
			}
			return writeKey(w, rkey.Elem().Interface())
		}
		return writeStringKey(w, fmt.Sprint(k))
	}
}

// keyRuneFilter returns r for all valid key runes, and -1 for all invalid key
// runes. When used as the mapping function for strings.Map and bytes.Map
// functions it causes them to remove invalid key runes from strings or byte
// slices respectively.
func keyRuneFilter(r rune) rune {
	if r <= ' ' || r == '=' || r == '"' || r == utf8.RuneError {
		return -1
	}
	return r
}

func writeStringKey(w io.Writer, key string) error {
	k := strings.Map(keyRuneFilter, key)
	if k == "" {
		return ErrInvalidKey
	}
	_, err := io.WriteString(w, k)
	return err
}

func writeBytesKey(w io.Writer, key []byte) error {
	k := bytes.Map(keyRuneFilter, key)
	if len(k) == 0 {
		return ErrInvalidKey
	}
	_, err := w.Write(k)
	return err
}

func writeValue(w io.Writer, value interface{}) error {
	switch v := value.(type) {
	case nil:
		return writeBytesValue(w, null)
	case string:
		return writeStringValue(w, v, true)
	case []byte:
		return writeBytesValue(w, v)
	case encoding.TextMarshaler:
		vb, err := safeMarshal(v)
		if err != nil {
			return err
		}
		if vb == nil {
			vb = null
		}
		return writeBytesValue(w, vb)
	case error:
		se, ok := safeError(v)
		return writeStringValue(w, se, ok)
	case fmt.Stringer:
		ss, ok := safeString(v)
		return writeStringValue(w, ss, ok)
	default:
		rvalue := reflect.ValueOf(value)
		switch rvalue.Kind() {
		case reflect.Array, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Struct:
			return ErrUnsupportedValueType
		case reflect.Ptr:
			if rvalue.IsNil() {
				return writeBytesValue(w, null)
			}
			return writeValue(w, rvalue.Elem().Interface())
		}
		return writeStringValue(w, fmt.Sprint(v), true)
	}
}

func needsQuotedValueRune(r rune) bool {
	return r <= ' ' || r == '=' || r == '"' || r == utf8.RuneError
}

func writeStringValue(w io.Writer, value string, ok bool) error {
	var err error
	if ok && value == "null" {
		_, err = io.WriteString(w, `"null"`)
	} else if strings.IndexFunc(value, needsQuotedValueRune) != -1 {
		_, err = writeQuotedString(w, value)
	} else {
		_, err = io.WriteString(w, value)
	}
	return err
}

func writeBytesValue(w io.Writer, value []byte) error {
	var err error
	if bytes.IndexFunc(value, needsQuotedValueRune) != -1 {
		_, err = writeQuotedBytes(w, value)
	} else {
		_, err = w.Write(value)
	}
	return err
}

// EndRecord writes a newline character to the stream and resets the encoder
// to the beginning of a new record.
func (enc *Encoder) EndRecord() error {
	_, err := enc.w.Write(newline)
	if err == nil {
		enc.needSep = false
	}
	return err
}

// Reset resets the encoder to the beginning of a new record.
func (enc *Encoder) Reset() {
	enc.needSep = false
}

func safeError(err error) (s string, ok bool) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if v := reflect.ValueOf(err); v.Kind() == reflect.Ptr && v.IsNil() {
				s, ok = "null", false
			} else {
				s, ok = fmt.Sprintf("PANIC:%v", panicVal), false
			}
		}
	}()
	s, ok = err.Error(), true
	return
}

func safeString(str fmt.Stringer) (s string, ok bool) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if v := reflect.ValueOf(str); v.Kind() == reflect.Ptr && v.IsNil() {
				s, ok = "null", false
			} else {
				s, ok = fmt.Sprintf("PANIC:%v", panicVal), true
			}
		}
	}()
	s, ok = str.String(), true
	return
}

func safeMarshal(tm encoding.TextMarshaler) (b []byte, err error) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if v := reflect.ValueOf(tm); v.Kind() == reflect.Ptr && v.IsNil() {
				b, err = nil, nil
			} else {
				b, err = nil, fmt.Errorf("panic when marshalling: %s", panicVal)
			}
		}
	}()
	b, err = tm.MarshalText()
	if err != nil {
		return nil, &MarshalerError{
			Type: reflect.TypeOf(tm),
			Err:  err,
		}
	}
	return
}
