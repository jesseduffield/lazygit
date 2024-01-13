package ansi

import (
	"regexp"
	"strings"
	"text/template"
)

// TemplateFuncMap contains a few useful template helpers
var (
	TemplateFuncMap = template.FuncMap{
		"Left": func(values ...interface{}) string {
			s := values[0].(string)
			n := values[1].(int)
			if n > len(s) {
				n = len(s)
			}

			return s[:n]
		},
		"Matches": func(values ...interface{}) bool {
			ok, _ := regexp.MatchString(values[1].(string), values[0].(string))
			return ok
		},
		"Mid": func(values ...interface{}) string {
			s := values[0].(string)
			l := values[1].(int)
			if l > len(s) {
				l = len(s)
			}

			if len(values) > 2 {
				r := values[2].(int)
				if r > len(s) {
					r = len(s)
				}
				return s[l:r]
			}
			return s[l:]
		},
		"Right": func(values ...interface{}) string {
			s := values[0].(string)
			n := len(s) - values[1].(int)
			if n < 0 {
				n = 0
			}

			return s[n:]
		},
		"Last": func(values ...interface{}) string {
			return values[0].([]string)[len(values[0].([]string))-1]
		},
		// strings functions
		"Compare":      strings.Compare, // 1.5+ only
		"Contains":     strings.Contains,
		"ContainsAny":  strings.ContainsAny,
		"Count":        strings.Count,
		"EqualFold":    strings.EqualFold,
		"HasPrefix":    strings.HasPrefix,
		"HasSuffix":    strings.HasSuffix,
		"Index":        strings.Index,
		"IndexAny":     strings.IndexAny,
		"Join":         strings.Join,
		"LastIndex":    strings.LastIndex,
		"LastIndexAny": strings.LastIndexAny,
		"Repeat":       strings.Repeat,
		"Replace":      strings.Replace,
		"Split":        strings.Split,
		"SplitAfter":   strings.SplitAfter,
		"SplitAfterN":  strings.SplitAfterN,
		"SplitN":       strings.SplitN,
		"Title":        strings.Title,
		"ToLower":      strings.ToLower,
		"ToTitle":      strings.ToTitle,
		"ToUpper":      strings.ToUpper,
		"Trim":         strings.Trim,
		"TrimLeft":     strings.TrimLeft,
		"TrimPrefix":   strings.TrimPrefix,
		"TrimRight":    strings.TrimRight,
		"TrimSpace":    strings.TrimSpace,
		"TrimSuffix":   strings.TrimSuffix,
	}
)
