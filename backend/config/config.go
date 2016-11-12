package config

// TODO: config file path should probably be read from env variable or command line argument

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var secretKey []byte
var tokenValidityPeriod int

func GetSecretKey() []byte {
	return secretKey
}

func GetTokenValidityPeriod() int {
	return tokenValidityPeriod
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$GOPATH/src/github.com/VirrageS/chirp/backend")

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("Error reading config file.")
	}

	configSecretKey := viper.GetString("secret_key")
	configValidityPeriod := viper.GetInt("token_validity_period")

	if configSecretKey == "" || configValidityPeriod <= 0 {
		log.WithFields(log.Fields{
			"secret key":      configSecretKey,
			"validity period": configValidityPeriod,
		}).Fatal("Config file doesn't contain valid data.")
	}

	secretKey = []byte(configSecretKey)
	tokenValidityPeriod = configValidityPeriod
}
