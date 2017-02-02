package config

import "time"

// TokenConfigProvider provides Token access configuration.
type TokenConfigProvider interface {
	GetSecretKey() []byte
	GetAuthTokenValidityPeriod() time.Duration
	GetRefreshTokenValidityPeriod() time.Duration
}

// PasswordConfigProvider provides Password access configuration.
type PasswordConfigProvider interface {
	GetRandomPasswordLength() int
}

// DatabaseConfigProvider provides general Database access configuration.
type DatabaseConfigProvider interface {
	GetUsername() string
	GetPassword() string
	GetHost() string
	GetPort() string
}

// PostgresConfigProvider provides Postgres access configuration.
type PostgresConfigProvider interface {
	DatabaseConfigProvider
}

// CacheConfigProvider provides general Cache access configuration.
type CacheConfigProvider interface {
	GetPassword() string
	GetHost() string
	GetPort() string
	GetExpirationTime() time.Duration
}

// RedisConfigProvider provides Redis access configuration.
type RedisConfigProvider interface {
	CacheConfigProvider
	GetDB() int
}

// AuthorizationGoogleConfigProvider provides Google authorization access configuration.
type AuthorizationGoogleConfigProvider interface {
	GetClientID() string
	GetClientSecret() string
	GetCallbackURI() string
	GetAuthURL() string
	GetTokenURL() string
}

// ElasticsearchConfigProvider provides Elasticsearch access configuration.
type ElasticsearchConfigProvider interface {
	GetUsername() string
	GetPassword() string
	GetHost() string
	GetPort() string
}
