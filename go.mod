module github.com/MikMuellerDev/QuickClip

go 1.17

require (
	github.com/MikMuellerDev/QuickClip/middleware v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/routes v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/sessions v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/templates v0.0.0-00010101000000-000000000000
	github.com/MikMuellerDev/QuickClip/utils v0.0.0-00010101000000-000000000000
)

require golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect

replace (
	github.com/MikMuellerDev/QuickClip/middleware => ./cmd/middleware
	github.com/MikMuellerDev/QuickClip/routes => ./cmd/routes
	github.com/MikMuellerDev/QuickClip/sessions => ./cmd/sessions
	github.com/MikMuellerDev/QuickClip/templates => ./cmd/templates
	github.com/MikMuellerDev/QuickClip/utils => ./cmd/utils
)
