package main

import (
	"fmt"
	"net/http"

	"github.com/MikMuellerDev/QuickClip/middleware"
	"github.com/MikMuellerDev/QuickClip/routes"
	"github.com/MikMuellerDev/QuickClip/sessions"
	"github.com/MikMuellerDev/QuickClip/templates"
	utils "github.com/MikMuellerDev/QuickClip/utils"
)

func main() {
	// Initialize all loggers
	log := utils.NewLogger()
	middleware.InitLogger(log)
	routes.InitLogger(log)
	sessions.InitLogger(log)
	templates.InitLogger(log)
	utils.InitLogger(log)
	log.Debug("Logging initialized.")

	utils.ReadClipFile()
	utils.RequestSave()
	utils.ReadConfigFile()

	config := utils.GetConfig()
	config.Version = "0.8.1"
	r := routes.NewRouter()
	middleware.InitializeLogin(config)

	sessions.Init(config.Production)
	templates.LoadTemplates("../templates/*.html")
	http.Handle("/", r)
	log.Info(fmt.Sprintf("\x1b[34mQuickClip [Version %s] [%s] is running on http://localhost:%d", config.Version, config.InstanceName, config.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
