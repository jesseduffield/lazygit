package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
)

// SplitLines takes a multiline string and splits it on newlines
// currently we are also stripping \r's which may have adverse effects for
// windows users (but no issues have been raised yet)
func SplitLines(multilineString string) []string {
	multilineString = strings.Replace(multilineString, "\r", "", -1)
	if multilineString == "" || multilineString == "\n" {
		return make([]string, 0)
	}
	lines := strings.Split(multilineString, "\n")
	if lines[len(lines)-1] == "" {
		return lines[:len(lines)-1]
	}
	return lines
}

// WithPadding pads a string as much as you want
func WithPadding(str string, padding int) string {
	if padding-len(str) < 0 {
		return str
	}
	return str + strings.Repeat(" ", padding-len(str))
}

// ColoredString takes a string and a colour attribute and returns a colored
// string with that attribute
func ColoredString(str string, colorAttribute color.Attribute) string {
	colour := color.New(colorAttribute)
	return ColoredStringDirect(str, colour)
}

// ColoredStringDirect used for aggregating a few color attributes rather than
// just sending a single one
func ColoredStringDirect(str string, colour *color.Color) string {
	return colour.SprintFunc()(fmt.Sprint(str))
}

// GetCurrentRepoName gets the repo's base name
func GetCurrentRepoName() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return filepath.Base(pwd)
}

// TrimTrailingNewline - Trims the trailing newline
// TODO: replace with `chomp` after refactor
func TrimTrailingNewline(str string) string {
	if strings.HasSuffix(str, "\n") {
		return str[:len(str)-1]
	}
	return str
}

// NormalizeLinefeeds - Removes all Windows and Mac style line feeds
func NormalizeLinefeeds(str string) string {
	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

// GetProjectRoot returns the path to the root of the project. Only to be used
// in testing contexts, as with binaries it's unlikely this path will exist on
// the machine
func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return strings.Split(dir, "lazygit")[0] + "lazygit"
}

// Loader dumps a string to be displayed as a loader
func Loader() string {
	characters := "|/-\\"
	now := time.Now()
	nanos := now.UnixNano()
	index := nanos / 50000000 % int64(len(characters))
	return characters[index : index+1]
}

// ResolvePlaceholderString populates a template with values
func ResolvePlaceholderString(str string, arguments map[string]string) string {
	for key, value := range arguments {
		str = strings.Replace(str, "{{"+key+"}}", value, -1)
	}
	return str
}

// Min returns the minimum of two integers
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

type Displayable interface {
	GetDisplayStrings() []string
}

// RenderList takes a slice of items, confirms they implement the Displayable
// interface, then generates a list of their displaystrings to write to a panel's
// buffer
func RenderList(slice interface{}) (string, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return "", errors.New("RenderList given a non-slice type")
	}

	displayables := make([]Displayable, s.Len())

	for i := 0; i < s.Len(); i++ {
		value, ok := s.Index(i).Interface().(Displayable)
		if !ok {
			return "", errors.New("item does not implement the Displayable interface")
		}
		displayables[i] = value
	}

	return renderDisplayableList(displayables)
}

// renderDisplayableList takes a list of displayable items, obtains their display
// strings via GetDisplayStrings() and then returns a single string containing
// each item's string representation on its own line, with appropriate horizontal
// padding between the item's own strings
func renderDisplayableList(items []Displayable) (string, error) {
	if len(items) == 0 {
		return "", nil
	}

	stringArrays := getDisplayStringArrays(items)

	if !displayArraysAligned(stringArrays) {
		return "", errors.New("Each item must return the same number of strings to display")
	}

	padWidths := getPadWidths(stringArrays)
	paddedDisplayStrings := getPaddedDisplayStrings(stringArrays, padWidths)

	return strings.Join(paddedDisplayStrings, "\n"), nil
}

func getPadWidths(stringArrays [][]string) []int {
	if len(stringArrays[0]) <= 1 {
		return []int{}
	}
	padWidths := make([]int, len(stringArrays[0])-1)
	for i := range padWidths {
		for _, strings := range stringArrays {
			if len(strings[i]) > padWidths[i] {
				padWidths[i] = len(strings[i])
			}
		}
	}
	return padWidths
}

func getPaddedDisplayStrings(stringArrays [][]string, padWidths []int) []string {
	paddedDisplayStrings := make([]string, len(stringArrays))
	for i, stringArray := range stringArrays {
		if len(stringArray) == 0 {
			continue
		}
		for j, padWidth := range padWidths {
			paddedDisplayStrings[i] += WithPadding(stringArray[j], padWidth) + " "
		}
		paddedDisplayStrings[i] += stringArray[len(padWidths)]
	}
	return paddedDisplayStrings
}

// displayArraysAligned returns true if every string array returned from our
// list of displayables has the same length
func displayArraysAligned(stringArrays [][]string) bool {
	for _, strings := range stringArrays {
		if len(strings) != len(stringArrays[0]) {
			return false
		}
	}
	return true
}

func getDisplayStringArrays(displayables []Displayable) [][]string {
	stringArrays := make([][]string, len(displayables))
	for i, item := range displayables {
		stringArrays[i] = item.GetDisplayStrings()
	}
	return stringArrays
}

// IncludesString if the list contains the string
func IncludesString(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
