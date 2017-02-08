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
	Database            DatabaseConfigProvider
	Redis               RedisConfigProvider
	AuthorizationGoogle AuthorizationGoogleConfigProvider
	Elasticsearch       ElasticsearchConfigProvider
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

	serverConfiguration := config.getServerConfig()
	return &Configuration{
		Token:               serverConfiguration,
		Password:            serverConfiguration,
		Database:            config.getDatabaseConfig(),
		Redis:               config.getRedisCacheConfig(),
		AuthorizationGoogle: config.getAuthorizationGoogleConfig(),
		Elasticsearch:       config.getElasticsearchConfig(),
	}
}
