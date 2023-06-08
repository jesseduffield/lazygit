package yaml_utils

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// takes a yaml document in bytes, a path to a key, and a value to set. The value must be a scalar.
func UpdateYamlValue(yamlBytes []byte, path []string, value string) ([]byte, error) {
	// Parse the YAML file.
	var node yaml.Node
	err := yaml.Unmarshal(yamlBytes, &node)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Empty document: need to create the top-level map ourselves
	if len(node.Content) == 0 {
		node.Content = append(node.Content, &yaml.Node{
			Kind: yaml.MappingNode,
		})
	}

	body := node.Content[0]

	updateYamlNode(body, path, value)

	// Convert the updated YAML node back to YAML bytes.
	updatedYAMLBytes, err := yaml.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML node to bytes: %w", err)
	}

	return updatedYAMLBytes, nil
}

// Recursive function to update the YAML node.
func updateYamlNode(node *yaml.Node, path []string, value string) {
	if len(path) == 0 {
		node.Value = value
		return
	}

	key := path[0]
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			updateYamlNode(node.Content[i+1], path[1:], value)
			return
		}
	}

	// if the key doesn't exist, we'll add it
	node.Content = append(node.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: key,
	}, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: value,
	})
}
