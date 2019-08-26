package internal

import (
	"bytes"
	"strings"
	"sync"
	gotemplate "text/template"
)

// Template stores the template for a string.
type Template struct {
	Src        string
	LeftDelim  string
	RightDelim string

	parseOnce      sync.Once
	parsedTemplate *gotemplate.Template
	parseError     error
}

func (t *Template) Execute(funcs gotemplate.FuncMap, data interface{}) (string, error) {
	leftDelim := t.LeftDelim
	if leftDelim == "" {
		leftDelim = "{{"
	}
	if !strings.Contains(t.Src, leftDelim) {
		// Fast path to avoid parsing a template that has no actions.
		return t.Src, nil
	}

	var gt *gotemplate.Template
	var err error
	if funcs == nil {
		t.parseOnce.Do(func() {
			// If funcs is nil, then we only need to parse this template once.
			t.parsedTemplate, t.parseError = gotemplate.New("").Delims(t.LeftDelim, t.RightDelim).Parse(t.Src)
		})
		gt, err = t.parsedTemplate, t.parseError
	} else {
		gt, err = gotemplate.New("").Delims(t.LeftDelim, t.RightDelim).Funcs(funcs).Parse(t.Src)
	}

	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := gt.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
