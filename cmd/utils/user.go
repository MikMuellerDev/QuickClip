package utils

type User struct {
	Name         string
	Password     string
	Permissions  []string
	WriteAllowed []string
}

type Password struct {
	Password string
}

func HasPermission(username string, permissionToCheck string) bool {
	if username == "admin" {
		return true
	}

	for _, user := range config.Users {
		if user.Name == username {
			for _, permission := range user.Permissions {
				if permission == permissionToCheck || permission == "*" {
					return true
				}
			}
		}
	}
	return false
}

func HasWritePermission(username string, permissionToCheck string) bool {
	if username == "admin" {
		return true
	}

	for _, user := range config.Users {
		if user.Name == username {
			for _, permission := range user.WriteAllowed {
				if permission == permissionToCheck || permission == "*" {
					return true
				}
			}
		}
	}
	return false
}

func DoesUserExist(username string) bool {
	for _, user := range config.Users {
		if user.Name == username {
			return true
		}
	}
	return false
}

// Returns a boolean for indicating the success
func AddUser(user User) bool {
	if DoesUserExist(user.Name) {
		return false
	}
	config.Users = append(config.Users, user)
	return WriteConfigFile()
}

func DeleteUser(username string) bool {
	if !DoesUserExist(username) {
		return false
	}
	var newUsers []User
	for _, v := range config.Users {
		if v.Name != username {
			newUsers = append(newUsers, v)
		}
	}
	config.Users = newUsers
	return WriteConfigFile()
}

func AlterUser(username string, newUser User) bool {
	if !DoesUserExist(username) {
		return false
	}
	var newUsers []User
	for _, v := range config.Users {
		if v.Name != username {
			newUsers = append(newUsers, v)
		} else {
			newUsers = append(newUsers, newUser)
		}
	}
	config.Users = newUsers
	return WriteConfigFile()
}

func AlterPassword(username string, newPassword string) bool {
	if !DoesUserExist(username) {
		return false
	}
	var newUsers []User
	for _, user := range config.Users {
		if user.Name != username {
			newUsers = append(newUsers, user)
		} else {
			user.Password = newPassword
			newUsers = append(newUsers, user)
		}
	}
	config.Users = newUsers
	return WriteConfigFile()
}

func GetUsers() []User {
	return config.Users
}
