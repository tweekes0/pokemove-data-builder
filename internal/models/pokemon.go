package models

import (
	"database/sql"
	"errors"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
)

const (
	insert = `INSERT INTO pokemon(poke_id, name, sprite, species) 
	VALUES ($1, $2, $3, $4)`
	deleteByID  = `DELETE FROM pokemon WHERE poke_id = $1`
	queryByID   = `SELECT poke_id, name, sprite, species FROM pokemon WHERE poke_id = $1`
	queryByName = `SELECT poke_id, name, sprite, species FROM pokemon WHERE name = $1`
	queryAll    = `SELECT poke_id, name, sprite, species FROM pokemon`
	exists      = `SELECT EXISTS(SELECT 1 FROM pokemon WHERE id = $1)`
)

type PokemonModel struct {
	DB *sql.DB
}

func (m *PokemonModel) PokemonInsert(p client.Pokemon) error {
	_, err := m.DB.Exec(insert, p.PokeID, p.Name, p.Sprite, p.Species)
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
	res, err := m.DB.Exec(deleteByID, pokeID)
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

	err := m.DB.QueryRow(exists, pokeID).Scan(&e)
	if err != nil {
		return false, err
	}

	return e, nil
}

func (m *PokemonModel) PokemonGet(pokeID int) (*client.Pokemon, error) {
	p := &client.Pokemon{}

	err := m.DB.QueryRow(queryByID, pokeID).Scan(&p.PokeID, &p.Name, &p.Sprite, &p.Species)
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

	rows, err := m.DB.Query(queryAll)
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

	return pokemon, nil
}