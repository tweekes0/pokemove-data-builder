package models

import (
	"database/sql"
	"errors"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
)

const (
	pokemonInsert = `INSERT INTO pokemon(poke_id, name, sprite, species) 
	VALUES ($1, $2, $3, $4)`
	pokemonDelete    = `DELETE FROM pokemon WHERE poke_id = $1`
	pokemonGetByID   = `SELECT poke_id, name, sprite, species FROM pokemon WHERE poke_id = $1`
	pokemonGetByName = `SELECT poke_id, name, sprite, species FROM pokemon WHERE name = $1`
	pokemonGetAll    = `SELECT poke_id, name, sprite, species FROM pokemon`
	pokemonExists    = `SELECT EXISTS(SELECT 1 FROM pokemon WHERE id = $1)`
)

type PokemonModel struct {
	DB *sql.DB
}

func (m *PokemonModel) PokemonInsert(p client.Pokemon) error {
	_, err := m.DB.Exec(pokemonInsert, p.PokeID, p.Name, p.Sprite, p.Species)
	if err != nil {
		return err
	}

	return nil
}

func (m *PokemonModel) PokemonBulkInsert(pokemon []client.Pokemon) error {
	tblInfo := []string{"pokemon", "poke_id", "name", "sprite", "species"}
	stmt, teardown := transactionSetup(m.DB, tblInfo)

	for _, p := range pokemon {
		_, err := stmt.Exec(p.PokeID, p.Name, p.Sprite, p.Species)
		if err != nil {
			return err
		}
	}

	if err := teardown(); err != nil {
		return err
	}

	return nil
}

func (m *PokemonModel) PokemonDelete(pokeID int) error {
	res, err := m.DB.Exec(pokemonDelete, pokeID)
	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if int(c) == 0 {
		return ErrDoesNotExist
	}

	return nil
}

func (m *PokemonModel) PokemonExists(pokeID int) (bool, error) {
	var e bool

	err := m.DB.QueryRow(pokemonExists, pokeID).Scan(&e)
	if err != nil {
		return false, err
	}

	return e, nil
}

func (m *PokemonModel) PokemonGet(pokeID int) (*client.Pokemon, error) {
	p := &client.Pokemon{}

	err := m.DB.QueryRow(pokemonGetByID, pokeID).Scan(&p.PokeID, &p.Name, &p.Sprite, &p.Species)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDoesNotExist
		}

		return nil, err
	}

	return p, nil
}

func (m *PokemonModel) PokemonGetAll() ([]*client.Pokemon, error) {
	pokemon := []*client.Pokemon{}

	rows, err := m.DB.Query(pokemonGetAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := &client.Pokemon{}

		if err = rows.Scan(&p.PokeID, &p.Name, &p.Sprite, &p.Species); err != nil {
			return nil, err
		}

		pokemon = append(pokemon, p)
	}

	if len(pokemon) == 0 {
		return nil, ErrDoesNotExist
	}

	return pokemon, nil
}
