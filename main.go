package main

import (
	"log"
)

const (
	MoveEndpoint    = "https://pokeapi.co/api/v2/move"
	PokemonEndpoint = "https://pokeapi.co/api/v2/pokemon"
)

type APIReceiver interface {
	AddWorker()
	FlattenEntries()
	Init(int)
	GetEntries(string, string, int)
	Wait()
	CsvEntries() []CsvEntry
}

func GetAPIData(recv APIReceiver, limit int, endpoint, lang string) error {
	basicResp, err := getBasicResponse(limit, MoveEndpoint)
	if err != nil {
		return err
	}

	recv.Init(basicResp.Count) 

	for i := 0; i < basicResp.Count; i++ {
		recv.AddWorker()
		go recv.GetEntries(basicResp.Results[i].Url, lang, i)
	}

	recv.Wait()
	recv.FlattenEntries()
	return nil
}

func main() { 
	moves := MovesReceiver{}
	// pokemon := PokemonReceiver{}
	lang := "en"
	limit := 1000

	if err := GetAPIData(&moves, limit, MoveEndpoint, lang); err != nil {
		log.Fatal(err.Error())
	}

	// if err := GetAPIData(&pokemon, limit, PokemonEndpoint, lang); err != nil {
	// 	log.Fatal(err.Error())
	// }

	movesCsv, err := createCsv("./data/moves.csv")
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = ToCsv(movesCsv, &moves); err != nil {
		log.Fatal(err.Error())
	}

}
