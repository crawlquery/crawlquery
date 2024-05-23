package node

import (
	"bytes"
	"crawlquery/node/dto"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

func (c *Client) Crawl(pageID, url string) (string, error) {
	req := dto.CrawlRequest{
		PageID: pageID,
		URL:    url,
	}

	jsonBody, err := json.Marshal(req)

	if err != nil {
		return "", err
	}

	res, err := http.Post(
		c.BaseURL+"/crawl-jobs",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return "", err
	}

	var crawlRes dto.CrawlResponse
	err = json.NewDecoder(res.Body).Decode(&crawlRes)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return crawlRes.PageHash, nil
}
