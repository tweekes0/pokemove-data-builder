package client

import (
	"encoding/json"
	"fmt"
	"sync"
)

const (
	CurrentGen = 9
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
	PokeID        int    `json:"poke_id"`
	OriginGen     int    `json:"generation"`
	GenTypeChange int    `json:"-"`
	Name          string `json:"name"`
	Sprite        string `json:"sprite"`
	ShinySprite   string `json:"shiny_sprite"`
	Species       string `json:"species"`
	PrimaryType   string `json:"primary_type"`
	SecondaryType string `json:"secondary_type,omitempty"`
	Moves         []move `json:"-"`
}

type PokemonBrief struct {
	PokeID int    `json:"poke_id"`
	Name   string `json:"name"`
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
func pokemonResponseToStruct(data PokemonResponse, lang string) []Pokemon {
	var pokes []Pokemon

	if len(data.PastTypes) > 0 {
		for _, t := range data.PastTypes {
			var pp Pokemon
			pp.PokeID = data.ID
			pp.OriginGen = getOriginGeneration(data.ID)
			pp.Name = data.Name
			pp.Species = data.Species.Name
			pp.Sprite = data.Sprite.Other["official-artwork"].FrontDefault
			pp.ShinySprite = data.Sprite.Other["official-artwork"].FrontShiny
			pp.Moves = data.Moves
			pp.GenTypeChange = getGeneration(t.Generation.Name)
			pp.PrimaryType = t.Types[0].Type.Name
			pp.SecondaryType = ""

			if len(t.Types) > 1 {
				pp.SecondaryType = t.Types[1].Type.Name
			}

			pokes = append(pokes, pp)
		}
	}

	var p Pokemon
	p.PokeID = data.ID
	p.OriginGen = getOriginGeneration(data.ID)
	p.Name = data.Name
	p.Species = data.Species.Name
	p.Sprite = data.Sprite.Other["official-artwork"].FrontDefault
	p.ShinySprite = data.Sprite.Other["official-artwork"].FrontShiny
	p.Moves = data.Moves
	p.GenTypeChange = CurrentGen
	p.PrimaryType = data.Types[0].Type.Name
	p.SecondaryType = ""

	if len(data.Types) > 1 {
		p.SecondaryType = data.Types[1].Type.Name
	}

	pokes = append(pokes, p)
	return pokes
}

type PokemonReceiver struct {
	wg          *sync.WaitGroup
	entryMatrix [][]Pokemon
	Entries     []Pokemon
	Relations   []PokemonMoveRelation
	Endpoint    string
}

func (p *PokemonReceiver) Init(n int) {
	p.wg = new(sync.WaitGroup)
	p.entryMatrix = make([][]Pokemon, n)
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
	for _, entry := range p.entryMatrix {
		p.Entries = append(p.Entries, entry...)
	}

	prev := 0

	for _, pokemon := range p.Entries {
		if prev == pokemon.PokeID {
			continue
		}

		prev = pokemon.PokeID

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

	p.entryMatrix[i] = pokemon
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
