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
