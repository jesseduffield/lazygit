package parser

import (
	"bytes"
	"io"
	"strconv"

	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var attrNameID = []byte("id")
var attrNameClass = []byte("class")

// An Attribute is an attribute of the markdown elements
type Attribute struct {
	Name  []byte
	Value interface{}
}

// An Attributes is a collection of attributes.
type Attributes []Attribute

// Find returns a (value, true) if an attribute correspond with given name is found, otherwise (nil, false).
func (as Attributes) Find(name []byte) (interface{}, bool) {
	for _, a := range as {
		if bytes.Equal(a.Name, name) {
			return a.Value, true
		}
	}
	return nil, false
}

func (as Attributes) findUpdate(name []byte, cb func(v interface{}) interface{}) bool {
	for i, a := range as {
		if bytes.Equal(a.Name, name) {
			as[i].Value = cb(a.Value)
			return true
		}
	}
	return false
}

// ParseAttributes parses attributes into a map.
// ParseAttributes returns a parsed attributes and true if could parse
// attributes, otherwise nil and false.
func ParseAttributes(reader text.Reader) (Attributes, bool) {
	savedLine, savedPosition := reader.Position()
	reader.SkipSpaces()
	if reader.Peek() != '{' {
		reader.SetPosition(savedLine, savedPosition)
		return nil, false
	}
	reader.Advance(1)
	attrs := Attributes{}
	for {
		if reader.Peek() == '}' {
			reader.Advance(1)
			return attrs, true
		}
		attr, ok := parseAttribute(reader)
		if !ok {
			reader.SetPosition(savedLine, savedPosition)
			return nil, false
		}
		if bytes.Equal(attr.Name, attrNameClass) {
			if !attrs.findUpdate(attrNameClass, func(v interface{}) interface{} {
				ret := make([]byte, 0, len(v.([]byte))+1+len(attr.Value.([]byte)))
				ret = append(ret, v.([]byte)...)
				return append(append(ret, ' '), attr.Value.([]byte)...)
			}) {
				attrs = append(attrs, attr)
			}
		} else {
			attrs = append(attrs, attr)
		}
		reader.SkipSpaces()
		if reader.Peek() == ',' {
			reader.Advance(1)
			reader.SkipSpaces()
		}
	}
}

func parseAttribute(reader text.Reader) (Attribute, bool) {
	reader.SkipSpaces()
	c := reader.Peek()
	if c == '#' || c == '.' {
		reader.Advance(1)
		line, _ := reader.PeekLine()
		i := 0
		// HTML5 allows any kind of characters as id, but XHTML restricts characters for id.
		// CommonMark is basically defined for XHTML(even though it is legacy).
		// So we restrict id characters.
		for ; i < len(line) && !util.IsSpace(line[i]) &&
			(!util.IsPunct(line[i]) || line[i] == '_' || line[i] == '-' || line[i] == ':' || line[i] == '.'); i++ {
		}
		name := attrNameClass
		if c == '#' {
			name = attrNameID
		}
		reader.Advance(i)
		return Attribute{Name: name, Value: line[0:i]}, true
	}
	line, _ := reader.PeekLine()
	if len(line) == 0 {
		return Attribute{}, false
	}
	c = line[0]
	if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		c == '_' || c == ':') {
		return Attribute{}, false
	}
	i := 0
	for ; i < len(line); i++ {
		c = line[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_' || c == ':' || c == '.' || c == '-') {
			break
		}
	}
	name := line[:i]
	reader.Advance(i)
	reader.SkipSpaces()
	c = reader.Peek()
	if c != '=' {
		return Attribute{}, false
	}
	reader.Advance(1)
	reader.SkipSpaces()
	value, ok := parseAttributeValue(reader)
	if !ok {
		return Attribute{}, false
	}
	if bytes.Equal(name, attrNameClass) {
		if _, ok = value.([]byte); !ok {
			return Attribute{}, false
		}
	}
	return Attribute{Name: name, Value: value}, true
}

