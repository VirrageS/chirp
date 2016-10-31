package config

import (
	"os"
	"fmt"
	"encoding/json"
	"runtime"
	"path"
)

type configuration struct {
	SecretKey string
}

var config *configuration = initializeConfiguration()

func GetSecretKey() string {
	return config.SecretKey
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
	if readConfig.SecretKey == "" {
		panic("Secret key was not read properly!")
	}

	return &readConfig
}
