package getter

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestHttpGetter_impl(t *testing.T) {
	var _ Getter = new(HttpGetter)
}

func TestHttpGetter_header(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	g := new(HttpGetter)
	dst := tempDir(t)
	defer os.RemoveAll(dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/header"

	// Get it!
	if err := g.Get(dst, &u); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestHttpGetter_requestHeader(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	g := new(HttpGetter)
	g.Header = make(http.Header)
	g.Header.Add("X-Foobar", "foobar")
	dst := tempDir(t)
	defer os.RemoveAll(dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/expect-header"
	u.RawQuery = "expected=X-Foobar"

	// Get it!
	if err := g.GetFile(dst, &u); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("err: %s", err)
	}
	assertContents(t, dst, "Hello\n")
}

func TestHttpGetter_meta(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	g := new(HttpGetter)
	dst := tempDir(t)
	defer os.RemoveAll(dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/meta"

	// Get it!
	if err := g.Get(dst, &u); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestHttpGetter_metaSubdir(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	g := new(HttpGetter)
	dst := tempDir(t)
	defer os.RemoveAll(dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/meta-subdir"

	// Get it!
	if err := g.Get(dst, &u); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "sub.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestHttpGetter_metaSubdirGlob(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	g := new(HttpGetter)
	dst := tempDir(t)
	defer os.RemoveAll(dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/meta-subdir-glob"

	// Get it!
	if err := g.Get(dst, &u); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "sub.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestHttpGetter_none(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	g := new(HttpGetter)
	dst := tempDir(t)
	defer os.RemoveAll(dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/none"

	// Get it!
	if err := g.Get(dst, &u); err == nil {
		t.Fatal("should error")
	}
}

func TestHttpGetter_resume(t *testing.T) {
	load := []byte(testHttpMetaStr)
	sha := sha256.New()
	if n, err := sha.Write(load); n != len(load) || err != nil {
		t.Fatalf("sha write failed: %d, %s", n, err)
	}
	checksum := hex.EncodeToString(sha.Sum(nil))
	downloadFrom := len(load) / 2

	ln := testHttpServer(t)
	defer ln.Close()

	dst := tempDir(t)
	defer os.RemoveAll(dst)

	dst = filepath.Join(dst, "..", "range")
	f, err := os.Create(dst)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if n, err := f.Write(load[:downloadFrom]); n != downloadFrom || err != nil {
		t.Fatalf("partial file write failed: %d, %s", n, err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close failed: %s", err)
	}

	u := url.URL{
		Scheme:   "http",
		Host:     ln.Addr().String(),
		Path:     "/range",
		RawQuery: "checksum=" + checksum,
	}
	t.Logf("url: %s", u.String())

	// Finish getting it!
	if err := GetFile(dst, u.String()); err != nil {
		t.Fatalf("finishing download should not error: %v", err)
	}

	b, err := ioutil.ReadFile(dst)
	if err != nil {
		t.Fatalf("readfile failed: %v", err)
	}

	if string(b) != string(load) {
		t.Fatalf("file differs: got:\n%s\n expected:\n%s\n", string(b), string(load))
	}

	// Get it again
	if err := GetFile(dst, u.String()); err != nil {
		t.Fatalf("should not error: %v", err)
	}
}

func TestHttpGetter_file(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	g := new(HttpGetter)
	dst := tempTestFile(t)
	defer os.RemoveAll(filepath.Dir(dst))

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/file"

	// Get it!
	if err := g.GetFile(dst, &u); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("err: %s", err)
	}
	assertContents(t, dst, "Hello\n")
}

func TestHttpGetter_auth(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	g := new(HttpGetter)
	dst := tempDir(t)
	defer os.RemoveAll(dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/meta-auth"
	u.User = url.UserPassword("foo", "bar")

	// Get it!
	if err := g.Get(dst, &u); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestHttpGetter_authNetrc(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	g := new(HttpGetter)
	dst := tempDir(t)
	defer os.RemoveAll(dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/meta"

	// Write the netrc file
	path, closer := tempFileContents(t, fmt.Sprintf(testHttpNetrc, ln.Addr().String()))
	defer closer()
	defer tempEnv(t, "NETRC", path)()

	// Get it!
	if err := g.Get(dst, &u); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// test round tripper that only returns an error
type errRoundTripper struct{}

func (errRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, errors.New("test round tripper")
}

// verify that the default httpClient no longer comes from http.DefaultClient
func TestHttpGetter_cleanhttp(t *testing.T) {
	ln := testHttpServer(t)
	defer ln.Close()

	// break the default http client
	http.DefaultClient.Transport = errRoundTripper{}
	defer func() {
		http.DefaultClient.Transport = http.DefaultTransport
	}()

	g := new(HttpGetter)
	dst := tempDir(t)
	defer os.RemoveAll(dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/header"

	// Get it!
	if err := g.Get(dst, &u); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testHttpServer(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/expect-header", testHttpHandlerExpectHeader)
	mux.HandleFunc("/file", testHttpHandlerFile)
	mux.HandleFunc("/header", testHttpHandlerHeader)
	mux.HandleFunc("/meta", testHttpHandlerMeta)
	mux.HandleFunc("/meta-auth", testHttpHandlerMetaAuth)
	mux.HandleFunc("/meta-subdir", testHttpHandlerMetaSubdir)
	mux.HandleFunc("/meta-subdir-glob", testHttpHandlerMetaSubdirGlob)
	mux.HandleFunc("/range", testHttpHandlerRange)

	var server http.Server
	server.Handler = mux
	go server.Serve(ln)

	return ln
}

func testHttpHandlerExpectHeader(w http.ResponseWriter, r *http.Request) {
	if expected, ok := r.URL.Query()["expected"]; ok {
		if r.Header.Get(expected[0]) != "" {
			w.Write([]byte("Hello\n"))
			return
		}
	}

	w.WriteHeader(400)
}

func testHttpHandlerFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello\n"))
}

func testHttpHandlerHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Terraform-Get", testModuleURL("basic").String())
	w.WriteHeader(200)
}

func testHttpHandlerMeta(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(testHttpMetaStr, testModuleURL("basic").String())))
}

func testHttpHandlerMetaAuth(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(401)
		return
	}

	if user != "foo" || pass != "bar" {
		w.WriteHeader(401)
		return
	}

	w.Write([]byte(fmt.Sprintf(testHttpMetaStr, testModuleURL("basic").String())))
}

func testHttpHandlerMetaSubdir(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(testHttpMetaStr, testModuleURL("basic//subdir").String())))
}

func testHttpHandlerMetaSubdirGlob(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(testHttpMetaStr, testModuleURL("basic//sub*").String())))
}

func testHttpHandlerNone(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(testHttpNoneStr))
}

func testHttpHandlerRange(w http.ResponseWriter, r *http.Request) {
	load := []byte(testHttpMetaStr)
	switch r.Method {
	case "HEAD":
		w.Header().Add("accept-ranges", "bytes")
		w.Header().Add("content-length", strconv.Itoa(len(load)))
	default:
		// request should have header "Range: bytes=0-1023"
		// or                         "Range: bytes=123-"
		rangeHeaderValue := strings.Split(r.Header.Get("Range"), "=")[1]
		rng, _ := strconv.Atoi(strings.Split(rangeHeaderValue, "-")[0])
		if rng < 1 || rng > len(load) {
			http.Error(w, "", http.StatusBadRequest)
		}
		w.Write(load[rng:])
	}
}

const testHttpMetaStr = `
<html>
<head>
<meta name="terraform-get" content="%s">
</head>
</html>
`

const testHttpNoneStr = `
<html>
<head>
</head>
</html>
`

const testHttpNetrc = `
machine %s
login foo
password bar
`
