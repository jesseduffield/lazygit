package commands

import (
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// Dir is a directory containing files
type Dir struct {
	Name    string
	Parrent *Dir
	Files   []*File
	SubDirs []*Dir
}

// Height returns the display height of this dir
func (d *Dir) Height() (height int) {
	if d.SubDirs != nil {
		height++
	}
	height += len(d.Files)
	for _, dir := range d.SubDirs {
		height += dir.Height()
	}
	return
}

// AddFile adds a file to the dir
func (d *Dir) AddFile(f *File) {
	f.InDir = d
	d.Files = append(d.Files, f)
}

// NewDir creates a new Sub Directory for d
func (d *Dir) NewDir(name string) *Dir {
	newDir := NewDir()
	newDir.Name = name
	newDir.Parrent = d
	d.SubDirs = append(d.SubDirs, newDir)
	return newDir
}

// MatchPath matches a path given and returns the file or dir depending on the end point
func (d *Dir) MatchPath(path []int) (*File, *Dir) {
	currentDir := d
	for _, key := range path {
		dirKey := key - len(currentDir.Files)
		if dirKey >= len(currentDir.SubDirs) {
			// Key is out of range, erset to the last entry
			key = len(currentDir.Files) + len(currentDir.SubDirs) - 1
		}

		if key < len(currentDir.Files) {
			// Selected a file
			return currentDir.Files[key], nil
		}

		currentDir = currentDir.SubDirs[dirKey]
	}
	return nil, currentDir
}

// Combine 2 dirs if d only has has 1 subDir and no files
// Instiad of:
//   dir1
//     dir2
//       file
// It will look like:
//   dir1/dir2
//     file
func (d *Dir) Combine() {
	if len(d.Files) == 0 && len(d.SubDirs) == 1 {
		originalName := d.Name
		*d = *d.SubDirs[0]
		d.Name = path.Join(d.Name, originalName)
		d.Combine()
		return
	}
	for _, subDir := range d.SubDirs {
		subDir.Combine()
	}
}

// Render returns a string to render on the screen
func (d *Dir) Render() string {
	return strings.Join(d.RenderAsList(), "\n")
}

// RenderAsList renders dir as a list
func (d *Dir) RenderAsList() []string {
	toReturn := []string{}
	add := func(in ...string) {
		toReturn = append(toReturn, in...)
	}

	for _, file := range d.Files {
		add(path.Base(file.Name))
	}

	for _, dir := range d.SubDirs {
		add(dir.Name + "/")
		list := dir.RenderAsList()
		for i, item := range list {
			list[i] = "  " + item
		}
		add(list...)
	}

	return toReturn
}

// NewDir creates a empty Dir
func NewDir() *Dir {
	return &Dir{
		Files:   make([]*File, 0),
		SubDirs: make([]*Dir, 0),
	}
}

func FilesToTree(log *logrus.Entry, files []*File) *Dir {
	root := NewDir()
	for _, file := range files {
		dir := path.Dir(file.Name)
		currentDir := root
	dirLoop:
		for _, folder := range strings.Split(dir, "/") {
			for _, subdir := range currentDir.SubDirs {
				if subdir.Name == folder {
					currentDir = subdir
					continue dirLoop
				}
			}
			currentDir = currentDir.NewDir(folder)
		}
		currentDir.AddFile(file)
	}
	return root
}

// File : A file from git status
// duplicating this for now
type File struct {
	Name                    string
	HasStagedChanges        bool
	HasUnstagedChanges      bool
	Tracked                 bool
	Deleted                 bool
	HasMergeConflicts       bool
	HasInlineMergeConflicts bool
	DisplayString           string
	Type                    string // one of 'file', 'directory', and 'other'
	ShortStatus             string // e.g. 'AD', ' A', 'M ', '??'
	InDir                   *Dir
}

// GetDisplayStrings returns the display string of a file
func (f *File) GetDisplayStrings(isFocused bool) []string {
	// potentially inefficient to be instantiating these color
	// objects with each render
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	if !f.Tracked && !f.HasStagedChanges {
		return []string{red.Sprint(f.DisplayString)}
	}

	output := green.Sprint(f.DisplayString[0:1])
	output += red.Sprint(f.DisplayString[1:3])
	if f.HasUnstagedChanges {
		output += red.Sprint(f.Name)
	} else {
		output += green.Sprint(f.Name)
	}
	return []string{output}
}
