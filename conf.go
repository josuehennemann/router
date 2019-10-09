package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	LogPath     string `json:"logPath"`     // path do diretorio de log
	DataPath    string `json:"dataPath"`    // path do diretorio data da aplicação
	HttpAddress string `json:"httpAddress"` //endereço http
	LibPath     string `json:"libPath"`
	IsDev       bool
}

func initConfig() error {
	config = &Config{}
	data, err := ioutil.ReadFile(*fileConf)
	if err != nil {
		return err
	}
	config.IsDev = true

	err = json.Unmarshal(data, config)
	if config.HttpAddress != "localhost" {
		config.IsDev = false
	}
	return err
}
