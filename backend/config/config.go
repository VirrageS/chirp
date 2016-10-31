package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
)

type configuration struct {
	SecretKey           string
	TokenValidityPeriod int
}

var config *configuration = initializeConfiguration()

func GetSecretKey() []byte {
	return []byte(config.SecretKey)
}

func GetTokenValidityPeriod() int {
	return config.TokenValidityPeriod
}

func initializeConfiguration() *configuration {
	_, filename, _, _ := runtime.Caller(1)
	filepath := path.Join(path.Dir(filename), "config.json")

	file, err := os.Open(filepath)
	if err != nil {
		panic("Cannot open config file!")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	readConfig := configuration{}

	err = decoder.Decode(&readConfig)
	if err != nil {
		panic(fmt.Sprintf("Couldn't decode config file! Error = %v.", err))
	}
	if readConfig.SecretKey == "" || readConfig.TokenValidityPeriod <= 0 {
		panic("Config files does not contain valid values!")
	}

	return &readConfig
}
