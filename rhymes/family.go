package rhymes

import (
	"algorhytm/language"
	"algorhytm/orderedset"
)

type chain struct {
	Node     *language.TrieNode
	Sequence []string
}

func FindFamily(lang *language.Language, phonemeSequence []string) orderedset.Set[string] {
	stack := []chain{{lang.Trie.GetRoot(), phonemeSequence}}

	words := orderedset.NewSet[string]()

	for len(stack) > 0 {
		// Pop the stack.
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		currentNode := current.Node
		currentSequence := current.Sequence

		// If the sequence is empty, add word references to the result.
		if len(currentSequence) == 0 {
			if currentNode.IsEndOfWord {
				words.Add(currentNode.WordReferences.ToSlice()...)
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
				return words
			}
			if family == "vowel" {
				stack = append(stack, chain{childNode, nextSequence})
			} else {
				phonemes, err := lang.PhoneticAlphabet.GetFamilyPhonemes(family)
				if err != nil {
					return words
				}

				for _, phoneme := range phonemes {
					if _, exists := childNode.Children[phoneme]; exists {
						stack = append(stack, chain{childNode, append([]string{phoneme}, nextSequence[1:]...)})
					}
				}
			}
		}
	}

	return words
}
