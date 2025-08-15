//go:build generate

// This program generates a property file in Go file from Unicode Character
// Database auxiliary data files. The command line arguments are as follows:
//
//  1. The name of the Unicode data file (just the filename, without extension).
//     Can be "-" (to skip) if the emoji flag is included.
//  2. The name of the locally generated Go file.
//  3. The name of the slice mapping code points to properties.
//  4. The name of the generator, for logging purposes.
//  5. (Optional) Flags, comma-separated. The following flags are available:
//     - "emojis=<property>": include the specified emoji properties (e.g.
//     "Extended_Pictographic").
//     - "gencat": include general category properties.
//
//go:generate go run gen_properties.go auxiliary/GraphemeBreakProperty graphemeproperties.go graphemeCodePoints graphemes emojis=Extended_Pictographic
//go:generate go run gen_properties.go auxiliary/WordBreakProperty wordproperties.go workBreakCodePoints words emojis=Extended_Pictographic
//go:generate go run gen_properties.go auxiliary/SentenceBreakProperty sentenceproperties.go sentenceBreakCodePoints sentences
//go:generate go run gen_properties.go LineBreak lineproperties.go lineBreakCodePoints lines gencat
//go:generate go run gen_properties.go EastAsianWidth eastasianwidth.go eastAsianWidth eastasianwidth
//go:generate go run gen_properties.go - emojipresentation.go emojiPresentation emojipresentation emojis=Emoji_Presentation
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// We want to test against a specific version rather than the latest. When the
// package is upgraded to a new version, change these to generate new tests.
const (
	propertyURL = `https://www.unicode.org/Public/15.0.0/ucd/%s.txt`
	emojiURL    = `https://unicode.org/Public/15.0.0/ucd/emoji/emoji-data.txt`
)

// The regular expression for a line containing a code point range property.
var propertyPattern = regexp.MustCompile(`^([0-9A-F]{4,6})(\.\.([0-9A-F]{4,6}))?\s*;\s*([A-Za-z0-9_]+)\s*#\s(.+)$`)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Not enough arguments, see code for details")
		os.Exit(1)
	}

	log.SetPrefix("gen_properties (" + os.Args[4] + "): ")
	log.SetFlags(0)

	// Parse flags.
	flags := make(map[string]string)
	if len(os.Args) >= 6 {
		for _, flag := range strings.Split(os.Args[5], ",") {
			flagFields := strings.Split(flag, "=")
			if len(flagFields) == 1 {
				flags[flagFields[0]] = "yes"
			} else {
				flags[flagFields[0]] = flagFields[1]
			}
		}
	}

	// Parse the text file and generate Go source code from it.
	_, includeGeneralCategory := flags["gencat"]
	var mainURL string
	if os.Args[1] != "-" {
		mainURL = fmt.Sprintf(propertyURL, os.Args[1])
	}
	src, err := parse(mainURL, flags["emojis"], includeGeneralCategory)
	if err != nil {
		log.Fatal(err)
	}

	// Format the Go code.
	formatted, err := format.Source([]byte(src))
	if err != nil {
		log.Fatal("gofmt:", err)
	}

	// Save it to the (local) target file.
	log.Print("Writing to ", os.Args[2])
	if err := ioutil.WriteFile(os.Args[2], formatted, 0644); err != nil {
		log.Fatal(err)
	}
}

