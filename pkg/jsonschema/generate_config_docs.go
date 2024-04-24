package jsonschema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/orderedmap"
	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/samber/lo"

	"gopkg.in/yaml.v3"
)

type Node struct {
	Name        string
	Description string
	Default     any
	Children    []*Node
}

const (
	IndentLevel                  = 2
	DocumentationCommentStart    = "<!-- START CONFIG YAML: AUTOMATICALLY GENERATED with `go generate ./..., DO NOT UPDATE MANUALLY -->\n"
	DocumentationCommentEnd      = "<!-- END CONFIG YAML -->"
	DocumentationCommentStartLen = len(DocumentationCommentStart)
)

func insertBlankLines(buffer bytes.Buffer) bytes.Buffer {
	lines := strings.Split(strings.TrimRight(buffer.String(), "\n"), "\n")

	var newBuffer bytes.Buffer

	previousIndent := -1
	wasComment := false

	for _, line := range lines {
		trimmedLine := strings.TrimLeft(line, " ")
		indent := len(line) - len(trimmedLine)
		isComment := strings.HasPrefix(trimmedLine, "#")
		if isComment && !wasComment && indent <= previousIndent {
			newBuffer.WriteString("\n")
		}
		newBuffer.WriteString(line)
		newBuffer.WriteString("\n")
		previousIndent = indent
		wasComment = isComment
	}

	return newBuffer
}

func prepareMarshalledConfig(buffer bytes.Buffer) []byte {
	buffer = insertBlankLines(buffer)

	// Remove all `---` lines
	lines := strings.Split(strings.TrimRight(buffer.String(), "\n"), "\n")

	var newBuffer bytes.Buffer

	for _, line := range lines {
		if strings.TrimSpace(line) != "---" {
			newBuffer.WriteString(line)
			newBuffer.WriteString("\n")
		}
	}

	config := newBuffer.Bytes()

	// Add markdown yaml block tag
	config = append([]byte("```yaml\n"), config...)
	config = append(config, []byte("```\n")...)

	return config
}

func setComment(yamlNode *yaml.Node, description string) {
	// Workaround for the way yaml formats the HeadComment if it contains
	// blank lines: it renders these without a leading "#", but we want a
	// leading "#" even on blank lines. However, yaml respects it if the
	// HeadComment already contains a leading "#", so we prefix all lines
	// (including blank ones) with "#".
	yamlNode.HeadComment = strings.Join(
		lo.Map(strings.Split(description, "\n"), func(s string, _ int) string {
			if s == "" {
				return "#" // avoid trailing space on blank lines
			}
			return "# " + s
		}),
		"\n")
}

func (n *Node) MarshalYAML() (interface{}, error) {
	node := yaml.Node{
		Kind: yaml.MappingNode,
	}

	keyNode := yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: n.Name,
	}
	if n.Description != "" {
		setComment(&keyNode, n.Description)
	}

	if n.Default != nil {
		valueNode := yaml.Node{
			Kind: yaml.ScalarNode,
		}
		err := valueNode.Encode(n.Default)
		if err != nil {
			return nil, err
		}
		node.Content = append(node.Content, &keyNode, &valueNode)
	} else if len(n.Children) > 0 {
		childrenNode := yaml.Node{
			Kind: yaml.MappingNode,
		}
		for _, child := range n.Children {
			childYaml, err := child.MarshalYAML()
			if err != nil {
				return nil, err
			}

			childKey := yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: child.Name,
			}
			if child.Description != "" {
				setComment(&childKey, child.Description)
			}
			childYaml = childYaml.(*yaml.Node)
			childrenNode.Content = append(childrenNode.Content, childYaml.(*yaml.Node).Content...)
		}
		node.Content = append(node.Content, &keyNode, &childrenNode)
	}

	return &node, nil
}

func getDescription(v *orderedmap.OrderedMap) string {
	description, ok := v.Get("description")
	if !ok {
		description = ""
	}
	return description.(string)
}

func getDefault(v *orderedmap.OrderedMap) (error, any) {
	defaultValue, ok := v.Get("default")
	if ok {
		return nil, defaultValue
	}

	dataType, ok := v.Get("type")
	if ok {
		dataTypeString := dataType.(string)
		if dataTypeString == "string" {
			return nil, ""
		}
	}

	return errors.New("Failed to get default value"), nil
}

func parseNode(parent *Node, name string, value *orderedmap.OrderedMap) {
	description := getDescription(value)
	err, defaultValue := getDefault(value)
	if err == nil {
		leaf := &Node{Name: name, Description: description, Default: defaultValue}
		parent.Children = append(parent.Children, leaf)
	}

	properties, ok := value.Get("properties")
	if !ok {
		return
	}

	orderedProperties := properties.(orderedmap.OrderedMap)

	node := &Node{Name: name, Description: description}
	parent.Children = append(parent.Children, node)

	keys := orderedProperties.Keys()
	for _, name := range keys {
		value, _ := orderedProperties.Get(name)
		typedValue := value.(orderedmap.OrderedMap)
		parseNode(node, name, &typedValue)
	}
}

func writeToConfigDocs(config []byte) error {
	configPath := utils.GetLazyRootDirectory() + "/docs/Config.md"
	markdown, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Error reading Config.md file %w", err)
	}

	startConfigSectionIndex := bytes.Index(markdown, []byte(DocumentationCommentStart))
	if startConfigSectionIndex == -1 {
		return errors.New("Default config starting comment not found")
	}

	endConfigSectionIndex := bytes.Index(markdown[startConfigSectionIndex+DocumentationCommentStartLen:], []byte(DocumentationCommentEnd))
	if endConfigSectionIndex == -1 {
		return errors.New("Default config closing comment not found")
	}

	endConfigSectionIndex = endConfigSectionIndex + startConfigSectionIndex + DocumentationCommentStartLen

	newMarkdown := make([]byte, 0, len(markdown)-endConfigSectionIndex+startConfigSectionIndex+len(config))
	newMarkdown = append(newMarkdown, markdown[:startConfigSectionIndex+DocumentationCommentStartLen]...)
	newMarkdown = append(newMarkdown, config...)
	newMarkdown = append(newMarkdown, markdown[endConfigSectionIndex:]...)

	if err := os.WriteFile(configPath, newMarkdown, 0o644); err != nil {
		return fmt.Errorf("Error writing to file %w", err)
	}
	return nil
}

func GenerateConfigDocs() {
	content, err := os.ReadFile(GetSchemaDir() + "/config.json")
	if err != nil {
		panic("Error reading config.json")
	}

	schema := orderedmap.New()

	err = json.Unmarshal(content, &schema)
	if err != nil {
		panic("Failed to unmarshal config.json")
	}

	root, ok := schema.Get("properties")
	if !ok {
		panic("properties key not found in schema")
	}
	orderedRoot := root.(orderedmap.OrderedMap)

	rootNode := Node{}
	for _, name := range orderedRoot.Keys() {
		value, _ := orderedRoot.Get(name)
		typedValue := value.(orderedmap.OrderedMap)
		parseNode(&rootNode, name, &typedValue)
	}

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(IndentLevel)

	for _, child := range rootNode.Children {
		err := encoder.Encode(child)
		if err != nil {
			panic("Failed to Marshal document")
		}
	}
	encoder.Close()

	config := prepareMarshalledConfig(buffer)

	err = writeToConfigDocs(config)
	if err != nil {
		panic(err)
	}
}
