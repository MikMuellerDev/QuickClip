module github.com/MikMuellerDev/QuickClip/routes

go 1.17

require (
	github.com/MikMuellerDev/QuickClip/middleware v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/sessions v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/templates v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/utils v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
)

require github.com/sirupsen/logrus v1.8.1

replace (
	github.com/MikMuellerDev/QuickClip/middleware => ../middleware
	github.com/MikMuellerDev/QuickClip/routes => ../routes
	github.com/MikMuellerDev/QuickClip/sessions => ../sessions
	github.com/MikMuellerDev/QuickClip/templates => ../templates
	github.com/MikMuellerDev/QuickClip/utils => ../utils
)
