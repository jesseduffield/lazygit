package internal

import (
	"strings"
	gotemplate "text/template"
)

// Template stores the template for a string.
type Template struct {
	Src      string
	Template *gotemplate.Template
	ParseErr *error
}

func (t *Template) parse(leftDelim, rightDelim string, funcs gotemplate.FuncMap) error {
	if t.ParseErr == nil {
		if strings.Contains(t.Src, leftDelim) {
			gt, err := gotemplate.New("").Funcs(funcs).Delims(leftDelim, rightDelim).Parse(t.Src)
			t.Template = gt
			t.ParseErr = &err
		} else {
			t.ParseErr = new(error)
		}
	}
	return *t.ParseErr
}
