package config

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

const (
	configFilePath = "../data/config.json"
)

type Config struct {
	PubKey    string `json:"OMISE_API_PUBLIC_KEY"`
	TokensUrl string `json:"OMISE_API_TOKENS_URL"`
}

func GetConfig() (*Config, error) {
	config := &Config{}

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
