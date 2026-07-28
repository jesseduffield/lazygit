package yaml_utils

import (
	"bytes"
	"errors"
	"fmt"
	"slices"

	"gopkg.in/yaml.v3"
)

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
			node.Content = slices.Delete(node.Content, i, i+2)
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
func RenameYamlKey(rootNode *yaml.Node, path []string, newKey string) (error, bool) {
	// Empty document: nothing to do.
	if len(rootNode.Content) == 0 {
		return nil, false
	}

	body := rootNode.Content[0]

	return renameYamlKey(body, path, newKey)
}

// Recursive function to rename the YAML key.
func renameYamlKey(node *yaml.Node, path []string, newKey string) (error, bool) {
	if node.Kind != yaml.MappingNode {
		return errors.New("yaml node in path is not a dictionary"), false
	}

	keyNode, valueNode := LookupKey(node, path[0])
	if keyNode == nil {
		return nil, false
	}

	// end of path reached: rename key
	if len(path) == 1 {
		// Check that new key doesn't exist yet
		if newKeyNode, _ := LookupKey(node, newKey); newKeyNode != nil {
			return fmt.Errorf("new key `%s' already exists", newKey), false
		}

		keyNode.Value = newKey
		return nil, true
	}

	return renameYamlKey(valueNode, path[1:], newKey)
}

// Takes the root node of a yaml document, the path to an existing key, and the
// path at which it should live instead. If the key exists, it (and its value)
// is moved to the new path, creating intermediate mapping nodes as needed, and
// any mapping nodes left empty behind it are removed. Does nothing if the key
// at oldPath doesn't exist. Returns an error if a key already exists at newPath,
// or if a node along either path exists but isn't a mapping.
func MoveYamlKey(rootNode *yaml.Node, oldPath []string, newPath []string) (error, bool) {
	// Empty document: nothing to do.
	if len(rootNode.Content) == 0 {
		return nil, false
	}

	body := rootNode.Content[0]

	// Bail out early if there's nothing to move.
	oldParent, err := findContainingMap(body, oldPath, false)
	if err != nil {
		return err, false
	}
	if oldParent == nil {
		return nil, false
	}
	keyNode, valueNode := LookupKey(oldParent, oldPath[len(oldPath)-1])
	if keyNode == nil {
		return nil, false
	}

	// Find or create the destination map, and make sure it's free.
	newParent, err := findContainingMap(body, newPath, true)
	if err != nil {
		return err, false
	}
	newKey := newPath[len(newPath)-1]
	if existing, _ := LookupKey(newParent, newKey); existing != nil {
		return fmt.Errorf("new key `%s' already exists", newKey), false
	}

	// Move the key, then prune any maps that became empty behind it. The
	// destination is populated first so that a map shared by both paths isn't
	// mistaken for empty during pruning.
	RemoveKey(oldParent, oldPath[len(oldPath)-1])
	keyNode.Value = newKey
	newParent.Content = append(newParent.Content, keyNode, valueNode)
	removeEmptyMaps(body, oldPath[:len(oldPath)-1])

	return nil, true
}

// Descends path (excluding its final element) and returns the mapping node that
// should directly contain that final element. With create set, missing
// intermediate maps are created; otherwise a missing intermediate yields a nil
// result. Returns an error if a node along the path exists but isn't a mapping.
func findContainingMap(node *yaml.Node, path []string, create bool) (*yaml.Node, error) {
	for _, key := range path[:len(path)-1] {
		if node.Kind != yaml.MappingNode {
			return nil, errors.New("yaml node in path is not a dictionary")
		}
		_, valueNode := LookupKey(node, key)
		if valueNode == nil {
			if !create {
				return nil, nil
			}
			valueNode = &yaml.Node{Kind: yaml.MappingNode}
			node.Content = append(node.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
				valueNode)
		}
		node = valueNode
	}
	if node.Kind != yaml.MappingNode {
		return nil, errors.New("yaml node in path is not a dictionary")
	}
	return node, nil
}

// Walks path from node and removes any mapping that is empty once its child has
// been removed, cascading upward. Stops at the first non-empty ancestor (which
// keeps every ancestor above it non-empty too).
func removeEmptyMaps(node *yaml.Node, path []string) {
	if len(path) == 0 {
		return
	}
	_, child := LookupKey(node, path[0])
	if child == nil {
		return
	}
	removeEmptyMaps(child, path[1:])
	if child.Kind == yaml.MappingNode && len(child.Content) == 0 {
		RemoveKey(node, path[0])
	}
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
		for i := range len(node.Content) {
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
