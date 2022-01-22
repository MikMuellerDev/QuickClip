package utils

type User struct {
	Name         string
	Password     string
	Permissions  []string
	WriteAllowed []string
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
