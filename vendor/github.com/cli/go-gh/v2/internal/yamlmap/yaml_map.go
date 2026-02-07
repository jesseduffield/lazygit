// Package yamlmap is a wrapper of gopkg.in/yaml.v3 for interacting
// with yaml data as if it were a map.
package yamlmap

import (
	"errors"

	"gopkg.in/yaml.v3"
)

const (
	modified = "modifed"
)

type Map struct {
	*yaml.Node
}

var ErrNotFound = errors.New("not found")
var ErrInvalidYaml = errors.New("invalid yaml")
var ErrInvalidFormat = errors.New("invalid format")

func StringValue(value string) *Map {
	return &Map{&yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
	}}
}

func MapValue() *Map {
	return &Map{&yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}}
}

func NullValue() *Map {
	return &Map{&yaml.Node{
		Kind: yaml.ScalarNode,
		Tag:  "!!null",
	}}
}

func Unmarshal(data []byte) (*Map, error) {
	var root yaml.Node
	err := yaml.Unmarshal(data, &root)
	if err != nil {
		return nil, ErrInvalidYaml
	}
	if len(root.Content) == 0 {
		return MapValue(), nil
	}
	if root.Content[0].Kind != yaml.MappingNode {
		return nil, ErrInvalidFormat
	}
	return &Map{root.Content[0]}, nil
}

func Marshal(m *Map) ([]byte, error) {
	return yaml.Marshal(m.Node)
}

func (m *Map) AddEntry(key string, value *Map) {
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: key,
	}
	m.Content = append(m.Content, keyNode, value.Node)
	m.SetModified()
}

func (m *Map) Empty() bool {
	return m.Content == nil || len(m.Content) == 0
}

func (m *Map) FindEntry(key string) (*Map, error) {
	// Note: The content slice of a yamlMap looks like [key1, value1, key2, value2, ...].
	// When iterating over the content slice we only want to compare the keys of the yamlMap.
	for i, v := range m.Content {
		if i%2 != 0 {
			continue
		}
		if v.Value == key {
			if i+1 < len(m.Content) {
				return &Map{m.Content[i+1]}, nil
			}
		}
	}
	return nil, ErrNotFound
}

func (m *Map) Keys() []string {
	// Note: The content slice of a yamlMap looks like [key1, value1, key2, value2, ...].
	// When iterating over the content slice we only want to select the keys of the yamlMap.
	keys := []string{}
	for i, v := range m.Content {
		if i%2 != 0 {
			continue
		}
		keys = append(keys, v.Value)
	}
	return keys
}

func (m *Map) RemoveEntry(key string) error {
	// Note: The content slice of a yamlMap looks like [key1, value1, key2, value2, ...].
	// When iterating over the content slice we only want to compare the keys of the yamlMap.
	// If we find they key to remove, remove the key and its value from the content slice.
	found, skipNext := false, false
	newContent := []*yaml.Node{}
	for i, v := range m.Content {
		if skipNext {
			skipNext = false
			continue
		}
		if i%2 != 0 || v.Value != key {
			newContent = append(newContent, v)
		} else {
			found = true
			skipNext = true
			m.SetModified()
		}
	}
	if !found {
		return ErrNotFound
	}
	m.Content = newContent
	return nil
}

func (m *Map) SetEntry(key string, value *Map) {
	// Note: The content slice of a yamlMap looks like [key1, value1, key2, value2, ...].
	// When iterating over the content slice we only want to compare the keys of the yamlMap.
	// If we find they key to set, set the next item in the content slice to the new value.
	m.SetModified()
	for i, v := range m.Content {
		if i%2 != 0 || v.Value != key {
			continue
		}
		if v.Value == key {
			if i+1 < len(m.Content) {
				m.Content[i+1] = value.Node
				return
			}
		}
	}
	m.AddEntry(key, value)
}

// Note: This is a hack to introduce the concept of modified/unmodified
// on top of gopkg.in/yaml.v3. This works by setting the Value property
// of a MappingNode to a specific value and then later checking if the
// node's Value property is that specific value. When a MappingNode gets
// output as a string the Value property is not used, thus changing it
// has no impact for our purposes.
func (m *Map) SetModified() {
	// Can not mark a non-mapping node as modified
	if m.Node.Kind != yaml.MappingNode && m.Node.Tag == "!!null" {
		m.Node.Kind = yaml.MappingNode
		m.Node.Tag = "!!map"
	}
	if m.Node.Kind == yaml.MappingNode {
		m.Node.Value = modified
	}
}

// Traverse map using BFS to set all nodes as unmodified.
func (m *Map) SetUnmodified() {
	i := 0
	queue := []*yaml.Node{m.Node}
	for {
		if i > (len(queue) - 1) {
			break
		}
		q := queue[i]
		i = i + 1
		if q.Kind != yaml.MappingNode {
			continue
		}
		q.Value = ""
		queue = append(queue, q.Content...)
	}
}

// Traverse map using BFS to searach for any nodes that have been modified.
func (m *Map) IsModified() bool {
	i := 0
	queue := []*yaml.Node{m.Node}
	for {
		if i > (len(queue) - 1) {
			break
		}
		q := queue[i]
		i = i + 1
		if q.Kind != yaml.MappingNode {
			continue
		}
		if q.Value == modified {
			return true
		}
		queue = append(queue, q.Content...)
	}
	return false
}

func (m *Map) String() string {
	data, err := Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}
