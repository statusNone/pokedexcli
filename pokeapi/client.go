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
	Pokedex map[string]Pokemon
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
	Name           string     `json:"name"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	Stats          []statInfo `json:"stats"`
	Types          []typeInfo `json:"types"`
	BaseExperience int        `json:"base_experience"`
	URL            string     `json:"url"`
}

type statInfo struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type typeInfo struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		BASEURL: "https://pokeapi.co/api/v2/",
		Cache:   *pokecache.NewCache(interval),
		Pokedex: make(map[string]Pokemon),
	}
}

// ListLocationAreas ...
func (c *Client) ListLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	url := c.BASEURL + "/location-area"
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
	url := c.BASEURL + "/location-area" + name

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

// AttemptCapture ...
func (c *Client) AttemptCapture(name string) (Pokemon, error) {
	url := c.BASEURL + "/pokemon/" + name

	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
