package direnv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDirenvExport(t *testing.T) {
	hello := "hello"
	empty := ""

	scenarios := []struct {
		name    string
		input   string
		want    map[string]*string
		wantErr bool
	}{
		{name: "empty stdout means no .envrc was loaded", input: "", want: nil},
		{name: "literal null from direnv means no delta", input: "null", want: nil},
		{name: "empty object means no delta", input: "{}", want: map[string]*string{}},
		{name: "string value is a set", input: `{"FOO":"hello"}`, want: map[string]*string{"FOO": &hello}},
		{name: "null value is an unset", input: `{"FOO":null}`, want: map[string]*string{"FOO": nil}},
		{
			name:  "set and unset can coexist",
			input: `{"FOO":"hello","BAR":null,"BAZ":""}`,
			want:  map[string]*string{"FOO": &hello, "BAR": nil, "BAZ": &empty},
		},
		{name: "malformed JSON is an error", input: `{not json`, wantErr: true},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			got, err := parseDirenvExport([]byte(s.input))
			if s.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, s.want, got)
			}
		})
	}
}
