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

func main() {
	err := saveLanguageFileToJson(i18n.EnglishTranslationSet(), "en.json")
	if err != nil {
		log.Fatal(err)
	}
}
