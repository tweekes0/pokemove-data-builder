package client

import (
	"encoding/json"
	"fmt"
	"sync"
)

// struct for meta data of pokemon abilities
type AbilityMetadata struct {
	PokeID    int
	AbilityID int
	Slot      int
	Hidden    bool
}

func (a AbilityMetadata) GetHeader() []string {
	var header []string
	header = append(header, "poke-id")
	header = append(header, "ability-id")
	header = append(header, "slot")
	header = append(header, "hidden")

	return header
}

func (a AbilityMetadata) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", a.PokeID))
	fields = append(fields, fmt.Sprintf("%v", a.AbilityID))
	fields = append(fields, fmt.Sprintf("%v", a.Slot))
	fields = append(fields, fmt.Sprintf("%v", a.Hidden))

	return fields
}

// struct for pokemon abilities
type Ability struct {
	AbilityID   int
	Name        string
	Description string
	Generation  int
	MainSeries  bool
	pokemon     []pokemonAbility
}

func (a Ability) GetHeader() []string {
	var header []string
	header = append(header, "ability-id")
	header = append(header, "name")
	header = append(header, "description")
	header = append(header, "generation")

	return header
}

func (a Ability) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", a.AbilityID))
	fields = append(fields, a.Name)
	fields = append(fields, a.Description)
	fields = append(fields, fmt.Sprintf("%v", a.Generation))

	return fields
}

func abilityResponseToStruct(data AbilityResponse, lang string) Ability {
	var ability Ability
	ability.AbilityID = data.ID
	ability.Name = data.Name
	ability.MainSeries = data.MainSeries
	ability.Generation = getGeneration(data.Generation.Name)
	ability.Description = getFlavorText(
		ability.Generation,
		lang,
		data.FlavorTexts,
	)
	ability.pokemon = data.Pokemon

	return ability
}

type AbilityReceiver struct {
	wg      *sync.WaitGroup
	entries []Ability
}

func (a *AbilityReceiver) Init(n int) {
	a.wg = new(sync.WaitGroup)
	a.entries = make([]Ability, n)
}

func (a *AbilityReceiver) AddWorker() {
	a.wg.Add(1)
}

func (a *AbilityReceiver) Wait() {
	a.wg.Wait()
}

func (a *AbilityReceiver) CsvEntries() []CsvEntry {
	var e []CsvEntry
	for _, entry := range a.entries {
		e = append(e, entry)
	}

	return e
}

func (a *AbilityReceiver) PostProcess() {
	var ab []Ability

	for _, ability := range a.entries {
		if ability.MainSeries {
			ab = append(ab, ability)
		}
	}

	a.entries = ab
}

func (a *AbilityReceiver) GetEntries(url, lang string, i int) {
	var resp AbilityResponse
	data, _ := getResponse(url)

	defer a.wg.Done()

	json.Unmarshal(data, &resp)

	ability := abilityResponseToStruct(resp, lang)
	a.entries[i] = ability
}

// Gets the relationship of Ability to Pokemon
// and returns a slice of CsvEntries
func (a *AbilityReceiver) GetRelations() []CsvEntry {
	var rels []CsvEntry

	for _, a := range a.entries {
		for _, p := range a.pokemon {
			meta := AbilityMetadata{}
			meta.AbilityID = a.AbilityID
			meta.PokeID = getUrlID(p.Pokemon.Url)
			meta.Hidden = p.Hidden
			meta.Slot = p.Slot

			if meta.PokeID != -1 {
				rels = append(rels, meta)
			}
		}
	}

	return rels
}
