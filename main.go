package main

import (
	"log"
)

const (
	MoveEndpoint    = "https://pokeapi.co/api/v2/move"
	PokemonEndpoint = "https://pokeapi.co/api/v2/pokemon"
)

type APIReceiver interface {
	GetAPIData(string, string) error
}

func main() { 
	moves := MovesReceiver{}
	if err := moves.GetAPIData("en"); err != nil {
		log.Fatal(err.Error())
	}

	if err := MovesToCsv("data/moves.csv", moves.moves); err != nil {
		log.Fatal(err.Error())
	}
}
