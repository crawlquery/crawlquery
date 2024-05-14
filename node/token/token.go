package token

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	goose "github.com/advancedlogic/GoOse"
	rake "github.com/afjoseph/RAKE.Go"
	"github.com/securisec/go-keywords"
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

func Keywords(data string) []string {
	k, _ := keywords.Extract(data, keywords.ExtractOptions{
		StripTags:        true,
		RemoveDuplicates: true,
		IgnorePattern:    "<.+>",
		Lowercase:        true,
	})

	return k
}

func RakeKeywords(data string) []string {
	candidates := rake.RunRake(data)

	words := make([]string, 0)

	for _, candidate := range candidates {
		words = append(words, candidate.Key)
	}

	return words
}

func Positions(tokens []string) map[string][]int {
	tokenPositions := make(map[string][]int)
	position := 0
	for _, token := range tokens {
		tokenPositions[token] = append(tokenPositions[token], position)
		position++
	}

	return tokenPositions
}

// tokenize takes HTML content, extracts text, and splits it into tokens.
func Tokenize(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		panic(err) // Proper error handling should replace panic in production
	}

	var textBuilder strings.Builder

	doc.Find("script, style, noscript").Remove()

	// Iterate over each text node within the body, ensuring proper spacing
	doc.Find("body").Contents().Each(func(i int, selection *goquery.Selection) {
		if goquery.NodeName(selection) == "#text" {
			text := strings.TrimSpace(selection.Text())
			if text != "" {
				textBuilder.WriteString(text + " ")
			}
		}
	})

	// Iterate over all elements in the body to extract text content
	doc.Find("body *").Each(func(i int, selection *goquery.Selection) {
		if goquery.NodeName(selection) != "script" && goquery.NodeName(selection) != "style" && goquery.NodeName(selection) != "noscript" {
			text := strings.TrimSpace(selection.Text())
			if text != "" {
				textBuilder.WriteString(text + " ")
			}
		}
	})

	// Consolidated text from the builder
	consolidatedText := textBuilder.String()

	// Normalize text: convert to lower case
	normalizedText := strings.ToLower(consolidatedText)

	// Remove punctuation using a regular expression
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	finalText := reg.ReplaceAllString(normalizedText, "")

	// Return the cleaned text
	return finalText

}

func Tokens() []string {
	g := goose.New()
	article, _ := g.ExtractFromURL("http://edition.cnn.com/2012/07/08/opinion/banzi-ted-open-source/index.html")
	println("title", article.Title)
	println("description", article.MetaDescription)
	println("keywords", article.MetaKeywords)
	println("content", article.CleanedText)
	println("url", article.FinalURL)
	println("top image", article.TopImage)
}
