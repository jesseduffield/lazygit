package config

import (
	"github.com/karimkhaleel/jsonschema"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// SidePanel is one entry in gui.sidePanels: a side panel made up of one or more
// tabs, written in YAML as a list of tab names (e.g. [files, worktrees]).
type SidePanel []string

// ValidSidePanelTabs lists every name that may appear in gui.sidePanels. Each
// names a list that can stand alone as a panel or be grouped with others as the
// tabs of one panel. The resolver in the gui package must handle every entry
// here; a test enforces that the two stay in sync.
var ValidSidePanelTabs = []string{
	"status",
	"repos",
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

func (p SidePanel) MarshalYAML() (any, error) {
	// Render in flow style (`[a, b]`) rather than the default block style, which
	// is more compact and reads better in the generated docs.
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

// JSONSchema describes a side panel as a list of tab names, restricted to the
// known names.
func (SidePanel) JSONSchema() *jsonschema.Schema {
	names := lo.Map(ValidSidePanelTabs, func(name string, _ int) any { return name })
	return &jsonschema.Schema{
		Type:  "array",
		Items: &jsonschema.Schema{Type: "string", Enum: names},
	}
}

// SidePanelsForMultiRepo puts repos in the place of status, or first if there
// is no status, because multi-repo mode focuses repos at startup.
func SidePanelsForMultiRepo(panels []SidePanel) []SidePanel {
	hasTab := func(name string) bool {
		return lo.SomeBy(panels, func(panel SidePanel) bool { return lo.Contains(panel, name) })
	}
	if hasTab("repos") {
		return panels
	}
	if !hasTab("status") {
		return append([]SidePanel{{"repos"}}, panels...)
	}

	return lo.Map(panels, func(panel SidePanel, _ int) SidePanel {
		return lo.Map(panel, func(name string, _ int) string {
			if name == "status" {
				return "repos"
			}
			return name
		})
	})
}
