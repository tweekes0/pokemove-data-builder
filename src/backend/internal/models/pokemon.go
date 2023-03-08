package models

import (
	"database/sql"
	"errors"

	"github.com/tweekes0/pokemoves/src/backend/internal/client"
)

const (
	pokemonExists  = `SELECT EXISTS(SELECT 1 FROM pokemon WHERE id = $1)`
	pokemonGetByID = `
	SELECT 
		poke_id, name, sprite, shiny_sprite, species, origin_gen, primary_type, secondary_type 
	FROM pokemon WHERE poke_id = $1 and gen_of_type_change >= $2
	`
	pokemonGetByName = `
	SELECT 
		poke_id, name, sprite, shiny_sprite, species, origin_gen, primary_type, secondary_type 
	FROM pokemon WHERE name = $1 AND gen_of_type_change >= $2 
	`
	pokemonGetAll = `
	SELECT DISTINCT ON (poke_id)
		poke_id, name, sprite, shiny_sprite, species, origin_gen,
		primary_type, secondary_type 
	FROM pokemon ORDER BY poke_id;
	`
	pokemonGetAllBrief = `SELECT poke_id, name FROM pokemon ORDER BY poke_id;`
	pokemonMovesJoin = `
	SELECT
		pm.move_id, pm.name, pm.accuracy, pm.power, pm.power_points,
		pm.type, pm.damage_type, pm.description,  
		pmr.learn_method, pmr.level_learned, pmr.game_name, pm.generation
	FROM pokemon_move_rels pmr
	JOIN pokemon p ON p.poke_id = pmr.poke_id
	JOIN (
		SELECT DISTINCT ON (move_id) pm.* FROM pokemon_moves pm 
		WHERE generation <= $1
		ORDER BY move_id, generation desc
	) pm ON pm.move_id = pmr.move_id 
	WHERE p.poke_id = $2 and pmr.generation = $3 
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
		"pokemon", "poke_id", "name", "sprite", "shiny_sprite", "species",
		"origin_gen", "gen_of_type_change", "primary_type", "secondary_type",
	}
	stmt, teardown := transactionSetup(m.DB, tblInfo)

	for _, p := range pokemon {
		_, err := stmt.Exec(
			p.(client.Pokemon).PokeID,
			p.(client.Pokemon).Name,
			p.(client.Pokemon).Sprite,
			p.(client.Pokemon).ShinySprite,
			p.(client.Pokemon).Species,
			p.(client.Pokemon).OriginGen,
			p.(client.Pokemon).GenTypeChange,
			p.(client.Pokemon).PrimaryType,
			p.(client.Pokemon).SecondaryType,
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

func (m *PokemonModel) PokemonExists(pokeID int) (bool, error) {
	var e bool

	err := m.DB.QueryRow(pokemonExists, pokeID).Scan(&e)
	if err != nil {
		return false, err
	}

	return e, nil
}

func (m *PokemonModel) PokemonGet(pokeID, gen int) (*client.Pokemon, error) {
	p := &client.Pokemon{}

	if gen < 0 || gen > client.CurrentGen {
		gen = client.CurrentGen
	} 

	err := m.DB.QueryRow(pokemonGetByID, pokeID, gen).Scan(
		&p.PokeID, &p.Name, &p.Sprite, &p.ShinySprite, &p.Species, &p.OriginGen,
		&p.PrimaryType, &p.SecondaryType,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDoesNotExist
		}

		return nil, err
	}

	return p, nil
}

func (m *PokemonModel) PokemonGetByName(name string, gen int) (*client.Pokemon, error) {
	p := &client.Pokemon{}

	err := m.DB.QueryRow(pokemonGetByName, name).Scan(
		&p.PokeID, &p.Name, &p.Sprite, &p.ShinySprite, &p.Species, &p.OriginGen,
		&p.PrimaryType, &p.SecondaryType,
	)
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

		err = rows.Scan(
			&p.PokeID, &p.Name, &p.Sprite, &p.ShinySprite, &p.Species, &p.OriginGen,
			&p.PrimaryType, &p.SecondaryType,
		)
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

func (m *PokemonModel) PokemonGetAllBrief() ([]*client.PokemonBrief, error) {
	pks :=[]*client.PokemonBrief{}

	rows, err := m.DB.Query(pokemonGetAllBrief)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		pk := &client.PokemonBrief{}

		if err := rows.Scan(&pk.PokeID, &pk.Name); err != nil {
			return nil, err
		}

		pks = append(pks, pk)
	}

	if len(pks) == 0 {
		return nil, ErrDoesNotExist
	}

	return pks, nil
}

func (m *PokemonModel) PokemonMovesJoinByGen(pokeID, gen int) ([]*MovesJoinRow, error) {
	mvs := []*MovesJoinRow{}

	rows, err := m.DB.Query(pokemonMovesJoin, gen, pokeID, gen)
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
