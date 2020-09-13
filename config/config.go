package config

import (
	"os"
	"sync"
	"io/ioutil"
	"encoding/json"
)

const (
	configFilePath = "../data/config.json"
)

type appConfig struct {
	PubKey    string `json:"OMISE_API_PUBLIC_KEY"`
	SecKey    string `json:"OMISE_API_SECRET_KEY"`
	TokensUrl string `json:"OMISE_API_TOKENS_URL"`
	ChargeUrl string `json:"OMISE_API_CHARGE_URL"`
}

var config *appConfig
var once sync.Once

func GetConfig() *appConfig {
	once.Do(func() {
		// TODO: maybe trowing here is the only option to hande the error
		config, _ = readConfig()
	})

	return config
}

func readConfig() (*appConfig, error) {
	config := &appConfig{}

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		return config, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, config)
	if err != nil {
		return config, err
	}

	return config, nil
}
