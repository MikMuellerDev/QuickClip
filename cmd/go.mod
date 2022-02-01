module github.com/MikMuellerDev/QuickClip

go 1.17

require (
	github.com/MikMuellerDev/QuickClip/middleware v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/routes v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/sessions v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/templates v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
)

replace (
	github.com/MikMuellerDev/QuickClip/middleware => ./middleware
	github.com/MikMuellerDev/QuickClip/routes => ./routes
	github.com/MikMuellerDev/QuickClip/sessions => ./sessions
	github.com/MikMuellerDev/QuickClip/templates => ./templates
	github.com/MikMuellerDev/QuickClip/utils => ./utils
)
