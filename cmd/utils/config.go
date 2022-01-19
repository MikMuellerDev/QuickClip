package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Version      string
	Production   bool
	Port         int
	Users        []User
	InstanceName string
}

var config Config

func ReadConfigFile() {
	path := "../config/config.json"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	log.Debug(fmt.Sprintf("Loaded QuickClip config File from %s", path))
}

func GetConfig() *Config {
	return &config
}

func GetVersion() (string, bool) {
	return config.Version, config.Production
}
