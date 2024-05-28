package node

import (
	"bytes"
	"crawlquery/node/dto"
	"encoding/json"
	"errors"
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

func (c *Client) Crawl(pageID, url string) (*dto.Page, error) {
	req := dto.CrawlRequest{
		PageID: pageID,
		URL:    url,
	}

	jsonBody, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	res, err := http.Post(
		c.BaseURL+"/crawl",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusBadRequest {

		var errRes dto.ErrorResponse

		err = json.NewDecoder(res.Body).Decode(&errRes)

		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		defer res.Body.Close()

		return nil, fmt.Errorf("unexpected status code: %d (%s)", res.StatusCode, errRes.Error)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var crawlRes dto.CrawlResponse
	err = json.NewDecoder(res.Body).Decode(&crawlRes)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return crawlRes.Page, nil
}

func (c *Client) Index(pageID string) error {

	endpoint := fmt.Sprintf("%s/pages/%s/index", c.BaseURL, pageID)
	res, err := http.Post(
		endpoint,
		"application/json",
		nil,
	)

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

	res, err := http.Post(fmt.Sprintf("%s/repair/get-index-metas", c.BaseURL), "application/json", bytes.NewBuffer(jsonBody))

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

func (c *Client) GetAllIndexMetas() ([]dto.IndexMeta, error) {

	res, err := http.Get(fmt.Sprintf("%s/repair/get-all-index-metas", c.BaseURL))

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

	res, err := http.Post(fmt.Sprintf("%s/repair/get-page-dumps", c.BaseURL), "application/json", bytes.NewBuffer(jsonBody))

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
