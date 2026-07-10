package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"turtle/core/lgr"
)

// trs stores translations for each language.
var trs map[string]map[string]interface{} = make(map[string]map[string]interface{})

// initT initializes the translations from the JSON files located in "./static/translation/".
func InitT() {
	defer Recover("Failed to load translations")

	folder := "./static/translation/"

	// Walk through the directory to find all the translation files.
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
			lang := info.Name()[:len(info.Name())-len(filepath.Ext(info.Name()))]
			file, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Load the translation data.
			var data map[string]interface{}
			if err := json.Unmarshal(file, &data); err != nil {
				return err
			}

			trs[lang] = data
		}

		return nil
	})

	if err != nil {
		lgr.Error("Failed to load languages: %v", err)
	}
}

// t translates the given key for the specified language, or returns the key if no translation is found.
func T(key, lang string) string {
	keyToTranslate := lang

	// Fallback to "en" if the language is not found.
	if _, ok := trs[keyToTranslate]; !ok {
		keyToTranslate = "en"
	}

	// If still not found, return the key.
	if _, ok := trs[keyToTranslate]; !ok {
		return key
	}

	langStorage := trs[keyToTranslate]

	if val, ok := langStorage[key]; ok {
		return fmt.Sprintf("%v", val) // Convert the value to string
	}

	return key
}

type Translator = func(string) string

// translatedT returns a function that translates keys based on the given language.
func TranslatedT(lang string) Translator {
	return func(key string) string {
		return T(key, lang)
	}
}
