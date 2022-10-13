package client

import "fmt"

// struct for pokemon abilities
type Ability struct {
	Name        string
	Description string
	Generation  int
}

func (a Ability) GetHeader() []string {
	var header []string
	header = append(header, "name")
	header = append(header, "description")
	header = append(header, "generation")

	return header
}

func (a Ability) ToSlice() []string {
	var fields []string
	fields = append(fields, a.Name)
	fields = append(fields, a.Description)
	fields = append(fields, fmt.Sprintf("%v", a.Generation))

	return fields
}