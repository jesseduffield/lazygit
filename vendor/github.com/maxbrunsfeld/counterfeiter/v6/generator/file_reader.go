package generator

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileReader interface {
	Get(cwd, path string) (content string, err error)
}

type Opener func(string) (io.ReadCloser, error)

var (
	defaultOpen Opener = func(p string) (io.ReadCloser, error) { return os.Open(p) }
)

func (open Opener) readString(path string) (string, error) {
	if open == nil {
		open = defaultOpen
	}

	f, err := open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type SimpleFileReader struct {
	Open Opener
}

var _ FileReader = &SimpleFileReader{}

func (r *SimpleFileReader) Get(cwd, path string) (string, error) {
	if path == "" {
		return "", nil
	}

	p := normalisePath(cwd, path)
	return r.Open.readString(p)
}

type CachedFileReader struct {
	Open  Opener
	cache map[string]string
}

var _ FileReader = &CachedFileReader{}

func (r *CachedFileReader) Get(cwd, path string) (string, error) {
	if path == "" {
		return "", nil
	}

	p := normalisePath(cwd, path)

	if s, ok := r.cache[p]; ok {
		return s, nil
	}

	s, err := r.Open.readString(p)
	if err != nil {
		return "", err
	}

	if r.cache == nil {
		r.cache = map[string]string{}
	}
	r.cache[p] = s
	return s, nil
}

func normalisePath(cwd, path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	return filepath.Clean(path)
}
