package html

import (
	"bytes"
	"crawlquery/html/dto"
	"encoding/json"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c Client) GetPage(hash string) ([]byte, error) {
	resp, err := http.Get(c.BaseURL + "/pages/" + hash)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var getPageResponse dto.GetPageResponse

	if err := json.NewDecoder(resp.Body).Decode(&getPageResponse); err != nil {
		return nil, err
	}

	return getPageResponse.HTML, nil
}

func (c Client) StorePage(hash string, html []byte) error {

	storePageRequest := dto.StorePageRequest{
		Hash: hash,
		HTML: html,
	}

	body, err := json.Marshal(storePageRequest)

	if err != nil {
		return err
	}

	resp, err := http.Post(c.BaseURL+"/pages", "application/json", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil

}
