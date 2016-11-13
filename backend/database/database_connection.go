package database

// TODO: find the best place/package to store this file (think about design)

import (
	"database/sql"

	_ "github.com/lib/pq"

	log "github.com/Sirupsen/logrus"
)

func NewDatabaseConnection() *sql.DB {
	// TODO: read user data, host and port from config file or something like that

	db, err := sql.Open("postgres", "user=postgres password=postgres host=localhost sslmode=disable")
	if err != nil {
		log.WithError(err).Fatal("Error opening database.")
	}

	err = db.Ping()
	if err != nil {
		log.WithError(err).Fatal("Error connecting to database.")
	}

	return db
}
