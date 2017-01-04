package config

// TODO: config file path should probably be read from env variable or command line argument
// TODO: rename this file

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	secretKey                  []byte
	authTokenValidityPeriod    int
	refreshTokenValidityPeriod int
}

func (config *ServerConfiguration) GetSecretKey() []byte {
	return config.secretKey
}

func (config *ServerConfiguration) GetAuthTokenValidityPeriod() int {
	return config.authTokenValidityPeriod
}

func (config *ServerConfiguration) GetRefreshTokenValidityPeriod() int {
	return config.refreshTokenValidityPeriod
}

type DatabaseConfiguration struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (dbConfig *DatabaseConfiguration) GetUsername() string {
	return dbConfig.Username
}

func (dbConfig *DatabaseConfiguration) GetPassword() string {
	return dbConfig.Password
}

func (dbConfig *DatabaseConfiguration) GetHost() string {
	return dbConfig.Host
}

func (dbConfig *DatabaseConfiguration) GetPort() string {
	return dbConfig.Port
}

type RedisCacheConfiguration struct {
	password            string
	host                string
	port                string
	db                  int
	cacheExpirationTime time.Duration
}

func (cacheConfig *RedisCacheConfiguration) GetPassword() string {
	return cacheConfig.password
}

func (cacheConfig *RedisCacheConfiguration) GetHost() string {
	return cacheConfig.host
}

func (cacheConfig *RedisCacheConfiguration) GetPort() string {
	return cacheConfig.port
}

func (cacheConfig *RedisCacheConfiguration) GetDB() int {
	return cacheConfig.db
}

func (cacheConfig *RedisCacheConfiguration) GetCacheExpirationTime() time.Duration {
	return cacheConfig.cacheExpirationTime
}

type AuthorizationGoogleConfiguration struct {
	clientID     string
	clientSecret string
	callbackURI  string
	authURL      string
	tokenURL     string
}

func (config *AuthorizationGoogleConfiguration) GetClientID() string {
	return config.clientID
}

func (config *AuthorizationGoogleConfiguration) GetClientSecret() string {
	return config.clientSecret
}

func (config *AuthorizationGoogleConfiguration) GetCallbackURI() string {
	return config.callbackURI
}

func (config *AuthorizationGoogleConfiguration) GetAuthURL() string {
	return config.authURL
}

func (config *AuthorizationGoogleConfiguration) GetTokenURL() string {
	return config.tokenURL
}

// TODO: Maybe read the config only once on init() or something and then return the global object?
func GetConfig(fileName string) (
	ServiceConfigProvider,
	DBConfigProvider,
	RedisConfigProvider,
	AuthorizationGoogleConfigurationProvider,
) {
	viper.AddConfigPath("$GOPATH/src/github.com/VirrageS/chirp/backend")
	viper.SetConfigName(fileName)

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("Error reading config file.")
	}

	serviceConfig := readServiceConfig()
	databaseConfig := readDatabaseConfig()
	cacheConfig := readCacheConfig()
	authorizationConfig := readAuthorizationConfig()
	return serviceConfig, databaseConfig, cacheConfig, authorizationConfig
}

func readServiceConfig() *ServerConfiguration {
	configSecretKey := viper.GetString("secret_key")
	configAuthTokenValidityPeriod := viper.GetInt("auth_token_validity_period")
	configRefreshTokenValidityPeriod := viper.GetInt("refresh_token_validity_period")

	if configSecretKey == "" || configAuthTokenValidityPeriod <= 0 || configRefreshTokenValidityPeriod <= 0 {
		log.WithFields(log.Fields{
			"secret key":              configSecretKey,
			"auth validity period":    configAuthTokenValidityPeriod,
			"refresh validity period": configRefreshTokenValidityPeriod,
		}).Fatal("Config file doesn't contain valid data.")
	}

	return &ServerConfiguration{
		secretKey:                  []byte(configSecretKey),
		authTokenValidityPeriod:    configAuthTokenValidityPeriod,
		refreshTokenValidityPeriod: configRefreshTokenValidityPeriod,
	}
}

func readDatabaseConfig() *DatabaseConfiguration {
	configDBUsername := viper.GetString("database.username")
	configDBPassword := viper.GetString("database.password")
	configDBHost := viper.GetString("database.host")
	configDBPort := viper.GetString("database.port")

	if configDBUsername == "" || configDBPassword == "" || configDBHost == "" || configDBPort == "" {
		log.WithFields(log.Fields{
			"username": configDBUsername,
			"password": configDBPassword,
			"host":     configDBHost,
			"port":     configDBPort,
		}).Fatal("Config file doesn't contain valid database access data.")
	}

	return &DatabaseConfiguration{
		Username: configDBUsername,
		Password: configDBPassword,
		Host:     configDBHost,
		Port:     configDBPort,
	}
}

func readCacheConfig() *RedisCacheConfiguration {
	password := viper.GetString("redis.password")
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	db := viper.GetInt("redis.db")
	cacheExpirationTime := viper.GetDuration("redis.cache_expiration_time")

	if host == "" || port == "" || db < 0 || cacheExpirationTime < 0 {
		log.WithFields(log.Fields{
			"password":        password,
			"host":            host,
			"port":            port,
			"db":              db,
			"expiration time": cacheExpirationTime,
		}).Fatal("Config file doesn't contain valid redis access data.")
	}

	return &RedisCacheConfiguration{
		password:            password,
		host:                host,
		port:                port,
		db:                  db,
		cacheExpirationTime: cacheExpirationTime,
	}
}

func readAuthorizationConfig() *AuthorizationGoogleConfiguration {
	configClientID := viper.GetString("authorization_google.client_id")
	configClientSecret := viper.GetString("authorization_google.client_secret")
	configCallbackURI := viper.GetString("authorization_google.callback_uri")
	configAuthURL := viper.GetString("authorization_google.auth_url")
	configTokenURL := viper.GetString("authorization_google.token_url")

	if configClientID == "" || configClientSecret == "" {
		log.WithFields(log.Fields{
			"client_id":     configClientID,
			"client_secret": configClientSecret,
			"callback_uri":  configCallbackURI,
			"auth_url":      configAuthURL,
			"token_url":     configTokenURL,
		}).Fatal("Config file doesn't contain valid data.")
	}

	return &AuthorizationGoogleConfiguration{
		clientID:     configClientID,
		clientSecret: configClientSecret,
		callbackURI:  configCallbackURI,
		authURL:      configAuthURL,
		tokenURL:     configTokenURL,
	}
}
