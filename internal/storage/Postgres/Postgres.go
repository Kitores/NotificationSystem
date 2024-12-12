package postgres

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
)

type PostgreSqlx struct {
	db *sqlx.DB
}
type Subscriber struct {
	Id         int
	Phone      string
	FirstName  string
	LastName   string
	TelegramId int
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

func (pg *PostgreSqlx) DeleteUser(firstName string, lastName string) error {
	const funcName = "storage/postgres/DeleteUser()"
	query := fmt.Sprintf("DELETE FROM users WHERE first_name = '%s' AND last_name = '%s'", firstName, lastName)
	_, err := pg.db.Exec(query)
	if err != nil {
		return fmt.Errorf("Error deleting user from PostgreSQL")
	}
	return err
}

func (pg *PostgreSqlx) GetUsers() ([]Subscriber, error) {
	const funcName = "storage/postgres/GetUsers()"
	query := "SELECT * FROM subs"
	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error getting users from PostgreSQL")
	}
	defer rows.Close()
	var subs []Subscriber
	for rows.Next() {
		var sub Subscriber
		if err = rows.Scan(&sub.Id, &sub.Phone, &sub.FirstName, &sub.LastName, &sub.TelegramId); err != nil {
			log.Println(err)
		}
		subs = append(subs, sub)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error getting users from PostgreSQL")
	}
	return subs, err
}
