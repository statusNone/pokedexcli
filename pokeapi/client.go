package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

// Client ...
type Client struct {
	BASEURL string
}

// LocationAreaResponse ...
type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		BASEURL: "https://pokeapi.co/api/v2/location-area",
	}
}

// ListLocationAreas ...
func (c *Client) ListLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	url := c.BASEURL
	if pageURL != nil {
		url = *pageURL
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	var result LocationAreaResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	return result, nil
}
