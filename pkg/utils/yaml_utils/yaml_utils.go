package yaml_utils

import (
	"bytes"
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
	updatedYAMLBytes, err := YamlMarshal(body)
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
	if _, valueNode := LookupKey(node, key); valueNode != nil {
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

func LookupKey(node *yaml.Node, key string) (*yaml.Node, *yaml.Node) {
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i], node.Content[i+1]
		}
	}

	return nil, nil
}

// Returns the key and value if they were present
func RemoveKey(node *yaml.Node, key string) (*yaml.Node, *yaml.Node) {
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			key, value := node.Content[i], node.Content[i+1]
			node.Content = append(node.Content[:i], node.Content[i+2:]...)
			return key, value
		}
	}

	return nil, nil
}

// Walks a yaml document from the root node to the specified path, and then applies the transformation to that node.
// If the requested path is not defined in the document, no changes are made to the document.
func TransformNode(rootNode *yaml.Node, path []string, transform func(node *yaml.Node) error) error {
	// Empty document: nothing to do.
	if len(rootNode.Content) == 0 {
		return nil
	}

	body := rootNode.Content[0]

	if err := transformNode(body, path, transform); err != nil {
		return err
	}

	return nil
}

// A recursive function to walk down the tree. See TransformNode for more details.
func transformNode(node *yaml.Node, path []string, transform func(node *yaml.Node) error) error {
	if len(path) == 0 {
		return transform(node)
	}

	keyNode, valueNode := LookupKey(node, path[0])
	if keyNode == nil {
		return nil
	}

	return transformNode(valueNode, path[1:], transform)
}

// Takes the root node of a yaml document, a path to a key, and a new name for the key.
// Will rename the key to the new name if it exists, and do nothing otherwise.
func RenameYamlKey(rootNode *yaml.Node, path []string, newKey string) error {
	// Empty document: nothing to do.
	if len(rootNode.Content) == 0 {
		return nil
	}

	body := rootNode.Content[0]

	if err := renameYamlKey(body, path, newKey); err != nil {
		return err
	}

	return nil
}

// Recursive function to rename the YAML key.
func renameYamlKey(node *yaml.Node, path []string, newKey string) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("yaml node in path is not a dictionary")
	}

	keyNode, valueNode := LookupKey(node, path[0])
	if keyNode == nil {
		return nil
	}

	// end of path reached: rename key
	if len(path) == 1 {
		// Check that new key doesn't exist yet
		if newKeyNode, _ := LookupKey(node, newKey); newKeyNode != nil {
			return fmt.Errorf("new key `%s' already exists", newKey)
		}

		keyNode.Value = newKey
		return nil
	}

	return renameYamlKey(valueNode, path[1:], newKey)
}

// Traverses a yaml document, calling the callback function for each node. The
// callback is expected to modify the node in place
func Walk(rootNode *yaml.Node, callback func(node *yaml.Node, path string)) error {
	// Empty document: nothing to do.
	if len(rootNode.Content) == 0 {
		return nil
	}

	body := rootNode.Content[0]

	if err := walk(body, "", callback); err != nil {
		return err
	}

	return nil
}

func walk(node *yaml.Node, path string, callback func(*yaml.Node, string)) error {
	callback(node, path)
	switch node.Kind {
	case yaml.DocumentNode:
		return errors.New("Unexpected document node in the middle of a yaml tree")
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			name := node.Content[i].Value
			childNode := node.Content[i+1]
			var childPath string
			if path == "" {
				childPath = name
			} else {
				childPath = fmt.Sprintf("%s.%s", path, name)
			}
			err := walk(childNode, childPath, callback)
			if err != nil {
				return err
			}
		}
	case yaml.SequenceNode:
		for i := 0; i < len(node.Content); i++ {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			err := walk(node.Content[i], childPath, callback)
			if err != nil {
				return err
			}
		}
	case yaml.ScalarNode:
		// nothing to do
	case yaml.AliasNode:
		return errors.New("Alias nodes are not supported")
	}

	return nil
}

func YamlMarshal(node *yaml.Node) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)

	err := encoder.Encode(node)
	return buffer.Bytes(), err
}
