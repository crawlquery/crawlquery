package token

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

// tokenize takes HTML content, extracts text, and splits it into tokens.
func Tokenize(htmlContent string) map[string][]int {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		panic(err) // Proper error handling should replace panic in production
	}

	var textBuilder strings.Builder

	if titleText := doc.Find("title").Text(); titleText != "" {
		textBuilder.WriteString(strings.TrimSpace(titleText) + " ")
	}

	if bodyText := doc.Find("body").Contents().Not("body > *").Text(); bodyText != "" {
		textBuilder.WriteString(strings.TrimSpace(bodyText) + " ")
	}
	// Iterate over each node within the body, ensuring proper spacing
	doc.Find("body *").Each(func(i int, selection *goquery.Selection) {
		nodeText := selection.Text()
		// Ensure each text node ends with a space to prevent merging
		trimmedText := strings.TrimSpace(nodeText) + " "
		textBuilder.WriteString(trimmedText)
	})

	// Consolidated text from the builder
	consolidatedText := textBuilder.String()
	fmt.Println("Consolidated Text:", consolidatedText)

	// Normalize text: convert to lower case
	normalizedText := strings.ToLower(consolidatedText)

	// Remove punctuation using a regular expression
	// Ensuring spaces are not removed by the regex
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	finalText := reg.ReplaceAllString(normalizedText, "")

	// Remove multiple spaces resulting from removals and edge cases
	spaceCleanedText := strings.Join(strings.Fields(finalText), " ")

	// Split text into words based on whitespace
	words := strings.Fields(spaceCleanedText)

	tokens := make(map[string][]int)
	position := 0
	for _, word := range words {
		tokens[word] = append(tokens[word], position)
		position++
	}

	return tokens
}