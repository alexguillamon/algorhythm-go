package english

import (
	"algorhytm/language"
	"algorhytm/phonetics"
	"encoding/json"
	"fmt"
	"os"
)

func GetProperties() (*language.LanguageProperties, error) {

	flatFile, err := os.Open("language/english/isle_cleaned_flat.json")
	if err != nil {
		return nil, fmt.Errorf("%s was not able to open file", err)
	}
	defer flatFile.Close()

	posFile, err := os.Open("language/english/isle_cleaned_POS.json")
	if err != nil {
		return nil, fmt.Errorf("%s was not able to open file", err)
	}
	defer posFile.Close()

	var dict language.PhonemeDictionary
	var posDict language.POSDictionary

	if err := json.NewDecoder(flatFile).Decode(&dict); err != nil {
		return nil, fmt.Errorf("there was an error decoding the flat json: %v", err)
	}

	if err := json.NewDecoder(posFile).Decode(&posDict); err != nil {
		return nil, fmt.Errorf("there was an error decoding the pos json: %v", err)
	}

	return &language.LanguageProperties{
		PhonemeDictionary: &dict,
		Alphabet:          phonetics.IPA,
		POSDictionary:     &posDict,
	}, nil

}
