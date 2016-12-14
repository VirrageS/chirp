package config

// TODO: config file path should probably be read from env variable or command line argument
// TODO: rename this file

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Interfaces that provide only those parameters that are required by different prats of the system

// Provides secret key
type SecretKeyProvider interface {
	GetSecretKey() []byte
}

// Provides configuration for ServiceProvider
type ServiceConfigProvider interface {
	SecretKeyProvider
	GetAuthTokenValidityPeriod() int
	GetRefreshTokenValidityPeriod() int
}

// Provides configuration for CacheProvider
type CacheConfigProvider interface {
	GetCacheExpirationTime() time.Duration
}

// Provides full configuration
type ConfigProvider interface {
	ServiceConfigProvider
	CacheConfigProvider
}

// Stores global configuration of the system
type Configuration struct {
	secretKey                  []byte
	authTokenValidityPeriod    int
	refreshTokenValidityPeriod int
	cacheExpirationTime        time.Duration
}

func (config *Configuration) GetSecretKey() []byte {
	return config.secretKey
}

func (config *Configuration) GetAuthTokenValidityPeriod() int {
	return config.authTokenValidityPeriod
}

func (config *Configuration) GetRefreshTokenValidityPeriod() int {
	return config.refreshTokenValidityPeriod
}

func (config *Configuration) GetCacheExpirationTime() time.Duration {
	return config.cacheExpirationTime
}

// TODO: Maybe read the config only once on init() or something and then return the global object?
func GetConfig(fileName string) *Configuration {
	return readServiceConfig(fileName)
}

func readServiceConfig(fileName string) *Configuration {
	viper.SetConfigName(fileName)
	viper.AddConfigPath("$GOPATH/src/github.com/VirrageS/chirp/backend")

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("Error reading config file.")
	}

	configSecretKey := viper.GetString("secret_key")
	configAuthTokenValidityPeriod := viper.GetInt("auth_token_validity_period")
	configRefreshTokenValidityPeriod := viper.GetInt("refresh_token_validity_period")
	configCacheExpirationTime := viper.GetDuration("cache_expiration_time")

	if configSecretKey == "" || configAuthTokenValidityPeriod <= 0 || configRefreshTokenValidityPeriod <= 0 {
		log.WithFields(log.Fields{
			"secret key":              configSecretKey,
			"auth validity period":    configAuthTokenValidityPeriod,
			"refresh validity period": configRefreshTokenValidityPeriod,
			"cache expiration time":   configCacheExpirationTime,
		}).Fatal("Config file doesn't contain valid data.")
	}

	return &Configuration{
		secretKey:                  []byte(configSecretKey),
		authTokenValidityPeriod:    configAuthTokenValidityPeriod,
		refreshTokenValidityPeriod: configRefreshTokenValidityPeriod,
		cacheExpirationTime:        configCacheExpirationTime,
	}
}
