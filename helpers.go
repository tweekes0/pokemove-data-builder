package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// resolves the pokeapi version group to a generation number
// https://pokeapi.co/docs/v2#versiongroup
func resolveVersionGroup(url string) int {
	group_id := strings.Split(url, "/")[6]
	id, err := strconv.Atoi(group_id)
	if err != nil {
		return -1
	}

	switch id {
	case 1, 2:
		return 1
	case 3, 4:
		return 2
	case 5, 6, 7, 12, 13:
		return 3
	case 8, 9, 10:
		return 4
	case 11, 14:
		return 5
	case 15, 16:
		return 6
	case 17, 18, 19:
		return 7
	case 20, 21, 22, 23, 24:
		return 8
	default:
		return 9
	}
}

func moveResponseToStruct(data VerboseMoveResponse) (PokeMove, error) {
	var move PokeMove
	move.MoveID = data.ID
	move.Accuracy = data.Accuracy
	move.Power = data.Power
	move.PowerPoints = data.PowerPoints
	move.Name = data.Name
	move.Type = data.Type.Name
	move.DamageType = data.DamageType.Name
	move.Generation = getGeneration(data.Generation.Name)
	move.Description = ""

	if len(data.EffectDescription) > 0 {
		move.Description = data.EffectDescription[0].ShortEffect
	}

	return move, nil
}

func getGeneration(generation string) int {
	switch generation {
	case "generation-i":
		return 1
	case "generation-ii":
		return 2
	case "generation-iii":
		return 3
	case "generation-iv":
		return 4
	case "generation-v":
		return 5
	case "generation-vi":
		return 6
	case "generation-vii":
		return 7
	case "generation-viii":
		return 8
	default:
		return -1
	}
}

func createCsv(path string, entries []CsvEntry) (*os.File, error) {
	csvFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	header := entries[0].GetHeader()
	w := csv.NewWriter(csvFile)	
	w.Comma = '|'
	
	if err = w.Write(header); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if err = w.Write(entry.ToSlice()); err != nil {
			return nil, err
		}
	}

	w.Flush()
	return csvFile, nil
}