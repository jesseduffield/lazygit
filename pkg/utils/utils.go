package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
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
	uncoloredStr := Decolorise(str)
	if padding < len(uncoloredStr) {
		return str
	}
	return str + strings.Repeat(" ", padding-len(uncoloredStr))
}

// ColoredString takes a string and a colour attribute and returns a colored
// string with that attribute
func ColoredString(str string, colorAttributes ...color.Attribute) string {
	colour := color.New(colorAttributes...)
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
		str = strings.Replace(str, "{{."+key+"}}", value, -1)
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

func RenderDisplayStrings(displayStringsArr [][]string) string {
	padWidths := getPadWidths(displayStringsArr)
	paddedDisplayStrings := getPaddedDisplayStrings(displayStringsArr, padWidths)

	return strings.Join(paddedDisplayStrings, "\n")
}

// Decolorise strips a string of color
func Decolorise(str string) string {
	re := regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]`)
	return re.ReplaceAllString(str, "")
}

func getPadWidths(stringArrays [][]string) []int {
	maxWidth := 0
	for _, stringArray := range stringArrays {
		if len(stringArray) > maxWidth {
			maxWidth = len(stringArray)
		}
	}
	if maxWidth-1 < 0 {
		return []int{}
	}
	padWidths := make([]int, maxWidth-1)
	for i := range padWidths {
		for _, strings := range stringArrays {
			uncoloredString := Decolorise(strings[i])
			if len(uncoloredString) > padWidths[i] {
				padWidths[i] = len(uncoloredString)
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
			if len(stringArray)-1 < j {
				continue
			}
			paddedDisplayStrings[i] += WithPadding(stringArray[j], padWidth) + " "
		}
		if len(stringArray)-1 < len(padWidths) {
			continue
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

// IncludesString if the list contains the string
func IncludesString(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IncludesInt if the list contains the Int
func IncludesInt(list []int, a int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// NextIndex returns the index of the element that comes after the given number
func NextIndex(numbers []int, currentNumber int) int {
	for index, number := range numbers {
		if number > currentNumber {
			return index
		}
	}
	return len(numbers) - 1
}

// PrevIndex returns the index that comes before the given number, cycling if we reach the end
func PrevIndex(numbers []int, currentNumber int) int {
	end := len(numbers) - 1
	for i := end; i >= 0; i-- {
		if numbers[i] < currentNumber {
			return i
		}
	}
	return 0
}

func AsJson(i interface{}) string {
	bytes, _ := json.MarshalIndent(i, "", "    ")
	return string(bytes)
}

// UnionInt returns the union of two int arrays
func UnionInt(a, b []int) []int {
	m := make(map[int]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			// this does not mutate the original a slice
			// though it does mutate the backing array I believe
			// but that doesn't matter because if you later want to append to the
			// original a it must see that the backing array has been changed
			// and create a new one
			a = append(a, item)
		}
	}
	return a
}

// DifferenceInt returns the difference of two int arrays
func DifferenceInt(a, b []int) []int {
	result := []int{}
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			result = append(result, item)
		}
	}
	return result
}

// used to keep a number n between 0 and max, allowing for wraparounds
func ModuloWithWrap(n, max int) int {
	if n >= max {
		return n % max
	} else if n < 0 {
		return max + n
	} else {
		return n
	}
}

// NextIntInCycle returns the next int in a slice, returning to the first index if we've reached the end
func NextIntInCycle(sl []int, current int) int {
	for i, val := range sl {
		if val == current {
			if i == len(sl)-1 {
				return sl[0]
			}
			return sl[i+1]
		}
	}
	return sl[0]
}

// PrevIntInCycle returns the prev int in a slice, returning to the first index if we've reached the end
func PrevIntInCycle(sl []int, current int) int {
	for i, val := range sl {
		if val == current {
			if i > 0 {
				return sl[i-1]
			}
			return sl[len(sl)-1]
		}
	}
	return sl[len(sl)-1]
}

// TruncateWithEllipsis returns a string, truncated to a certain length, with an ellipsis
func TruncateWithEllipsis(str string, limit int) string {
	if limit == 1 && len(str) > 1 {
		return "."
	}

	if limit == 2 && len(str) > 2 {
		return ".."
	}

	ellipsis := "..."
	if len(str) <= limit {
		return str
	}

	remainingLength := limit - len(ellipsis)
	return str[0:remainingLength] + "..."
}

func FindStringSubmatch(str string, regexpStr string) (bool, []string) {
	re := regexp.MustCompile(regexpStr)
	match := re.FindStringSubmatch(str)
	return len(match) > 0, match
}

func StringArraysOverlap(strArrA []string, strArrB []string) bool {
	for _, first := range strArrA {
		for _, second := range strArrB {
			if first == second {
				return true
			}
		}
	}

	return false
}

func MustConvertToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func ResolveTemplate(templateStr string, object interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, object); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Safe will close termbox if a panic occurs so that we don't end up in a malformed
// terminal state
func Safe(f func()) {
	panicking := true
	defer func() {
		if panicking && gocui.Screen != nil {
			gocui.Screen.Fini()
		}
	}()

	f()

	panicking = false
}
