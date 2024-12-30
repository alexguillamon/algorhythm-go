package rhymes

import (
	"algorhytm/language"
	"strings"
	"unsafe"

	mapset "github.com/deckarep/golang-set/v2"
)

type taggedChain struct {
	Node                  *language.TrieNode
	Sequence              []string
	consecutiveConsonants int
}

func tagPhonemeSequence(phonemeSequence []string) []string {
	taggedSequence := make([]string, len(phonemeSequence))
	for i, phoneme := range phonemeSequence {
		taggedSequence[i] = "-" + phoneme
	}
	return taggedSequence
}

func FindAdditive(lang *language.Language, phonemeSequence []string) (*mapset.Set[string], error) {
	cache := mapset.NewSet[StateKey]()

	stack := []taggedChain{{lang.Trie.GetRoot(), tagPhonemeSequence(phonemeSequence), 0}}

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
		childNode, exists := currentNode.Children[lang.PhoneticAlphabet.CleanPhoneme(currentPhoneme)]
		if exists {
			if len(nextSequence) == 0 {
				stack = append(stack, taggedChain{childNode, nextSequence, consecutiveConsonants})
				continue
			}

			nextPhoneme := nextSequence[0]
			nextPhonemeNode, nextPhonemeExists := childNode.Children[lang.PhoneticAlphabet.CleanPhoneme(nextPhoneme)]

			family, err := lang.PhoneticAlphabet.GetPhonemeFamily(nextPhoneme)
			if err != nil {
				return nil, err
			}
			if family == "vowel" {
				stack = append(stack, taggedChain{childNode, nextSequence, 0})
			} else {
				if consecutiveConsonants > 4 {
					continue
				}

				if strings.Contains(nextPhoneme, "-") {
					familyPhonemes, err := lang.PhoneticAlphabet.GetFamilyPhonemes(family)
					if err != nil {
						return nil, err
					}

					for _, familyPhoneme := range familyPhonemes {
						if _, exists := childNode.Children[familyPhoneme]; exists {
							stack = append(stack, taggedChain{
								childNode,
								append(
									[]string{"-" + familyPhoneme},
									nextSequence[1:]...,
								),
								consecutiveConsonants + 1})
						}

						for _, consonant := range lang.PhoneticAlphabet.GetConsonants() {
							if nextPhonemeExists {
								_, consonantExistInNext := nextPhonemeNode.Children[consonant]
								if consonantExistInNext {
									stack = append(
										stack,
										taggedChain{
											childNode,
											append(
												[]string{"-" + familyPhoneme, consonant},
												nextSequence[1:]...,
											),
											consecutiveConsonants + 1,
										},
									)
								}
							}

							if _, consonantInNode := childNode.Children[consonant]; consonantInNode {
								stack = append(
									stack,
									taggedChain{
										childNode,
										append(
											[]string{consonant, "-" + familyPhoneme},
											nextSequence[1:]...,
										),
										consecutiveConsonants + 1,
									},
								)
							}

						}
					}

				} else {
					if nextPhonemeExists {
						stack = append(
							stack,
							taggedChain{
								childNode,
								nextSequence,
								consecutiveConsonants + 1,
							},
						)
					}

					for _, consonant := range lang.PhoneticAlphabet.GetConsonants() {
						if nextPhonemeExists {
							_, consonantExistInNext := nextPhonemeNode.Children[consonant]
							if consonantExistInNext {
								stack = append(
									stack,
									taggedChain{
										childNode,
										append(
											[]string{nextPhoneme, consonant},
											nextSequence[1:]...,
										),
										consecutiveConsonants + 1,
									},
								)
							}
						}

						if _, consonantInNode := childNode.Children[consonant]; consonantInNode {
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
				}
			}
		}
	}

	return &words, nil
}
