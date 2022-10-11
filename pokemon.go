package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// struct for pokemon
type Pokemon struct {
	PokeID int
	Name   string
	Sprite string
}

func (p Pokemon) GetHeader() []string {
	var header []string
	header = append(header, "pokeID")
	header = append(header, "name")
	header = append(header, "sprite")

	return header
}

func (p Pokemon) ToSlice() []string {
	var s []string 
	s = append(s, fmt.Sprintf("%v", p.PokeID))
	s = append(s, p.Name)
	s = append(s, p.Sprite)

	return s
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

func (p *PokemonReceiver) FlattenEntries() { }

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