package api

import (
	"bytes"
	"crawlquery/api/dto"
	"encoding/json"
	"errors"
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