func parseAttributeValue(reader text.Reader) (interface{}, bool) {
	reader.SkipSpaces()
	c := reader.Peek()
	var value interface{}
	ok := false
	switch c {
	case text.EOF:
		return Attribute{}, false
	case '{':
		value, ok = ParseAttributes(reader)
	case '[':
		value, ok = parseAttributeArray(reader)
	case '"':
		value, ok = parseAttributeString(reader)
	default:
		if c == '-' || c == '+' || util.IsNumeric(c) {
			value, ok = parseAttributeNumber(reader)
		} else {
			value, ok = parseAttributeOthers(reader)
		}
	}
	if !ok {
		return nil, false
	}
	return value, true
}

func parseAttributeArray(reader text.Reader) ([]interface{}, bool) {
	reader.Advance(1) // skip [
	ret := []interface{}{}
	for i := 0; ; i++ {
		c := reader.Peek()
		comma := false
		if i != 0 && c == ',' {
			reader.Advance(1)
			comma = true
		}
		if c == ']' {
			if !comma {
				reader.Advance(1)
				return ret, true
			}
			return nil, false
		}
		reader.SkipSpaces()
		value, ok := parseAttributeValue(reader)
		if !ok {
			return nil, false
		}
		ret = append(ret, value)
		reader.SkipSpaces()
	}
}

func parseAttributeString(reader text.Reader) ([]byte, bool) {
	reader.Advance(1) // skip "
	line, _ := reader.PeekLine()
	i := 0
	l := len(line)
	var buf bytes.Buffer
	for i < l {
		c := line[i]
		if c == '\\' && i != l-1 {
			n := line[i+1]
			switch n {
			case '"', '/', '\\':
				buf.WriteByte(n)
				i += 2
			case 'b':
				buf.WriteString("\b")
				i += 2
			case 'f':
				buf.WriteString("\f")
				i += 2
			case 'n':
				buf.WriteString("\n")
				i += 2
			case 'r':
				buf.WriteString("\r")
				i += 2
			case 't':
				buf.WriteString("\t")
				i += 2
			default:
				buf.WriteByte('\\')
				i++
			}
			continue
		}
		if c == '"' {
			reader.Advance(i + 1)
			return buf.Bytes(), true
		}
		buf.WriteByte(c)
		i++
	}
	return nil, false
}

func scanAttributeDecimal(reader text.Reader, w io.ByteWriter) {
	for {
		c := reader.Peek()
		if util.IsNumeric(c) {
			w.WriteByte(c)
		} else {
			return
		}
		reader.Advance(1)
	}
}

func parseAttributeNumber(reader text.Reader) (float64, bool) {
	sign := 1
	c := reader.Peek()
	if c == '-' {
		sign = -1
		reader.Advance(1)
	} else if c == '+' {
		reader.Advance(1)
	}
	var buf bytes.Buffer
	if !util.IsNumeric(reader.Peek()) {
		return 0, false
	}
	scanAttributeDecimal(reader, &buf)
	if buf.Len() == 0 {
		return 0, false
	}
	c = reader.Peek()
	if c == '.' {
		buf.WriteByte(c)
		reader.Advance(1)
		scanAttributeDecimal(reader, &buf)
	}
	c = reader.Peek()
	if c == 'e' || c == 'E' {
		buf.WriteByte(c)
		reader.Advance(1)
		c = reader.Peek()
		if c == '-' || c == '+' {
			buf.WriteByte(c)
			reader.Advance(1)
		}
		scanAttributeDecimal(reader, &buf)
	}
	f, err := strconv.ParseFloat(buf.String(), 10)
	if err != nil {
		return 0, false
	}
	return float64(sign) * f, true
}

var bytesTrue = []byte("true")
var bytesFalse = []byte("false")
var bytesNull = []byte("null")

func parseAttributeOthers(reader text.Reader) (interface{}, bool) {
	line, _ := reader.PeekLine()
	c := line[0]
	if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		c == '_' || c == ':') {
		return nil, false
	}
	i := 0
	for ; i < len(line); i++ {
		c := line[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_' || c == ':' || c == '.' || c == '-') {
			break
		}
	}
	value := line[:i]
	reader.Advance(i)
	if bytes.Equal(value, bytesTrue) {
		return true, true
	}
	if bytes.Equal(value, bytesFalse) {
		return false, true
	}
	if bytes.Equal(value, bytesNull) {
		return nil, true
	}
	return value, true
}
