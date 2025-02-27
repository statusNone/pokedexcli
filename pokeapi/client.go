package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"pokedexcli/pokecache"
	"time"
)

const interval = 5 * time.Second

// Client ...
type Client struct {
	BASEURL string
	Cache   pokecache.Cache
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

// LocationArea ...
type LocationArea struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

// PokemonEncounter ...
type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

// Pokemon ...
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		BASEURL: "https://pokeapi.co/api/v2/location-area",
		Cache:   *pokecache.NewCache(interval),
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

// GetLocationArea ...
func (c *Client) GetLocationArea(name string) (LocationArea, error) {
	url := c.BASEURL + "/" + name

	if cached, ok := c.Cache.Get(url); ok {
		var locationArea LocationArea
		err := json.Unmarshal(cached, &locationArea)
		if err != nil {
			return LocationArea{}, err
		}
		return locationArea, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationArea{}, err
	}

	c.Cache.Add(url, body)

	var locationArea LocationArea
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return LocationArea{}, err
	}

	return locationArea, nil
}
