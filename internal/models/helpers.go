package models

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

func transactionSetup(db *sql.DB, tblInfo []string) (*sql.Stmt, func() error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf(err.Error())
	}

	stmt, err := tx.Prepare(pq.CopyIn(tblInfo[0], tblInfo[1:]...))
	if err != nil {
		log.Fatalf(err.Error())
	}

	teardown := func() error {
		_, e := stmt.Exec()
		if e != nil {
			return e
		}
	
		e = stmt.Close()
		if e != nil {
			return e
		}
	
		e = tx.Commit()
		if err != nil {
			return e
		}

		return nil
	}
	
	return stmt, teardown
}