package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"

	"github.com/VirrageS/chirp/backend/cache"
)

// Struct that implements DatabaseAccessor
type Database struct {
	UserDataAccessor
	TweetDataAccessor
}

// Constructs new Database that uses given sql.DB connection
func NewDatabase(databaseConnection *sql.DB, cache cache.CacheProvider) DatabaseAccessor {
	return &Database{
		NewUserDB(databaseConnection, cache),
		NewTweetDB(databaseConnection, cache),
	}
}

// Returns new connection to DB specified in config file. Panics when unrecoverable error occurs.
// For now it takes port as parameter so we can redirect tests to testing database
func NewConnection(port string) *sql.DB {
	// TODO: read user data, host and port from config file

	db, err := sql.Open("postgres", "user=postgres password=postgres host=localhost sslmode=disable port="+port)
	if err != nil {
		log.WithError(err).Fatal("Error opening database.")
	}

	err = db.Ping()
	if err != nil {
		log.WithError(err).Fatal("Error connecting to database.")
	}

	return db
}
