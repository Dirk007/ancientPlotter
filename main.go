package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Dirk007/ancientPlotter/pkg/jobs"
	"github.com/Dirk007/ancientPlotter/pkg/network"

	"github.com/fred1268/go-clap/clap"
)

const defaultPort = 11175

var contextDeps = jobs.NewContextDependencies()

type config struct {
	DryRun       bool    `clap:"--dry-run,d"`
	PrintOnly    bool    `clap:"--print-only,p"`
	SerialDevice *string `clap:"--serial-device,s"`
	Filename     string  `clap:"--filename,f,mandatory"`
	Serve        bool    `clap:"--serve,S"`
	Port         int     `clap:"--port,P"`
}

func spinUpLogPrinter(ctx context.Context) {
	go func() {
		idLog, logChannel := contextDeps.Logs.Register()
		idStats, statsChannel := contextDeps.Stats.Register()
		defer contextDeps.Logs.Remove(idLog)
		defer contextDeps.Stats.Remove(idStats)

		for {
			select {
			case log, ok := <-logChannel:
				if !ok {
					fmt.Println("logs channel closed")
					return
				}
				fmt.Println(log)
			case stats, ok := <-statsChannel:
				if !ok {
					fmt.Println("stats channel closed")
					return
				}
				percent := (100.0 / float64(stats.Total)) * float64(stats.Line)
				msg := fmt.Sprintf("%3d/%3d [#%d:%d/%d], %.1f%%\r", stats.Line, stats.Total, stats.CurrentTry, stats.CurrentTotal-stats.CurrentRest, stats.CurrentTotal, percent)
				fmt.Print(msg)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func errorExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func serve(_ context.Context, config *config) {
	if config.Port == 0 {
		config.Port = defaultPort
	}

	fmt.Printf("Starting web server on port %d...\n", config.Port)
	go network.Serve(contextDeps, config.Port, jobs.JobConfig{
		DryRun:       config.DryRun,
		PrintOnly:    config.PrintOnly,
		SerialDevice: config.SerialDevice,
	})
}

func alive(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				contextDeps.Alive.Broadcast(ctx, "Online")
				time.Sleep(time.Second * 1)
			}
		}
	}()
}

func log(ctx context.Context, content string) {
	contextDeps.Logs.Broadcast(ctx, content)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &config{}
	_, err := clap.Parse(os.Args, config)
	if err != nil {
		errorExit(err)
	}

	spinUpLogPrinter(ctx)
	alive(ctx)

	if config.Serve {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		serve(ctx, config)
		wg.Wait()
	}

	go func() {
		job := jobs.NewPlotJob(config.Filename)
		job.Run(ctx, contextDeps, jobs.JobConfig{
			DryRun:       config.DryRun,
			PrintOnly:    config.PrintOnly,
			SerialDevice: config.SerialDevice,
		})
	}()

	log(ctx, "Hit return *after* the plotter has completed the job.")
	var dontcare string
	_, _ = fmt.Scanln(&dontcare)

	log(ctx, "canceling...")
	cancel()
	ctx = context.Background()
	log(ctx, "done")

	time.Sleep(time.Second * 5) // Allow the last log message to be printed before exiting

	// -> Vereinfachen, Kombinieren, Vereinigung <--
	// -> Vereinigung, Vereinfachen <-- Besser

	log(ctx, "Done.")
}
