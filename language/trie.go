package language

import (
	"algorhytm/phonetics"

	mapset "github.com/deckarep/golang-set/v2"
)

type trie struct {
	root *trieNode
}

type trieNode struct {
	children       map[string]*trieNode
	isEndOfWord    bool
	wordReferences mapset.Set[string]
}

func (t *trie) Search(sequence []string) mapset.Set[string] {
	node := t.root
	for _, phoneme := range sequence {
		_, ok := node.children[phoneme]
		if !ok {
			return mapset.NewSet[string]()
		}
		node = node.children[phoneme]
	}
	return node.wordReferences
}

func (t *trie) insert(sequence []string, wordReference string) {
	node := t.root
	for _, phoneme := range sequence {
		_, ok := node.children[phoneme]
		if !ok {
			node.children[phoneme] = &trieNode{
				children:       make(map[string]*trieNode),
				wordReferences: mapset.NewSet[string](),
			}
		}
		node = node.children[phoneme]
	}
	node.isEndOfWord = true
	node.wordReferences.Add(wordReference)
}

func buildTrie(
	phonemeDictionary *PhonemeDictionary,
	alphabet *phonetics.Alphabet,
	t *trie,
) *trie {
	if t == nil {
		t = &trie{
			root: &trieNode{
				children:       make(map[string]*trieNode),
				wordReferences: mapset.NewSet[string](),
			},
		}

	}

	for word, pronunciations := range *phonemeDictionary {
		for _, pronunciation := range pronunciations {

			stress_idx := alphabet.FindStress(pronunciation)
			if stress_idx == -1 {
				continue
			}
			phoneme_sequence := pronunciation[stress_idx:]
			t.insert(phoneme_sequence, word)
		}
	}
	return t

}
