package service

import "crawlquery/pkg/client/html"

type Service struct {
	htmlClient *html.Client
}

func NewService(htmlClient *html.Client) *Service {
	return &Service{
		htmlClient: htmlClient,
	}
}

func (s *Service) Save(pageID string, html []byte) error {
	return s.htmlClient.StorePage(pageID, html)
}

func (s *Service) Get(pageID string) ([]byte, error) {
	return s.htmlClient.GetPage(pageID)
}