// parse parses the Unicode Properties text files located at the given URLs and
// returns their equivalent Go source code to be used in the uniseg package. If
// "emojiProperty" is not an empty string, emoji code points for that emoji
// property (e.g. "Extended_Pictographic") will be included. In those cases, you
// may pass an empty "propertyURL" to skip parsing the main properties file. If
// "includeGeneralCategory" is true, the Unicode General Category property will
// be extracted from the comments and included in the output.
func parse(propertyURL, emojiProperty string, includeGeneralCategory bool) (string, error) {
	if propertyURL == "" && emojiProperty == "" {
		return "", errors.New("no properties to parse")
	}

	// Temporary buffer to hold properties.
	var properties [][4]string

	// Open the first URL.
	if propertyURL != "" {
		log.Printf("Parsing %s", propertyURL)
		res, err := http.Get(propertyURL)
		if err != nil {
			return "", err
		}
		in1 := res.Body
		defer in1.Close()

		// Parse it.
		scanner := bufio.NewScanner(in1)
		num := 0
		for scanner.Scan() {
			num++
			line := strings.TrimSpace(scanner.Text())

			// Skip comments and empty lines.
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}

			// Everything else must be a code point range, a property and a comment.
			from, to, property, comment, err := parseProperty(line)
			if err != nil {
				return "", fmt.Errorf("%s line %d: %v", os.Args[4], num, err)
			}
			properties = append(properties, [4]string{from, to, property, comment})
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	}

	// Open the second URL.
	if emojiProperty != "" {
		log.Printf("Parsing %s", emojiURL)
		res, err := http.Get(emojiURL)
		if err != nil {
			return "", err
		}
		in2 := res.Body
		defer in2.Close()

		// Parse it.
		scanner := bufio.NewScanner(in2)
		num := 0
		for scanner.Scan() {
			num++
			line := scanner.Text()

			// Skip comments, empty lines, and everything not containing
			// "Extended_Pictographic".
			if strings.HasPrefix(line, "#") || line == "" || !strings.Contains(line, emojiProperty) {
				continue
			}

			// Everything else must be a code point range, a property and a comment.
			from, to, property, comment, err := parseProperty(line)
			if err != nil {
				return "", fmt.Errorf("emojis line %d: %v", num, err)
			}
			properties = append(properties, [4]string{from, to, property, comment})
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	}

	// Avoid overflow during binary search.
	if len(properties) >= 1<<31 {
		return "", errors.New("too many properties")
	}

	// Sort properties.
	sort.Slice(properties, func(i, j int) bool {
		left, _ := strconv.ParseUint(properties[i][0], 16, 64)
		right, _ := strconv.ParseUint(properties[j][0], 16, 64)
		return left < right
	})

	// Header.
	var (
		buf          bytes.Buffer
		emojiComment string
	)
	columns := 3
	if includeGeneralCategory {
		columns = 4
	}
	if emojiURL != "" {
		emojiComment = `
// and
// ` + emojiURL + `
// ("Extended_Pictographic" only)`
	}
	buf.WriteString(`// Code generated via go generate from gen_properties.go. DO NOT EDIT.

package uniseg

// ` + os.Args[3] + ` are taken from
// ` + propertyURL + emojiComment + `
// on ` + time.Now().Format("January 2, 2006") + `. See https://www.unicode.org/license.html for the Unicode
// license agreement.
var ` + os.Args[3] + ` = [][` + strconv.Itoa(columns) + `]int{
	`)

	// Properties.
	for _, prop := range properties {
		if includeGeneralCategory {
			generalCategory := "gc" + prop[3][:2]
			if generalCategory == "gcL&" {
				generalCategory = "gcLC"
			}
			prop[3] = prop[3][3:]
			fmt.Fprintf(&buf, "{0x%s,0x%s,%s,%s}, // %s\n", prop[0], prop[1], translateProperty("pr", prop[2]), generalCategory, prop[3])
		} else {
			fmt.Fprintf(&buf, "{0x%s,0x%s,%s}, // %s\n", prop[0], prop[1], translateProperty("pr", prop[2]), prop[3])
		}
	}

	// Tail.
	buf.WriteString("}")

	return buf.String(), nil
}

// parseProperty parses a line of the Unicode properties text file containing a
// property for a code point range and returns it along with its comment.
func parseProperty(line string) (from, to, property, comment string, err error) {
	fields := propertyPattern.FindStringSubmatch(line)
	if fields == nil {
		err = errors.New("no property found")
		return
	}
	from = fields[1]
	to = fields[3]
	if to == "" {
		to = from
	}
	property = fields[4]
	comment = fields[5]
	return
}

// translateProperty translates a property name as used in the Unicode data file
// to a variable used in the Go code.
func translateProperty(prefix, property string) string {
	return prefix + strings.ReplaceAll(property, "_", "")
}
