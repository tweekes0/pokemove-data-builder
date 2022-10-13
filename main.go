package main

import (
	"log"

	"github.com/tweekes0/pokemonmoves-backend/client"
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

	// create csv files
	movesCsv, err := client.CreateFile("./data/", "moves.csv")
	handleError(err)

	pokemonCsv, err := client.CreateFile("./data/", "pokemon.csv")
	handleError(err)

	abilityCsv, err := client.CreateFile("./data/", "ability.csv")
	handleError(err)

	abilityRelCsv, err := client.CreateFile("./data", "ability-relations.csv")
	handleError(err)

	moveRelCsv, err := client.CreateFile("./data", "move-relations.csv")
	handleError(err)

	// write csv files
	err = client.ToCsv(movesCsv, moves.CsvEntries())
	handleError(err)

	err = client.ToCsv(pokemonCsv, pokemon.CsvEntries())
	handleError(err)

	err = client.ToCsv(abilityCsv, ability.CsvEntries())
	handleError(err)

	err = client.ToCsv(abilityRelCsv, ability.GetRelations())
	handleError(err)

	err = client.ToCsv(moveRelCsv, pokemon.GetRelations())
	handleError(err)
}
