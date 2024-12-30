package rhymes

import (
	"algorhytm/language"

	mapset "github.com/deckarep/golang-set/v2"
)

type RhymeService struct {
	language          *language.Language
	rhymesIdentity    mapset.Set[string]
	rhymesPerfect     mapset.Set[string]
	rhymesFamily      mapset.Set[string]
	rhymesAdditive    mapset.Set[string]
	rhymesSubtractive mapset.Set[string]
	rhymesAssonance   mapset.Set[string]
	resultsClose      mapset.Set[string]
	resultsMedium     mapset.Set[string]
}

// NewRhymeService initializes a new RhymeService
func NewRhymeService(language *language.Language) *RhymeService {
	return &RhymeService{
		language:          language,
		rhymesIdentity:    mapset.NewSet[string](),
		rhymesPerfect:     mapset.NewSet[string](),
		rhymesFamily:      mapset.NewSet[string](),
		rhymesAdditive:    mapset.NewSet[string](),
		rhymesSubtractive: mapset.NewSet[string](),
		rhymesAssonance:   mapset.NewSet[string](),
		resultsClose:      mapset.NewSet[string](),
		resultsMedium:     mapset.NewSet[string](),
	}
}

// FindRhymeStart determines the starting index of the rhyme

func (rs *RhymeService) rhymesFindAll(word string) {
	raw, exists := (*rs.language.PhonemeDictionary)[word]
	if !exists {

		return
	}

	phoneticSequence := rs.language.PhoneticAlphabet.FindStressSquence(raw[0])

	rs.rhymesPerfect = rs.language.Trie.Search(phoneticSequence)
	rs.rhymesFamily = FindFamily(rs.language, phoneticSequence)
	rs.rhymesAdditive = FindAdditive(rs.language, phoneticSequence)
	rs.rhymesSubtractive = FindSubstractive(rs.language, phoneticSequence)
	rs.rhymesAssonance = FindAssonance(rs.language, phoneticSequence)

	rs.resultsClose = mapset.NewSet[string]()
	rs.resultsClose = rs.resultsClose.Union(rs.rhymesPerfect).
		Union(rs.rhymesFamily).
		Union(rs.rhymesSubtractive).
		Union(rs.rhymesAdditive).
		Difference(rs.rhymesIdentity)

	rs.resultsMedium = rs.rhymesAssonance.Difference(rs.resultsClose).Difference(rs.rhymesIdentity)
}

func (rs *RhymeService) RhymesFind(word string) map[string][]string {
	rs.rhymesFindAll(word)

	return map[string][]string{
		"close":  rs.resultsClose.ToSlice(),
		"medium": rs.resultsMedium.ToSlice(),
	}
}
