package node

import (
	"bytes"
	"context"
	"crawlquery/node/dto"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Client struct {
	hostname string
	port     uint
	context  context.Context
}

type Option func(*Client)

func WithHostname(hostname string) Option {
	return func(c *Client) {
		c.hostname = hostname
	}
}

func WithPort(port uint) Option {
	return func(c *Client) {
		c.port = port
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *Client) {
		c.context = ctx
	}
}

func NewClient(opts ...Option) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	if c.context == nil {
		c.context = context.Background()
	}

	return c
}

func (c *Client) buildUrl(path string) string {
	return fmt.Sprintf("http://%s:%d%s", c.hostname, c.port, path)
}

func (c *Client) SendRequest(method, path string, body []byte) (*http.Response, error) {

	req, err := http.NewRequestWithContext(c.context, method, c.buildUrl(path), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func (c *Client) Crawl(pageID, url string) (*dto.CrawlResponse, error) {
	req := dto.CrawlRequest{
		PageID: pageID,
		URL:    url,
	}

	jsonBody, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	res, err := c.SendRequest("POST", "/crawl", jsonBody)

	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {

		var errRes dto.ErrorResponse

		err = json.NewDecoder(res.Body).Decode(&errRes)

		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		defer res.Body.Close()

		return nil, fmt.Errorf("unexpected status code: %d (%s)", res.StatusCode, errRes.Error)
	}

	var crawlRes dto.CrawlResponse
	err = json.NewDecoder(res.Body).Decode(&crawlRes)

	if err != nil {
		return nil, err
	}

	return &crawlRes, nil
}

func (c *Client) Index(pageID string, url string, contentHash string) error {

	req := dto.IndexRequest{
		PageID:      pageID,
		URL:         url,
		ContentHash: contentHash,
	}

	jsonBody, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := c.SendRequest("POST", "/index", jsonBody)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {

		var errRes dto.ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return fmt.Errorf("unexpected status code: %d (%s)", res.StatusCode, errRes.Error)
		}

		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var indexRes dto.IndexResponse

	if err := json.NewDecoder(res.Body).Decode(&indexRes); err != nil {
		return err
	}

	if indexRes.Success {
		return nil
	}

	return errors.New("indexing returned success=false")
}

func (c *Client) GetIndexMetas(pageIDs []string) ([]dto.IndexMeta, error) {

	var req dto.GetIndexMetasRequest

	req.PageIDs = append(req.PageIDs, pageIDs...)

	jsonBody, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	res, err := c.SendRequest("POST", "/repair/get-index-metas", jsonBody)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var indexMetaResponse dto.GetIndexMetasResponse

	if err := json.NewDecoder(res.Body).Decode(&indexMetaResponse); err != nil {
		return nil, err
	}

	return indexMetaResponse.IndexMetas, nil
}

func (c *Client) GetAllIndexMetas() ([]dto.IndexMeta, error) {

	res, err := c.SendRequest("GET", "/repair/get-all-index-metas", nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var indexMetaResponse dto.GetIndexMetasResponse

	if err := json.NewDecoder(res.Body).Decode(&indexMetaResponse); err != nil {
		return nil, err
	}

	return indexMetaResponse.IndexMetas, nil
}

func (c *Client) GetPageDumps(pageIDs []string) ([]dto.PageDump, error) {

	var req dto.GetPageDumpsRequest

	req.PageIDs = append(req.PageIDs, pageIDs...)

	jsonBody, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	res, err := c.SendRequest("POST", "/repair/get-page-dumps", jsonBody)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var pageDumpResponse dto.GetPageDumpsResponse

	if err := json.NewDecoder(res.Body).Decode(&pageDumpResponse); err != nil {
		return nil, err
	}

	return pageDumpResponse.PageDumps, nil
}
