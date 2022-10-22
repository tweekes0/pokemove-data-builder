package models

import (
	"database/sql"
	"errors"
	// "errors"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
)

const (
	moveInsert = `INSERT INTO pokemon_moves(move_id, name, accuracy, power,
		powerpoints, generation, type, damagetype, description) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	moveDelete  = `DELETE FROM pokemon_moves WHERE move_id = $1`
	moveGetByID = `SELECT move_id, name, accuracy, power, powerpoints,
		generation, type, damagetype, description 
		FROM pokemon_moves WHERE move_id = $1 AND generation = $2`
	moveGetByName = `SELECT move_id, name, accuracy, power, powerpoints,
		generation, type, damagetype, description FROM pokemon_moves WHERE name = $1`
	moveGetAll = `SELECT move_id, name, accuracy, power, powerpoints,
		generation, type, damagetype, description FROM pokemon_moves`
	moveExists = `SELECT EXISTS(SELECT 1 FROM pokemon_moves WHERE id = $1)`
)

type MovesModel struct {
	DB *sql.DB
}

func (m *MovesModel) MoveInsert(mv client.PokemonMove) error {
	_, err := m.DB.Exec(
		moveInsert,
		mv.MoveID, mv.Name, mv.Accuracy, mv.Power, mv.PowerPoints,
		mv.Generation, mv.Type, mv.DamageType, mv.Description,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *MovesModel) MoveGet(moveID, gen int) (*client.PokemonMove, error) {
	mv := &client.PokemonMove{}

	err := m.DB.QueryRow(moveGetByID, moveID, gen).Scan(
		&mv.MoveID, &mv.Name, &mv.Accuracy, &mv.Power, &mv.PowerPoints,
		&mv.Generation, &mv.Type, &mv.DamageType, &mv.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDoesNotExist
		}

		return nil, err
	}

	return mv, nil
}
