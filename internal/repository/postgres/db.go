package postgres

import (
	"database/sql"
	"fmt"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "Katy314"
	DB_NAME     = "FilmLibrary"
)

// DB set up
func SetupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	return db
}
