package phonetics

import (
	"fmt"
	"regexp"
	"strings"
)

type accentChars struct {
	Primary   string
	Secondary string
	None      string
}

// type IPhoneticAlphabet interface {
// 	GetPhonemeFamily(phoneme string) (string, error)
// 	GetFamilyPhonemes(family string) []string
// 	GetConsonants() []string
// 	IsVowel(phoneme string) (bool, error)
// 	IsConsonant(phoneme string) (bool, error)
// 	CleanPhoneme(phoneme string) string
// 	MatchPhonemeStress(origin, target string) string
// }

type Alphabet struct {
	accentChars     accentChars
	phonemeToFamily map[string]string
	familyToPhoneme map[string][]string
	consonants      map[string]string
}

func buildAlphabet(phonemeToFamily map[string]string, accent accentChars) *Alphabet {
	pa := &Alphabet{
		accentChars:     accent,
		phonemeToFamily: phonemeToFamily,
	}
	pa.familyToPhoneme = pa.createFamilyToPhonemes()
	pa.consonants = pa.createConsonants()
	return pa
}

func (pa *Alphabet) createFamilyToPhonemes() map[string][]string {
	familyMap := make(map[string][]string)
	for phone, fam := range pa.phonemeToFamily {
		familyMap[fam] = append(familyMap[fam], phone)
	}
	return familyMap
}

func (pa *Alphabet) createConsonants() map[string]string {
	consonants := make(map[string]string)
	for phone, fam := range pa.phonemeToFamily {
		if fam != "vowel" {
			consonants[phone] = fam
		}
	}
	return consonants
}

func (pa *Alphabet) GetPrimaryStressChar() string {
	return pa.accentChars.Primary
}

func (pa *Alphabet) GetSecondaryStressChar() string {
	return pa.accentChars.Secondary
}

func (pa *Alphabet) GetPhonemeFamily(phoneme string) (string, error) {
	cleaned := pa.CleanPhoneme(phoneme)
	family, ok := pa.phonemeToFamily[cleaned]
	if !ok {
		return "", fmt.Errorf("phoneme %q not found in the alphabet", phoneme)
	}
	return family, nil
}

func (pa *Alphabet) GetFamilyPhonemes(family string) ([]string, error) {
	phonemes, ok := pa.familyToPhoneme[family]
	if !ok {
		return nil, fmt.Errorf("family %q not found in the alphabet", family)
	}
	return phonemes, nil
}

func (pa *Alphabet) GetConsonants() []string {
	keys := make([]string, 0, len(pa.consonants))
	for c := range pa.consonants {
		keys = append(keys, c)
	}
	return keys
}

var reDigits = regexp.MustCompile(`\d`)

func (pa *Alphabet) CleanPhoneme(phoneme string) string {
	noDigits := reDigits.ReplaceAllString(phoneme, "")
	return strings.TrimPrefix(noDigits, "-")
}

func (pa *Alphabet) IsVowel(phoneme string) (bool, error) {
	family, err := pa.GetPhonemeFamily(phoneme)
	if err != nil {
		return false, err
	}
	return family == "vowel", nil
}

func (pa *Alphabet) IsConsonant(phoneme string) (bool, error) {
	family, err := pa.GetPhonemeFamily(phoneme)
	if err != nil {
		return false, err
	}
	return family != "vowel", nil
}

func (pa *Alphabet) MatchPhonemeStress(origin, target string) string {
	targetCleaned := pa.CleanPhoneme(target)

	switch {
	case strings.Contains(origin, pa.accentChars.Primary):
		return targetCleaned + pa.accentChars.Primary
	case strings.Contains(origin, pa.accentChars.Secondary):
		return targetCleaned + pa.accentChars.Secondary
	case strings.Contains(origin, pa.accentChars.None):
		return targetCleaned + pa.accentChars.None
	default:
		return targetCleaned
	}
}

func (pa *Alphabet) FindStress(sequence []string) int {
	stress_idx := -1
	for i := len(sequence) - 1; i >= 0; i-- {
		phoneme := sequence[i]
		if strings.Contains(phoneme, pa.GetPrimaryStressChar()) {
			return len(sequence) - i - 1
		}

	}
	return stress_idx
}
