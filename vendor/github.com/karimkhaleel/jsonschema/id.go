package jsonschema

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ID represents a Schema ID type which should always be a URI.
// See draft-bhutton-json-schema-00 section 8.2.1
type ID string

// EmptyID is used to explicitly define an ID with no value.
const EmptyID ID = ""

// Validate is used to check if the ID looks like a proper schema.
// This is done by parsing the ID as a URL and checking it has all the
// relevant parts.
func (id ID) Validate() error {
	u, err := url.Parse(id.String())
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Hostname() == "" {
		return errors.New("missing hostname")
	}
	if !strings.Contains(u.Hostname(), ".") {
		return errors.New("hostname does not look valid")
	}
	if u.Path == "" {
		return errors.New("path is expected")
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return errors.New("unexpected schema")
	}
	return nil
}

// Anchor sets the anchor part of the schema URI.
func (id ID) Anchor(name string) ID {
	b := id.Base()
	return ID(b.String() + "#" + name)
}

// Def adds or replaces a definition identifier.
func (id ID) Def(name string) ID {
	b := id.Base()
	return ID(b.String() + "#/$defs/" + name)
}

// Add appends the provided path to the id, and removes any
// anchor data that might be there.
func (id ID) Add(path string) ID {
	b := id.Base()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return ID(b.String() + path)
}

// Base removes any anchor information from the schema
func (id ID) Base() ID {
	s := id.String()
	i := strings.LastIndex(s, "#")
	if i != -1 {
		s = s[0:i]
	}
	s = strings.TrimRight(s, "/")
	return ID(s)
}

// String provides string version of ID
func (id ID) String() string {
	return string(id)
}
