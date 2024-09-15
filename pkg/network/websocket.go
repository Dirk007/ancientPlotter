package network

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/Dirk007/ancientPlotter/pkg/feeder"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
)

func RenderTemplate(filename string, data any) (string, error) {
	completeName := fmt.Sprintf("assets/templates/%s.html", filename)
	tmpl, err := template.ParseFiles(completeName)
	if err != nil {
		return "", err
	}
	target := new(strings.Builder)
	err = tmpl.Execute(target, data)
	return target.String(), err
}

func fmtDateTime(dt time.Time) string {
	return dt.Format("02.01.2006 15:04:05")
}

func handleAlive(ws *websocket.Conn, msg string) error {
	m := map[string]any{
		"datetime": fmtDateTime(time.Now()),
		"state":    msg,
	}

	line, err := RenderTemplate("alive", m)
	if err != nil {
		return err
	}
	err = websocket.Message.Send(ws, line)
	return err
}

func handleLog(ws *websocket.Conn, log string) error {
	line := fmt.Sprintf("<span hx-swap-oob='afterbegin:#notifications'>%s: %s<br/></span>", fmtDateTime(time.Now()), log)
	err := websocket.Message.Send(ws, line)
	return err
}

func handleStat(ws *websocket.Conn, stat feeder.Stats) error {
	if stat.FatalError != nil {
		line, err := RenderTemplate("onerror", map[string]string{"error": stat.FatalError.Error()})
		if err != nil {
			return err
		}
		err = websocket.Message.Send(ws, line)
		_ = handleLog(ws, fmt.Sprintf("!!!!! ERROR: %+v", stat.FatalError))
		return err
	}

	if stat.Line == 1 {
		line, err := RenderTemplate("onstarted", map[string]any{"jobID": stat.JobID})
		if err != nil {
			return err
		}
		_ = websocket.Message.Send(ws, line)
	}

	percent := 100.0 / float64(stat.Total) * float64(stat.Line)
	m := map[string]any{
		"jobID":        stat.JobID,
		"percent":      percent,
		"line":         stat.Line,
		"total":        stat.Total,
		"currentTry":   stat.CurrentTry,
		"currentTotal": stat.CurrentTotal,
		"currentRest":  stat.CurrentRest,
		"datetime":     fmtDateTime(time.Now()),
	}

	line, err := RenderTemplate("statusbar", m)
	if err != nil {
		return err
	}
	err = websocket.Message.Send(ws, line)
	return err
}

func handleWebsocket(c echo.Context) error {
	logContext, ok := c.(Context)
	if !ok {
		return fmt.Errorf("invalid context of type %T. This is a bug", c)
	}

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		idAlive, alive := logContext.Alive.Register()
		idLog, logs := logContext.Logs.Register()
		idStats, stats := logContext.Stats.Register()
		defer logContext.Alive.Remove(idAlive)
		defer logContext.Logs.Remove(idLog)
		defer logContext.Stats.Remove(idStats)
		for {
			var err error
			select {
			case <-c.Request().Context().Done():
				return
			case alive := <-alive:
				err = handleAlive(ws, alive)
			case log := <-logs:
				err = handleLog(ws, log)
			case stats := <-stats:
				err = handleStat(ws, stats)
			}
			if err != nil {
				c.Logger().Error(err)
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
