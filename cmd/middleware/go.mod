module github.com/MikMuellerDev/QuickClip/midddleware

go 1.17

replace (
	github.com/MikMuellerDev/QuickClip/sessions => ../sessions
	github.com/MikMuellerDev/QuickClip/utils => ../utils
)

require (
	github.com/MikMuellerDev/QuickClip/sessions v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/utils v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.8.1
)

require (
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)
