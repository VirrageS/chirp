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
	username string
	password string
	host     string
	port     string
}

func (dbConfig *DatabaseConfiguration) GetUsername() string {
	return dbConfig.username
}

func (dbConfig *DatabaseConfiguration) GetPassword() string {
	return dbConfig.password
}

func (dbConfig *DatabaseConfiguration) GetHost() string {
	return dbConfig.host
}

func (dbConfig *DatabaseConfiguration) GetPort() string {
	return dbConfig.port
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

type ElasticsearchConfiguration struct {
	username string
	password string
	host     string
	port     string
}

func (esConfig *ElasticsearchConfiguration) GetUsername() string {
	return esConfig.username
}

func (esConfig *ElasticsearchConfiguration) GetPassword() string {
	return esConfig.password
}

func (esConfig *ElasticsearchConfiguration) GetHost() string {
	return esConfig.host
}

func (esConfig *ElasticsearchConfiguration) GetPort() string {
	return esConfig.port
}

// TODO: Maybe read the config only once on init() or something and then return the global object?
func GetConfig(fileName string) (
	ServiceConfigProvider,
	DBConfigProvider,
	RedisConfigProvider,
	AuthorizationGoogleConfigurationProvider,
	ElasticsearchConfigProvider,
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
	elasticsearchConfig := readElasticsearchConfig()

	return serviceConfig, databaseConfig, cacheConfig, authorizationConfig, elasticsearchConfig
}

func readServiceConfig() *ServerConfiguration {
	secretKey := viper.GetString("secret_key")
	authTokenValidityPeriod := viper.GetInt("auth_token_validity_period")
	refreshTokenValidityPeriod := viper.GetInt("refresh_token_validity_period")

	if secretKey == "" || authTokenValidityPeriod <= 0 || refreshTokenValidityPeriod <= 0 {
		log.WithFields(log.Fields{
			"secret key":              secretKey,
			"auth validity period":    authTokenValidityPeriod,
			"refresh validity period": refreshTokenValidityPeriod,
		}).Fatal("Config file doesn't contain valid data.")
	}

	return &ServerConfiguration{
		secretKey:                  []byte(secretKey),
		authTokenValidityPeriod:    authTokenValidityPeriod,
		refreshTokenValidityPeriod: refreshTokenValidityPeriod,
	}
}

func readDatabaseConfig() *DatabaseConfiguration {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")

	if username == "" || password == "" || host == "" || port == "" {
		log.WithFields(log.Fields{
			"username": username,
			"password": password,
			"host":     host,
			"port":     port,
		}).Fatal("Config file doesn't contain valid database access data.")
	}

	return &DatabaseConfiguration{
		username: username,
		password: password,
		host:     host,
		port:     port,
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
	clientID := viper.GetString("authorization_google.client_id")
	clientSecret := viper.GetString("authorization_google.client_secret")
	callbackURI := viper.GetString("authorization_google.callback_uri")
	authURL := viper.GetString("authorization_google.auth_url")
	tokenURL := viper.GetString("authorization_google.token_url")

	if clientID == "" || clientSecret == "" {
		log.WithFields(log.Fields{
			"client_id":     clientID,
			"client_secret": clientSecret,
			"callback_uri":  callbackURI,
			"auth_url":      authURL,
			"token_url":     tokenURL,
		}).Fatal("Config file doesn't contain valid data.")
	}

	return &AuthorizationGoogleConfiguration{
		clientID:     clientID,
		clientSecret: clientSecret,
		callbackURI:  callbackURI,
		authURL:      authURL,
		tokenURL:     tokenURL,
	}
}

func readElasticsearchConfig() *ElasticsearchConfiguration {
	username := viper.GetString("elasticsearch.username")
	password := viper.GetString("elasticsearch.password")
	host := viper.GetString("elasticsearch.host")
	port := viper.GetString("elasticsearch.port")

	if username == "" || password == "" || host == "" || port == "" {
		log.WithFields(log.Fields{
			"username": username,
			"password": password,
			"host":     host,
			"port":     port,
		}).Fatal("Config file doesn't contain valid elasticsearch access data.")
	}

	return &ElasticsearchConfiguration{
		username: username,
		password: password,
		host:     host,
		port:     port,
	}
}
