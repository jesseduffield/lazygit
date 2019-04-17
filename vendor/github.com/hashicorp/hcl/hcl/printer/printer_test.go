package printer

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/hcl/parser"
)

var update = flag.Bool("update", false, "update golden files")

const (
	dataDir = "testdata"
)

type entry struct {
	source, golden string
}

// Use go test -update to create/update the respective golden files.
var data = []entry{
	{"complexhcl.input", "complexhcl.golden"},
	{"list.input", "list.golden"},
	{"list_comment.input", "list_comment.golden"},
	{"comment.input", "comment.golden"},
	{"comment_crlf.input", "comment.golden"},
	{"comment_aligned.input", "comment_aligned.golden"},
	{"comment_array.input", "comment_array.golden"},
	{"comment_end_file.input", "comment_end_file.golden"},
	{"comment_multiline_indent.input", "comment_multiline_indent.golden"},
	{"comment_multiline_no_stanza.input", "comment_multiline_no_stanza.golden"},
	{"comment_multiline_stanza.input", "comment_multiline_stanza.golden"},
	{"comment_newline.input", "comment_newline.golden"},
	{"comment_object_multi.input", "comment_object_multi.golden"},
	{"comment_standalone.input", "comment_standalone.golden"},
	{"empty_block.input", "empty_block.golden"},
	{"list_of_objects.input", "list_of_objects.golden"},
	{"multiline_string.input", "multiline_string.golden"},
	{"object_singleline.input", "object_singleline.golden"},
	{"object_with_heredoc.input", "object_with_heredoc.golden"},
}

func TestFiles(t *testing.T) {
	for _, e := range data {
		source := filepath.Join(dataDir, e.source)
		golden := filepath.Join(dataDir, e.golden)
		t.Run(e.source, func(t *testing.T) {
			check(t, source, golden)
		})
	}
}

func check(t *testing.T, source, golden string) {
	src, err := ioutil.ReadFile(source)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := format(src)
	if err != nil {
		t.Error(err)
		return
	}

	// update golden files if necessary
	if *update {
		if err := ioutil.WriteFile(golden, res, 0644); err != nil {
			t.Error(err)
		}
		return
	}

	// get golden
	gld, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Error(err)
		return
	}

	// formatted source and golden must be the same
	if err := diff(source, golden, res, gld); err != nil {
		t.Error(err)
		return
	}
}

// diff compares a and b.
func diff(aname, bname string, a, b []byte) error {
	var buf bytes.Buffer // holding long error message

	// compare lengths
	if len(a) != len(b) {
		fmt.Fprintf(&buf, "\nlength changed: len(%s) = %d, len(%s) = %d", aname, len(a), bname, len(b))
	}

	// compare contents
	line := 1
	offs := 1
	for i := 0; i < len(a) && i < len(b); i++ {
		ch := a[i]
		if ch != b[i] {
			fmt.Fprintf(&buf, "\n%s:%d:%d: %q", aname, line, i-offs+1, lineAt(a, offs))
			fmt.Fprintf(&buf, "\n%s:%d:%d: %q", bname, line, i-offs+1, lineAt(b, offs))
			fmt.Fprintf(&buf, "\n\n")
			break
		}
		if ch == '\n' {
			line++
			offs = i + 1
		}
	}

	if buf.Len() > 0 {
		return errors.New(buf.String())
	}
	return nil
}

// format parses src, prints the corresponding AST, verifies the resulting
// src is syntactically correct, and returns the resulting src or an error
// if any.
func format(src []byte) ([]byte, error) {
	formatted, err := Format(src)
	if err != nil {
		return nil, err
	}

	// make sure formatted output is syntactically correct
	if _, err := parser.Parse(formatted); err != nil {
		return nil, fmt.Errorf("parse: %s\n%s", err, formatted)
	}

	return formatted, nil
}

// lineAt returns the line in text starting at offset offs.
func lineAt(text []byte, offs int) []byte {
	i := offs
	for i < len(text) && text[i] != '\n' {
		i++
	}
	return text[offs:i]
}

// TestFormatParsable ensures that the output of Format() is can be parsed again.
func TestFormatValidOutput(t *testing.T) {
	cases := []string{
		"#\x00",
		"#\ue123t",
		"x=//\n0y=<<_\n_\n",
		"y=[1,//\n]",
		"Y=<<4\n4/\n\n\n/4/@=4/\n\n\n/4000000004\r\r\n00004\n",
		"x=<<_\n_\r\r\n_\n",
		"X=<<-\n\r\r\n",
	}

	for _, c := range cases {
		f, err := Format([]byte(c))
		if err != nil {
			// ignore these failures, not all inputs are valid HCL.
			t.Logf("Format(%q) = %v", c, err)
			continue
		}

		if _, err := parser.Parse(f); err != nil {
			t.Errorf("Format(%q) = %q; Parse(%q) = %v", c, f, f, err)
			continue
		}
	}
}
