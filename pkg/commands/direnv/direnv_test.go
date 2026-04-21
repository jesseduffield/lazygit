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

func TestParseDirenvStatus(t *testing.T) {
	scenarios := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no .envrc found",
			input: `{"state":{"foundRC":null}}`,
			want:  "",
		},
		{
			name:  "found and allowed (0)",
			input: `{"state":{"foundRC":{"allowed":0,"path":"/repo/.envrc"}}}`,
			want:  "",
		},
		{
			name:  "found but not allowed (1) — eligible for approval",
			input: `{"state":{"foundRC":{"allowed":1,"path":"/repo/.envrc"}}}`,
			want:  "/repo/.envrc",
		},
		{
			name:  "found but denied (2) — user already said no",
			input: `{"state":{"foundRC":{"allowed":2,"path":"/repo/.envrc"}}}`,
			want:  "",
		},
		{
			name:  "malformed JSON",
			input: `{not json`,
			want:  "",
		},
		{
			name:  "empty input",
			input: "",
			want:  "",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.want, parseDirenvStatus([]byte(s.input)))
		})
	}
}
