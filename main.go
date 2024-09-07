package main

import (
	"fmt"
	"os"

	"github.com/Dirk007/ancientPlotter/pkg/feeder"
	"github.com/Dirk007/ancientPlotter/pkg/normalizer"
	"github.com/Dirk007/ancientPlotter/pkg/opcodes"
	"github.com/Dirk007/ancientPlotter/pkg/serial"

	"github.com/fred1268/go-clap/clap"
)

type config struct {
	DryRun       bool    `clap:"--dry-run,d"`
	PrintOnly    bool    `clap:"--print-only,p"`
	SerialDevice *string `clap:"--serial-device,s"`
	Filename     string  `clap:"--filename,f,mandatory"`
}

func getWriterFor(config *config) (feeder.WriterFn, error) {
	if config.PrintOnly {
		return serial.PrintConsole, nil
	}

	var portName string
	if config.SerialDevice == nil {
		fmt.Print("Trying to guess serial port... ")
		guess, err := serial.GuessPortName()
		if err != nil {
			return nil, err
		}
		portName = guess
		fmt.Println(portName)
	} else {
		portName = *config.SerialDevice
	}

	serialWriter, err := serial.TryNew(portName)
	if err != nil {
		return nil, err
	}
	return serialWriter.Write, nil
}

func printStats(line int, total int, currentTry int, currentTotal int, currentRest int) error {
	percent := (100.0 / float64(total)) * float64(line)
	s := fmt.Sprintf("%3d/%3d [#%d:%d/%d], %.1f%%\r", line, total, currentTry, currentTotal-currentRest, currentTotal, percent)
	fmt.Print(s)
	return nil
}

func errorExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	config := &config{}
	_, err := clap.Parse(os.Args, config)
	if err != nil {
		errorExit(err)
	}

	writeFn, err := getWriterFor(config)
	if err != nil {
		errorExit(err)
	}

	content, err := os.ReadFile(config.Filename)
	if err != nil {
		errorExit(err)
	}

	defaultNormalizer := normalizer.Default()

	if config.DryRun {
		defaultNormalizer = defaultNormalizer.
			WithOpcodeReplacement(opcodes.OpcodePenDown, opcodes.OpcodePenUp)
	}

	normalized, err := defaultNormalizer.Normalize(string(content))
	if err != nil {
		errorExit(err)
	}

	defaultFeeder := feeder.Default()
	err = defaultFeeder.Feed(writeFn, normalized, printStats)
	if err != nil {
		errorExit(err)
	}

	fmt.Println("\r\nHit return *after* the plotter has completed the job.")
	var dontcare string
	_, _ = fmt.Scanln(&dontcare)

	// -> Vereinfachen, Kombinieren, Vereinigung <--
	// -> Vereinigung, Vereinfachen <-- Besser

	fmt.Println("\r\nDone.")
}
