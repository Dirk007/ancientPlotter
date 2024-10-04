package jobs

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Dirk007/ancientPlotter/pkg/feeder"
	"github.com/Dirk007/ancientPlotter/pkg/normalizer"
	"github.com/Dirk007/ancientPlotter/pkg/opcodes"
	"github.com/Dirk007/ancientPlotter/pkg/serial"
	"github.com/google/uuid"
)

type JobConfig struct {
	CLIJob       bool
	DryRun       bool
	PrintOnly    bool
	SerialDevice *string
}

type PlotJob struct {
	ID     string
	Path   string
	State  JobState
	Cancel *context.CancelFunc
}

func NewPlotJob(path string) *PlotJob {
	return &PlotJob{
		ID:    uuid.New().String(),
		Path:  path,
		State: JobStateNew,
	}
}

func getWriterFor(config *JobConfig, deps *ContextDependencies) (serial.Writer, error) {
	if config.PrintOnly {
		return &serial.ConsoleWriter{}, nil
	}

	var portName string
	if config.SerialDevice == nil {
		deps.Logs.Broadcast(context.Background(), "Trying to guess serial port... ")
		guess, err := serial.GuessPortName()
		if err != nil {
			return nil, err
		}
		deps.Logs.Broadcast(context.Background(), fmt.Sprintf("Guess: %s", guess))
		portName = guess
		fmt.Println(portName)
	} else {
		portName = *config.SerialDevice
	}

	serialWriter, err := serial.TryNew(portName)
	if err != nil {
		return nil, err
	}
	return serialWriter, nil
}

func sendErrorStat(deps *ContextDependencies, err error) {
	stat := feeder.Stats{
		FatalError:   err,
		Total:        0,
		Line:         0,
		CurrentTry:   0,
		CurrentTotal: 0,
		CurrentRest:  0,
	}
	deps.Stats.Broadcast(context.Background(), stat)
}

func (j PlotJob) Run(ctx context.Context, deps *ContextDependencies, config JobConfig) error {
	defer func() {
		if config.CLIJob {
			return
		}
		os.Remove(j.Path)
	}()

	deps.Logs.Broadcast(ctx, fmt.Sprintf("Run job: %+v with config %+v", j, config))

	ext := filepath.Ext(j.Path)
	if ext != ".hpgl" {
		deps.Logs.Broadcast(ctx, fmt.Sprintf("Trying to convert %s to hpgl...", j.Path))
		cmd := exec.Command(
			"/usr/bin/inkscape",
			"--export-type=hpgl",
			j.Path,
		)
		cmd.Dir = "/"
		deps.Logs.Broadcast(ctx, fmt.Sprintf("Trying by %s...", cmd))
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		deps.Logs.Broadcast(ctx, fmt.Sprintf("Convert Err: %v", err))
		if err != nil {
			deps.Logs.Broadcast(ctx, fmt.Sprintf("Unable to convert: %v", err))
			return err
		}

		deps.Logs.Broadcast(ctx, fmt.Sprintf("Convert: StdOut %s, StdErr %s", out.String(), stderr.String()))

		j.Path = strings.Trim(j.Path, filepath.Ext(j.Path)) + ".hpgl"
		defer os.Remove(j.Path)
	}

	content, err := os.ReadFile(j.Path)
	if err != nil {
		sendErrorStat(deps, err)
		return err
	}

	writer, err := getWriterFor(&config, deps)
	if err != nil {
		sendErrorStat(deps, err)
		return err
	}

	defaultNormalizer := normalizer.Default()

	if config.DryRun {
		defaultNormalizer = defaultNormalizer.
			WithOpcodeReplacement(opcodes.OpcodePenDown, opcodes.OpcodePenUp)
	}

	normalized, err := defaultNormalizer.Normalize(string(content))
	if err != nil {
		sendErrorStat(deps, err)
		return err
	}

	defaultFeeder := feeder.New(feeder.DefaultMaxTries, feeder.DefaultBackoff, j.ID, writer)
	err = defaultFeeder.Feed(ctx, normalized, func(stat feeder.Stats) error {
		deps.Stats.Broadcast(context.Background(), stat)
		return nil
	})
	if err != nil {
		deps.Logs.Broadcast(ctx, fmt.Sprintf("Feeder error %v", err))
		sendErrorStat(deps, err)
		return err
	}

	if !config.CLIJob {
		deps.Logs.Broadcast(ctx, "Waiting for final cancel after plot is finished")
		_ = <-ctx.Done()
	}
	deps.Logs.Broadcast(context.Background(), "Finished job: "+j.ID)

	return nil
}
