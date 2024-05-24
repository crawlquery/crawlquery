package phrase

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/jdkato/prose/v2"
)

func TestWordClasses(t *testing.T) {
	sentence := "He run quickly."
	doc, _ := prose.NewDocument(sentence)

	tokens := doc.Tokens()

	fmt.Printf("Tokens: %v\n", tokens)

	t.Fail()
}

func sortPhrases(phrases [][]string) {
	sort.Slice(phrases, func(i, j int) bool {
		return phraseString(phrases[i]) < phraseString(phrases[j])
	})
}

func phraseString(phrase []string) string {
	return strings.Join(phrase, " ")
}
