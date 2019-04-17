// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipeline

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/text/collate"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/internal"
	"golang.org/x/text/internal/catmsg"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
)

var transRe = regexp.MustCompile(`messages\.(.*)\.json`)

// Generate writes a Go file with the given package name to w, which defines a
// Catalog with translated messages.
func Generate(w io.Writer, pkg string, extracted *Locale, trans ...*Locale) (n int, err error) {
	// TODO: add in external input. Right now we assume that all files are
	// manually created and stored in the textdata directory.

	// Build up index of translations and original messages.
	translations := map[language.Tag]map[string]Message{}
	languages := []language.Tag{}
	langVars := []string{}
	usedKeys := map[string]int{}

	for _, loc := range trans {
		tag := loc.Language
		if _, ok := translations[tag]; !ok {
			translations[tag] = map[string]Message{}
			languages = append(languages, tag)
		}
		for _, m := range loc.Messages {
			if !m.Translation.IsEmpty() {
				for _, id := range m.ID {
					if _, ok := translations[tag][id]; ok {
						logf("Duplicate translation in locale %q for message %q", tag, id)
					}
					translations[tag][id] = m
				}
			}
		}
	}

	// Verify completeness and register keys.
	internal.SortTags(languages)

	for _, tag := range languages {
		langVars = append(langVars, strings.Replace(tag.String(), "-", "_", -1))
		dict := translations[tag]
		for _, msg := range extracted.Messages {
			for _, id := range msg.ID {
				if trans, ok := dict[id]; ok && !trans.Translation.IsEmpty() {
					if _, ok := usedKeys[msg.Key]; !ok {
						usedKeys[msg.Key] = len(usedKeys)
					}
					break
				}
				// TODO: log missing entry.
				logf("%s: Missing entry for %q.", tag, id)
			}
		}
	}

	cw := gen.NewCodeWriter()

	x := &struct {
		Fallback  language.Tag
		Languages []string
	}{
		Fallback:  extracted.Language,
		Languages: langVars,
	}

	if err := lookup.Execute(cw, x); err != nil {
		return 0, wrap(err, "error")
	}

	keyToIndex := []string{}
	for k := range usedKeys {
		keyToIndex = append(keyToIndex, k)
	}
	sort.Strings(keyToIndex)
	fmt.Fprint(cw, "var messageKeyToIndex = map[string]int{\n")
	for _, k := range keyToIndex {
		fmt.Fprintf(cw, "%q: %d,\n", k, usedKeys[k])
	}
	fmt.Fprint(cw, "}\n\n")

	for i, tag := range languages {
		dict := translations[tag]
		a := make([]string, len(usedKeys))
		for _, msg := range extracted.Messages {
			for _, id := range msg.ID {
				if trans, ok := dict[id]; ok && !trans.Translation.IsEmpty() {
					m, err := assemble(&msg, &trans.Translation)
					if err != nil {
						return 0, wrap(err, "error")
					}
					// TODO: support macros.
					data, err := catmsg.Compile(tag, nil, m)
					if err != nil {
						return 0, wrap(err, "error")
					}
					key := usedKeys[msg.Key]
					if d := a[key]; d != "" && d != data {
						logf("Duplicate non-consistent translation for key %q, picking the one for message %q", msg.Key, id)
					}
					a[key] = string(data)
					break
				}
			}
		}
		index := []uint32{0}
		p := 0
		for _, s := range a {
			p += len(s)
			index = append(index, uint32(p))
		}

		cw.WriteVar(langVars[i]+"Index", index)
		cw.WriteConst(langVars[i]+"Data", strings.Join(a, ""))
	}
	return cw.WriteGo(w, pkg, "")
}

func assemble(m *Message, t *Text) (msg catmsg.Message, err error) {
	keys := []string{}
	for k := range t.Var {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var a []catmsg.Message
	for _, k := range keys {
		t := t.Var[k]
		m, err := assemble(m, &t)
		if err != nil {
			return nil, err
		}
		a = append(a, &catmsg.Var{Name: k, Message: m})
	}
	if t.Select != nil {
		s, err := assembleSelect(m, t.Select)
		if err != nil {
			return nil, err
		}
		a = append(a, s)
	}
	if t.Msg != "" {
		sub, err := m.Substitute(t.Msg)
		if err != nil {
			return nil, err
		}
		a = append(a, catmsg.String(sub))
	}
	switch len(a) {
	case 0:
		return nil, errorf("generate: empty message")
	case 1:
		return a[0], nil
	default:
		return catmsg.FirstOf(a), nil

	}
}

func assembleSelect(m *Message, s *Select) (msg catmsg.Message, err error) {
	cases := []string{}
	for c := range s.Cases {
		cases = append(cases, c)
	}
	sortCases(cases)

	caseMsg := []interface{}{}
	for _, c := range cases {
		cm := s.Cases[c]
		m, err := assemble(m, &cm)
		if err != nil {
			return nil, err
		}
		caseMsg = append(caseMsg, c, m)
	}

	ph := m.Placeholder(s.Arg)

	switch s.Feature {
	case "plural":
		// TODO: only printf-style selects are supported as of yet.
		return plural.Selectf(ph.ArgNum, ph.String, caseMsg...), nil
	}
	return nil, errorf("unknown feature type %q", s.Feature)
}

func sortCases(cases []string) {
	// TODO: implement full interface.
	sort.Slice(cases, func(i, j int) bool {
		if cases[j] == "other" && cases[i] != "other" {
			return true
		}
		// the following code relies on '<' < '=' < any letter.
		return cmpNumeric(cases[i], cases[j]) == -1
	})
}

var cmpNumeric = collate.New(language.Und, collate.Numeric).CompareString

var lookup = template.Must(template.New("gen").Parse(`
import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p := messageKeyToIndex[key]
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		{{range .Languages}}"{{.}}": &dictionary{index: {{.}}Index, data: {{.}}Data },
		{{end}}
	}
	fallback := language.MustParse("{{.Fallback}}")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

`))
