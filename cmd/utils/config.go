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

func writeEmergencyConfigFile() {
	var users []User
	var permissionsAdmin []string
	var writePermissionsAdmin []string

	permissionsAdmin = append(permissionsAdmin, "*")
	writePermissionsAdmin = append(writePermissionsAdmin, "*")

	users = append(users,
		User{
			Name:         "admin",
			Password:     "password",
			Permissions:  permissionsAdmin,
			WriteAllowed: writePermissionsAdmin,
		})

	users = append(users,
		User{
			Name:         "default",
			Password:     "password",
			Permissions:  make([]string, 1),
			WriteAllowed: make([]string, 1),
		})

	config = Config{Production: true, Port: 80, Users: users, InstanceName: "QuickClip"}
	var jsonBlob = []byte(`{}`)
	err := json.Unmarshal(jsonBlob, &config)
	if err != nil {
		log.Fatal("[Write] Error during unmarshal", err.Error())
	}

	configJson, _ := json.MarshalIndent(config, "", "    ")
	err = ioutil.WriteFile("../config/config.json", configJson, 0644)
	if err != nil {
		log.Fatal("[Write] Error writing config: %s", err.Error())
	}
	log.Debug("Written emergency config contents to config.json.")
}

func GetConfig() *Config {
	return &config
}

func GetVersion() (string, bool) {
	return config.Version, config.Production
}
