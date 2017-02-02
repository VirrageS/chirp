package config

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/VirrageS/chirp/backend/utils"
)

// Configuration keeps all config providers
type Configuration struct {
	Token               TokenConfigProvider
	Password            PasswordConfigProvider
	Postgres            PostgresConfigProvider
	Redis               RedisConfigProvider
	Elasticsearch       ElasticsearchConfigProvider
	AuthorizationGoogle AuthorizationGoogleConfigProvider
}

// New reads and creates configuration from path provided in env `$CHIRP_CONFIG_PATH`
// or from default path `$GOPATH/src/github.com/VirrageS/chirp/backend`.
//
// By setting `$CHIRP_CONFIG_NAME` variable you can specify the name
// of config file. Default name is `config`.
//
// By setting `$CHIRP_CONFIG_TYPE` variable you can specify which config type will
// be chosen: `development`, `production` or `test`. Default is `development`.
func New() *Configuration {
	v := viper.New()
	if cp := os.Getenv(`CHIRP_CONFIG_PATH`); cp != "" {
		v.AddConfigPath(cp)
	}
	v.AddConfigPath(`$GOPATH/src/github.com/VirrageS/chirp/backend`)

	cn := utils.GetenvOrDefault(`CHIRP_CONFIG_NAME`, "config")
	v.SetConfigName(cn)

	err := v.ReadInConfig()
	if err != nil {
		log.WithError(err).Error("Error reading config file.")
		return nil
	}

	ct := utils.GetenvOrDefault(`CHIRP_CONFIG_TYPE`, "development")
	subConfig := v.Sub(ct)
	if subConfig == nil {
		log.Errorf("Failed to read '%s' config.", ct)
		return nil
	}

	config := &generalConfig{subConfig}
	return &Configuration{
		Token:               config.getServerConfig(),
		Password:            config.getServerConfig(),
		Postgres:            config.getPostgresConfig(),
		Redis:               config.getRedisConfig(),
		Elasticsearch:       config.getElasticsearchConfig(),
		AuthorizationGoogle: config.getAuthorizationGoogleConfig(),
	}
}
