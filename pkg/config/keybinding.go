package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/karimkhaleel/jsonschema"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// Keybinding represents the value of a single keybinding entry in the user's
// config. It's a slice of key strings to allow alternates, but for backward
// compatibility (and because most bindings only have one key) it can be
// written in YAML/JSON as either a single scalar string or as a sequence of
// strings.
type Keybinding []string

func (k *Keybinding) UnmarshalYAML(node *yaml.Node) error {
	var ss []string
	switch node.Kind {
	case yaml.ScalarNode:
		var s string
		if err := node.Decode(&s); err != nil {
			return err
		}
		ss = []string{s}
	case yaml.SequenceNode:
		if err := node.Decode(&ss); err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected a string or a sequence of strings for keybinding, got %v", node.Tag)
	}
	// Drop empty and <disabled> entries so clients never have to special-case
	// them: an empty Keybinding means "no key bound", a non-empty one is
	// guaranteed to contain only real keys.
	*k = lo.Filter(ss, func(s string, _ int) bool {
		return s != "" && s != "<disabled>"
	})
	return nil
}

func (k Keybinding) MarshalYAML() (any, error) {
	if len(k) == 1 {
		return k[0], nil
	}
	// Render multi-key bindings in flow style (`[a, b]`) rather than the default
	// block style, which is more compact and reads better in the generated docs.
	node := &yaml.Node{
		Kind:  yaml.SequenceNode,
		Style: yaml.FlowStyle,
	}
	for _, s := range k {
		node.Content = append(node.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: s,
		})
	}
	return node, nil
}

func (k Keybinding) MarshalJSON() ([]byte, error) {
	if len(k) == 1 {
		return json.Marshal(k[0])
	}
	return json.Marshal([]string(k))
}

// String renders the keybinding as a human-readable label, joining
// alternates with " or " for use in help text.
func (k Keybinding) String() string {
	return strings.Join(k, " or ")
}

// JSONSchema lets the schema generator describe this type as a union of a
// string and an array of strings instead of just an array.
func (Keybinding) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{Type: "string"},
			{Type: "array", Items: &jsonschema.Schema{Type: "string"}},
		},
	}
}

// mergeLegacyAlt folds a deprecated `*Alt*` field into the corresponding
// multi-key main field.
func mergeLegacyAlt(main *Keybinding, alt Keybinding) {
	*main = lo.Union(*main, alt)
}
