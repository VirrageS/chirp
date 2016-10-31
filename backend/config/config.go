package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var secretKey []byte
var tokenValidityPeriod int

func GetSecretKey() []byte {
	return secretKey
}

func GetTokenValidityPeriod() int {
	return tokenValidityPeriod
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %v\n", err))
	}

	configSecretKey := viper.GetString("secret_key")
	configValidityPeriod := viper.GetInt("token_validity_period")
	if configSecretKey == "" || configValidityPeriod <= 0 {
		panic(fmt.Sprintf("Config file contains invalid data! secretKey = %s, validityPeriod = %d",
			configSecretKey, configValidityPeriod))
	}

	secretKey = []byte(configSecretKey)
	tokenValidityPeriod = configValidityPeriod
}
