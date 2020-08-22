package commands

import (
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

const RENAME_SEPARATOR = " -> "

func (f *File) IsRename() bool {
	return strings.Contains(f.Name, RENAME_SEPARATOR)
}

// Names returns an array containing just the filename, or in the case of a rename, the after filename and the before filename
func (f *File) Names() []string {
	return strings.Split(f.Name, RENAME_SEPARATOR)
}

// returns true if the file names are the same or if a a file rename includes the filename of the other
func (f *File) Matches(f2 *File) bool {
	return utils.StringArraysOverlap(f.Names(), f2.Names())
}

// Dir is a directory containing files
type Dir struct {
	Name                    string
	Parrent                 *Dir
	Files                   []*File
	SubDirs                 []*Dir
	ShortStatus             string // e.g. "AD", " A", "M ", "??"
	Untracked               bool
	HasStagedChanges        bool
	HasUnstagedChanges      bool
	Tracked                 bool
	Deleted                 bool
	HasNoStagedChanges      bool
	HasMergeConflicts       bool
	HasInlineMergeConflicts bool
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
		bindTo := *check.bindTo
		switch check.check {
		case ' ':
			// Ignored for now
		case '?':
			switch bindTo {
			case " ":
				*check.bindTo = "?"
			}
		case 'M':
			switch bindTo {
			case " ", "A", "D", "R":
				*check.bindTo = "M"
			}
		case 'A':
			switch bindTo {
			case " ":
				*check.bindTo = "A"
			}
		case 'D':
			switch bindTo {
			case " ":
				*check.bindTo = "D"
			}
		case 'R':
			switch bindTo {
			case " ":
				*check.bindTo = "R"
			}
		case 'C':
			switch bindTo {
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
		lastStatus = MergeGITStatus(lastStatus, current.ShortStatus)

		stagedChange := rune(lastStatus[0])
		unstagedChange := rune(lastStatus[1])

		current.ShortStatus = lastStatus
		current.HasNoStagedChanges = strings.ContainsRune(" U?", stagedChange)
		current.HasStagedChanges = !current.HasNoStagedChanges
		current.HasUnstagedChanges = unstagedChange != ' '
		current.Deleted = unstagedChange == 'D' || stagedChange == 'D'
		current.Untracked = utils.IncludesString([]string{"??", "A ", "AM"}, lastStatus)
		current.Tracked = !current.Untracked
		current.HasMergeConflicts = utils.IncludesString([]string{"DD", "AA", "UU", "AU", "UA", "UD", "DU"}, lastStatus)
		current.HasInlineMergeConflicts = utils.IncludesString([]string{"UU", "AA"}, lastStatus)

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
func (d *Dir) Combine(log *logrus.Entry) *Dir {
	if len(d.SubDirs) == 0 {
		return d
	}

	for {
		if len(d.Files) != 0 || len(d.SubDirs) != 1 || d.Name == "" {
			break
		}
		toMerge := d.SubDirs[0]
		d.Name = path.Join(d.Name, toMerge.Name)
		d.SubDirs = toMerge.SubDirs
		d.Files = toMerge.Files
		for _, dir := range d.SubDirs {
			dir.Parrent = d
		}
		for _, file := range d.Files {
			file.InDir = d
		}
	}

	for _, subDir := range d.SubDirs {
		subDir.Combine(log)
	}

	return d
}

// Render returns a string to render on the screen
func (d *Dir) Render(focusedFile *File, focusedDir *Dir) string {
	return strings.Join(d.RenderAsList(focusedFile, focusedDir), "\n")
}

// AllFiles returns all files in d
func (d *Dir) AllFiles() []*File {
	response := d.Files

	for _, subDir := range d.SubDirs {
		response = append(response, subDir.AllFiles()...)
	}

	return response
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
	return root.Combine(log)
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

// AbsolutePath returns the absolute path of the dir relative to the git repo
func (d *Dir) AbsolutePath() string {
	dir := d
	path := d.Name

	if dir.Parrent == nil {
		return path
	}

	for {
		dir = d.Parrent
		if dir.Parrent == nil {
			// This checks if we don't include the root dir
			break
		}
		path = dir.Name + "/" + path
	}

	return path
}

// GetTreeDisplayString returns the display string of a file for the tree view
func (f *File) GetTreeDisplayString(focused bool) string {
	parts := strings.Split(f.Name, "/")
	name := parts[len(parts)-1]

	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	if !f.Tracked && !f.HasStagedChanges {
		return red.Sprint(f.DisplayString[0:3] + name)
	}

	output := green.Sprint(f.DisplayString[0:1])
	output += red.Sprint(f.DisplayString[1:3])
	if f.HasUnstagedChanges {
		output += red.Sprint(name)
	} else {
		output += green.Sprint(name)
	}
	return output
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

		for _, dir := range parrent.SubDirs {
			if dir == current {
				count += 1
				break
			}
			count += dir.Height()
		}

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
