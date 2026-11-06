package jsonschema

import (
	"testing"

	"github.com/karimkhaleel/jsonschema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfigDocsZeroIntegerDefault(t *testing.T) {
	type integerConfig struct {
		Width int `yaml:"width"`
	}

	r := &jsonschema.Reflector{FieldNameTag: "yaml"}
	schema := r.Reflect(&integerConfig{})
	setDefaultVals(schema, schema.Definitions["integerConfig"], integerConfig{})
	root := &Node{Name: "gui"}
	recurseOverSchema(schema, schema.Definitions["integerConfig"], root, nil)
	output, err := yaml.Marshal(root)
	assert.NoError(t, err)
	assert.Equal(t, "gui:\n    width: 0\n", string(output))
}
