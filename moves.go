package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// struct for pokemon move models
type PokemonMove struct {
	MoveID      int
	Accuracy    int
	Power       int
	PowerPoints int
	Generation  int
	Name        string
	Type        string
	DamageType  string
	Description string
}

func (p PokemonMove) GetHeader() []string {
	var header []string
	header = append(header, "moveID")
	header = append(header, "accuracy")
	header = append(header, "power")
	header = append(header, "pp")
	header = append(header, "generation")
	header = append(header, "name")
	header = append(header, "type")
	header = append(header, "damage-type")
	header = append(header, "description")

	return header
}

func (p PokemonMove) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", p.MoveID))
	fields = append(fields, fmt.Sprintf("%v", p.Accuracy))
	fields = append(fields, fmt.Sprintf("%v", p.Power))
	fields = append(fields, fmt.Sprintf("%v", p.PowerPoints))
	fields = append(fields, fmt.Sprintf("%v", p.Generation))
	fields = append(fields, p.Name)
	fields = append(fields, p.Type)
	fields = append(fields, p.DamageType)
	fields = append(fields, p.Description)

	return fields
}

// api receive for the /moves endpoint
type MovesReceiver struct {
	// a slice of slices since the number of moves per response is variable
	entryMatrix [][]PokemonMove
	entries     []PokemonMove
	wg          *sync.WaitGroup
}

func (m *MovesReceiver) Init(n int) {
	m.wg = new(sync.WaitGroup)
	m.entryMatrix = make([][]PokemonMove, n)
}

func (m *MovesReceiver) AddWorker() {
	m.wg.Add(1)
}

func (m *MovesReceiver) Wait() {
	m.wg.Wait()
}

func (m *MovesReceiver) CsvEntries() []CsvEntry {
	var e []CsvEntry
	for _, entry := range m.entries {
		e = append(e, entry)
	}

	return e
}

func (m *MovesReceiver) GetEntries(url, lang string, i int) {
	resp := MoveResponse{}
	moves := []PokemonMove{}
	data, _ := getResponse(url)

	defer m.wg.Done()

	json.Unmarshal(data, &resp)

	gen := getGeneration(resp.Generation.Name)
	if len(resp.PastValues) > 0 {
		for _, value := range resp.PastValues {
			oldMove := moveResponseToStruct(resp, lang)

			if value.Accuracy != 0 {
				oldMove.Accuracy = value.Accuracy
			}

			if value.Power != 0 {
				oldMove.Power = value.Power
			}

			if value.PowerPoints != 0 {
				oldMove.PowerPoints = value.PowerPoints
			}

			if value.Type.Name != "" {
				oldMove.Type = value.Type.Name
			}

			oldMove.Generation = gen
			oldMove.Description = getFlavorText(gen, lang, resp.FlavorTexts)
			gen = resolveVersionGroup(value.VersionGroup.Url)

			moves = append(moves, oldMove)
		}
	}

	move := moveResponseToStruct(resp, lang)
	move.Generation = gen
	move.Description = getFlavorText(gen, lang, resp.FlavorTexts)

	moves = append(moves, move)

	m.entryMatrix[i] = moves
}

func (m *MovesReceiver) FlattenEntries() {
	for _, entry := range m.entryMatrix {
		m.entries = append(m.entries, entry...)
	}
}
