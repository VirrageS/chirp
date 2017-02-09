package database

import (
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/config"
)

// Returns new connection to DB specified in config file. Panics when unrecoverable error occurs.
func NewConnection(config config.DatabaseConfigProvider) *sql.DB {
	username := config.GetUsername()
	password := config.GetPassword()
	host := config.GetHost()
	port := config.GetPort()

	accessString := fmt.Sprintf("user=%s password=%s host=%s sslmode=disable port=%s", username, password, host, port)
	db, err := sql.Open("postgres", accessString)
	if err != nil {
		log.WithError(err).Fatal("Error opening database.")
	}

	err = db.Ping()
	if err != nil {
		log.WithError(err).Fatal("Error connecting to database.")
	}

	return db
}
