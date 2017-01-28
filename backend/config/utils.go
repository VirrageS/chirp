package config

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type serverConfig struct {
	secretKey                  []byte
	authTokenValidityPeriod    time.Duration
	refreshTokenValidityPeriod time.Duration
	randomPasswordLength       int
}

func (config *serverConfig) GetSecretKey() []byte {
	return config.secretKey
}

func (config *serverConfig) GetAuthTokenValidityPeriod() time.Duration {
	return config.authTokenValidityPeriod
}

func (config *serverConfig) GetRefreshTokenValidityPeriod() time.Duration {
	return config.refreshTokenValidityPeriod
}

func (config *serverConfig) GetRandomPasswordLength() int {
	return config.randomPasswordLength
}

type databaseConfig struct {
	username string
	password string
	host     string
	port     string
}

func (config *databaseConfig) GetUsername() string {
	return config.username
}

func (config *databaseConfig) GetPassword() string {
	return config.password
}

func (config *databaseConfig) GetHost() string {
	return config.host
}

func (config *databaseConfig) GetPort() string {
	return config.port
}

type redisCacheConfig struct {
	password       string
	host           string
	port           string
	db             int
	expirationTime time.Duration
}

func (config *redisCacheConfig) GetPassword() string {
	return config.password
}

func (config *redisCacheConfig) GetHost() string {
	return config.host
}

func (config *redisCacheConfig) GetPort() string {
	return config.port
}

func (config *redisCacheConfig) GetDB() int {
	return config.db
}

func (config *redisCacheConfig) GetExpirationTime() time.Duration {
	return config.expirationTime
}

type authorizationGoogleConfig struct {
	clientID     string
	clientSecret string
	callbackURI  string
	authURL      string
	tokenURL     string
}

func (config *authorizationGoogleConfig) GetClientID() string {
	return config.clientID
}

func (config *authorizationGoogleConfig) GetClientSecret() string {
	return config.clientSecret
}

func (config *authorizationGoogleConfig) GetCallbackURI() string {
	return config.callbackURI
}

func (config *authorizationGoogleConfig) GetAuthURL() string {
	return config.authURL
}

func (config *authorizationGoogleConfig) GetTokenURL() string {
	return config.tokenURL
}

type elasticsearchConfig struct {
	username string
	password string
	host     string
	port     string
}

func (config *elasticsearchConfig) GetUsername() string {
	return config.username
}

func (config *elasticsearchConfig) GetPassword() string {
	return config.password
}

func (config *elasticsearchConfig) GetHost() string {
	return config.host
}

func (config *elasticsearchConfig) GetPort() string {
	return config.port
}

type generalConfig struct {
	*viper.Viper
}

func (config *generalConfig) getServerConfig() *serverConfig {
	secretKey := config.GetString("secret_key")
	authTokenValidityPeriod := config.GetDuration("auth_token_validity_period")
	refreshTokenValidityPeriod := config.GetDuration("refresh_token_validity_period")
	randomPasswordLength := config.GetInt("random_password_length")

	if secretKey == "" || authTokenValidityPeriod <= 0 || refreshTokenValidityPeriod <= 0 || randomPasswordLength <= 0 {
		log.WithFields(log.Fields{
			"secret key":              secretKey,
			"auth validity period":    authTokenValidityPeriod,
			"refresh validity period": refreshTokenValidityPeriod,
			"random password length":  randomPasswordLength,
		}).Fatal("Config file doesn't contain valid data.")
	}

	return &serverConfig{
		secretKey:                  []byte(secretKey),
		authTokenValidityPeriod:    authTokenValidityPeriod,
		refreshTokenValidityPeriod: refreshTokenValidityPeriod,
		randomPasswordLength:       randomPasswordLength,
	}
}

func (config *generalConfig) getDatabaseConfig() *databaseConfig {
	username := config.GetString("database.username")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	port := config.GetString("database.port")

	if username == "" || password == "" || host == "" || port == "" {
		log.WithFields(log.Fields{
			"username": username,
			"password": password,
			"host":     host,
			"port":     port,
		}).Fatal("Config file doesn't contain valid database access data.")
	}

	return &databaseConfig{
		username: username,
		password: password,
		host:     host,
		port:     port,
	}
}

func (config *generalConfig) getRedisCacheConfig() *redisCacheConfig {
	password := config.GetString("redis.password")
	host := config.GetString("redis.host")
	port := config.GetString("redis.port")
	db := config.GetInt("redis.db")
	expirationTime := config.GetDuration("redis.expiration_time")

	if host == "" || port == "" || db < 0 || expirationTime < 0 {
		log.WithFields(log.Fields{
			"password":        password,
			"host":            host,
			"port":            port,
			"db":              db,
			"expiration time": expirationTime,
		}).Fatal("Config file doesn't contain valid redis access data.")
	}

	return &redisCacheConfig{
		password:       password,
		host:           host,
		port:           port,
		db:             db,
		expirationTime: expirationTime,
	}
}

func (config *generalConfig) getAuthorizationGoogleConfig() *authorizationGoogleConfig {
	clientID := config.GetString("authorization_google.client_id")
	clientSecret := config.GetString("authorization_google.client_secret")
	callbackURI := config.GetString("authorization_google.callback_uri")
	authURL := config.GetString("authorization_google.auth_url")
	tokenURL := config.GetString("authorization_google.token_url")

	if clientID == "" || clientSecret == "" || callbackURI == "" || authURL == "" || tokenURL == "" {
		log.WithFields(log.Fields{
			"client_id":     clientID,
			"client_secret": clientSecret,
			"callback_uri":  callbackURI,
			"auth_url":      authURL,
			"token_url":     tokenURL,
		}).Fatal("Config file doesn't contain valid data.")
	}

	return &authorizationGoogleConfig{
		clientID:     clientID,
		clientSecret: clientSecret,
		callbackURI:  callbackURI,
		authURL:      authURL,
		tokenURL:     tokenURL,
	}
}

func (config *generalConfig) getElasticsearchConfig() *elasticsearchConfig {
	username := config.GetString("elasticsearch.username")
	password := config.GetString("elasticsearch.password")
	host := config.GetString("elasticsearch.host")
	port := config.GetString("elasticsearch.port")

	if username == "" || password == "" || host == "" || port == "" {
		log.WithFields(log.Fields{
			"username": username,
			"password": password,
			"host":     host,
			"port":     port,
		}).Fatal("Config file doesn't contain valid elasticsearch access data.")
	}

	return &elasticsearchConfig{
		username: username,
		password: password,
		host:     host,
		port:     port,
	}
}
