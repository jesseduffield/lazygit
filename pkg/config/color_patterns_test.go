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

func TestColorPatternsExpandMergeKeys(t *testing.T) {
	scenarios := []struct {
		name     string
		content  string
		expected ColorPatterns
	}{
		{
			name: "Merge key with an alias",
			content: "base: &base\n" +
				"  '^a': red\n" +
				"  '^b': green\n" +
				"patterns:\n" +
				"  '^x': yellow\n" +
				"  <<: *base\n" +
				"  '^b': blue\n",
			expected: ColorPatterns{
				{Pattern: "^x", Color: "yellow"},
				{Pattern: "^a", Color: "red"},
				{Pattern: "^b", Color: "blue"},
			},
		},
		{
			name: "Merge key with a list of aliases",
			content: "first: &first\n" +
				"  '^a': red\n" +
				"  '^c': white\n" +
				"second: &second\n" +
				"  '^a': green\n" +
				"  '^b': blue\n" +
				"patterns:\n" +
				"  <<: [*first, *second]\n" +
				"  '^z': black\n",
			expected: ColorPatterns{
				{Pattern: "^a", Color: "red"},
				{Pattern: "^c", Color: "white"},
				{Pattern: "^b", Color: "blue"},
				{Pattern: "^z", Color: "black"},
			},
		},
		{
			name: "Merge key with a mapping written in place",
			content: "patterns:\n" +
				"  <<: {'^a': red}\n" +
				"  '^b': blue\n",
			expected: ColorPatterns{
				{Pattern: "^a", Color: "red"},
				{Pattern: "^b", Color: "blue"},
			},
		},
		{
			name: "Merged mapping with a merge key of its own",
			content: "base: &base\n" +
				"  '^a': red\n" +
				"extended: &extended\n" +
				"  <<: *base\n" +
				"  '^b': green\n" +
				"patterns:\n" +
				"  <<: *extended\n" +
				"  '^c': blue\n",
			expected: ColorPatterns{
				{Pattern: "^a", Color: "red"},
				{Pattern: "^b", Color: "green"},
				{Pattern: "^c", Color: "blue"},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			var config colorPatternsConfig
			err := yaml.Unmarshal([]byte(s.content), &config)
			assert.NoError(t, err)
			assert.Equal(t, s.expected, config.Patterns)
		})
	}
}

func TestColorPatternsRejectMergeKeyThatMergesItsOwnMapping(t *testing.T) {
	var config colorPatternsConfig
	err := yaml.Unmarshal([]byte("patterns: &patterns\n"+
		"  '^a': red\n"+
		"  <<: *patterns\n"), &config)
	assert.ErrorContains(t, err, "anchor 'patterns' value contains itself")
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
