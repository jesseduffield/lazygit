package jsonschema

import (
	"testing"

	"github.com/karimkhaleel/jsonschema"
	"github.com/stretchr/testify/assert"
)

type testColors struct {
	Border []string `yaml:"border"`
}

type testGui struct {
	Theme     testColors `yaml:"theme"`
	DarkTheme testColors `yaml:"darkTheme"`
}

type testConfig struct {
	Gui        testGui    `yaml:"gui"`
	OtherTheme testColors `yaml:"otherTheme"`
}

func TestInlineSharedStructDefinitions(t *testing.T) {
	r := &jsonschema.Reflector{FieldNameTag: "yaml"}
	schema := r.Reflect(&testConfig{})

	repeatedStructs := inlineSharedStructDefinitions(schema, "testConfig")
	setDefaultVals(schema, schema.Definitions["testConfig"], testConfig{
		Gui: testGui{Theme: testColors{Border: []string{"green"}}},
	})

	gui := getSubSchema(schema, schema.Definitions["testConfig"], "gui")
	theme := getSubSchema(schema, gui, "theme")
	darkTheme := getSubSchema(schema, gui, "darkTheme")
	otherTheme := getSubSchema(schema, schema.Definitions["testConfig"], "otherTheme")

	assert.Equal(t, []*jsonschema.Schema{darkTheme, otherTheme}, repeatedStructs)
	assert.NotContains(t, schema.Definitions, "testColors")

	assert.Equal(t, []string{"green"}, getSubSchema(schema, theme, "border").Default)
	assert.Nil(t, getSubSchema(schema, darkTheme, "border").Default)
	assert.Nil(t, getSubSchema(schema, otherTheme, "border").Default)
}
