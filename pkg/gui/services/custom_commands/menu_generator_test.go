package custom_commands

import (
	"testing"

	"github.com/lobes/lazytask/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestMenuGenerator(t *testing.T) {
	type scenario struct {
		testName    string
		cmdOut      string
		filter      string
		valueFormat string
		labelFormat string
		test        func([]*commandMenuItem, error)
	}

	scenarios := []scenario{
		{
			"Extract remote branch name",
			"upstream/pr-1",
			"(?P<remote>[a-z_]+)/(?P<branch>.*)",
			"{{ .branch }}",
			"Remote: {{ .remote }}",
			func(actualEntry []*commandMenuItem, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "pr-1", actualEntry[0].value)
				assert.EqualValues(t, "Remote: upstream", actualEntry[0].label)
			},
		},
		{
			"Multiple named groups with empty labelFormat",
			"upstream/pr-1",
			"(?P<remote>[a-z]*)/(?P<branch>.*)",
			"{{ .branch }}|{{ .remote }}",
			"",
			func(actualEntry []*commandMenuItem, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "pr-1|upstream", actualEntry[0].value)
				assert.EqualValues(t, "pr-1|upstream", actualEntry[0].label)
			},
		},
		{
			"Multiple named groups with group ids",
			"upstream/pr-1",
			"(?P<remote>[a-z]*)/(?P<branch>.*)",
			"{{ .group_2 }}|{{ .group_1 }}",
			"Remote: {{ .group_1 }}",
			func(actualEntry []*commandMenuItem, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "pr-1|upstream", actualEntry[0].value)
				assert.EqualValues(t, "Remote: upstream", actualEntry[0].label)
			},
		},
		{
			"No named groups",
			"upstream/pr-1",
			"([a-z]*)/(.*)",
			"{{ .group_2 }}|{{ .group_1 }}",
			"Remote: {{ .group_1 }}",
			func(actualEntry []*commandMenuItem, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "pr-1|upstream", actualEntry[0].value)
				assert.EqualValues(t, "Remote: upstream", actualEntry[0].label)
			},
		},
		{
			"No filter",
			"upstream/pr-1",
			"",
			"",
			"",
			func(actualEntry []*commandMenuItem, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "upstream/pr-1", actualEntry[0].value)
				assert.EqualValues(t, "upstream/pr-1", actualEntry[0].label)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			s.test(NewMenuGenerator(utils.NewDummyCommon()).call(s.cmdOut, s.filter, s.valueFormat, s.labelFormat))
		})
	}
}
