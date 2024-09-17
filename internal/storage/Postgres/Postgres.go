package postgres

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgreSqlx struct {
	db *sqlx.DB
}

func NewPg(connString string) (*PostgreSqlx, error) {
	Db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		fmt.Errorf("Error connecting to PostgreSQL")
		return nil, err
	}
	postgresInit := &PostgreSqlx{db: Db}
	if postgresInit.db == nil {
		fmt.Errorf("Error init to PostgreSQL")
		return nil, err
	}
	return postgresInit, nil
}

func (pg *PostgreSqlx) SaveNewUser(firstName, lastName string, phoneNumber int) error {
	const funcName = "storage/postgres/SaveNewUser()"
	query := fmt.Sprintf("INSERT INTO users(phone_number, first_name, last_name) VALUES('%s', '%s', '%s')", phoneNumber, firstName, lastName)
	_, err := pg.db.Exec(query)
	if err != nil {
		return fmt.Errorf("Error inserting new user into PostgreSQL")
	}
	return err
}
