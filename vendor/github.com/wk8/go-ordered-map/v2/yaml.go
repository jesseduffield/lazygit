package orderedmap

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	_ yaml.Marshaler   = &OrderedMap[int, any]{}
	_ yaml.Unmarshaler = &OrderedMap[int, any]{}
)

// MarshalYAML implements the yaml.Marshaler interface.
func (om *OrderedMap[K, V]) MarshalYAML() (interface{}, error) {
	if om == nil {
		return []byte("null"), nil
	}

	node := yaml.Node{
		Kind: yaml.MappingNode,
	}

	for pair := om.Oldest(); pair != nil; pair = pair.Next() {
		key, value := pair.Key, pair.Value

		keyNode := &yaml.Node{}

		// serialize key to yaml, then deserialize it back into the node
		// this is a hack to get the correct tag for the key
		if err := keyNode.Encode(key); err != nil {
			return nil, err
		}

		valueNode := &yaml.Node{}
		if err := valueNode.Encode(value); err != nil {
			return nil, err
		}

		node.Content = append(node.Content, keyNode, valueNode)
	}

	return &node, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (om *OrderedMap[K, V]) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("pipeline must contain YAML mapping, has %v", value.Kind)
	}

	if om.list == nil {
		om.initialize(0)
	}

	for index := 0; index < len(value.Content); index += 2 {
		var key K
		var val V

		if err := value.Content[index].Decode(&key); err != nil {
			return err
		}
		if err := value.Content[index+1].Decode(&val); err != nil {
			return err
		}

		om.Set(key, val)
	}

	return nil
}
