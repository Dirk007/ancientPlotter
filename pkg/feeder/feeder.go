package feeder

import (
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DefaultMaxTries = 5
	DefaultBackoff  = 100 * time.Millisecond
)

type (
	WriterFn = func(s string) (written int, err error)
	StatFn   = func(line int, total int, currentTry int, currentTotal int, currentRest int) (err error)
)

type Feeder struct {
	maxTries int
	backoff  time.Duration
}

func New(maxTries int, backoff time.Duration) *Feeder {
	return &Feeder{
		maxTries: maxTries,
		backoff:  backoff,
	}
}

func Default() *Feeder {
	return New(DefaultMaxTries, DefaultBackoff)
}

func (f *Feeder) Feed(writeFn WriterFn, instructions []string, statFn StatFn) error {
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
			if err = statFn(index+1, len(instructions), tries, total, rest); err != nil {
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
