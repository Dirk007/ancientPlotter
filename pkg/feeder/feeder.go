package feeder

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DefaultMaxTries = 5
	DefaultBackoff  = 100 * time.Millisecond
)

var _ json.Marshaler = Stats{}

type Stats struct {
	JobID        string
	FatalError   string
	Line         int
	Total        int
	CurrentTry   int
	CurrentTotal int
	CurrentRest  int
}

func (s Stats) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"line":          s.Line,
		"total":         s.Total,
		"current_try":   s.CurrentTry,
		"current_total": s.CurrentTotal,
		"current_rest":  s.CurrentRest,
	})
}

type (
	WriterFn = func(s string) (written int, err error)
	StatFn   = func(stats Stats) (err error)
)

type Feeder struct {
	maxTries int
	backoff  time.Duration
	jobID    string
}

func New(maxTries int, backoff time.Duration, jobID string) *Feeder {
	return &Feeder{
		maxTries: maxTries,
		backoff:  backoff,
		jobID:    jobID,
	}
}

func (f *Feeder) Feed(ctx context.Context, writeFn WriterFn, instructions []string, statFn StatFn) error {
	for index, instruction := range instructions {
		instruction := strings.TrimSpace(instruction)
		if len(instruction) == 0 {
			logrus.Warnf("skipping empty instruction at position %d/%d", index, len(instructions))
			continue
		}
		instruction += ";"
		rest := len(instruction)
		total := len(instruction)
		tries := 0
		for {
			tries++
			written, err := writeFn(instruction)
			if err != nil {
				return err
			}
			rest -= written
			if rest < 0 {
				return errors.New("write overflow. Written more bytes than feeded. Can not continue")
			}
			stats := Stats{
				JobID:        f.jobID,
				Line:         index + 1,
				Total:        len(instructions),
				CurrentTry:   tries,
				CurrentTotal: total,
				CurrentRest:  rest,
			}
			if err = statFn(stats); err != nil {
				return err
			}
			if err = ctx.Err(); err != nil {
				return err
			}
			if rest == 0 {
				break
			}
			instruction = instruction[written:]
			time.Sleep(f.backoff)
			if tries >= f.maxTries {
				return errors.New("failed to send instruction after maximum tries")
			}
		}
	}
	return nil
}
