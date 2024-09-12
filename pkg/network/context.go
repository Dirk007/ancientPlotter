package network

import (
	"github.com/Dirk007/ancientPlotter/pkg/jobs"
	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
	*jobs.ContextDependencies
	Config jobs.JobConfig
}

type LoggingContext struct {
	deps   *jobs.ContextDependencies
	config jobs.JobConfig
}

func NewLoggingContext(deps *jobs.ContextDependencies, config jobs.JobConfig) *LoggingContext {
	return &LoggingContext{
		deps,
		config,
	}
}

func (lc *LoggingContext) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := Context{
			Context:             c,
			ContextDependencies: lc.deps,
			Config:              lc.config,
		}
		return next(ctx)
	}
}
