module github.com/Dirk007/ancientPlotter

go 1.23.0

// Waiting for https://github.com/bugst/go-serial/pull/195 to be merged upstream
replace go.bug.st/serial => github.com/Dirk007/go-serial v0.0.0-20240907112601-b7483e31a79c

require (
	github.com/fred1268/go-clap v1.2.1
	github.com/sirupsen/logrus v1.9.3
	go.bug.st/serial v1.6.2
)

require (
	github.com/creack/goselect v0.1.2 // indirect
	golang.org/x/sys v0.19.0 // indirect
)
