package client

import (
	"encoding/json"
	"fmt"
	"sync"
)

// struct for pokemon
type Pokemon struct {
	PokeID  int
	Name    string
	Sprite  string
	Species string
}

func (p Pokemon) GetHeader() []string {
	var header []string
	header = append(header, "pokeID")
	header = append(header, "name")
	header = append(header, "sprite")
	header = append(header, "species")

	return header
}

func (p Pokemon) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", p.PokeID))
	fields = append(fields, p.Name)
	fields = append(fields, p.Species)
	fields = append(fields, p.Sprite)

	return fields
}

// struct that receives data from the pokeapi pokemon endpoint
func pokemonResponseToStruct(data PokemonResponse, lang string) Pokemon {
	var p Pokemon
	p.PokeID = data.ID
	p.Name = data.Name
	p.Species = data.Species.Name
	p.Sprite = data.Sprite.Other["official-artwork"].FrontDefault

	return p
}

type PokemonReceiver struct {
	wg      *sync.WaitGroup
	entries []Pokemon
}

func (p *PokemonReceiver) Init(n int) {
	p.wg = new(sync.WaitGroup)
	p.entries = make([]Pokemon, n)
}

func (p *PokemonReceiver) AddWorker() {
	p.wg.Add(1)
}

func (p *PokemonReceiver) Wait() {
	p.wg.Wait()
}

func (p *PokemonReceiver) PostProcess() {}

func (p *PokemonReceiver) CsvEntries() []CsvEntry {
	var e []CsvEntry
	for _, entry := range p.entries {
		e = append(e, entry)
	}

	return e
}

func (p *PokemonReceiver) GetEntries(url, lang string, i int) {
	resp := PokemonResponse{}
	data, _ := getResponse(url)

	defer p.wg.Done()

	json.Unmarshal(data, &resp)

	pokemon := pokemonResponseToStruct(resp, lang)

	p.entries[i] = pokemon
}
