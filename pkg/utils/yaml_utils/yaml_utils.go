package yaml_utils

import (
	"errors"
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

	if body.Kind != yaml.MappingNode {
		return yamlBytes, errors.New("yaml document is not a dictionary")
	}

	if didChange, err := updateYamlNode(body, path, value); err != nil || !didChange {
		return yamlBytes, err
	}

	// Convert the updated YAML node back to YAML bytes.
	updatedYAMLBytes, err := yaml.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML node to bytes: %w", err)
	}

	return updatedYAMLBytes, nil
}

// Recursive function to update the YAML node.
func updateYamlNode(node *yaml.Node, path []string, value string) (bool, error) {
	if len(path) == 0 {
		if node.Kind != yaml.ScalarNode {
			return false, errors.New("yaml node is not a scalar")
		}
		if node.Value != value {
			node.Value = value
			return true, nil
		}
		return false, nil
	}

	if node.Kind != yaml.MappingNode {
		return false, errors.New("yaml node in path is not a dictionary")
	}

	key := path[0]
	if _, valueNode := lookupKey(node, key); valueNode != nil {
		return updateYamlNode(valueNode, path[1:], value)
	}

	// if the key doesn't exist, we'll add it

	// at end of path: add the new key, done
	if len(path) == 1 {
		node.Content = append(node.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: key,
		}, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: value,
		})
		return true, nil
	}

	// otherwise, create the missing intermediate node and continue
	newNode := &yaml.Node{
		Kind: yaml.MappingNode,
	}
	node.Content = append(node.Content, &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: key,
	}, newNode)
	return updateYamlNode(newNode, path[1:], value)
}

func lookupKey(node *yaml.Node, key string) (*yaml.Node, *yaml.Node) {
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i], node.Content[i+1]
		}
	}

	return nil, nil
}

// takes a yaml document in bytes, a path to a key, and a new name for the key.
// Will rename the key to the new name if it exists, and do nothing otherwise.
func RenameYamlKey(yamlBytes []byte, path []string, newKey string) ([]byte, error) {
	// Parse the YAML file.
	var node yaml.Node
	err := yaml.Unmarshal(yamlBytes, &node)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Empty document: nothing to do.
	if len(node.Content) == 0 {
		return yamlBytes, nil
	}

	body := node.Content[0]

	if didRename, err := renameYamlKey(body, path, newKey); err != nil || !didRename {
		return yamlBytes, err
	}

	// Convert the updated YAML node back to YAML bytes.
	updatedYAMLBytes, err := yaml.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML node to bytes: %w", err)
	}

	return updatedYAMLBytes, nil
}

// Recursive function to rename the YAML key.
func renameYamlKey(node *yaml.Node, path []string, newKey string) (bool, error) {
	if node.Kind != yaml.MappingNode {
		return false, errors.New("yaml node in path is not a dictionary")
	}

	keyNode, valueNode := lookupKey(node, path[0])
	if keyNode == nil {
		return false, nil
	}

	// end of path reached: rename key
	if len(path) == 1 {
		// Check that new key doesn't exist yet
		if newKeyNode, _ := lookupKey(node, newKey); newKeyNode != nil {
			return false, fmt.Errorf("new key `%s' already exists", newKey)
		}

		keyNode.Value = newKey
		return true, nil
	}

	return renameYamlKey(valueNode, path[1:], newKey)
}
