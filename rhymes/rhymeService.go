package rhymes

import (
	"algorhytm/language"

	mapset "github.com/deckarep/golang-set/v2"
)

// TODO: find a way to order results (is there a way to have an ordered set?, also if not we could probably live with just deconstructing the rhymes in order and and allow the inside of each to not have order)
// TODO: fix error handling
// - add panic to functions that should never error out but are possible to error out (most things in the alphabet, the program should never pass things that don't exist or work for these)
// - fix the error handling of not finding words
// TODO: figure out how should response objects be created and where would they go
// TODO: write unit test for rhyme algorithms and for rhyming service
// TODO: write unit test for trie and trie creation
// TODO: write unit test for language?
// TODO: add hot reload
// TODO: create make file for project
// TODO: revisit the structure of the project and also look to make the rhyming function private

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
	resultsFar        mapset.Set[string]
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
		resultsFar:        mapset.NewSet[string](),
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

	rs.resultsClose = rs.rhymesPerfect.
		Union(rs.rhymesFamily).
		Difference(rs.rhymesIdentity)

	rs.resultsMedium = rs.rhymesAdditive.
		Union(rs.rhymesSubtractive).
		Difference(rs.resultsClose).
		Difference(rs.rhymesIdentity)
	rs.resultsFar = rs.rhymesAssonance.
		Difference(rs.resultsClose).
		Difference(rs.resultsMedium).
		Difference(rs.rhymesIdentity)
}

func (rs *RhymeService) RhymesFind(word string) map[string][]string {
	rs.rhymesFindAll(word)

	return map[string][]string{
		"close":  rs.resultsClose.ToSlice(),
		"medium": rs.resultsMedium.ToSlice(),
		"far":    rs.resultsFar.ToSlice(),
	}
}
