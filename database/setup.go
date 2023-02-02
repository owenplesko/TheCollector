package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Setup(url string, port string, user string, password string, dbName string) error {
	var err error
	connection := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", url, port, user, password, dbName)
	db, err = sql.Open("postgres", connection)
	return err
}
