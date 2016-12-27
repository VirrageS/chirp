package config

import "time"

// Interfaces that provide only those parameters that are required by different parts of the system

// Provides secret key
type SecretKeyProvider interface {
	GetSecretKey() []byte
}

// Provides configuration for CacheProvider
type CacheConfigProvider interface {
	GetCacheExpirationTime() time.Duration
}

// Provides DB access configuration
type DBAccessConfigProvider interface {
	GetDBUsername() string
	GetDBPassword() string
	GetDBHost() string
	GetDBPort() string
}

// Provides Redis access configuration
type RedisAccessConfigProvider interface {
	GetRedisPassword() string
	GetRedisHost() string
	GetRedisPort() string
	GetRedisDB() int
}

// Provides full redis configuration
type RedisConfigProvider interface {
	RedisAccessConfigProvider
	CacheConfigProvider
}

// Provides configuration for ServiceProvider
type ServiceConfigProvider interface {
	SecretKeyProvider
	GetAuthTokenValidityPeriod() int
	GetRefreshTokenValidityPeriod() int
}
