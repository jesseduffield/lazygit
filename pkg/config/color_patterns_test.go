package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type colorPatternsConfig struct {
	Patterns ColorPatterns `yaml:"patterns"`
}

func TestColorPatternsKeepTheirOrder(t *testing.T) {
	var config colorPatternsConfig
	err := yaml.Unmarshal([]byte("patterns:\n"+
		"  '^b': red\n"+
		"  '^a': '#00ff00'\n"+
		"  '^c': blue\n"), &config)
	assert.NoError(t, err)
	assert.Equal(t, ColorPatterns{
		{Pattern: "^b", Color: "red"},
		{Pattern: "^a", Color: "#00ff00"},
		{Pattern: "^c", Color: "blue"},
	}, config.Patterns)
}

func TestColorPatternsOfALaterFileComeFirst(t *testing.T) {
	var config colorPatternsConfig
	err := yaml.Unmarshal([]byte("patterns:\n"+
		"  '^a': red\n"+
		"  '^b': green\n"), &config)
	assert.NoError(t, err)
	err = yaml.Unmarshal([]byte("patterns:\n"+
		"  '^c': blue\n"+
		"  '^b': yellow\n"), &config)
	assert.NoError(t, err)
	assert.Equal(t, ColorPatterns{
		{Pattern: "^c", Color: "blue"},
		{Pattern: "^b", Color: "yellow"},
		{Pattern: "^a", Color: "red"},
	}, config.Patterns)
}

func TestColorPatternsMustBeAMapping(t *testing.T) {
	var config colorPatternsConfig
	err := yaml.Unmarshal([]byte("patterns: 5\n"), &config)
	assert.ErrorContains(t, err, "cannot unmarshal !!int `5` into map[string]string")
}

func TestColorPatternsSurviveMarshalling(t *testing.T) {
	config := colorPatternsConfig{Patterns: ColorPatterns{
		{Pattern: "^b", Color: "red"},
		{Pattern: "^a", Color: "#00ff00"},
		{Pattern: "true", Color: "blue"},
	}}
	content, err := yaml.Marshal(config)
	assert.NoError(t, err)

	var roundTripped colorPatternsConfig
	err = yaml.Unmarshal(content, &roundTripped)
	assert.NoError(t, err)
	assert.Equal(t, config, roundTripped)
}
