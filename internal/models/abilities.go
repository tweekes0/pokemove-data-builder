package models

import (
	"database/sql"
	"errors"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
)

const (
	abilityInsert = `INSERT INTO pokemon_abilities(ability_id, name,
		description, generation) VALUES ($1, $2, $3, $4)`
	abilityGetById = `SELECT ability_id, name, description, generation 
	FROM pokemon_abilities WHERE ability_id = $1`
	abilityGetAll = `SELECT ability_id, name, description, generation FROM pokemon_abilities`
)

type AbilitiesModel struct {
	DB *sql.DB
}

func (m *AbilitiesModel) AbilityInsert(a client.PokemonAbility) error {
	_, err := m.DB.Exec(
		abilityInsert,
		a.AbilityID, a.Name, a.Description, a.Generation,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *AbilitiesModel) AbilityBulkInsert(ab []client.PokemonAbility) error {
	tblInfo := []string{
		"pokemon_abilities", "ability_id", "name", "description",
		"generation",
	}
	stmt, teardown := transactionSetup(m.DB, tblInfo)

	for _, a := range ab {
		_, err := stmt.Exec(
			a.AbilityID, a.Name, a.Description, a.Generation,
		)

		if err != nil {
			return err
		}
	}

	if err := teardown(); err != nil {
		return err
	}

	return nil
}

func (m *AbilitiesModel) AbilityGet(a_id int) (*client.PokemonAbility, error) {
	a := &client.PokemonAbility{}

	err := m.DB.QueryRow(abilityGetById, a_id).Scan(
		&a.AbilityID, &a.Name, &a.Description, &a.Generation,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDoesNotExist
		}

		return nil, err
	}

	return a, nil
}

func (m *AbilitiesModel) AbilityGetAll() ([]*client.PokemonAbility, error) {
	var abs []*client.PokemonAbility

	rows, err := m.DB.Query(abilityGetAll)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		a := client.PokemonAbility{}

		err := rows.Scan(&a.AbilityID, &a.Name, &a.Description, &a.Generation,)

		if err != nil {
			return nil, err	
		}

		abs = append(abs, &a)
	}

	return abs, nil
}
