package parse

import goose "github.com/advancedlogic/GoOse"

type ParseResult struct {
	Title       string
	Description string
	Keywords    string
	Content     string
	URL         string
	TopImage    string
}

func Parse(htmlContent string, url string) (*ParseResult, error) {
	g := goose.New()
	article, err := g.ExtractFromRawHTML(htmlContent, url)

	if err != nil {
		return nil, err
	}

	return &ParseResult{
		Title:       article.Title,
		Description: article.MetaDescription,
		Keywords:    article.MetaKeywords,
		Content:     article.CleanedText,
		URL:         article.FinalURL,
		TopImage:    article.TopImage,
	}, nil
}
