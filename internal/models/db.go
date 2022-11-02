package models

import (
	"database/sql"

	"github.com/tweekes0/pokemonmoves-backend/internal/client"
	"github.com/tweekes0/pokemonmoves-backend/internal/config"
)

const (
	countQuery = "SELECT count(id) FROM pokemon"
)

type DBConn struct {
	*sql.DB
	AbilitiesModel
	MovesModel
	PokemonModel
}

func (c *DBConn) getModels() []Model {
	return []Model{
		&c.AbilitiesModel,
		&c.MovesModel,
		&c.PokemonModel,
	}
}

func NewDBConn() (*DBConn, error) {
	conf, err := config.LoadDBConfig()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", conf.GetDBN())
	if err != nil {
		return nil, err
	}

	return &DBConn{
		db,
		AbilitiesModel{db},
		MovesModel{db},
		PokemonModel{db},
	}, nil
}

func (c *DBConn) CheckDB() (bool, error) {	
	var count int

	if err := c.QueryRow(countQuery).Scan(&count); err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (c *DBConn) PopulateDB(recv ...client.APIReceiver) error {
	for i, m := range c.getModels() {
		if err := m.BulkInsert(recv[i].GetEntries()); err != nil {
			return err
		}

		if err := m.RelationsBulkInsert(recv[i].GetRelations()); err != nil {
			return err
		}
	}

	return nil
}