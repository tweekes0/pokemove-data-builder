package models

import (
	"database/sql"
	"errors"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
)

const (
	pokemonInsert = `INSERT INTO pokemon(poke_id, name, sprite, species, origin_gen) 
	VALUES ($1, $2, $3, $4 $5)`
	pokemonDelete    = `DELETE FROM pokemon WHERE poke_id = $1`
	pokemonGetByID   = `SELECT poke_id, name, sprite, species, origin_gen FROM pokemon WHERE poke_id = $1`
	pokemonGetByName = `SELECT poke_id, name, sprite, species, origin_gen FROM pokemon WHERE name = $1`
	pokemonGetAll    = `SELECT poke_id, name, sprite, species, origin_gen FROM pokemon`
	pokemonExists    = `SELECT EXISTS(SELECT 1 FROM pokemon WHERE id = $1)`
	pokemonMovesJoin = `
	SELECT DISTINCT 
	pm.move_id, pm.name, pm.accuracy, pm.power, pm.power_points,
	pm.type, pm.damage_type, pm.description,  
	pmr.learn_method, pmr.level_learned, pmr.game_name, pmr.generation
	FROM pokemon p 
	JOIN pokemon_move_rels pmr ON p.poke_id = pmr.poke_id
	JOIN pokemon_moves pm ON pm.move_id = pmr.move_id
	WHERE p.poke_id = $1 and pmr.generation = $2;
	`
)

type PokemonModel struct {
	DB *sql.DB
}

type MovesJoinRow struct {
	MoveID       int    `json:"id"`
	Accuracy     int    `json:"accuracy"`
	Power        int    `json:"power"`
	PowerPoints  int    `json:"power_points"`
	Generation   int    `json:"generation"`
	LevelLearned int    `json:"level_learned"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	DamageType   string `json:"damage_type"`
	Description  string `json:"description"`
	LearnMethod  string `json:"learn_method"`
	GameName     string `json:"game_name"`
}

func (m *PokemonModel) BulkInsert(pokemon []interface{}) error {
	tblInfo := []string{
		"pokemon", "poke_id", "name", "sprite", "species", "origin_gen",
	}
	stmt, teardown := transactionSetup(m.DB, tblInfo)

	for _, p := range pokemon {
		_, err := stmt.Exec(
			p.(client.Pokemon).PokeID,
			p.(client.Pokemon).Name,
			p.(client.Pokemon).Sprite,
			p.(client.Pokemon).Species,
			p.(client.Pokemon).OriginGen,
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

func (m *PokemonModel) RelationsBulkInsert(rels []interface{}) error {
	tblInfo := []string{
		"pokemon_move_rels", "poke_id", "move_id", "generation",
		"level_learned", "learn_method", "game_name",
	}
	stmt, teardown := transactionSetup(m.DB, tblInfo)

	for _, rel := range rels {
		_, err := stmt.Exec(
			rel.(client.PokemonMoveRelation).PokeID, 
			rel.(client.PokemonMoveRelation).MoveID, 
			rel.(client.PokemonMoveRelation).Generation,
			rel.(client.PokemonMoveRelation).LevelLearned,
			rel.(client.PokemonMoveRelation).LearnMethod, 
			rel.(client.PokemonMoveRelation).GameName,
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

func (m *PokemonModel) PokemonInsert(p client.Pokemon) error {
	_, err := m.DB.Exec(pokemonInsert, p.PokeID, p.Name, p.Sprite, p.Species, p.OriginGen)
	if err != nil {
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

		err = rows.Scan(&p.PokeID, &p.Name, &p.Sprite, &p.Species, &p.OriginGen)
		if err != nil {
			return nil, err
		}

		pokemon = append(pokemon, p)
	}

	if len(pokemon) == 0 {
		return nil, ErrDoesNotExist
	}

	return pokemon, nil
}

func (m *PokemonModel) PokemonMovesJoinByGen(pokeID, gen int) ([]*MovesJoinRow, error) {
	mvs := []*MovesJoinRow{}

	rows, err := m.DB.Query(pokemonMovesJoin, pokeID, gen)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		mv := &MovesJoinRow{}

		err := rows.Scan(
			&mv.MoveID, &mv.Name, &mv.Accuracy, &mv.Power, &mv.PowerPoints,
			&mv.Type, &mv.DamageType, &mv.Description,
			&mv.LearnMethod, &mv.LevelLearned, &mv.GameName, &mv.Generation,
		)

		if err != nil {
			return nil, err
		}

		mvs = append(mvs, mv)
	}

	return mvs, nil
}
