package yaml_utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestUpdateYamlValue(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		path        []string
		value       string
		expectedOut string
		expectedErr string
	}{
		{
			name:        "update value",
			in:          "foo: bar\n",
			path:        []string{"foo"},
			value:       "baz",
			expectedOut: "foo: baz\n",
			expectedErr: "",
		},
		{
			name:        "add new key and value",
			in:          "foo: bar\n",
			path:        []string{"foo2"},
			value:       "baz",
			expectedOut: "foo: bar\nfoo2: baz\n",
			expectedErr: "",
		},
		{
			name:        "add new key and value when document was empty",
			in:          "",
			path:        []string{"foo"},
			value:       "bar",
			expectedOut: "foo: bar\n",
			expectedErr: "",
		},
		{
			name:        "preserve inline comment",
			in:          "foo: bar # my comment\n",
			path:        []string{"foo2"},
			value:       "baz",
			expectedOut: "foo: bar # my comment\nfoo2: baz\n",
			expectedErr: "",
		},
		{
			name:        "nested update",
			in:          "foo:\n  bar: baz\n",
			path:        []string{"foo", "bar"},
			value:       "qux",
			expectedOut: "foo:\n  bar: qux\n",
			expectedErr: "",
		},
		{
			name:        "nested where parents doesn't exist yet",
			in:          "",
			path:        []string{"foo", "bar", "baz"},
			value:       "qux",
			expectedOut: "foo:\n  bar:\n    baz: qux\n",
			expectedErr: "",
		},
		{
			name:        "don't rewrite file if value didn't change",
			in:          "foo:\n  bar: baz\n",
			path:        []string{"foo", "bar"},
			value:       "baz",
			expectedOut: "foo:\n  bar: baz\n",
			expectedErr: "",
		},

		// Error cases
		{
			name:        "existing document is not a dictionary",
			in:          "42\n",
			path:        []string{"foo"},
			value:       "bar",
			expectedOut: "42\n",
			expectedErr: "yaml document is not a dictionary",
		},
		{
			name:        "trying to update a note that is not a scalar",
			in:          "foo: [1, 2, 3]\n",
			path:        []string{"foo"},
			value:       "bar",
			expectedOut: "foo: [1, 2, 3]\n",
			expectedErr: "yaml node is not a scalar",
		},
		{
			name:        "not all path elements are dictionaries",
			in:          "foo:\n  bar: [1, 2, 3]\n",
			path:        []string{"foo", "bar", "baz"},
			value:       "qux",
			expectedOut: "foo:\n  bar: [1, 2, 3]\n",
			expectedErr: "yaml node in path is not a dictionary",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, actualErr := UpdateYamlValue([]byte(test.in), test.path, test.value)
			if test.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.EqualError(t, actualErr, test.expectedErr)
			}

			assert.Equal(t, test.expectedOut, string(out))
		})
	}
}

func TestRenameYamlKey(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		path        []string
		newKey      string
		expectedOut string
		expectedErr string
	}{
		{
			name:        "rename key",
			in:          "foo: 5\n",
			path:        []string{"foo"},
			newKey:      "bar",
			expectedOut: "bar: 5\n",
			expectedErr: "",
		},
		{
			name:        "rename key, nested",
			in:          "foo:\n  bar: 5\n",
			path:        []string{"foo", "bar"},
			newKey:      "baz",
			expectedOut: "foo:\n  baz: 5\n",
			expectedErr: "",
		},
		{
			name:        "rename non-scalar key",
			in:          "foo:\n  bar: 5\n",
			path:        []string{"foo"},
			newKey:      "qux",
			expectedOut: "qux:\n  bar: 5\n",
			expectedErr: "",
		},
		{
			name:        "don't rewrite file if value didn't change",
			in:          "foo:\n  bar: 5\n",
			path:        []string{"nonExistingKey"},
			newKey:      "qux",
			expectedOut: "foo:\n  bar: 5\n",
			expectedErr: "",
		},

		// Error cases
		{
			name:        "existing document is not a dictionary",
			in:          "42\n",
			path:        []string{"foo"},
			newKey:      "bar",
			expectedOut: "42\n",
			expectedErr: "yaml node in path is not a dictionary",
		},
		{
			name:        "not all path elements are dictionaries",
			in:          "foo:\n  bar: [1, 2, 3]\n",
			path:        []string{"foo", "bar", "baz"},
			newKey:      "qux",
			expectedOut: "foo:\n  bar: [1, 2, 3]\n",
			expectedErr: "yaml node in path is not a dictionary",
		},
		{
			name:        "new key exists",
			in:          "foo: 5\nbar: 7\n",
			path:        []string{"foo"},
			newKey:      "bar",
			expectedOut: "foo: 5\nbar: 7\n",
			expectedErr: "new key `bar' already exists",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node := unmarshalForTest(t, test.in)
			actualErr := RenameYamlKey(&node, test.path, test.newKey)
			if test.expectedErr == "" {
				assert.NoError(t, actualErr)
			} else {
				assert.EqualError(t, actualErr, test.expectedErr)
			}
			out := marshalForTest(t, &node)

			assert.Equal(t, test.expectedOut, out)
		})
	}
}

