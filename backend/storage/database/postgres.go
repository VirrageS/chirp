package database

import (
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/config"
)

// NewPostgresDatabase returns new connection to PostgreSQL database.
// All configurations can be specified in config file.
func NewPostgresDatabase(config config.DatabaseConfigProvider) *Connection {
	username := config.GetUsername()
	password := config.GetPassword()
	host := config.GetHost()
	port := config.GetPort()

	accessString := fmt.Sprintf("user=%s password=%s host=%s sslmode=disable port=%s", username, password, host, port)
	db, err := sql.Open("postgres", accessString)
	if err != nil {
		log.WithError(err).Error("Error opening database.")
		return nil
	}

	if err = db.Ping(); err != nil {
		log.WithError(err).Error("Error connecting to database.")
		return nil
	}

	return &Connection{db}
}
