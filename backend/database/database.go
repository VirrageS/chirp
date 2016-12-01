package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

// Struct that implements DatabaseAccessor
type Database struct {
	UserDataAccessor
	TweetDataAccessor
}

// Constructs new Database that uses given sql.DB connection
func NewDatabase(databaseConnection *sql.DB) DatabaseAccessor {
	return &Database{
		NewUserDB(databaseConnection),
		NewTweetDB(databaseConnection),
	}
}

// Returns new connection to DB specified in config file. Panics when unrecoverable error occurs.
func NewConnection() *sql.DB {
	// TODO: read user data, host and port from config file

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
