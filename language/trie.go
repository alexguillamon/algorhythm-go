package language

import (
	"algorhytm/orderedset"
	"algorhytm/phonetics"
	"strings"
)

type trie struct {
	root *TrieNode
}

type TrieNode struct {
	Children       map[string]*TrieNode
	IsEndOfWord    bool
	WordReferences orderedset.Set[string]
}

func (t *trie) GetRoot() *TrieNode {
	return t.root
}

func (t *trie) Search(sequence []string) orderedset.Set[string] {
	node := t.root
	for _, phoneme := range sequence {
		_, ok := node.Children[phoneme]
		if !ok {
			return orderedset.NewSet[string]()
		}
		node = node.Children[phoneme]
	}
	return node.WordReferences
}

func (t *trie) insert(sequence []string, wordReference string) {
	node := t.root
	for _, phoneme := range sequence {
		_, ok := node.Children[phoneme]
		if !ok {
			node.Children[phoneme] = &TrieNode{
				Children:       make(map[string]*TrieNode),
				WordReferences: orderedset.NewSet[string](),
			}
		}
		node = node.Children[phoneme]
	}
	node.IsEndOfWord = true
	node.WordReferences.Add(wordReference)
}

func buildTrie(
	phonemeDictionary *PhonemeDictionary,
	alphabet *phonetics.Alphabet,
	t *trie,
) *trie {
	if t == nil {
		t = &trie{
			root: &TrieNode{
				Children:       make(map[string]*TrieNode),
				WordReferences: orderedset.NewSet[string](),
			},
		}

	}

	for word, pronunciations := range *phonemeDictionary {
		if strings.Contains(word, "'") || strings.Contains(word, "-") || strings.ContainsAny(word, "_'-") {
			continue
		}

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