func TestWalk_paths(t *testing.T) {
	tests := []struct {
		name          string
		document      string
		expectedPaths []string
	}{
		{
			name:          "empty document",
			document:      "",
			expectedPaths: []string{},
		},
		{
			name:          "scalar",
			document:      "x: 5",
			expectedPaths: []string{"", "x"}, // called with an empty path for the root node
		},
		{
			name:          "nested",
			document:      "foo:\n  x: 5",
			expectedPaths: []string{"", "foo", "foo.x"},
		},
		{
			name:          "deeply nested",
			document:      "foo:\n  bar:\n    baz: 5",
			expectedPaths: []string{"", "foo", "foo.bar", "foo.bar.baz"},
		},
		{
			name:          "array",
			document:      "foo:\n  bar: [3, 7]",
			expectedPaths: []string{"", "foo", "foo.bar", "foo.bar[0]", "foo.bar[1]"},
		},
		{
			name:          "nested arrays",
			document:      "foo:\n  bar: [[3, 7], [8, 9]]",
			expectedPaths: []string{"", "foo", "foo.bar", "foo.bar[0]", "foo.bar[0][0]", "foo.bar[0][1]", "foo.bar[1]", "foo.bar[1][0]", "foo.bar[1][1]"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node := unmarshalForTest(t, test.document)
			paths := []string{}
			err := Walk(&node, func(node *yaml.Node, path string) {
				paths = append(paths, path)
			})

			assert.NoError(t, err)
			assert.Equal(t, test.expectedPaths, paths)
		})
	}
}

func TestWalk_inPlaceChanges(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		callback    func(node *yaml.Node, path string)
		expectedOut string
	}{
		{
			name:     "no change",
			in:       "x: 5",
			callback: func(node *yaml.Node, path string) {},
		},
		{
			name: "change value",
			in:   "x: 5\ny: 3",
			callback: func(node *yaml.Node, path string) {
				if path == "x" {
					node.Value = "7"
				}
			},
			expectedOut: "x: 7\ny: 3\n",
		},
		{
			name: "change nested value",
			in:   "x:\n  y: 5",
			callback: func(node *yaml.Node, path string) {
				if path == "x.y" {
					node.Value = "7"
				}
			},
			expectedOut: "x:\n  y: 7\n",
		},
		{
			name: "change array value",
			in:   "x:\n  - y: 5",
			callback: func(node *yaml.Node, path string) {
				if path == "x[0].y" {
					node.Value = "7"
				}
			},
			expectedOut: "x:\n  - y: 7\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node := unmarshalForTest(t, test.in)
			err := Walk(&node, test.callback)
			assert.NoError(t, err)
			if test.expectedOut == "" {
				unmodifiedOriginal := unmarshalForTest(t, test.in)
				assert.Equal(t, unmodifiedOriginal, node)
			} else {
				result := marshalForTest(t, &node)
				assert.Equal(t, test.expectedOut, result)
			}
		})
	}
}

func TestTransformNode(t *testing.T) {
	transformIntValueToString := func(node *yaml.Node) error {
		if node.Kind == yaml.ScalarNode {
			if node.ShortTag() == "!!int" {
				node.Tag = "!!str"
				return nil
			} else if node.ShortTag() == "!!str" {
				// We have already transformed it,
				return nil
			} else {
				return fmt.Errorf("Node was of bad type")
			}
		} else {
			return fmt.Errorf("Node was not a scalar")
		}
	}

	tests := []struct {
		name        string
		in          string
		path        []string
		transform   func(node *yaml.Node) error
		expectedOut string
	}{
		{
			name:      "Path not present",
			in:        "foo: 1",
			path:      []string{"bar"},
			transform: transformIntValueToString,
		},
		{
			name: "Part of path present",
			in: `
foo:
  bar: 2`,
			path:      []string{"foo", "baz"},
			transform: transformIntValueToString,
		},
		{
			name: "Successfully Transforms to string",
			in: `
foo:
  bar: 2`,
			path:      []string{"foo", "bar"},
			transform: transformIntValueToString,
			expectedOut: `foo:
  bar: "2"
`, // Note the trailing newline changes because of how it re-marshalls
		},
		{
			name: "Does nothing when already transformed",
			in: `
foo:
  bar: "2"`,
			path:      []string{"foo", "bar"},
			transform: transformIntValueToString,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node := unmarshalForTest(t, test.in)
			err := TransformNode(&node, test.path, test.transform)
			if err != nil {
				t.Fatal(err)
			}
			if test.expectedOut == "" {
				unmodifiedOriginal := unmarshalForTest(t, test.in)
				assert.Equal(t, unmodifiedOriginal, node)
			} else {
				result := marshalForTest(t, &node)
				assert.Equal(t, test.expectedOut, result)
			}
		})
	}
}

func unmarshalForTest(t *testing.T, input string) yaml.Node {
	t.Helper()
	var node yaml.Node
	err := yaml.Unmarshal([]byte(input), &node)
	if err != nil {
		t.Fatal(err)
	}
	return node
}

func marshalForTest(t *testing.T, node *yaml.Node) string {
	t.Helper()
	result, err := YamlMarshal(node)
	if err != nil {
		t.Fatal(err)
	}
	return string(result)
}
