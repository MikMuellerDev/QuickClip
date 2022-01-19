package routes

type ResponseStruct struct {
	Success   bool
	ErrorCode int
	Title     string
	Message   string
}

type VersionStruct struct {
	Version    string
	Production bool
}

type UserResponse struct {
	User     string
	LoggedIn bool
}
