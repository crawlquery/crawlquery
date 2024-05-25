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
		c.BaseURL+"/crawl",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return "", err
	}

	if res.StatusCode == http.StatusBadRequest {

		var errRes dto.ErrorResponse

		err = json.NewDecoder(res.Body).Decode(&errRes)

		if err != nil {
			return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		defer res.Body.Close()

		return "", fmt.Errorf("unexpected status code: %d (%s)", res.StatusCode, errRes.Error)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var crawlRes dto.CrawlResponse
	err = json.NewDecoder(res.Body).Decode(&crawlRes)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	return crawlRes.PageHash, nil
}
