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

type moveData struct {
	Move client.PokemonMove
	Rel  client.PokemonMoveRelation
}

type MovesJoin struct {
	P client.Pokemon

	Moves []moveData
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

func (m *PokemonModel) MoveRelationsBulkInsert(rels []client.PokemonMoveRelation) error {
	tblInfo := []string{
		"pokemon_move_rels", "poke_id", "move_id", "generation",
		"level_learned", "learn_method", "game_name",
	}
	stmt, teardown := transactionSetup(m.DB, tblInfo)

	for _, rel := range rels {
		_, err := stmt.Exec(
			rel.PokeID, rel.MoveID, rel.Generation, rel.LevelLearned,
			rel.LearnMethod, rel.GameName,
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

func (m *PokemonModel) PokemonMovesJoinByGen(pokeID, gen int) (*MovesJoin, error) {
	mj := &MovesJoin{}

	p, err := m.PokemonGet(pokeID)
	if err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(pokemonMovesJoin, pokeID, gen)
	if err != nil {
		return nil, err
	}

	mvs := []moveData{}
	for rows.Next() {
		var mv client.PokemonMove
		var rel client.PokemonMoveRelation

		err := rows.Scan(
			&mv.MoveID, &mv.Name, &mv.Accuracy, &mv.Power, &mv.PowerPoints,
			&mv.Type, &mv.DamageType, &mv.Description,
			&rel.LearnMethod, &rel.LevelLearned, &rel.GameName, &rel.Generation,
		)

		if err != nil {
			return nil, err
		}

		md := moveData{
			Move: mv,
			Rel: rel,
		}

		mvs = append(mvs, md)
	}

	mj.P = *p
	mj.Moves = mvs

	return mj, nil
}
