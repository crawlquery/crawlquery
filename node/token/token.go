package token

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	rake "github.com/afjoseph/RAKE.Go"
	"github.com/securisec/go-keywords"
)

func Positions(tokens []string) map[string][]int {
	tokenPositions := make(map[string][]int)
	position := 0
	for _, token := range tokens {
		tokenPositions[token] = append(tokenPositions[token], position)
		position++
	}

	return tokenPositions
}

func TokenizeTerm(term string) []string {
	// Normalize text: convert to lower case
	normalizedText := strings.ToLower(term)

	// Remove punctuation using a regular expression
	// Ensuring spaces are not removed by the regex
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	finalText := reg.ReplaceAllString(normalizedText, "")

	// Remove multiple spaces resulting from removals and edge cases
	spaceCleanedText := strings.Join(strings.Fields(finalText), " ")

	// Split text into words based on whitespace
	words := strings.Fields(spaceCleanedText)

	return words
}

// removeUnwantedElements removes script, style, and comments from the document
func removeUnwantedElements(doc *goquery.Document) {
	doc.Find("script").Remove()
	doc.Find("style").Remove()
	doc.Find("noscript").Remove()
	doc.Find("comment").Remove()
	doc.Find("head").Remove()
	doc.Find("meta").Remove()
}

// extractTextRecursively extracts text from a node and its children, ensuring proper spacing
func extractTextRecursively(selection *goquery.Selection, textBuilder *strings.Builder) {
	if goquery.NodeName(selection) == "#text" {
		text := strings.TrimSpace(selection.Text())
		if text != "" {
			textBuilder.WriteString(text + " ")
		}
	} else {
		selection.Contents().Each(func(i int, child *goquery.Selection) {
			extractTextRecursively(child, textBuilder)
		})
	}
}

func RakeKeywords(data string) []string {
	candidates := rake.RunRake(data)

	words := make([]string, 0)

	for _, candidate := range candidates {
		words = append(words, candidate.Key)
	}

	return words
}

func TopKeywords(positions map[string][]int) []string {
	// Find the most used keywords
	max := 0
	for _, v := range positions {
		if len(v) > max {
			max = len(v)
		}
	}

	// Find the keywords that are used the most
	topKeywords := make([]string, 0)
	for k, v := range positions {
		if len(v) == max {
			topKeywords = append(topKeywords, k)
		}
	}

	return topKeywords
}

func Keywords(doc *goquery.Document) []string {
	removeUnwantedElements(doc)

	var textBuilder strings.Builder

	// Recursively extract text from the body
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		extractTextRecursively(s, &textBuilder)
	})

	// Consolidated text from the builder
	consolidatedText := textBuilder.String()

	// Normalize text: convert to lower case
	normalizedText := strings.ToLower(consolidatedText)

	// Remove punctuation using a regular expression, ensuring spaces are not removed by the regex
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	finalText := reg.ReplaceAllString(normalizedText, "")

	k, _ := keywords.Extract(finalText, keywords.ExtractOptions{})

	return k
}
