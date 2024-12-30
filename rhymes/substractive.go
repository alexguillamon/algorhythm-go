package rhymes

import (
	"algorhytm/language"
	"strings"
	"unsafe"

	mapset "github.com/deckarep/golang-set/v2"
)

type StateKey struct {
	NodeAddr uintptr // Unique memory address of the node
	Sequence string  // Hash of the current sequence
}

func FindSubstractive(lang *language.Language, phonemeSequence []string) (*mapset.Set[string], error) {
	cache := mapset.NewSet[StateKey]()
	stack := []chain{{lang.Trie.GetRoot(), phonemeSequence}}

	words := mapset.NewSet[string]()

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		currentNode := current.Node
		currentSequence := current.Sequence

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
				stack = append(stack, chain{childNode, nextSequence})
				continue
			}

			family, err := lang.PhoneticAlphabet.GetPhonemeFamily(nextSequence[0])
			if err != nil {
				return nil, err
			}
			if family == "vowel" {
				stack = append(stack, chain{childNode, nextSequence})
			} else {
				phonemes, err := lang.PhoneticAlphabet.GetFamilyPhonemes(family)
				if err != nil {
					return nil, err
				}

				for _, phoneme := range phonemes {
					if _, exists := childNode.Children[phoneme]; exists {
						stack = append(stack, chain{childNode, append([]string{phoneme}, nextSequence[1:]...)})
					}
				}
				stack = append(stack, chain{childNode, nextSequence[1:]})
				stack = append(stack, chain{currentNode, nextSequence[1:]})

			}
		}

	}
	return &words, nil
}
