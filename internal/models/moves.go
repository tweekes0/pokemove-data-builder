package models

import (
	"database/sql"
	"errors"
	// "errors"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
)

const (
	moveInsert = `INSERT INTO pokemon_moves(move_id, name, accuracy, power,
		power_points, generation, type, damage_type, description) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	moveDelete   = `DELETE FROM pokemon_moves WHERE move_id = $1`
	moveGetByGen = `SELECT move_id, name, accuracy, power, power_points,
		generation, type, damage_type, description 
		FROM pokemon_moves WHERE move_id = $1 AND generation = $2`
	moveGetByID = `SELECT move_id, name, accuracy, power, power_points,
		generation, type, damage_type, description 
		FROM pokemon_moves WHERE move_id = $1`
	moveGetByName = `SELECT move_id, name, accuracy, power, power_points,
		generation, type, damage_type, description FROM pokemon_moves WHERE name = $1`
	moveGetAll = `SELECT move_id, name, accuracy, power, power_points,
		generation, type, damage_type, description FROM pokemon_moves`
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

func (m *PokemonModel) MoveBulkInsert(moves []client.PokemonMove) error {
	tblInfo := []string{
		"pokemon_moves", "move_id", "name", "accuracy",
		"power", "power_points", "generation", "type", "damage_type",
		"description",
	}
	stmt, teardown := transactionSetup(m.DB, tblInfo)

	for _, mv := range moves {
		_, err := stmt.Exec(
			mv.MoveID, mv.Name, mv.Accuracy, mv.Power, mv.PowerPoints,
			mv.Generation, mv.Type, mv.DamageType, mv.Description,
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

func (m *MovesModel) MoveGet(moveID int) ([]*client.PokemonMove, error) {
	moves := []*client.PokemonMove{}

	rows, err := m.DB.Query(moveGetByID, moveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		mv := &client.PokemonMove{}

		err = rows.Scan(
			&mv.MoveID, &mv.Name, &mv.Accuracy, &mv.Power, &mv.PowerPoints,
			&mv.Generation, &mv.Type, &mv.DamageType, &mv.Description,
		)
		if err != nil {
			return nil, err
		}

		moves = append(moves, mv)
	}

	if len(moves) == 0 {
		return nil, ErrDoesNotExist
	}

	return moves, nil
}

func (m *MovesModel) MoveGetByGeneration(moveID, gen int) (*client.PokemonMove, error) {
	mv := &client.PokemonMove{}

	err := m.DB.QueryRow(moveGetByGen, moveID, gen).Scan(
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
