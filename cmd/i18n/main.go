package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/jesseduffield/lazygit/pkg/i18n"
)

func saveLanguageFileToJson(tr *i18n.TranslationSet, filepath string) error {
	jsonData, err := json.MarshalIndent(tr, "", "  ")
	if err != nil {
		return err
	}

	jsonData = append(jsonData, '\n')
	return os.WriteFile(filepath, jsonData, 0o644)
}

func saveNonEnglishLanguageFilesToJson() error {
	translationSets, _ := i18n.GetTranslationSets()
	for lang, tr := range translationSets {
		if lang == "en" {
			continue
		}

		err := saveLanguageFileToJson(tr, "pkg/i18n/translations/"+lang+".json")
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	err := saveNonEnglishLanguageFilesToJson()
	if err != nil {
		log.Fatal(err)
	}
}
