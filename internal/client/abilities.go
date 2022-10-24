package client

import (
	"encoding/json"
	"fmt"
	"sync"
)

// struct for meta data of pokemon abilities
type PokemonAbilityRelation struct {
	PokeID    int
	AbilityID int
	Slot      int
	Hidden    bool
}

func (a PokemonAbilityRelation) GetHeader() []string {
	var header []string
	header = append(header, "poke-id")
	header = append(header, "ability-id")
	header = append(header, "slot")
	header = append(header, "hidden")

	return header
}

func (a PokemonAbilityRelation) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", a.PokeID))
	fields = append(fields, fmt.Sprintf("%v", a.AbilityID))
	fields = append(fields, fmt.Sprintf("%v", a.Slot))
	fields = append(fields, fmt.Sprintf("%v", a.Hidden))

	return fields
}

// struct for pokemon abilities
type PokemonAbility struct {
	AbilityID   int
	Name        string
	Description string
	Generation  int
	MainSeries  bool
	pokemon     []pokemonAbility
}

func (a PokemonAbility) GetHeader() []string {
	var header []string
	header = append(header, "ability-id")
	header = append(header, "name")
	header = append(header, "description")
	header = append(header, "generation")

	return header
}

func (a PokemonAbility) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", a.AbilityID))
	fields = append(fields, a.Name)
	fields = append(fields, a.Description)
	fields = append(fields, fmt.Sprintf("%v", a.Generation))

	return fields
}

func abilityResponseToStruct(data AbilityResponse, lang string) PokemonAbility {
	var ability PokemonAbility
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
	wg        *sync.WaitGroup
	Entries   []PokemonAbility
	Relations []PokemonAbilityRelation
}

func (a *AbilityReceiver) Init(n int) {
	a.wg = new(sync.WaitGroup)
	a.Entries = make([]PokemonAbility, n)
}

func (a *AbilityReceiver) AddWorker() {
	a.wg.Add(1)
}

func (a *AbilityReceiver) Wait() {
	a.wg.Wait()
}

func (a *AbilityReceiver) CsvEntries() []CsvEntry {
	var e []CsvEntry
	for _, entry := range a.Entries {
		e = append(e, entry)
	}

	return e
}

// Add main series abilities to entries and 
// populate the pokemon to ability relations slice
func (a *AbilityReceiver) PostProcess() {
	var ab []PokemonAbility

	for _, ability := range a.Entries {
		if ability.MainSeries {
			ab = append(ab, ability)
		}
	}

	a.Entries = ab

	for _, entry := range a.Entries {
		for _, p := range entry.pokemon {
			var meta PokemonAbilityRelation
			meta.AbilityID = entry.AbilityID
			meta.PokeID = getUrlID(p.Pokemon.Url)
			meta.Hidden = p.Hidden
			meta.Slot = p.Slot

			if meta.PokeID != -1 {
				a.Relations = append(a.Relations, meta)
			}
		}
	}
}

func (a *AbilityReceiver) FetchEntries(url, lang string, i int) {
	var resp AbilityResponse
	data, _ := getResponse(url)

	defer a.wg.Done()

	json.Unmarshal(data, &resp)

	ability := abilityResponseToStruct(resp, lang)
	a.Entries[i] = ability
}

// Gets the relationship of Ability to Pokemon
// and returns a slice of CsvEntries
func (a *AbilityReceiver) GetRelations() []CsvEntry {
	var rels []CsvEntry

	for _, rel := range a.Relations {
		rels = append(rels, rel)
	}

	return rels
}
