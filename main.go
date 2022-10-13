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

func main() { 
	ability := client.AbilityReceiver{}
	moves := client.MovesReceiver{}
	pokemon := client.PokemonReceiver{}

	lang := "en"
	limit := 2000

	if err := client.GetAPIData(&moves, limit, MoveEndpoint, lang); err != nil {
		log.Fatal(err.Error())
	}

	if err := client.GetAPIData(&pokemon, limit, PokemonEndpoint, lang); err != nil {
		log.Fatal(err.Error())
	}

	if err := client.GetAPIData(&ability, limit, AbilityEndpoint, lang); err != nil {
		log.Fatal(err.Error())
	}

	movesCsv, err := client.CreateFile("./data/", "moves.csv")
	if err != nil {
		log.Fatal(err.Error())
	}

	pokemonCsv, err := client.CreateFile("./data/", "pokemon.csv")
	if err != nil {
		log.Fatal(err.Error())
	}

	abilityCsv, err := client.CreateFile("./data/", "ability.csv")
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = client.ToCsv(movesCsv, moves.CsvEntries()); err != nil {
		log.Fatal(err.Error())
	}

	if err = client.ToCsv(pokemonCsv, pokemon.CsvEntries()); err != nil {
		log.Fatal(err.Error())
	}

	if err = client.ToCsv(abilityCsv, ability.CsvEntries()); err != nil {
		log.Fatal(err.Error())
	}
}
