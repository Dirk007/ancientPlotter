package network

import (
	"fmt"

	"github.com/Dirk007/ancientPlotter/pkg/jobs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Serve(deps *jobs.ContextDependencies, port int, config jobs.JobConfig) {
	fmt.Printf("Starting server with jobconfig %+v\n", config)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(NewLoggingContext(deps, config).Middleware)
	e.Static("/", "assets/public")
	e.POST("/upload", handleUpload)
	e.POST("/action/start", handleStartJob)
	e.POST("/action/stop", handleStopJob)
	e.GET("/ws", handleWebsocket)
	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
