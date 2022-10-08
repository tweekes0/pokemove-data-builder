package main

import "fmt"

type CsvEntry interface {
	GetHeader() []string
	ToSlice() []string
}

// struct for pokemon move models
type PokeMove struct {
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

func (p PokeMove) GetHeader() []string {
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

func (p PokeMove) ToSlice() []string {
	var fields []string
	fields = append(fields, fmt.Sprintf("%v", p.MoveID))
	fields = append(fields, fmt.Sprintf("%v", p.Accuracy))
	fields = append(fields, fmt.Sprintf("%v", p.Power))
	fields = append(fields, fmt.Sprintf("%v", p.PowerPoints))
	fields = append(fields, fmt.Sprintf("%v",p.Generation))
	fields = append(fields, p.Name)
	fields = append(fields, p.Type)
	fields = append(fields, p.DamageType)
	fields = append(fields, p.Description)

	return fields
}
