package utils

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

func ResolveTemplate(templateStr string, object any, funcs template.FuncMap) (string, error) {
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

// ResolveCommandTemplate is like ResolveTemplate but shell-quotes all string
// values in the object using quoteFn before template execution, preventing
// command injection from git-derived data such as branch names or file paths.
func ResolveCommandTemplate(templateStr string, object any, funcs template.FuncMap, quoteFn func(string) string) (string, error) {
	return ResolveTemplate(templateStr, shellQuoteAll(object, quoteFn), funcs)
}

func shellQuoteAll(v any, quoteFn func(string) string) any {
	if v == nil {
		return nil
	}
	result := shellQuoteValue(reflect.ValueOf(v), quoteFn)
	if result.IsValid() {
		return result.Interface()
	}
	return nil
}

func shellQuoteValue(v reflect.Value, quoteFn func(string) string) reflect.Value {
	switch v.Kind() {
	case reflect.String:
		return reflect.ValueOf(quoteFn(v.String())).Convert(v.Type())
	case reflect.Ptr:
		if v.IsNil() {
			return v
		}
		result := reflect.New(v.Type().Elem())
		result.Elem().Set(shellQuoteValue(v.Elem(), quoteFn))
		return result
	case reflect.Struct:
		result := reflect.New(v.Type()).Elem()
		for i := range v.NumField() {
			if result.Field(i).CanSet() {
				result.Field(i).Set(shellQuoteValue(v.Field(i), quoteFn))
			}
		}
		return result
	case reflect.Slice:
		if v.IsNil() {
			return v
		}
		result := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := range v.Len() {
			result.Index(i).Set(shellQuoteValue(v.Index(i), quoteFn))
		}
		return result
	case reflect.Map:
		if v.IsNil() {
			return v
		}
		result := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			result.SetMapIndex(key, shellQuoteValue(v.MapIndex(key), quoteFn))
		}
		return result
	default:
		return v
	}
}

// ResolvePlaceholderString populates a template with values
func ResolvePlaceholderString(str string, arguments map[string]string) string {
	oldnews := make([]string, 0, len(arguments)*4)
	for key, value := range arguments {
		oldnews = append(oldnews,
			"{{"+key+"}}", value,
			"{{."+key+"}}", value,
		)
	}
	return strings.NewReplacer(oldnews...).Replace(str)
}
