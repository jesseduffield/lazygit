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
	// other maps in the config, and gives each pattern the color that yaml's
	// rules for merge keys (<<) decide on.
	var colors map[string]string
	if err := node.Decode(&colors); err != nil {
		return err
	}

	keys, err := mappingKeys(node)
	if err != nil {
		return err
	}

	patterns := ColorPatterns(lo.Map(keys, func(pattern string, _ int) ColorPattern {
		return ColorPattern{Pattern: pattern, Color: colors[pattern]}
	}))
	*p = patterns.over(*p)
	return nil
}

// mappingKeys returns the keys of a mapping that node is or refers to, in the
// order in which they are written. A merge key (<<) is replaced by the keys
// that it merges, except those that the mapping has itself, which keep their
// own place. The node must have been decoded successfully, which rules out
// malformed merge keys and aliases that contain themselves.
func mappingKeys(node *yaml.Node) ([]string, error) {
	if node.Kind == yaml.AliasNode {
		return mappingKeys(node.Alias)
	}

	var ownKeys []string
	mergePosition := 0
	var mergeValue *yaml.Node
	for i := 0; i < len(node.Content)-1; i += 2 {
		if isMergeKey(node.Content[i]) {
			// A mapping has at most one merge key, because yaml rejects
			// duplicate keys
			mergePosition, mergeValue = len(ownKeys), node.Content[i+1]
			continue
		}

		var key string
		if err := node.Content[i].Decode(&key); err != nil {
			return nil, err
		}
		ownKeys = append(ownKeys, key)
	}

	if mergeValue == nil {
		return ownKeys, nil
	}

	mergedKeys, err := mergedMappingKeys(mergeValue)
	if err != nil {
		return nil, err
	}
	return slices.Concat(ownKeys[:mergePosition], lo.Without(mergedKeys, ownKeys...), ownKeys[mergePosition:]), nil
}

// mergedMappingKeys returns the keys that the value of a merge key brings in:
// that of a mapping, or of each mapping in a list. A key of several mappings
// in a list takes its place from the first of them, which also gives it its
// value.
func mergedMappingKeys(mergeValue *yaml.Node) ([]string, error) {
	mappings := []*yaml.Node{mergeValue}
	if mergeValue.Kind == yaml.SequenceNode {
		mappings = mergeValue.Content
	}

	var keys []string
	for _, mapping := range mappings {
		keysOfMapping, err := mappingKeys(mapping)
		if err != nil {
			return nil, err
		}
		keys = append(keys, keysOfMapping...)
	}
	return lo.Uniq(keys), nil
}

// isMergeKey reports whether node is a merge key (<<), by the same rule as
// yaml's decoder.
func isMergeKey(node *yaml.Node) bool {
	return node.Kind == yaml.ScalarNode && node.Value == "<<" &&
		(node.Tag == "" || node.Tag == "!" || node.ShortTag() == "!!merge")
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
