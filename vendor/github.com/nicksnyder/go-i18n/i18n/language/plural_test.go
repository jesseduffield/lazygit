package language

import (
	"testing"
)

func TestNewPlural(t *testing.T) {
	tests := []struct {
		src    string
		plural Plural
		err    bool
	}{
		{"zero", Zero, false},
		{"one", One, false},
		{"two", Two, false},
		{"few", Few, false},
		{"many", Many, false},
		{"other", Other, false},
		{"asdf", Invalid, true},
	}
	for _, test := range tests {
		plural, err := NewPlural(test.src)
		wrongErr := (err != nil && !test.err) || (err == nil && test.err)
		if plural != test.plural || wrongErr {
			t.Errorf("NewPlural(%#v) returned %#v,%#v; expected %#v", test.src, plural, err, test.plural)
		}
	}
}
