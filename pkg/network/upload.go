package network

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Dirk007/ancientPlotter/pkg/jobs"
	"github.com/labstack/echo"
)

func handleUpload(c echo.Context) error {
	depContext, ok := c.(Context)
	if !ok {
		return fmt.Errorf("invalid context of type %T. This is a bug", c)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	depContext.Logs.Broadcast(context.Background(), fmt.Sprintf("Incoming file %s", file.Filename))
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	tempDir, err := os.MkdirTemp("", "plot.*.hpgl")
	if err != nil {
		return err
	}

	depContext.Logs.Broadcast(context.Background(), fmt.Sprintf("Saving to %s/%s", tempDir, file.Filename))
	dst, err := os.Create(fmt.Sprintf("%s/%s", tempDir, file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	job := jobs.NewPlotJob(dst.Name())
	depContext.Jobs.Add(job.ID, job.Path)

	depContext.Logs.Broadcast(context.Background(), fmt.Sprintf("Job ID %s saved", job.ID))

	m := map[string]any{
		"job":      job,
		"datetime": fmtDateTime(time.Now()),
		"config":   depContext.Config,
		"file":     file,
	}

	line, err := RenderTemplate("onuploaded", m)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, line)
}
