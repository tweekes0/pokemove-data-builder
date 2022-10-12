package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// interface for writing structs to CSV files
type CsvEntry interface {
	GetHeader() []string
	ToSlice() []string
}

// resolves the pokeapi version group to a generation number
// https://pokeapi.co/docs/v2#versiongroup
func resolveVersionGroup(url string) int {
	id := getVersionGroupID(url)

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
		return -1
	}
}

func getFlavorText(gen int, lang string, texts []flavorText) string {
	defaultText := getDefaultFlavorText(lang, texts)

	for _, text := range texts {
		id := resolveVersionGroup(text.VersionGroup.Url)
		if gen == id && lang == text.Language.Name {
			ret := strings.ReplaceAll(text.Text, "\n", " ")
			return ret
		}
	}

	return defaultText
}

func getDefaultFlavorText(lang string, texts []flavorText) string {
	for _, text := range texts {
		if lang == text.Language.Name {
			ret := strings.ReplaceAll(text.Text, "\n", " ")
			return ret
		}
	}

	return ""
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

func getVersionGroupID(url string) int {
	group_id := strings.Split(url, "/")[6]
	id, err := strconv.Atoi(group_id)
	if err != nil {
		return -1
	}

	return id
}

func createDataDir() error {
	_, err := os.Stat("./data")
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir("./data", 0755); err != nil { 
				return err
			}
		}
	}

	return nil
}

func createCsv(path string) (*os.File, error) {
	if err := createDataDir(); err != nil {
		return nil, err
	}

	csvFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	return csvFile, nil
}

func writeCsvEntry(w *csv.Writer, entry CsvEntry) error {
	if err := w.Write(entry.ToSlice()); err != nil {
		return err
	}

	return nil
}

// func to write APIReceivers to a csv file
func ToCsv(csvFile *os.File, recv APIReceiver) error {
	if len(recv.CsvEntries()) == 0 {
		return ErrEmptyCsv
	}

	file, err := os.OpenFile(csvFile.Name(), os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = '|'
	defer w.Flush()

	if err = w.Write(recv.CsvEntries()[0].GetHeader()); err != nil {
		return err
	}

	for _, entry := range recv.CsvEntries() {
		if err = writeCsvEntry(w, entry); err != nil {
			return err
		}
	}

	return nil
}