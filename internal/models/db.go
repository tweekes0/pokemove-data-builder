package models

import (
	"database/sql"

	"github.com/tweekes0/pokemonmoves-backend/internal/config"
)

type DBConn struct {
	*sql.DB
	AbilitiesModel
	MovesModel
	PokemonModel
}

func NewDBConn() (*DBConn, error) {
	conf, err := config.ReadDBConfig()
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
