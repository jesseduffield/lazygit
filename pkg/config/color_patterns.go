package config

import (
	"slices"

	"github.com/karimkhaleel/jsonschema"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// ColorPatterns assigns colors to the names that match regular expressions.
// It's written in YAML as a mapping from pattern to color, and keeps the
// patterns in the order in which they are written, so that the first pattern
// that matches a name can decide its color.
type ColorPatterns []ColorPattern

type ColorPattern struct {
	Pattern string
	Color   string
}

// UnmarshalYAML puts the patterns it reads in front of the ones that are there
// already, which come from config files that were loaded earlier.
func (p *ColorPatterns) UnmarshalYAML(node *yaml.Node) error {
	// Decoding into a map reports malformed input the same way as for the
	// other maps in the config.
	var colors map[string]string
	if err := node.Decode(&colors); err != nil {
		return err
	}

	patterns := make(ColorPatterns, 0, len(colors))
	for i := 0; i < len(node.Content)-1; i += 2 {
		var pattern string
		if err := node.Content[i].Decode(&pattern); err != nil {
			return err
		}
		patterns = append(patterns, ColorPattern{Pattern: pattern, Color: colors[pattern]})
	}

	*p = patterns.over(*p)
	return nil
}

func (p ColorPatterns) MarshalYAML() (any, error) {
	node := &yaml.Node{Kind: yaml.MappingNode}
	for _, pattern := range p {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: pattern.Pattern},
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: pattern.Color},
		)
	}
	return node, nil
}

// JSONSchema describes the patterns as the mapping they are written as.
func (ColorPatterns) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:                 "object",
		AdditionalProperties: &jsonschema.Schema{Type: "string"},
	}
}

// over returns p followed by the patterns of lower that p doesn't have.
func (p ColorPatterns) over(lower ColorPatterns) ColorPatterns {
	return slices.Concat(p, lo.Reject(lower, func(l ColorPattern, _ int) bool {
		return slices.ContainsFunc(p, func(c ColorPattern) bool { return c.Pattern == l.Pattern })
	}))
}
