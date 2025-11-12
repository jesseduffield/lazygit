package jsonschema

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/karimkhaleel/jsonschema"
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

func wrapLine(line string, maxLineLength int) []string {
	result := []string{}
	startOfLine := 0
	lastSpaceIdx := -1
	for i, r := range line {
		// Don't break on "See https://..." lines
		if r == ' ' && line[startOfLine:i] != "See" {
			lastSpaceIdx = i + 1
		} else if i-startOfLine >= maxLineLength && lastSpaceIdx != -1 {
			result = append(result, line[startOfLine:lastSpaceIdx-1])
			startOfLine = lastSpaceIdx
			lastSpaceIdx = -1
		}
	}
	result = append(result, line[startOfLine:])
	return result
}

func setComment(yamlNode *yaml.Node, description string) {
	lines := strings.Split(description, "\n")
	wrappedLines := lo.Flatten(lo.Map(lines,
		func(line string, _ int) []string { return wrapLine(line, 78) }))

	// Workaround for the way yaml formats the HeadComment if it contains
	// blank lines: it renders these without a leading "#", but we want a
	// leading "#" even on blank lines. However, yaml respects it if the
	// HeadComment already contains a leading "#", so we prefix all lines
	// (including blank ones) with "#".
	yamlNode.HeadComment = strings.Join(
		lo.Map(wrappedLines, func(s string, _ int) string {
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

	if len(n.Children) > 0 {
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
	} else {
		valueNode := yaml.Node{
			Kind: yaml.ScalarNode,
		}
		err := valueNode.Encode(n.Default)
		if err != nil {
			return nil, err
		}
		node.Content = append(node.Content, &keyNode, &valueNode)
	}

	return &node, nil
}

func writeToConfigDocs(config []byte) error {
	configPath := utils.GetLazyRootDirectory() + "/docs-master/Config.md"
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

func GenerateConfigDocs(schema *jsonschema.Schema) {
	rootNode := &Node{
		Children: make([]*Node, 0),
	}

	recurseOverSchema(schema, schema.Definitions["UserConfig"], rootNode)

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

	err := writeToConfigDocs(config)
	if err != nil {
		panic(err)
	}
}

func recurseOverSchema(rootSchema, schema *jsonschema.Schema, parent *Node) {
	if schema == nil || schema.Properties == nil || schema.Properties.Len() == 0 {
		return
	}

	for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		subSchema := getSubSchema(rootSchema, schema, pair.Key)

		if strings.Contains(strings.ToLower(subSchema.Description), "deprecated") {
			continue
		}

		node := Node{
			Name:        pair.Key,
			Description: subSchema.Description,
			Default:     getZeroValue(subSchema.Default, subSchema.Type),
		}
		parent.Children = append(parent.Children, &node)
		recurseOverSchema(rootSchema, subSchema, &node)
	}
}

func getZeroValue(val any, t string) any {
	if !isZeroValue(val) {
		return val
	}

	switch t {
	case "string":
		return ""
	case "boolean":
		return false
	case "object":
		return map[string]any{}
	case "array":
		return []any{}
	default:
		return nil
	}
}
