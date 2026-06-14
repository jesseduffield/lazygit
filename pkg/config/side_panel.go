package config

import (
	"fmt"

	"github.com/karimkhaleel/jsonschema"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// SidePanel is one entry in gui.sidePanels: a side panel made up of one or more
// tabs. For the common single-tab case it can be written in YAML as a single
// scalar string; a panel with tabs is written as a sequence of names.
type SidePanel []string

// ValidSidePanelTabs lists every name that may appear in gui.sidePanels. Each
// names a list that can stand alone as a panel or be grouped with others as the
// tabs of one panel. The resolver in the gui package must handle every entry
// here; a test enforces that the two stay in sync.
var ValidSidePanelTabs = []string{
	"status",
	"files",
	"worktrees",
	"submodules",
	"branches",
	"remotes",
	"tags",
	"commits",
	"reflog",
	"stash",
}

func (p *SidePanel) UnmarshalYAML(node *yaml.Node) error {
	var ss []string
	switch node.Kind {
	case yaml.ScalarNode:
		var s string
		if err := node.Decode(&s); err != nil {
			return err
		}
		ss = []string{s}
	case yaml.SequenceNode:
		if err := node.Decode(&ss); err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected a string or a sequence of strings for a side panel, got %v", node.Tag)
	}
	*p = ss
	return nil
}

func (p SidePanel) MarshalYAML() (any, error) {
	if len(p) == 1 {
		return p[0], nil
	}
	// Render multi-tab panels in flow style (`[a, b]`) rather than the default
	// block style, which is more compact and reads better in the generated docs.
	node := &yaml.Node{
		Kind:  yaml.SequenceNode,
		Style: yaml.FlowStyle,
	}
	for _, s := range p {
		node.Content = append(node.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: s,
		})
	}
	return node, nil
}

// JSONSchema describes a side panel as either a single tab name or a list of
// them, restricted to the known names.
func (SidePanel) JSONSchema() *jsonschema.Schema {
	names := lo.Map(ValidSidePanelTabs, func(name string, _ int) any { return name })
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{Type: "string", Enum: names},
			{Type: "array", Items: &jsonschema.Schema{Type: "string", Enum: names}},
		},
	}
}
