// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message_test

import (
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func Example_http() {
	// languages supported by this service:
	matcher := language.NewMatcher(message.DefaultCatalog.Languages())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lang, _ := r.Cookie("lang")
		accept := r.Header.Get("Accept-Language")
		fallback := "en"
		tag, _ := language.MatchStrings(matcher, lang.String(), accept, fallback)

		p := message.NewPrinter(tag)

		p.Fprintln(w, "User language is", tag)
	})
}

func ExamplePrinter_numbers() {
	for _, lang := range []string{"en", "de", "de-CH", "fr", "bn"} {
		p := message.NewPrinter(language.Make(lang))
		p.Printf("%-6s %g\n", lang, 123456.78)
	}

	// Output:
	// en     123,456.78
	// de     123.456,78
	// de-CH  123’456.78
	// fr     123 456,78
	// bn     ১,২৩,৪৫৬.৭৮
}
