package routes

import (
	"fmt"
	"net/http"

	middleware "github.com/MikMuellerDev/QuickClip/middleware"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var log *logrus.Logger

func InitLogger(logger *logrus.Logger) {
	log = logger
}

// Initializes a new Router, used in main.go
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", middleware.LogRequest(indexGetHandler)).Methods("GET")
	r.HandleFunc("/admin", middleware.AuthRequired(middleware.LogRequest(adminGetHandler))).Methods("GET")
	r.HandleFunc("/dash", middleware.LogRequest(dashGetHandler)).Methods("GET")
	r.HandleFunc("/edit/{id}", middleware.LogRequest(editGetHandler)).Methods("GET")

	r.HandleFunc("/api/save", middleware.LogRequest(saveToFile)).Methods("PUT")
	r.HandleFunc("/api/version", middleware.LogRequest(getVersion)).Methods("GET")
	r.HandleFunc("/api/clips", middleware.LogRequest(getClips)).Methods("GET")
	r.HandleFunc("/api/clip/{id}", middleware.LogRequest(getClipById)).Methods("GET")
	r.HandleFunc("/api/clip/probe/{id}", middleware.LogRequest(probeWriteAccess)).Methods("GET")
	r.HandleFunc("/api/clip/refresh/{id}", getClipById).Methods("GET")

	r.HandleFunc("/api/user", middleware.LogRequest(getApiUser)).Methods("GET")
	r.HandleFunc("/api/users", middleware.AdminAuthRequired(middleware.LogRequest(getUserList))).Methods("GET")
	r.HandleFunc("/api/user", middleware.AdminAuthRequired(middleware.LogRequest(createUser))).Methods("POST")
	r.HandleFunc("/api/user/{username}", middleware.AdminAuthRequired(middleware.LogRequest())).Methods("PUT")
	r.HandleFunc("/api/user/{username}", middleware.AdminAuthRequired(middleware.LogRequest(getApiUser))).Methods("DELETE")
	r.HandleFunc("/api/password", middleware.ApiAuthRequired(middleware.LogRequest(getApiUser))).Methods("PUT")

	r.HandleFunc("/api/clips/update/{id}", middleware.ApiAuthRequired(middleware.LogRequest(removeClip))).Methods("DELETE")
	r.HandleFunc("/api/clips/edit/{id}", middleware.LogRequest(editClip)).Methods("PUT")
	r.HandleFunc("/api/clips/update", middleware.ApiAuthRequired(middleware.LogRequest(modClip))).Methods("PUT")
	r.HandleFunc("/api/clips/add", middleware.ApiAuthRequired(middleware.LogRequest(addClip))).Methods("POST")

	r.HandleFunc("/login", middleware.LogRequest(loginGetHandler)).Methods("GET")
	r.HandleFunc("/login", middleware.LogRequest(loginPostHandler)).Methods("POST")
	r.HandleFunc("/logout", middleware.LogRequest(logoutGetHandler)).Methods("GET")

	r.HandleFunc("/robots.txt", middleware.LogRequest(robotsTxtHandler)).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	filepath := "../static"
	pathPrefix := "/static/"
	fs := http.FileServer(http.Dir(filepath))
	r.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, fs))
	log.Debug(fmt.Sprintf("Initialized FileServer for directory: %s. with replacement prefix: %s", filepath, pathPrefix))
	log.Debug("Initialized Router.")
	return r
}
