package feeder

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Dirk007/ancientPlotter/pkg/serial"
	"github.com/sirupsen/logrus"
)

const (
	DefaultMaxTries = 5
	DefaultBackoff  = 100 * time.Millisecond
)

type Stats struct {
	JobID        string
	FatalError   error
	Line         int
	Total        int
	CurrentTry   int
	CurrentTotal int
	CurrentRest  int
}

type (
	StatFn = func(stats Stats) (err error)
)

type Feeder struct {
	maxTries int
	backoff  time.Duration
	jobID    string
	writer   serial.Writer
}

func New(maxTries int, backoff time.Duration, jobID string, writer serial.Writer) *Feeder {
	return &Feeder{
		maxTries: maxTries,
		backoff:  backoff,
		jobID:    jobID,
		writer:   writer,
	}
}

func (f *Feeder) WriteInstruction(ctx context.Context, instruction string, statFn StatFn, currentStats *Stats) error {
	rest := len(instruction)
	total := len(instruction)
	tries := 0
	stats := *currentStats
	for {
		tries++
		written, err := f.writer.Write(instruction)
		if err != nil {
			return err
		}
		rest -= written
		if rest < 0 {
			return errors.New("write overflow. Written more bytes than feeded. Can not continue")
		}

		stats.CurrentTry = tries
		stats.CurrentTotal = total
		stats.CurrentRest = rest
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
	return nil
}

func (f *Feeder) Feed(ctx context.Context, instructions []string, statFn StatFn) error {
	for index, instruction := range instructions {
		instruction := strings.TrimSpace(instruction)
		if len(instruction) == 0 {
			logrus.Warnf("skipping empty instruction at position %d/%d", index, len(instructions))
			continue
		}
		instruction += ";"

		stats := Stats{
			JobID: f.jobID,
			Line:  index + 1,
			Total: len(instructions),
		}

		if err := f.WriteInstruction(ctx, instruction, statFn, &stats); err != nil {
			return err
		}
	}
	return nil
}
