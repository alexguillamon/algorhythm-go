package rhymes

import (
	"algorhytm/language"
	"strings"
	"unsafe"

	mapset "github.com/deckarep/golang-set/v2"
)

func FindAssonance(lang *language.Language, phonemeSequence []string) (*mapset.Set[string], error) {
	cache := mapset.NewSet[StateKey]()

	stack := []taggedChain{{lang.Trie.GetRoot(), phonemeSequence, 0}}

	words := mapset.NewSet[string]()

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		currentNode := current.Node
		currentSequence := current.Sequence
		consecutiveConsonants := current.consecutiveConsonants

		currentSequenceHash := strings.Join(currentSequence, ",")
		stateKey := StateKey{
			NodeAddr: uintptr(unsafe.Pointer(currentNode)),
			Sequence: currentSequenceHash,
		}

		// Skip if the state has already been visited
		if cache.Contains(stateKey) {
			continue
		}

		// Mark the state as visited
		cache.Add(stateKey)

		// If the sequence is empty, add word references to the result.
		if len(currentSequence) == 0 {
			if currentNode.IsEndOfWord {
				words.Append(currentNode.WordReferences.ToSlice()...)
			}
			continue
		}

		currentPhoneme := currentSequence[0]
		nextSequence := currentSequence[1:]

		// Check if the child node exists for the current phoneme.
		childNode, exists := currentNode.Children[currentPhoneme]
		if exists {
			if len(nextSequence) == 0 {
				stack = append(stack, taggedChain{childNode, nextSequence, consecutiveConsonants})
				continue
			}

			if lang.PhoneticAlphabet.IsConsonant(nextSequence[0]) {
				if consecutiveConsonants > 4 {
					continue
				}
				stack = append(
					stack,
					taggedChain{
						childNode,
						nextSequence[1:],
						consecutiveConsonants + 1,
					},
				)
				stack = expandConsonantPhoneme(
					stack,
					childNode,
					nextSequence,
					consecutiveConsonants,
					lang,
				)

				if len(nextSequence) > 1 && lang.PhoneticAlphabet.IsVowel(nextSequence[1]) {
					stack = expandVowelPhoneme(
						stack,
						childNode,
						nextSequence[1:],
						lang,
					)
				}

			} else {
				stack = expandVowelPhoneme(
					stack,
					childNode,
					nextSequence,
					lang,
				)

			}

		}
	}

	return &words, nil
}

func expandVowelPhoneme(
	stack []taggedChain,
	childNode *language.TrieNode,
	nextSequence []string,
	lang *language.Language,
) []taggedChain {
	// For all the vowels
	phonemes, err := lang.PhoneticAlphabet.GetFamilyPhonemes("vowel")
	if err != nil {
		return stack
	}

	for _, phoneme := range phonemes {
		stressedPhoneme := lang.PhoneticAlphabet.MatchPhonemeStress(nextSequence[0], phoneme)
		if _, exists := childNode.Children[stressedPhoneme]; exists {
			// Try all in position
			stack = append(
				stack,
				taggedChain{
					childNode,
					append([]string{stressedPhoneme}, nextSequence[1:]...),
					0,
				},
			)

		}

	}
	return stack
}

func expandConsonantPhoneme(
	stack []taggedChain,
	childNode *language.TrieNode,
	nextSequence []string,
	consecutiveConsonants int,
	lang *language.Language,
) []taggedChain {
	// For all the consonants
	for _, consonant := range lang.PhoneticAlphabet.GetConsonants() {
		if _, exists := childNode.Children[consonant]; exists {
			// Try all in current position
			stack = append(
				stack,
				taggedChain{
					childNode,
					append([]string{consonant}, nextSequence[1:]...),
					consecutiveConsonants + 1,
				},
			)
			// Try all before the current phoneme
			stack = append(
				stack,
				taggedChain{
					childNode,
					append([]string{consonant}, nextSequence...),
					consecutiveConsonants + 1,
				},
			)
		}

	}
	return stack
}
