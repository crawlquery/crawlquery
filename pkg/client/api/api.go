package api

import (
	"bytes"
	"crawlquery/api/dto"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Client struct {
	BaseURL string
	logger  *zap.SugaredLogger
}

func NewClient(baseURL string, logger *zap.SugaredLogger) *Client {
	return &Client{
		BaseURL: baseURL,
		logger:  logger,
	}
}

func (c *Client) ListNodesByShardID(shardID uint16) ([]*dto.PublicNode, error) {

	endpoint := fmt.Sprintf("%s/shards/%d/nodes", c.BaseURL, shardID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		c.logger.Errorf("error creating request: %v", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Errorf("error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("error response: %v", resp.Status)
		return nil, errors.New("could not list nodes")
	}

	var listRes dto.ListNodesByShardResponse
	if err := json.NewDecoder(resp.Body).Decode(&listRes); err != nil {
		c.logger.Errorf("error decoding response: %v", err)
		return nil, err
	}

	return listRes.Nodes, nil
}

func (c *Client) AuthenticateNode(key string) (*dto.Node, error) {

	authenticateNodeRequest := &dto.AuthenticateNodeRequest{
		Key: key,
	}

	encoded, err := json.Marshal(authenticateNodeRequest)

	if err != nil {
		c.logger.Errorf("error encoding request: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/auth/node", bytes.NewBuffer(encoded))
	if err != nil {
		c.logger.Errorf("error creating request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Errorf("error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("error response: %v", resp.Status)
		return nil, errors.New("could not authenticate node")
	}

	var authRes dto.AuthenticateNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&authRes); err != nil {
		c.logger.Errorf("error decoding response: %v", err)
		return nil, err
	}

	if authRes.Node.ID == "" {
		return nil, errors.New("could not authenticate node")
	}

	return authRes.Node, nil
}

func (c *Client) CreateCrawlJob(url string) (*dto.CrawlJob, error) {

	crawlRequest := &dto.CreateCrawlJobRequest{
		URL: url,
	}

	encoded, err := json.Marshal(crawlRequest)

	if err != nil {
		c.logger.Errorf("error encoding request: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/crawl-jobs", bytes.NewBuffer(encoded))
	if err != nil {
		c.logger.Errorf("error creating request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Errorf("error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("error response: %v", resp.Status)
		return nil, errors.New("could not request crawl")
	}

	var crawlRes dto.CreateCrawlJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&crawlRes); err != nil {
		c.logger.Errorf("error decoding response: %v", err)
		return nil, err
	}

	return &crawlRes.CrawlJob, nil
}

func (c *Client) CreateLink(src, dst string) error {

	linkRequest := &dto.CreateLinkRequest{
		Src: src,
		Dst: dst,
	}

	encoded, err := json.Marshal(linkRequest)

	if err != nil {
		c.logger.Errorf("error encoding request: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/links", bytes.NewBuffer(encoded))
	if err != nil {
		c.logger.Errorf("error creating request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Errorf("error sending request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		c.logger.Errorf("error response: %v", resp.Status)
		return errors.New("could not create link")
	}

	return nil
}
