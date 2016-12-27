package config

// TODO: config file path should probably be read from env variable or command line argument
// TODO: rename this file

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Stores global configuration of the system. Implements all config interfaces
type Configuration struct {
	secretKey                  []byte
	authTokenValidityPeriod    int
	refreshTokenValidityPeriod int
	cacheExpirationTime        time.Duration

	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     string

	redisPassword string
	redisHost     string
	redisPort     string
	redisDB       int
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

func (config *Configuration) GetDBUsername() string {
	return config.DBUsername
}

func (config *Configuration) GetDBPassword() string {
	return config.DBPassword
}

func (config *Configuration) GetDBHost() string {
	return config.DBHost
}

func (config *Configuration) GetDBPort() string {
	return config.DBPort
}

func (config *Configuration) GetRedisPassword() string {
	return config.redisPassword
}

func (config *Configuration) GetRedisHost() string {
	return config.redisHost
}

func (config *Configuration) GetRedisPort() string {
	return config.redisPort
}

func (config *Configuration) GetRedisDB() int {
	return config.redisDB
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

	configDBUsername := viper.GetString("database.username")
	configDBPassword := viper.GetString("database.password")
	configDBHost := viper.GetString("database.host")
	configDBPort := viper.GetString("database.port")

	configRedisPassword := viper.GetString("redis.password")
	configRedisHost := viper.GetString("redis.host")
	configRedisPort := viper.GetString("redis.port")
	configRedisDB := viper.GetInt("redis.db")

	if configSecretKey == "" || configAuthTokenValidityPeriod <= 0 || configRefreshTokenValidityPeriod <= 0 ||
		configRedisPort == "" {
		log.WithFields(log.Fields{
			"secret key":              configSecretKey,
			"auth validity period":    configAuthTokenValidityPeriod,
			"refresh validity period": configRefreshTokenValidityPeriod,
			"cache expiration time":   configCacheExpirationTime,
		}).Fatal("Config file doesn't contain valid data.")
	}

	if configDBUsername == "" || configDBPassword == "" || configDBHost == "" || configDBPort == "" {
		log.WithFields(log.Fields{
			"username": configDBUsername,
			"password": configDBPassword,
			"host":     configDBHost,
			"port":     configDBPort,
		}).Fatal("Config file doesn't contain valid database access data.")
	}

	if configRedisHost == "" || configRedisPort == "" || configRedisDB < 0 {
		log.WithFields(log.Fields{
			"password": configRedisPassword,
			"host":     configRedisHost,
			"port":     configRedisPort,
			"db":       configRedisDB,
		}).Fatal("Config file doesn't contain valid redis access data.")
	}

	return &Configuration{
		secretKey:                  []byte(configSecretKey),
		authTokenValidityPeriod:    configAuthTokenValidityPeriod,
		refreshTokenValidityPeriod: configRefreshTokenValidityPeriod,
		cacheExpirationTime:        configCacheExpirationTime,

		DBUsername: configDBUsername,
		DBPassword: configDBPassword,
		DBHost:     configDBHost,
		DBPort:     configDBPort,

		redisPassword: configRedisPassword,
		redisHost:     configRedisHost,
		redisPort:     configRedisPort,
		redisDB:       configRedisDB,
	}
}
