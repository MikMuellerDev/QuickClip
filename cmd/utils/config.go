package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type Config struct {
	Version      string
	Production   bool
	Port         int
	Users        []User
	InstanceName string
}

var config Config
var configMutex sync.RWMutex

func ReadConfigFile() {
	path := "../config/config.json"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Error("Error during Unmarshal(), Invalid Json Config file: ", err)
		if string(content) == "" || string(content) == " " {
			writeEmergencyConfigFile()
			log.Warn("[Config file empty] Loaded QuickClip config File from recovery preset.")
		} else {
			log.Fatal("Malformed (non-empty) config file, halting server.")
		}
		return
	}
	log.Debug(fmt.Sprintf("Loaded QuickClip config File from %s", path))
}

// The admin user (with a random password) is created afterwards by EnsureAdminUser()
func writeEmergencyConfigFile() {
	config = Config{Production: true, Port: 80, Users: []User{}, InstanceName: "QuickClip"}
	if !writeConfig() {
		log.Fatal("[Write] Error writing emergency config.")
	}
	log.Debug("Written emergency config contents to config.json.")
}

// Caller must hold configMutex
func writeConfig() bool {
	configJson, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Error("Error during marshal: ", err.Error())
		return false
	}
	err = ioutil.WriteFile("../config/config.json", configJson, 0600)
	if err != nil {
		log.Error(fmt.Sprintf("Error writing configuration file: %s", err.Error()))
		return false
	}
	log.Debug("Written configuration file to config.json")
	return true
}

func WriteConfigFile() bool {
	configMutex.Lock()
	defer configMutex.Unlock()
	return writeConfig()
}

func GetConfig() *Config {
	return &config
}

func GetVersion() (string, bool) {
	return config.Version, config.Production
}
