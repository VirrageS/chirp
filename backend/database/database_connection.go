package database

// TODO: find the best place/package to store this file (think about design)

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var Database *sql.DB = initDB()

func initDB() *sql.DB {
	// TODO: read user data, host and port from config file or something like that
	db, err := sql.Open("postgres", "user=postgres password=postgres host=localhost sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("Couldn't open database! Error: %v.", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect to database! Error: %v.", err))
	}

	return db
}
