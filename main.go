package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/pokeapi"
	"strings"
)

type config struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

type cliCommand struct {
	Name        string
	Description string
	Callback    func(c *config, args []string) error
}

var cfg *config
var commands map[string]cliCommand
var client *pokeapi.Client

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Displays the names of 20 location areas in the Pokemon World",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous 20 location areas",
			Callback:    commandMapb,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore an area of the map",
			Callback:    commandExplore,
		},
	}
}

func init() {
	client = pokeapi.NewClient()
	baseURL := client.BASEURL
	cfg = &config{
		Next:     &baseURL,
		Previous: nil,
	}
	commands = getCommands()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		words := cleanInput(input)

		if len(words) > 0 {
			if command, ok := commands[words[0]]; ok {
				err := command.Callback(cfg, words[1:])
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

func commandHelp(c *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for cmdName, cmd := range commands {
		fmt.Printf("%s: %s\n", cmdName, cmd.Description)
	}
	return nil
}

func commandExit(c *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(c *config, args []string) error {
	resp, err := client.ListLocationAreas(c.Next)
	if err != nil {
		return err
	}

	c.Next = resp.Next
	c.Previous = resp.Previous

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapb(c *config, args []string) error {
	if c.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	resp, err := client.ListLocationAreas(c.Previous)
	if err != nil {
		return err
	}

	c.Next = resp.Next
	c.Previous = resp.Previous

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandExplore(c *config, args []string) error {

	locationAreaName := args[0]
	fmt.Printf("Exploring %s...\n", locationAreaName)

	locationArea, err := client.GetLocationArea(locationAreaName)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}

	return nil
}
