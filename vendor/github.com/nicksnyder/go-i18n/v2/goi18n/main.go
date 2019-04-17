// Command goi18n manages message files used by the i18n package.
//
//     go get -u github.com/nicksnyder/go-i18n/v2/goi18n
//     goi18n -help
//
// Use `goi18n extract` to create a message file that contains the messages defined in your Go source files.
//     # en.toml
//     [PersonCats]
//     description = "The number of cats a person has"
//     one = "{{.Name}} has {{.Count}} cat."
//     other = "{{.Name}} has {{.Count}} cats."
//
// Use `goi18n merge` to create message files for translation.
//     # translate.es.toml
//     [PersonCats]
//     description = "The number of cats a person has"
//     hash = "sha1-f937a0e05e19bfe6cd70937c980eaf1f9832f091"
//     one = "{{.Name}} has {{.Count}} cat."
//     other = "{{.Name}} has {{.Count}} cats."
//
// Use `goi18n merge` to merge translated message files with your existing message files.
//     # active.es.toml
//     [PersonCats]
//     description = "The number of cats a person has"
//     hash = "sha1-f937a0e05e19bfe6cd70937c980eaf1f9832f091"
//     one = "{{.Name}} tiene {{.Count}} gato."
//     other = "{{.Name}} tiene {{.Count}} gatos."
//
// Load the active messages into your bundle.
//     bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
//     bundle.MustLoadMessageFile("active.es.toml")
package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/text/language"
)

func mainUsage() {
	fmt.Fprintf(os.Stderr, `goi18n (v2) is a tool for managing message translations.

Usage:

	goi18n command [arguments]

The commands are:

	merge		merge message files
	extract		extract messages from Go files

Workflow:

	Use 'goi18n extract' to create a message file that contains the messages defined in your Go source files.

		# en.toml
		[PersonCats]
		description = "The number of cats a person has"
		one = "{{.Name}} has {{.Count}} cat."
		other = "{{.Name}} has {{.Count}} cats."

	Use 'goi18n merge' to create message files for translation.

		# translate.es.toml
		[PersonCats]
		description = "The number of cats a person has"
		hash = "sha1-f937a0e05e19bfe6cd70937c980eaf1f9832f091"
		one = "{{.Name}} has {{.Count}} cat."
		other = "{{.Name}} has {{.Count}} cats."

	Use 'goi18n merge' to merge translated message files with your existing message files.

		# active.es.toml
		[PersonCats]
		description = "The number of cats a person has"
		hash = "sha1-f937a0e05e19bfe6cd70937c980eaf1f9832f091"
		one = "{{.Name}} tiene {{.Count}} gato."
		other = "{{.Name}} tiene {{.Count}} gatos."

	Load the active messages into your bundle.

		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		bundle.MustLoadMessageFile("active.es.toml")
`)
}

type command interface {
	name() string
	parse(arguments []string)
	execute() error
}

func main() {
	os.Exit(testableMain(os.Args[1:]))
}

func testableMain(args []string) int {
	flags := flag.NewFlagSet("goi18n", flag.ContinueOnError)
	flags.Usage = mainUsage
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 2
		}
		return 1
	}
	if flags.NArg() == 0 {
		mainUsage()
		return 2
	}
	commands := []command{
		&mergeCommand{},
		&extractCommand{},
	}
	cmdName := flags.Arg(0)
	for _, cmd := range commands {
		if cmd.name() == cmdName {
			cmd.parse(flags.Args()[1:])
			if err := cmd.execute(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
			return 0
		}
	}
	fmt.Fprintf(os.Stderr, "goi18n: unknown subcommand %s\n", cmdName)
	return 1
}

type languageTag language.Tag

func (lt languageTag) String() string {
	return lt.Tag().String()
}

func (lt *languageTag) Set(value string) error {
	t, err := language.Parse(value)
	if err != nil {
		return err
	}
	*lt = languageTag(t)
	return nil
}

func (lt languageTag) Tag() language.Tag {
	tag := language.Tag(lt)
	if tag.IsRoot() {
		return language.English
	}
	return tag
}
