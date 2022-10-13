package client

import (
	"encoding/json"
	"fmt"
	"sync"
)

// struct for meta data of pokemon abilities
type AbilityMeta struct {
	PokeID    int
	AbilityId int
	Slot      int
	Hidden    bool
}

// struct for pokemon abilities
type Ability struct {
	ID          int
	Name        string
	Description string
	Generation  int
	MainSeries  bool
}

func (a Ability) GetHeader() []string {
	var header []string
	header = append(header, "id")
	header = append(header, "name")
	header = append(header, "description")
	header = append(header, "generation")

	return header
}

func (a Ability) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", a.ID))
	fields = append(fields, a.Name)
	fields = append(fields, a.Description)
	fields = append(fields, fmt.Sprintf("%v", a.Generation))

	return fields
}

func abilityResponseToStruct(data AbilityResponse, lang string) Ability {
	var ability Ability
	ability.ID = data.ID
	ability.Name = data.Name
	ability.Generation = getGeneration(data.Generation.Name)
	ability.Description = getFlavorText(
		ability.Generation,
		lang,
		data.FlavorTexts,
	)

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

func (a *AbilityReceiver) FlattenEntries() {}

func (a *AbilityReceiver) GetEntries(url, lang string, i int) {
	resp := AbilityResponse{}
	data, _ := getResponse(url)

	defer a.wg.Done()

	json.Unmarshal(data, &resp)

	ability := abilityResponseToStruct(resp, lang)
	a.entries[i] = ability
}
