package network

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func handleStartJob(c echo.Context) error {
	depContext, ok := c.(Context)
	if !ok {
		return fmt.Errorf("invalid context of type %T. This is a bug", c)
	}

	jobID := c.FormValue("job_id")
	job, err := depContext.Jobs.Get(jobID)
	if err != nil {
		return err
	}

	depContext.Logs.Broadcast(context.Background(), fmt.Sprintf("Starting job: %+v", job))

	ctx, cancel := context.WithCancel(context.Background())
	job.Cancel = &cancel
	go job.Run(ctx, depContext.ContextDependencies, depContext.Config)

	// m := map[string]any{
	// 	"job":      job,
	// 	"datetime": fmtDateTime(time.Now()),
	// 	"config":   depContext.Config,
	// }

	// line, err := RenderTemplate("onstarted", m)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func handleStopJob(c echo.Context) error {
	depContext, ok := c.(Context)
	if !ok {
		return fmt.Errorf("invalid context of type %T. This is a bug", c)
	}

	jobID := c.FormValue("job_id")
	job, err := depContext.Jobs.Get(jobID)
	if err != nil {
		return err
	}

	if job.Cancel != nil {
		depContext.Logs.Broadcast(context.Background(), fmt.Sprintf("Calling cancel at %p for job by context", job.Cancel))
		(*job.Cancel)()
	}

	depContext.Logs.Broadcast(context.Background(), fmt.Sprintf("Stopped job: %+v", job))

	m := map[string]any{
		"job":      job,
		"datetime": fmtDateTime(time.Now()),
		"config":   depContext.Config,
	}

	line, err := RenderTemplate("onstopped", m)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, line)
}
