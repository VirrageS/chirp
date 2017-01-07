package config

import "time"

// Interfaces that provide only those parameters that are required by different parts of the system

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

// Provides DB access configuration
type DBConfigProvider interface {
	GetUsername() string
	GetPassword() string
	GetHost() string
	GetPort() string
}

// Provides configuration for CacheProvider
type CacheConfigProvider interface {
	GetPassword() string
	GetHost() string
	GetPort() string
	GetCacheExpirationTime() time.Duration
}

// Provides full redis configuration
type RedisConfigProvider interface {
	CacheConfigProvider
	GetDB() int
}

// Provides Elasticsearch access configuration
type ElasticsearchConfigProvider interface {
	GetUsername() string
	GetPassword() string
	GetHost() string
	GetPort() string
}

type AuthorizationGoogleConfigurationProvider interface {
	GetClientID() string
	GetClientSecret() string
	GetCallbackURI() string
	GetAuthURL() string
	GetTokenURL() string
}
