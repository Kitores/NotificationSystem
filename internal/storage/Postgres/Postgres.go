package Postgres

import "github.com/jmoiron/sqlx"

type PostgreSqlx struct {
	sqlx *sqlx.DB
}

func NewPg(connString string) (*PostgreSqlx, error) {

}
