package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// struct for pokeapi NamedAPIResource
// https://pokeapi.co/docs/v2#namedapiresource
type namedResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// struct for pokeapi flavor text entries
//https://pokeapi.co/docs/v2#moveflavortext
type flavorText struct {
	Text         string        `json:"flavor_text"`
	Language     namedResource `json:"language"`
	VersionGroup namedResource `json:"version_group"`
}

// struct for pokeapi pastMoveValues
// https://pokeapi.co/docs/v2#moves
type pastMoveValue struct {
	Accuracy      int             `json:"accuracy"`
	EffectChance  int             `json:"effect_chance"`
	Power         int             `json:"power"`
	PowerPoints   int             `json:"pp"`
	Type          namedResource   `json:"type"`
	VersionGroup  namedResource   `json:"version_group"`
}

// struct for pokeapi endpoints when not supplied with
// a parameter
type BasicResponse struct {
	Count   int             `json:"count"`
	Results []namedResource `json:"results"`
}

func (r *BasicResponse) Print() {
	fmt.Printf("Count: %v\n", r.Count)
	fmt.Println("Results:")
	for _, r := range r.Results {
		fmt.Printf("\tName:%v\n", r.Name)
		fmt.Printf("\tUrl:%v\n", r.Url)
	}
}

// struct for pokeapi Moves endpoint response
type MoveResponse struct {
	ID          int             `json:"id"`
	Accuracy    int             `json:"accuracy"`
	Power       int             `json:"power"`
	PowerPoints int             `json:"pp"`
	Name        string          `json:"name"`
	DamageType  namedResource   `json:"damage_class"`
	Type        namedResource   `json:"type"`
	Generation  namedResource   `json:"generation"`
	FlavorTexts []flavorText    `json:"flavor_text_entries"`
	PastValues  []pastMoveValue `json:"past_values"`
}

func (r *MoveResponse) Print() {
	fmt.Printf("Move ID: %v\n", r.ID)
	fmt.Printf("Move Accuracy: %v\n", r.Accuracy)
	fmt.Printf("Move Power: %v\n", r.Power)
	fmt.Printf("Move Name: %v\n", r.Name)
	fmt.Printf("Move Type: %v\n", r.Type)
	fmt.Printf("Move Generation: %v\n", r.Generation)
	fmt.Println()
}

func getResponse(url string) ([]byte, error) {
	var data []byte

	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	return data, nil
}

func getBasicMoveResponse(limit int, endpoint string) (BasicResponse, error) {
	var basicResp BasicResponse

	url := fmt.Sprintf("%v?limit=%v", endpoint, limit)

	data, err := getResponse(url)
	if err != nil {
		return basicResp, err
	}
	json.Unmarshal(data, &basicResp)

	return basicResp, nil
}
