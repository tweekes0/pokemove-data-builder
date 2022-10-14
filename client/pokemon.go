package client

import (
	"encoding/json"
	"fmt"
	"sync"
)

type PokemonMoveMetadata struct {
	PokeID      int
	MoveID      int
	Generation  int
	LearnLevel  int
	LearnMethod string
	GameName    string
}

func (p PokemonMoveMetadata) GetHeader() []string {
	var header []string
	header = append(header, "poke-id")
	header = append(header, "move-id")
	header = append(header, "generation")
	header = append(header, "level-learned")
	header = append(header, "learn-method")
	header = append(header, "game")

	return header
}

func (p PokemonMoveMetadata) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", p.PokeID))
	fields = append(fields, fmt.Sprintf("%v", p.MoveID))
	fields = append(fields, fmt.Sprintf("%v", p.Generation))
	fields = append(fields, fmt.Sprintf("%v", p.LearnLevel))
	fields = append(fields, p.LearnMethod)
	fields = append(fields, p.GameName)

	return fields
}

// struct for pokemon
type Pokemon struct {
	PokeID  int
	Name    string
	Sprite  string
	Species string
	Moves   []move
}

func (p Pokemon) GetHeader() []string {
	var header []string
	header = append(header, "poke-id")
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
	p.Moves = data.Moves

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
	var resp PokemonResponse
	data, _ := getResponse(url)

	defer p.wg.Done()

	json.Unmarshal(data, &resp)

	pokemon := pokemonResponseToStruct(resp, lang)

	p.entries[i] = pokemon
}

// Gets the relationship of Move to Pokemon
// and returns a slice of CsvEntries
func (p *PokemonReceiver) GetRelations() []CsvEntry {
	var rels []CsvEntry

	for _, pokemon := range p.entries {
		for _, m := range pokemon.Moves {
			for _, detail := range m.Details {
				meta := PokemonMoveMetadata{}
				meta.PokeID = pokemon.PokeID
				meta.MoveID = getUrlID(m.Name.Url)
				meta.Generation = resolveVersionGroup(detail.VersionGroup.Url)
				meta.LearnLevel = detail.LevelLearned
				meta.LearnMethod = detail.MethodLearned.Name
				meta.GameName = detail.VersionGroup.Name

				rels = append(rels, meta)
			}
		}
	}

	return rels
}
