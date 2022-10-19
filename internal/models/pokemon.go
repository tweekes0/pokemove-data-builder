package models

import (
	"database/sql"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
)

const (
	insert = `INSERT INTO pokemon(poke_id, name, sprite, species) 
	VALUES ($1, $2, $3, $4)`
	deleteByID  = `DELETE FROM pokemon WHERE poke_id = $1`
	queryByID   = `SELECT * FROM pokemon WHERE poke_id = $1`
	queryByName = `SELECT * FROM pokemon WHERE name = $1`
	queryAll    = `SELECT * FROM pokemon`
	exists = `SELECT EXISTS(SELECT 1 FROM pokemon WHERE id = $1)`
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

func (m *PokemonModel) PokemonDelete(id int) error {
	res, err := m.DB.Exec(deleteByID, id)
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

func (m *PokemonModel) PokemonExists(id int) (bool, error) {
	var e bool

	err := m.DB.QueryRow(exists, id).Scan(&e)
	if err != nil {
		return false, err
	}

	return e, nil
}
