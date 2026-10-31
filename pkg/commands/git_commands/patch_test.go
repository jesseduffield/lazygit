package git_commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitObjectWithHeader(t *testing.T) {
	scenarios := []struct {
		name         string
		commitObject string
		headerName   string
		value        string
		expected     string
	}{
		{
			name: "replace parent",
			commitObject: "tree 1111\nparent 2222\n" +
				"author A U Thor <author@example.com> 1700000000 +0100\n" +
				"committer C O Mitter <committer@example.com> 1700000001 +0100\n" +
				"\nsubject\n\nbody\n",
			headerName: "parent",
			value:      "3333",
			expected: "tree 1111\nparent 3333\n" +
				"author A U Thor <author@example.com> 1700000000 +0100\n" +
				"committer C O Mitter <committer@example.com> 1700000001 +0100\n" +
				"\nsubject\n\nbody\n",
		},
		{
			name: "replace tree, keeping other headers",
			commitObject: "tree 1111\nparent 2222\n" +
				"author A U Thor <author@example.com> 1700000000 +0100\n" +
				"committer C O Mitter <committer@example.com> 1700000001 +0100\n" +
				"encoding ISO-8859-1\n" +
				"\nsubject\n",
			headerName: "tree",
			value:      "3333",
			expected: "tree 3333\nparent 2222\n" +
				"author A U Thor <author@example.com> 1700000000 +0100\n" +
				"committer C O Mitter <committer@example.com> 1700000001 +0100\n" +
				"encoding ISO-8859-1\n" +
				"\nsubject\n",
		},
		{
			name: "drop signature",
			commitObject: "tree 1111\nparent 2222\n" +
				"author A U Thor <author@example.com> 1700000000 +0100\n" +
				"committer C O Mitter <committer@example.com> 1700000001 +0100\n" +
				"gpgsig -----BEGIN PGP SIGNATURE-----\n \n iQEz\n -----END PGP SIGNATURE-----\n" +
				"gpgsig-sha256 -----BEGIN SSH SIGNATURE-----\n U1NI\n -----END SSH SIGNATURE-----\n" +
				"some-other-header value\n" +
				"\nsubject\n",
			headerName: "parent",
			value:      "3333",
			expected: "tree 1111\nparent 3333\n" +
				"author A U Thor <author@example.com> 1700000000 +0100\n" +
				"committer C O Mitter <committer@example.com> 1700000001 +0100\n" +
				"some-other-header value\n" +
				"\nsubject\n",
		},
		{
			name:         "empty message",
			commitObject: "tree 1111\nparent 2222\n\n",
			headerName:   "parent",
			value:        "3333",
			expected:     "tree 1111\nparent 3333\n\n",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, commitObjectWithHeader(s.commitObject, s.headerName, s.value))
		})
	}
}
