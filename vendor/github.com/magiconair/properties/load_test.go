// Copyright 2018 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestEncoding(t *testing.T) {
	if got, want := utf8Default, Encoding(0); got != want {
		t.Fatalf("got encoding %d want %d", got, want)
	}
	if got, want := UTF8, Encoding(1); got != want {
		t.Fatalf("got encoding %d want %d", got, want)
	}
	if got, want := ISO_8859_1, Encoding(2); got != want {
		t.Fatalf("got encoding %d want %d", got, want)
	}
}

func TestLoadFailsWithNotExistingFile(t *testing.T) {
	_, err := LoadFile("doesnotexist.properties", ISO_8859_1)
	assert.Equal(t, err != nil, true, "")
	assert.Matches(t, err.Error(), "open.*no such file or directory")
}

func TestLoadFilesFailsOnNotExistingFile(t *testing.T) {
	_, err := LoadFile("doesnotexist.properties", ISO_8859_1)
	assert.Equal(t, err != nil, true, "")
	assert.Matches(t, err.Error(), "open.*no such file or directory")
}

func TestLoadFilesDoesNotFailOnNotExistingFileAndIgnoreMissing(t *testing.T) {
	p, err := LoadFiles([]string{"doesnotexist.properties"}, ISO_8859_1, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, p.Len(), 0)
}

func TestLoadString(t *testing.T) {
	x := "key=äüö"
	p1 := MustLoadString(x)
	p2 := must(Load([]byte(x), UTF8))
	assert.Equal(t, p1, p2)
}

func TestLoadMap(t *testing.T) {
	// LoadMap does not guarantee the same import order
	// of keys every time since map access is randomized.
	// Therefore, we need to compare the generated maps.
	m := map[string]string{"key": "value", "abc": "def"}
	assert.Equal(t, LoadMap(m).Map(), m)
}

func TestLoadFile(t *testing.T) {
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	filename := tf.makeFile("key=value")
	p := MustLoadFile(filename, ISO_8859_1)

	assert.Equal(t, p.Len(), 1)
	assertKeyValues(t, "", p, "key", "value")
}

func TestLoadFiles(t *testing.T) {
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	filename := tf.makeFile("key=value")
	filename2 := tf.makeFile("key2=value2")
	p := MustLoadFiles([]string{filename, filename2}, ISO_8859_1, false)
	assertKeyValues(t, "", p, "key", "value", "key2", "value2")
}

func TestLoadExpandedFile(t *testing.T) {
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	if err := os.Setenv("_VARX", "some-value"); err != nil {
		t.Fatal(err)
	}
	filename := tf.makeFilePrefix(os.Getenv("_VARX"), "key=value")
	filename = strings.Replace(filename, os.Getenv("_VARX"), "${_VARX}", -1)
	p := MustLoadFile(filename, ISO_8859_1)
	assertKeyValues(t, "", p, "key", "value")
}

func TestLoadFilesAndIgnoreMissing(t *testing.T) {
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	filename := tf.makeFile("key=value")
	filename2 := tf.makeFile("key2=value2")
	p := MustLoadFiles([]string{filename, filename + "foo", filename2, filename2 + "foo"}, ISO_8859_1, true)
	assertKeyValues(t, "", p, "key", "value", "key2", "value2")
}

func TestLoadURL(t *testing.T) {
	srv := testServer()
	defer srv.Close()
	p := MustLoadURL(srv.URL + "/a")
	assertKeyValues(t, "", p, "key", "value")
}

func TestLoadURLs(t *testing.T) {
	srv := testServer()
	defer srv.Close()
	p := MustLoadURLs([]string{srv.URL + "/a", srv.URL + "/b"}, false)
	assertKeyValues(t, "", p, "key", "value", "key2", "value2")
}

func TestLoadURLsAndFailMissing(t *testing.T) {
	srv := testServer()
	defer srv.Close()
	p, err := LoadURLs([]string{srv.URL + "/a", srv.URL + "/c"}, false)
	assert.Equal(t, p, (*Properties)(nil))
	assert.Matches(t, err.Error(), ".*returned 404.*")
}

func TestLoadURLsAndIgnoreMissing(t *testing.T) {
	srv := testServer()
	defer srv.Close()
	p := MustLoadURLs([]string{srv.URL + "/a", srv.URL + "/b", srv.URL + "/c"}, true)
	assertKeyValues(t, "", p, "key", "value", "key2", "value2")
}

func TestLoadURLEncoding(t *testing.T) {
	srv := testServer()
	defer srv.Close()

	uris := []string{"/none", "/utf8", "/plain", "/latin1", "/iso88591"}
	for i, uri := range uris {
		p := MustLoadURL(srv.URL + uri)
		assert.Equal(t, p.GetString("key", ""), "äöü", fmt.Sprintf("%d", i))
	}
}

func TestLoadURLFailInvalidEncoding(t *testing.T) {
	srv := testServer()
	defer srv.Close()

	p, err := LoadURL(srv.URL + "/json")
	assert.Equal(t, p, (*Properties)(nil))
	assert.Matches(t, err.Error(), ".*invalid content type.*")
}

func TestLoadAll(t *testing.T) {
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	filename := tf.makeFile("key=value")
	filename2 := tf.makeFile("key2=value3")
	filename3 := tf.makeFile("key=value4")
	srv := testServer()
	defer srv.Close()
	p := MustLoadAll([]string{filename, filename2, srv.URL + "/a", srv.URL + "/b", filename3}, UTF8, false)
	assertKeyValues(t, "", p, "key", "value4", "key2", "value2")
}

type tempFiles []string

func (tf *tempFiles) removeAll() {
	for _, path := range *tf {
		err := os.Remove(path)
		if err != nil {
			fmt.Printf("os.Remove: %v", err)
		}
	}
}

func (tf *tempFiles) makeFile(data string) string {
	return tf.makeFilePrefix("properties", data)
}

func (tf *tempFiles) makeFilePrefix(prefix, data string) string {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		panic("ioutil.TempFile: " + err.Error())
	}

	// remember the temp file so that we can remove it later
	*tf = append(*tf, f.Name())

	n, err := fmt.Fprint(f, data)
	if err != nil {
		panic("fmt.Fprintln: " + err.Error())
	}
	if n != len(data) {
		panic(fmt.Sprintf("Data size mismatch. expected=%d wrote=%d\n", len(data), n))
	}

	err = f.Close()
	if err != nil {
		panic("f.Close: " + err.Error())
	}

	return f.Name()
}

func testServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		send := func(data []byte, contentType string) {
			w.Header().Set("Content-Type", contentType)
			if _, err := w.Write(data); err != nil {
				panic(err)
			}
		}

		utf8 := []byte("key=äöü")
		iso88591 := []byte{0x6b, 0x65, 0x79, 0x3d, 0xe4, 0xf6, 0xfc} // key=äöü

		switch r.RequestURI {
		case "/a":
			send([]byte("key=value"), "")
		case "/b":
			send([]byte("key2=value2"), "")
		case "/none":
			send(utf8, "")
		case "/utf8":
			send(utf8, "text/plain; charset=utf-8")
		case "/json":
			send(utf8, "application/json; charset=utf-8")
		case "/plain":
			send(iso88591, "text/plain")
		case "/latin1":
			send(iso88591, "text/plain; charset=latin1")
		case "/iso88591":
			send(iso88591, "text/plain; charset=iso-8859-1")
		default:
			w.WriteHeader(404)
		}
	}))
}
