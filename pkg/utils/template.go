package utils

import (
	"bytes"
	"strings"
	"text/template"
)

func ResolveTemplate(templateStr string, object interface{}, funcs template.FuncMap) (string, error) {
	tmpl, err := template.New("template").Funcs(funcs).Option("missingkey=error").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, object); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ResolvePlaceholderString populates a template with values
func ResolvePlaceholderString(str string, arguments map[string]string) string {
	for key, value := range arguments {
		str = strings.Replace(str, "{{"+key+"}}", value, -1)
		str = strings.Replace(str, "{{."+key+"}}", value, -1)
	}
	return str
}
