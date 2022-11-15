package client

import (
	"encoding/json"
	"fmt"
	"sync"
)

type PokemonMoveRelation struct {
	PokeID       int
	MoveID       int
	Generation   int
	LevelLearned int
	LearnMethod  string
	GameName     string
}

func (p PokemonMoveRelation) GetHeader() []string {
	var header []string
	header = append(header, "poke-id")
	header = append(header, "move-id")
	header = append(header, "generation")
	header = append(header, "level-learned")
	header = append(header, "learn-method")
	header = append(header, "game")

	return header
}

func (p PokemonMoveRelation) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", p.PokeID))
	fields = append(fields, fmt.Sprintf("%v", p.MoveID))
	fields = append(fields, fmt.Sprintf("%v", p.Generation))
	fields = append(fields, fmt.Sprintf("%v", p.LevelLearned))
	fields = append(fields, p.LearnMethod)
	fields = append(fields, p.GameName)

	return fields
}

// struct for pokemon
type Pokemon struct {
	PokeID    int    `json:"poke_id"`
	OriginGen int    `json:"generation"`
	Name      string `json:"name"`
	Sprite    string `json:"sprite"`
	Species   string `json:"species"`
	Moves     []move `json:"-"`
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
	p.OriginGen = getOriginGeneration(data.ID)
	p.Name = data.Name
	p.Species = data.Species.Name
	p.Sprite = data.Sprite.Other["official-artwork"].FrontDefault
	p.Moves = data.Moves

	return p
}

type PokemonReceiver struct {
	wg        *sync.WaitGroup
	Entries   []Pokemon
	Relations []PokemonMoveRelation
	Endpoint  string
}

func (p *PokemonReceiver) Init(n int) {
	p.wg = new(sync.WaitGroup)
	p.Entries = make([]Pokemon, n)
}

func (p *PokemonReceiver) AddWorker() {
	p.wg.Add(1)
}

func (p *PokemonReceiver) Wait() {
	p.wg.Wait()
}

func (p *PokemonReceiver) GetEndpoint() string {
	return p.Endpoint
}

// When all data is fetched from api populate
// pokemon to move relationship slice
func (p *PokemonReceiver) PostProcess() {
	for _, pokemon := range p.Entries {
		for _, m := range pokemon.Moves {
			for _, detail := range m.Details {
				var meta PokemonMoveRelation
				meta.PokeID = pokemon.PokeID
				meta.MoveID = getUrlID(m.Name.Url)
				meta.Generation = resolveVersionGroup(detail.VersionGroup.Url)
				meta.LevelLearned = detail.LevelLearned
				meta.LearnMethod = detail.MethodLearned.Name
				meta.GameName = detail.VersionGroup.Name

				p.Relations = append(p.Relations, meta)
			}
		}
	}
}

func (p *PokemonReceiver) CsvEntries() []CsvEntry {
	var e []CsvEntry
	for _, entry := range p.Entries {
		e = append(e, entry)
	}

	return e
}

func (p *PokemonReceiver) FetchEntries(url, lang string, i int) {
	var resp PokemonResponse
	data, _ := getResponse(url)

	defer p.wg.Done()

	json.Unmarshal(data, &resp)

	pokemon := pokemonResponseToStruct(resp, lang)

	p.Entries[i] = pokemon
}

// Gets the relationship of Move to Pokemon
// and returns a slice of CsvEntries
func (p *PokemonReceiver) GetCsvRelations() []CsvEntry {
	var rels []CsvEntry

	for _, rel := range p.Relations {
		rels = append(rels, rel)
	}

	return rels
}

func (p *PokemonReceiver) GetEntries() []interface{} {
	var entries []interface{}

	for _, e := range p.Entries {
		entries = append(entries, e)
	}

	return entries
}

func (p *PokemonReceiver) GetRelations() []interface{} {
	var rels []interface{}

	for _, r := range p.Relations {
		rels = append(rels, r)
	}

	return rels
}