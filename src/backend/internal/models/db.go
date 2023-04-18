package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/tweekes0/pokemoves/src/backend/internal/client"
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

func getDBN() string {
	return fmt.Sprintf(
		"dbname=%v host=%v user=%v password=%v sslmode=%v",
		os.Getenv("POSTGRES_DB"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"disable",
	)
}

func NewDBConn() (*DBConn, error) {
	connString := getDBN()
	db, err := sql.Open("postgres", connString)
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
