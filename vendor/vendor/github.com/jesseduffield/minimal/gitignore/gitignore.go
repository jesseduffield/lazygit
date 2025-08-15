// Copyright 2018 iriri. All rights reserved. Use of this source code is
// governed by a BSD-style license which can be found in the LICENSE file.

// Package gitignore can be used to parse .gitignore-style files into lists of
// globs that can be used to test against paths or selectively walk a file
// tree. Gobwas's glob package is used for matching because it is faster than
// using regexp, which is overkill, and supports globstars (**), unlike
// filepath.Match.
package gitignore

import (
	"bufio"
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

type ignoreFile struct {
	globs   []glob.Glob
	abspath []string
}

type IgnoreList struct {
	files []ignoreFile
	cwd   []string
}

func toSplit(path string) []string {
	return strings.Split(filepath.ToSlash(path), "/")
}

func fromSplit(path []string) string {
	return filepath.FromSlash(strings.Join(path, "/"))
}

// New creates a new ignore list.
func New() (IgnoreList, error) {
	cwd, err := filepath.Abs(".")
	if err != nil {
		return IgnoreList{}, err
	}
	files := make([]ignoreFile, 1, 4)
	files[0].globs = make([]glob.Glob, 0, 16)
	return IgnoreList{
		files,
		toSplit(cwd),
	}, nil
}

// From creates a new ignore list and populates the first entry with the
// contents of the specified file.
func From(path string) (IgnoreList, error) {
	ign, err := New()
	if err == nil {
		err = ign.append(path, nil)
	}
	return ign, err
}

// FromGit finds the root directory of the current git repository and creates a
// new ignore list with the contents of all .gitignore files in that git
// repository.
func FromGit() (IgnoreList, error) {
	ign, err := New()
	if err == nil {
		err = ign.AppendGit()
	}
	return ign, err
}

func clean(s string) string {
	i := len(s) - 1
	for ; i >= 0; i-- {
		if s[i] != ' ' || i > 0 && s[i-1] == '\\' {
			return s[:i+1]
		}
	}
	return ""
}

// AppendGlob appends a single glob as a new entry in the ignore list. The root
// (relevant for matching patterns that begin with "/") is assumed to be the
// current working directory.
func (ign *IgnoreList) AppendGlob(s string) error {
	g, err := glob.Compile(clean(s), '/')
	if err == nil {
		ign.files[0].globs = append(ign.files[0].globs, g)
	}
	return err
}

func toRelpath(s string, dir, cwd []string) string {
	if s != "" {
		if s[0] != '/' {
			return s
		}
		if dir == nil || cwd == nil {
			return s[1:]
		}
		dir = append(dir, toSplit(s[1:])...)
	}

	i := 0
	min := len(cwd)
	if len(dir) < min {
		min = len(dir)
	}
	for ; i < min; i++ {
		if dir[i] != cwd[i] {
			break
		}
	}
	if i == min && len(cwd) == len(dir) {
		return "."
	}

	ss := make([]string, (len(cwd)-i)+(len(dir)-i))
	j := 0
	for ; j < len(cwd)-i; j++ {
		ss[j] = ".."
	}
	for k := 0; j < len(ss); j, k = j+1, k+1 {
		ss[j] = dir[i+k]
	}
	return fromSplit(ss)
}

func (ign *IgnoreList) append(path string, dir []string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var ignf *ignoreFile
	if dir != nil {
		ignf = &ign.files[0]
	} else {
		d, err := filepath.Abs(filepath.Dir(path))
		if err != nil {
			return err
		}
		if d != fromSplit(ign.cwd) {
			dir = toSplit(d)
			ignf = &ignoreFile{
				make([]glob.Glob, 0, 16),
				dir,
			}
		} else {
			ignf = &ign.files[0]
		}
	}
	scn := bufio.NewScanner(bufio.NewReader(f))
	for scn.Scan() {
		s := scn.Text()
		if s == "" || s[0] == '#' {
			continue
		}
		g, err := glob.Compile(toRelpath(clean(s), dir, ign.cwd), '/')
		if err != nil {
			return err
		}
		ignf.globs = append(ignf.globs, g)
	}
	ign.files = append(ign.files, *ignf)
	return nil
}

// Append appends the globs in the specified file to the ignore list. Files are
// expected to have the same format as .gitignore files.
func (ign *IgnoreList) Append(path string) error {
	return ign.append(path, nil)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func findGitRoot(cwd []string) (string, error) {
	p := fromSplit(cwd)
	for !exists(p + "/.git") {
		if len(cwd) == 1 {
			return "", errors.New("not in a git repository")
		}
		cwd = cwd[:len(cwd)-1]
		p = fromSplit(cwd)
	}
	return p, nil
}

func (ign *IgnoreList) appendAll(fname, root string) error {
	return filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Base(path) == fname {
				ign.append(path, nil)
			}
			return nil
		})
}

// AppendGit finds the root directory of the current git repository and appends
// the contents of all .gitignore files in that git repository to the ignore
// list.
func (ign *IgnoreList) AppendGit() error {
	gitRoot, err := findGitRoot(ign.cwd)
	if err != nil {
		return err
	}
	if err = ign.AppendGlob(".git"); err != nil {
		return err
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	if gg := filepath.Join(usr.HomeDir, ".gitignore_global"); exists(gg) {
		if err = ign.append(gg, toSplit(gitRoot)); err != nil {
			return err
		}
	}
	return ign.appendAll(".gitignore", gitRoot)
}

func isPrefix(abspath, dir []string) bool {
	if len(abspath) > len(dir) {
		return false
	}
	for i := range abspath {
		if abspath[i] != dir[i] {
			return false
		}
	}
	return true
}

func (ign *IgnoreList) match(path string, info os.FileInfo) bool {
	if path == "." {
		return false
	}
	ss := make([]string, 0, 4)
	base := filepath.Base(path)
	ss = append(ss, path)
	if base != path {
		ss = append(ss, base)
	} else {
		ss = append(ss, "./"+path)
	}
	if info != nil && info.IsDir() {
		ss = append(ss, path+"/")
		if base != path {
			ss = append(ss, base+"/")
		} else {
			ss = append(ss, "./"+path+"/")
		}
	}

	d, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return false
	}
	dir := toSplit(d)
	for _, f := range ign.files {
		if isPrefix(f.abspath, dir) || len(f.abspath) == 0 {
			for _, g := range f.globs {
				for _, s := range ss {
					if g.Match(s) {
						return true
					}
				}
			}
		}
	}
	return false
}

// Match returns whether any of the globs in the ignore list match the
// specified path. Uses the same matching rules as .gitignore files.
func (ign *IgnoreList) Match(path string) bool {
	return ign.match(path, nil)
}

// Walk walks the file tree with the specified root and calls fn on each file
// or directory. Files and directories that match any of the globs in the
// ignore list are skipped.
func (ign *IgnoreList) Walk(root string, fn filepath.WalkFunc) error {
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	return filepath.Walk(
		toRelpath("", toSplit(abs), ign.cwd),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if ign.match(path, info) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return err
			}
			return fn(path, info, err)
		})
}
