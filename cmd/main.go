package main

import (
	"log"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
)

const (
	AbilityEndpoint = "https://pokeapi.co/api/v2/ability"
	MoveEndpoint    = "https://pokeapi.co/api/v2/move"
	PokemonEndpoint = "https://pokeapi.co/api/v2/pokemon"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	ability := client.AbilityReceiver{}
	moves := client.MovesReceiver{}
	pokemon := client.PokemonReceiver{}

	lang := "en"
	limit := 2000

	// fetch api data
	err := client.GetAPIData(&moves, limit, MoveEndpoint, lang)
	handleError(err)

	err = client.GetAPIData(&pokemon, limit, PokemonEndpoint, lang)
	handleError(err)

	err = client.GetAPIData(&ability, limit, AbilityEndpoint, lang)
	handleError(err)

	// Generate CSV files of fetched API data
	// client.generateCsvs(pokemon, moves, ability)
}
