package commands

import (
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

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

// Dir is a directory containing files
type Dir struct {
	Name        string
	Parrent     *Dir
	Files       []*File
	SubDirs     []*Dir
	ShortStatus string // e.g. 'AD', ' A', 'M ', '??'
}

// MergeGITStatus merges 2 git status together
func MergeGITStatus(a, b string) string {
	toFormat := []*string{&a, &b}
	for _, item := range toFormat {
		switch len(*item) {
		case 0:
			*item += "  "
		case 1:
			*item += " "
		case 2:
		default:
			*item = string((*item)[:2])
		}
	}

	a = strings.ToUpper(a)
	b = strings.ToUpper(b)

	Cx, Cy := " ", " "

	toCheck := []struct {
		check  rune
		bindTo *string
	}{
		{rune(a[0]), &Cx},
		{rune(a[1]), &Cy},
		{rune(b[0]), &Cx},
		{rune(b[1]), &Cy},
	}
	for _, check := range toCheck {
		switch check.check {
		case ' ':
			// Ignored for now
		case '?':
			switch *check.bindTo {
			case " ":
				*check.bindTo = "?"
			}
		case 'M':
			switch *check.bindTo {
			case " ", "A", "D", "R":
				*check.bindTo = "M"
			}
		case 'A':
			switch *check.bindTo {
			case " ":
				*check.bindTo = "A"
			}
		case 'D':
			switch *check.bindTo {
			case " ":
				*check.bindTo = "D"
			}
		case 'R':
			switch *check.bindTo {
			case " ":
				*check.bindTo = "R"
			}
		case 'C':
			switch *check.bindTo {
			case " ":
				*check.bindTo = "C"
			}
		case 'U':
			*check.bindTo = "U"
		}
	}

	return Cx + Cy
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

	lastStatus := f.ShortStatus
	current := d
	for {
		current.ShortStatus = MergeGITStatus(lastStatus, current.ShortStatus)
		lastStatus = current.ShortStatus

		current = current.Parrent
		if current == nil {
			break
		}
	}

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
		if key < len(currentDir.Files) {
			return currentDir.Files[key], nil
		}

		key -= len(currentDir.Files)
		if len(currentDir.SubDirs) == 0 {
			if len(currentDir.Files) > 0 {
				return currentDir.Files[len(currentDir.Files)-1], nil
			}
			return nil, currentDir
		}
		if key >= len(currentDir.SubDirs) {
			key = len(currentDir.SubDirs) - 1
		}
		currentDir = currentDir.SubDirs[key]
	}
	if currentDir.Parrent == nil {
		// We are at the root here, this issn't a line on the screen
		return nil, nil
	}
	return nil, currentDir
}

// Combine 2 dirs if d only has has 1 subDir and no files
// Instiad of:
//   dir1/
//     dir2/
//       file
// It will look like:
//   dir1/dir2/
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
func (d *Dir) Render(focusedFile *File, focusedDir *Dir) string {
	return strings.Join(d.RenderAsList(focusedFile, focusedDir), "\n")
}

// RenderAsList renders dir as a list
func (d *Dir) RenderAsList(focusedFile *File, focusedDir *Dir) []string {
	toReturn := []string{}
	add := func(in ...string) {
		toReturn = append(toReturn, in...)
	}

	for _, file := range d.Files {
		add(file.GetTreeDisplayString(focusedFile == file))
	}

	for _, dir := range d.SubDirs {
		add(dir.GetTreeDisplayString(focusedDir == dir))
		list := dir.RenderAsList(focusedFile, focusedDir)
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

// FilesToTree changes a list of files into a dir
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

// GetTreeDisplayString returns the display string of a dir for the tree view
func (d *Dir) GetTreeDisplayString(focused bool) string {
	dirName := d.Name + "/"

	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)

	x, y := string(d.ShortStatus[0]), string(d.ShortStatus[1])

	xColor := green.Sprint(x)
	if x == "?" {
		xColor = red.Sprint(x)
	}
	yColor := red.Sprint(y)

	output := xColor + yColor + " "
	if y != " " {
		output += red.Sprint(dirName)
	} else {
		output += green.Sprint(dirName)
	}
	return output
}

// GetTreeDisplayString returns the display string of a file for the tree view
func (f *File) GetTreeDisplayString(focused bool) string {
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	if !f.Tracked && !f.HasStagedChanges {
		return red.Sprint(f.DisplayString)
	}

	output := green.Sprint(f.DisplayString[0:1])
	output += red.Sprint(f.DisplayString[1:3])
	name := path.Base(f.Name)
	if f.HasUnstagedChanges {
		output += red.Sprint(name)
	} else {
		output += green.Sprint(name)
	}
	return output
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

// GetY returns the dir it's y position
func (d *Dir) GetY() int {
	count := -1
	current := d
	parrent := d.Parrent
	for {
		if parrent == nil {
			break
		}
		for i, dir := range parrent.SubDirs {
			if dir == current {
				count += i
				break
			}
			count += len(dir.SubDirs) + len(dir.Files)
		}
		count += len(parrent.Files)
		count += 1

		current = parrent
		parrent = current.Parrent
	}
	return count
}

// GetY returns the file it's y position
func (f *File) GetY() int {
	dir := f.InDir
	count := dir.GetY()
	for i, file := range dir.Files {
		if file == f {
			count += i + 1
			break
		}
	}
	return count
}
