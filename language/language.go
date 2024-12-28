package language

import "algorhytm/phonetics"

type Language struct {
	*PhonemeDictionary
	Trie             *trie
	PhoneticAlphabet *phonetics.Alphabet
	*POSDictionary
}

type PhonemeDictionary map[string][][]string

type POSDictionary map[string]map[string][][]string

type LanguageProperties struct {
	*PhonemeDictionary
	*POSDictionary
	Alphabet *phonetics.Alphabet
}

func Initialize(properties *LanguageProperties) *Language {
	return &Language{
		PhonemeDictionary: properties.PhonemeDictionary,
		Trie:              buildTrie(properties.PhonemeDictionary, properties.Alphabet, nil),
		PhoneticAlphabet:  properties.Alphabet,
		POSDictionary:     properties.POSDictionary,
	}
}
