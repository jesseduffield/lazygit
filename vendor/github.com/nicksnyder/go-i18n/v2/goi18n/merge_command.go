package main

import (
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicksnyder/go-i18n/v2/internal"
	"github.com/nicksnyder/go-i18n/v2/internal/plural"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

func usageMerge() {
	fmt.Fprintf(os.Stderr, `usage: goi18n merge [options] [message files]

Merge reads all messages in the message files and produces two files per language.

	xx-yy.active.format
		This file contains messages that should be loaded at runtime.

	xx-yy.translate.format
		This file contains messages that are empty and should be translated.

Message file names must have a suffix of a supported format (e.g. ".json") and
contain a valid language tag as defined by RFC 5646 (e.g. "en-us", "fr", "zh-hant", etc.).

To add support for a new language, create an empty translation file with the
appropriate name and pass it in to goi18n merge.

Flags:

	-sourceLanguage tag
		Translate messages from this language (e.g. en, en-US, zh-Hant-CN)
 		Default: en

	-outdir directory
		Write message files to this directory.
		Default: .

	-format format
		Output message files in this format.
		Supported formats: json, toml, yaml
		Default: toml
`)
}

type mergeCommand struct {
	messageFiles   []string
	sourceLanguage languageTag
	outdir         string
	format         string
}

func (mc *mergeCommand) name() string {
	return "merge"
}

func (mc *mergeCommand) parse(args []string) {
	flags := flag.NewFlagSet("merge", flag.ExitOnError)
	flags.Usage = usageMerge

	flags.Var(&mc.sourceLanguage, "sourceLanguage", "en")
	flags.StringVar(&mc.outdir, "outdir", ".", "")
	flags.StringVar(&mc.format, "format", "toml", "")
	flags.Parse(args)

	mc.messageFiles = flags.Args()
}

func (mc *mergeCommand) execute() error {
	if len(mc.messageFiles) < 1 {
		return fmt.Errorf("need at least one message file to parse")
	}
	inFiles := make(map[string][]byte)
	for _, path := range mc.messageFiles {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		inFiles[path] = content
	}
	ops, err := merge(inFiles, mc.sourceLanguage.Tag(), mc.outdir, mc.format)
	if err != nil {
		return err
	}
	for path, content := range ops.writeFiles {
		if err := ioutil.WriteFile(path, content, 0666); err != nil {
			return err
		}
	}
	for _, path := range ops.deleteFiles {
		// Ignore error since it isn't guaranteed to exist.
		os.Remove(path)
	}
	return nil
}

type fileSystemOp struct {
	writeFiles  map[string][]byte
	deleteFiles []string
}

func merge(messageFiles map[string][]byte, sourceLanguageTag language.Tag, outdir, outputFormat string) (*fileSystemOp, error) {
	unmerged := make(map[language.Tag][]map[string]*internal.MessageTemplate)
	sourceMessageTemplates := make(map[string]*internal.MessageTemplate)
	unmarshalFuncs := map[string]internal.UnmarshalFunc{
		"json": json.Unmarshal,
		"toml": toml.Unmarshal,
		"yaml": yaml.Unmarshal,
	}
	for path, content := range messageFiles {
		mf, err := internal.ParseMessageFileBytes(content, path, unmarshalFuncs)
		if err != nil {
			return nil, fmt.Errorf("failed to load message file %s: %s", path, err)
		}
		templates := map[string]*internal.MessageTemplate{}
		for _, m := range mf.Messages {
			templates[m.ID] = internal.NewMessageTemplate(m)
		}
		if mf.Tag == sourceLanguageTag {
			for _, template := range templates {
				if sourceMessageTemplates[template.ID] != nil {
					return nil, fmt.Errorf("multiple source translations for id %s", template.ID)
				}
				template.Hash = hash(template)
				sourceMessageTemplates[template.ID] = template
			}
		}
		unmerged[mf.Tag] = append(unmerged[mf.Tag], templates)
	}

	if len(sourceMessageTemplates) == 0 {
		return nil, fmt.Errorf("no messages found for source locale %s", sourceLanguageTag)
	}

	pluralRules := plural.DefaultRules()
	all := make(map[language.Tag]map[string]*internal.MessageTemplate)
	all[sourceLanguageTag] = sourceMessageTemplates
	for _, srcTemplate := range sourceMessageTemplates {
		for dstLangTag, messageTemplates := range unmerged {
			if dstLangTag == sourceLanguageTag {
				continue
			}
			pluralRule := pluralRules.Rule(dstLangTag)
			if pluralRule == nil {
				// Non-standard languages not supported because
				// we don't know if translations are complete or not.
				continue
			}
			if all[dstLangTag] == nil {
				all[dstLangTag] = make(map[string]*internal.MessageTemplate)
			}
			dstMessageTemplate := all[dstLangTag][srcTemplate.ID]
			if dstMessageTemplate == nil {
				dstMessageTemplate = &internal.MessageTemplate{
					Message: &i18n.Message{
						ID:          srcTemplate.ID,
						Description: srcTemplate.Description,
						Hash:        srcTemplate.Hash,
					},
					PluralTemplates: make(map[plural.Form]*internal.Template),
				}
				all[dstLangTag][srcTemplate.ID] = dstMessageTemplate
			}

			// Check all unmerged message templates for this message id.
			for _, messageTemplates := range messageTemplates {
				unmergedTemplate := messageTemplates[srcTemplate.ID]
				if unmergedTemplate == nil {
					continue
				}
				// Ignore empty hashes for v1 backward compatibility.
				if unmergedTemplate.Hash != "" && unmergedTemplate.Hash != srcTemplate.Hash {
					// This was translated from different content so discard.
					continue
				}

				// Merge in the translated messages.
				for pluralForm := range pluralRule.PluralForms {
					dt := unmergedTemplate.PluralTemplates[pluralForm]
					if dt != nil && dt.Src != "" {
						dstMessageTemplate.PluralTemplates[pluralForm] = dt
					}
				}
			}
		}
	}

	translate := make(map[language.Tag]map[string]*internal.MessageTemplate)
	active := make(map[language.Tag]map[string]*internal.MessageTemplate)
	for langTag, messageTemplates := range all {
		active[langTag] = make(map[string]*internal.MessageTemplate)
		if langTag == sourceLanguageTag {
			active[langTag] = messageTemplates
			continue
		}
		pluralRule := pluralRules.Rule(langTag)
		if pluralRule == nil {
			// Non-standard languages not supported because
			// we don't know if translations are complete or not.
			continue
		}
		for _, messageTemplate := range messageTemplates {
			srcMessageTemplate := sourceMessageTemplates[messageTemplate.ID]
			activeMessageTemplate, translateMessageTemplate := activeDst(srcMessageTemplate, messageTemplate, pluralRule)
			if translateMessageTemplate != nil {
				if translate[langTag] == nil {
					translate[langTag] = make(map[string]*internal.MessageTemplate)
				}
				translate[langTag][messageTemplate.ID] = translateMessageTemplate
			}
			if activeMessageTemplate != nil {
				active[langTag][messageTemplate.ID] = activeMessageTemplate
			}
		}
	}

	writeFiles := make(map[string][]byte, len(translate)+len(active))
	for langTag, messageTemplates := range translate {
		path, content, err := writeFile(outdir, "translate", langTag, outputFormat, messageTemplates, false)
		if err != nil {
			return nil, err
		}
		writeFiles[path] = content
	}
	deleteFiles := []string{}
	for langTag, messageTemplates := range active {
		path, content, err := writeFile(outdir, "active", langTag, outputFormat, messageTemplates, langTag == sourceLanguageTag)
		if err != nil {
			return nil, err
		}
		if len(content) > 0 {
			writeFiles[path] = content
		} else {
			deleteFiles = append(deleteFiles, path)
		}
	}
	return &fileSystemOp{writeFiles: writeFiles, deleteFiles: deleteFiles}, nil
}

// activeDst returns the active part of the dst and whether dst is a complete translation of src.
func activeDst(src, dst *internal.MessageTemplate, pluralRule *plural.Rule) (active *internal.MessageTemplate, translateMessageTemplate *internal.MessageTemplate) {
	pluralForms := pluralRule.PluralForms
	if len(src.PluralTemplates) == 1 {
		pluralForms = map[plural.Form]struct{}{
			plural.Other: {},
		}
	}
	for pluralForm := range pluralForms {
		dt := dst.PluralTemplates[pluralForm]
		if dt == nil || dt.Src == "" {
			if translateMessageTemplate == nil {
				translateMessageTemplate = &internal.MessageTemplate{
					Message: &i18n.Message{
						ID:          src.ID,
						Description: src.Description,
						Hash:        src.Hash,
					},
					PluralTemplates: make(map[plural.Form]*internal.Template),
				}
			}
			translateMessageTemplate.PluralTemplates[pluralForm] = src.PluralTemplates[plural.Other]
			continue
		}
		if active == nil {
			active = &internal.MessageTemplate{
				Message: &i18n.Message{
					ID:          src.ID,
					Description: src.Description,
					Hash:        src.Hash,
				},
				PluralTemplates: make(map[plural.Form]*internal.Template),
			}
		}
		active.PluralTemplates[pluralForm] = dt
	}
	return
}

func hash(t *internal.MessageTemplate) string {
	h := sha1.New()
	io.WriteString(h, t.Description)
	io.WriteString(h, t.PluralTemplates[plural.Other].Src)
	return fmt.Sprintf("sha1-%x", h.Sum(nil))
}
