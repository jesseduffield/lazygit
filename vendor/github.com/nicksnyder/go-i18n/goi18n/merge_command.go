package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sort"

	"gopkg.in/yaml.v2"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	toml "github.com/pelletier/go-toml"
)

type mergeCommand struct {
	translationFiles []string
	sourceLanguage   string
	outdir           string
	format           string
	flat             bool
}

func (mc *mergeCommand) execute() error {
	if len(mc.translationFiles) < 1 {
		return fmt.Errorf("need at least one translation file to parse")
	}

	if lang := language.Parse(mc.sourceLanguage); lang == nil {
		return fmt.Errorf("invalid source locale: %s", mc.sourceLanguage)
	}

	bundle := bundle.New()
	for _, tf := range mc.translationFiles {
		if err := bundle.LoadTranslationFile(tf); err != nil {
			return fmt.Errorf("failed to load translation file %s: %s\n", tf, err)
		}
	}

	translations := bundle.Translations()
	sourceLanguageTag := language.NormalizeTag(mc.sourceLanguage)
	sourceTranslations := translations[sourceLanguageTag]
	if sourceTranslations == nil {
		return fmt.Errorf("no translations found for source locale %s", sourceLanguageTag)
	}
	for translationID, src := range sourceTranslations {
		for _, localeTranslations := range translations {
			if dst := localeTranslations[translationID]; dst == nil || reflect.TypeOf(src) != reflect.TypeOf(dst) {
				localeTranslations[translationID] = src.UntranslatedCopy()
			}
		}
	}

	for localeID, localeTranslations := range translations {
		lang := language.MustParse(localeID)[0]
		all := filter(localeTranslations, func(t translation.Translation) translation.Translation {
			return t.Normalize(lang)
		})
		if err := mc.writeFile("all", all, localeID); err != nil {
			return err
		}

		untranslated := filter(localeTranslations, func(t translation.Translation) translation.Translation {
			if t.Incomplete(lang) {
				return t.Normalize(lang).Backfill(sourceTranslations[t.ID()])
			}
			return nil
		})
		if err := mc.writeFile("untranslated", untranslated, localeID); err != nil {
			return err
		}
	}
	return nil
}

func (mc *mergeCommand) parse(arguments []string) {
	flags := flag.NewFlagSet("merge", flag.ExitOnError)
	flags.Usage = usageMerge

	sourceLanguage := flags.String("sourceLanguage", "en-us", "")
	outdir := flags.String("outdir", ".", "")
	format := flags.String("format", "json", "")
	flat := flags.Bool("flat", true, "")

	flags.Parse(arguments)

	mc.translationFiles = flags.Args()
	mc.sourceLanguage = *sourceLanguage
	mc.outdir = *outdir
	mc.format = *format
	if *format == "toml" {
		mc.flat = true
	} else {
		mc.flat = *flat
	}
}

func (mc *mergeCommand) SetArgs(args []string) {
	mc.translationFiles = args
}

func (mc *mergeCommand) writeFile(label string, translations []translation.Translation, localeID string) error {
	sort.Sort(translation.SortableByID(translations))

	var convert func([]translation.Translation) interface{}
	if mc.flat {
		convert = marshalFlatInterface
	} else {
		convert = marshalInterface
	}

	buf, err := mc.marshal(convert(translations))
	if err != nil {
		return fmt.Errorf("failed to marshal %s strings to %s: %s", localeID, mc.format, err)
	}

	filename := filepath.Join(mc.outdir, fmt.Sprintf("%s.%s.%s", localeID, label, mc.format))

	if err := ioutil.WriteFile(filename, buf, 0666); err != nil {
		return fmt.Errorf("failed to write %s: %s", filename, err)
	}
	return nil
}

func filter(translations map[string]translation.Translation, f func(translation.Translation) translation.Translation) []translation.Translation {
	filtered := make([]translation.Translation, 0, len(translations))
	for _, translation := range translations {
		if t := f(translation); t != nil {
			filtered = append(filtered, t)
		}
	}
	return filtered

}

func marshalFlatInterface(translations []translation.Translation) interface{} {
	mi := make(map[string]interface{}, len(translations))
	for _, translation := range translations {
		mi[translation.ID()] = translation.MarshalFlatInterface()
	}
	return mi
}

func marshalInterface(translations []translation.Translation) interface{} {
	mi := make([]interface{}, len(translations))
	for i, translation := range translations {
		mi[i] = translation.MarshalInterface()
	}
	return mi
}

func (mc mergeCommand) marshal(v interface{}) ([]byte, error) {
	switch mc.format {
	case "json":
		return json.MarshalIndent(v, "", "  ")
	case "toml":
		return marshalTOML(v)
	case "yaml":
		return yaml.Marshal(v)
	}
	return nil, fmt.Errorf("unsupported format: %s\n", mc.format)
}

func marshalTOML(v interface{}) ([]byte, error) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid format for marshaling to TOML")
	}
	tree, err := toml.TreeFromMap(m)
	if err != nil {
		return nil, err
	}
	s, err := tree.ToTomlString()
	return []byte(s), err
}

func usageMerge() {
	fmt.Printf(`Merge translation files.

Usage:

    goi18n merge [options] [files...]

Translation files:

    A translation file contains the strings and translations for a single language.

    Translation file names must have a suffix of a supported format (e.g. .json) and
    contain a valid language tag as defined by RFC 5646 (e.g. en-us, fr, zh-hant, etc.).

    For each language represented by at least one input translation file, goi18n will produce 2 output files:

        xx-yy.all.format
            This file contains all strings for the language (translated and untranslated).
            Use this file when loading strings at runtime.

        xx-yy.untranslated.format
            This file contains the strings that have not been translated for this language.
            The translations for the strings in this file will be extracted from the source language.
            After they are translated, merge them back into xx-yy.all.format using goi18n.

Merging:

    goi18n will merge multiple translation files for the same language.
    Duplicate translations will be merged into the existing translation.
    Non-empty fields in the duplicate translation will overwrite those fields in the existing translation.
    Empty fields in the duplicate translation are ignored.

Adding a new language:

    To produce translation files for a new language, create an empty translation file with the
    appropriate name and pass it in to goi18n.

Options:

    -sourceLanguage tag
        goi18n uses the strings from this language to seed the translations for other languages.
        Default: en-us

    -outdir directory
        goi18n writes the output translation files to this directory.
        Default: .

    -format format
        goi18n encodes the output translation files in this format.
        Supported formats: json, toml, yaml
        Default: json

    -flat
        goi18n writes the output translation files in flat format.
        Usage of '-format toml' automitically sets this flag.
        Default: true

`)
}
