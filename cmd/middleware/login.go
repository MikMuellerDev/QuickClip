package middleware

import (
	"fmt"

	"github.com/MikMuellerDev/QuickClip/utils"
)

var users []utils.User

func InitializeLogin(config *utils.Config) {
	users = config.Users
}

func TestCredentials(user string, password string, silent bool) bool {
	if !silent {
		log.Trace(fmt.Sprintf("Testing login credentials for User %s", user))
	}
	for _, v := range users {
		if v.Name == user && v.Password == password {
			if !silent {
				log.Debug(fmt.Sprintf("Login successful for User %s", user))
			}
			return true
		}
	}
	if !silent {
		log.Debug(fmt.Sprintf("Login failed for User %s with Password: %s", user, password))
	}
	return false
}
