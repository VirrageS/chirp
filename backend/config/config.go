package config

// TODO: config file path should probably be read from env variable or command line argument
// TODO: rename this file

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Interfaces that provide only those parameters that are required by different prats of the system

// Provides configuration for ServiceProvider
type ServiceConfigProvider interface {
	GetSecretKey() []byte
	GetTokenValidityPeriod() int
}

// Provides secret key
type SecretKeyProvider interface {
	GetSecretKey() []byte
}

// Stores global configuration of the system
type Configuration struct {
	secretKey           []byte
	tokenValidityPeriod int
}

func (config *Configuration) GetSecretKey() []byte {
	return config.secretKey
}

func (config *Configuration) GetTokenValidityPeriod() int {
	return config.tokenValidityPeriod
}

// TODO: Maybe read the config only once on init() or something and then return the global object?
func GetConfig() *Configuration {
	return readServiceConfig()
}

func readServiceConfig() *Configuration {
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

	return &Configuration{
		secretKey:           []byte(configSecretKey),
		tokenValidityPeriod: configValidityPeriod,
	}
}
