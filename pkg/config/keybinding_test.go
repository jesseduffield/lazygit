package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestKeybindingUnmarshalYAML(t *testing.T) {
	scenarios := []struct {
		name     string
		input    string
		expected Keybinding
		wantErr  bool
	}{
		{
			name:     "scalar string",
			input:    `q`,
			expected: Keybinding{"q"},
		},
		{
			name:     "scalar with special characters",
			input:    `<esc>`,
			expected: Keybinding{"<esc>"},
		},
		{
			name:     "sequence with one element",
			input:    `[q]`,
			expected: Keybinding{"q"},
		},
		{
			name:     "sequence with multiple elements",
			input:    `["q", "<esc>"]`,
			expected: Keybinding{"q", "<esc>"},
		},
		{
			name:     "empty sequence",
			input:    `[]`,
			expected: Keybinding{},
		},
		{
			name:     "scalar <disabled> decodes to empty",
			input:    `<disabled>`,
			expected: Keybinding{},
		},
		{
			name:     "scalar empty string decodes to empty",
			input:    `""`,
			expected: Keybinding{},
		},
		{
			name:     "<disabled> entries are filtered out of a sequence",
			input:    `["q", "<disabled>", "<esc>"]`,
			expected: Keybinding{"q", "<esc>"},
		},
		{
			name:    "mapping is rejected",
			input:   `{key: q}`,
			wantErr: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			var k Keybinding
			err := yaml.Unmarshal([]byte(s.input), &k)
			if s.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, s.expected, k)
		})
	}
}

func TestKeybindingMarshalYAML(t *testing.T) {
	scenarios := []struct {
		name     string
		input    Keybinding
		expected string
	}{
		{
			name:     "single key emits a scalar",
			input:    Keybinding{"q"},
			expected: "q\n",
		},
		{
			name:     "multiple keys emit a flow sequence",
			input:    Keybinding{"q", "<esc>"},
			expected: "[q, <esc>]\n",
		},
		{
			name:     "empty keybinding emits an empty sequence",
			input:    Keybinding{},
			expected: "[]\n",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			out, err := yaml.Marshal(s.input)
			assert.NoError(t, err)
			assert.Equal(t, s.expected, string(out))
		})
	}
}

func TestKeybindingMarshalJSON(t *testing.T) {
	scenarios := []struct {
		name     string
		input    Keybinding
		expected string
	}{
		{
			name:     "single key emits a string",
			input:    Keybinding{"q"},
			expected: `"q"`,
		},
		{
			name:     "multiple keys emit an array",
			input:    Keybinding{"q", "esc"},
			expected: `["q","esc"]`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			out, err := json.Marshal(s.input)
			assert.NoError(t, err)
			assert.Equal(t, s.expected, string(out))
		})
	}
}

func TestMergeLegacyAltKeybindings(t *testing.T) {
	scenarios := []struct {
		name     string
		quit     Keybinding
		quitAlt1 Keybinding
		expected Keybinding
	}{
		{
			name:     "alt is folded into main",
			quit:     Keybinding{"q"},
			quitAlt1: Keybinding{"<ctrl+c>"},
			expected: Keybinding{"q", "<ctrl+c>"},
		},
		{
			name:     "alt is not appended if already present",
			quit:     Keybinding{"q", "<ctrl+c>"},
			quitAlt1: Keybinding{"<ctrl+c>"},
			expected: Keybinding{"q", "<ctrl+c>"},
		},
		{
			name:     "empty alt is ignored",
			quit:     Keybinding{"q"},
			quitAlt1: nil,
			expected: Keybinding{"q"},
		},
		{
			name:     "user-supplied multi-key main is preserved",
			quit:     Keybinding{"q", "<esc>"},
			quitAlt1: Keybinding{"<ctrl+c>"},
			expected: Keybinding{"q", "<esc>", "<ctrl+c>"},
		},
		{
			name:     "multi-key alt is folded element by element",
			quit:     Keybinding{"q"},
			quitAlt1: Keybinding{"<ctrl+c>", "<esc>"},
			expected: Keybinding{"q", "<ctrl+c>", "<esc>"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			cfg := KeybindingConfig{
				Universal: KeybindingUniversalConfig{
					Quit:     s.quit,
					QuitAlt1: s.quitAlt1,
				},
			}
			cfg.MergeLegacyAltKeybindings()
			assert.Equal(t, s.expected, cfg.Universal.Quit)
		})
	}
}

func TestKeybindingYAMLRoundTrip(t *testing.T) {
	scenarios := []Keybinding{
		{"q"},
		{"q", "<esc>"},
		{"<ctrl+c>", "<ctrl+d>", "<esc>"},
	}
	for _, original := range scenarios {
		out, err := yaml.Marshal(original)
		assert.NoError(t, err)
		var decoded Keybinding
		assert.NoError(t, yaml.Unmarshal(out, &decoded))
		assert.Equal(t, original, decoded)
	}
}

func TestKeybindingConfigYAMLAcceptsBothForms(t *testing.T) {
	scenarios := []struct {
		name     string
		yaml     string
		expected Keybinding
	}{
		{
			name:     "scalar form",
			yaml:     "quit: q\n",
			expected: Keybinding{"q"},
		},
		{
			name:     "sequence form",
			yaml:     "quit: [q, <esc>]\n",
			expected: Keybinding{"q", "<esc>"},
		},
		{
			name:     "block sequence form",
			yaml:     "quit:\n  - q\n  - <esc>\n",
			expected: Keybinding{"q", "<esc>"},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			var cfg KeybindingUniversalConfig
			assert.NoError(t, yaml.Unmarshal([]byte(s.yaml), &cfg))
			assert.Equal(t, s.expected, cfg.Quit)
		})
	}
}

func TestJumpToBlockYAMLAcceptsMixedForms(t *testing.T) {
	yamlInput := `
jumpToBlock:
  - "1"
  - ["2", "@"]
  - "3"
  - "4"
  - "5"
`
	var cfg KeybindingUniversalConfig
	assert.NoError(t, yaml.Unmarshal([]byte(yamlInput), &cfg))
	expected := []Keybinding{{"1"}, {"2", "@"}, {"3"}, {"4"}, {"5"}}
	assert.Equal(t, expected, cfg.JumpToBlock)
}
