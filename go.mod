module github.com/Dirk007/ancientPlotter

go 1.23.0

// Waiting for https://github.com/bugst/go-serial/pull/195 to be merged upstream
replace go.bug.st/serial => github.com/Dirk007/go-serial v0.0.0-20240907112601-b7483e31a79c

require (
	github.com/fred1268/go-clap v1.2.1
	github.com/google/uuid v1.6.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/sirupsen/logrus v1.9.3
	go.bug.st/serial v1.6.2
	golang.org/x/net v0.29.0
)

require (
	github.com/creack/goselect v0.1.2 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
