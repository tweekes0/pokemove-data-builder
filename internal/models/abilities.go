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
	abilityGetAll        = `SELECT ability_id, name, description, generation FROM pokemon_abilities`
	pokemonAbilitiesJoin = `
	SELECT DISTINCT
	pa.ability_id, pa.name, pa.description, pa.generation,
	par.slot, par.hidden
	FROM pokemon p
	JOIN pokemon_ability_rels par on p.poke_id = par.poke_id
	JOIN pokemon_abilities pa on pa.ability_id = par.ability_id
	WHERE p.poke_id = $1
	`
)

type AbilitiesModel struct {
	DB *sql.DB
}

type AbilityJoinRow struct {
	ID          int
	Name        string
	Description string
	Generation  int
	Slot        int
	Hidden      bool
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

		err := rows.Scan(&a.AbilityID, &a.Name, &a.Description, &a.Generation)

		if err != nil {
			return nil, err
		}

		abs = append(abs, &a)
	}

	return abs, nil
}

func (m *AbilitiesModel) AbilityRelationsBulkInsert(
	rels []client.PokemonAbilityRelation) error {

	tblInfo := []string{
		"pokemon_ability_rels", "poke_id", "ability_id", "slot", "hidden",
	}
	stmt, teardown := transactionSetup(m.DB, tblInfo)

	for _, rel := range rels {
		_, err := stmt.Exec(rel.PokeID, rel.AbilityID, rel.Slot, rel.Hidden)
		if err != nil {
			return err
		}
	}

	if err := teardown(); err != nil {
		return err
	}

	return nil
}

func (m *AbilitiesModel) PokemonAbilitiesJoin(pokeID int) ([]*AbilityJoinRow, error) {
	abs := []*AbilityJoinRow{}

	rows, err := m.DB.Query(pokemonAbilitiesJoin, pokeID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		ab := &AbilityJoinRow{}
		err := rows.Scan(
			&ab.ID, &ab.Name, &ab.Description,
			&ab.Generation, &ab.Slot, &ab.Hidden,
		)

		if err != nil {
			return nil, err
		}

		abs = append(abs, ab)
	}

	return abs, nil
}
