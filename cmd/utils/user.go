package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

type User struct {
	Name         string
	Password     string
	Permissions  []string
	WriteAllowed []string
}

type Password struct {
	Password string
}

func randomPassword() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Makes sure an admin user exists and warns about default passwords.
// Must be called once at startup, after ReadConfigFile().
func EnsureAdminUser() {
	configMutex.Lock()
	defer configMutex.Unlock()

	if !userExists("admin") {
		password, err := randomPassword()
		if err != nil {
			log.Fatal("Could not generate admin password: ", err)
		}
		config.Users = append(config.Users, User{Name: "admin", Password: password, Permissions: []string{"*"}, WriteAllowed: []string{"*"}})
		// Printed to stdout only, so that it does not end up in the log files
		fmt.Printf("\x1b[33mNo admin user was configured. Created user 'admin' with password: %s\nChange this password after logging in.\x1b[0m\n", password)
		if !writeConfig() {
			log.Fatal("Could not save the admin user.")
		}
	}

	for _, user := range config.Users {
		if user.Password == "password" {
			log.Warn(fmt.Sprintf("\x1b[33mUser %q uses the default password 'password'. Change it immediately.", user.Name))
		}
	}
}

func CheckCredentials(username string, password string) bool {
	configMutex.RLock()
	defer configMutex.RUnlock()
	for _, user := range config.Users {
		if user.Name == username {
			return subtle.ConstantTimeCompare([]byte(user.Password), []byte(password)) == 1
		}
	}
	return false
}

func hasPermission(username string, permissionToCheck string, write bool) bool {
	if username == "admin" {
		return true
	}

	configMutex.RLock()
	defer configMutex.RUnlock()
	for _, user := range config.Users {
		if user.Name == username {
			permissions := user.Permissions
			if write {
				permissions = user.WriteAllowed
			}
			for _, permission := range permissions {
				if permission == permissionToCheck || permission == "*" {
					return true
				}
			}
		}
	}
	return false
}

func HasPermission(username string, permissionToCheck string) bool {
	return hasPermission(username, permissionToCheck, false)
}

func HasWritePermission(username string, permissionToCheck string) bool {
	return hasPermission(username, permissionToCheck, true)
}

// Caller must hold configMutex
func userExists(username string) bool {
	for _, user := range config.Users {
		if user.Name == username {
			return true
		}
	}
	return false
}

func DoesUserExist(username string) bool {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return userExists(username)
}

// Returns a boolean for indicating the success
func AddUser(user User) bool {
	configMutex.Lock()
	defer configMutex.Unlock()
	if userExists(user.Name) {
		return false
	}
	config.Users = append(config.Users, user)
	return writeConfig()
}

func DeleteUser(username string) bool {
	if username == "admin" {
		return false
	}
	configMutex.Lock()
	defer configMutex.Unlock()
	if !userExists(username) {
		return false
	}
	var newUsers []User
	for _, v := range config.Users {
		if v.Name != username {
			newUsers = append(newUsers, v)
		}
	}
	config.Users = newUsers
	return writeConfig()
}

// Renaming users is not supported, the name of newUser is ignored.
// If the new User's password is "?" or empty, don't change it.
func AlterUser(username string, newUser User) bool {
	newUser.Name = username
	keepPassword := newUser.Password == "?" || newUser.Password == ""

	configMutex.Lock()
	defer configMutex.Unlock()
	if !userExists(username) {
		return false
	}
	var newUsers []User
	for _, v := range config.Users {
		if v.Name != username {
			newUsers = append(newUsers, v)
		} else {
			if keepPassword {
				newUser.Password = v.Password
			}
			newUsers = append(newUsers, newUser)
		}
	}
	config.Users = newUsers
	return writeConfig()
}

func AlterPassword(username string, newPassword string) bool {
	configMutex.Lock()
	defer configMutex.Unlock()
	if !userExists(username) {
		return false
	}
	for i, user := range config.Users {
		if user.Name == username {
			config.Users[i].Password = newPassword
		}
	}
	return writeConfig()
}

// Returns a copy of all users with their passwords removed
func GetUsers() []User {
	configMutex.RLock()
	defer configMutex.RUnlock()
	users := make([]User, len(config.Users))
	for i, user := range config.Users {
		users[i] = User{
			Name:         user.Name,
			Permissions:  append([]string{}, user.Permissions...),
			WriteAllowed: append([]string{}, user.WriteAllowed...),
		}
	}
	return users
}
