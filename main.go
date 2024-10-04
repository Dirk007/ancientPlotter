package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Dirk007/ancientPlotter/pkg/jobs"
	"github.com/Dirk007/ancientPlotter/pkg/network"

	"github.com/Dirk007/clapper/pkg/clapper"
)

const defaultPort = 11175

var contextDeps = jobs.NewContextDependencies()

type config struct {
	DryRun       bool    `clapper:"long,help='Keep the cutter floating'"`
	PrintOnly    bool    `clapper:"long,short,help='Just dump the plot to the console'"`
	SerialDevice *string `clapper:"long,short,help='Serial device to use for cutting'"`
	Serve        bool    `clapper:"long,short=S,help='Start a server to receive and display plots'"`
	Port         int     `clapper:"long,short=P,default=11175,help='Port to listen on for plots'"`

	Help bool `clapper:"long,short=,help='Show this help and exit'"`
}

func spinUpLogPrinter(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
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

func alive(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
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

func showHelp() {
	fmt.Fprintln(os.Stderr, "Usage: ancient-plotter [options] <input-files>")
	s, _ := clapper.HelpDefault(new(config))
	fmt.Fprintln(os.Stderr, s)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &config{}
	filenames, err := clapper.Parse(config)
	if err != nil {
		errorExit(err)
	}

	if config.Help {
		showHelp()
		os.Exit(0)
	}

	wg := &sync.WaitGroup{}
	spinUpLogPrinter(ctx, wg)
	alive(ctx, wg)

	if config.Serve {
		wg.Add(1)
		serve(ctx, config)
		wg.Wait()
		os.Exit(0)
	}

	if len(filenames) == 0 {
		errorExit(fmt.Errorf("no input files provided"))
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		job := jobs.NewPlotJob(filenames[0])
		err = job.Run(ctx, contextDeps, jobs.JobConfig{
			CLIJob:       true,
			DryRun:       config.DryRun,
			PrintOnly:    config.PrintOnly,
			SerialDevice: config.SerialDevice,
		})
		if err != nil {
			errorExit(err)
		}
		log(ctx, "Hit return *after* the plotter has completed the job.")
	}()

	var dontcare string
	_, _ = fmt.Scanln(&dontcare)

	log(ctx, "canceling...")
	cancel()
	wg.Wait()
	ctx = context.Background()
	log(ctx, "done")

	log(ctx, "Done.")
}
