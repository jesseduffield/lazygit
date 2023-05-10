package yaml_utils

import "testing"

func TestUpdateYaml(t *testing.T) {
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
			name:        "preserve inline comment",
			in:          "foo: bar # my comment\n",
			path:        []string{"foo2"},
			value:       "baz",
			expectedOut: "foo: bar # my comment\nfoo2: baz\n",
			expectedErr: "",
		},
		{
			name:  "nested update",
			in:    "foo:\n  bar: baz\n",
			path:  []string{"foo", "bar"},
			value: "qux",
			// indentation is not preserved. See https://github.com/go-yaml/yaml/issues/899
			expectedOut: "foo:\n    bar: qux\n",
			expectedErr: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			out, err := UpdateYaml([]byte(test.in), test.path, test.value)
			if test.expectedErr != "" {
				if err == nil {
					t.Errorf("expected error %q but got none", test.expectedErr)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if string(out) != test.expectedOut {
				t.Errorf("expected %q but got %q", test.expectedOut, string(out))
			}
		})
	}
}
